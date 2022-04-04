package cmd

import (
	"fmt"

	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/spf13/cobra"
)

var remotesCmd = &cobra.Command{
	Short: "Interact with firefly remotes",
	Use:   "remote",
}
var remotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Remotes",
	Run: func(cmd *cobra.Command, args []string) {
		remotesMap := config.GetRemotes()
		for name, config := range remotesMap {
			fmt.Println("remote: ", name)
			fmt.Println("    url: ", config.Url)
		}
	},
}

// trainCmd represents the train command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config",
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(remotesCmd)

	//remote subcommands
	remotesCmd.AddCommand(remotesListCmd)
}
