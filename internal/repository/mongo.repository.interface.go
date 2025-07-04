package repository

import (
	"context"

	"github.com/praction-networks/acs-proxy/internal/models"
)

type DeviceRepository interface {
	GetBySn(ctx context.Context, deviceSn string) (*models.DeviceModel, error)
	GetAllDevices(ctx context.Context) ([]models.DeviceModel, error)
}
