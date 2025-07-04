package mongodb

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/common/logger"
)

// Global MongoDB client
var (
	clientInstance     *mongo.Client
	clientInstanceOnce sync.Once
)

// InitializeMongo initializes the MongoDB client and sets up indexes
func InitializeMongo(ctx context.Context, clientConfig config.MongoConfig) (*mongo.Client, error) {
	var err error

	clientInstanceOnce.Do(func() {
		var uri string
		if clientConfig.URI != "" {
			uri = clientConfig.URI
		} else {
			uri = fmt.Sprintf("mongodb://%s:%s", clientConfig.Host, clientConfig.Port)
		}

		clientOptions := options.Client().ApplyURI(uri).SetMaxPoolSize(100).SetMinPoolSize(10)

		clientInstance, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Fatal("Failed to connect to MongoDB", err)
			return
		}

		pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err = clientInstance.Ping(pingCtx, nil); err != nil {
			logger.Fatal("Failed to ping MongoDB", err)
			return
		}

		logger.Info("MongoDB connected successfully")
		// Now create collections if not exists
		db := clientInstance.Database(clientConfig.Database)
		for _, collectionName := range CollectionNames {
			err := db.CreateCollection(ctx, collectionName)
			if err != nil && !isNamespaceExistsError(err) {
				logger.Error("Collection creation failed", "collection", collectionName, "error", err)
			}
		}
	})

	if err != nil {
		return nil, err
	}

	// Initialize collections
	db := clientInstance.Database(clientConfig.Database)
	for name, collectionName := range CollectionNames {
		Collections[name] = db.Collection(collectionName)
		logger.Info(fmt.Sprintf("Initialized collection: %s", collectionName))
	}

	return clientInstance, nil
}

// GetClient returns the MongoDB client instance
func GetClient() *mongo.Client {
	if clientInstance == nil {
		logger.Fatal("MongoDB client not initialized")
	}
	return clientInstance
}

// CloseClient closes the MongoDB connection
func CloseClient(ctx context.Context) error {
	if clientInstance != nil {
		if err := clientInstance.Disconnect(ctx); err != nil {
			logger.Error("Error closing MongoDB connection", err)
			return err
		} else {
			logger.Info("MongoDB connection closed successfully")
			return nil
		}
	}
	return nil
}

// GetCollection retrieves a MongoDB collection by its identifier
func GetCollection(name string) *mongo.Collection {
	collection, exists := Collections[name]
	if !exists {
		logger.Fatal(fmt.Sprintf("Collection %s not initialized", name))
	}
	return collection
}

func isNamespaceExistsError(err error) bool {
	if err == nil {
		return false
	}
	if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
		// MongoDB error code 48 = NamespaceExists
		return true
	}
	return false
}
