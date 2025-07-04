package models

type Logging struct {
	LogLevel string `json:"logLevel" validate:"required,oneof=trace debug info warn error"`
}
