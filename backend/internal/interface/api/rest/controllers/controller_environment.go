package controllers

import (
	"net/http"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/model"
)

func (sc *apiControllerInterface) GetEnvironment(c *gin.Context) {
	logger.Info("Init GetEnvironment controller",
		zap.String("journey", "GetEnvironment"),
	)

	environment, err := sc.services.Environment.GetEnvironment(c.Request.Context())
	if err != nil {
		logger.Error("Error to insert device",
			err,
			zap.String("journey", "GetEnvironment"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to insert device",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToEnvironmentResponse(environment))
}
