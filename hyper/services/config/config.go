package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/gohypergiant/hyperdrive/hyper/types"
	"os"

	awssdkconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

func GetConfig() types.Configuration {
	var config types.Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config
}
func GetNamedProfileConfig(s3AwsProfile string) types.NamedProfileConfiguration {
	var namedProfileConfig types.NamedProfileConfiguration
	awsConfigFilePath := awssdkconfig.DefaultSharedConfigFilename()
	if _, errFile := os.Stat(awsConfigFilePath); errFile == nil {
		// AWS config file exists at $HOME/.aws/config. We're good.
	} else if errors.Is(errFile, os.ErrNotExist) {
		fmt.Println("Error:", awsConfigFilePath, "does not exist. Please create one.")
		os.Exit(1)
	}

	ctx := context.TODO()
	cfg, errConfig := awssdkconfig.LoadDefaultConfig(ctx,
		awssdkconfig.WithSharedConfigProfile(s3AwsProfile))
	if errConfig != nil {
		fmt.Println("Error:", errConfig)
		os.Exit(1)
	}

	creds, errCreds := cfg.Credentials.Retrieve(ctx)
	if errCreds != nil {
		fmt.Println("Error:", errCreds)
		os.Exit(1)
	}

	namedProfileConfig.AccessKey = creds.AccessKeyID
	namedProfileConfig.Secret = creds.SecretAccessKey
	namedProfileConfig.Token = creds.SessionToken
	namedProfileConfig.Region = cfg.Region
	return namedProfileConfig
}
func GetComputeRemotes() map[string]types.ComputeRemoteConfiguration {
	var remotesMap map[string]types.ComputeRemoteConfiguration
	err := viper.UnmarshalKey("compute_remotes", &remotesMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return remotesMap
}
func GetComputeRemote(name string) types.ComputeRemoteConfiguration {
	remotes := GetComputeRemotes()
	return remotes[name]
}
func GetWorkspacePersistenceRemotes() map[string]types.WorkspacePersistenceRemoteConfiguration {
	var remotesMap map[string]types.WorkspacePersistenceRemoteConfiguration
	err := viper.UnmarshalKey("workspace_remotes", &remotesMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return remotesMap
}
func GetWorkspacePersistenceRemote(name string) types.WorkspacePersistenceRemoteConfiguration {
	remotes := GetWorkspacePersistenceRemotes()
	return remotes[name]
}
func UpdateComputeRemote(name string, configuration types.ComputeRemoteConfiguration) {
	var config types.Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config.ComputeRemotes == nil {
		config.ComputeRemotes = make(map[string]types.ComputeRemoteConfiguration)
	}
	config.ComputeRemotes[name] = configuration
	viper.Set("compute_remotes", config.ComputeRemotes)
	viper.WriteConfig()

}
func UpdateWorkspaceRemote(name string, configuration types.WorkspacePersistenceRemoteConfiguration) {
	var config types.Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config.WorkspacePersistenceRemotes == nil {
		config.WorkspacePersistenceRemotes = make(map[string]types.WorkspacePersistenceRemoteConfiguration)
	}
	config.WorkspacePersistenceRemotes[name] = configuration
	viper.Set("workspace_remotes", config.WorkspacePersistenceRemotes)
	viper.WriteConfig()

}
