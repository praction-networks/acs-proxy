package config

import (
	"fmt"
)

func ValidateAppConfig() error {
	if AppConfig.ServerEnv.Port == "" {
		return fmt.Errorf("server port is not configured")
	}
	if AppConfig.MongoDBEnv.Host == "" || AppConfig.MongoDBEnv.Port == "" {
		return fmt.Errorf("MongoDB host or port is not configured")
	}
	return nil
}
