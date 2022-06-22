package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	awssdkconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/spf13/viper"
)

type RemoteType string

const (
	EC2     RemoteType = "ec2"
	Firefly RemoteType = "firefly"
)

var ValidRemoteTypes = []RemoteType{
	Firefly,
	EC2,
}

type RemoteConfiguration struct {
	Type                 RemoteType                 `mapstructure:"remote_type" json:"remote_type"`
	FireflyConfiguration FireflyRemoteConfiguration `mapstructure:"firefly" json:"firefly"`
	EC2Configuration     EC2RemoteConfiguration     `mapstructure:"ec2" json:"ec2"`
	JupyterAPIKey        string                     `mapstructure:"jupyter_api_key" json:"jupyter_api_key"`
}
type FireflyRemoteConfiguration struct {
	Url      string `mapstructure:"url" json:"url"`
	HubToken string `mapstructure:"hub_token" json:"hub_token"`
	Username string `mapstructure:"username" json:"username"`
}
type EC2RemoteConfiguration struct {
	Profile   string `mapstructure:"profile" json:"profile"`
	AccessKey string `mapstructure:"access_key" json:"access_key"`
	Secret    string `mapstructure:"secret" json:"secret"`
	Region    string `mapstructure:"region" json:"region"`
	Token     string
}
type Configuration struct {
	SchemaVersion string                         `mapstructure:"schema_version" json:"schema_version"`
	Remotes       map[string]RemoteConfiguration `mapstructure:"remotes" json:"remotes"`
}
type NamedProfileConfiguration struct {
	AccessKey string
	Secret string
	Token string
	Region string
}

func GetConfig() Configuration {
	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return config
}
func GetNamedProfileConfig(s3AwsProfile string) NamedProfileConfiguration {
	var namedProfileConfig NamedProfileConfiguration
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
func GetRemotes() map[string]RemoteConfiguration {
	var remotesMap map[string]RemoteConfiguration
	err := viper.UnmarshalKey("remotes", &remotesMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return remotesMap
}
func GetRemote(name string) RemoteConfiguration {
	remotes := GetRemotes()
	return remotes[name]
}
func UpdateRemote(name string, configuration RemoteConfiguration) {
	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config.Remotes == nil {
		config.Remotes = make(map[string]RemoteConfiguration)
	}
	config.Remotes[name] = configuration
	viper.Set("remotes", config.Remotes)
	viper.WriteConfig()

}
