package hyperpackage

type IHyperpackageService interface {
	Build(dockerfileSavePath string, imageTags []string)
	Run(imageTag string)
	BuildAndRun(dockerfileSavePath string, imageTags []string)
	Import(importModelFileName string, modelFlavor string, trainShape string)
	List()
	Stop(name string)
}

func HyperpackageService(hyperpackagePath string, manifestPath string) IHyperpackageService {
	return LocalHyperpackageService{
		HyperpackagePath: hyperpackagePath,
		ManifestPath:     manifestPath,
	}
}

const HYPERPACK_CONTAINER_PREFIX = "hyperpackage"
