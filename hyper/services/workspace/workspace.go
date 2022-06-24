package workspace

import (
	"fmt"
	config2 "github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"os"
)

func WorkspaceService(remoteName string, manifestPath string) types.IWorkspaceService {
	config := config2.GetWorkspacePersistenceRemote(remoteName)
	if config.Type == types.S3 {
		return S3WorkspaceService{
			S3Configuration: config.S3Configuration,
			ManifestPath:    manifestPath,
		}
	}
	fmt.Println("invalid workspace remote specified")
	os.Exit(1)
	return nil
}
