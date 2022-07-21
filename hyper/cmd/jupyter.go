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
	"fmt"

	"github.com/gohypergiant/hyperdrive/hyper/client/cli"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var (
	image           string
	jupyterBrowser  bool
	mountPoint      string
	pullImage       bool
	requirements    bool
	s3AccessKey     string
	s3AccessSecret  string
	s3Region        string
	s3AwsProfile    string
	ec2InstanceType string
	amiID           string
	hostPort        string
	jupyterApiKey   string
)

func getPort(isRemote bool) int {
	defaultPort := "8888"
	defaultPortOpen := true
	if hostPort == "" && isRemote {
		hostPort = defaultPort
	} else if hostPort == "" && !isRemote {
		dockerClient := cli.NewDockerClient()
		nowRunningContainers, _ := dockerClient.ListAllRunningContainers()
		defaultPortUInt64, _ := strconv.ParseUint(defaultPort, 10, 64)
		defaultPortUInt16 := uint16(defaultPortUInt64)
		for _, runningContainer := range nowRunningContainers {
			if runningContainer.Ports[0].PublicPort == defaultPortUInt16 {
				defaultPortOpen = false
				fmt.Println("Default port 8888 is in use. Will assign a random port for the container.")
				break
			}
		}
		if defaultPortOpen {
			hostPort = defaultPort
		} else {
			fmt.Println("random port assign")
		}
	}
	port, err := strconv.Atoi(hostPort)
	if err != nil {
		log.Fatal("Couldn't parse port")
	}
	return port
}

// jupyterCmd represents the jupyter command
var jupyterCmd = &cobra.Command{
	Use:   "jupyter",
	Short: "Run a local jupyter server",
	Run: func(cmd *cobra.Command, args []string) {

		launchOptions := types.JupyterLaunchOptions{
			Flavor:        image,
			PullImage:     pullImage,
			LaunchBrowser: jupyterBrowser,
			Requirements:  requirements,
			RestartAlways: false,
			APIKey:        jupyterApiKey,
			S3AwsProfile:  s3AwsProfile,
			HostPort:      getPort(RemoteName != ""),
		}
		notebook.NotebookService(
			RemoteName,
			manifestPath,
			s3AccessKey,
			s3AccessSecret,
			s3Region).Start(
			launchOptions,
			types.EC2StartOptions{InstanceType: ec2InstanceType, AmiId: amiID},
			getWorkspaceSyncOptions(),
		)
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
var jupyterRemoteHost = &cobra.Command{
	Use:   "remoteHost",
	Short: "start server on remote host",
	Run: func(cmd *cobra.Command, args []string) {
		launchOptions := types.JupyterLaunchOptions{
			Flavor:        image,
			PullImage:     pullImage,
			LaunchBrowser: jupyterBrowser,
			Requirements:  requirements,
			HostPort:      getPort(true),
			RestartAlways: true,
			APIKey:        jupyterApiKey,
			S3AwsProfile:  s3AwsProfile,
		}
		notebook.NotebookService(
			RemoteName,
			manifestPath,
			s3AccessKey,
			s3AccessSecret,
			s3Region).Start(
			launchOptions,
			types.EC2StartOptions{InstanceType: ec2InstanceType, AmiId: amiID},
			getWorkspaceSyncOptions(),
		)
	},
}

func init() {
	rootCmd.AddCommand(jupyterCmd)
	jupyterCmd.AddCommand(jupyterListCmd)
	jupyterCmd.AddCommand(jupyterStopCmd)
	jupyterCmd.AddCommand(jupyterRemoteHost)

	jupyterCmd.Flags().BoolVarP(&jupyterBrowser, "browser", "", false, "Open jupyter in a browser after launching")
	jupyterCmd.PersistentFlags().BoolVarP(&pullImage, "pull", "", false, "Pull latest image before running")
	jupyterCmd.PersistentFlags().BoolVarP(&requirements, "requirements", "", false, "Install more packages from a requirements.txt file")
	jupyterCmd.PersistentFlags().StringVar(&image, "image", "pytorch", "Image to be used [huggingface-pytorch|huggingface-tensorflow|pytorch|spark|tensorflow|xgboost]")
	jupyterCmd.PersistentFlags().StringVar(&s3AccessKey, "s3AccessKey", "", "S3 Access Key to use")
	jupyterCmd.PersistentFlags().StringVar(&s3AccessSecret, "s3AccessSecret", "", "S3 Secret to use")
	jupyterCmd.PersistentFlags().StringVar(&s3Region, "s3Region", "", "S3 Region")
	jupyterCmd.PersistentFlags().StringVar(&s3AwsProfile, "s3AwsProfile", "", "Named AWS profile")
	jupyterCmd.PersistentFlags().StringVar(&ec2InstanceType, "ec2InstanceType", "", "The type of EC2 instance to be created")
	jupyterCmd.PersistentFlags().StringVar(&amiID, "amiId", "", "The ID of the AMI")
	jupyterCmd.PersistentFlags().StringVar(&jupyterApiKey, "apiKey", "", "API key to use for the jupyter instance")
	jupyterCmd.PersistentFlags().StringVar(&hostPort, "hostPort", "", "Host port for container")
	jupyterCmd.PersistentFlags().StringVarP(&workspaceRemoteName, "workspaceRemote", "r", "", "name of the jupyter remote to use for syncing")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3Profile, "workspaceS3Profile", "", "Named AWS profile to use (from ~/.aws/config) [Overrides workspaceRemote]")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3AccessKey, "workspaceS3AccessKey", "", "AWS Access Key for accessing S3 buckets [Overrides workspaceRemote]")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3Secret, "workspaceS3Secret", "", "AWS Secret for accessing S3 buckets [Overrides workspaceRemote]")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3Token, "workspaceS3Token", "", "AWS Token for accessing S3 buckets [Overrides workspaceRemote]")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3Region, "workspaceS3Region", "", "AWS Region for accessing S3 buckets [Overrides workspaceRemote]")
	jupyterCmd.PersistentFlags().StringVar(&workspaceS3BucketName, "workspaceS3BucketName", "", "Bucket name for accessing S3 buckets [Overrides workspaceRemote]")
	jupyterStopCmd.Flags().StringVar(&mountPoint, "mountPoint", "", "Mount Point of Jupyter Server to be stopped")
}
