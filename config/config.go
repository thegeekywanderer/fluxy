// Package config package ensures fluxy service is configured properly
package config

import (
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
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Error("error to decode, %v", err)
		return err
	}

	return nil
}
