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
			fmt.Println("    url: ", config.FireflyConfiguration.Url)
		}
	},
}

// trainCmd represents the train command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config",
}

var remoteName string
var url string
var username string
var token string
var remoteType config.RemoteType

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
	//TODO: Set default value to hosted backend

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
func getFireflyConfig() RemoteConfiguration {

		if url == "" {
			url = getUrl()
		}
		if username == "" {
			username = getUsername()
		}
		if token == "" {
			token = getToken()
		}
		return config.RemoteConfiguration{
			Type:                 config.Firefly,
			FireflyConfiguration: config.FireflyRemoteConfiguration{Url: url, Username: username, HubToken: token},
		}

}
func getConfigType() config.RemoteType {
	return config.Firefly
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Config",
	Run: func(cmd *cobra.Command, args []string) {
		var remoteConfig config.RemoteConfiguration 
		if remoteName == "" {
			remoteName = getRemoteName()
		}
		if remoteType == "" {
			remoteName = getConfigType()
		}
		switch remoteType {
		case config.Firefly:
		default:
			remoteConfig = getFireflyConfig()
			break;
		}

		config.UpdateRemote(remoteName, remoteConfig)
		fmt.Printf("Added %s remote at %s", remoteName, url)
	},
}

func init() {
	initCmd.Flags().StringVar(&remoteName, "remoteName", "", "Name of the remote for the config")
	initCmd.Flags().StringVar(&remoteType, "remoteType")
	initCmd.Flags().StringVar(&url, "url", "", "URL to the remote")
	initCmd.Flags().StringVar(&username, "username", "", "Username for the remote")
	initCmd.Flags().StringVar(&token, "token", "", "token for the remote")
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(remotesCmd)

	//remote subcommands
	remotesCmd.AddCommand(remotesListCmd)
}
