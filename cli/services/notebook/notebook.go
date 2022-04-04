package notebook

import "github.com/gohypergiant/hyperdrive/cli/services/config"

type INotebookService interface {
	Start(image string, pullImage bool, jupyterBrowser bool)
	List()
	Stop(mountPointOrIdentifier string)
	UploadTrainingJobData()
	WaitForTrainingToComplete(timeout int)
	DownloadHyperpack()
}

func NotebookService(remoteName string, manifestPath string) INotebookService {

	if remoteName == "" {
		return LocalNotebookService{
			ManifestPath: manifestPath,
		}
	} else {
		return RemoteNotebookService{
			RemoteConfiguration: config.GetRemote(remoteName),
			ManifestPath:        manifestPath,
		}
	}
}
