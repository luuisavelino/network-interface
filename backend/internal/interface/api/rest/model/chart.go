package model

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
)

func ToChartResponse(chart entities.Chart) ChartResponse {
	chartResponse := make(ChartResponse)

	for key, value := range chart {
		chartResponse[key] = CoverageArea{
			X: value.X,
			Y: value.Y,
			R: value.R,
		}
	}

	return chartResponse
}

type ChartResponse map[string]CoverageArea

type CoverageArea struct {
	X int     `json:"x" binding:"required"`
	Y int     `json:"y" binding:"required"`
	R float64 `json:"r"`
}

func (r CoverageArea) ToDomain() entities.CoverageArea {
	return entities.CoverageArea{
		X: r.X,
		Y: r.Y,
		R: r.R,
	}
}
