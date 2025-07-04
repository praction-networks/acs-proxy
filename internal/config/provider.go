package config

import "log"

// ProvideConfig is a Wire provider function clearly loading EnvConfig.
func ProvideConfig() EnvConfig {
	if err := LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return AppConfig
}

func ProvideGeniacsConfig(cfg EnvConfig) *GeniacsConfig {
	return &cfg.GenieacsEnv
}
