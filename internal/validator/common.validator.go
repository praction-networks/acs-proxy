package validator

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/praction-networks/common/response"
)

// NewValidator returns a validator instance with all custom rules registered
func NewValidator() *validator.Validate {
	v := validator.New()

	if err := RegisterCustomValidations(v); err != nil {
		log.Fatalf("Validator registration error: %v", err)
	}

	return v
}

// extractValidationErrors extracts validation errors and formats them into a response.
func extractValidationErrors(err error) []response.ErrorDetail {
	var errors []response.ErrorDetail
	if err == nil {
		return errors
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			errors = append(errors, response.ErrorDetail{
				Field:   e.StructNamespace(),
				Message: formatValidationMessage(e),
			})
		}
	} else {
		errors = append(errors, response.ErrorDetail{
			Field:   "general",
			Message: err.Error(), // fallback if it's another kind of validation error
		})
	}

	return errors
}

// formatValidationMessage maps validation tags to readable error messages
func formatValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "The " + e.Field() + " field is required."
	case "required_if":
		return "The " + e.Field() + " field is required when " + e.Param() + " is present."
	case "required_with":
		return "The " + e.Field() + " field is required when " + e.Param() + " is provided."
	case "oneof":
		return "The " + e.Field() + " field must be one of the following values: " + e.Param()
	case "email":
		return "The " + e.Field() + " field must be a valid email address."
	case "url":
		return "The " + e.Field() + " must be a valid URL."
	case "gst":
		return "The " + e.Field() + " field must be a valid GST number."
	case "pan":
		return "The " + e.Field() + " field must be a valid PAN number."
	case "pincode":
		return "The " + e.Field() + " field must be a valid 6-digit pincode."
	case "cuid2":
		return "The " + e.Field() + " field must be a valid CUID2 string."
	case "fqdn_or_ip":
		return "The " + e.Field() + " must be a valid IP or domain name."
	case "min":
		return "The " + e.Field() + " must have a minimum length of " + e.Param() + " characters."
	case "max":
		return "The " + e.Field() + " must have a maximum length of " + e.Param() + " characters."
	case "dive":
		return "One or more values in the " + e.Field() + " list are invalid."
	case "atleastone":
		return "At least one of" + e.Field() + "or the other related field must be provided"
	case "alphanum":
		return "The " + e.Field() + " field must be alphanumeric."
	case "singleword":
		return "The " + e.Field() + " field must be a single word."
	default:
		return "Validation failed on field " + e.Field() + " with constraint '" + e.Tag() + "'"
	}
}
