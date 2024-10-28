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
	GetDevices(ctx context.Context) (entities.Devices, error)
	InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error)
	UpdateRoutingTable(ctx context.Context, deviceLabel string) error
	GetDevice(ctx context.Context, deviceLabel string) (entities.Device, error)
	GetRoute(tx context.Context, sourceId string, targetId string) ([]entities.Route, error)
}

func (rs deviceService) GetDevices(ctx context.Context) (entities.Devices, error) {
	logger.Info("Init GetDevices service",
		zap.String("journey", "GetDevices"),
	)

	return rs.environment.GetDevices(), nil
}

func (rs deviceService) InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error) {
	logger.Info("Init InsertDevice service",
		zap.String("journey", "InsertDevice"),
	)

	device.RoutingTable = make(map[uuid.UUID]entities.Routing, 0)

	rs.environment.AddDevice(&device)

	// rs.ScheduleWalk(device.GetDeviceLabel())
	// rs.ScheduleUpdateRoutingTable(device.GetDeviceLabel())

	rs.scheduler.Start()

	return device, nil
}

func (rs deviceService) GetDevice(ctx context.Context, deviceLabel string) (entities.Device, error) {
	logger.Info("Init GetDevice service",
		zap.String("journey", "GetDevice"),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		return entities.Device{}, errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	return *device, nil
}

func (rs deviceService) UpdateRoutingTable(ctx context.Context, deviceLabel string) error {
	logger.Info("Init UpdateRoutingTable service",
		zap.String("journey", "UpdateRoutingTable"),
		zap.String("deviceLabel", deviceLabel),
	)

	currentDevice := rs.environment.GetDeviceByLabel(deviceLabel)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	fmt.Println("currentDevice :: ", currentDevice)

	routingTableToAdd := make(map[uuid.UUID]entities.Routing, 0)
	for _, message := range currentDevice.GetUnreadMessages() {
		if message.Topic == "update-routing" && message.Destination == deviceLabel {
			table := message.Content.(map[uuid.UUID]entities.Routing)
			for routeUuid, routing := range table {
				routingTableToAdd[routeUuid] = routing
			}

			message.Read()
		}
	}

	currentDevice.AddRouting(routingTableToAdd)
	currentDevice.RemoveFromTableRoutesWith(deviceLabel)

	devicePosition := rs.environment.GetDeviceInChart(deviceLabel)
	if devicePosition == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	devicesWithCommunication := rs.environment.ScanDevicesWithCommunication(deviceLabel)

	routingTable := make(map[uuid.UUID]entities.Routing, 0)
	for _, device := range devicesWithCommunication {
		position := rs.environment.GetDeviceInChart(device.Label)

		weight := rs.environment.GetDistanceTo(
			devicePosition.X, devicePosition.Y, position.X, position.Y,
		)

		routingTable[uuid.New()] = entities.Routing{
			Source: currentDevice.GetDeviceLabel(),
			Target: device.GetDeviceLabel(),
			Weight: weight,
		}
	}

	currentDevice.AddRouting(routingTable)

	if len(currentDevice.GetRoutingTable()) == 0 {
		return nil
	}


	currentDevice.PrintPrettyTable()

	for _, device := range devicesWithCommunication {
		message := entities.NewMessage(
			"update-routing",
			deviceLabel,
			device.GetDeviceLabel(),
			routingTable,
		)
		
		device.AddMessageToReceived(&message)
		currentDevice.AddMessageToSent(&message)
	}

	fmt.Println("currentDevice ADDED NEW MESSAGES", currentDevice)

	return nil
}

func (rs deviceService) GetRoute(tx context.Context, sourceId string, targetId string) ([]entities.Route, error) {
	logger.Info("Init GetRoute service",
		zap.String("journey", "GetRoute"),
	)

	sourceDevice := rs.environment.GetDeviceByLabel(sourceId)

	if sourceDevice == nil {
		return nil, errors.New(fmt.Sprintf("device not found: %s", sourceId))
	}

	graph := &entities.Graph{}
	for _, routing := range sourceDevice.GetRoutingTable() {
		graph.AddEdge(routing.Source, routing.Target, routing.Weight)
	}

	bestPaths := graph.DijkstraKBest(sourceId, targetId, 3)
	if len(bestPaths) == 0 {
		return nil, errors.New("no path found")
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

func (rs deviceService) Send(deviceLabel string) {}

func (rs deviceService) ScheduleWalk(deviceLabel string) {
	rs.scheduler.NewJob(
		gocron.CronJob(
			"*/30 * * * * *",
			true,
		),
		gocron.NewTask(
			rs.environment.Walk,
			deviceLabel,
		),
	)
}

func (rs deviceService) ScheduleUpdateRoutingTable(deviceLabel string) {
	rs.scheduler.NewJob(
		gocron.CronJob(
			fmt.Sprintf("*/%d * * * * *", rs.environment.GetDeviceByLabel(deviceLabel).MessageFreq),
			true,
		),
		gocron.NewTask(
			rs.UpdateRoutingTable,
			context.Background(),
			deviceLabel,
		),
	)
}
