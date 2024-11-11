package model

import "github.com/luuisavelino/network-interface/internal/domain/entities"

type EnvironmentRequest struct{}

func (r EnvironmentRequest) ToDomain() entities.Environment {
	return entities.Environment{}
}

func ToEnvironmentResponse(d entities.Environment) EnvironmentResponse {
	var deviceResponse []DeviceResponse
	chartResponse := ToChartResponse(d.GetChart())

	for _, device := range d.Devices {
		deviceResponse = append(deviceResponse, ToDeviceResponse(*device))
	}

	return EnvironmentResponse{
		Devices: deviceResponse,
		Chart:   chartResponse,
	}
}

type EnvironmentResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Chart   ChartResponse    `json:"chart"`
}
