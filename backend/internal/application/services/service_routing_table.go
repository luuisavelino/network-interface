package services

import (
	"context"

	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
)

func NewRoutingTableService(environment *entities.Environment) RoutingTableService {
	return routingTableService{
		environment: environment,
	}
}

type routingTableService struct {
	environment *entities.Environment
}

type RoutingTableService interface {
	GetTable(ctx context.Context) ([]entities.Routing, error)
}

func (rs routingTableService) GetTable(ctx context.Context) ([]entities.Routing, error) {
	logger.Info("Init GetTable service",
		zap.String("journey", "GetTable"),
	)

	return nil, nil
}
