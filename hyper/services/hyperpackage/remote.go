package hyperpackage

import (
	"fmt"

	"github.com/gohypergiant/hyperdrive/hyper/client/aws"
	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
)

type RemoteHyperpackageService struct {
	HyperpackagePath    string
	ManifestPath        string
	RemoteConfiguration types.ComputeRemoteConfiguration
}

func (s RemoteHyperpackageService) BuildAndRun(dockerfileSavePath string, imageTags []string, jupyterOptions types.JupyterLaunchOptions, ec2Options types.EC2StartOptions, syncOptions types.WorkspaceSyncOptions) {
	studyName := manifest.GetName(s.ManifestPath)
	if len(imageTags) == 0 {
		imageTags = []string{fmt.Sprintf("%s:latest", studyName)}
	}

	if s.RemoteConfiguration.Type == types.Firefly {
		fmt.Println("Firefly does not support deployment of hyperpackage")
	} else if s.RemoteConfiguration.Type == types.EC2 {
		if jupyterOptions.S3AwsProfile != "" {
			fmt.Printf("Using AWS named profile '%s' to retrieve AWS creds\n", jupyterOptions.S3AwsProfile)
			namedProfileConfig := config.GetNamedProfileConfig(jupyterOptions.S3AwsProfile)
			s.RemoteConfiguration.EC2Configuration.AccessKey = namedProfileConfig.AccessKey
			s.RemoteConfiguration.EC2Configuration.Secret = namedProfileConfig.Secret
			s.RemoteConfiguration.EC2Configuration.Region = namedProfileConfig.Region
			s.RemoteConfiguration.EC2Configuration.Token = namedProfileConfig.Token
		}
		jupyterOptions.HostPort = 8888
		aws.StartJupyterEC2(s.ManifestPath, s.RemoteConfiguration.EC2Configuration, ec2Options.InstanceType, ec2Options.AmiId, jupyterOptions, syncOptions, aws.DeployEC2)
	} else {
		fmt.Println("Not Implemented")
	}
}
func (s RemoteHyperpackageService) Build(dockerfileSavePath string, imageTags []string) {
}
func (s RemoteHyperpackageService) Run(imageTag string, hostPort string) {
}
func (s RemoteHyperpackageService) Import(importModelFileName string, modelFlavor string, trainShape string) {
}
func (s RemoteHyperpackageService) List()            {}
func (s RemoteHyperpackageService) Stop(name string) {}
