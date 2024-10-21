package controllers

import (
	"net/http"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (sc *apiControllerInterface) GetTable(c *gin.Context) {
	logger.Info("Init GetTable controller",
		zap.String("journey", "GetTable"),
	)

	tableList, err := sc.services.RoutingTable.GetTable(c.Request.Context())
	if err != nil {
		logger.Error("Error to list tables",
			err,
			zap.String("journey", "GetTable"),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error", "message": "Error to list tables",
		})

		return
	}

	c.JSON(http.StatusOK, tableList)
}
