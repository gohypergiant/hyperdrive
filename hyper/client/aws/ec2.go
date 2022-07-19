package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"time"

	config2 "github.com/gohypergiant/hyperdrive/hyper/services/config"

	"github.com/docker/distribution/uuid"
	hyperdriveTypes "github.com/gohypergiant/hyperdrive/hyper/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/client/ssh"
)

const HYPERDRIVE_TYPE_TAG string = "hyperdrive-type"
const HYPERDRIVE_NAME_TAG string = "hyperdrive-name"
const HYPERDRIVE_SECURITY_GROUP_NAME string = "-SecurityGroup"

type EC2Type int64

const (
	NotebookEC2 EC2Type = iota
	DeployEC2
)

// TODO, we should get this dynamically
const version string = "0.0.32"

func GetInstances(c context.Context, api hyperdriveTypes.EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}
func GetEC2Client(remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) *ec2.Client {

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
func ListServers(remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) {

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
func GetHyperdriveInstances(remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) ([]types.Instance, error) {

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
func CreateInternetGateway(client *ec2.Client, projectName string) string {
	input := &ec2.CreateInternetGatewayInput{
		TagSpecifications: getTagSpecification(projectName, types.ResourceTypeInternetGateway),
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

		inputMakeVPC := &ec2.CreateVpcInput{
			CidrBlock:         aws.String("10.0.0.0/16"),
			TagSpecifications: getTagSpecification(projectName, types.ResourceTypeVpc),
		}

		resultMakeVPC, err := MakeVpc(context.TODO(), client, inputMakeVPC)
		if err != nil {
			panic("error creating VPC," + err.Error())
		}
		vpcID = *resultMakeVPC.Vpc.VpcId

		internetGatewayID := CreateInternetGateway(client, projectName)
		fmt.Println("Internet Gateway ID:", internetGatewayID)

		AddInternetGateway(client, vpcID, internetGatewayID)

		inputMakeRouteTable := &ec2.CreateRouteTableInput{
			VpcId:             aws.String(vpcID),
			TagSpecifications: getTagSpecification(projectName, types.ResourceTypeRouteTable),
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

func getTagSpecification(projectName string, resourceType types.ResourceType) []types.TagSpecification {
	tagSpecification := []types.TagSpecification{
		{
			ResourceType: resourceType,
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
	return tagSpecification
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

	subnetMakeInput := &ec2.CreateSubnetInput{
		CidrBlock:         aws.String("10.0.1.0/24"),
		VpcId:             aws.String(vID),
		AvailabilityZone:  aws.String(region + "a"),
		TagSpecifications: getTagSpecification(projectName, types.ResourceTypeSubnet),
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

func getOrCreateSecurityGroup(client *ec2.Client, vID string, projectName string, httpPort int) string {
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

	scInput := &ec2.CreateSecurityGroupInput{
		GroupName:         aws.String(projectName + HYPERDRIVE_SECURITY_GROUP_NAME),
		Description:       aws.String("Security group for EC2 instances provisioned by hyper"),
		VpcId:             aws.String(vID),
		TagSpecifications: getTagSpecification(projectName, types.ResourceTypeSecurityGroup),
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
				FromPort:   aws.Int32(int32(httpPort)),
				ToPort:     aws.Int32(int32(httpPort)),
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
func WritePublicKey(client *ec2.Client, keyName string, projectName string, publicKeyBytes []byte) error {

	importKeyPairInput := &ec2.ImportKeyPairInput{
		KeyName:           aws.String(keyName),
		PublicKeyMaterial: publicKeyBytes,
		TagSpecifications: getTagSpecification(projectName, types.ResourceTypeKeyPair),
	}

	_, err := ImportKeyPair(context.TODO(), client, importKeyPairInput)
	if err != nil {
		return err
	}

	return nil
}
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
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

	if keyName == "" {

		keyName = projectName
		sshFolderPath := path.Join(UserHomeDir(), "/.ssh")
		privateKeyPath := path.Join(sshFolderPath, fmt.Sprintf("/%s", keyName))

		if _, err = os.Stat(sshFolderPath); os.IsNotExist(err) {
			os.Mkdir(sshFolderPath, os.FileMode(ssh.SSH_FOLDER_FILE_MODE))
		}

		var publicKeyBytes, privateKeyBytes []byte
		originalDir, err := os.Getwd()
		if err != nil {
			panic("error changing working directory: " + err.Error())
		}
		os.Chdir(sshFolderPath)

		if _, err = os.Stat(privateKeyPath); os.IsNotExist(err) {
			privateKeyBytes, publicKeyBytes = ssh.CreateRSAKeyPair(keyName)
			err = ssh.WriteKey(privateKeyPath, privateKeyBytes, fs.FileMode(ssh.PRIVATE_KEY_FILE_MODE))
			if err != nil {
				panic("error writing private key " + err.Error())
			}

			err = ssh.AddKeySshAgent(privateKeyPath)
			if err != nil {
				fmt.Println("ssh-agent not available")

				_, err = os.Stat(ssh.DEFAULT_KEY)
				if err != nil {
					fmt.Println("Writing key to default key value: id_rsa")
					os.Rename(keyName, ssh.DEFAULT_KEY)
					keyName = ssh.DEFAULT_KEY
				}
			}

		} else {
			publicKeyBytes = ssh.GetPublicKeyBytes(keyName)
		}
		os.Chdir(originalDir)

		err = WritePublicKey(client, keyName, projectName, publicKeyBytes)
		if err != nil {
			panic("error importing Key Pair," + err.Error())
		}
	}
	return keyName
}
func IsStudyInstance(i types.Instance, studyName string) bool {
	for _, t := range i.Tags {
		if *t.Key == HYPERDRIVE_NAME_TAG && *t.Value == studyName {
			return true
		}
	}
	return false
}
func GetInstanceForStudy(studyName string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) (types.Instance, error) {
	client := GetEC2Client(remoteCfg)

	ec2DescribeInput := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("instance-state-code"),
				Values: []string{"16"},
			},
		},
	}

	ec2DescribeResult, err := GetInstances(context.TODO(), client, ec2DescribeInput)
	if err != nil {
		return types.Instance{}, err
	}

	for _, r := range ec2DescribeResult.Reservations {
		for _, i := range r.Instances {
			if IsHyperdriveInstance(i) && IsStudyInstance(i, studyName) {
				return i, nil
			}
		}
	}
	return types.Instance{}, nil
}
func IsStructureEmpty(i types.Instance) bool {
	return reflect.DeepEqual(i, types.Instance{})
}
func StartJupyterEC2(manifestPath string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration, ec2Type string, amiID string, jupyterLaunchOptions hyperdriveTypes.JupyterLaunchOptions, syncOptions hyperdriveTypes.WorkspaceSyncOptions, serverType EC2Type) {
	startupScript := getEc2StartScript(version, jupyterLaunchOptions, syncOptions, remoteCfg, serverType)
	StartServer(manifestPath, remoteCfg, ec2Type, amiID, startupScript, jupyterLaunchOptions.HostPort, serverType)
}
func StartServer(manifestPath string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration, ec2Type string, amiID string, startupScript string, hostPort int, serverType EC2Type) {

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

	hyperInstance, err := GetInstanceForStudy(projectName, remoteCfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !IsStructureEmpty(hyperInstance) {
		message := fmt.Sprintf(`Hyper instance already running.
Instance size: %s
You can access jupyter lab at http://%s:8888/lab
If you want to change the instance size, stop the current running instance:

	hyper jupyter stop --remote=<REMOTE_PROFILE_NAME>
`, hyperInstance.InstanceType, *hyperInstance.PublicIpAddress)
		fmt.Println(message)
		return
	}

	vpcID, rtID := getOrCreateVPC(client, projectName)
	fmt.Println("VPC ID:", vpcID)
	fmt.Println("Route Table ID:", rtID)

	subnetID := getOrCreateSubnet(client, vpcID, remoteCfg.Region, projectName, rtID)
	fmt.Println("Subnet ID:", subnetID)

	securityGroupID := getOrCreateSecurityGroup(client, vpcID, projectName, hostPort)
	fmt.Println("Security group ID:", securityGroupID)

	keyName := getOrCreateKeyPair(client, projectName)
	fmt.Println("Key name:", keyName)

	minMaxCount := int32(1)
	ec2Input := &ec2.RunInstancesInput{
		ImageId:           aws.String(amiID),
		InstanceType:      types.InstanceType(*aws.String(ec2Type)),
		MinCount:          aws.Int32(minMaxCount),
		MaxCount:          aws.Int32(minMaxCount),
		SecurityGroupIds:  []string{securityGroupID},
		SubnetId:          aws.String(subnetID),
		KeyName:           aws.String(keyName),
		TagSpecifications: getTagSpecification(projectName, types.ResourceTypeInstance),
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
		ip, err = getInstanceIpAddress(*result.Instances[0].InstanceId, remoteCfg)
		if err != nil {

			fmt.Println("Provisioned instance but cannot get publicIP")
			os.Exit(1)
			return
		}

	}
	fmt.Println("")
	fmt.Println("EC2 instance provisioned. You can access via ssh by running:")

	if keyName == ssh.DEFAULT_KEY {
		fmt.Println("ssh ec2-user@" + *ip)
	} else {
		fmt.Println("ssh -i ~/.ssh/" + keyName + " ec2-user@" + *ip)
	}
	fmt.Println("")

	if serverType == NotebookEC2 {
		fmt.Println("In a few minutes, you should be able to access jupyter lab at http://" + *ip + ":8888/lab")
	} else if serverType == DeployEC2 {
		fmt.Println("Deploy completed, preditions avaliable at http://" + *ip + ":" + strconv.Itoa(hostPort))
	}

}

func getEc2StartScript(version string, jupyterLaunchOptions hyperdriveTypes.JupyterLaunchOptions, syncOptions hyperdriveTypes.WorkspaceSyncOptions, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration, serverType EC2Type) string {

	if syncOptions.S3Config.Profile != "" {

		namedProfileConfig := config2.GetNamedProfileConfig(syncOptions.S3Config.Profile)
		syncOptions.S3Config.AccessKey = namedProfileConfig.AccessKey
		syncOptions.S3Config.Secret = namedProfileConfig.Secret
		syncOptions.S3Config.Token = namedProfileConfig.Token
	}

	if serverType == NotebookEC2 {
		syncParameters := fmt.Sprintf("--s3AccessKey %s --s3Secret %s --s3Token %s --s3Region %s --s3BucketName %s -n %s", syncOptions.S3Config.AccessKey, syncOptions.S3Config.Secret, syncOptions.S3Config.Token, syncOptions.S3Config.Region, syncOptions.S3Config.BucketName, syncOptions.StudyName)
		syncCommand := fmt.Sprintf("hyper workspace sync %s -w", syncParameters)
		pullCommand := fmt.Sprintf("hyper workspace pull %s", syncParameters)
		s3Parameters := fmt.Sprintf("--s3AccessKey %s --s3AccessSecret %s --s3Region %s", remoteCfg.AccessKey, remoteCfg.Secret, remoteCfg.Region)

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
sudo -u ec2-user %s
sudo -u ec2-user nohup %s &
chown -R ec2-user:ec2-user .
sudo -u ec2-user bash -c 'hyper jupyter remoteHost --hostPort %d --apiKey %s %s &'
`, version, version, pullCommand, syncCommand, jupyterLaunchOptions.HostPort, jupyterLaunchOptions.APIKey, s3Parameters)

		return startupScript

	} else if serverType == DeployEC2 {
		syncParameters := fmt.Sprintf("--s3AccessKey %s --s3Secret %s --s3Token %s --s3Region %s --s3BucketName %s -n %s", syncOptions.S3Config.AccessKey, syncOptions.S3Config.Secret, syncOptions.S3Config.Token, syncOptions.S3Config.Region, syncOptions.S3Config.BucketName, syncOptions.StudyName)
		packCommand := fmt.Sprintf("hyper workspace pack %s", syncParameters)
		runParameters := fmt.Sprintf("--hyperpackagePath %s.hyperpack.zip --hostPort %d", syncOptions.StudyName, jupyterLaunchOptions.HostPort)
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
sudo -u ec2-user %s
chown -R ec2-user:ec2-user .
sudo -u ec2-user bash -c 'hyper hyperpackage run %s &'
`, version, version, packCommand, runParameters)

		return startupScript
	} else {
		fmt.Println("EC2 server type not implemented")
	}
	return ""
}
func getInstanceIpAddress(instanceId string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) (*string, error) {

	instances, err := GetHyperdriveInstances(remoteCfg)
	if err != nil {
		return nil, err
	}

	for _, i := range instances {
		if *i.InstanceId == instanceId {
			return i.PublicIpAddress, nil
		}
	}
	return nil, errors.New("Could not find public IP for instance")

}

//TODO: Refactor this in to a series of smaller, well-named functions for readability
func StopServer(manifestPath string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration) {
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

		sshFolderPath := path.Join(UserHomeDir(), "/.ssh")

		if keyName != ssh.DEFAULT_KEY {
			originalDir, err := os.Getwd()
			if err != nil {
				panic("error changing working directory: " + err.Error())
			}
			os.Chdir(sshFolderPath)

			_, err = os.Stat(keyName)
			if err == nil {
				_ = ssh.RemoveKeySshAgent(keyName)
				fmt.Println("Delete your key ~/.ssh/" + keyName + " if it is no longer in use")
			}
			os.Chdir(originalDir)
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
					panic("error disassociating Route Table," + err.Error())
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

func WriteFileToEC2(instanceIp string, remoteCfg hyperdriveTypes.EC2ComputeRemoteConfiguration, projectName string, filePath string) {

	keyName := projectName
	sshFolderPath := path.Join(UserHomeDir(), "/.ssh")
	privateKeyPath := path.Join(sshFolderPath, fmt.Sprintf("/%s", keyName))

	err := ssh.CopyToRemote("ec2-user", privateKeyPath, instanceIp, filePath, "./")
	if err != nil {
		privateKeyPath = path.Join(sshFolderPath, fmt.Sprintf("/%s", ssh.DEFAULT_KEY))
		err = ssh.CopyToRemote("ec2-user", privateKeyPath, instanceIp, filePath, "./")
		if err != nil {
			fmt.Println("Cannot copy file to EC2 server")
			os.Exit(1)
			return
		}
	}
}
