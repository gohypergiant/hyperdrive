package hyperpackage

import "github.com/gohypergiant/hyperdrive/hyper/types"

type RemoteHyperpackageService struct {
	HyperpackagePath    string
	ManifestPath        string
	RemoteConfiguration types.ComputeRemoteConfiguration
}

func (s RemoteHyperpackageService) BuildAndRun(dockerfileSavePath string, imageTags []string) {}
func (s RemoteHyperpackageService) Build(dockerfileSavePath string, imageTags []string)       {}
func (s RemoteHyperpackageService) Run(imageTag string)                                       {}
func (s RemoteHyperpackageService) Import(importModelFileName string, modelFlavor string, trainShape string) {
}
func (s RemoteHyperpackageService) List()            {}
func (s RemoteHyperpackageService) Stop(name string) {}
