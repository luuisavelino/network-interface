package main

import (
	"fmt"
	"log"
	"time"

	"github.com/luuisavelino/network-interface/internal/application/services"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/controllers"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/middleware"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/routes"
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/envs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
)

func main() {
	environment := entities.NewEnvironment()

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}

	services := services.ApiServices{
		RoutingTable:  services.NewRoutingTableService(environment),
		Device:        services.NewDeviceService(environment, s),
		Environment:   services.NewEnvironmentService(environment),
	}

	apiController := controllers.NewApiController(services)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PATCH"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())
	routes.InitRoutes(&router.RouterGroup, apiController)

	if err := router.Run(fmt.Sprintf(":%d", envs.Api.Port)); err != nil {
		log.Fatal(err)
	}
}
