package services

import (
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
	"context"
	"github.com/google/uuid"
	"github.com/go-co-op/gocron/v2"
	"fmt"
	"errors"
)

func NewDeviceService(environment entities.Environment, scheduler gocron.Scheduler) DeviceService {
	return deviceService{
		environment: environment,
		scheduler: scheduler,
	}
}

type deviceService struct {
	environment entities.Environment
	scheduler gocron.Scheduler
}

type DeviceService interface {
	InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error)
	UpdateRoutingTable(ctx context.Context, deviceId int)
	GetDevice(ctx context.Context, deviceId int) (entities.Device, error)
	GetRoute(tx context.Context, sourceId int, targetId int) ([]entities.Route, error)
}

func (rs deviceService) InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error) {
	logger.Info("Init InsertDevice service",
		zap.String("journey", "InsertDevice"),
	)

	device.RoutingTable = make(map[uuid.UUID]entities.Routing, 0)

	rs.environment.AddDevice(&device)

	// rs.ScheduleWalk(device.GetDeviceID())
	// rs.ScheduleUpdateRoutingTable(device.GetDeviceID())

	rs.scheduler.Start()

	return device, nil
}

func (rs deviceService) GetDevice(ctx context.Context, deviceId int) (entities.Device, error) {
	logger.Info("Init GetDevice service",
		zap.String("journey", "GetDevice"),
	)

	return *rs.environment.GetDeviceById(deviceId), nil
}

func (rs deviceService) UpdateRoutingTable(ctx context.Context, deviceId int) () {
	logger.Info("Init UpdateRoutingTable service",
		zap.String("journey", "UpdateRoutingTable"),
		zap.Int("deviceId", deviceId),
	)

	currentDevice := rs.environment.GetDeviceById(deviceId)

	routingTableToAdd := make(map[uuid.UUID]entities.Routing, 0)
	for _, message := range currentDevice.GetUnreadMessages() {
		if message.Topic == "update-routing" && message.Destination == deviceId {
			table := message.Content.(map[uuid.UUID]entities.Routing)
			for routeUuid, routing := range table {
				routingTableToAdd[routeUuid] = routing
			}

			message.Read()
		}
	}

	currentDevice.RemoveFromTableRoutesWith(deviceId)

	devicesWithCommunication := rs.environment.ScanDevicesWithCommunication(deviceId)
	routingTable := make(map[uuid.UUID]entities.Routing, 0)
	for _, device := range devicesWithCommunication {
		weight := currentDevice.GetDistanceTo(device.PosX, device.PosY)
		routingTable[uuid.New()] = entities.Routing{
			Source: currentDevice.GetDeviceID(),
			Target: device.GetDeviceID(),
			Weight: weight,
		}
	}

	currentDevice.AddRouting(routingTable)

	if len(currentDevice.GetRoutingTable()) == 0 {
		return
	}

	for _, device := range devicesWithCommunication {
		message := entities.NewMessage(
			"update-routing",
			deviceId,
			device.GetDeviceID(),
			routingTable,
		)
		
		device.AddMessageToReceived(&message)
		currentDevice.AddMessageToSent(&message)
	}
}

func (rs deviceService) GetRoute(tx context.Context, sourceId int, targetId int) ([]entities.Route, error) {
	logger.Info("Init GetRoute service",
		zap.String("journey", "GetRoute"),
	)

	sourceDevice := rs.environment.GetDeviceById(sourceId)

	if sourceDevice == nil {
		logger.Error("Source device not found",
			errors.New("source device not found"),
			zap.String("journey", "GetRoute"),
		)
		return nil, nil
	}

	graph := &entities.Graph{}
	for _, routing := range sourceDevice.GetRoutingTable() {
		graph.AddEdge(routing.Source, routing.Target, routing.Weight)
	}

	bestPaths := graph.DijkstraKBest(sourceId, targetId, 3)
	if len(bestPaths) == 0 {
		logger.Error("No path found",
			errors.New("no path found"),
			zap.String("journey", "GetRoute"),
		)
		return nil, nil
	}

	routes := make([]entities.Route, 0)

	for _, route := range bestPaths[0].Path {
		routes = append(routes, entities.Route{
			Source: route.Source,
			Target: route.Target,
		})
	}

	return routes, nil
}

func (rs deviceService) Send(deviceId int) {

}


func (rs deviceService) ScheduleWalk(deviceId int) {
	rs.scheduler.NewJob(
		gocron.CronJob(
			"*/30 * * * * *",
			true,
		),
		gocron.NewTask(
			rs.environment.GetDeviceById(deviceId).Walk,
		),
	)
}

func (rs deviceService) ScheduleUpdateRoutingTable(deviceId int) {
	rs.scheduler.NewJob(
		gocron.CronJob(
			fmt.Sprintf("*/%d * * * * *", rs.environment.GetDeviceById(deviceId).MessageFreq),
			true,
		),
		gocron.NewTask(
			rs.UpdateRoutingTable,
			context.Background(),
			deviceId,
		),
	)
}