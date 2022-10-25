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

// func GetConfig() {
// 	// Set the file name of the configurations file
// 	viper.SetConfigName("config")

// 	// Set the path to look for the configurations file
// 	viper.AddConfigPath(".")

// 	// // Enable VIPER to read Environment Variables
// 	// viper.AutomaticEnv()

// 	viper.SetConfigType("yml")
// 	var configuration Configurations

// 	if err := viper.ReadInConfig(); err != nil {
// 		fmt.Printf("Error reading config file, %s", err)
// 	}

// 	err := viper.Unmarshal(&configuration)
// 	if err != nil {
// 		fmt.Printf("Unable to decode into struct, %v", err)
// 	}

// 	// Reading variables using the model
// 	fmt.Println("Reading variables using the model..")
// 	fmt.Println("Database is\t", configuration.Stmp.Host)
// 	fmt.Println("Port is\t\t", configuration.Stmp.Pass)
// 	fmt.Println("EXAMPLE_PATH is\t", configuration.Stmp.Port)
// 	fmt.Println("EXAMPLE_VAR is\t", configuration.Stmp.User)
// }
