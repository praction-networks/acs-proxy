package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/praction-networks/acs-proxy/cmd/wire"
	"github.com/praction-networks/acs-proxy/internal/app"
	"github.com/praction-networks/common/logger"
)

func main() {
	// Use signal.NotifyContext to cancel on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize DI container
	container, err := wire.InitializeContainer()
	if err != nil {
		logger.Fatal("Failed to initialize container", err)
		os.Exit(1)
	}

	// Initialize application
	appInstance, err := app.New(container)
	if err != nil {
		logger.Fatal("Application initialization error", err)
		os.Exit(1)
	}

	// Start app in a goroutine
	go func() {
		if err := appInstance.Start(ctx, nil); err != nil {
			logger.Error("Application runtime error", err)
			stop() // stop on error
		}
	}()

	// Block until context is canceled (e.g., SIGINT)
	<-ctx.Done()
	logger.Warn("Shutdown signal received, initiating graceful shutdown")

	// Graceful shutdown with timeout
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := appInstance.Stop(ctxShutdown); err != nil {
		logger.Fatal("Graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Application shutdown gracefully")
}
