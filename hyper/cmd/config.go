package cmd

import (
	"errors"
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"log"
	"os"
	"strings"

	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/manifoldco/promptui"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

var remotesCmd = &cobra.Command{
	Short: "Interact with firefly remotes",
	Use:   "remote",
}
var remotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ComputeRemotes",
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

var computeRemoteName string
var computeRemoteTypeInput string
var computeRemoteJupyterAPIKey string
var fireflyUrl string
var fireflyUsername string
var fireflyToken string
var ec2Profile string
var ec2AccessKey string
var ec2Secret string
var ec2Region string

func getValidatedString(message string, validate promptui.ValidateFunc) string {
	prompt := promptui.Prompt{
		Label:    message,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}
func getOptionalString(message string) string {

	prompt := promptui.Prompt{
		Label: message,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return result
}
func getFireflyConfig() types.ComputeRemoteConfiguration {

	if fireflyUrl == "" {
		fireflyUrl = getOptionalString("Enter the remote URL [default: Use Hypergiant hosted Hyperdrive]")
	}
	if fireflyUsername == "" {
		fireflyUsername = getValidatedString("Enter your username", func(input string) error {
			if len(input) <= 0 {
				return errors.New("must provide a username")
			}
			return nil
		})
	}
	if fireflyToken == "" {
		fireflyToken = getValidatedString("Enter your firefly API token",
			func(input string) error {
				if len(input) <= 0 {
					return errors.New("must provide an API token")
				}
				return nil
			})
	}
	return types.ComputeRemoteConfiguration{
		Type:                 types.Firefly,
		FireflyConfiguration: types.FireflyComputeRemoteConfiguration{Url: fireflyUrl, Username: fireflyUsername, HubToken: fireflyToken},
	}
}
func getEC2Config() types.ComputeRemoteConfiguration {

	if ec2Profile == "" {
		ec2Profile = getOptionalString("Enter the name of the configured AWS profile (leave blank to enter a key pair)")
	}

	// If the user has left the profile blank, prompt for keypair
	if ec2Profile == "" {
		if ec2AccessKey == "" {
			ec2AccessKey = getValidatedString("Enter AWS Access Key for provisioning EC2 instances", func(input string) error {
				if len(input) <= 0 {
					return errors.New("must provide an Access Key")
				}
				return nil
			})

		}
		if ec2Secret == "" {
			ec2Secret = getValidatedString("Enter AWS Secret for provisioning EC2 instances", func(input string) error {
				if len(input) <= 0 {
					return errors.New("must provide an Access Secret")
				}
				return nil
			})
		}
	}

	if ec2Region == "" {
		ec2Region = getValidatedString("Enter the region you wish to provision EC2 instances in", func(input string) error {
			if len(input) <= 0 {
				return errors.New("must provide a region")
			}
			return nil
		})
	}

	return types.ComputeRemoteConfiguration{
		Type: types.EC2,
		EC2Configuration: types.EC2ComputeRemoteConfiguration{
			Profile:   ec2Profile,
			AccessKey: ec2AccessKey,
			Secret:    ec2Secret,
			Region:    ec2Region,
		},
	}
}
func getComputeRemoteType() types.ComputeRemoteType {
	if computeRemoteTypeInput == "" {
		prompt := promptui.Select{
			Label: "Choose a remote type",
			Items: types.ValidRemoteTypes,
		}
		_, remoteTypeInput, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return types.ComputeRemoteType(remoteTypeInput)
	}
	if computeRemoteTypeInput == string(types.Firefly) {
		return types.Firefly
	}
	if computeRemoteTypeInput == string(types.EC2) {
		return types.EC2
	}

	fmt.Println("Invalid or unsupported remote type")
	os.Exit(1)
	return types.ComputeRemoteType(computeRemoteTypeInput)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Config",
	Run: func(cmd *cobra.Command, args []string) {
		initializeDeployRemoteConfig()
		initializeWorkspacePersistenceRemoteConfig()
	},
}

func initializeWorkspacePersistenceRemoteConfig() {

}
func initializeDeployRemoteConfig() {
	var remoteConfig types.ComputeRemoteConfiguration
	if computeRemoteName == "" {
		computeRemoteName = getValidatedString("Enter a name for this remote", func(input string) error {
			if len(input) <= 0 {
				return errors.New("must provide a name")
			}
			return nil
		})
	}
	remoteType := getComputeRemoteType()
	// if
	switch remoteType {
	case types.EC2:
		remoteConfig = getEC2Config()
	case types.Firefly:
		fallthrough
	default:
		remoteConfig = getFireflyConfig()
		fmt.Printf("Adding %s remote at %s", computeRemoteName, fireflyUrl)
		break
	}

	if computeRemoteJupyterAPIKey == "" {
		computeRemoteJupyterAPIKey = getOptionalString("Enter a Jupyter token to use for remote instances [leave blank to generate one]")
		if computeRemoteJupyterAPIKey == "" {

			pass, err := password.Generate(64, 10, 0, true, true)
			if err != nil {
				log.Fatal(err)
			}
			computeRemoteJupyterAPIKey = strings.ToUpper(pass)
			log.Printf("A Jupyter Token of %s has been generated. You will need it to access the UI on remote instances. If you need to find this later you can find it in your ~/.hyperdrive file", computeRemoteJupyterAPIKey)
		}
	}
	remoteConfig.JupyterAPIKey = computeRemoteJupyterAPIKey
	config.UpdateRemote(computeRemoteName, remoteConfig)
}

func init() {
	initCmd.Flags().StringVar(&computeRemoteName, "computeRemoteName", "", "Name of the remote for the config")
	initCmd.Flags().StringVar(&computeRemoteTypeInput, "computeRemoteType", "", "Remote type [firefly|ec2]")
	initCmd.Flags().StringVar(&computeRemoteJupyterAPIKey, "computeRemoteJupyterAPIKey", "", "API key to use on jupyter instances that get created")
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
