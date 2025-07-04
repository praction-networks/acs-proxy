package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/common/response"
)

func ValidateLogLevelUpdate(logLevel *models.Logging) []response.ErrorDetail {
	v := validator.New()
	if err := v.Struct(logLevel); err != nil {
		return extractValidationErrors(err)
	}
	return nil
}
