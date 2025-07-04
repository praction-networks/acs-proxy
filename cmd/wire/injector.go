//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/acs-proxy/internal/database/mongodb"
	"github.com/praction-networks/acs-proxy/internal/dependency"
	"github.com/praction-networks/acs-proxy/internal/genieacs"
	"github.com/praction-networks/acs-proxy/internal/handlers"
	"github.com/praction-networks/acs-proxy/internal/logging"
	"github.com/praction-networks/acs-proxy/internal/monitoring"
	"github.com/praction-networks/acs-proxy/internal/services"

	"go.mongodb.org/mongo-driver/mongo"
)

var mongoDBSet = wire.NewSet(
	mongodb.ProvideMongoDB,     // Provides `*mongo.Client`
	mongodb.ProvideMongoDBName, // Provides `string` (database name)

)

var genieACSClientSet = wire.NewSet(
	config.ProvideGeniacsConfig,
	genieacs.NewClient,
)

var repositorySet = wire.NewSet(
	mongodb.ProvideMongoSkeletonRepository,
)

// ✅ Service and Handler Providers
var serviceSet = wire.NewSet(
	services.NewDeviceService,
)

var handlerSet = wire.NewSet(
	handlers.NewDeviceHandler,
)

var monitoringSet = wire.NewSet(
	ProvideConnectionMonitor,
)

func ProvideConnectionMonitor(
	mongo *mongo.Client,
) *monitoring.ConnectionMonitor {
	return monitoring.New(mongo)
}

// Production Container Initialization
func InitializeContainer() (*dependency.AppContainer, error) {
	wire.Build(
		config.ProvideConfig,
		logging.ProvideLogger,
		mongoDBSet,
		genieACSClientSet,
		repositorySet,
		serviceSet,
		handlerSet,
		monitoringSet,
		wire.Struct(new(dependency.AppContainer), "*"),
	)
	return &dependency.AppContainer{}, nil
}

// ✅ Test Configuration Struct
type TestConfig struct {
	MongoURI string
	NatsURI  string
}

// ✅ Provide Test Configuration
func ProvideTestConfig(testCfg TestConfig) config.EnvConfig {
	return config.EnvConfig{
		ServerEnv: config.ServerConfig{
			Port: "3000",
		},
		MongoDBEnv: config.MongoConfig{
			URI:      testCfg.MongoURI,
			Database: "testdb",
		},
		LoggerEnv: config.LoggerConfig{
			LogLevel: "debug",
		},
	}
}

// ✅ Initialize Test Container
func InitializeTestContainer(testCfg TestConfig) (*dependency.AppContainer, error) {
	wire.Build(
		config.ProvideConfig,
		logging.ProvideLogger,
		mongoDBSet,
		genieACSClientSet,
		repositorySet,
		serviceSet,
		handlerSet,
		monitoringSet,
		wire.Struct(new(dependency.AppContainer), "*"),
	)
	return &dependency.AppContainer{}, nil
}
