package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type RemoteConfiguration struct {
	Url      string `mapstructure:"url"`
	HubToken string `mapstructure:"hub_token"`
	Username string `mapstructure:"username"`
}

type Configuration struct {
	Remotes map[string]RemoteConfiguration `mapstructure:"remotes"`
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
