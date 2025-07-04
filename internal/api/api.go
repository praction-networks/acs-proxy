package api

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/praction-networks/acs-proxy/internal/config"
	"github.com/praction-networks/acs-proxy/internal/dependency"
	"github.com/praction-networks/acs-proxy/internal/handlers"
	"github.com/praction-networks/acs-proxy/internal/monitoring"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"

	middlewares "github.com/praction-networks/acs-proxy/internal/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/swaggo/swag/example/basic/docs"
)

//	@title			Domain Nat Core Service
//	@version		1.0.0
//	@description	This is i9 GenieACS Proxy for API documents
//	@termsOfService	https://praction.in/terms/

//	@contact.name	I9 API Support
//	@contact.url	http://www.praction.in/support
//	@contact.email	support@praction.in

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host						acsproxy.praction.in
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and then your key

//	@securityDefinitions.apikey	CookieAuth
//	@in							cookie
//	@name						auth
//	@description				Authentication via cookie. Set by the server after successful login.

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api

func SetupRouter(container *dependency.AppContainer) *chi.Mux {

	SetupSwaggerInfo(container.Config)
	r := chi.NewRouter()

	// ✅ Chi default middleware
	r.Use(middleware.RealIP)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(60 * time.Second))

	// Your custom logger + Prometheus + tracing middleware
	r.Use(helpers.RequestIDMiddleware)
	r.Use(helpers.LoggingMiddleware)
	r.Use(helpers.MetricsMiddleware)

	r.Use(cors.Handler(
		cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://ispm.praction.in",
				"https://ispm.praction.in",
				"http://dashboard.praction.in",
				"https://dashboard.praction.in",
				"http://dash.praction.in",
				"https://dash.praction.in",
				"http://acsproxy.praction.in",
				"https://acsproxy.praction.in",
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

	setupSwaggerRoutes(r, container.Config)

	r.Route("/api/v1/acs-proxy", func(r chi.Router) {
		r.Use(middlewares.APIKeyAuthMiddleware(container.Config.AuthEnv.APIKey))

		// Health Check
		healthHandler := &monitoring.HealthHandler{
			MongoClient: container.MongoClient,
		}
		r.Get("/health", healthHandler.GetHealth)

		logLevelHandler := &handlers.LogLevelHandler{}
		r.Post("/log-level", logLevelHandler.LogLevelHandler)

		// Domain Routes
		setupDeviceRoutes(r, container.DeviceHandler)
		setupTaskRoutes(r, container.TaskHandler)

	})

	logger.Info("Router initialized")
	return r

}

// setupDeviceRoutes defines all /devices related routes
func setupDeviceRoutes(r chi.Router, handler *handlers.DeviceHandler) {
	r.Route("/devices", func(d chi.Router) {
		d.Get("/{sn}", handler.GetOnt)
		d.Get("/{id}/projection", handler.GetDeviceProjection)
		d.Get("/{id}/tasks", handler.GetDeviceTasks)
		d.Get("/last-inform", handler.GetDevicesByLastInform)
		d.Post("/pppoe", handler.SetPPPoECred)
		d.Post("/wifi", handler.SetWifiCred)
		d.Post("/{id}/refresh", handler.RefreshDevice)
		r.Post("/{id}/get-parameters", handler.GetParameterValues)
		r.Post("/{id}/set-parameters", handler.SetParameterValues)
		r.Post("/{id}/refresh-object", handler.RefreshObject)
		r.Post("/{id}/add-object", handler.AddObject)
		r.Post("/{id}/delete-object", handler.DeleteObject)
		r.Post("/{id}/reboot", handler.TriggerReboot)
		r.Post("/{id}/factory-reset", handler.TriggerFactoryReset)
	})
}

// setupTaskRoutes defines /tasks routes
func setupTaskRoutes(r chi.Router, handler *handlers.TaskHandler) {
	r.Route("/tasks", func(t chi.Router) {
		t.Post("/{task_id}/retry", handler.RetryTask)
		t.Delete("/{task_id}", handler.DeleteTask)
	})
}

// setupSwaggerRoutes configures Swagger routes.
func setupSwaggerRoutes(r chi.Router, cfg config.EnvConfig) {
	var apiUrl string
	if cfg.EnvironmentEnv.Host == "" {
		apiUrl = "http://ispm.praction.in/swagger/json"
	} else {
		apiUrl = fmt.Sprintf("http://%s/swagger/json", cfg.EnvironmentEnv.Host)
	}

	r.Get("/swagger-api/*", httpSwagger.Handler(httpSwagger.URL(apiUrl)))
	swaggerHandler := &handlers.SwaggerHandler{}
	r.Get("/swagger/json", swaggerHandler.GetSwaggerJson)
	r.Get("/swagger/yaml", swaggerHandler.GetSwaggerYaml)
}

func SetupSwaggerInfo(cfg config.EnvConfig) {
	if cfg.EnvironmentEnv.Version == "" || cfg.EnvironmentEnv.Host == "" {
		logger.Error("Invalid Swagger configuration: Version or Host is missing")
	}

	docs.SwaggerInfo.Title = "ACS Poroxy Service"
	docs.SwaggerInfo.Description = "This is i9 ACS Proxy for API documents"
	docs.SwaggerInfo.Version = cfg.EnvironmentEnv.Version
	docs.SwaggerInfo.Host = cfg.EnvironmentEnv.Host
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	// Add security definitions and security settings to the Swagger template
	docs.SwaggerInfo.SwaggerTemplate = `
	{
		"swagger": "2.0",
		"info": {
			"title": "{{.Title}}",
			"description": "{{.Description}}",
			"version": "{{.Version}}"
		},
		"host": "{{.Host}}",
		"basePath": "{{.BasePath}}",
		"schemes": {{marshal .Schemes}},
		"paths": {},
		"securityDefinitions": {
			"BearerAuth": {
				"type": "apiKey",
				"name": "Authorization",
				"in": "header",
				"description": "Type 'Bearer' followed by a space and then your token"
			},
			"CookieAuth": {
				"type": "apiKey",
				"name": "auth",
				"in": "cookie",
				"description": "Authentication via cookie. Set by the server after successful login."
			}
		},
		"security": [
			{
				"BearerAuth": [],
				"CookieAuth": []
			}
		]
	}`
}
