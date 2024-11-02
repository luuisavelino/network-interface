package controllers

import (
	"net/http"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/model"
)

func (sc *apiControllerInterface) GetDevices(c *gin.Context) {
	logger.Info("Init GetDevices controller",
		zap.String("journey", "GetDevices"),
	)

	devices, err := sc.services.Device.GetDevices(c.Request.Context())
	if err != nil {
		logger.Error("Error to get devices",
			err,
			zap.String("journey", "GetDevices"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to get devices",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToDevicesResponse(devices))
}

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

	deviceLabel := c.Param("label")

	device, err := sc.services.Device.GetDevice(c.Request.Context(), deviceLabel)
	if err != nil {
		logger.Error("Error to get device",
			err,
			zap.String("journey", "GetDevice"),
			zap.String("deviceLabel", deviceLabel),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, model.ToDeviceResponse(device))
}

func (sc *apiControllerInterface) UpdateRoutingTable(c *gin.Context) {
	logger.Info("Init UpdateRoutingTable controller",
		zap.String("journey", "UpdateRoutingTable"),
	)

	deviceLabel := c.Param("label")

	err := sc.services.Device.UpdateRoutingTable(c.Request.Context(), deviceLabel)
	if err != nil {
		logger.Error("Error to update routing table",
			err,
			zap.String("journey", "UpdateRoutingTable"),
			zap.String("deviceLabel", deviceLabel),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to update routing table",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success", "message": "Routing table updated",
	})
}

func (sc *apiControllerInterface) GetRoute(c *gin.Context) {
	logger.Info("Init GetRoute controller",
		zap.String("journey", "GetRoute"),
	)

	source := c.Param("source")
	target := c.Param("target")

	routes, err := sc.services.Device.GetRoute(c.Request.Context(), source, target)
	if err != nil {
		logger.Error("Error to get route",
			err,
			zap.String("journey", "GetRoute"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, model.ToRouteResponse(routes))
}

func (sc *apiControllerInterface) DeleteDevice(c *gin.Context) {
	logger.Info("Init DeleteDevice controller",
		zap.String("journey", "DeleteDevice"),
	)

	deviceLabel := c.Param("label")

	err := sc.services.Device.DeleteDevice(c.Request.Context(), deviceLabel)
	if err != nil {
		logger.Error("Error to delete device",
			err,
			zap.String("journey", "DeleteDevice"),
			zap.String("deviceLabel", deviceLabel),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to delete device",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success", "message": "Device deleted",
	})
}

func (sc *apiControllerInterface) SendMessage(c *gin.Context) {
	logger.Info("Init SendMessage controller",
		zap.String("journey", "SendMessage"),
	)

	var message model.MessageRequest
	if err := c.BindJSON(&message); err != nil {
		logger.Error("Error to bind message",
			err,
			zap.String("journey", "SendMessage"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to bind message",
		})

		return
	}

	err := sc.services.Device.SendUserMessage(c.Request.Context(), message.ToDomain())
	if err != nil {
		logger.Error("Error to send message",
			err,
			zap.String("journey", "SendMessage"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to send message",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success", "message": "Message sent",
	})
}