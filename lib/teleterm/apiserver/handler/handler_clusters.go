// Copyright 2021 Gravitational, Inc
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

package handler

import (
	"context"

	"github.com/gravitational/trace"

	api "github.com/gravitational/teleport/gen/proto/go/teleport/lib/teleterm/v1"
	"github.com/gravitational/teleport/lib/teleterm/clusters"
)

// ListRootClusters lists root clusters
func (s *Handler) ListRootClusters(ctx context.Context, r *api.ListClustersRequest) (*api.ListClustersResponse, error) {
	clusters, err := s.DaemonService.ListRootClusters(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	result := []*api.Cluster{}
	for _, cluster := range clusters {
		result = append(result, newAPIRootCluster(cluster))
	}

	return &api.ListClustersResponse{
		Clusters: result,
	}, nil
}

// ListLeafClusters lists leaf clusters
func (s *Handler) ListLeafClusters(ctx context.Context, req *api.ListLeafClustersRequest) (*api.ListClustersResponse, error) {
	leaves, err := s.DaemonService.ListLeafClusters(ctx, req.ClusterUri)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	response := &api.ListClustersResponse{}
	for _, leaf := range leaves {
		response.Clusters = append(response.Clusters, newAPILeafCluster(leaf))
	}

	return response, nil
}

// AddCluster creates a new cluster
func (s *Handler) AddCluster(ctx context.Context, req *api.AddClusterRequest) (*api.Cluster, error) {
	cluster, err := s.DaemonService.AddCluster(ctx, req.Name)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return newAPIRootCluster(cluster), nil
}

// RemoveCluster removes a cluster from local system
func (s *Handler) RemoveCluster(ctx context.Context, req *api.RemoveClusterRequest) (*api.EmptyResponse, error) {
	if err := s.DaemonService.RemoveCluster(ctx, req.ClusterUri); err != nil {
		return nil, trace.Wrap(err)
	}

	return &api.EmptyResponse{}, nil
}

// GetCluster returns a cluster
func (s *Handler) GetCluster(ctx context.Context, req *api.GetClusterRequest) (*api.Cluster, error) {
	cluster, _, err := s.DaemonService.ResolveClusterWithDetails(ctx, req.ClusterUri)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	apiRootClusterWithDetails, err := newAPIRootClusterWithDetails(cluster)

	return apiRootClusterWithDetails, trace.Wrap(err)
}

func newAPIRootCluster(cluster *clusters.Cluster) *api.Cluster {
	loggedInUser := cluster.GetLoggedInUser()

	apiCluster := &api.Cluster{
		Uri:       cluster.URI.String(),
		Name:      cluster.Name,
		ProxyHost: cluster.GetProxyHost(),
		Connected: cluster.Connected(),
		LoggedInUser: &api.LoggedInUser{
			Name:           loggedInUser.Name,
			SshLogins:      loggedInUser.SSHLogins,
			Roles:          loggedInUser.Roles,
			ActiveRequests: loggedInUser.ActiveRequests,
		},
	}

	return apiCluster
}

func newAPIRootClusterWithDetails(cluster *clusters.ClusterWithDetails) (*api.Cluster, error) {
	apiCluster := newAPIRootCluster(cluster.Cluster)

	apiCluster.Features = &api.Features{
		AdvancedAccessWorkflows: cluster.Features.GetAdvancedAccessWorkflows(),
		IsUsageBasedBilling:     cluster.Features.GetIsUsageBased(),
	}
	apiCluster.LoggedInUser.RequestableRoles = cluster.RequestableRoles
	apiCluster.LoggedInUser.SuggestedReviewers = cluster.SuggestedReviewers
	apiCluster.AuthClusterId = cluster.AuthClusterID
	apiCluster.LoggedInUser.Acl = cluster.ACL
	userType, err := clusters.UserTypeFromString(cluster.UserType)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	apiCluster.LoggedInUser.UserType = userType
	apiCluster.ProxyVersion = cluster.ProxyVersion

	return apiCluster, nil
}

func newAPILeafCluster(leaf clusters.LeafCluster) *api.Cluster {
	return &api.Cluster{
		Name:      leaf.Name,
		Uri:       leaf.URI.String(),
		Connected: leaf.Connected,
		Leaf:      true,
		LoggedInUser: &api.LoggedInUser{
			Name:      leaf.LoggedInUser.Name,
			SshLogins: leaf.LoggedInUser.SSHLogins,
			Roles:     leaf.LoggedInUser.Roles,
		},
	}
}
