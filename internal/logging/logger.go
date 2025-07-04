package logging

import (
	"fmt"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/common/logger"

	"go.uber.org/zap"
)

func ProvideLogger(cfg config.EnvConfig) *zap.Logger {
	loggerConfig := logger.LoggerConfig{
		LogLevel: cfg.LoggerEnv.LogLevel,
	}

	err := logger.InitializeLogger(loggerConfig)
	if err != nil {
		panic(fmt.Sprintf("Logger initialization failed: %v", err))
	}

	return logger.GetGlobalLogger()
}
