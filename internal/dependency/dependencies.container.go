// dependency/dependency.go
package dependency

import (
	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/acs-proxy/internal/handlers"
	"github.com/praction-networks/acs-proxy/internal/monitoring"
	"github.com/praction-networks/acs-proxy/internal/repository"
	"github.com/praction-networks/acs-proxy/internal/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AppContainer struct {
	Config            config.EnvConfig
	Logger            *zap.Logger
	MongoClient       *mongo.Client
	ConnectionMonitor *monitoring.ConnectionMonitor

	// Repositories
	DeviceRepository repository.DeviceRepository

	// Services
	DeviceService services.DeviceService
	TaskService   services.TaskService

	// Handlers
	DeviceHandler *handlers.DeviceHandler
	TaskHandler   *handlers.TaskHandler
}
