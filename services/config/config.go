package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	DefaultHost = "http://localhost:4001"
	DefaultPort = 4001
)

type WasabiConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Region          string
	Endpoint        string
}

// Configuration is the app configuration stored in a json file.
type Configuration struct {
	Host         string       `json:"host" mapstructure:"host"`
	Port         int          `json:"port" mapstructure:"port"`
	WasabiConfig WasabiConfig `json:"wasabiConfig" mapstructure:"wasabiConfig"`
}

// ReadConfigFile read the configuration from the filesystem.
func ReadConfigFile(configFilePath string) (*Configuration, error) {
	if configFilePath == "" {
		viper.SetConfigFile("./config.json")
	} else {
		viper.SetConfigFile(configFilePath)
	}

	viper.SetEnvPrefix("sharelo")
	viper.AutomaticEnv() // read config values from env
	viper.SetDefault("Host", DefaultHost)
	viper.SetDefault("Port", DefaultPort)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}

	configuration := Configuration{}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		return nil, err
	}

	log.Println("readConfigFile")
	log.Printf("%+v", removeSecurityData(configuration))

	return &configuration, nil
}

func removeSecurityData(config Configuration) Configuration {
	clean := config
	return clean
}
