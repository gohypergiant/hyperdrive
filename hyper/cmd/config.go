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
var fireflyUrl string
var fireflyUsername string
var fireflyToken string
var remoteTypeInput string

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
func getFireflyConfig() config.RemoteConfiguration {

	if fireflyUrl == "" {
		fireflyUrl = getUrl()
	}
	if fireflyUsername == "" {
		fireflyUsername = getUsername()
	}
	if fireflyToken == "" {
		fireflyToken = getToken()
	}
	return config.RemoteConfiguration{
		Type:                 config.Firefly,
		FireflyConfiguration: config.FireflyRemoteConfiguration{Url: fireflyUrl, Username: fireflyUsername, HubToken: fireflyToken},
	}

}
func getConfigType() config.RemoteType {
	if remoteTypeInput == "" {
		prompt := promptui.Select{
			Label: "Choose a remote type",
			Items: config.ValidRemoteTypes,
		}
		_, remoteTypeInput, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return config.RemoteType(remoteTypeInput)
	}
	if remoteTypeInput == string(config.Firefly) {
		return config.Firefly
	}
	if remoteTypeInput == string(config.EC2) {
		return config.EC2
	}

	fmt.Println("Invalid or unsupported remote type")
	os.Exit(1)
	return (config.RemoteType(remoteTypeInput))
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Config",
	Run: func(cmd *cobra.Command, args []string) {
		var remoteConfig config.RemoteConfiguration
		if remoteName == "" {
			remoteName = getRemoteName()
		}
		remoteType := getConfigType()
		switch remoteType {
		case config.Firefly:
			fallthrough
		default:
			remoteConfig = getFireflyConfig()
			fmt.Printf("Adding %s remote at %s", remoteName, fireflyUrl)
			break
		}

		config.UpdateRemote(remoteName, remoteConfig)
	},
}

func init() {
	initCmd.Flags().StringVar(&remoteName, "remoteName", "", "Name of the remote for the config")
	initCmd.Flags().StringVarP(&remoteTypeInput, "remoteType", "r", "", "Remote type [firefly|ec2]")
	initCmd.Flags().StringVar(&fireflyUrl, "fireflyUrl", "", "URL to the firefly remote")
	initCmd.Flags().StringVar(&fireflyUsername, "fireflyUsername", "", "Username for the firefly remote")
	initCmd.Flags().StringVar(&fireflyToken, "fireflyToken", "", "token for the firefly remote")
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(remotesCmd)

	//remote subcommands
	remotesCmd.AddCommand(remotesListCmd)
}
