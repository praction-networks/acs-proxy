package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"
)

type SwaggerHandler struct{}

// GetSwaggerJson godoc
//
//	@Summary		Get Swagger JSON documentation
//	@Description	Serve the OpenAPI documentation in JSON format
//	@Tags			Swagger
//	@Produce		json
//	@Success		200	{string}	string				"Swagger JSON content"
//	@Failure		404	{object}	models.BaseError	"Swagger JSON file not found"
//	@Failure		503	{object}	models.BaseError	"Service Unavailable - Failed to read or serve the file"
//	@Router			/swagger/json [get]
//	@Security		ApiKeyAuth
func (s *SwaggerHandler) GetSwaggerJson(w http.ResponseWriter, r *http.Request) {
	filePath := "docs/swagger.json"
	serveSwaggerFile(w, filePath, "JSON")
}

// GetSwaggerYaml godoc
//
//	@Summary		Get Swagger YAML documentation
//	@Description	Serve the OpenAPI documentation in YAML format
//	@Tags			Swagger
//	@Produce		plain
//	@Success		200	{string}	string				"Swagger YAML content"
//	@Failure		404	{object}	models.BaseError	"Swagger YAML file not found"
//	@Failure		503	{object}	models.BaseError	"Service Unavailable - Failed to read or serve the file"
//	@Router			/swagger/yaml [get]
//	@Security		ApiKeyAuth
func (s *SwaggerHandler) GetSwaggerYaml(w http.ResponseWriter, r *http.Request) {
	filePath := "docs/swagger.yaml"
	serveSwaggerFile(w, filePath, "YAML")
}

// serveSwaggerFile is a helper to serve Swagger files (JSON or YAML)
func serveSwaggerFile(w http.ResponseWriter, filePath string, fileType string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logger.Error(fmt.Sprintf("%s file not found", fileType), err)
			response.Send404NotFound(w, fileType+" file not found")
		} else {
			logger.Error(fmt.Sprintf("Failed to read %s file", fileType), err)
			response.Send503ServiceUnavailable(w, "Failed to read "+fileType+" file")
		}
		return
	}

	// Set the Content-Type header and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(data)

	if writeErr != nil {
		logger.Error(fmt.Sprintf("Failed to write %s response", fileType), err)
		response.Send503ServiceUnavailable(w, "Failed to write "+fileType+" response")
	}
}
