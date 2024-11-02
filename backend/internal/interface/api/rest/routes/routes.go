package routes

import (
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, controller controllers.ApiControllerInterface) {
	v1 := r.Group("/api/v1")

	devices := v1.Group("/devices")
	{
			devices.GET("", controller.GetDevices)
			devices.POST("", controller.InsertDevice)
			devices.PATCH("/:label", controller.UpdateRoutingTable)
			devices.GET("/:label", controller.GetDevice)
			devices.DELETE("/:label", controller.DeleteDevice)
			devices.GET("/route/:source/:target", controller.GetRoute)
			devices.POST("/message", controller.SendMessage)
	}

	chart := v1.Group("/chart")
	{
		chart.GET("", controller.GetChart)
		chart.POST("/:deviceLabel", controller.SetDeviceInChart)
	}

	environment := v1.Group("/environment")
	{
		environment.GET("", controller.GetEnvironment)
	}
}
