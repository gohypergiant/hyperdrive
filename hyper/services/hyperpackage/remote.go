package hyperpackage

type RemoteHyperpackageService struct {
	HyperpackagePath string
	ManifestPath     string
}

func (s RemoteHyperpackageService) Build(dockerfileSavePath string, imageTags []string)       {}
func (s RemoteHyperpackageService) Run(imageTag string)                                       {}
func (s RemoteHyperpackageService) BuildAndRun(dockerfileSavePath string, imageTags []string) {}
func (s RemoteHyperpackageService) Import(importModelFileName string, modelFlavor string, trainShape string) {
}
func (s RemoteHyperpackageService) List()            {}
func (s RemoteHyperpackageService) Stop(name string) {}
