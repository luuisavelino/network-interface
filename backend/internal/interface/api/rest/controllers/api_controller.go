package controllers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/luuisavelino/network-interface/internal/application/services"
)

type apiControllerInterface struct {
	services services.ApiServices
}

type ApiControllerInterface interface {
	RoutingsControllerInterface
	DevicesControllerInterface
	EnvironmentControllerInterface
	ChartControllerInterface
}

type RoutingsControllerInterface interface {
	GetTable(c *gin.Context)
}

type DevicesControllerInterface interface {
	GetDevices(c *gin.Context)
	InsertDevice(c *gin.Context)
	GetDevice(c *gin.Context)
	UpdateRoutingTable(c *gin.Context)
	DeleteDevice(c *gin.Context)
	GetRoute(c *gin.Context)
	SendRequest(c *gin.Context)
}

type EnvironmentControllerInterface interface {
	GetEnvironment(c *gin.Context)
}

type ChartControllerInterface interface {
	GetChart(c *gin.Context)
	SetDeviceInChart(c *gin.Context)
}

func NewApiController(apiServices services.ApiServices) ApiControllerInterface {
	return &apiControllerInterface{
		services: apiServices,
	}
}
