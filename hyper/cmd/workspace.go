package cmd

import (
	"fmt"

	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/gohypergiant/hyperdrive/hyper/services/workspace"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/spf13/cobra"
)

var (
	watchSync           = false
	localWorkspacePath  string
	workspaceRemoteName string
	workspaceS3Token    string
	studyName           string
	remotePackPath      string
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
		workspaceSyncOptions := getWorkspaceSyncOptions()
		workspace.WorkspaceService(workspaceRemoteName, manifestPath, workspaceSyncOptions.S3Config).Sync(localWorkspacePath, watchSync, studyName)
	},
}
var workspacePullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull",
	Run: func(cmd *cobra.Command, args []string) {
		//notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).List()
		workspaceSyncOptions := getWorkspaceSyncOptions()
		workspace.WorkspaceService(workspaceRemoteName, manifestPath, workspaceSyncOptions.S3Config).Pull(localWorkspacePath, studyName)
	},
}
var workspacePackCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack",
	Run: func(cmd *cobra.Command, args []string) {
		//notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).List()
		workspaceSyncOptions := getWorkspaceSyncOptions()
		workspace.WorkspaceService(workspaceRemoteName, manifestPath, workspaceSyncOptions.S3Config).Pack(studyName, remotePackPath)
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

func getWorkspaceSyncOptions() types.WorkspaceSyncOptions {
	workpaceSyncOptions := types.WorkspaceSyncOptions{}
	if studyName == "" {
		studyName = notebook.GetNotebookName(manifestPath)
	}
	workpaceSyncOptions.StudyName = studyName

	if workspaceRemoteName != "" {
		remoteConfig := config.GetWorkspacePersistenceRemote(workspaceRemoteName)
		workpaceSyncOptions.S3Config = remoteConfig.S3Configuration
	}

	if workspaceS3IsManuallySpecified() {

		workpaceSyncOptions.S3Config = types.S3WorkspacePersistenceRemoteConfiguration{
			Secret:     workspaceS3Secret,
			AccessKey:  workspaceS3AccessKey,
			Token:      workspaceS3Token,
			Profile:    workspaceS3Profile,
			BucketName: workspaceS3BucketName,
			Region:     workspaceS3Region,
		}

	} else {
		fmt.Println("Warning: workspace sync not configured")
	}
	return workpaceSyncOptions
}
func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceSyncCmd)
	workspaceCmd.AddCommand(workspacePullCmd)
	workspaceCmd.AddCommand(workspacePackCmd)

	workspaceSyncCmd.Flags().BoolVarP(&watchSync, "watch", "w", false, "Run sync in watch mode")
	workspaceSyncCmd.Flags().StringVarP(&localWorkspacePath, "localWorkspacePath", "l", "", "Local workspace path to sync")
	workspacePullCmd.Flags().StringVarP(&localWorkspacePath, "localWorkspacePath", "l", "", "Local workspace path to sync")
	workspacePackCmd.Flags().StringVarP(&remotePackPath, "remotePackPath", "", "", "Path to pack zip file")

	workspaceCmd.PersistentFlags().StringVarP(&workspaceRemoteName, "remote", "r", "", "name of the workspace remote to use for syncing")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Profile, "s3Profile", "", "Named AWS profile to use (from ~/.aws/config) [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3AccessKey, "s3AccessKey", "", "AWS Access Key for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Secret, "s3Secret", "", "AWS Secret for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Token, "s3Token", "", "AWS Token for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3Region, "s3Region", "", "AWS Region for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVar(&workspaceS3BucketName, "s3BucketName", "", "Bucket name for accessing S3 buckets [Overrides workspaceRemote]")
	workspaceCmd.PersistentFlags().StringVarP(&studyName, "studyName", "n", "", "Bucket name for accessing S3 buckets [Overrides workspaceRemote]")
}
