package model

import "github.com/luuisavelino/network-interface/internal/domain/entities"

// import "github.com/luuisavelino/network-interface/internal/domain/entities"

type RoutingTableResponse struct {
	RoutingTable []RoutingResponse `json:"routing_table"`
}

type RoutingResponse struct {
	Source string  `json:"source"`
	Target string  `json:"target"`
	Weight float64 `json:"weight"`
	Type   string  `json:"type"`
}

func ToRoutingTableResponse(routingTable entities.Routing) RoutingTableResponse {
	routing := make([]RoutingResponse, 0)

	for routingType, sources := range routingTable {
		for source, targets := range sources {
			for target, weight := range targets {
				routing = append(routing, RoutingResponse{
					Source: source,
					Target: target,
					Weight: weight,
					Type:   routingType,
				})
			}
		}
	}

	return RoutingTableResponse{
		RoutingTable: routing,
	}
}
