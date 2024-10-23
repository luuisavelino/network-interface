package services

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
	"context"
	// "fmt"
)

func NewEnvironmentService(environment entities.Environment) EnvironmentService {
	return environmentService{
		environment: environment,
	}
}

type environmentService struct {
	environment entities.Environment
}

type EnvironmentService interface {
	GetEnvironment(ctx context.Context) (entities.Environment, error)
}

func (rs environmentService) GetEnvironment(ctx context.Context) (entities.Environment, error) {
	logger.Info("Init GetEnvironment service",
		zap.String("journey", "GetEnvironment"),
	)

	return rs.environment.GetEnvironment(), nil
}
