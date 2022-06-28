package workspace

import (
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/client/aws"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"os"
	"time"
)

type S3WorkspaceService struct {
	ManifestPath    string
	S3Configuration types.S3WorkspacePersistenceRemoteConfiguration
}

func (s S3WorkspaceService) Sync(localPath string, watch bool, studyName string) {

	if studyName == "" {
		studyName = notebook.GetNotebookName(s.ManifestPath)
	}
	if localPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		localPath = cwd
	}
	if watch {
		s.watchSync(localPath, studyName)

	} else {
		s.syncOnce(localPath, studyName)
	}
}
func (s S3WorkspaceService) syncOnce(localPath string, studyName string) {
	remotePath := s.GetS3Url(studyName)
	fmt.Println(remotePath)

	fmt.Println("syncing local to remote")
	aws.SyncDirectory(s.S3Configuration, localPath, remotePath)
	fmt.Println("syncing remote to local")
	aws.SyncDirectory(s.S3Configuration, remotePath, localPath)

}
func (s S3WorkspaceService) watchSync(localPath string, studyName string) {
	s.syncOnce(localPath, studyName)
	for range time.Tick(time.Second * 10) {
		s.syncOnce(localPath, studyName)
	}
}

func (s S3WorkspaceService) GetS3Url(studyName string) string {

	return fmt.Sprintf("s3://%s/%s", s.S3Configuration.BucketName, studyName)
}
