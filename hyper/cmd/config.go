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
var remoteTypeInput string
var fireflyUrl string
var fireflyUsername string
var fireflyToken string
var ec2Profile string
var ec2AccessKey string
var ec2Secret string
var ec2Region string

func getFireflyUsername() string {

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
func getFireflyUrl() string {

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

func getFireflyToken() string {

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
		fireflyUrl = getFireflyUrl()
	}
	if fireflyUsername == "" {
		fireflyUsername = getFireflyUsername()
	}
	if fireflyToken == "" {
		fireflyToken = getFireflyToken()
	}
	return config.RemoteConfiguration{
		Type:                 config.Firefly,
		FireflyConfiguration: config.FireflyRemoteConfiguration{Url: fireflyUrl, Username: fireflyUsername, HubToken: fireflyToken},
	}
}
func getEC2Config() config.RemoteConfiguration {

	if ec2Profile == "" {
		ec2Profile = getEC2Profile()
	}

	// If the user has left the profile blank, prompt for keypair
	if ec2Profile == "" {
		if ec2AccessKey == "" {
			ec2AccessKey = getEC2AccessKey()
		}
		if ec2Secret == "" {
			ec2Secret = getEC2Secret()
		}
	}

	if ec2Region == "" {
		ec2Region = getEC2Region()
	}

	return config.RemoteConfiguration{
		Type: config.EC2,
		EC2Configuration: config.EC2RemoteConfiguration{
			Profile:   ec2Profile,
			AccessKey: ec2AccessKey,
			Secret:    ec2Secret,
			Region:    ec2Region,
		},
	}
}
func getEC2Profile() string {
	prompt := promptui.Prompt{
		Label: "Enter the name of the configured AWS profile (leave blank to enter a key pair)",
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}
func getEC2AccessKey() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide an Access Key")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter AWS Access Key for provisioning EC2 instances",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}
func getEC2Secret() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide an Secret")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter AWS Secret for provisioning EC2 instances",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}
func getEC2Region() string {

	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New("must provide a region")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "Enter the region you wish to provision EC2 instances in",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
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
		case config.EC2:
			remoteConfig = getEC2Config()
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
	/*
	* Firefly flags
	 */
	initCmd.Flags().StringVar(&fireflyUrl, "fireflyUrl", "", "URL to the firefly remote")
	initCmd.Flags().StringVar(&fireflyUsername, "fireflyUsername", "", "Username for the firefly remote")
	initCmd.Flags().StringVar(&fireflyToken, "fireflyToken", "", "token for the firefly remote")
	/*
	* EC2 flags
	 */
	initCmd.Flags().StringVar(&ec2AccessKey, "ec2Profile", "", "Named AWS profile to use (from ~/.aws/config)")
	initCmd.Flags().StringVar(&ec2AccessKey, "ec2AccessKey", "", "AWS Access Key for provisioning EC2 instances")
	initCmd.Flags().StringVar(&ec2Secret, "ec2Secret", "", "AWS Secret for provisioning EC2 instances")
	initCmd.Flags().StringVar(&ec2Region, "ec2Region", "", "AWS Region for provisioning EC2 instances")

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(remotesCmd)

	//remote subcommands
	remotesCmd.AddCommand(remotesListCmd)
}
