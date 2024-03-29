package config

import (
	"github.com/spf13/viper"
	"fmt"
	"os"
)

func LoadServerConfig() error {
	// Check if config.env path exists in the root directory.
	if _, err := os.Stat("config.env"); os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist\n", "env.config")
		return err
	} else {
		viper.SetConfigFile("config.env")
	}

	// Check if local config.env file can be read-in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Unable to read in config.env")	
		return err
	}

	// Check if OpenWeather API key is set.
	variable := viper.GetString( "OW_API_KEY")	
	if variable == "" {
		fmt.Println("The config variable, OW_API_KEY does not exist")
		return err
	}

	//fmt.Println("...Finished successfully loading config")
	return nil
}

func ConfigGetString(variable string) string {
	return viper.GetString(variable)
}