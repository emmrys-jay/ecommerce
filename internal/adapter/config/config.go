package config

import (
	"log"

	"github.com/spf13/viper"
)

// Setup initialize configuration
var (
	// Params ParamsConfiguration
	config *Configuration
)

// Params = getConfig.Params
func Setup() *Configuration {
	var configuration *Configuration

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Params = configuration.Params
	config = configuration
	log.Println("Configurations loading successfully")

	return config
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return config
}
