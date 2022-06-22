package notebook

import (
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
)

func NotebookService(remoteName string, manifestPath string, s3AccessKey string, s3AccessSecret string, s3Region string) types.INotebookService {

	s3Creds := types.S3Credentials{
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
			RemoteConfiguration: config.GetComputeRemote(remoteName),
			ManifestPath:        manifestPath,
		}
	}
}
