package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/praction-networks/common/helpers"
	"github.com/praction-networks/common/logger"
	"github.com/praction-networks/common/response"

	"github.com/praction-networks/acs-proxy/internal/models"
	"github.com/praction-networks/acs-proxy/internal/services"
	"github.com/praction-networks/acs-proxy/internal/validator"
	"github.com/praction-networks/common/appError"
)

type DeviceHandler struct {
	DeviceService services.DeviceService
}

func NewDeviceHandler(deviceService services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		DeviceService: deviceService,
	}
}

// @Summary      Get Device Info
// @Description  Retrieve device by its serial number
// @Tags         Devices
// @Produce      json
// @Param        sn   path      string  true  "Device Serial Number: minimum 4 last characters"
// @Success      200  {object}  models.DeviceModel
// @Failure      400  {object}  models.BaseError
// @Failure      401  {object}  models.BaseError
// @Failure      403  {object}  models.BaseError
// @Failure      404  {object}  models.BaseError
// @Failure      500  {object}  models.BaseError
// @Router       /acs-proxy/devices/{sn} [get]
// @Security     ApiKeyAuth
func (d *DeviceHandler) GetOnt(w http.ResponseWriter, r *http.Request) {

	ID := chi.URLParam(r, "sn")

	if ID == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "Device ID is required", http.StatusBadRequest, nil))
		return
	}

	logger.Info("Fetching Device", "partial or complete Serial Number", ID)

	var device = models.DeviceSearchID{
		DeviceSN: ID,
	}

	if !helpers.ValidateRequestAndRespond(w, validator.ValidateDeviceSearch(&device), "Validation failed for device search attributes") {
		logger.Error("Validation failed for department attributes")
		return
	}

	deviceResp, err := d.DeviceService.GetOne(r.Context(), device.DeviceSN)
	if err != nil {
		logger.Error("Failed to get Device", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("device fetched successfully", "id", device.DeviceSN)
	response.Send200OK(w, "device fetched successfully", deviceResp)
}

// @Summary      Set PPPoE Credentials
// @Description  Update PPPoE username and password on the device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        data  body      models.SetPPPoECred  true  "PPPoE Credentials"
// @Success      200   {object}  models.DeviceResponseModel
// @Failure      400   {object}  models.BaseError
// @Failure      401   {object}  models.BaseError
// @Failure      500   {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/pppoe [post]
func (d *DeviceHandler) SetPPPoECred(w http.ResponseWriter, r *http.Request) {
	logger.Info("Setting PPPoE credentials")

	var pppoeCred models.SetPPPoECred
	if !helpers.ParseRequestBodyAndRespond(r, w, &pppoeCred) {
		logger.Error("Failed to parse PPPoE credentials")
		return
	}

	if !helpers.ValidateRequestAndRespond(w, validator.ValidatePPPoECred(&pppoeCred), "Validation failed for PPPoE credentials") {
		logger.Error("Validation failed for PPPoE credentials")
		return
	}

	err := d.DeviceService.SetPPPoECredintials(r.Context(), &pppoeCred)
	if err != nil {
		logger.Error("Failed to set PPPoE credentials", err)
		helpers.HandleAppError(w, err)
		return
	}

	if err != nil {
		logger.Error("Failed to set PPPoE credentials", err)
		helpers.HandleAppError(w, err)
		return

	}

	logger.Info("PPPoE credentials set successfully")
	response.Send200OK(w, "PPPoE credentials set successfully", nil)

}

// @Summary      Set WiFi Credentials
// @Description  Update WiFi SSID and Password on the device
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        data  body      models.SetWirelessCred  true  "WiFi Credentials"
// @Success      200   {object}  models.DeviceResponseModel
// @Failure      400   {object}  models.BaseError
// @Failure      401   {object}  models.BaseError
// @Failure      500   {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/wifi [post]
func (d *DeviceHandler) SetWifiCred(w http.ResponseWriter, r *http.Request) {
	logger.Info("Setting PPPoE credentials")

	var wifiCred models.SetWirelessCred
	if !helpers.ParseRequestBodyAndRespond(r, w, &wifiCred) {
		logger.Error("Failed to parse PPPoE credentials")
		return
	}

	if !helpers.ValidateRequestAndRespond(w, validator.ValidateWiFiCred(&wifiCred), "Validation failed for PPPoE credentials") {
		logger.Error("Validation failed for PPPoE credentials")
		return
	}

	err := d.DeviceService.SetWifiCredintials(r.Context(), &wifiCred)
	if err != nil {
		logger.Error("Failed to set PPPoE credentials", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("Wifi credentials set successfully")
	response.Send200OK(w, "Wifi credentials set successfully", nil)
}

// @Summary      Get Devices Not Informed Recently
// @Description  Get devices that haven't sent inform in X days (based on timestamp)
// @Tags         Device
// @Produce      json
// @Param        timestamp  query     string  true  "UTC timestamp in format YYYY-MM-DD HH:mm:ss +0000"
// @Success      200        {object}  models.DeviceResponseModel
// @Failure      400        {object}  models.BaseError
// @Failure      401        {object}  models.BaseError
// @Failure      500        {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/last-inform [get]
func (d *DeviceHandler) GetDevicesByLastInform(w http.ResponseWriter, r *http.Request) {
	timestamp := r.URL.Query().Get("timestamp")
	if timestamp == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "timestamp query param is required", http.StatusBadRequest, nil))
		return
	}

	data, err := d.DeviceService.GetDevicesByLastInformBefore(r.Context(), timestamp)
	if err != nil {
		helpers.HandleAppError(w, err)
		return
	}

	response.Send200OK(w, "Devices fetched", json.RawMessage(data))
}

// @Summary      Get Pending Tasks for Device
// @Description  Fetch all scheduled tasks for a given device
// @Tags         Device
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      200  {object}  models.DeviceResponseModel
// @Failure      400  {object}  models.BaseError
// @Failure      401  {object}  models.BaseError
// @Failure      500  {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/{id}/tasks [get]
func (d *DeviceHandler) GetDeviceTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "Device ID is required", http.StatusBadRequest, nil))
		return
	}

	data, err := d.DeviceService.GetDeviceTasks(r.Context(), id)
	if err != nil {
		helpers.HandleAppError(w, err)
		return
	}

	response.Send200OK(w, "Pending tasks fetched", json.RawMessage(data))
}

// @Summary      Get Specific Device Parameters
// @Description  Fetch projected fields for a given device
// @Tags         Device
// @Produce      json
// @Param        id          path      string  true  "Device ID"
// @Param        projection  query     string  true  "Comma-separated list of fields"
// @Success      200         {object}  models.DeviceResponseModel
// @Failure      400         {object}  models.BaseError
// @Failure      401         {object}  models.BaseError
// @Failure      500         {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/{id}/projection [get]
func (d *DeviceHandler) GetDeviceProjection(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	projection := r.URL.Query().Get("projection")
	if id == "" || projection == "" {
		helpers.HandleAppError(w, appError.New(appError.InvalidInputError, "Device ID and projection are required", http.StatusBadRequest, nil))
		return
	}

	data, err := d.DeviceService.GetDeviceProjection(r.Context(), id, projection)
	if err != nil {
		helpers.HandleAppError(w, err)
		return
	}

	response.Send200OK(w, "Device projection fetched", json.RawMessage(data))
}

// // @Summary      Reboot Device
// // @Description  Reboot a given CPE/ONT device
// // @Tags         Device
// // @Produce      json
// // @Param        id   path      string  true  "Device ID"
// // @Success      200  {object}  models.DeviceResponseModel
// // @Failure      400  {object}  models.BaseError
// // @Failure      401  {object}  models.BaseError
// // @Failure      500  {object}  models.BaseError
// // @Security     ApiKeyAuth
// // @Router       /acs-proxy/devices/{id}/reboot [post]
// func (d *DeviceHandler) RebootDevice(w http.ResponseWriter, r *http.Request) {
// 	deviceID := chi.URLParam(r, "id")
// 	logger.Info("Received reboot request", "deviceID", deviceID)

// 	err := d.DeviceService.Reboot(r.Context(), deviceID)
// 	if err != nil {
// 		logger.Error("Reboot failed", err)
// 		helpers.HandleAppError(w, err)
// 		return
// 	}

// 	response.Send200OK(w, "Device reboot triggered", nil)
// }

// @Summary      Refresh Device
// @Description  Send inform now to refresh device parameters
// @Tags         Device
// @Produce      json
// @Param        id   path      string  true  "Device ID"
// @Success      200  {object}  models.DeviceResponseModel
// @Failure      400  {object}  models.BaseError
// @Failure      401  {object}  models.BaseError
// @Failure      500  {object}  models.BaseError
// @Security     ApiKeyAuth
// @Router       /acs-proxy/devices/{id}/refresh [post]
func (d *DeviceHandler) RefreshDevice(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Received refresh request", "deviceID", deviceID)

	err := d.DeviceService.Refresh(r.Context(), deviceID)
	if err != nil {
		logger.Error("Refresh failed", err)
		helpers.HandleAppError(w, err)
		return
	}

	response.Send200OK(w, "Device refresh triggered", nil)
}

// @Summary      Trigger GetParameterValues task
// @Description  Requests the CPE to return values for listed parameters
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Param        body body models.GetParameterValuesRequest true "Parameters to retrieve"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/get-parameters [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) GetParameterValues(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering GetParameterValues task", "deviceID", deviceID)

	var req models.GetParameterValuesRequest
	if !helpers.ParseRequestBodyAndRespond(r, w, &req) ||
		!helpers.ValidateRequestAndRespond(w, validator.ValidateGetParameterValues(&req), "Invalid parameter list") {
		logger.Error("Validation failed for GetParameterValues request")
		return
	}

	if err := h.DeviceService.GetParameterValues(r.Context(), deviceID, &req); err != nil {
		logger.Error("Failed to trigger GetParameterValues task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("GetParameterValues task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Parameter fetch task submitted", nil)
}

// @Summary      Trigger SetParameterValues task
// @Description  Set multiple configuration parameters on the CPE
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Param        body body models.SetParameterValuesRequest true "Parameters to set"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/set-parameters [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) SetParameterValues(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering SetParameterValues task", "deviceID", deviceID)

	var req models.SetParameterValuesRequest
	if !helpers.ParseRequestBodyAndRespond(r, w, &req) ||
		!helpers.ValidateRequestAndRespond(w, validator.ValidateSetParameterValues(&req), "Invalid parameter values") {
		logger.Error("Validation failed for SetParameterValues request")
		return
	}

	if err := h.DeviceService.SetParameterValues(r.Context(), deviceID, &req); err != nil {
		logger.Error("Failed to trigger SetParameterValues task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("SetParameterValues task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Set parameter task submitted", nil)
}

// @Summary      Trigger RefreshObject task
// @Description  Refreshes an object in the CPE
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Param        body body models.RefreshObjectRequest true "Object to refresh"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/refresh-object [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) RefreshObject(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering RefreshObject task", "deviceID", deviceID)

	var req models.RefreshObjectRequest
	if !helpers.ParseRequestBodyAndRespond(r, w, &req) ||
		!helpers.ValidateRequestAndRespond(w, validator.ValidateRefreshObject(&req), "Invalid object name") {
		logger.Error("Validation failed for RefreshObject request")
		return
	}

	if err := h.DeviceService.RefreshObject(r.Context(), deviceID, &req); err != nil {
		logger.Error("Failed to trigger RefreshObject task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("RefreshObject task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Refresh task submitted", nil)
}

// @Summary      Trigger AddObject task
// @Description  Adds a new object to the CPE
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Param        body body models.AddObjectRequest true "Object to add"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/add-object [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) AddObject(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering AddObject task", "deviceID", deviceID)

	var req models.AddObjectRequest
	if !helpers.ParseRequestBodyAndRespond(r, w, &req) ||
		!helpers.ValidateRequestAndRespond(w, validator.ValidateAddObject(&req), "Invalid object name") {
		logger.Error("Validation failed for AddObject request")
		return
	}

	if err := h.DeviceService.AddObject(r.Context(), deviceID, &req); err != nil {
		logger.Error("Failed to trigger AddObject task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("AddObject task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Add object task submitted", nil)
}

// @Summary      Trigger DeleteObject task
// @Description  Deletes an object from the CPE
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Param        body body models.DeleteObjectRequest true "Object to delete"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/delete-object [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering DeleteObject task", "deviceID", deviceID)

	var req models.DeleteObjectRequest
	if !helpers.ParseRequestBodyAndRespond(r, w, &req) ||
		!helpers.ValidateRequestAndRespond(w, validator.ValidateDeleteObject(&req), "Invalid object name") {
		logger.Error("Validation failed for DeleteObject request")
		return
	}

	if err := h.DeviceService.DeleteObject(r.Context(), deviceID, &req); err != nil {
		logger.Error("Failed to trigger DeleteObject task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("DeleteObject task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Delete object task submitted", nil)
}

// @Summary      Trigger Reboot task
// @Description  Sends a reboot task to the device
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/reboot [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) TriggerReboot(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering Reboot task", "deviceID", deviceID)

	if err := h.DeviceService.RebootDevice(r.Context(), deviceID); err != nil {
		logger.Error("Failed to trigger Reboot task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("Reboot task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Reboot task submitted", nil)
}

// @Summary      Trigger Factory Reset task
// @Description  Sends a factory reset task to the device
// @Tags         Device Tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Device ID"
// @Success      200 {object} models.BaseSuccess
// @Failure      400,500 {object} models.BaseError
// @Router       /acs-proxy/devices/{id}/factory-reset [post]
// @Security     ApiKeyAuth
func (h *DeviceHandler) TriggerFactoryReset(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "id")
	logger.Info("Triggering FactoryReset task", "deviceID", deviceID)

	if err := h.DeviceService.FactoryResetDevice(r.Context(), deviceID); err != nil {
		logger.Error("Failed to trigger FactoryReset task", err)
		helpers.HandleAppError(w, err)
		return
	}

	logger.Info("FactoryReset task submitted successfully", "deviceID", deviceID)
	response.Send200OK(w, "Factory reset task submitted", nil)
}
