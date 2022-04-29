package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type RemoteConfiguration interface {
	RemoteType string `mapstructure:"remote_type"`
}
type FireflyRemoteConfiguration struct {
	Url      string `mapstructure:"url"`
	HubToken string `mapstructure:"hub_token"`
	Username string `mapstructure:"username"`
}
type Ec2RemoteConfiguration struct {
	AccessKey string `mapstructure:"access_key"`
	Username string `mapstructure:"secret"`
}
type Configuration struct {
	SchemaVersion string `mapstructure:"configSchemaVersion"`
	Remotes struct {
		Firefly map[string]FireflyRemoteConfiguration `mapstructure:"firefly"`
		Ec2     map[string]Ec2RemoteConfiguration     `mapstructure:"ec2"`
	} `mapstructure:"remotes"`
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

func GetRemotes() map[string]FireflyRemoteConfiguration {
	var remotesMap map[string]FireflyRemoteConfiguration
	err := viper.UnmarshalKey("remotes", &remotesMap)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return remotesMap
}
func GetRemote(name string) FireflyRemoteConfiguration {
	remotes := GetRemotes()
	return remotes[name]
}
func UpdateRemote(name string, configuration FireflyRemoteConfiguration) {
	var config Configuration
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if config.Remotes == nil {
		config.Remotes = make(map[string]FireflyRemoteConfiguration)
	}
	config.Remotes[name] = configuration
	viper.Set("remotes", config.Remotes)
	viper.WriteConfig()

}
