package config

import (
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Sftp SftpConfigurations
}

// ServerConfigurations exported
type SftpConfigurations struct {
	Host string
	User string
	Pass string
	Port int
}

func GetConfig() (config Configurations, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type

	// viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
