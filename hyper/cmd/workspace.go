package cmd

import (
	"github.com/gohypergiant/hyperdrive/hyper/services/workspace"
	"github.com/spf13/cobra"
)

var watchSync = false

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Commands for interacting with remote workspaces",
}
var workspaceSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync",
	Run: func(cmd *cobra.Command, args []string) {
		//notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).List()
		println("syncing...")
		workspace.WorkspaceService("dev", manifestPath).Sync("", watchSync)
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceSyncCmd.Flags().BoolVarP(&watchSync, "watch", "w", false, "Run sync in watch mode")
	workspaceCmd.AddCommand(workspaceSyncCmd)
}
