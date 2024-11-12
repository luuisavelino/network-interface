package model

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
)

type DeviceRequest struct {
	Label        string `json:"label" binding:"required"`
	Power        int    `json:"power" binding:"required"`
}

func (r DeviceRequest) ToDomain() entities.Device {
	return entities.Device{
		Label:        r.Label,
		Power:        r.Power,
	}
}

func ToDeviceResponse(d entities.Device) DeviceResponse {
	routingTable := make([]RoutingResponse, 0)

	routingTable = ToRoutingTableResponse(d.GetRoutingTable()).RoutingTable

	return DeviceResponse{
		Label:        d.Label,
		Power:        d.Power,
		Battery:      100,
		Status:       "active",
		Requests:     ToRequestsResponse(d.Requests),
		RoutingTable: routingTable,
	}
}

func ToDevicesResponse(d entities.Devices) []DeviceResponse {
	devices := make([]DeviceResponse, 0)

	for _, device := range d {
		devices = append(devices, ToDeviceResponse(*device))
	}

	return devices
}

type DeviceResponse struct {
	Label        string            `json:"label"`
	Power        int               `json:"power"`
	Battery      int               `json:"battery"`
	Status       string            `json:"status"`
	Requests     RequestsResponse  `json:"requests"`
	RoutingTable []RoutingResponse `json:"routing_table"`
}
