package model

import "github.com/luuisavelino/network-interface/internal/domain/entities"

type RouteResponse struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

func ToRouteResponse(routes []entities.Route) []RouteResponse {
	routesResponse := make([]RouteResponse, 0)
	for _, route := range routes {
		routesResponse = append(routesResponse, RouteResponse{
			Source: route.Source,
			Target: route.Target,
		})
	}
	return routesResponse
}
