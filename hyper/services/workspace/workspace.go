package workspace

import (
	"fmt"
	config2 "github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"os"
)

func WorkspaceService(remoteName string, manifestPath string, workspaceS3Configuration types.S3WorkspacePersistenceRemoteConfiguration) types.IWorkspaceService {

	workspaceRemoteType := types.S3

	if (types.S3WorkspacePersistenceRemoteConfiguration{}) == workspaceS3Configuration {
		workspaceS3Configuration = config2.GetWorkspacePersistenceRemote(remoteName).S3Configuration
	}

	if workspaceRemoteType == types.S3 {
		return S3WorkspaceService{
			S3Configuration: workspaceS3Configuration,
			ManifestPath:    manifestPath,
		}
	}
	fmt.Println("invalid workspace remote specified")
	os.Exit(1)
	return nil
}
