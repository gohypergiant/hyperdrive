package aws

import (
	"context"
	"fmt"
	"os"

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
				{
					Key:   aws.String("owner"),
					Value: aws.String("rodolfo-sekijima"),
				},
			},
		},
	}
	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: tagSpecification,
	}

	result, err := MakeInternetGateway(context.TODO(), client, input)
	if err != nil {
		panic("error when creating Internet Gateway: " + err.Error())
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
		panic("error when attaching the internet gateway to a VPC: " + err.Error())
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

func getOrCreateVPC(client *ec2.Client, projectName string) string {
	vpcInput := &ec2.DescribeVpcsInput{}
	result, err := GetVpcs(context.TODO(), client, vpcInput)
	if err != nil {
		panic("error when fetching VPCs, " + err.Error())
	}

	vID := GetVpcId(result, projectName)

	if vID == "" {
		fmt.Println("No Hyperdrive VPC found")
		fmt.Println("Creating Hyperdrive VPC")

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
					{
						Key:   aws.String("owner"),
						Value: aws.String("rodolfo-sekijima"),
					},
				},
			},
		}

		input := &ec2.CreateVpcInput{
			CidrBlock:         aws.String("10.0.0.0/16"),
			TagSpecifications: tagSpecification,
		}

		result, err := MakeVpc(context.TODO(), client, input)
		if err != nil {
			panic("error when creating a VPC " + err.Error())
		}
		vID = *result.Vpc.VpcId
		fmt.Println("Created VPC with ID:", vID)

		internetGatewayID := CreateInternetGateway(client, projectName)
		fmt.Println("Create Internet Gateway with ID:", internetGatewayID)

		AddInternetGateway(client, vID, internetGatewayID)

	}
	return vID
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

func getOrCreateSubnet(client *ec2.Client, vID string, region string, projectName string) string {

	snDescribeInput := &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vID},
			},
		},
	}
	snDescribeResult, err := GetSubnets(context.TODO(), client, snDescribeInput)
	if err != nil {
		panic("error when fetching subnets: " + err.Error())
	}

	snID := GetSubnetID(snDescribeResult, projectName)

	if snID == "" {
		fmt.Println("No Subnet found for VPC", vID)
		fmt.Println("Creating subnet for project", projectName)

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
					{
						Key:   aws.String("owner"),
						Value: aws.String("rodolfo-sekijima"),
					},
				},
			},
		}

		snInput := &ec2.CreateSubnetInput{
			CidrBlock:         aws.String("10.0.1.0/24"),
			VpcId:             aws.String(vID),
			AvailabilityZone:  aws.String(region + "a"),
			TagSpecifications: tagSpecification,
		}

		snResult, err := MakeSubnet(context.TODO(), client, snInput)
		if err != nil {
			panic("error when creating subnet: " + err.Error())
		}

		snID = *snResult.Subnet.SubnetId

		snChangeInput := &ec2.ModifySubnetAttributeInput{
			SubnetId: aws.String(snID),
			MapPublicIpOnLaunch: &types.AttributeBooleanValue{
				Value: aws.Bool(true),
			},
		}

		_, err = ChangeSubnet(context.TODO(), client, snChangeInput)
		if err != nil {
			panic("error when modifying subnet attribute: " + err.Error())
		}

	}

	return snID
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
	var scID string

	scDescribeInput := &ec2.DescribeSecurityGroupsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vID},
			},
		},
	}
	scDescribeResult, err := GetSecurityGroups(context.TODO(), client, scDescribeInput)
	if err != nil {
		panic("error when fetching Security Groups: " + err.Error())
	}

	scID = GetSecurityGroupId(scDescribeResult, projectName)

	if scID == "" {
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
					{
						Key:   aws.String("owner"),
						Value: aws.String("rodolfo-sekijima"),
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

		scMakeResult, err := MakeSecurityGroup(context.TODO(), client, scInput)
		if err != nil {
			panic("error when creating the security group: " + err.Error())
		}

		scID = *scMakeResult.GroupId

		scPermissionsInput := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: aws.String(scID),
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
			},
		}

		_, err = MakeSecurityGroupPermissions(context.TODO(), client, scPermissionsInput)
		if err != nil {
			panic("error when adding permissions to the security group: " + err.Error())
		}
	}
	return scID
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

	describeInput := &ec2.DescribeKeyPairsInput{
		IncludePublicKey: aws.Bool(true),
	}

	kpDescribeResult, err := GetKeyPairs(context.TODO(), client, describeInput)
	if err != nil {
		panic("error when fetching key pairs: " + err.Error())
	}
	keyName := getKeyPairName(kpDescribeResult, projectName)

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
					{
						Key:   aws.String("owner"),
						Value: aws.String("rodolfo-sekijima"),
					},
				},
			},
		}

		keyName = projectName + "-hyperdrive"
		makeInput := &ec2.CreateKeyPairInput{
			KeyName:           aws.String(keyName),
			TagSpecifications: tagSpecification,
		}

		makeResult, err := MakeKeyPair(context.TODO(), client, makeInput)
		if err != nil {
			panic("error when creating key pair: " + err.Error())
		}
		err = WriteKey(os.Getenv("HOME")+"/aws_ec2_key.pem", makeResult.KeyMaterial)
		if err != nil {
			fmt.Printf("Couldn't write key pair to file: %v", err)
		}

	}
	return keyName
}

func StartServer(manifestPath string, remoteCfg HyperConfig.EC2RemoteConfiguration, ec2Type string) {

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

	vpcID := getOrCreateVPC(client)
	fmt.Println("VPC ID:", vpcID)

	subnetID := getOrCreateSubnet(client, vpcID, remoteCfg.Region, projectName)
	fmt.Println("Subnet ID:", subnetID)

	securityGroupID := getOrCreateSecurityGroup(client, vpcID, projectName)
	fmt.Println("Security group ID:", securityGroupID)

	minMaxCount := int32(1)

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
				{
					Key:   aws.String("owner"),
					Value: aws.String("rodolfo-sekijima"),
				},
			},
		},
	}

	ec2Input := &ec2.RunInstancesInput{
		DryRun:            aws.Bool(true),
		ImageId:           aws.String("ami-e7527ed7"),
		InstanceType:      types.InstanceType(*aws.String(ec2Type)),
		MinCount:          &minMaxCount,
		MaxCount:          &minMaxCount,
		SecurityGroupIds:  []string{securityGroupID},
		SubnetId:          &subnetID,
		KeyName:           aws.String(keyName),
		TagSpecifications: tagSpecification,
	}

	_, err := MakeInstance(context.TODO(), client, ec2Input)
	if err != nil {
		fmt.Println("Got an error creating an instance:")
		fmt.Println(err)
		return
	}
}
