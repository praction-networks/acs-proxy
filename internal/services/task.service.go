// File: internal/services/task.go
package services

import (
	"context"
	"net/http"

	"github.com/praction-networks/acs-proxy/internal/genieacs"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
)

type TaskServiceImpl struct {
	GenieACSClient *genieacs.Client
	TimeProvider   TimeProvider
}

func NewTaskService(
	acsClient *genieacs.Client,
) DeviceService {
	return &DeviceServiceImpl{
		TimeProvider:   &RealTimeProvider{},
		GenieACSClient: acsClient,
	}
}

func (s *TaskServiceImpl) RetryTask(ctx context.Context, taskID string) error {
	logger.Info("Retrying task", "taskID", taskID)
	resp, err := s.GenieACSClient.RetryTask(taskID)
	if err != nil {
		logger.Error("Failed to retry task", err)
		return appError.New(appError.ExternalServiceError, "Failed to retry task", http.StatusBadGateway, err)
	}
	logger.Info("Task retry successful", "status", resp.StatusCode())
	return nil
}

func (s *TaskServiceImpl) DeleteTask(ctx context.Context, taskID string) error {
	logger.Info("Deleting task", "taskID", taskID)
	resp, err := s.GenieACSClient.DeleteTask(taskID)
	if err != nil {
		logger.Error("Failed to delete task", err)
		return appError.New(appError.ExternalServiceError, "Failed to delete task", http.StatusBadGateway, err)
	}
	logger.Info("Task deletion successful", "status", resp.StatusCode())
	return nil
}
