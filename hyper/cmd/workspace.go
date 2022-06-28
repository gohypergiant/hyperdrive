package cmd

import (
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/services/workspace"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/spf13/cobra"
)

var (
	watchSync           = false
	localWorkspacePath  string
	workspaceRemoteName string
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Commands for interacting with remote workspaces",
}
var workspaceSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync",
	Run: func(cmd *cobra.Command, args []string) {
		//notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).List()
		s3Config := types.S3WorkspacePersistenceRemoteConfiguration{}
		if workspaceS3IsManuallySpecified() {
			fmt.Println("Using s3 credentials provided as arguments for workspace sync")
			s3Config.Secret = workspaceS3Secret
			s3Config.AccessKey = workspaceS3AccessKey
			s3Config.Profile = workspaceS3Profile
			s3Config.BucketName = workspaceS3BucketName
			s3Config.Region = workspaceS3Region
		} else {

			fmt.Printf("Using configured remote %s for workspace sync \n")
		}
		workspace.WorkspaceService(workspaceRemoteName, manifestPath, s3Config).Sync(localWorkspacePath, watchSync)
	},
}

func workspaceS3IsManuallySpecified() bool {
	if workspaceRemoteName != "" {
		return false
	}
	if workspaceS3BucketName == "" {
		return false
	}
	if workspaceS3Region == "" {
		return false
	}
	if workspaceS3Profile != "" {
		return true
	}
	if workspaceS3AccessKey != "" && workspaceS3Secret != "" {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceSyncCmd)

	workspaceSyncCmd.Flags().BoolVarP(&watchSync, "watch", "w", false, "Run sync in watch mode")
	workspaceSyncCmd.Flags().StringVarP(&localWorkspacePath, "localWorkspacePath", "l", "", "Local workspace path to sync")

	workspaceCmd.PersistentFlags().StringVarP(&workspaceRemoteName, "remote", "r", "", "name of the workspace remote to use for syncing")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Profile, "s3Profile", "", "Named AWS profile to use (from ~/.aws/config) [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3AccessKey, "s3AccessKey", "", "AWS Access Key for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Secret, "s3Secret", "", "AWS Secret for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Region, "s3Region", "", "AWS Region for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3BucketName, "s3BucketName", "", "Bucket name for accessing S3 buckets [Overrides workspaceRemote]")
}
