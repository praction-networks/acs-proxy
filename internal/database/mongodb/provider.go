package mongodb

import (
	"context"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/acs-proxy/internal/repository"
	"github.com/praction-networks/common/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

// ✅ Now correctly returns (*mongo.Client, string, error)// ✅ **Fixed Function: Provide Only `*mongo.Client, error`**
func ProvideMongoDB(cfg config.EnvConfig) (*mongo.Client, error) {
	client, err := InitializeMongo(context.Background(), cfg.MongoDBEnv)
	if err != nil {
		logger.Error("MongoDB initialization failed:", err)
		return nil, err // ✅ Return error instead of panicking
	}
	return client, nil
}

// **New Function to Provide Database Name Separately**
func ProvideMongoDBName(cfg config.EnvConfig) string {
	return cfg.MongoDBEnv.Database
}

//  Provide the struct with both collections

type MongoCollections struct {
	SkeletonCollection *mongo.Collection
}

// ✅ Provide `MongoDomainRepository` directly
func ProvideMongoSkeletonRepository(client *mongo.Client, dbName string) repository.DeviceRepository {
	collection := client.Database(dbName).Collection("devices")
	return &repository.MongoDeviceRepository{Collection: collection}
}
