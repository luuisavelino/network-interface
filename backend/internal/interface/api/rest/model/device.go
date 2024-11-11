package model

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
)

type DeviceRequest struct {
	Label        string `json:"label" binding:"required"`
	Power        int    `json:"power" binding:"required"`
	WalkingSpeed int    `json:"walking_speed" binding:"required"`
	MessageFreq  int    `json:"message_freq" binding:"required"`
}

func (r DeviceRequest) ToDomain() entities.Device {
	return entities.Device{
		Label:        r.Label,
		Power:        r.Power,
		WalkingSpeed: r.WalkingSpeed,
		MessageFreq:  r.MessageFreq,
	}
}

func ToDeviceResponse(d entities.Device) DeviceResponse {
	routingTable := make([]RoutingResponse, 0)

	for _, route := range d.RoutingTable {
		routingTable = append(routingTable, RoutingResponse{
			Source: route.Source,
			Target: route.Target,
			Weight: route.Weight,
		})
	}

	return DeviceResponse{
		Label:        d.Label,
		Power:        d.Power,
		Battery:      100,
		Status:       "active",
		Messages:     ToMessagesResponse(d.Messages),
		WalkingSpeed: d.WalkingSpeed,
		MessageFreq:  d.MessageFreq,
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
	Messages     MessagesResponse  `json:"messages"`
	WalkingSpeed int               `json:"walking_speed"`
	MessageFreq  int               `json:"message_freq"`
	RoutingTable []RoutingResponse `json:"routing_table"`
}
