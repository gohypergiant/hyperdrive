package workspace

import (
	"fmt"
	"os"
	"time"

	"github.com/gohypergiant/hyperdrive/hyper/client/aws"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/rogpeppe/go-internal/lockedfile"
)

const LOCKFILE_NAME string = "hyperdrive-workspace.lock"

type S3WorkspaceService struct {
	ManifestPath    string
	S3Configuration types.S3WorkspacePersistenceRemoteConfiguration
}

func (s S3WorkspaceService) Pull(localPath string, studyName string) {

	studyName, localPath, err := s.determinePathAndName(localPath, studyName)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.pull(localPath, studyName)
}

func (s S3WorkspaceService) determinePathAndName(localPath string, studyName string) (string, string, error) {
	if studyName == "" {
		studyName = notebook.GetNotebookName(s.ManifestPath)
	}
	if localPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
		localPath = cwd
	}
	return studyName, localPath, nil
}
func (s S3WorkspaceService) Sync(localPath string, watch bool, studyName string) {

	studyName, localPath, err := s.determinePathAndName(localPath, studyName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if watch {
		s.watchSync(localPath, studyName)

	} else {
		s.syncOnce(localPath, studyName)
	}
}
func (s S3WorkspaceService) pull(localPath string, studyName string) {
	remotePath := s.GetS3Url(studyName)
	fmt.Println(remotePath)

	fmt.Println("Pulling from remote")
	aws.SyncDirectory(s.S3Configuration, remotePath, localPath)
}
func (s S3WorkspaceService) syncOnce(localPath string, studyName string) {
	remotePath := s.GetS3Url(studyName)
	fmt.Println(remotePath)

	lockfile, err := lockedfile.Create(LOCKFILE_NAME)
	if err != nil {
		fmt.Println("could not aquire lock to sync")
		return
	}
	fmt.Println("syncing local to remote")
	aws.SyncDirectory(s.S3Configuration, localPath, remotePath)
	fmt.Println("syncing remote to local")
	aws.SyncDirectory(s.S3Configuration, remotePath, localPath)
	lockfile.Close()
	os.Remove(LOCKFILE_NAME)
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

func (s S3WorkspaceService) Pack(studyName string, packFile string) {
	var packPath string = packFile
	if packFile == "" {
		packPath = studyName + "/_jobs/" + studyName + "/" + studyName + ".hyperpack.zip"
	}

	err := aws.DownloadObject(s.S3Configuration, studyName+".hyperpack.zip", packPath)

	if err != nil {
		fmt.Println("Error pulling from S3: ", err)
	}
}
