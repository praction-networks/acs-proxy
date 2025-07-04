package app

import (
	"net/http"

	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/metrics"

	"github.com/praction-networks/acs-proxy/internal/api"
	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/acs-proxy/internal/dependency"
	"github.com/praction-networks/acs-proxy/internal/monitoring"

	"go.mongodb.org/mongo-driver/mongo"
)

// App encapsulates application dependencies and state.
type App struct {
	Config            config.EnvConfig
	router            http.Handler
	mongoClient       *mongo.Client
	server            *http.Server
	metricsServer     *http.Server
	ConnectionMonitor *monitoring.ConnectionMonitor
	container         *dependency.AppContainer
}

// New initializes the application with dependencies clearly from Wire.
func New(container *dependency.AppContainer) (*App, error) {
	metrics.RegisterAllMetrics()

	app := &App{
		Config:            container.Config,
		mongoClient:       container.MongoClient,
		container:         container,
		ConnectionMonitor: container.ConnectionMonitor,
	}

	app.router = api.SetupRouter(app.container)
	app.server = initializeHttpServer(app.router)
	app.metricsServer = StartMetricsServer(app.mongoClient)

	logger.Info("Application initialized successfully")
	return app, nil
}
