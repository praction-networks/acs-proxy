package app

import (
	"context"

	"github.com/praction-networks/acs-proxy/internal/database/mongodb"
	"github.com/praction-networks/common/logger"
)

func (a *App) Stop(ctx context.Context) error {
	logger.Info("Starting graceful shutdown...")

	if err := mongodb.CloseClient(ctx); err != nil {
		logger.Error("Error during MongoDB shutdown", "error", err)
		return err
	}

	logger.Info("All resources cleaned up successfully for Domain Service")
	return nil
}
