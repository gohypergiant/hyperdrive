package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/docker/distribution/uuid"

	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	HyperConfig "github.com/gohypergiant/hyperdrive/hyper/services/config"
)

const HYPERDRIVE_TYPE_TAG string = "hyperdrive-type"
const HYPERDRIVE_NAME_TAG string = "hyperdrive-name"
const HYPERDRIVE_SECURITY_GROUP_NAME string = "-SecurityGroup"

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

func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}
func GetEC2Client(remoteCfg HyperConfig.EC2RemoteConfiguration) *ec2.Client {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(remoteCfg.Profile))
	cfg.Region = remoteCfg.Region
	if remoteCfg.AccessKey != "" && remoteCfg.Secret != "" {
		cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(remoteCfg.AccessKey, remoteCfg.Secret, uuid.Generate().String()))
	}
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	return ec2.NewFromConfig(cfg)

}
func ListServers(remoteCfg HyperConfig.EC2RemoteConfiguration) {

	result, err := GetHyperdriveInstances(remoteCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	for _, i := range result {
		fmt.Println("   " + GetHyperName(i))
		fmt.Println("")
	}

}
func GetHyperdriveInstances(remoteCfg HyperConfig.EC2RemoteConfiguration) ([]types.Instance, error) {

	client := GetEC2Client(remoteCfg)
	input := &ec2.DescribeInstancesInput{}

	result, err := GetInstances(context.TODO(), client, input)
	instances := []types.Instance{}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return nil, err
	}

	for _, r := range result.Reservations {
		fmt.Println("Instance IDs:")
		for _, i := range r.Instances {
			if IsHyperdriveInstance(i) {
				fmt.Println("   " + GetHyperName(i))
				instances = append(instances, i)
			}
		}

		fmt.Println("")
	}
	return instances, nil
}
func IsHyperdriveInstance(i types.Instance) bool {
	for _, t := range i.Tags {
		if *t.Key == HYPERDRIVE_TYPE_TAG {
			return true
		}
	}
	return false
}
func GetHyperName(i types.Instance) string {
	for _, t := range i.Tags {
		if *t.Key == HYPERDRIVE_NAME_TAG {
			return *t.Value
		}
	}
	return ""
}
func MakeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}
func DeleteInstances(c context.Context, api TerminateInstancesAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}
func GetVpcs(c context.Context, api DescribeVpcsAPI, input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return api.DescribeVpcs(c, input)
}
func MakeVpc(c context.Context, api CreateVpcAPI, input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	return api.CreateVpc(c, input)
}
func DeleteVpc(c context.Context, api DeleteVpcAPI, input *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	return api.DeleteVpc(c, input)
}
func GetSubnets(c context.Context, api DescribeSubnetsAPI, input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return api.DescribeSubnets(c, input)
}
func MakeSubnet(c context.Context, api CreateSubnetAPI, input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	return api.CreateSubnet(c, input)
}
func DeleteSubnet(c context.Context, api DeleteSubnetAPI, input *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	return api.DeleteSubnet(c, input)
}
func ChangeSubnet(c context.Context, api ModifySubnetAttributeAPI, input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	return api.ModifySubnetAttribute(c, input)
}
func GetInternetGateways(c context.Context, api DescribeInternetGatewaysAPI, input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return api.DescribeInternetGateways(c, input)
}
func MakeInternetGateway(c context.Context, api CreateInternetGatewayAPI, input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	return api.CreateInternetGateway(c, input)
}
func DeleteInternetGateway(c context.Context, api DeleteInternetGatewayAPI, input *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	return api.DeleteInternetGateway(c, input)
}
func AttachInternetGateway(c context.Context, api AttachInternetGatewayAPI, input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	return api.AttachInternetGateway(c, input)
}
func DetachInternetGateway(c context.Context, api DetachInternetGatewayAPI, input *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	return api.DetachInternetGateway(c, input)
}
func DescribeRouteTables(c context.Context, api DescribeRouteTablesAPI, input *ec2.DescribeRouteTablesInput) (*ec2.DescribeRouteTablesOutput, error) {
	return api.DescribeRouteTables(c, input)
}
func MakeRouteTable(c context.Context, api CreateRouteTableAPI, input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	return api.CreateRouteTable(c, input)
}
func DeleteRouteTable(c context.Context, api DeleteRouteTableAPI, input *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	return api.DeleteRouteTable(c, input)
}
func AddRouteTable(c context.Context, api AssociateRouteTableAPI, input *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	return api.AssociateRouteTable(c, input)
}
func DisassociateRouteTable(c context.Context, api DisassociateRouteTableAPI, input *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	return api.DisassociateRouteTable(c, input)
}
func AddRoute(c context.Context, api CreateRouteAPI, input *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	return api.CreateRoute(c, input)
}
func GetSecurityGroups(c context.Context, api DescribeSecurityGroupsAPI, input *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	return api.DescribeSecurityGroups(c, input)
}
func MakeSecurityGroup(c context.Context, api CreateSecurityGroupAPI, input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	return api.CreateSecurityGroup(c, input)
}
func DeleteSecurityGroup(c context.Context, api DeleteSecurityGroupAPI, input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	return api.DeleteSecurityGroup(c, input)
}
func MakeSecurityGroupPermissions(c context.Context, api AddSecurityGroupPermissionsAPI, input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return api.AuthorizeSecurityGroupIngress(c, input)
}
func GetKeyPairs(c context.Context, api DescribeKeyPairsAPI, input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return api.DescribeKeyPairs(c, input)
}
func MakeKeyPair(c context.Context, api CreateKeyPairAPI, input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	return api.CreateKeyPair(c, input)
}
func DeleteKeyPair(c context.Context, api DeleteKeyPairAPI, input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	return api.DeleteKeyPair(c, input)
}
func CreateInternetGateway(client *ec2.Client, projectName string) string {
	tagSpecification := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeInternetGateway,
			Tags: []types.Tag{
				{
					Key:   aws.String(HYPERDRIVE_TYPE_TAG),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String(HYPERDRIVE_NAME_TAG),
					Value: aws.String(projectName),
				},
			},
		},
	}
	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: tagSpecification,
	}

	result, err := MakeInternetGateway(context.TODO(), client, input)
	if err != nil {
		panic("error creating Internet Gateway: " + err.Error())
	}

	return *result.InternetGateway.InternetGatewayId
}
func AddInternetGateway(client *ec2.Client, vID string, igID string) {

	input := &ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String(igID),
		VpcId:             aws.String(vID),
	}

	_, err := AttachInternetGateway(context.TODO(), client, input)
	if err != nil {
		panic("error attaching the Internet Gateway to VPC: " + err.Error())
	}
}
func GetVpcId(r *ec2.DescribeVpcsOutput, projectName string) string {

	for _, v := range r.Vpcs {
		for _, t := range v.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *v.VpcId
			}
		}
	}

	return ""
}

func getOrCreateVPC(client *ec2.Client, projectName string) (string, string) {
	var routeTableID string

	vpcDescribeInput := &ec2.DescribeVpcsInput{}
	result, err := GetVpcs(context.TODO(), client, vpcDescribeInput)
	if err != nil {
		panic("error when fetching VPCs, " + err.Error())
	}

	vpcID := GetVpcId(result, projectName)

	if vpcID == "" {
		fmt.Println("Creating VPC")

		tagSpecification := []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
				Tags: []types.Tag{
					{
						Key:   aws.String(HYPERDRIVE_TYPE_TAG),
						Value: aws.String("true"),
					},
					{
						Key:   aws.String(HYPERDRIVE_NAME_TAG),
						Value: aws.String(projectName),
					},
				},
			},
		}

		inputMakeVPC := &ec2.CreateVpcInput{
			CidrBlock:         aws.String("10.0.0.0/16"),
			TagSpecifications: tagSpecification,
		}

		resultMakeVPC, err := MakeVpc(context.TODO(), client, inputMakeVPC)
		if err != nil {
			panic("error creating VPC," + err.Error())
		}
		vpcID = *resultMakeVPC.Vpc.VpcId

		internetGatewayID := CreateInternetGateway(client, projectName)
		fmt.Println("Internet Gateway ID:", internetGatewayID)

		AddInternetGateway(client, vpcID, internetGatewayID)

		tagSpecification = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeRouteTable,
				Tags: []types.Tag{
					{
						Key:   aws.String(HYPERDRIVE_TYPE_TAG),
						Value: aws.String("true"),
					},
					{
						Key:   aws.String(HYPERDRIVE_NAME_TAG),
						Value: aws.String(projectName),
					},
				},
			},
		}

		inputMakeRouteTable := &ec2.CreateRouteTableInput{
			VpcId:             aws.String(vpcID),
			TagSpecifications: tagSpecification,
		}

		resultMakeRouteTable, err := MakeRouteTable(context.TODO(), client, inputMakeRouteTable)
		if err != nil {
			panic("error creating Route Table," + err.Error())
		}

		routeTableID = *resultMakeRouteTable.RouteTable.RouteTableId

		inputAddRoute := &ec2.CreateRouteInput{
			RouteTableId:         aws.String(routeTableID),
			DestinationCidrBlock: aws.String("0.0.0.0/0"),
			GatewayId:            aws.String(internetGatewayID),
		}

		_, err = AddRoute(context.TODO(), client, inputAddRoute)
		if err != nil {
			panic("error adding Route to Route Table," + err.Error())
		}
	}
	return vpcID, routeTableID
}

func GetSubnetID(r *ec2.DescribeSubnetsOutput, projectName string) string {

	for _, s := range r.Subnets {
		for _, t := range s.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *s.SubnetId
			}
		}
	}
	return ""
}
func setSubnetToProvisionPublicIP(subnetID string, client *ec2.Client) {
	subnetChangeInput := &ec2.ModifySubnetAttributeInput{
		SubnetId: aws.String(subnetID),
		MapPublicIpOnLaunch: &types.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	}

	_, err := ChangeSubnet(context.TODO(), client, subnetChangeInput)
	if err != nil {
		panic("error modifying Subnet attribute," + err.Error())
	}

}
func getOrCreateSubnet(client *ec2.Client, vID string, region string, projectName string, rtID string) string {

	subnetDescribeInput := &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vID},
			},
		},
	}
	subnetDescribeResult, err := GetSubnets(context.TODO(), client, subnetDescribeInput)
	if err != nil {
		panic("error fetching Subnets," + err.Error())
	}

	subnetID := GetSubnetID(subnetDescribeResult, projectName)

	if subnetID != "" {
		setSubnetToProvisionPublicIP(subnetID, client)
		return subnetID
	}
	fmt.Println("No Subnet found for VPC", vID)
	fmt.Println("Creating Subnet")

	tagSpecification := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeSubnet,
			Tags: []types.Tag{
				{
					Key:   aws.String(HYPERDRIVE_TYPE_TAG),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String(HYPERDRIVE_NAME_TAG),
					Value: aws.String(projectName),
				},
			},
		},
	}

	subnetMakeInput := &ec2.CreateSubnetInput{
		CidrBlock:         aws.String("10.0.1.0/24"),
		VpcId:             aws.String(vID),
		AvailabilityZone:  aws.String(region + "a"),
		TagSpecifications: tagSpecification,
	}

	subnetMakeResult, err := MakeSubnet(context.TODO(), client, subnetMakeInput)
	if err != nil {
		panic("error creating Subnet," + err.Error())
	}

	subnetID = *subnetMakeResult.Subnet.SubnetId

	setSubnetToProvisionPublicIP(subnetID, client)
	inputAddRouteTable := &ec2.AssociateRouteTableInput{
		RouteTableId: aws.String(rtID),
		SubnetId:     aws.String(subnetID),
	}

	_, err = AddRouteTable(context.TODO(), client, inputAddRouteTable)
	if err != nil {
		panic("error associating Route Table to Subnet," + err.Error())
	}

	return subnetID
}

func GetSecurityGroupId(r *ec2.DescribeSecurityGroupsOutput, projectName string) string {

	for _, s := range r.SecurityGroups {
		for _, t := range s.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *s.GroupId
			}
		}
	}
	return ""
}

func getOrCreateSecurityGroup(client *ec2.Client, vID string, projectName string) string {
	var securityGroupID string

	securityGroupDescribeInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vID},
			},
		},
	}
	securityGroupDescribeResult, err := GetSecurityGroups(context.TODO(), client, securityGroupDescribeInput)
	if err != nil {
		panic("error fetching Security Groups," + err.Error())
	}

	securityGroupID = GetSecurityGroupId(securityGroupDescribeResult, projectName)

	if securityGroupID != "" {
		return securityGroupID
	}
	fmt.Println("No Security Group found on VPC", vID)
	fmt.Println("Creating Security Group for project", projectName)

	tagSpecification := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeSecurityGroup,
			Tags: []types.Tag{
				{
					Key:   aws.String(HYPERDRIVE_TYPE_TAG),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String(HYPERDRIVE_NAME_TAG),
					Value: aws.String(projectName),
				},
			},
		},
	}

	scInput := &ec2.CreateSecurityGroupInput{
		GroupName:         aws.String(projectName + HYPERDRIVE_SECURITY_GROUP_NAME),
		Description:       aws.String("Security group for EC2 instances provisioned by hyper"),
		VpcId:             aws.String(vID),
		TagSpecifications: tagSpecification,
	}

	securityGroupMakeResult, err := MakeSecurityGroup(context.TODO(), client, scInput)
	if err != nil {
		panic("error creating Security Group," + err.Error())
	}

	securityGroupID = *securityGroupMakeResult.GroupId

	securityGroupPermissionsInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: aws.String(securityGroupID),
		IpPermissions: []types.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(22),
				ToPort:     aws.Int32(22),
				IpRanges: []types.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int32(8888),
				ToPort:     aws.Int32(8888),
				IpRanges: []types.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
		},
	}

	_, err = MakeSecurityGroupPermissions(context.TODO(), client, securityGroupPermissionsInput)
	if err != nil {
		panic("error adding permissions to the Security Group," + err.Error())
	}
	return securityGroupID
}
func getKeyPairName(r *ec2.DescribeKeyPairsOutput, projectName string) string {

	for _, s := range r.KeyPairs {
		for _, t := range s.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *s.KeyName
			}
		}
	}
	return ""
}
func WriteKey(fileName string, fileData *string) error {
	err := os.WriteFile(fileName, []byte(*fileData), 0400)
	return err
}
func getOrCreateKeyPair(client *ec2.Client, projectName string) string {

	keyPairDescribeInput := &ec2.DescribeKeyPairsInput{
		IncludePublicKey: aws.Bool(true),
	}

	keyPairDescribeResult, err := GetKeyPairs(context.TODO(), client, keyPairDescribeInput)
	if err != nil {
		panic("error fetching Key Pairs: " + err.Error())
	}
	keyName := getKeyPairName(keyPairDescribeResult, projectName)
	fmt.Printf("keyName: %s", keyName)

	if keyName == "" {
		tagSpecification := []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeKeyPair,
				Tags: []types.Tag{
					{
						Key:   aws.String(HYPERDRIVE_TYPE_TAG),
						Value: aws.String("true"),
					},
					{
						Key:   aws.String(HYPERDRIVE_NAME_TAG),
						Value: aws.String(projectName),
					},
				},
			},
		}

		keyName = projectName + "-hyperdrive"
		keyPairMakeInput := &ec2.CreateKeyPairInput{
			KeyName:           aws.String(keyName),
			TagSpecifications: tagSpecification,
		}

		keyPairMakeResult, err := MakeKeyPair(context.TODO(), client, keyPairMakeInput)
		if err != nil {
			panic("error creating Key Pair," + err.Error())
		}
		keyPath := path.Join(os.Getenv("HOME"), fmt.Sprintf("%s.pem", keyName))
		err = WriteKey(keyPath, keyPairMakeResult.KeyMaterial) // todo fix path for windows
		if err != nil {
			fmt.Printf("Couldn't write key pair to file: %v", err)
		}

	}
	return keyName
}

func StartServer(manifestPath string, remoteCfg HyperConfig.EC2RemoteConfiguration, ec2Type string, amiID string) {

	if ec2Type == "" {
		fmt.Println("EC2InstanceTypeNotFound: please specify a EC2 instance type using the flag --ec2InstanceType")
		return
	}
	projectName := manifest.GetProjectName(manifestPath)
	if projectName == "" {
		fmt.Println("ProjectNameNotFound: please specify a project_name on the manifest (", manifestPath, ")")
		return
	}
	fmt.Println("Project name is:", projectName)
	client := GetEC2Client(remoteCfg)

	vpcID, rtID := getOrCreateVPC(client, projectName)
	fmt.Println("VPC ID:", vpcID)
	fmt.Println("Route Table ID:", rtID)

	subnetID := getOrCreateSubnet(client, vpcID, remoteCfg.Region, projectName, rtID)
	fmt.Println("Subnet ID:", subnetID)

	securityGroupID := getOrCreateSecurityGroup(client, vpcID, projectName)
	fmt.Println("Security group ID:", securityGroupID)

	keyName := getOrCreateKeyPair(client, projectName)
	fmt.Println("Key name:", keyName)

	tagSpecification := []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{
				{
					Key:   aws.String(HYPERDRIVE_TYPE_TAG),
					Value: aws.String("true"),
				},
				{
					Key:   aws.String(HYPERDRIVE_NAME_TAG),
					Value: aws.String(projectName),
				},
			},
		},
	}

	version := "0.0.2"
	minMaxCount := int32(1)
	startupScript := fmt.Sprintf(`
#!/bin/bash -xe
#yum update -y
service docker start
mkdir -p /tmp/hyperdrive/project
curl -fsSL https://github.com/gohypergiant/hyperdrive/releases/download/%s/hyperdrive_%s_Linux_x86_64.tar.gz -o /tmp/hyperdrive/hyper.tar
tar -xvf /tmp/hyperdrive/hyper.tar -C /tmp/hyperdrive
mv /tmp/hyperdrive/hyper /usr/bin/hyper
sudo chown ec2-user:ec2-user /tmp/hyperdrive/project
cd /tmp/hyperdrive/project
sudo -u ec2-user bash -c 'hyper jupyter remoteHost --hostPort 8888 &'
`, version, version)
	ec2Input := &ec2.RunInstancesInput{
		ImageId:           aws.String(amiID),
		InstanceType:      types.InstanceType(*aws.String(ec2Type)),
		MinCount:          aws.Int32(minMaxCount),
		MaxCount:          aws.Int32(minMaxCount),
		SecurityGroupIds:  []string{securityGroupID},
		SubnetId:          aws.String(subnetID),
		KeyName:           aws.String(keyName),
		TagSpecifications: tagSpecification,
		UserData:          aws.String(base64.StdEncoding.EncodeToString([]byte(startupScript))),
	}

	result, err := MakeInstance(context.TODO(), client, ec2Input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}

	ip := result.Instances[0].PublicIpAddress
	if ip == nil {

		instances, err := GetHyperdriveInstances(remoteCfg)
		if err != nil {

			fmt.Println("Provisioned instance but cannot get publicIP")
			return
		}
		for _, i := range instances {
			if *i.InstanceId == *result.Instances[0].InstanceId {
				ip = i.PublicIpAddress
				break
			}
		}

	}
	if ip == nil {

		fmt.Println("Provisioned instance but cannot get publicIP")
		return
	}
	fmt.Println("")
	fmt.Println("EC2 instance provisioned. You can access via ssh by running:")
	fmt.Println("ssh -i " + keyName + ".pem ec2-user@" + *ip)
	fmt.Println("")
	fmt.Println("In a few minutes, you should be able to access jupyter lab at http://" + *ip + ":8888/lab")
}
func GetRouteTableID(r *ec2.DescribeRouteTablesOutput, projectName string) string {

	for _, r := range r.RouteTables {
		for _, t := range r.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *r.RouteTableId
			}
		}
	}
	return ""
}
func GetAssociationID(r *ec2.DescribeRouteTablesOutput, subnetID string) string {
	for _, rt := range r.RouteTables {
		for _, a := range rt.Associations {
			if a.SubnetId == &subnetID {
				return *a.RouteTableAssociationId
			}
		}
	}
	return ""
}
func GetInternetGatewayID(r *ec2.DescribeInternetGatewaysOutput, projectName string) string {
	for _, i := range r.InternetGateways {
		for _, t := range i.Tags {
			if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == projectName {
				return *i.InternetGatewayId
			}
		}
	}
	return ""
}
func StopServer(manifestPath string, remoteCfg HyperConfig.EC2RemoteConfiguration) {
	projectName := manifest.GetProjectName(manifestPath)
	if projectName == "" {
		fmt.Println("ProjectNameNotFound: please specify a project_name on the manifest (", manifestPath, ")")
		return
	}

	client := GetEC2Client(remoteCfg)

	vpcDescribeInput := &ec2.DescribeVpcsInput{}
	vpcDescribeResult, err := GetVpcs(context.TODO(), client, vpcDescribeInput)
	if err != nil {
		panic("error when fetching VPCs, " + err.Error())
	}

	vpcID := GetVpcId(vpcDescribeResult, projectName)

	if vpcID == "" {
		fmt.Println("No VPC associated with this project found")
		return
	}

	ec2DescribeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
			{
				Name:   aws.String("instance-state-code"),
				Values: []string{"16"},
			},
		},
	}

	ec2DescribeResult, err := GetInstances(context.TODO(), client, ec2DescribeInput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	for _, r := range ec2DescribeResult.Reservations {
		for _, i := range r.Instances {
			if IsHyperdriveInstance(i) {
				instanceTerminateInput := &ec2.TerminateInstancesInput{
					InstanceIds: []string{*i.InstanceId},
				}

				_, err = DeleteInstances(context.TODO(), client, instanceTerminateInput)
				if err != nil {
					panic("error deleting Instances: " + err.Error())
				}
			}

			instanceStatusWaiter := ec2.NewInstanceTerminatedWaiter(client)
			instanceStatusWaitInput := &ec2.DescribeInstancesInput{
				Filters: []types.Filter{
					{
						Name:   aws.String("instance-id"),
						Values: []string{*i.InstanceId},
					},
				},
			}
			_, err := instanceStatusWaiter.WaitForOutput(context.TODO(), instanceStatusWaitInput, time.Duration(240*time.Second))
			if err != nil {
				panic("error waiting to delete Instance: " + err.Error())
			}

			fmt.Println("Instance deleted:", *i.InstanceId)
		}
	}

	keyPairDescribeInput := &ec2.DescribeKeyPairsInput{
		IncludePublicKey: aws.Bool(true),
	}

	keyPairDescribeResult, err := GetKeyPairs(context.TODO(), client, keyPairDescribeInput)
	if err != nil {
		panic("error fetching Key Pairs: " + err.Error())
	}
	keyName := getKeyPairName(keyPairDescribeResult, projectName)

	if keyName != "" {
		keyPairDeleteInput := &ec2.DeleteKeyPairInput{
			KeyName: aws.String(keyName),
		}

		_, err = DeleteKeyPair(context.TODO(), client, keyPairDeleteInput)
		if err != nil {
			panic("error deleting Key Pair: " + err.Error())
		}
		fmt.Println("Key Pair deleted:", keyName)

		keyPath := path.Join(os.Getenv("HOME"), fmt.Sprintf("%s.pem", keyName))
		_, err := os.Stat(keyPath)
		if err == nil {
			err = os.Remove(keyPath)
			if err != nil {
				panic("error deleting local pem file: " + err.Error())
			}
		}
	}

	securityGroupDescribeInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	}
	securityGroupDescribeResult, err := GetSecurityGroups(context.TODO(), client, securityGroupDescribeInput)
	if err != nil {
		panic("error fetching Security Groups," + err.Error())
	}

	securityGroupID := GetSecurityGroupId(securityGroupDescribeResult, projectName)

	if securityGroupID != "" {
		securityGroupDeleteInput := &ec2.DeleteSecurityGroupInput{
			GroupId: aws.String(securityGroupID),
		}

		_, err = DeleteSecurityGroup(context.TODO(), client, securityGroupDeleteInput)
		if err != nil {
			panic("error deleting Security Group," + err.Error())
		}
		fmt.Println("Security Group deleted:", securityGroupID)
	}

	subnetDescribeInput := &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcID},
			},
		},
	}
	subnetDescribeResult, err := GetSubnets(context.TODO(), client, subnetDescribeInput)
	if err != nil {
		panic("error fetching Subnets," + err.Error())
	}

	subnetID := GetSubnetID(subnetDescribeResult, projectName)

	if subnetID != "" {
		routeTableDescribeInput := &ec2.DescribeRouteTablesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{vpcID},
				},
				{
					Name:   aws.String("association.subnet-id"),
					Values: []string{subnetID},
				},
			},
		}

		routeTableDescribeOutput, err := DescribeRouteTables(context.TODO(), client, routeTableDescribeInput)
		if err != nil {
			panic("error fetching Route Tables," + err.Error())
		}

		for _, rt := range routeTableDescribeOutput.RouteTables {
			routeTableID := *rt.RouteTableId

			for _, a := range rt.Associations {
				associationID := *a.RouteTableAssociationId

				routeTableDisassociateInput := &ec2.DisassociateRouteTableInput{
					AssociationId: aws.String(associationID),
				}
				_, err = DisassociateRouteTable(context.TODO(), client, routeTableDisassociateInput)
				if err != nil {
					panic("error disassociationg Route Table," + err.Error())
				}
			}
			routeTableDeleteInput := &ec2.DeleteRouteTableInput{
				RouteTableId: aws.String(routeTableID),
			}
			_, err = DeleteRouteTable(context.TODO(), client, routeTableDeleteInput)
			if err != nil {
				panic("error deleting Route Table," + err.Error())
			}
			fmt.Println("Route Table deleted:", routeTableID)
		}
	}

	internetGatewayDescribeInput := &ec2.DescribeInternetGatewaysInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("attachment.vpc-id"),
				Values: []string{vpcID},
			},
		},
	}
	internetGatewayDescribeResult, err := GetInternetGateways(context.TODO(), client, internetGatewayDescribeInput)
	if err != nil {
		panic("error fetching Internet Gateways," + err.Error())
	}
	internetGatewayID := GetInternetGatewayID(internetGatewayDescribeResult, projectName)

	if internetGatewayID != "" {
		internetGatewayDetachInput := &ec2.DetachInternetGatewayInput{
			InternetGatewayId: aws.String(internetGatewayID),
			VpcId:             aws.String(vpcID),
		}
		_, err = DetachInternetGateway(context.TODO(), client, internetGatewayDetachInput)
		if err != nil {
			panic("error detaching Internet Gateway," + err.Error())
		}

		internetGatewayDeleteInput := &ec2.DeleteInternetGatewayInput{
			InternetGatewayId: aws.String(internetGatewayID),
		}

		_, err = DeleteInternetGateway(context.TODO(), client, internetGatewayDeleteInput)
		if err != nil {
			panic("error deleting Internet Gateway," + err.Error())
		}
		fmt.Println("Internet Gateway deleted:", internetGatewayID)
	}

	if subnetID != "" {
		subnetDeleteInput := &ec2.DeleteSubnetInput{
			SubnetId: aws.String(subnetID),
		}
		_, err = DeleteSubnet(context.TODO(), client, subnetDeleteInput)
		if err != nil {
			panic("error deleting Subnet," + err.Error())
		}
		fmt.Println("Subnet deleted:", subnetID)
	}

	if vpcID != "" {
		vpcDeleteInput := &ec2.DeleteVpcInput{
			VpcId: aws.String(vpcID),
		}
		_, err = DeleteVpc(context.TODO(), client, vpcDeleteInput)
		if err != nil {
			panic("error deleting VPC," + err.Error())
		}
		fmt.Println("VPC deleted:", vpcID)
	}
}
