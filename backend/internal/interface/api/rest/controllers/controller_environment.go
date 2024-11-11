package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luuisavelino/network-interface/internal/interface/api/rest/model"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
)

func (sc *apiControllerInterface) GetEnvironment(c *gin.Context) {
	logger.Info("Init GetEnvironment controller",
		zap.String("journey", "GetEnvironment"),
	)

	environment, err := sc.services.Environment.GetEnvironment(c.Request.Context())
	if err != nil {
		logger.Error("Error to get device",
			err,
			zap.String("journey", "GetEnvironment"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to get device",
		})

		return
	}

	c.JSON(http.StatusOK, model.ToEnvironmentResponse(environment))
}
