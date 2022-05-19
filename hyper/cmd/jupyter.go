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
	"github.com/gohypergiant/hyperdrive/hyper/client/aws"
	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/spf13/cobra"
)

var (
	id              string
	image           string
	jupyterBrowser  bool
	mountPoint      string
	pullImage       bool
	repoTag         string
	publicPort      uint16
	s3AccessKey     string
	s3AccessSecret  string
	s3Region        string
	ec2InstanceType string
	remoteProfile   string
)

// jupyterCmd represents the jupyter command
var jupyterCmd = &cobra.Command{
	Use:   "jupyter",
	Short: "Run a local jupyter server",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).Start(image, pullImage, jupyterBrowser)
	},
}

var jupyterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running local jupyter servers",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).List()
	},
}

var jupyterStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop and remove a currently running local jupyter server",
	Run: func(cmd *cobra.Command, args []string) {
		notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region).Stop(mountPoint)
	},
}

var jupyterCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create EC2 instance to run a remote jupyter server",
	Run: func(cmd *cobra.Command, args []string) {
		remoteConfiguration := config.GetRemote(remoteProfile)
		aws.StartServer(manifestPath, remoteConfiguration.EC2Configuration, ec2InstanceType)
	},
}

func init() {
	rootCmd.AddCommand(jupyterCmd)
	jupyterCmd.AddCommand(jupyterListCmd)
	jupyterCmd.AddCommand(jupyterStopCmd)
	jupyterCmd.AddCommand(jupyterCreateCmd)

	jupyterCmd.Flags().BoolVarP(&jupyterBrowser, "browser", "", false, "Open jupyter in a browser after launching")
	jupyterCmd.Flags().BoolVarP(&pullImage, "pull", "", false, "Pull latest image before running")
	jupyterCmd.Flags().StringVar(&image, "image", "pytorch", "Image to be used [huggingface-pytorch|huggingface-tensorflow|pytorch|spark|tensorflow|xgboost]")
	jupyterCmd.Flags().StringVar(&s3AccessKey, "s3AccessKey", "", "S3 Access Key to use")
	jupyterCmd.Flags().StringVar(&s3AccessSecret, "s3AccessSecret", "", "S3 Secret to use")
	jupyterCmd.Flags().StringVar(&s3Region, "s3Region", "", "S3 Region")
	jupyterStopCmd.Flags().StringVar(&mountPoint, "mountPoint", "", "Mount Point of Jupyter Server to be stopped")
	jupyterCreateCmd.Flags().StringVar(&ec2InstanceType, "ec2InstanceType", "", "The type of EC2 instance to be created")
	jupyterCreateCmd.Flags().StringVar(&remoteProfile, "remoteProfile", "", "Name of the remote profiile to be used to create a EC2 instance")
}
