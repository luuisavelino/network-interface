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

	label := device.GetDeviceLabel()

	rs.ScheduleReadMessages(label)
	// rs.ScheduleWalk(label)
	// rs.ScheduleUpdateRoutingTable(label)

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

func (rs deviceService) ReadMessages(ctx context.Context, deviceLabel string) error {
	logger.Info("Init ReadMessages service",
		zap.String("journey", "ReadMessages"),
		zap.String("deviceLabel", deviceLabel),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	routingTableToAdd := make(map[uuid.UUID]entities.Routing, 0)

	for _, message := range device.GetUnreadMessages() {
		switch message.Topic {
		case "new-connection":
			fmt.Println("new-connection")
			// check if the device is in the chart
			// answer with new-connection-ack
		case "new-connection-ack":
			// add the device to the routing table
			// send confirm-connection
			fmt.Println("new-connection-ack")
		case "confirm-connection":
			// add the device to the routing table
			// send confirm-connection-ack
			fmt.Println("confirm-connection")
		case "update-routing":
			table := message.Content.(map[uuid.UUID]entities.Routing)
			for routeUuid, routing := range table {
				routingTableToAdd[routeUuid] = routing
			}
		case "user-message":
			// send user-message-ack
			fmt.Println("user-message")
		case "user-message-ack":
			fmt.Println("user-message-ack")
		default:
			continue
		}

		message.Read()
	}

	device.AddRouting(routingTableToAdd)

	return nil
}

// "update-routing"
// "user-message"
// "user-message-ack"

func (rs deviceService) NewConnection(ctx context.Context, deviceLabel, target string) error {
	logger.Info("Init NewConnection service",
		zap.String("journey", "NewConnection"),
		zap.String("deviceLabel", deviceLabel),
	)

	currentDevice := rs.environment.GetDeviceByLabel(deviceLabel)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	isNearby := rs.environment.CheckIfDeviceIsNearby(deviceLabel, target)
	if !isNearby {
		return nil
	}

	message := entities.NewMessage(
		"new-connection-ack",
		deviceLabel,
		target,
		nil,
	)

	rs.environment.SendMessage(currentDevice, device, message)
	return nil
}

func (rs deviceService) NewConnectionAck(ctx context.Context, deviceLabel, target string) error {
	logger.Info("Init NewConnectionAck service",
		zap.String("journey", "NewConnectionAck"),
		zap.String("deviceLabel", deviceLabel),
	)

	currentDevice := rs.environment.GetDeviceByLabel(deviceLabel)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	targetDevice := rs.environment.GetDeviceByLabel(target)
	if targetDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", target))
	}

	weight := rs.environment.GetDistanceTo(
		rs.environment.GetDeviceInChart(deviceLabel).X,
		rs.environment.GetDeviceInChart(deviceLabel).Y,
		rs.environment.GetDeviceInChart(target).X,
		rs.environment.GetDeviceInChart(target).Y,
	)

	routingTable := make(map[uuid.UUID]entities.Routing, 0)
	routingTable[uuid.New()] = entities.Routing{
		Source: deviceLabel,
		Target: target,
		Weight: weight,
	}

	currentDevice.AddRouting(routingTable)

	message := entities.NewMessage(
		"confirm-connection",
		deviceLabel,
		target,
		routingTable,
	)

	rs.environment.SendMessage(currentDevice, targetDevice, message)

	return nil
}

func (rs deviceService) ConfirmConnection(ctx context.Context, deviceLabel, target string, routingTable map[uuid.UUID]entities.Routing) error {
	logger.Info("Init ConfirmConnection service",
		zap.String("journey", "ConfirmConnection"),
		zap.String("deviceLabel", deviceLabel),
	)

	currentDevice := rs.environment.GetDeviceByLabel(deviceLabel)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	currentDevice.AddRouting(routingTable)
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

	currentDevice.RemoveFromTableRoutesWith(deviceLabel)

	devicePosition := rs.environment.GetDeviceInChart(deviceLabel)
	if devicePosition == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	devicesNearby := rs.environment.ScanDeviceNearby(deviceLabel)
	if len(devicesNearby) == 0 {
		return nil
	}

	for _, device := range devicesNearby {
		message := entities.NewMessage(
			"new-connection",
			deviceLabel,
			device.GetDeviceLabel(),
			routingTable,
		)
		
		rs.environment.SendMessage(currentDevice, device, message)
	}

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

func (rs deviceService) SendMessage(currentDevice *entities.Device, targetDevice *entities.Device, message entities.Message) {
	logger.Info("Init SendMessage service",
		zap.String("journey", "SendMessage"),
		zap.String("source", currentDevice.GetDeviceLabel()),
		zap.String("target", targetDevice.GetDeviceLabel()),
		zap.String("message", message.Topic),
	)

	targetDevice.AddMessageToReceived(&message)
	currentDevice.AddMessageToSent(&message)
}

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

func (rs deviceService) ScheduleReadMessages(deviceLabel string) {
	rs.scheduler.NewJob(
		gocron.CronJob(
			fmt.Sprintf("*/10 * * * * *"),
			true,
		),
		gocron.NewTask(
			rs.ReadMessages,
			context.Background(),
			deviceLabel,
		),
	)
}
