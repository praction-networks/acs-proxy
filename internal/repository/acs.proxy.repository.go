package repository

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDeviceRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewMongoDeviceRepository(collection *mongo.Collection) DeviceRepository {
	return &MongoDeviceRepository{
		Collection: collection,
	}
}

// ✅ Interface assertion
var _ DeviceRepository = (*MongoDeviceRepository)(nil)

func (r *MongoDeviceRepository) GetBySn(ctx context.Context, deviceSn string) (*models.DeviceModel, error) {
	logFields := []interface{}{"partial_sn", deviceSn}
	logger.Info("Starting GetDevice operation", logFields...)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		logger.Error("Failed to list device IDs", "error", err)
		return nil, appError.New(appError.DBTransactionError, "Failed to query devices", http.StatusInternalServerError, err)
	}
	defer cursor.Close(ctx)

	var matchingID string
	for cursor.Next(ctx) {
		var doc struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&doc); err != nil {
			logger.Error("Failed to decode device ID", "error", err)
			continue
		}

		if strings.Contains(doc.ID, deviceSn) {
			matchingID = doc.ID
			break
		}
	}

	if matchingID == "" {
		logger.Warn("Device not found for partial SN", logFields...)
		return nil, appError.New(appError.EntityNotFound, "Device not found", http.StatusNotFound, nil)
	}

	// Fetch full document by matching ID
	var device models.DeviceModel
	err = r.Collection.FindOne(ctx, bson.M{"_id": matchingID}).Decode(&device)
	if err != nil {
		logger.Error("Failed to fetch device by full ID", "error", err)
		return nil, appError.New(appError.DBTransactionError, "Failed to fetch device", http.StatusInternalServerError, err)
	}

	logFields = append(logFields, "resolved_id", matchingID)
	logger.Info("Device fetched successfully", logFields...)
	return &device, nil
}

func (r *MongoDeviceRepository) GetAllDevices(ctx context.Context) ([]models.DeviceModel, error) {
	logFields := []interface{}{}
	logger.Info("Fetching all devices (debug)", logFields...)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to fetch devices", "error", err)
		return nil, appError.New(appError.DBTransactionError, "Failed to fetch devices", http.StatusInternalServerError, err)
	}
	defer cursor.Close(ctx)

	var devices []models.DeviceModel
	for cursor.Next(ctx) {
		var d models.DeviceModel
		if err := cursor.Decode(&d); err != nil {
			logger.Error("Failed to decode device", "error", err)
			continue
		}

		// Log serial number for debugging
		logger.Debug("Device found", "id", d.ID, "serial", d.DeviceID.SerialNumber)
		devices = append(devices, d)
	}

	if err := cursor.Err(); err != nil {
		logger.Error("Cursor error", "error", err)
		return nil, appError.New(appError.DBTransactionError, "Cursor error", http.StatusInternalServerError, err)
	}

	logger.Info("Total devices fetched", "count", len(devices))
	return devices, nil
}
