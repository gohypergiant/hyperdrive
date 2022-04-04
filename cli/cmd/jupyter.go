/*
Copyright © 2022 Hypergiant, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/gohypergiant/hyperdrive/cli/services/notebook"
	"github.com/spf13/cobra"
)

var (
	id             string
	image          string
	jupyterBrowser bool
	mountPoint     string
	pullImage      bool
	repoTag        string
	publicPort     uint16
)

// jupyterCmd represents the jupyter command
var jupyterCmd = &cobra.Command{
	Use:   "jupyter",
	Short: "Run a local jupyter server",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath).Start(image, pullImage, jupyterBrowser)
	},
}

var jupyterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running local jupyter servers",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath).List()
	},
}

var jupyterStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop and remove a currently running local jupyter server",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath).Stop(mountPoint)
	},
}

func init() {
	rootCmd.AddCommand(jupyterCmd)
	jupyterCmd.AddCommand(jupyterListCmd)
	jupyterCmd.AddCommand(jupyterStopCmd)

	jupyterCmd.Flags().BoolVarP(&jupyterBrowser, "browser", "", false, "Open jupyter in a browser after launching")
	jupyterCmd.Flags().BoolVarP(&pullImage, "pull", "", false, "Pull latest image before running")
	jupyterCmd.Flags().StringVar(&image, "image", "pytorch", "Image to be used [huggingface-pytorch|huggingface-tensorflow|pytorch|spark|tensorflow|xgboost]")
	jupyterStopCmd.Flags().StringVar(&mountPoint, "mountPoint", "", "Mount Point of Jupyter Server to be stopped")
}
