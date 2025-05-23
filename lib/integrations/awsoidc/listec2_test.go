/*
Copyright 2023 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package awsoidc

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gravitational/trace"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/types"
)

type mockListEC2Client struct {
	pageSize     int
	accountID    string
	ec2Instances []ec2Types.Instance
}

// GetCallerIdentity returns information about the caller identity.
func (m *mockListEC2Client) GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: &m.accountID,
	}, nil
}

// Returns information about ec2 instances.
// This API supports pagination.
func (m mockListEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	requestedPage := 1

	stateFilter := false
	platformFilter := false
	for _, filter := range params.Filters {
		if aws.ToString(filter.Name) == "instance-state-name" && len(filter.Values) == 1 && filter.Values[0] == "running" {
			stateFilter = true
		}
		if aws.ToString(filter.Name) == "platform-details" && len(filter.Values) == 1 && filter.Values[0] == "Linux/UNIX" {
			platformFilter = true
		}
	}
	if !stateFilter || !platformFilter {
		return nil, trace.BadParameter("instance-state-name and platform-details filters were not included")
	}

	totalInstances := len(m.ec2Instances)

	if params.NextToken != nil {
		currentMarker, err := strconv.Atoi(*params.NextToken)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		requestedPage = currentMarker
	}

	sliceStart := m.pageSize * (requestedPage - 1)
	sliceEnd := m.pageSize * requestedPage
	if sliceEnd > totalInstances {
		sliceEnd = totalInstances
	}

	ret := &ec2.DescribeInstancesOutput{
		Reservations: []ec2Types.Reservation{{
			Instances: m.ec2Instances[sliceStart:sliceEnd],
		}},
	}

	if sliceEnd < totalInstances {
		nextToken := strconv.Itoa(requestedPage + 1)
		ret.NextToken = stringPointer(nextToken)
	}

	return ret, nil
}

func TestListEC2(t *testing.T) {
	ctx := context.Background()

	noErrorFunc := func(err error) bool {
		return err == nil
	}

	const pageSize = 100
	t.Run("pagination", func(t *testing.T) {
		totalEC2s := 203

		allInstances := make([]ec2Types.Instance, 0, totalEC2s)
		for i := 0; i < totalEC2s; i++ {
			allInstances = append(allInstances, ec2Types.Instance{
				PrivateDnsName:   aws.String("my-private-dns.compute.aws"),
				InstanceId:       aws.String(fmt.Sprintf("i-%d", i)),
				VpcId:            aws.String("vpc-abcd"),
				SubnetId:         aws.String("subnet-123"),
				PrivateIpAddress: aws.String("172.31.1.1"),
			})
		}

		mockListClient := &mockListEC2Client{
			pageSize:     pageSize,
			accountID:    "123456789012",
			ec2Instances: allInstances,
		}

		// First page must return pageSize number of Servers
		resp, err := ListEC2(ctx, mockListClient, ListEC2Request{
			Region:      "us-east-1",
			Integration: "myintegration",
			NextToken:   "",
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.NextToken)
		require.Len(t, resp.Servers, pageSize)
		nextPageToken := resp.NextToken
		require.Equal(t, resp.Servers[0].GetCloudMetadata().AWS.InstanceID, "i-0")

		// Second page must return pageSize number of Servers
		resp, err = ListEC2(ctx, mockListClient, ListEC2Request{
			Region:      "us-east-1",
			Integration: "myintegration",
			NextToken:   nextPageToken,
		})
		require.NoError(t, err)
		require.NotEmpty(t, resp.NextToken)
		require.Len(t, resp.Servers, pageSize)
		nextPageToken = resp.NextToken
		require.Equal(t, resp.Servers[0].GetCloudMetadata().AWS.InstanceID, "i-100")

		// Third page must return only the remaining Servers and an empty nextToken
		resp, err = ListEC2(ctx, mockListClient, ListEC2Request{
			Region:      "us-east-1",
			Integration: "myintegration",
			NextToken:   nextPageToken,
		})
		require.NoError(t, err)
		require.Empty(t, resp.NextToken)
		require.Len(t, resp.Servers, 3)
		require.Equal(t, resp.Servers[0].GetCloudMetadata().AWS.InstanceID, "i-200")
	})

	for _, tt := range []struct {
		name          string
		req           ListEC2Request
		mockInstances []ec2Types.Instance
		errCheck      func(error) bool
		respCheck     func(*testing.T, *ListEC2Response)
	}{
		{
			name: "valid for listing instances",
			req: ListEC2Request{
				Region:      "us-east-1",
				Integration: "myintegration",
				NextToken:   "",
			},
			mockInstances: []ec2Types.Instance{{
				PrivateDnsName:   aws.String("my-private-dns.compute.aws"),
				InstanceId:       aws.String("i-123456789abcedf"),
				VpcId:            aws.String("vpc-abcd"),
				SubnetId:         aws.String("subnet-123"),
				PrivateIpAddress: aws.String("172.31.1.1"),
			},
			},
			respCheck: func(t *testing.T, ldr *ListEC2Response) {
				require.Len(t, ldr.Servers, 1, "expected 1 server, got %d", len(ldr.Servers))
				require.Empty(t, ldr.NextToken, "expected an empty NextToken")

				expectedServer := &types.ServerV2{
					Kind:    "node",
					Version: "v2",
					SubKind: "openssh-ec2-ice",
					Metadata: types.Metadata{
						Labels: map[string]string{
							"account-id":               "123456789012",
							"region":                   "us-east-1",
							"teleport.dev/instance-id": "i-123456789abcedf",
						},
						Namespace: "default",
					},
					Spec: types.ServerSpecV2{
						Addr:     "172.31.1.1:22",
						Hostname: "my-private-dns.compute.aws",
						CloudMetadata: &types.CloudMetadata{
							AWS: &types.AWSInfo{
								AccountID:   "123456789012",
								InstanceID:  "i-123456789abcedf",
								Region:      "us-east-1",
								VPCID:       "vpc-abcd",
								Integration: "myintegration",
								SubnetID:    "subnet-123",
							},
						},
					},
				}
				require.Empty(t, cmp.Diff(expectedServer, ldr.Servers[0], cmpopts.IgnoreFields(types.ServerV2{}, "Metadata.Name")))
			},
			errCheck: noErrorFunc,
		},
		{
			name: "no region",
			req: ListEC2Request{
				Integration: "myintegration",
			},
			errCheck: trace.IsBadParameter,
		},
		{
			name: "no integration",
			req: ListEC2Request{
				Region: "us-east-1",
			},
			errCheck: trace.IsBadParameter,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			mockListClient := &mockListEC2Client{
				pageSize:     pageSize,
				accountID:    "123456789012",
				ec2Instances: tt.mockInstances,
			}
			resp, err := ListEC2(ctx, mockListClient, tt.req)
			require.True(t, tt.errCheck(err), "unexpected err: %v", err)
			if tt.respCheck != nil {
				tt.respCheck(t, resp)
			}
		})
	}
}
