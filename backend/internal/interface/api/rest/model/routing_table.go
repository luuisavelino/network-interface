package model

// import "github.com/luuisavelino/network-interface/internal/domain/entities"


type RoutingTableResponse struct {
	RoutingTable []RoutingResponse `json:"routing_table"`
}

type RoutingResponse struct {
	Source  int    `json:"source"`
	Target  int    `json:"target"`
	Weight  float64 `json:"weight"`
}

// func ToResponse(t entities.RoutingTable) RoutingTableResponse {
// 	table := make([]RouteResponse, 0)
// 	for _, route := range t.RoutingTable {
// 		table = append(table, RouteResponse{
// 			Source:   route.Source,
// 			Target:   route.Target,
// 			Weight:   route.Weight,
// 		})
// 	}
// 	return RoutingTableResponse{
// 		RoutingTable: table,
// 	}
// }
