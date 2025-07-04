package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var AppConfig EnvConfig

func LoadConfig() error {
	viper.AddConfigPath("/internal/config")
	viper.SetConfigName("environment")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		return fmt.Errorf("error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error unmarshaling config file: %v", err)
		return fmt.Errorf("error unmarshaling config file: %v", err)
	}

	log.Println("Configuration loaded successfully.")
	return nil
}
