package cmd

import (
	"errors"
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"

	"github.com/gohypergiant/hyperdrive/hyper/services/config"
	"github.com/manifoldco/promptui"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

var computeRemotesCmd = &cobra.Command{
	Short: "Interact with firefly remotes",
	Use:   "computeRemote",
}
var computeRemotesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Compute Remotes",
	Run: func(cmd *cobra.Command, args []string) {
		remotesMap := config.GetComputeRemotes()
		for name, config := range remotesMap {
			fmt.Println("remote: ", name)
			fmt.Println("    url: ", config.FireflyConfiguration.Url)
		}
	},
}
var computeRemotesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Workspace Remote",
	Run: func(cmd *cobra.Command, args []string) {
		initializeComputeRemoteConfig()
	},
}
var workspaceRemotesCmd = &cobra.Command{
	Short: "Interact with firefly remotes",
	Use:   "workspaceRemote",
}
var workspaceRemotesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Workspace Remote",
	Run: func(cmd *cobra.Command, args []string) {
		initializeWorkspacePersistenceRemoteConfig()
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
var workspacePersistenceRemoteName string
var workspacePersistenceRemoteTypeInput string
var workspacePersistenceRemoteJupyterAPIKey string
var workspaceS3Profile string
var workspaceS3AccessKey string
var workspaceS3Secret string
var workspaceS3Region string
var workspaceS3BucketName string

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
func getS3Config() types.WorkspacePersistenceRemoteConfiguration {

	if workspaceS3Profile == "" {
		workspaceS3Profile = getOptionalString("Enter the name of the configured AWS profile (leave blank to enter a key pair)")
	}

	// If the user has left the profile blank, prompt for keypair
	if workspaceS3Profile == "" {
		if workspaceS3AccessKey == "" {
			workspaceS3AccessKey = getValidatedString("Enter AWS Access Key for provisioning S3 buckets", func(input string) error {
				if len(input) <= 0 {
					return errors.New("must provide an Access Key")
				}
				return nil
			})

		}
		if workspaceS3Secret == "" {
			workspaceS3Secret = getValidatedString("Enter AWS Secret for provisioning S3 buckets", func(input string) error {
				if len(input) <= 0 {
					return errors.New("must provide an Access Secret")
				}
				return nil
			})
		}
	}

	if workspaceS3Region == "" {
		workspaceS3Region = getValidatedString("Enter the region you wish to provision S3 buckets in", func(input string) error {
			if len(input) <= 0 {
				return errors.New("must provide a region")
			}
			return nil
		})
	}

	return types.WorkspacePersistenceRemoteConfiguration{
		Type: types.S3,
		S3Configuration: types.S3WorkspacePersistenceRemoteConfiguration{
			Profile:    workspaceS3Profile,
			AccessKey:  workspaceS3AccessKey,
			Secret:     workspaceS3Secret,
			Region:     workspaceS3Region,
			BucketName: getWorkspaceBucketName(),
		},
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
func getWorkspacePersistenceRemoteType() types.WorkspacePersistenceRemoteType {
	//For now we only support one persistence remote type, but in the future update this to add prompts to configure it
	return types.S3
}
func getComputeRemoteType() types.ComputeRemoteType {
	if computeRemoteTypeInput == "" {
		prompt := promptui.Select{
			Label: "Choose a remote type",
			Items: types.ValidComputeRemoteTypes,
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
		initializeComputeRemoteConfig()
		initializeWorkspacePersistenceRemoteConfig()
	},
}

func initializeWorkspacePersistenceRemoteConfig() {

	var remoteConfig types.WorkspacePersistenceRemoteConfiguration
	if workspacePersistenceRemoteName == "" {
		workspacePersistenceRemoteName = getValidatedString("Enter a name for this remote", func(input string) error {
			if len(input) <= 0 {
				return errors.New("must provide a name")
			}
			return nil
		})
	}
	remoteType := getWorkspacePersistenceRemoteType()
	switch remoteType {
	case types.S3:
		fallthrough
	default:
		remoteConfig = getS3Config()
		fmt.Printf("Adding %s workspace remote", workspacePersistenceRemoteName)
		break
	}

	config.UpdateWorkspaceRemote(workspacePersistenceRemoteName, remoteConfig)
}
func getWorkspaceBucketName() string {
	if workspaceS3BucketName == "" {
		workspaceS3BucketName = getOptionalString("Enter the name of the S3 bucket to use. Bucket names must be globally unique. If it doesn't exist we will attempt to create it the first time we sync. (Leave blank to let us generate one)")
		if workspaceS3BucketName == "" {

			workspaceS3BucketName = uuid.NewString()
			log.Printf("A bucket named %s will be created on the first sync if it doesn't exist.", workspaceS3BucketName)
		}
	}
	return workspaceS3BucketName
}
func initializeComputeRemoteConfig() {
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

	remoteConfig.JupyterAPIKey = getJupyterAPIKey()
	config.UpdateComputeRemote(computeRemoteName, remoteConfig)
}

func getJupyterAPIKey() string {
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
	return computeRemoteJupyterAPIKey
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
	initCmd.Flags().StringVar(&ec2Profile, "ec2Profile", "", "Named AWS profile to use (from ~/.aws/config)")
	initCmd.Flags().StringVar(&ec2AccessKey, "ec2AccessKey", "", "AWS Access Key for provisioning EC2 instances")
	initCmd.Flags().StringVar(&ec2Secret, "ec2Secret", "", "AWS Secret for provisioning EC2 instances")
	initCmd.Flags().StringVar(&ec2Region, "ec2Region", "", "AWS Region for provisioning EC2 instances")

	/*
	* Workspace S3 flags
	 */
	initCmd.Flags().StringVar(&workspaceS3Profile, "workspaceS3Profile", "", "Named AWS profile to use (from ~/.aws/config)")
	initCmd.Flags().StringVar(&workspaceS3AccessKey, "workspaceS3AccessKey", "", "AWS Access Key for provisioning S3 instances")
	initCmd.Flags().StringVar(&workspaceS3Secret, "workspaceS3Secret", "", "AWS Secret for provisioning S3 instances")
	initCmd.Flags().StringVar(&workspaceS3Region, "workspaceS3Region", "", "AWS Region for provisioning S3 instances")
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(initCmd)
	configCmd.AddCommand(computeRemotesCmd)
	configCmd.AddCommand(workspaceRemotesCmd)

	//remote subcommands
	computeRemotesCmd.AddCommand(computeRemotesListCmd)
	computeRemotesCmd.AddCommand(computeRemotesAddCmd)

	//
	//workspaceRemotesCmd.AddCommand(workspaceRemotesListCmd)
	workspaceRemotesCmd.AddCommand(workspaceRemotesAddCmd)
}
