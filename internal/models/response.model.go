package models

// Response is a unified structure for both single and multiple object responses.
// @Description Generic API response wrapper for both single and list results.
type Response[T any] struct {
	// PaginationMeta contains metadata like total, limit, offset.
	Pagination *MetaMode `json:"pagination,omitempty"`

	// Data holds the payload — can be an object or a list depending on `IsArray`.
	Data T `json:"data"`
}

// MetaMode contains pagination details.
type MetaMode struct {
	// Total number of matching documents.
	// Example: 100
	Total int `json:"total,omitempty" example:"100"`

	// Limit applied to this query.
	// Example: 10
	Limit int `json:"limit,omitempty" example:"10"`

	// Offset/Skip used in this query.
	// Example: 0
	Offset int `json:"offset,omitempty" example:"0"`
}

// BaseSuccess represents a standard API success response.
// @Description Generic success response wrapper.
type BaseSuccess struct {
	// Status of the response (e.g., success, info).
	// Example: success
	Status string `json:"status" example:"success"`

	// StatusCode is the HTTP status code.
	// Example: 200
	StatusCode int `json:"status_code" example:"200"`

	// Message is a human-readable message.
	// Example: Operation completed successfully
	Message string `json:"message,omitempty" example:"Operation completed successfully"`

	// Data holds the response payload. It can be any type.
	Data interface{} `json:"data,omitempty"`
}

// BaseError represents a standard API error response.
// @Description Generic error response wrapper.
type BaseError struct {
	// Status of the response (e.g., error).
	// Example: error
	Status string `json:"status" example:"error"`

	// StatusCode is the HTTP status code.
	// Example: 400
	StatusCode int `json:"status_code" example:"400"`

	// Message is a human-readable error message.
	// Example: Invalid input
	Message string `json:"message,omitempty" example:"Invalid input"`

	// Errors is a list of detailed field-level errors.
	Errors []ErrorDetail `json:"errors,omitempty"`
}

// ErrorDetail provides details about validation or field errors.
type ErrorDetail struct {
	// Field that caused the error.
	// Example: email
	Field string `json:"field" example:"email"`

	// Message describing the error.
	// Example: Email is required
	Message string `json:"message" example:"Email is required"`
}

type DeviceResponseModel struct {
	// Status of the API call
	// Example: success
	Status string `json:"status" example:"success"`

	// HTTP status code
	// Example: 200
	StatusCode int `json:"status_code" example:"200"`

	// Descriptive message
	// Example: Configs fetched successfully
	Message string `json:"message,omitempty" example:"Configs fetched successfully"`

	// Data contains pagination and list of app messengers
	Data struct {
		Pagination *MetaMode     `json:"pagination,omitempty"`
		Data       []DeviceModel `json:"data"`
	} `json:"data"`
}
