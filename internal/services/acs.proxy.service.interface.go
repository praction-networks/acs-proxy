package services

import (
	"context"

	"github.com/praction-networks/acs-proxy/internal/models"
)

type DeviceService interface {
	GetOne(ctx context.Context, deviceSN string) (models.Response[*models.DeviceModel], error)
	GetAll(ctx context.Context) (models.Response[[]models.DeviceModel], error)
	SetPPPoECredintials(ctx context.Context, pppoeCredentials *models.SetPPPoECred) error
	SetWifiCredintials(ctx context.Context, wifiCredentials *models.SetWirelessCred) error
	GetDeviceProjection(ctx context.Context, deviceID, projection string) ([]byte, error)
	GetDeviceTasks(ctx context.Context, deviceID string) ([]byte, error)
	GetDevicesByLastInformBefore(ctx context.Context, timestamp string) ([]byte, error)
	Reboot(ctx context.Context, deviceID string) error
	Refresh(ctx context.Context, deviceID string) error
	GetParameterValues(ctx context.Context, deviceID string, req *models.GetParameterValuesRequest) error
	SetParameterValues(ctx context.Context, deviceID string, req *models.SetParameterValuesRequest) error
	RefreshObject(ctx context.Context, deviceID string, req *models.RefreshObjectRequest) error
	AddObject(ctx context.Context, deviceID string, req *models.AddObjectRequest) error
	DeleteObject(ctx context.Context, deviceID string, req *models.DeleteObjectRequest) error
	RebootDevice(ctx context.Context, deviceID string) error
	FactoryResetDevice(ctx context.Context, deviceID string) error
}

type TaskService interface {
	RetryTask(ctx context.Context, taskID string) error
	DeleteTask(ctx context.Context, taskID string) error
}
