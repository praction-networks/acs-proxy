package models

import "time"

type DeviceID struct {
	Manufacturer string `json:"_Manufacturer" bson:"_Manufacturer"`
	OUI          string `json:"_OUI" bson:"_OUI"`
	ProductClass string `json:"_ProductClass" bson:"_ProductClass"`
	SerialNumber string `json:"_SerialNumber" bson:"_SerialNumber"`
}

type DeviceModel struct {
	ID                string                 `json:"_id" bson:"_id"`
	InternetGateway   map[string]interface{} `json:"InternetGatewayDevice" bson:"InternetGatewayDevice"`
	FactoryReset      map[string]interface{} `json:"FactoryReset" bson:"FactoryReset"`
	Reboot            map[string]interface{} `json:"Reboot" bson:"Reboot"`
	VirtualParameters map[string]interface{} `json:"VirtualParameters" bson:"VirtualParameters"`
	DeviceID          DeviceID               `json:"_deviceId" bson:"_deviceId"`
	LastInform        time.Time              `json:"_lastInform" bson:"_lastInform"`
	Registered        time.Time              `json:"_registered" bson:"_registered"`
	LastBoot          time.Time              `json:"_lastBoot" bson:"_lastBoot"`
	Timestamp         time.Time              `json:"_timestamp" bson:"_timestamp"`
}

type DeviceSearchID struct {
	DeviceSN string `json:"SerialNumber" validate:"required,alphanum,min=4,max=16"`
}

type SetPPPoECred struct {
	DeviceID      string `json:"deviceID" validate:"required"`
	Manufacturer  string `json:"manufacturer" validate:"required"`
	PPPoEUsername string `json:"PPPoEUsername" validate:"required,min=2,max=100,singleword"`
	PPPoEPassword string `json:"PPPoEPassword" validate:"required,min=2,max=100,singleword"`
}

type SetWirelessCred struct {
	DeviceID         string `json:"deviceID" validate:"required"`
	Manufacturer     string `json:"manufacturer" validate:"required"`
	WirelessUsername string `json:"WirelessUsername" validate:"required,min=8,max=100,singleword"`
	WirelessPassword string `json:"WirelessPassword" validate:"required,min=8,max=100,singleword"`
}
