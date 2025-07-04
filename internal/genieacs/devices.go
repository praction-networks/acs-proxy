package genieacs

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// GetDeviceByID retrieves a device by its full Device ID
func (c *Client) GetDeviceByID(deviceID string) (*resty.Response, error) {
	query := map[string]interface{}{"_id": deviceID}
	return c.http.R().
		SetQueryParam("query", toJSONString(query)).
		Get("/devices")
}

func (c *Client) FindDeviceByMAC(mac string) (*resty.Response, error) {
	query := map[string]string{
		"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANIPConnection.1.MACAddress": mac,
	}
	queryJSON, _ := json.Marshal(query)

	return c.http.R().
		SetQueryParam("query", string(queryJSON)).
		Get("/devices/")
}

// DeleteDevice removes a device by ID
func (c *Client) DeleteDevice(deviceID string) (*resty.Response, error) {
	return c.http.R().
		Delete("/devices/" + url.PathEscape(deviceID))
}

// FindDevicesByLastInformBefore returns devices that haven't informed since a given time
func (c *Client) FindDevicesByLastInformBefore(timestamp string) (*resty.Response, error) {
	query := map[string]interface{}{
		"_lastInform": map[string]interface{}{
			"$lt": timestamp,
		},
	}
	return c.http.R().
		SetQueryParam("query", toJSONString(query)).
		Get("/devices")
}

// GetPendingTasksForDevice returns all tasks for a given device
func (c *Client) GetPendingTasksForDevice(deviceID string) (*resty.Response, error) {
	query := map[string]interface{}{
		"device": deviceID,
	}
	return c.http.R().
		SetQueryParam("query", toJSONString(query)).
		Get("/tasks")
}

// GetDeviceProjection returns a device with selected projection fields
func (c *Client) GetDeviceProjection(deviceID string, projection string) (*resty.Response, error) {
	query := map[string]interface{}{"_id": deviceID}
	return c.http.R().
		SetQueryParam("query", toJSONString(query)).
		SetQueryParam("projection", projection).
		Get("/devices")
}

func (c *Client) TriggerTask(deviceID string, task map[string]any) (*resty.Response, error) {
	url := fmt.Sprintf("/devices/%s/tasks?timeout=3000&connection_request", url.PathEscape(deviceID))
	return c.http.R().
		SetHeader("Content-Type", "application/json").
		SetBody(task).
		Post(url)
}
