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
	"github.com/gohypergiant/hyperdrive/hyper/services/status"
	"github.com/spf13/cobra"
)

var (
	statusEndpointPort		string
	statusFilePath				string
)

var remoteStatusCmd = &cobra.Command{
	Use:   "remoteStatus",
	Short: "Summons an endpoint to obtain the status of the remote server",
	Run: func(cmd *cobra.Command, args []string) {
		status.StartEndpoint(statusEndpointPort, statusFilePath)
	},
}

var remoteStatusUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the status of the remote server",
	Run: func(cmd *cobra.Command, args []string) {
		status.UpdateStatus(args, statusFilePath)
	},
}

func init() {
	rootCmd.AddCommand(remoteStatusCmd)
	remoteStatusCmd.AddCommand(remoteStatusUpdateCmd)
	
	remoteStatusCmd.Flags().StringVar(&statusEndpointPort, "port", "3001", "Override the default remotestatus port")
	remoteStatusCmd.Flags().StringVar(&statusFilePath, "statusFile", "/statusfile.json", "Override the default statusfile path")

}
