package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/common/response"
)

func ValidateDeviceSearch(req *models.DeviceSearchID) []response.ErrorDetail {
	v := NewValidator()

	var validationErrors []response.ErrorDetail

	err := v.Struct(req)
	// Validate general structure
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.StructNamespace(),
				Message: formatValidationMessage(e),
			})
		}
	}

	return validationErrors
}

func ValidatePPPoECred(req *models.SetPPPoECred) []response.ErrorDetail {
	v := NewValidator()

	var validationErrors []response.ErrorDetail

	err := v.Struct(req)
	// Validate general structure
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.StructNamespace(),
				Message: formatValidationMessage(e),
			})
		}
	}

	return validationErrors
}

func ValidateWiFiCred(req *models.SetWirelessCred) []response.ErrorDetail {
	v := NewValidator()

	var validationErrors []response.ErrorDetail

	err := v.Struct(req)
	// Validate general structure
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.StructNamespace(),
				Message: formatValidationMessage(e),
			})
		}
	}

	return validationErrors
}

func ValidateGetParameterValues(req *models.GetParameterValuesRequest) []response.ErrorDetail {
	return validateStruct(req)
}

func ValidateSetParameterValues(req *models.SetParameterValuesRequest) []response.ErrorDetail {
	return validateStruct(req)
}

func ValidateRefreshObject(req *models.RefreshObjectRequest) []response.ErrorDetail {
	return validateStruct(req)
}

func ValidateAddObject(req *models.AddObjectRequest) []response.ErrorDetail {
	return validateStruct(req)
}

func ValidateDeleteObject(req *models.DeleteObjectRequest) []response.ErrorDetail {
	return validateStruct(req)
}

func validateStruct(req interface{}) []response.ErrorDetail {
	v := NewValidator()
	var validationErrors []response.ErrorDetail

	err := v.Struct(req)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, response.ErrorDetail{
				Field:   e.StructNamespace(),
				Message: formatValidationMessage(e),
			})
		}
	}
	return validationErrors
}
