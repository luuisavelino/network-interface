package model

import "github.com/luuisavelino/network-interface/internal/domain/entities"

type DeviceRequest struct {
	ID           int     `json:"id" binding:"required"`
	Power        int     `json:"power" binding:"required"`
	PosX         int     `json:"pos_x" binding:"required"`
	PosY         int     `json:"pos_y" binding:"required"`
	WalkingSpeed int     `json:"walking_speed" binding:"required"`
	MessageFreq  int     `json:"message_freq" binding:"required"`
}

func (r DeviceRequest) ToDomain() entities.Device {
	return entities.Device{
		ID:           r.ID,
		Power:        r.Power,
		PosX:         r.PosX,
		PosY:         r.PosY,
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
		ID:           d.ID,
		Power:        d.Power,
		PosX:         d.PosX,
		PosY:         d.PosY,
		WalkingSpeed: d.WalkingSpeed,
		MessageFreq:  d.MessageFreq,
		RoutingTable: routingTable,
	}
}

type DeviceResponse struct {
	ID           int     `json:"label"`
	Power        int     `json:"r"`
	PosX         int     `json:"x"`
	PosY         int     `json:"y"`
	WalkingSpeed int     `json:"walking_speed"`
	MessageFreq  int     `json:"message_freq"`
	RoutingTable []RoutingResponse `json:"routing_table"`
}
