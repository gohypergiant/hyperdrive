package types

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}
type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
}
type TerminateInstancesAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}
type DescribeVpcsAPI interface {
	DescribeVpcs(ctx context.Context,
		params *ec2.DescribeVpcsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
}
type CreateVpcAPI interface {
	CreateVpc(ctx context.Context,
		params *ec2.CreateVpcInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateVpcOutput, error)
}
type DeleteVpcAPI interface {
	DeleteVpc(ctx context.Context,
		params *ec2.DeleteVpcInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteVpcOutput, error)
}
type DescribeSubnetsAPI interface {
	DescribeSubnets(ctx context.Context,
		params *ec2.DescribeSubnetsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
}
type CreateSubnetAPI interface {
	CreateSubnet(ctx context.Context,
		params *ec2.CreateSubnetInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateSubnetOutput, error)
}
type DeleteSubnetAPI interface {
	DeleteSubnet(ctx context.Context,
		params *ec2.DeleteSubnetInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteSubnetOutput, error)
}
type ModifySubnetAttributeAPI interface {
	ModifySubnetAttribute(ctx context.Context,
		params *ec2.ModifySubnetAttributeInput,
		optFns ...func(*ec2.Options)) (*ec2.ModifySubnetAttributeOutput, error)
}
type DescribeInternetGatewaysAPI interface {
	DescribeInternetGateways(ctx context.Context,
		params *ec2.DescribeInternetGatewaysInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error)
}
type CreateInternetGatewayAPI interface {
	CreateInternetGateway(ctx context.Context,
		params *ec2.CreateInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateInternetGatewayOutput, error)
}
type DeleteInternetGatewayAPI interface {
	DeleteInternetGateway(ctx context.Context,
		params *ec2.DeleteInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteInternetGatewayOutput, error)
}
type AttachInternetGatewayAPI interface {
	AttachInternetGateway(ctx context.Context,
		params *ec2.AttachInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.AttachInternetGatewayOutput, error)
}
type DetachInternetGatewayAPI interface {
	DetachInternetGateway(ctx context.Context,
		params *ec2.DetachInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.DetachInternetGatewayOutput, error)
}
type DescribeRouteTablesAPI interface {
	DescribeRouteTables(ctx context.Context,
		params *ec2.DescribeRouteTablesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeRouteTablesOutput, error)
}
type CreateRouteTableAPI interface {
	CreateRouteTable(ctx context.Context,
		params *ec2.CreateRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateRouteTableOutput, error)
}
type DeleteRouteTableAPI interface {
	DeleteRouteTable(ctx context.Context,
		params *ec2.DeleteRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteRouteTableOutput, error)
}
type AssociateRouteTableAPI interface {
	AssociateRouteTable(ctx context.Context,
		params *ec2.AssociateRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.AssociateRouteTableOutput, error)
}
type DisassociateRouteTableAPI interface {
	DisassociateRouteTable(ctx context.Context,
		params *ec2.DisassociateRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.DisassociateRouteTableOutput, error)
}
type CreateRouteAPI interface {
	CreateRoute(ctx context.Context,
		params *ec2.CreateRouteInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateRouteOutput, error)
}
type DescribeSecurityGroupsAPI interface {
	DescribeSecurityGroups(ctx context.Context,
		params *ec2.DescribeSecurityGroupsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeSecurityGroupsOutput, error)
}
type CreateSecurityGroupAPI interface {
	CreateSecurityGroup(ctx context.Context,
		params *ec2.CreateSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
}
type DeleteSecurityGroupAPI interface {
	DeleteSecurityGroup(ctx context.Context,
		params *ec2.DeleteSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
}
type AddSecurityGroupPermissionsAPI interface {
	AuthorizeSecurityGroupIngress(ctx context.Context,
		params *ec2.AuthorizeSecurityGroupIngressInput,
		optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
}
type DescribeKeyPairsAPI interface {
	DescribeKeyPairs(ctx context.Context,
		params *ec2.DescribeKeyPairsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error)
}
type CreateKeyPairAPI interface {
	CreateKeyPair(ctx context.Context,
		params *ec2.CreateKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error)
}
type DeleteKeyPairAPI interface {
	DeleteKeyPair(ctx context.Context,
		params *ec2.DeleteKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteKeyPairOutput, error)
}
type EC2StartOptions struct {
	InstanceType string
	AmiId        string
}
