package app

import (
	"fmt"
	"net/http"

	"github.com/praction-networks/common/logger"

	"github.com/praction-networks/acs-proxy/internal/config"
)

func initializeHttpServer(r http.Handler) *http.Server {
	port := config.AppConfig.ServerEnv.Port
	if port == "" {
		port = "3030" // Default port
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	logger.Debug("HTTP Server initialized successfully for Domain Service", "port", port)
	return server
}
