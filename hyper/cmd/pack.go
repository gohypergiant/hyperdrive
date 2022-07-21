package cmd

import (
	"fmt"
	"strconv"

	"github.com/gohypergiant/hyperdrive/hyper/client/manifest"
	"github.com/gohypergiant/hyperdrive/hyper/services/hyperpackage"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/spf13/cobra"
)

var (
	hyperpackagePath          string
	hyperpackageContainerName string
	dockerfileSavePath        string
	imageTags                 []string
	importModelFileName       string
	modelFlavor               string
	trainShape                string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a hyperpack",
	Run: func(cmd *cobra.Command, args []string) {
		var portInt int
		var err error

		studyName := manifest.GetName(manifestPath)
		if hyperpackagePath == "" {
			hyperpackagePath = fmt.Sprintf("./%s.hyperpack.zip", studyName)
		}
		fmt.Println("🚀 Deploying")
		if dockerfileSavePath == "" {
			dockerfileSavePath = fmt.Sprintf("./%s.Dockerfile", studyName)
		}
		if hostPort != "" {
			portInt, err = strconv.Atoi(hostPort)

			if err != nil {
				fmt.Println("--hostPort not a integer")
				return
			}
		}

		fmt.Println("🚀 Building and Running Hyperpack")
		hyperpackage.HyperpackageService(hyperpackagePath,
			manifestPath,
			RemoteName).BuildAndRun(
			dockerfileSavePath,
			imageTags,
			types.JupyterLaunchOptions{HostPort: portInt},
			types.EC2StartOptions{InstanceType: ec2InstanceType, AmiId: amiID},
			getWorkspaceSyncOptions())
	},
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "builds a hyperpack (but doesn't run it)",
	Run: func(cmd *cobra.Command, args []string) {
		studyName := manifest.GetName(manifestPath)
		if hyperpackagePath == "" {
			hyperpackagePath = fmt.Sprintf("./%s.hyperpack.zip", studyName)
		}
		if dockerfileSavePath == "" {
			dockerfileSavePath = fmt.Sprintf("./%s.Dockerfile", studyName)
		}
		fmt.Println("🚀 Building hyperpackage %s. Dockerfile will be saved to ", dockerfileSavePath)
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath, RemoteName).Build(dockerfileSavePath, imageTags)
	},
}

// importCmd represents the import model command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "imports a trained model",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Importing a trained model...")
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath, RemoteName).Import(importModelFileName, modelFlavor, trainShape)
		fmt.Println("Importing complete.")
	},
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists hyperpackage containers that are currently running",
	Run: func(cmd *cobra.Command, args []string) {
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath, RemoteName).List()
	},
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops a hyperpackage container that is currently running",
	Run: func(cmd *cobra.Command, args []string) {
		studyName := manifest.GetName(manifestPath)
		if hyperpackageContainerName == "" {
			hyperpackageContainerName = studyName
		}
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath, RemoteName).Stop(hyperpackageContainerName)
	},
}
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "",
}

func init() {
	rootCmd.AddCommand(packCmd)
	packCmd.AddCommand(runCmd)
	packCmd.AddCommand(buildCmd)
	importCmd.Flags().StringVar(&importModelFileName, "filename", "", "import model filename")
	importCmd.Flags().StringVar(&modelFlavor, "modelFlavor", "sklearn", "model flavor")
	importCmd.Flags().StringVar(&trainShape, "shape", "", "Training shape of data, specifically the number of columns.")
	packCmd.AddCommand(importCmd)
	packCmd.AddCommand(listCmd)
	stopCmd.Flags().StringVar(&hyperpackageContainerName, "hyperpackagePath", "", "name of container to stop")
	packCmd.PersistentFlags().StringVarP(&hyperpackagePath, "hyperpackagePath", "p", "", "path to hyperpackage.zip")
	packCmd.PersistentFlags().StringVarP(&dockerfileSavePath, "dockerfileSavePath", "o", "", "path to save Dockerfile")
	packCmd.PersistentFlags().StringArrayVarP(&imageTags, "imageTags", "t", []string{}, "tag for resulting docker image")
	packCmd.PersistentFlags().StringVarP(&workspaceRemoteName, "workspaceRemote", "r", "", "name of the jupyter remote to use for syncing")
	packCmd.PersistentFlags().StringVar(&workspaceS3Profile, "workspaceS3Profile", "", "Named AWS profile to use (from ~/.aws/config) [Overrides workspaceRemote]")
	packCmd.PersistentFlags().StringVar(&workspaceS3AccessKey, "workspaceS3AccessKey", "", "AWS Access Key for accessing S3 buckets [Overrides workspaceRemote]")
	packCmd.PersistentFlags().StringVar(&workspaceS3Secret, "workspaceS3Secret", "", "AWS Secret for accessing S3 buckets [Overrides workspaceRemote]")
	packCmd.PersistentFlags().StringVar(&workspaceS3Token, "workspaceS3Token", "", "AWS Token for accessing S3 buckets [Overrides workspaceRemote]")
	packCmd.PersistentFlags().StringVar(&workspaceS3Region, "workspaceS3Region", "", "AWS Region for accessing S3 buckets [Overrides workspaceRemote]")
	packCmd.PersistentFlags().StringVar(&workspaceS3BucketName, "workspaceS3BucketName", "", "Bucket name for accessing S3 buckets [Overrides workspaceRemote]")
	runCmd.PersistentFlags().StringVar(&ec2InstanceType, "ec2InstanceType", "", "The type of EC2 instance to be created")
	runCmd.PersistentFlags().StringVar(&amiID, "amiId", "", "The ID of the AMI")
	runCmd.PersistentFlags().StringVar(&hostPort, "hostPort", "", "Host port for container")
	packCmd.AddCommand(stopCmd)
}
