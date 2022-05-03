package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)
type RemoteType string
const (
	EC2 RemoteType = "ec2"
	Firefly RemoteType = "firefly"
)
var ValidRemoteTypes = []RemoteType {
	Firefly, 
	EC2,
}

type RemoteConfiguration struct {
	Type RemoteType `mapstructure:"remote_type" json:"remote_type"`
	FireflyConfiguration FireflyRemoteConfiguration `mapstructure:"firefly" json:"firefly"`
	Ec2Configuration Ec2RemoteConfiguration `mapstructure:"ec2" json:"ec2"`
}
type FireflyRemoteConfiguration struct {
	Url      string `mapstructure:"url" json:"url"`
	HubToken string `mapstructure:"hub_token" json:"hub_token"`
	Username string `mapstructure:"username" json:"username"`
}
type Ec2RemoteConfiguration struct {
	AccessKey string `mapstructure:"access_key" json:"access_key"`
	Secret string `mapstructure:"secret" json:"secret"`
	Region string `mapstructure:"region" json:"region"`
}
type Configuration struct {
	SchemaVersion string `mapstructure:"schema_version" json:"schema_version"`
	Remotes map[string]RemoteConfiguration `mapstructure:"remotes" json:"remotes"`
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
