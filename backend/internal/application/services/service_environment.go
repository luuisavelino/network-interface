package services

import (
	"context"
	"errors"

	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
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
	GetChart(ctx context.Context) (entities.Chart, error)
	SetDeviceInChart(ctx context.Context, deviceLabel string, coverageArea entities.CoverageArea)
}

func (rs environmentService) GetEnvironment(ctx context.Context) (entities.Environment, error) {
	logger.Info("Init GetEnvironment service",
		zap.String("journey", "GetEnvironment"),
	)

	return rs.environment.GetEnvironment(), nil
}

func (rs environmentService) GetChart(ctx context.Context) (entities.Chart, error) {
	logger.Info("Init GetChart service",
		zap.String("journey", "GetChart"),
	)

	return rs.environment.GetChart(), nil
}

func (rs environmentService) SetDeviceInChart(ctx context.Context, deviceLabel string, coverageArea entities.CoverageArea) {
	logger.Info("Init SetDeviceInChart service",
		zap.String("journey", "SetDeviceInChart"),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		logger.Error("Device not found",
			errors.New("Device not found"),
			zap.String("deviceLabel", deviceLabel),
		)
		return
	}

	coverageArea.R = float64(device.Power)

	rs.environment.SetDeviceInChart(deviceLabel, coverageArea)
}
