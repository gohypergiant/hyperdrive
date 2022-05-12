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
	"os"
	"os/exec"
	"github.com/gohypergiant/hyperdrive/hyper/services/notebook"

	"github.com/spf13/cobra"
)

var (
	fetchTimeout int
)

// trainCmd represents the train command
var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Train a model",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚂choo choo🚂")
		notebookService := notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region)
		notebookService.UploadTrainingJobData()
		if RemoteName == "" {
			fmt.Println("TEST: should be local here")
			pmill, errPaper := exec.Command("papermill", "/home/jovyan/.executor/notebooks/executor-low-code.ipynb", "/home/jovyan/_jobs/threat_detection/outs.ipynb", "-p", "features", "data/object_data.json", "-p", "target", "data/label_data.csv", "-p", "job_name", "threat_detection", "-p", "study_yaml", "/home/jovyan/_jobs/threat_detection/_study.yaml").Output()
			if errPaper != nil {
				fmt.Println("Papermill Error.")
				os.Exit(1)
			}
			fmt.Println(pmill)
		}
		fmt.Println("Training data uploaded, to look for a completed hyperpackage, use the fetch subcommand.")
	},
}
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch resulting hyperpackage from training session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🥎🐕 Fetching")
		notebookService := notebook.NotebookService(RemoteName, manifestPath, s3AccessKey, s3AccessSecret, s3Region)
		notebookService.WaitForTrainingToComplete(fetchTimeout)
		notebookService.DownloadHyperpack()
	},
}

func init() {
	fetchCmd.Flags().IntVarP(&fetchTimeout, "fetchTimeout", "t", 3600, "Timeout in seconds to wait for training to complete (default 3600)")
	trainCmd.AddCommand(fetchCmd)
	trainCmd.Flags().StringVar(&image, "image", "pytorch", "Image to be used [huggingface-pytorch|huggingface-tensorflow|pytorch|spark|tensorflow|xgboost]")
	trainCmd.Flags().StringVar(&s3AccessKey, "s3AccessKey", "", "S3 Access Key to use")
	trainCmd.Flags().StringVar(&s3AccessSecret, "s3AccessSecret", "", "S3 Secret to use")
	trainCmd.Flags().StringVar(&s3Region, "s3Region", "", "S3 Region")
	rootCmd.AddCommand(trainCmd)
}
