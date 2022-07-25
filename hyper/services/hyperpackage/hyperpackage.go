package hyperpackage

import (
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
)

type IHyperpackageService interface {
	Build(dockerfileSavePath string, imageTags []string)
	Run(imageTag string, dockerOptions types.DockerOptions)
	BuildAndRun(dockerfileSavePath string, imageTags []string, jupyterOptions types.JupyterLaunchOptions, ec2Options types.EC2StartOptions, syncOptions types.WorkspaceSyncOptions, dockerOptions types.DockerOptions)
	Import(importModelFileName string, modelFlavor string, trainShape string)
	List()
	Stop(name string)
}

func HyperpackageService(hyperpackagePath string, manifestPath string, remoteName string) IHyperpackageService {

	if remoteName == "" {
		return LocalHyperpackageService{
			HyperpackagePath: hyperpackagePath,
			ManifestPath:     manifestPath,
		}
	} else {
		return RemoteHyperpackageService{
			HyperpackagePath:    hyperpackagePath,
			ManifestPath:        manifestPath,
			RemoteConfiguration: config.GetComputeRemote(remoteName),
		}
	}
}

const HYPERPACK_CONTAINER_PREFIX = "hyperpackage"
