package validator

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations registers all custom validation functions
func RegisterCustomValidations(v *validator.Validate) error {
	if err := v.RegisterValidation("gst", validateGSTNumber); err != nil {
		return fmt.Errorf("failed to register GST validator: %w", err)
	}
	if err := v.RegisterValidation("pan", validatePANNumber); err != nil {
		return fmt.Errorf("failed to register PAN validator: %w", err)
	}
	if err := v.RegisterValidation("pincode", validatePincode); err != nil {
		return fmt.Errorf("failed to register Pincode validator: %w", err)
	}
	if err := v.RegisterValidation("cuid2", validateCUID2); err != nil {
		return fmt.Errorf("failed to register CUID2 validator: %w", err)
	}
	if err := v.RegisterValidation("fqdn_or_ip", validateFQDNorIP); err != nil {
		return fmt.Errorf("failed to register FQDN/IP validator: %w", err)
	}
	if err := v.RegisterValidation("url", validateURL); err != nil {
		return fmt.Errorf("failed to register URL validator: %w", err)
	}
	if err := v.RegisterValidation("singleword", validateSingleWord); err != nil {
		return fmt.Errorf("failed to register singleword validator: %w", err)
	}

	return nil
}

// validateGSTNumber checks if the value is a valid GST number
func validateGSTNumber(fl validator.FieldLevel) bool {
	pattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, fl.Field().String())
	return matched
}

// validatePANNumber checks if the value is a valid PAN number
func validatePANNumber(fl validator.FieldLevel) bool {
	pattern := `^[A-Z]{5}[0-9]{4}[A-Z]{1}$`
	matched, _ := regexp.MatchString(pattern, fl.Field().String())
	return matched
}

// validatePincode checks if the value is a valid Indian 6-digit pincode
func validatePincode(fl validator.FieldLevel) bool {
	pincodeStr := fmt.Sprintf("%d", fl.Field().Int())
	pattern := `^[1-9][0-9]{5}$`
	matched, _ := regexp.MatchString(pattern, pincodeStr)
	return matched
}

// validateCUID2 checks if the value is a valid CUID2 string
func validateCUID2(fl validator.FieldLevel) bool {
	cuid2Regex := regexp.MustCompile(`^c[a-z0-9]{8,}$`)
	return cuid2Regex.MatchString(fl.Field().String())
}

// validateFQDNorIP checks if the value is a valid IP or FQDN
func validateFQDNorIP(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if net.ParseIP(value) != nil {
		return true
	}
	fqdnRegex := regexp.MustCompile(`^(?i:[a-z0-9]+(?:[-.][a-z0-9]+)*)\.[a-z]{2,}$`)
	return fqdnRegex.MatchString(value)
}

// validateURL checks if the value is a valid URL
func validateURL(fl validator.FieldLevel) bool {
	str := fl.Field().String()
	_, err := url.ParseRequestURI(str)
	return err == nil
}

func validateSingleWord(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return !regexp.MustCompile(`[\s\r\n\t]`).MatchString(value)
}
