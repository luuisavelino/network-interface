package controllers

import (
	"net/http"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/model"
	"strconv"
)

func (sc *apiControllerInterface) InsertDevice(c *gin.Context) {
	logger.Info("Init InsertDevice controller",
		zap.String("journey", "InsertDevice"),
	)

	var device model.DeviceRequest
	if err := c.BindJSON(&device); err != nil {
		logger.Error("Error to bind device",
			err,
			zap.String("journey", "InsertDevice"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to bind device",
		})

		return
	}

	deviceCreated, err := sc.services.Device.InsertDevice(c.Request.Context(), device.ToDomain())
	if err != nil {
		logger.Error("Error to insert device",
			err,
			zap.String("journey", "InsertDevice"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to insert device",
		})

		return
	}

	c.JSON(http.StatusCreated, model.ToDeviceResponse(deviceCreated))
}

func (sc *apiControllerInterface) GetDevice(c *gin.Context) {
	logger.Info("Init GetDevice controller",
		zap.String("journey", "GetDevice"),
	)

	id := c.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error to get device id",
			err,
			zap.String("journey", "GetDevice"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to get device id",
		})

		return
	}

	device, err := sc.services.Device.GetDevice(c.Request.Context(), deviceId)
	if err != nil {
		logger.Error("Error to get device",
			err,
			zap.String("journey", "GetDevice"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to get device",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToDeviceResponse(device))
}

func (sc *apiControllerInterface) UpdateRoutingTable(c *gin.Context) {
	logger.Info("Init UpdateRoutingTable controller",
		zap.String("journey", "UpdateRoutingTable"),
	)

	id := c.Param("id")
	deviceId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error to get device id",
			err,
			zap.String("journey", "UpdateRoutingTable"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to get device id",
		})

		return
	}

	sc.services.Device.UpdateRoutingTable(c.Request.Context(), deviceId)

	c.JSON(http.StatusOK, gin.H{
		"status": "success", "message": "Routing table updated",
	})
}

func (sc *apiControllerInterface) GetRoute(c *gin.Context) {
	logger.Info("Init GetRoute controller",
		zap.String("journey", "GetRoute"),
	)

	id := c.Param("sourceId")
	sourceId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error to get source id",
			err,
			zap.String("journey", "GetRoute"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to get device id",
		})

		return
	}

	id = c.Param("targetId")
	targetId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Error to get target id",
			err,
			zap.String("journey", "GetRoute"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to get device id",
		})

		return
	}

	routes, err := sc.services.Device.GetRoute(c.Request.Context(), sourceId, targetId)
	if err != nil {
		logger.Error("Error to get route",
			err,
			zap.String("journey", "GetRoute"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to get route",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToRouteResponse(routes))
}
