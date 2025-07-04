package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/praction-networks/common/logger"
)

func (a *App) Start(ctx context.Context, cancel context.CancelFunc) error {
	// Verify NATS connection clearly

	if err := a.mongoClient.Ping(ctx, nil); err != nil {
		logger.Error("MongoDB connection not available", err)
		return fmt.Errorf("MongoDB connection not available: %w", err)
	}

	go a.ConnectionMonitor.Start(ctx, 3)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				if !a.ConnectionMonitor.IsHealthy() {
					logger.Error("Detected disconnection from MongoDB or NATS. Shutting down.")
					// Cancel the main context to trigger shutdown
					cancel()
					return
				}
			}
		}
	}()

	logger.Info("Starting HTTP Server", "address", a.server.Addr)

	errCh := make(chan error, 1)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server failed: %w", err)
		}
		close(errCh)
	}()

	logger.Info("Starting Metrics Prometheus Scraper", "address", a.metricsServer.Addr)

	go func() {
		if err := a.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Metrics server error: %v\n", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		logger.Error("Server encountered error", "error", err)
		return err
	case <-ctx.Done():
		logger.Warn("Server stopping due to context cancellation...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server graceful shutdown failed", "error", err)
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		if err := a.metricsServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Metrics server graceful shutdown failed", "error", err)
			return fmt.Errorf("metrics server shutdown failed: %w", err)
		}

		logger.Info("Server shutdown completed gracefully")
		return nil
	}
}
