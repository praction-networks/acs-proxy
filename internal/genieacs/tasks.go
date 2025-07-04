package genieacs

import (
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/praction-networks/acs-proxy/internal/models"
)

// RefreshObject triggers a full parameter refresh
func (c *Client) RefreshObject(deviceID string) (*resty.Response, error) {
	payload := map[string]interface{}{
		"name":       "refreshObject",
		"objectName": "",
	}
	return c.http.R().
		SetBody(payload).
		Post("/devices/" + url.PathEscape(deviceID) + "/tasks?connection_request")
}

// SetWiFiCredentials sets SSID and password
func (c *Client) SetWiFiCredentials(wifiCred *models.SetWirelessCred) (*resty.Response, error) {
	ssidPrefix := "BitFiber_" + wifiCred.WirelessUsername
	password := wifiCred.WirelessPassword
	manufacturer := strings.ToUpper(wifiCred.Manufacturer)

	var parameterValues [][]interface{}

	switch manufacturer {
	case "HWTC", "REALTEK":
		parameterValues = [][]interface{}{
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
		}

	case "DRAGONPATH", "MONU", "ASFT", "DIXON":
		parameterValues = [][]interface{}{
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
		}

	case "ADOPT", "PON":
		parameterValues = [][]interface{}{
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.6.SSID", ssidPrefix, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.6.PreSharedKey.1.KeyPassphrase", password, "xsd:string"},
		}

	default:
		parameterValues = [][]interface{}{
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.SSID", ssidPrefix + "_2.4G", "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.1.PreSharedKey.1.PreSharedKey", password, "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.SSID", ssidPrefix + "_5G", "xsd:string"},
			{"InternetGatewayDevice.LANDevice.1.WLANConfiguration.5.PreSharedKey.1.PreSharedKey", password, "xsd:string"},
		}
	}

	payload := map[string]interface{}{
		"name":            "setParameterValues",
		"parameterValues": parameterValues,
	}

	return c.http.R().
		SetBody(payload).
		Post("/devices/" + url.PathEscape(wifiCred.DeviceID) + "/tasks?connection_request")
}

func (c *Client) SetPPPoECredentials(credentials *models.SetPPPoECred) (*resty.Response, error) {
	// Step 1: Set PPPoE username and password
	payload := map[string]interface{}{
		"name": "setParameterValues",
		"parameterValues": [][]interface{}{
			{"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANPPPConnection.1.Username", credentials.PPPoEUsername, "xsd:string"},
			{"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANPPPConnection.1.Password", credentials.PPPoEPassword, "xsd:string"},
		},
	}

	resp, err := c.http.R().
		SetBody(payload).
		Post("/devices/" + url.PathEscape(credentials.DeviceID) + "/tasks?connection_request")
	if err != nil {
		return resp, err
	}

	// Step 2: Bounce the connection
	var bouncePayload map[string]interface{}

	if credentials.Manufacturer == "ASFT" {
		// Send reboot task
		bouncePayload = map[string]interface{}{
			"name": "reboot",
		}
	} else {
		// Send WANPPPConnection reset
		bouncePayload = map[string]interface{}{
			"name": "setParameterValues",
			"parameterValues": [][]interface{}{
				{"InternetGatewayDevice.WANDevice.1.WANConnectionDevice.1.WANPPPConnection.1.Reset", true, "xsd:boolean"},
			},
		}
	}

	bounceResp, bounceErr := c.http.R().
		SetBody(bouncePayload).
		Post("/devices/" + url.PathEscape(credentials.DeviceID) + "/tasks?timeout=3000&connection_request")

	return bounceResp, bounceErr
}

// RebootDevice issues a reboot task
func (c *Client) RebootDevice(deviceID string) (*resty.Response, error) {
	payload := map[string]interface{}{
		"name": "reboot",
	}
	return c.http.R().
		SetBody(payload).
		Post("/devices/" + url.PathEscape(deviceID) + "/tasks?connection_request")
}

func (c *Client) RetryTask(taskID string) (*resty.Response, error) {
	return c.http.R().Post("/tasks/" + url.PathEscape(taskID) + "/retry")
}

func (c *Client) DeleteTask(taskID string) (*resty.Response, error) {
	return c.http.R().Delete("/tasks/" + url.PathEscape(taskID))
}
