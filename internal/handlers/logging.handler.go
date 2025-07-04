package handlers

import (
	"net/http"

	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/acs-proxy/internal/validator"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"
)

// LogLevelHandler handles HTTP requests related to updating the log level of the domain-service.
type LogLevelHandler struct{}

// LogLevelHandler sets up the log level for the domain-service.
//
//	@Summary		Setup log level for domain-service
//	@Description	With this endpoint you can setup log level for domain-service
//	@Tags			logs
//	@Accept			json
//	@Produce		json
//	@Param			message	body		models.Logging	true	"Log Message"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Router			/acs-proxy/log-level [post]
func (ll *LogLevelHandler) LogLevelHandler(w http.ResponseWriter, r *http.Request) {

	var logLevel *models.Logging
	// Parse the request body

	if !helpers.ParseRequestBodyAndRespond(r, w, &logLevel) {
		return
	}

	if !helpers.ValidateRequestAndRespond(w, validator.ValidateLogLevelUpdate(logLevel), "Invalid Log Level") {
		return
	}

	if err := logger.UpdateLogLevel(logLevel.LogLevel); err != nil {
		logger.Error("Error updating log level", err)
		helpers.HandleAppError(w, err)
	}

	logger.Info("Log level updated", "new_level", logLevel.LogLevel)
	response.Send200OK(w, "Log level updated sussfully", nil)

}
