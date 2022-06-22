package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gohypergiant/hyperdrive/hyper/types"
)

/*
* No logic here
 */

func MakeInstance(c context.Context, api types.EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}
func DeleteInstances(c context.Context, api types.TerminateInstancesAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}
func GetSubnets(c context.Context, api types.DescribeSubnetsAPI, input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return api.DescribeSubnets(c, input)
}
func MakeSubnet(c context.Context, api types.CreateSubnetAPI, input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	return api.CreateSubnet(c, input)
}
func DeleteSubnet(c context.Context, api types.DeleteSubnetAPI, input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	return api.DeleteSubnet(c, input)
}
func ChangeSubnet(c context.Context, api types.ModifySubnetAttributeAPI, input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	return api.ModifySubnetAttribute(c, input)
}
func GetInternetGateways(c context.Context, api types.DescribeInternetGatewaysAPI, input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return api.DescribeInternetGateways(c, input)
}
func MakeInternetGateway(c context.Context, api types.CreateInternetGatewayAPI, input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	return api.CreateInternetGateway(c, input)
}
func DeleteInternetGateway(c context.Context, api types.DeleteInternetGatewayAPI, input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	return api.DeleteInternetGateway(c, input)
}
func AttachInternetGateway(c context.Context, api types.AttachInternetGatewayAPI, input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	return api.AttachInternetGateway(c, input)
}
func DetachInternetGateway(c context.Context, api types.DetachInternetGatewayAPI, input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	return api.DetachInternetGateway(c, input)
}
func DescribeRouteTables(c context.Context, api types.DescribeRouteTablesAPI, input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return api.DescribeRouteTables(c, input)
}
func MakeRouteTable(c context.Context, api types.CreateRouteTableAPI, input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	return api.CreateRouteTable(c, input)
}
func DeleteRouteTable(c context.Context, api types.DeleteRouteTableAPI, input *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	return api.DeleteRouteTable(c, input)
}
func AddRouteTable(c context.Context, api types.AssociateRouteTableAPI, input *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	return api.AssociateRouteTable(c, input)
}
func DisassociateRouteTable(c context.Context, api types.DisassociateRouteTableAPI, input *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	return api.DisassociateRouteTable(c, input)
}
func AddRoute(c context.Context, api types.CreateRouteAPI, input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	return api.CreateRoute(c, input)
}
func GetSecurityGroups(c context.Context, api types.DescribeSecurityGroupsAPI, input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return api.DescribeSecurityGroups(c, input)
}
func MakeSecurityGroup(c context.Context, api types.CreateSecurityGroupAPI, input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	return api.CreateSecurityGroup(c, input)
}
func DeleteSecurityGroup(c context.Context, api types.DeleteSecurityGroupAPI, input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	return api.DeleteSecurityGroup(c, input)
}
func MakeSecurityGroupPermissions(c context.Context, api types.AddSecurityGroupPermissionsAPI, input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return api.AuthorizeSecurityGroupIngress(c, input)
}
func GetKeyPairs(c context.Context, api types.DescribeKeyPairsAPI, input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return api.DescribeKeyPairs(c, input)
}
func MakeKeyPair(c context.Context, api types.CreateKeyPairAPI, input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	return api.CreateKeyPair(c, input)
}
func ImportKeyPair(c context.Context, api types.ImportKeyPairAPI, input *ec2.ImportKeyPairInput) (*ec2.ImportKeyPairOutput, error) {
	return api.ImportKeyPair(c, input)
}
func DeleteKeyPair(c context.Context, api types.DeleteKeyPairAPI, input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	return api.DeleteKeyPair(c, input)
}
func GetVpcs(c context.Context, api types.DescribeVpcsAPI, input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return api.DescribeVpcs(c, input)
}
func MakeVpc(c context.Context, api types.CreateVpcAPI, input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	return api.CreateVpc(c, input)
}
func DeleteVpc(c context.Context, api types.DeleteVpcAPI, input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	return api.DeleteVpc(c, input)
}
