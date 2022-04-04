package cmd

import (
	"fmt"

	"github.com/gohypergiant/hyperdrive/cli/client/manifest"
	"github.com/gohypergiant/hyperdrive/cli/services/hyperpackage"
	"github.com/spf13/cobra"
)

var (
	hyperpackagePath          string
	hyperpackageContainerName string
	dockerfileSavePath        string
	imageTags                 []string
)

// runCmd represents the deploy command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a hyperpack",
	Run: func(cmd *cobra.Command, args []string) {
		studyName := manifest.GetName(manifestPath)
		if hyperpackagePath == "" {
			hyperpackagePath = fmt.Sprintf("./%s.hyperpack.zip", studyName)
		}
		fmt.Println("🚀 Deploying")
		if dockerfileSavePath == "" {
			dockerfileSavePath = fmt.Sprintf("./%s.Dockerfile", studyName)
		}
		fmt.Println("🚀 Building and Running Hyperpack")
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath).BuildAndRun(dockerfileSavePath, imageTags)
	},
}

// runCmd represents the deploy command
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
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath).Build(dockerfileSavePath, imageTags)
	},
}

// listCmd represents the deploy command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists hyperpackage containers that are currently running",
	Run: func(cmd *cobra.Command, args []string) {
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath).List()
	},
}

// stopCmd represents the deploy command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stops a hyperpackage container that is currently running",
	Run: func(cmd *cobra.Command, args []string) {
		studyName := manifest.GetName(manifestPath)
		if hyperpackageContainerName == "" {
			hyperpackageContainerName = studyName
		}
		hyperpackage.HyperpackageService(hyperpackagePath, manifestPath).Stop(hyperpackageContainerName)
	},
}
var hyperpackageCmd = &cobra.Command{
	Use:   "hyperpackage",
	Short: "",
}

func init() {
	rootCmd.AddCommand(hyperpackageCmd)
	hyperpackageCmd.AddCommand(runCmd)
	hyperpackageCmd.AddCommand(buildCmd)
	hyperpackageCmd.AddCommand(listCmd)
	stopCmd.Flags().StringVar(&hyperpackageContainerName, "hyperpackagePath", "", "name of container to stop")
	hyperpackageCmd.AddCommand(stopCmd)
	hyperpackageCmd.PersistentFlags().StringVarP(&hyperpackagePath, "hyperpackagePath", "p", "", "path to hyperpackage.zip")
	hyperpackageCmd.PersistentFlags().StringVarP(&dockerfileSavePath, "dockerfileSavePath", "o", "", "path to save Dockerfile")
	hyperpackageCmd.PersistentFlags().StringArrayVarP(&imageTags, "imageTags", "t", []string{}, "tag for resulting docker image")
}
