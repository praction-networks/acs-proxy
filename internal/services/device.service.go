package services

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/praction-networks/acs-proxy/internal/genieacs"
	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/acs-proxy/internal/repository"
	"github.com/praction-networks/common/appError"
	"github.com/praction-networks/common/logger"
)

type TimeProvider interface {
	Now() time.Time
}

type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

type DeviceServiceImpl struct {
	DeviceRepo     repository.DeviceRepository
	GenieACSClient *genieacs.Client
	TimeProvider   TimeProvider
}

func NewDeviceService(
	deviceRepo repository.DeviceRepository,
	acsClient *genieacs.Client,
) DeviceService {
	return &DeviceServiceImpl{
		DeviceRepo:     deviceRepo,
		TimeProvider:   &RealTimeProvider{},
		GenieACSClient: acsClient,
	}
}

func (s *DeviceServiceImpl) GetOne(ctx context.Context, deviceSN string) (models.Response[*models.DeviceModel], error) {
	logger.Info("Getting domain by ID", "DeviceSN", deviceSN)

	deviceResp, err := s.DeviceRepo.GetBySn(ctx, deviceSN)
	if err != nil {
		logger.Error("Failed to get device serial", err)
		return models.Response[*models.DeviceModel]{}, err
	}
	if deviceResp == nil {
		logger.Warn("Domain not found", "deviceSN", deviceSN)
		return models.Response[*models.DeviceModel]{}, appError.New(appError.EntityNotFound, "Device SN not found", http.StatusNotFound, nil)
	}

	logger.Info("Domain found", "CUID", deviceSN)

	return models.Response[*models.DeviceModel]{
		Data: deviceResp,
	}, nil
}

func (s *DeviceServiceImpl) GetAll(ctx context.Context) (models.Response[[]models.DeviceModel], error) {
	logger.Info("Getting all domains")

	devices, err := s.DeviceRepo.GetAllDevices(ctx)
	if err != nil {
		logger.Error("Failed to get all domains", err)
		return models.Response[[]models.DeviceModel]{}, err
	}

	return models.Response[[]models.DeviceModel]{
		Data: devices,
	}, nil
}

func (s *DeviceServiceImpl) SetPPPoECredintials(ctx context.Context, pppoeCredentials *models.SetPPPoECred) error {
	logger.Info("Setting PPPoE credentials", "DeviceID", pppoeCredentials.DeviceID)

	resp, err := s.GenieACSClient.SetPPPoECredentials(pppoeCredentials)
	if err != nil {
		logger.Error("Failed to send PPPoE credentials to GenieACS", err)
		return appError.New(appError.ExternalServiceError, "Failed to send PPPoE credentials", http.StatusBadGateway, err)
	}

	logger.Info("PPPoE credentials successfully set", "Status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) SetWifiCredintials(ctx context.Context, wifiCredentials *models.SetWirelessCred) error {
	logger.Info("Setting WiFi credentials", "DeviceID", wifiCredentials.DeviceID)

	resp, err := s.GenieACSClient.SetWiFiCredentials(wifiCredentials)
	if err != nil {
		logger.Error("Failed to send WiFi credentials to GenieACS", err)
		return appError.New(appError.ExternalServiceError, "Failed to send WiFi credentials", http.StatusBadGateway, err)
	}

	logger.Info("WiFi credentials successfully set", "Status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) GetDevicesByLastInformBefore(ctx context.Context, timestamp string) ([]byte, error) {
	logger.Info("Querying devices with _lastInform before", "timestamp", timestamp)

	resp, err := s.GenieACSClient.FindDevicesByLastInformBefore(timestamp)
	if err != nil {
		logger.Error("GenieACS query failed", err)
		return nil, appError.New(appError.ExternalServiceError, "GenieACS query failed", http.StatusBadGateway, err)
	}
	return resp.Body(), nil
}

func (s *DeviceServiceImpl) GetDeviceTasks(ctx context.Context, deviceID string) ([]byte, error) {
	logger.Info("Querying tasks for device", "deviceID", deviceID)

	resp, err := s.GenieACSClient.GetPendingTasksForDevice(deviceID)
	if err != nil {
		logger.Error("GenieACS task query failed", err)
		return nil, appError.New(appError.ExternalServiceError, "GenieACS task query failed", http.StatusBadGateway, err)
	}
	return resp.Body(), nil
}

func (s *DeviceServiceImpl) GetDeviceProjection(ctx context.Context, deviceID, projection string) ([]byte, error) {
	logger.Info("Querying projection for device", "deviceID", deviceID, "projection", projection)

	resp, err := s.GenieACSClient.GetDeviceProjection(deviceID, projection)
	if err != nil {
		logger.Error("Projection fetch failed", err)
		return nil, appError.New(appError.ExternalServiceError, "Projection fetch failed", http.StatusBadGateway, err)
	}
	return resp.Body(), nil
}

func (s *DeviceServiceImpl) Reboot(ctx context.Context, deviceID string) error {
	if strings.TrimSpace(deviceID) == "" {
		return appError.New(appError.InvalidInputError, "Device ID is required", http.StatusBadRequest, nil)
	}

	logger.Info("DeviceService: Triggering device reboot", "deviceID", deviceID)
	_, err := s.GenieACSClient.RebootDevice(deviceID)
	if err != nil {
		logger.Error("DeviceService: Failed to reboot device", err)
		return appError.New(appError.InvalidOperation, "Failed to reboot device", http.StatusInternalServerError, err)
	}

	return nil
}

func (s *DeviceServiceImpl) Refresh(ctx context.Context, deviceID string) error {
	if strings.TrimSpace(deviceID) == "" {
		return appError.New(appError.InvalidInputError, "Device ID is required", http.StatusBadRequest, nil)
	}

	logger.Info("DeviceService: Refreshing device", "deviceID", deviceID)
	_, err := s.GenieACSClient.RefreshObject(deviceID)
	if err != nil {
		logger.Error("DeviceService: Failed to refresh device", err)
		return appError.New(appError.InvalidOperation, "Failed to refresh device", http.StatusInternalServerError, err)
	}

	return nil
}

func (s *DeviceServiceImpl) GetParameterValues(ctx context.Context, deviceID string, req *models.GetParameterValuesRequest) error {
	logger.Info("DeviceService: Triggering GetParameterValues", "deviceID", deviceID)

	task := map[string]any{
		"name":           "getParameterValues",
		"parameterNames": req.ParameterNames,
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger GetParameterValues", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger GetParameterValues", http.StatusBadGateway, err)
	}

	logger.Info("GetParameterValues task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) SetParameterValues(ctx context.Context, deviceID string, req *models.SetParameterValuesRequest) error {
	logger.Info("DeviceService: Triggering SetParameterValues", "deviceID", deviceID)

	task := map[string]any{
		"name":            "setParameterValues",
		"parameterValues": req.ParameterValues,
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger SetParameterValues", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger SetParameterValues", http.StatusBadGateway, err)
	}

	logger.Info("SetParameterValues task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) RefreshObject(ctx context.Context, deviceID string, req *models.RefreshObjectRequest) error {
	logger.Info("DeviceService: Triggering RefreshObject", "deviceID", deviceID)

	task := map[string]any{
		"name":       "refreshObject",
		"objectName": req.ObjectName,
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger RefreshObject", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger RefreshObject", http.StatusBadGateway, err)
	}

	logger.Info("RefreshObject task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) AddObject(ctx context.Context, deviceID string, req *models.AddObjectRequest) error {
	logger.Info("DeviceService: Triggering AddObject", "deviceID", deviceID)

	task := map[string]any{
		"name":       "addObject",
		"objectName": req.ObjectName,
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger AddObject", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger AddObject", http.StatusBadGateway, err)
	}

	logger.Info("AddObject task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) DeleteObject(ctx context.Context, deviceID string, req *models.DeleteObjectRequest) error {
	logger.Info("DeviceService: Triggering DeleteObject", "deviceID", deviceID)

	task := map[string]any{
		"name":       "deleteObject",
		"objectName": req.ObjectName,
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger DeleteObject", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger DeleteObject", http.StatusBadGateway, err)
	}

	logger.Info("DeleteObject task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) RebootDevice(ctx context.Context, deviceID string) error {
	logger.Info("DeviceService: Triggering RebootDevice", "deviceID", deviceID)

	task := map[string]any{
		"name": "reboot",
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger RebootDevice", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger RebootDevice", http.StatusBadGateway, err)
	}

	logger.Info("RebootDevice task triggered", "status", resp.StatusCode())
	return nil
}

func (s *DeviceServiceImpl) FactoryResetDevice(ctx context.Context, deviceID string) error {
	logger.Info("DeviceService: Triggering FactoryResetDevice", "deviceID", deviceID)

	task := map[string]any{
		"name": "factoryReset",
	}

	resp, err := s.GenieACSClient.TriggerTask(deviceID, task)
	if err != nil {
		logger.Error("Failed to trigger FactoryResetDevice", err)
		return appError.New(appError.ExternalServiceError, "Failed to trigger FactoryResetDevice", http.StatusBadGateway, err)
	}

	logger.Info("FactoryResetDevice task triggered", "status", resp.StatusCode())
	return nil
}
