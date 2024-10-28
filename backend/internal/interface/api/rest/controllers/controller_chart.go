package controllers

import (
	"net/http"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/model"
)

func (sc *apiControllerInterface) GetChart(c *gin.Context) {
	logger.Info("Init GetChart controller",
		zap.String("journey", "GetChart"),
	)

	chart, err := sc.services.Environment.GetChart(c.Request.Context())
	if err != nil {
		logger.Error("Error to get device",
			err,
			zap.String("journey", "GetChart"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to get device",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToChartResponse(chart))
}

func (sc *apiControllerInterface) SetDeviceInChart(c *gin.Context) {
	logger.Info("Init SetDeviceInChart controller",
		zap.String("journey", "SetDeviceInChart"),
	)

	deviceLabel := c.Param("deviceLabel")

	var coverageArea model.CoverageArea
	if err := c.BindJSON(&coverageArea); err != nil {
		logger.Error("Error to bind coverage area",
			err,
			zap.String("journey", "SetDeviceInChart"),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error", "message": "Error to bind coverage area",
		})

		return
	}

	sc.services.Environment.SetDeviceInChart(c.Request.Context(), deviceLabel, coverageArea.ToDomain())

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
	})
}
