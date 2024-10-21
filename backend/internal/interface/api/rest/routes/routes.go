package routes

import (
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup, controller controllers.ApiControllerInterface) {
	v1 := r.Group("/api/v1")

	devices := v1.Group("/devices")
	{
			devices.POST("", controller.InsertDevice)
			devices.PATCH("/:id", controller.UpdateRoutingTable)
			devices.GET("/:id", controller.GetDevice)
			devices.GET("/route/:sourceId/:targetId", controller.GetRoute) // Changed the route to avoid conflict
	}

	environment := v1.Group("/environment")
	{
		environment.GET("", controller.GetEnvironment)
	}
}
