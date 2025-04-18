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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/gravitational/trace"
)

const (
	// allProtocols is a sentinel value used to identify a rule which allows all IP protocols.
	allProtocols = "all"
)

// ListSecurityGroupsRequest contains the required fields to list VPC Security Groups.
type ListSecurityGroupsRequest struct {
	// VPCID is the VPC to filter Security Groups.
	VPCID string

	// NextToken is the token to be used to fetch the next page.
	// If empty, the first page is fetched.
	NextToken string
}

// CheckAndSetDefaults checks if the required fields are present.
func (req *ListSecurityGroupsRequest) CheckAndSetDefaults() error {
	if req.VPCID == "" {
		return trace.BadParameter("vpc id is required")
	}

	return nil
}

// SecurityGroup is the Teleport representation of an EC2 Instance Connect Endpoint
type SecurityGroup struct {
	// Name is the Security Group name.
	// This is just a friendly name and should not be used for further API calls
	Name string `json:"name"`

	// ID is the security group ID.
	// This is the value that should be used when doing further API calls.
	ID string `json:"id"`

	// Description is a small description of the Security Group.
	// Might be empty.
	Description string `json:"description"`

	// InboundRules describe the Security Group Inbound Rules.
	// The CIDR of each rule represents the source IP that the rule applies to.
	InboundRules []SecurityGroupRule `json:"inboundRules"`

	// OutboundRules describe the Security Group Outbound Rules.
	// The CIDR of each rule represents the destination IP that the rule applies to.
	OutboundRules []SecurityGroupRule `json:"outboundRules"`
}

// SecurityGroupRule is a SecurityGroup role.
// It describes which protocol, port range and a list of IPs the rule applies to.
type SecurityGroupRule struct {
	// IPProtocol is the protocol used to describe the rule.
	// If the rule applies to all protocols, the "all" value is used.
	// The IP protocol name ( tcp , udp , icmp , icmpv6 ) or number (see Protocol
	// Numbers (http://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml)).
	IPProtocol string `json:"ipProtocol"`

	// FromPort is the inclusive start of the Port range for the Rule.
	FromPort int `json:"fromPort"`

	// ToPort is the inclusive end of the Port range for the Rule.
	ToPort int `json:"toPort"`

	// CIDRs contains a list of IP ranges that this rule applies to and a description for the value.
	CIDRs []CIDR `json:"cidrs"`
}

// CIDR has a CIDR (IP Range) and a description for the value.
type CIDR struct {
	// CIDR is the IP range using CIDR notation.
	CIDR string `json:"cidr"`
	// Description contains a small text describing the CIDR.
	Description string `json:"description"`
}

// ListSecurityGroupsResponse contains a page of SecurityGroups.
type ListSecurityGroupsResponse struct {
	// SecurityGroups contains the page of VPC Security Groups.
	SecurityGroups []SecurityGroup `json:"securityGroups"`

	// NextToken is used for pagination.
	// If non-empty, it can be used to request the next page.
	NextToken string `json:"nextToken"`
}

// ListSecurityGroupsClient describes the required methods to List Security Groups a 3rd Party API.
type ListSecurityGroupsClient interface {
	// DescribeSecurityGroups describes the specified security groups or all of your security groups.
	DescribeSecurityGroups(ctx context.Context, params *ec2.DescribeSecurityGroupsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
}

type defaultListSecurityGroupsClient struct {
	*ec2.Client
}

// NewListSecurityGroupsClient creates a new ListSecurityGroupsClient using a AWSClientRequest.
func NewListSecurityGroupsClient(ctx context.Context, req *AWSClientRequest) (ListSecurityGroupsClient, error) {
	ec2Client, err := newEC2Client(ctx, req)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &defaultListSecurityGroupsClient{
		Client: ec2Client,
	}, nil
}

// ListSecurityGroups calls the following AWS API:
// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeSecurityGroups.html
// It returns a list of VPC Security Groups and an optional NextToken that can be used to fetch the next page
func ListSecurityGroups(ctx context.Context, clt ListSecurityGroupsClient, req ListSecurityGroupsRequest) (*ListSecurityGroupsResponse, error) {
	if err := req.CheckAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	describeSecurityGroups := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2Types.Filter{{
			Name:   aws.String("vpc-id"),
			Values: []string{req.VPCID},
		}},
	}
	if req.NextToken != "" {
		describeSecurityGroups.NextToken = &req.NextToken
	}

	securityGroupsResp, err := clt.DescribeSecurityGroups(ctx, describeSecurityGroups)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &ListSecurityGroupsResponse{
		NextToken:      aws.ToString(securityGroupsResp.NextToken),
		SecurityGroups: convertAWSSecurityGroups(securityGroupsResp.SecurityGroups),
	}, nil
}

func convertAWSSecurityGroups(awsSG []ec2Types.SecurityGroup) []SecurityGroup {
	ret := make([]SecurityGroup, 0, len(awsSG))
	for _, sg := range awsSG {
		ret = append(ret, SecurityGroup{
			Name:          aws.ToString(sg.GroupName),
			ID:            aws.ToString(sg.GroupId),
			Description:   aws.ToString(sg.Description),
			InboundRules:  convertAWSIPPermissions(sg.IpPermissions),
			OutboundRules: convertAWSIPPermissions(sg.IpPermissionsEgress),
		})
	}

	return ret
}

func convertAWSIPPermissions(permissions []ec2Types.IpPermission) []SecurityGroupRule {
	rules := make([]SecurityGroupRule, 0, len(permissions))
	for _, permission := range permissions {
		ipProtocol := allProtocols
		// From AWS Docs:
		// > Use -1 to specify all protocols.
		if aws.ToString(permission.IpProtocol) != "-1" {
			ipProtocol = aws.ToString(permission.IpProtocol)
		}

		cidrs := make([]CIDR, 0, len(permission.IpRanges))
		for _, r := range permission.IpRanges {
			cidrs = append(cidrs, CIDR{
				CIDR:        aws.ToString(r.CidrIp),
				Description: aws.ToString(r.Description),
			})
		}

		fromPort := int(aws.ToInt32(permission.FromPort))
		toPort := int(aws.ToInt32(permission.ToPort))

		rules = append(rules, SecurityGroupRule{
			IPProtocol: ipProtocol,
			FromPort:   fromPort,
			ToPort:     toPort,
			CIDRs:      cidrs,
		})
	}

	return rules
}
