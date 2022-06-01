package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"

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
type ModifySubnetAttributeAPI interface {
	ModifySubnetAttribute(ctx context.Context,
		params *ec2.ModifySubnetAttributeInput,
		optFns ...func(*ec2.Options)) (*ec2.ModifySubnetAttributeOutput, error)
}
type CreateInternetGatewayAPI interface {
	CreateInternetGateway(ctx context.Context,
		params *ec2.CreateInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateInternetGatewayOutput, error)
}
type AttachInternetGatewayAPI interface {
	AttachInternetGateway(ctx context.Context,
		params *ec2.AttachInternetGatewayInput,
		optFns ...func(*ec2.Options)) (*ec2.AttachInternetGatewayOutput, error)
}
type CreateRouteTableAPI interface {
	CreateRouteTable(ctx context.Context,
		params *ec2.CreateRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateRouteTableOutput, error)
}
type AssociateRouteTableAPI interface {
	AssociateRouteTable(ctx context.Context,
		params *ec2.AssociateRouteTableInput,
		optFns ...func(*ec2.Options)) (*ec2.AssociateRouteTableOutput, error)
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

	client := GetEC2Client(remoteCfg)
	input := &ec2.DescribeInstancesInput{}

	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	for _, r := range result.Reservations {
		fmt.Println("Instance IDs:")
		for _, i := range r.Instances {
			if IsHyperdriveInstance(i) {
				fmt.Println("   " + GetHyperName(i))
			}
		}

		fmt.Println("")
	}

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
func GetVpcs(c context.Context, api DescribeVpcsAPI, input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return api.DescribeVpcs(c, input)
}
func MakeVpc(c context.Context, api CreateVpcAPI, input *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	return api.CreateVpc(c, input)
}
func GetSubnets(c context.Context, api DescribeSubnetsAPI, input *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return api.DescribeSubnets(c, input)
}
func MakeSubnet(c context.Context, api CreateSubnetAPI, input *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	return api.CreateSubnet(c, input)
}
func ChangeSubnet(c context.Context, api ModifySubnetAttributeAPI, input *ec2.ModifySubnetAttributeInput) (*ec2.ModifySubnetAttributeOutput, error) {
	return api.ModifySubnetAttribute(c, input)
}
func MakeInternetGateway(c context.Context, api CreateInternetGatewayAPI, input *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	return api.CreateInternetGateway(c, input)
}
func AttachInternetGateway(c context.Context, api AttachInternetGatewayAPI, input *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	return api.AttachInternetGateway(c, input)
}
func MakeRouteTable(c context.Context, api CreateRouteTableAPI, input *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	return api.CreateRouteTable(c, input)
}
func AddRouteTable(c context.Context, api AssociateRouteTableAPI, input *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	return api.AssociateRouteTable(c, input)
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
func MakeSecurityGroupPermissions(c context.Context, api AddSecurityGroupPermissionsAPI, input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return api.AuthorizeSecurityGroupIngress(c, input)
}
func GetKeyPairs(c context.Context, api DescribeKeyPairsAPI, input *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return api.DescribeKeyPairs(c, input)
}
func MakeKeyPair(c context.Context, api CreateKeyPairAPI, input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	return api.CreateKeyPair(c, input)
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
				FromPort:   aws.Int32(8080),
				ToPort:     aws.Int32(8080),
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

	minMaxCount := int32(1)
startupScript := `
#!/bin/bash
mkdir -p /tmp/hyperdrive
curl -fsSL https://github.com/gohypergiant/hyperdrive/releases/download/v0.0.0-troubleshoot/hyperdrive_0.0.0-troubleshoot_Linux_x86_64.tar.gz -o /tmp/hyperdrive/hyper.tar
tar -xvf /tmp/hyperdrive/hyper.tar -C /tmp/hyperdrive
mv /tmp/hyperdrive/hyper /usr/bin/hyper
hyper jupyter remoteHost
`
	ec2Input := &ec2.RunInstancesInput{
		ImageId:           aws.String(amiID),
		InstanceType:      types.InstanceType(*aws.String(ec2Type)),
		MinCount:          aws.Int32(minMaxCount),
		MaxCount:          aws.Int32(minMaxCount),
		SecurityGroupIds:  []string{securityGroupID},
		SubnetId:          aws.String(subnetID),
		KeyName:           aws.String(keyName),
		TagSpecifications: tagSpecification,
		UserData: aws.String(base64.StdEncoding.EncodeToString([]byte(startupScript))),
	}

	result, err := MakeInstance(context.TODO(), client, ec2Input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}

	if result.Instances[0].PublicIpAddress == nil {

		fmt.Println("Provisioned instance but cannot get publicIP")
		return
	}
	fmt.Print("EC2 instance provisioned. You can access via ssh by running:")
	fmt.Print("ssh -i " + keyName + ".pem ec2-user@" + *result.Instances[0].PublicIpAddress)
}
