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

package integration

import (
	"context"
	"fmt"
	"net"
	"os/user"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gravitational/teleport/api/defaults"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/api/utils"
	api "github.com/gravitational/teleport/gen/proto/go/teleport/lib/teleterm/v1"
	dbhelpers "github.com/gravitational/teleport/integration/db"
	"github.com/gravitational/teleport/integration/helpers"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/client"
	"github.com/gravitational/teleport/lib/service"
	"github.com/gravitational/teleport/lib/service/servicecfg"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/teleterm/api/uri"
	"github.com/gravitational/teleport/lib/teleterm/apiserver/handler"
	"github.com/gravitational/teleport/lib/teleterm/clusters"
	"github.com/gravitational/teleport/lib/teleterm/daemon"
	libutils "github.com/gravitational/teleport/lib/utils"
)

func TestTeleterm(t *testing.T) {
	pack := dbhelpers.SetupDatabaseTest(t,
		dbhelpers.WithListenerSetupDatabaseTest(helpers.SingleProxyPortSetup),
		dbhelpers.WithLeafConfig(func(config *servicecfg.Config) {
			config.Auth.NetworkingConfig.SetProxyListenerMode(types.ProxyListenerMode_Multiplex)
		}),
		dbhelpers.WithRootConfig(func(config *servicecfg.Config) {
			config.Auth.NetworkingConfig.SetProxyListenerMode(types.ProxyListenerMode_Multiplex)
		}),
	)
	pack.WaitForLeaf(t)

	creds, err := helpers.GenerateUserCreds(helpers.UserCredsRequest{
		Process:  pack.Root.Cluster.Process,
		Username: pack.Root.User.GetName(),
	})
	require.NoError(t, err)

	t.Run("adding root cluster", func(t *testing.T) {
		t.Parallel()

		testAddingRootCluster(t, pack, creds)
	})

	t.Run("ListRootClusters returns logged in user", func(t *testing.T) {
		t.Parallel()

		testListRootClustersReturnsLoggedInUser(t, pack, creds)
	})

	t.Run("GetCluster returns properties from auth server", func(t *testing.T) {
		t.Parallel()

		testGetClusterReturnsPropertiesFromAuthServer(t, pack)
	})

	t.Run("Test headless watcher", func(t *testing.T) {
		t.Parallel()

		testHeadlessWatcher(t, pack, creds)
	})

	t.Run("CreateConnectMyComputerRole", func(t *testing.T) {
		t.Parallel()
		testCreateConnectMyComputerRole(t, pack)
	})

	t.Run("CreateAndDeleteConnectMyComputerToken", func(t *testing.T) {
		t.Parallel()
		testCreatingAndDeletingConnectMyComputerToken(t, pack)
	})

	t.Run("WaitForConnectMyComputerNodeJoin", func(t *testing.T) {
		t.Parallel()
		testWaitForConnectMyComputerNodeJoin(t, pack, creds)
	})

	t.Run("DeleteConnectMyComputerNode", func(t *testing.T) {
		t.Parallel()
		testDeleteConnectMyComputerNode(t, pack)
	})
}

func testAddingRootCluster(t *testing.T, pack *dbhelpers.DatabasePack, creds *helpers.UserCreds) {
	t.Helper()

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                t.TempDir(),
		InsecureSkipVerify: true,
	})
	require.NoError(t, err)

	daemonService, err := daemon.New(daemon.Config{
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      t.TempDir(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	addedCluster, err := daemonService.AddCluster(context.Background(), pack.Root.Cluster.Web)
	require.NoError(t, err)

	clusters, err := daemonService.ListRootClusters(context.Background())
	require.NoError(t, err)

	clusterURIs := make([]uri.ResourceURI, 0, len(clusters))
	for _, cluster := range clusters {
		clusterURIs = append(clusterURIs, cluster.URI)
	}
	require.ElementsMatch(t, clusterURIs, []uri.ResourceURI{addedCluster.URI})
}

func testListRootClustersReturnsLoggedInUser(t *testing.T, pack *dbhelpers.DatabasePack, creds *helpers.UserCreds) {
	tc := mustLogin(t, pack.Root.User.GetName(), pack, creds)

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
	})
	require.NoError(t, err)

	daemonService, err := daemon.New(daemon.Config{
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      t.TempDir(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	handler, err := handler.New(
		handler.Config{
			DaemonService: daemonService,
		},
	)
	require.NoError(t, err)

	response, err := handler.ListRootClusters(context.Background(), &api.ListClustersRequest{})
	require.NoError(t, err)

	require.Equal(t, 1, len(response.Clusters))
	require.Equal(t, pack.Root.User.GetName(), response.Clusters[0].LoggedInUser.Name)
}

func testGetClusterReturnsPropertiesFromAuthServer(t *testing.T, pack *dbhelpers.DatabasePack) {
	authServer := pack.Root.Cluster.Process.GetAuthServer()

	// Use random names to not collide with other tests.
	uuid := uuid.NewString()
	suggestedReviewer := "suggested-reviewer"
	requestableRoleName := fmt.Sprintf("%s-%s", "requested-role", uuid)
	userName := fmt.Sprintf("%s-%s", "user", uuid)
	roleName := fmt.Sprintf("%s-%s", "get-cluster-role", uuid)

	requestableRole, err := types.NewRole(requestableRoleName, types.RoleSpecV6{})
	require.NoError(t, err)

	// Create user role with ability to request role
	userRole, err := types.NewRole(roleName, types.RoleSpecV6{
		Options: types.RoleOptions{},
		Allow: types.RoleConditions{
			Logins: []string{
				userName,
			},
			NodeLabels: types.Labels{types.Wildcard: []string{types.Wildcard}},
			Request: &types.AccessRequestConditions{
				Roles:              []string{requestableRoleName},
				SuggestedReviewers: []string{suggestedReviewer},
			},
		},
	})
	require.NoError(t, err)

	// add role that user can request
	_, err = authServer.UpsertRole(context.Background(), requestableRole)
	require.NoError(t, err)

	// add role that allows to request "requestableRole"
	_, err = authServer.UpsertRole(context.Background(), userRole)
	require.NoError(t, err)

	user, err := types.NewUser(userName)
	user.AddRole(userRole.GetName())
	require.NoError(t, err)

	_, err = authServer.UpsertUser(context.Background(), user)
	require.NoError(t, err)

	creds, err := helpers.GenerateUserCreds(helpers.UserCredsRequest{
		Process:  pack.Root.Cluster.Process,
		Username: userName,
	})
	require.NoError(t, err)

	tc := mustLogin(t, userName, pack, creds)

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
	})
	require.NoError(t, err)

	daemonService, err := daemon.New(daemon.Config{
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      t.TempDir(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	handler, err := handler.New(
		handler.Config{
			DaemonService: daemonService,
		},
	)
	require.NoError(t, err)

	rootClusterName, _, err := net.SplitHostPort(pack.Root.Cluster.Web)
	require.NoError(t, err)

	response, err := handler.GetCluster(context.Background(), &api.GetClusterRequest{
		ClusterUri: uri.NewClusterURI(rootClusterName).String(),
	})
	require.NoError(t, err)

	require.Equal(t, userName, response.LoggedInUser.Name)
	require.ElementsMatch(t, []string{requestableRoleName}, response.LoggedInUser.RequestableRoles)
	require.ElementsMatch(t, []string{suggestedReviewer}, response.LoggedInUser.SuggestedReviewers)
}

func testHeadlessWatcher(t *testing.T, pack *dbhelpers.DatabasePack, creds *helpers.UserCreds) {
	t.Helper()
	ctx := context.Background()

	tc := mustLogin(t, pack.Root.User.GetName(), pack, creds)

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
	})
	require.NoError(t, err)

	cluster, _, err := storage.Add(ctx, tc.WebProxyAddr)
	require.NoError(t, err)

	daemonService, err := daemon.New(daemon.Config{
		Storage: storage,
		CreateTshdEventsClientCredsFunc: func() (grpc.DialOption, error) {
			return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
		},
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      t.TempDir(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	expires := pack.Root.Cluster.Config.Clock.Now().Add(time.Minute)
	ha, err := types.NewHeadlessAuthentication(pack.Root.User.GetName(), "uuid", expires)
	require.NoError(t, err)
	ha.State = types.HeadlessAuthenticationState_HEADLESS_AUTHENTICATION_STATE_PENDING

	// Start the tshd event service and connect the daemon to it.

	tshdEventsService, addr := newMockTSHDEventsServiceServer(t)
	err = daemonService.UpdateAndDialTshdEventsServerAddress(addr)
	require.NoError(t, err)

	// Stop and restart the watcher twice to simulate logout + login + relogin. Ensure the watcher catches events.

	err = daemonService.StopHeadlessWatcher(cluster.URI.String())
	require.NoError(t, err)
	err = daemonService.StartHeadlessWatcher(cluster.URI.String(), false /* waitInit */)
	require.NoError(t, err)
	err = daemonService.StartHeadlessWatcher(cluster.URI.String(), true /* waitInit */)
	require.NoError(t, err)

	// Ensure the watcher catches events and sends them to the Electron App.

	err = pack.Root.Cluster.Process.GetAuthServer().UpsertHeadlessAuthentication(ctx, ha)
	assert.NoError(t, err)

	assert.Eventually(t,
		func() bool {
			return tshdEventsService.sendPendingHeadlessAuthenticationCount.Load() == 1
		},
		10*time.Second,
		500*time.Millisecond,
		"Expected tshdEventService to receive 1 SendPendingHeadlessAuthentication message but got %v",
		tshdEventsService.sendPendingHeadlessAuthenticationCount.Load(),
	)
}

func testCreateConnectMyComputerRole(t *testing.T, pack *dbhelpers.DatabasePack) {
	systemUser, err := user.Current()
	require.NoError(t, err)

	tests := []struct {
		name               string
		userAlreadyHasRole bool
		existingRole       func(userName string) types.RoleV6
	}{
		{
			name: "role does not exist",
		},
		{
			name: "role exists and includes current system username",
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{userName},
							},
							Logins: []string{systemUser.Username},
						},
					},
				}
			},
		},
		{
			name: "role exists and does not include current system username",
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{userName},
							},
							Logins: []string{fmt.Sprintf("bogus-login-%v", uuid.NewString())},
						},
					},
				}
			},
		},
		{
			name: "role exists and has no logins",
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{userName},
							},
							Logins: []string{},
						},
					},
				}
			},
		},
		{
			name: "role exists and owner node label was changed",
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{"bogus-username"},
							},
							Logins: []string{systemUser.Username},
						},
					},
				}
			},
		},
		{
			name:               "user already has existing role that includes current system username",
			userAlreadyHasRole: true,
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{userName},
							},
							Logins: []string{systemUser.Username},
						},
					},
				}
			},
		},
		{
			name:               "user already has existing role that does not include current system username",
			userAlreadyHasRole: true,
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{userName},
							},
							Logins: []string{fmt.Sprintf("bogus-login-%v", uuid.NewString())},
						},
					},
				}
			},
		},
		{
			name:               "user already has existing role with modified owner node label",
			userAlreadyHasRole: true,
			existingRole: func(userName string) types.RoleV6 {
				return types.RoleV6{
					Spec: types.RoleSpecV6{
						Allow: types.RoleConditions{
							NodeLabels: types.Labels{
								types.ConnectMyComputerNodeOwnerLabel: []string{"bogus-username"},
							},
							Logins: []string{fmt.Sprintf("bogus-login-%v", uuid.NewString())},
						},
					},
				}
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			authServer := pack.Root.Cluster.Process.GetAuthServer()
			uuid := uuid.NewString()
			userName := fmt.Sprintf("user-cmc-%s", uuid)
			roleName := fmt.Sprintf("connect-my-computer-%v", userName)

			var existingRole *types.RoleV6

			// Prepare an existing role if present.
			if test.existingRole != nil {
				role := test.existingRole(userName)
				role.SetMetadata(types.Metadata{
					Name: roleName,
				})
				existingRole = &role
				_, err := authServer.UpsertRole(ctx, &role)
				require.NoError(t, err)
			}

			// Prepare a role with rules required to call CreateConnectMyComputerRole.
			ruleWithAllowRules, err := types.NewRole(fmt.Sprintf("cmc-allow-rules-%v", uuid),
				types.RoleSpecV6{
					Allow: types.RoleConditions{
						Rules: []types.Rule{
							types.NewRule(types.KindUser, services.RW()),
							types.NewRule(types.KindRole, services.RW()),
						},
					},
				})
			require.NoError(t, err)
			userRoles := []types.Role{ruleWithAllowRules}

			// Create a new user to avoid colliding with other tests.
			// Assign to the user the role with allow rules and the existing role if present.
			if test.userAlreadyHasRole {
				if existingRole == nil {
					t.Log("userAlreadyHasRole must be used together with existingRole")
					t.Fail()
					return
				}
				userRoles = append(userRoles, existingRole)
			}
			_, err = auth.CreateUser(ctx, authServer, userName, userRoles...)
			require.NoError(t, err)

			userPassword := uuid
			require.NoError(t, authServer.UpsertPassword(userName, []byte(userPassword)))

			// Prepare daemon.Service.
			storage, err := clusters.NewStorage(clusters.Config{
				Dir:                t.TempDir(),
				InsecureSkipVerify: true,
			})
			require.NoError(t, err)

			daemonService, err := daemon.New(daemon.Config{
				Storage:        storage,
				KubeconfigsDir: t.TempDir(),
				AgentsDir:      t.TempDir(),
			})
			require.NoError(t, err)
			t.Cleanup(func() {
				daemonService.Stop()
			})
			handler, err := handler.New(
				handler.Config{
					DaemonService: daemonService,
				},
			)
			require.NoError(t, err)

			rootClusterName, _, err := net.SplitHostPort(pack.Root.Cluster.Web)
			require.NoError(t, err)
			rootClusterURI := uri.NewClusterURI(rootClusterName).String()

			// Log in as the new user.
			// It's important to use the actual login handler rather than mustLogin. mustLogin completely
			// skips the actual login flow and saves valid certs to disk. We already had a regression that
			// was not caught by this test because the test did not trigger certain code paths because it
			// was using mustLogin as a shortcut.
			_, err = handler.AddCluster(ctx, &api.AddClusterRequest{Name: pack.Root.Cluster.Web})
			require.NoError(t, err)
			_, err = handler.Login(ctx, &api.LoginRequest{
				ClusterUri: rootClusterURI,
				Params: &api.LoginRequest_Local{
					Local: &api.LoginRequest_LocalParams{User: userName, Password: userPassword},
				},
			})
			require.NoError(t, err)

			// Call CreateConnectMyComputerRole.
			response, err := handler.CreateConnectMyComputerRole(ctx, &api.CreateConnectMyComputerRoleRequest{
				RootClusterUri: rootClusterURI,
			})
			require.NoError(t, err)

			if test.userAlreadyHasRole {
				require.False(t, response.CertsReloaded,
					"expected the handler to signal that the certs were not reloaded since the user was already assigned the role")
			} else {
				require.True(t, response.CertsReloaded,
					"expected the handler to signal that the certs were reloaded since the user was just assigned a new role")
			}

			// Verify that the role exists.
			role, err := authServer.GetRole(ctx, roleName)
			require.NoError(t, err)

			// Verify that the role grants expected privileges.
			require.Contains(t, role.GetNodeLabels(types.Allow), types.ConnectMyComputerNodeOwnerLabel)
			expectedNodeLabelValue := utils.Strings{userName}
			actualNodeLabelValue := role.GetNodeLabels(types.Allow)[types.ConnectMyComputerNodeOwnerLabel]
			require.Equal(t, expectedNodeLabelValue, actualNodeLabelValue)
			require.Contains(t, role.GetLogins(types.Allow), systemUser.Username)

			// Verify that the certs have been reloaded and that the user is assigned the role.
			//
			// GetCluster reads data from the cert. If the certs were not reloaded properly, GetCluster
			// will not return the role that's just been assigned to the user.
			clusterDetails, err := handler.GetCluster(ctx, &api.GetClusterRequest{
				ClusterUri: rootClusterURI,
			})
			require.NoError(t, err)
			require.Contains(t, clusterDetails.LoggedInUser.Roles, roleName,
				"the user certs don't include the freshly added role; the certs might have not been reloaded properly")
		})
	}
}

func testCreatingAndDeletingConnectMyComputerToken(t *testing.T, pack *dbhelpers.DatabasePack) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	authServer := pack.Root.Cluster.Process.GetAuthServer()
	uuid := uuid.NewString()
	userName := fmt.Sprintf("user-cmc-%s", uuid)

	// Prepare a role with rules required to call CreateConnectMyComputerNodeToken.
	ruleWithAllowRules, err := types.NewRole(fmt.Sprintf("cmc-allow-rules-%v", uuid),
		types.RoleSpecV6{
			Allow: types.RoleConditions{
				Rules: []types.Rule{
					types.NewRule(types.KindToken, services.RW()),
				},
			},
		})
	require.NoError(t, err)
	userRoles := []types.Role{ruleWithAllowRules}

	_, err = auth.CreateUser(ctx, authServer, userName, userRoles...)
	require.NoError(t, err)

	// Log in as the new user.
	creds, err := helpers.GenerateUserCreds(helpers.UserCredsRequest{
		Process:  pack.Root.Cluster.Process,
		Username: userName,
	})
	require.NoError(t, err)
	tc := mustLogin(t, userName, pack, creds)

	fakeClock := clockwork.NewFakeClock()

	// Prepare daemon.Service.
	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
		Clock:              fakeClock,
	})
	require.NoError(t, err)

	daemonService, err := daemon.New(daemon.Config{
		Clock:          fakeClock,
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      t.TempDir(),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})
	handler, err := handler.New(
		handler.Config{
			DaemonService: daemonService,
		},
	)
	require.NoError(t, err)

	// Call CreateConnectMyComputerNodeToken.
	rootClusterName, _, err := net.SplitHostPort(pack.Root.Cluster.Web)
	require.NoError(t, err)
	rootClusterURI := uri.NewClusterURI(rootClusterName).String()
	requestCreatedAt := fakeClock.Now()
	createdTokenResponse, err := handler.CreateConnectMyComputerNodeToken(ctx, &api.CreateConnectMyComputerNodeTokenRequest{
		RootClusterUri: rootClusterURI,
	})
	require.NoError(t, err)
	require.Equal(t, &api.Label{
		Name:  types.ConnectMyComputerNodeOwnerLabel,
		Value: userName,
	}, createdTokenResponse.GetLabels()[0])

	// Verify that token exists
	tokenFromAuthServer, err := authServer.GetToken(ctx, createdTokenResponse.GetToken())
	require.NoError(t, err)

	// Verify that the token can be used to join nodes...
	require.Equal(t, types.SystemRoles{types.RoleNode}, tokenFromAuthServer.GetRoles())
	// ...and is valid for no longer than 5 minutes.
	require.LessOrEqual(t, tokenFromAuthServer.Expiry(), requestCreatedAt.Add(5*time.Minute))

	// watcher waits for the token deletion
	watcher, err := authServer.NewWatcher(ctx, types.Watch{
		Kinds: []types.WatchKind{
			{Kind: types.KindToken},
		},
	})
	require.NoError(t, err)
	defer watcher.Close()

	select {
	case <-time.After(time.Second * 10):
		t.Fatalf("Timeout waiting for event.")
	case event := <-watcher.Events():
		if event.Type != types.OpInit {
			t.Fatalf("Unexpected event type.")
		}
		require.Equal(t, event.Type, types.OpInit)
	case <-watcher.Done():
		t.Fatal(watcher.Error())
	}

	// Call DeleteConnectMyComputerToken.
	_, err = handler.DeleteConnectMyComputerToken(ctx, &api.DeleteConnectMyComputerTokenRequest{
		RootClusterUri: rootClusterURI,
		Token:          createdTokenResponse.GetToken(),
	})
	require.NoError(t, err)

	waitForResourceToBeDeleted(t, watcher, types.KindToken, createdTokenResponse.GetToken())

	_, err = authServer.GetToken(ctx, createdTokenResponse.GetToken())

	// The token should no longer exist.
	require.True(t, trace.IsNotFound(err))
}

func testWaitForConnectMyComputerNodeJoin(t *testing.T, pack *dbhelpers.DatabasePack, creds *helpers.UserCreds) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	t.Cleanup(cancel)

	tc := mustLogin(t, pack.Root.User.GetName(), pack, creds)

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
	})
	require.NoError(t, err)

	agentsDir := t.TempDir()
	daemonService, err := daemon.New(daemon.Config{
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      agentsDir,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	handler, err := handler.New(
		handler.Config{
			DaemonService: daemonService,
		},
	)
	require.NoError(t, err)

	profileName, _, err := net.SplitHostPort(pack.Root.Cluster.Web)
	require.NoError(t, err)

	waitForNodeJoinErr := make(chan error)

	go func() {
		_, err := handler.WaitForConnectMyComputerNodeJoin(ctx, &api.WaitForConnectMyComputerNodeJoinRequest{
			RootClusterUri: uri.NewClusterURI(profileName).String(),
		})
		waitForNodeJoinErr <- err
	}()

	// Start the new node.
	nodeConfig := newNodeConfig(t, pack.Root.Cluster.Config.Auth.ListenAddr, "token", types.JoinMethodToken)
	nodeConfig.DataDir = filepath.Join(agentsDir, profileName, "data")
	nodeConfig.Log = libutils.NewLoggerForTests()
	nodeSvc, err := service.NewTeleport(nodeConfig)
	require.NoError(t, err)
	require.NoError(t, nodeSvc.Start())
	t.Cleanup(func() { require.NoError(t, nodeSvc.Close()) })

	_, err = nodeSvc.WaitForEventTimeout(10*time.Second, service.TeleportReadyEvent)
	require.NoError(t, err, "timeout waiting for node readiness")

	// Verify that WaitForConnectMyComputerNodeJoin returned with no errors.
	require.NoError(t, <-waitForNodeJoinErr)
}

func testDeleteConnectMyComputerNode(t *testing.T, pack *dbhelpers.DatabasePack) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	t.Cleanup(cancel)

	authServer := pack.Root.Cluster.Process.GetAuthServer()
	uuid := uuid.NewString()
	userName := fmt.Sprintf("user-cmc-%s", uuid)

	// Prepare a role with rules required to call DeleteConnectMyComputerNode.
	ruleWithAllowRules, err := types.NewRole(fmt.Sprintf("cmc-allow-rules-%v", uuid),
		types.RoleSpecV6{
			Allow: types.RoleConditions{
				Rules: []types.Rule{
					types.NewRule(types.KindNode, services.RW()),
				},
			},
		})
	require.NoError(t, err)
	userRoles := []types.Role{ruleWithAllowRules}

	_, err = auth.CreateUser(ctx, authServer, userName, userRoles...)
	require.NoError(t, err)

	// Log in as the new user.
	creds, err := helpers.GenerateUserCreds(helpers.UserCredsRequest{
		Process:  pack.Root.Cluster.Process,
		Username: userName,
	})
	require.NoError(t, err)
	tc := mustLogin(t, userName, pack, creds)

	storage, err := clusters.NewStorage(clusters.Config{
		Dir:                tc.KeysDir,
		InsecureSkipVerify: tc.InsecureSkipVerify,
	})
	require.NoError(t, err)

	agentsDir := t.TempDir()
	daemonService, err := daemon.New(daemon.Config{
		Storage:        storage,
		KubeconfigsDir: t.TempDir(),
		AgentsDir:      agentsDir,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		daemonService.Stop()
	})

	handler, err := handler.New(
		handler.Config{
			DaemonService: daemonService,
		},
	)
	require.NoError(t, err)

	profileName, _, err := net.SplitHostPort(pack.Root.Cluster.Web)
	require.NoError(t, err)

	// Start the new node.
	nodeConfig := newNodeConfig(t, pack.Root.Cluster.Config.Auth.ListenAddr, "token", types.JoinMethodToken)
	nodeConfig.DataDir = filepath.Join(agentsDir, profileName, "data")
	nodeConfig.Log = libutils.NewLoggerForTests()
	nodeSvc, err := service.NewTeleport(nodeConfig)
	require.NoError(t, err)
	require.NoError(t, nodeSvc.Start())
	t.Cleanup(func() { require.NoError(t, nodeSvc.Close()) })

	// waits for the node to be added
	require.Eventually(t, func() bool {
		_, err := authServer.GetNode(ctx, defaults.Namespace, nodeConfig.HostUUID)
		return err == nil
	}, time.Minute, time.Second, "waiting for node to join cluster")

	//  stop the node before attempting to remove it, to more closely resemble what's going to happen in production
	err = nodeSvc.Close()
	require.NoError(t, err)

	// test
	_, err = handler.DeleteConnectMyComputerNode(ctx, &api.DeleteConnectMyComputerNodeRequest{
		RootClusterUri: uri.NewClusterURI(profileName).String(),
	})
	require.NoError(t, err)

	// waits for the node to be deleted
	require.Eventually(t, func() bool {
		_, err := authServer.GetNode(ctx, defaults.Namespace, nodeConfig.HostUUID)
		return trace.IsNotFound(err)
	}, time.Minute, time.Second, "waiting for node to be deleted")
}

// mustLogin logs in as the given user by completely skipping the actual login flow and saving valid
// certs to disk. clusters.Storage can then be pointed to tc.KeysDir and daemon.Service can act as
// if the user was successfully logged in.
//
// This is faster than going through the actual process, but keep in mind that it might skip some
// vital steps. It should be used only for tests which don't depend on complex user setup and do not
// reissue certs or modify them in some other way.
func mustLogin(t *testing.T, userName string, pack *dbhelpers.DatabasePack, creds *helpers.UserCreds) *client.TeleportClient {
	tc, err := pack.Root.Cluster.NewClientWithCreds(helpers.ClientConfig{
		Login:   userName,
		Cluster: pack.Root.Cluster.Secrets.SiteName,
	}, *creds)
	require.NoError(t, err)
	// Save the profile yaml file to disk as NewClientWithCreds doesn't do that by itself.
	err = tc.SaveProfile(false /* makeCurrent */)
	require.NoError(t, err)
	return tc
}

type mockTSHDEventsService struct {
	*api.UnimplementedTshdEventsServiceServer
	sendPendingHeadlessAuthenticationCount atomic.Uint32
}

func newMockTSHDEventsServiceServer(t *testing.T) (service *mockTSHDEventsService, addr string) {
	tshdEventsService := &mockTSHDEventsService{}

	ls, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcServer := grpc.NewServer()
	api.RegisterTshdEventsServiceServer(grpcServer, tshdEventsService)

	serveErr := make(chan error)
	go func() {
		serveErr <- grpcServer.Serve(ls)
	}()

	t.Cleanup(func() {
		grpcServer.GracefulStop()

		// For test cases that did not send any grpc calls, test may finish
		// before grpcServer.Serve is called and grpcServer.Serve will return
		// grpc.ErrServerStopped.
		err := <-serveErr
		if err != grpc.ErrServerStopped {
			assert.NoError(t, err)
		}
	})

	return tshdEventsService, ls.Addr().String()
}

func (c *mockTSHDEventsService) SendPendingHeadlessAuthentication(context.Context, *api.SendPendingHeadlessAuthenticationRequest) (*api.SendPendingHeadlessAuthenticationResponse, error) {
	c.sendPendingHeadlessAuthenticationCount.Add(1)
	return &api.SendPendingHeadlessAuthenticationResponse{}, nil
}

func waitForResourceToBeDeleted(t *testing.T, watcher types.Watcher, kind, name string) {
	timeout := time.After(time.Second * 15)
	for {
		select {
		case <-timeout:
			t.Fatalf("Timeout waiting for event.")
		case event := <-watcher.Events():
			if event.Type != types.OpDelete {
				continue
			}
			if event.Resource.GetKind() == kind && event.Resource.GetMetadata().Name == name {
				return
			}
		case <-watcher.Done():
			t.Fatalf("Watcher error %s.", watcher.Error())
		}
	}
}
