// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connectmycomputer

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"

	"github.com/gravitational/teleport"
	apidefaults "github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	apiutils "github.com/gravitational/teleport/api/utils"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/teleterm/clusters"
	"github.com/gravitational/teleport/lib/utils"
)

type RoleSetup struct {
	cfg *RoleSetupConfig
}

func NewRoleSetup(cfg *RoleSetupConfig) (*RoleSetup, error) {
	err := cfg.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	return &RoleSetup{cfg: cfg}, nil
}

type RoleSetupResult struct {
	CertsReloaded bool
}

// Run ensures that at the end of its execution the user has their own individual Connect My
// Computer role and that the role includes the current system username in allowed logins.
//
// If the role list of the user got updated, the return value has CertsReloaded set to true.
func (s *RoleSetup) Run(ctx context.Context, accessAndIdentity AccessAndIdentity, certManager CertManager, cluster *clusters.Cluster) (RoleSetupResult, error) {
	noCertsReloaded := RoleSetupResult{}
	if !cluster.URI.IsRoot() {
		return noCertsReloaded, trace.BadParameter("Connect My Computer works only with root clusters")
	}

	// Do not use GetCurrentUser – it returns the current view of the user given the certs, not merely
	// the resource from the backend. This means that GetCurrentUser includes any roles granted
	// through access requests. We don't want that since we're later going to update the user.
	clusterUser, err := accessAndIdentity.GetUser(ctx, cluster.GetLoggedInUser().Name, false /* withSecrets */)
	if err != nil {
		return noCertsReloaded, trace.Wrap(err)
	}

	userType := clusterUser.GetUserType()
	if userType != types.UserTypeLocal {
		return noCertsReloaded,
			trace.BadParameter("Connect My Computer works only with local users, the user %v was created by %v connector",
				clusterUser.GetName(), clusterUser.GetCreatedBy().Connector.Type)
	}

	roleName := fmt.Sprintf("%v%v", teleport.ConnectMyComputerRoleNamePrefix, clusterUser.GetName())

	doesRoleExist := true
	existingRole, err := accessAndIdentity.GetRole(ctx, roleName)
	if err != nil {
		if trace.IsNotFound(err) {
			doesRoleExist = false
		} else {
			return noCertsReloaded, trace.Wrap(err)
		}
	}

	systemUser, err := user.Current()
	if err != nil {
		return noCertsReloaded, trace.Wrap(err)
	}

	if !doesRoleExist {
		s.cfg.Log.Infof("Creating the role %v.", roleName)

		role, err := types.NewRole(roleName, types.RoleSpecV6{
			Allow: types.RoleConditions{
				NodeLabels: types.Labels{
					types.ConnectMyComputerNodeOwnerLabel: []string{clusterUser.GetName()},
				},
				Logins: []string{systemUser.Username},
			},
		})
		if err != nil {
			return noCertsReloaded, trace.Wrap(err)
		}
		if role, err = accessAndIdentity.UpsertRole(ctx, role); err != nil {
			return noCertsReloaded, trace.Wrap(err, "creating role %v", role.GetName())
		}
	} else {
		s.cfg.Log.Infof("The role %v already exists", roleName)
		isRoleDirty := false

		// Ensure that the current system username is in the role.
		//
		// This is to account for a use case where the same cluster user attempts to set up the role on
		// two different devices, potentially using two different system usernames. Since the role is
		// scoped per cluster user, it must include both system username
		allowedLogins := existingRole.GetLogins(types.Allow)

		if !slices.Contains(allowedLogins, systemUser.Username) {
			s.cfg.Log.Infof("Adding %v to the logins of the role %v.", systemUser.Username, roleName)

			existingRole.SetLogins(types.Allow, append(allowedLogins, systemUser.Username))
			isRoleDirty = true
		}

		// Ensure that the owner label has the expected value.
		//
		// This can happen only if someone has manually edited the role. Ensuring it has the expected
		// value will make sure that the user is able to connect to relevant nodes. This is done more to
		// reduce the support load than to make the feature more secure.
		allowedNodeLabels := existingRole.GetNodeLabels(types.Allow)
		ownerNodeLabelValue := allowedNodeLabels[types.ConnectMyComputerNodeOwnerLabel]
		expectedOwnerNodeLabelValue := []string{clusterUser.GetName()}

		if !slices.Equal(ownerNodeLabelValue, expectedOwnerNodeLabelValue) {
			s.cfg.Log.Infof("Overwriting the owner node label in the role %v.", roleName)

			allowedNodeLabels[types.ConnectMyComputerNodeOwnerLabel] = expectedOwnerNodeLabelValue
			isRoleDirty = true
		}

		if isRoleDirty {
			timeoutCtx, cancel := context.WithTimeout(ctx, resourceUpdateTimeout)
			defer cancel()
			err = s.syncResourceUpdate(timeoutCtx, accessAndIdentity, existingRole, func(ctx context.Context) error {
				existingRole, err := accessAndIdentity.UpsertRole(ctx, existingRole)
				return trace.Wrap(err, "updating role %v", existingRole.GetName())
			})
			if err != nil {
				return noCertsReloaded, trace.Wrap(err)
			}
		}
	}

	hasCMCRole := slices.Contains(clusterUser.GetRoles(), roleName)

	if hasCMCRole {
		s.cfg.Log.Infof("The user %v already has the role %v.", clusterUser.GetName(), roleName)
		return noCertsReloaded, nil
	}

	s.cfg.Log.Infof("Adding the role %v to the user %v.", roleName, clusterUser.GetName())
	clusterUser.AddRole(roleName)
	timeoutCtx, cancel := context.WithTimeout(ctx, resourceUpdateTimeout)
	defer cancel()
	err = s.syncResourceUpdate(timeoutCtx, accessAndIdentity, clusterUser, func(ctx context.Context) error {
		_, err := accessAndIdentity.UpdateUser(ctx, clusterUser)
		return trace.Wrap(err, "updating user %v", clusterUser.GetName())
	})
	if err != nil {
		return noCertsReloaded, trace.Wrap(err)
	}

	s.cfg.Log.Info("Reissuing certs.")
	// ReissueUserCerts called with CertCacheDrop and a bogus access request ID in DropAccessRequests
	// allows us to refresh the role list in the certs without forcing the user to relogin.
	//
	// Sending bogus request IDs is not documented but it is covered by tests. Refreshing roles based
	// on the server state is necessary for tsh request drop to work.
	//
	// If passing bogus request IDs ever needs to be removed, then there are two options here:
	// * Pass a wildcard instead. This will break setups where people use access requests to make
	//   Connect My Computer work. Most users will probably not use access requests for that though.
	// * Invalidate the stored certs somehow to force the user to relogin. If Connect makes a request
	//   after role setup and [client.IsErrorResolvableWithRelogin] returns true for the error from
	//   the response, Connect will ask the user to relogin.
	//
	// TODO(ravicious): Expand auth.ServerWithRoles.GenerateUserCerts to support refreshing role
	// list without having to send a bogus request ID.
	err = certManager.ReissueUserCerts(ctx, client.CertCacheDrop, client.ReissueParams{
		RouteToCluster:     cluster.Name,
		DropAccessRequests: []string{fmt.Sprintf("bogus-request-id-%v", uuid.NewString())},
	})
	return RoleSetupResult{CertsReloaded: true}, trace.Wrap(err)
}

const resourceUpdateTimeout = 15 * time.Second

// syncResourceUpdate calls a function which updates the given resource and then waits until the
// cache propagates the change.
func (s *RoleSetup) syncResourceUpdate(ctx context.Context, accessAndIdentity AccessAndIdentity, resource types.Resource, updateFunc func(context.Context) error) error {
	watcher, err := initializeWatcher(ctx, accessAndIdentity, resource.GetKind())
	if err != nil {
		return trace.Wrap(err)
	}
	defer watcher.Close()

	err = updateFunc(ctx)
	if err != nil {
		return trace.Wrap(err, "calling update function")
	}

	_, err = waitForOpPut(ctx, watcher, resource.GetKind(), resource.GetName())
	return trace.Wrap(err)
}

// AccessAndIdentity represents services.Access, services.Identity, services.Presence and auth.Cache
// methods used by [RoleSetup]. During a normal operation, auth.ClientI is passed as this interface.
type AccessAndIdentity interface {
	// See services.Access.GetRole.
	GetRole(ctx context.Context, name string) (types.Role, error)
	// See services.Access.UpsertRole.
	UpsertRole(context.Context, types.Role) (types.Role, error)
	// See auth.Cache.NewWatcher.
	NewWatcher(ctx context.Context, watch types.Watch) (types.Watcher, error)

	// See services.Identity.GetUser.
	GetUser(ctx context.Context, name string, withSecrets bool) (types.User, error)
	// See services.Identity.UpdateUser.
	UpdateUser(context.Context, types.User) (types.User, error)

	// See services.Presence.GetNode.
	GetNode(ctx context.Context, namespace, name string) (types.Server, error)
}

// CertManager enables the usage of only select methods from [client.ProxyClient] so that there
// is no need to mock the whole ProxyClient interface in tests.
type CertManager interface {
	// See [client.ProxyClient.ReissueUserCerts].
	ReissueUserCerts(context.Context, client.CertCachePolicy, client.ReissueParams) error
}

type RoleSetupConfig struct {
	Log *logrus.Entry
}

func (c *RoleSetupConfig) CheckAndSetDefaults() error {
	if c.Log == nil {
		c.Log = logrus.NewEntry(logrus.StandardLogger()).WithField(trace.Component, "CMC role")
	}

	return nil
}

type TokenProvisioner struct {
	cfg *TokenProvisionerConfig
}

func NewTokenProvisioner(cfg *TokenProvisionerConfig) *TokenProvisioner {
	cfg.checkAndSetDefaults()

	return &TokenProvisioner{cfg: cfg}
}

// CreateNodeToken creates a node join token that is valid for 5 minutes.
func (t *TokenProvisioner) CreateNodeToken(ctx context.Context, provisioner Provisioner, cluster *clusters.Cluster) (*NodeToken, error) {
	tokenName, err := utils.CryptoRandomHex(auth.TokenLenBytes)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	var req types.ProvisionTokenSpecV2
	req.Roles = types.SystemRoles{types.RoleNode}
	expires := t.cfg.Clock.Now().UTC().Add(5 * time.Minute)

	provisionToken, err := types.NewProvisionTokenFromSpec(tokenName, expires, req)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	err = provisioner.CreateToken(ctx, provisionToken)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &NodeToken{
		Token: tokenName,
		Labels: types.Labels{
			types.ConnectMyComputerNodeOwnerLabel: apiutils.Strings{cluster.GetLoggedInUser().Name},
		},
	}, nil
}

// DeleteToken deletes a join token
func (t *TokenProvisioner) DeleteToken(ctx context.Context, provisioner Provisioner, token string) error {
	err := provisioner.DeleteToken(ctx, token)
	return trace.Wrap(err)
}

type TokenProvisionerConfig struct {
	Clock clockwork.Clock
}

type NodeToken struct {
	Token  string
	Labels types.Labels
}

func (c *TokenProvisionerConfig) checkAndSetDefaults() {
	if c.Clock == nil {
		c.Clock = clockwork.NewRealClock()
	}
}

// Provisioner represents services.Provisioner methods used by TokenProvisioner.
// During a normal operation, auth.ClientI is passed as this interface.
type Provisioner interface {
	// See services.Provisioner.CreateToken.
	CreateToken(ctx context.Context, token types.ProvisionToken) error
	// See services.Provisioner.DeleteToken.
	DeleteToken(ctx context.Context, token string) error
}

type NodeJoinWait struct {
	cfg *NodeJoinWaitConfig
}

func NewNodeJoinWait(cfg *NodeJoinWaitConfig) (*NodeJoinWait, error) {
	err := cfg.CheckAndSetDefaults()
	if err != nil {
		return nil, err
	}

	return &NodeJoinWait{cfg: cfg}, nil
}

// Run grabs the host UUID of an agent from disk and then waits for the node with the given name to
// show up in the cluster.
//
// The Electron app calls this method soon after starting the agent process.
func (n *NodeJoinWait) Run(ctx context.Context, accessAndIdentity AccessAndIdentity, cluster *clusters.Cluster) (clusters.Server, error) {
	nodeName, err := n.getNodeNameFromHostUUIDFile(ctx, cluster)
	if err != nil {
		return clusters.Server{}, err
	}

	server, err := n.waitForNode(ctx, accessAndIdentity, cluster, nodeName)
	if err != nil {
		return clusters.Server{}, trace.Wrap(err)
	}

	// The default config generated by `teleport node config` during the setup of Connect My Computer
	// includes a command label with a hostname. Immediately after the node joins the cluster, the
	// label is most likely going to be empty. It takes a couple of seconds for it to update with the
	// actual hostname of the device.
	//
	// To work around that, we fill it out with os.Hostname if it's empty.
	err = n.fillOutHostnameLabelIfBlank(server)
	if err != nil {
		return clusters.Server{}, trace.Wrap(err)
	}

	return clusters.Server{
		URI:    cluster.URI.AppendServer(server.GetName()),
		Server: server,
	}, nil
}

func (n *NodeJoinWait) getNodeNameFromHostUUIDFile(ctx context.Context, cluster *clusters.Cluster) (string, error) {
	dataDir := getAgentDataDir(n.cfg.AgentsDir, cluster.ProfileName)

	// NodeJoinWait gets executed when the agent is booting up, so the host UUID file might not exist
	// on disk yet. Use a ticker to periodically check for its existence.
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// We're reading the host UUID file by ourselves rather than using utils.ReadHostUUID, because
			// that function returns NotFound both when the file doesn't exist and when the host UUID in
			// the file is empty.
			//
			// Here we need to be able to distinguish between both of those two cases.
			out, err := utils.ReadPath(utils.GetHostUUIDPath(dataDir))
			if err != nil {
				if trace.IsNotFound(err) {
					continue
				}
				return "", trace.Wrap(err)
			}

			id := strings.TrimSpace(string(out))
			if id == "" {
				return "", trace.NotFound("host UUID is empty")
			}

			return id, nil
		case <-ctx.Done():
			return "", trace.Wrap(ctx.Err(), "waiting for host UUID file to be created")
		}
	}
}

func (n *NodeJoinWait) waitForNode(ctx context.Context, accessAndIdentity AccessAndIdentity,
	cluster *clusters.Cluster, nodeName string) (types.Server, error) {
	watcher, err := initializeWatcher(ctx, accessAndIdentity, types.KindNode)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	defer watcher.Close()

	// Attempt to fetch the node from the cluster manually, in case it's joined the cluster before we
	// started the watcher.
	//
	// This means that we might return immediately if the node is still in the cache, even if
	// technically the agent has not joined the cluster yet. We're fine with this edge case.
	server, err := accessAndIdentity.GetNode(ctx, apidefaults.Namespace, nodeName)
	if err != nil {
		if !trace.IsNotFound(err) {
			return nil, trace.Wrap(err)
		}
		// Continue in case of NotFound error.
	} else {
		return server, nil
	}

	resource, err := waitForOpPut(ctx, watcher, types.KindNode, nodeName)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	server, ok := resource.(*types.ServerV2)
	if !ok {
		return nil, trace.Errorf("cannot cast event resource to server")
	}

	return server, nil
}

func (n *NodeJoinWait) fillOutHostnameLabelIfBlank(server types.Server) error {
	labels := server.GetCmdLabels()
	hostnameLabel, ok := labels[defaults.HostnameLabel]
	if ok && hostnameLabel.GetResult() == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return trace.Wrap(err)
		}

		hostnameLabel.SetResult(hostname)
		server.SetCmdLabels(labels)
	}

	return nil
}

type NodeJoinWaitConfig struct {
	// AgentsDir contains agent config files and data directories for Connect My Computer.
	AgentsDir string
}

func (c *NodeJoinWaitConfig) CheckAndSetDefaults() error {
	if c.AgentsDir == "" {
		return trace.BadParameter("missing agents dir")
	}

	return nil
}

// initializeWatcher creates a new watcher and waits for OpInit. The caller must remember to close
// the watcher.
func initializeWatcher(ctx context.Context, accessAndIdentity AccessAndIdentity, kind string) (types.Watcher, error) {
	watcher, err := accessAndIdentity.NewWatcher(ctx, types.Watch{
		Kinds: []types.WatchKind{
			{Kind: kind},
		},
	})
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// Wait for OpInit.
	select {
	case <-ctx.Done():
		watcher.Close()
		return nil, trace.Wrap(ctx.Err(), "waiting for OpInit event")
	case <-watcher.Done():
		return nil, trace.Wrap(watcher.Error(), "waiting for OpInit event")
	case event := <-watcher.Events():
		if event.Type != types.OpInit {
			watcher.Close()
			return nil, trace.Errorf("unexpected event type %q received from %s watcher", event.Type, kind)
		}
	}

	return watcher, nil
}

// waitForOpPut blocks until the watcher receives an OpPut event with a resource watching the given
// kind and name.
func waitForOpPut(ctx context.Context, watcher types.Watcher, kind string, name string) (types.Resource, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, trace.Wrap(ctx.Err(), "waiting for OpPut event for %v", name)
		case <-watcher.Done():
			return nil, trace.Wrap(watcher.Error(), "waiting for OpPut event for %v", name)
		case event := <-watcher.Events():
			if event.Type != types.OpPut {
				continue
			}

			// Kind + name combo is enough to uniquely identify a resource within a single cluster.
			if event.Resource.GetKind() == kind && event.Resource.GetName() == name {
				return event.Resource, nil
			}
		}
	}
}

type NodeDelete struct {
	cfg *NodeDeleteConfig
}

// Run grabs the host UUID of an agent from a disk and deletes the node with that name.
func (n *NodeDelete) Run(ctx context.Context, presence Presence, cluster *clusters.Cluster) error {
	hostUUID, err := utils.ReadHostUUID(getAgentDataDir(n.cfg.AgentsDir, cluster.ProfileName))
	if trace.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return trace.Wrap(err)
	}
	err = presence.DeleteNode(ctx, apidefaults.Namespace, hostUUID)
	if trace.IsNotFound(err) {
		return nil
	}
	return trace.Wrap(err)
}

func NewNodeDelete(cfg *NodeDeleteConfig) (*NodeDelete, error) {
	err := cfg.checkAndSetDefaults()
	if err != nil {
		return nil, err
	}

	return &NodeDelete{cfg: cfg}, nil
}

type NodeDeleteConfig struct {
	// AgentsDir contains agent config files and data directories for Connect My Computer.
	AgentsDir string
}

func (n *NodeDeleteConfig) checkAndSetDefaults() error {
	if n.AgentsDir == "" {
		return trace.BadParameter("missing agents dir")
	}

	return nil
}

// Presence represents services.Presence methods used by [NodeDelete].
// During a normal operation, auth.ClientI is passed as this interface.
type Presence interface {
	// See services.Presence.GetNode.
	DeleteNode(ctx context.Context, namespace, name string) error
}

type NodeName struct {
	cfg *NodeNameConfig
}

// Get returns the host UUID of the agent from a disk.
func (n *NodeName) Get(cluster *clusters.Cluster) (string, error) {
	hostUUID, err := utils.ReadHostUUID(getAgentDataDir(n.cfg.AgentsDir, cluster.ProfileName))
	return hostUUID, trace.Wrap(err)
}

func NewNodeName(cfg *NodeNameConfig) (*NodeName, error) {
	err := cfg.checkAndSetDefaults()
	if err != nil {
		return nil, err
	}

	return &NodeName{cfg: cfg}, nil
}

type NodeNameConfig struct {
	// AgentsDir contains agent config files and data directories for Connect My Computer.
	AgentsDir string
}

func (n *NodeNameConfig) checkAndSetDefaults() error {
	if n.AgentsDir == "" {
		return trace.BadParameter("missing agents dir")
	}

	return nil
}

func getAgentDataDir(agentsDir, profileName string) string {
	return filepath.Join(agentsDir, profileName, "data")
}
