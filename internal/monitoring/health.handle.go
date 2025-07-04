package monitoring

import (
	"context"
	"net/http"
	"time"

	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"

	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	MongoClient *mongo.Client
}

// GetHealth returns the health status of the service
//
//	@Summary		Health check endpoint
//	@Description	Returns the health status of the service
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Health status"
//	@Router			/api/v1/domain/health [get]
func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	mongoStatus := "healthy"
	overallStatus := "healthy"
	httpStatus := http.StatusOK

	var errorDetails []response.ErrorDetail

	// ✅ Check Mongo separately
	if h.MongoClient == nil || h.MongoClient.Ping(ctx, nil) != nil {
		mongoStatus = "down"
		overallStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable

		errorDetails = append(errorDetails, response.ErrorDetail{
			Field:   "mongo",
			Message: "MongoDB connection is down",
		})
	}

	resp := HealthStatusResponse{
		Service: "domain-service",
		Status:  overallStatus,
		Mongo:   mongoStatus,
	}

	logger.Info("Health check status",
		"status", resp.Status,
		"mongo", resp.Mongo,
	)

	if overallStatus == "unhealthy" {
		response.SendCustomError(w, "Service health check failed", errorDetails, httpStatus)
		return
	}

	response.Send200OK(w, "Service is healthy", resp)
}

type HealthStatusResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Mongo   string `json:"mongo"`
	NATS    string `json:"nats"`
}
