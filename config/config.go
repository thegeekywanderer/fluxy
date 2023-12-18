// Package config package ensures fluxy service is configured properly
package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/thegeekywanderer/fluxy/logger"
)

// Configuration for fluxy
type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	// Check if the .env file exists
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(".env not found")
		viper.AutomaticEnv() // Allow Viper to read from environment variables
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Error("error to decode, %v", err)
		return err
	}

	return nil
}
