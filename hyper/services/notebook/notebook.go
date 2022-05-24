package notebook

import "github.com/gohypergiant/hyperdrive/hyper/services/config"

type INotebookService interface {
	Start(image string, pullImage bool, jupyterBrowser bool, ec2InstanceType string, amiID string)
	List()
	Stop(mountPointOrIdentifier string)
	UploadTrainingJobData()
	WaitForTrainingToComplete(timeout int)
	DownloadHyperpack()
}
type S3Credentials struct {
	AccessKey    string
	AccessSecret string
	Region       string
}

func NotebookService(remoteName string, manifestPath string, s3AccessKey string, s3AccessSecret string, s3Region string) INotebookService {

	s3Creds := S3Credentials{
		AccessKey:    s3AccessKey,
		AccessSecret: s3AccessSecret,
		Region:       s3Region,
	}
	if remoteName == "" {
		return LocalNotebookService{
			ManifestPath:  manifestPath,
			S3Credentials: s3Creds,
		}
	} else {
		return RemoteNotebookService{
			RemoteConfiguration: config.GetRemote(remoteName),
			ManifestPath:        manifestPath,
		}
	}
}
