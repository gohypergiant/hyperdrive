package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/manifoldco/promptui"
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

func getUsername() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide a username")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter your username",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result

}
func getUrl() string {

	prompt := promptui.Prompt{
		Label: "Enter the remote URL [default: Use Hypergiant hosted Hyperdrive]",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result

}
func getRemoteName() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide a name for the remote")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter remote name",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result

}

func getToken() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide an API token")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter your API token",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result

}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Config",
	Run: func(cmd *cobra.Command, args []string) {
		remoteName := getRemoteName()
		url := getUrl()
		username := getUsername()
		token := getToken()
		fmt.Printf("%s %s %s %s", remoteName, url, username, token)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(remotesCmd)

	//remote subcommands
	remotesCmd.AddCommand(remotesListCmd)
}
