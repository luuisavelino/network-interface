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
	"time"
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
	rs.ScheduleUpdateRoutingTable(label)
	// rs.ScheduleWalk(label)

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
		fmt.Println("device not found: ", deviceLabel)
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	for _, message := range device.GetUnreadMessages() {
		switch message.Topic {
		case "new-connection":
			rs.NewConnection(ctx, deviceLabel, message.Sender)
		case "new-connection-ack":
			rs.NewConnectionAck(ctx, deviceLabel, message.Sender)
		case "confirm-connection":
			rs.ConfirmConnection(ctx, deviceLabel, message.Sender)
		case "update-routing":
			rs.UpdateRouting(ctx, deviceLabel, message.Sender, message.Content.(map[uuid.UUID]entities.Routing))
		case "user-message":
			rs.UserMessage(ctx, deviceLabel, message.Sender, message.Content.(string))
		case "user-message-ack":
			rs.UserMessageAck(ctx, deviceLabel, message.Sender)
		default:
			continue
		}

		message.Read()
	}

	return nil
}

func (rs deviceService) NewConnection(ctx context.Context, current, sender string) error {
	logger.Info("Init NewConnection service",
		zap.String("journey", "NewConnection"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		fmt.Println("device not found: ", current)
		return errors.New(fmt.Sprintf("device not found: %s", current))
	}

	isNearby := rs.environment.CheckIfDeviceIsNearby(current, sender)
	if !isNearby {
		fmt.Println("device not nearby: ", sender)
		return nil
	}

	senderDevice := rs.environment.GetDeviceByLabel(sender)

	message := entities.NewMessage(
		"new-connection-ack",
		current,
		sender,
		nil,
	)

	rs.SendMessage(currentDevice, senderDevice, message)

	return nil
}

func (rs deviceService) NewConnectionAck(ctx context.Context, current, sender string) error {
	logger.Info("Init NewConnectionAck service",
		zap.String("journey", "NewConnectionAck"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", current))
	}

	senderDevice := rs.environment.GetDeviceByLabel(sender)
	if senderDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", sender))
	}

	message := entities.NewMessage(
		"confirm-connection",
		current,
		sender,
		nil,
	)

	rs.SendMessage(currentDevice, senderDevice, message)

	return nil
}

func (rs deviceService) ConfirmConnection(ctx context.Context, current, sender string) error {
	logger.Info("Init ConfirmConnection service",
		zap.String("journey", "ConfirmConnection"),
		zap.String("current", current),
		zap.String("sender", sender),
	)

	fmt.Println("confirm-connection")

	return nil
}

func (rs deviceService) UpdateRouting(ctx context.Context, current, sender string, routingTable map[uuid.UUID]entities.Routing) error {
	logger.Info("Init UpdateRouting service",
		zap.String("journey", "UpdateRouting"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", current))
	}

	currentDevice.RemoveFromTableRoutesWith(sender)
	currentDevice.AddRouting(routingTable)

	return nil
}

func (rs deviceService) UserMessage(ctx context.Context, current, sender string, message string) error {
	logger.Info("Init UserMessage service",
		zap.String("journey", "UserMessage"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", current))
	}

	senderDevice := rs.environment.GetDeviceByLabel(sender)
	if senderDevice == nil {
		return errors.New(fmt.Sprintf("device not found: %s", sender))
	}

	newMessage := entities.NewMessage(
		"user-message-ack",
		current,
		sender,
		nil,
	)

	rs.SendMessage(currentDevice, senderDevice, newMessage)

	return nil
}

func (rs deviceService) UserMessageAck(ctx context.Context, current, sender string) error {
	logger.Info("Init UserMessageAck service",
		zap.String("journey", "UserMessageAck"),
		zap.String("current", current),
	)

	fmt.Println("UserMessageAck")

	return nil
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

	currentDevice.SetScanningDevices(true)
	currentDevice.ResetDeviceConn()

	devicePosition := rs.environment.GetDeviceInChart(deviceLabel)
	if devicePosition == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	devicesNearby := rs.environment.ScanDeviceNearby(deviceLabel)
	if len(devicesNearby) == 0 {
		currentDevice.RemoveFromTableRoutesWith(deviceLabel)
		return nil
	}

	messagesSent := 0
	for _, device := range devicesNearby {
		message := entities.NewMessage(
			"new-connection",
			deviceLabel,
			device.GetDeviceLabel(),
			nil,
		)
		
		rs.SendMessage(currentDevice, device, message)

		messagesSent++
	}

	timeout := time.After(30 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			logger.Info("Timeout while waiting for devices to connect",
				zap.String("deviceLabel", deviceLabel),
			)

			rs.PropagateRoutingTable(ctx, currentDevice)

			currentDevice.SetScanningDevices(false)

			return nil
		case <-tick:
			withConn := currentDevice.GetDevicesWithConn()

			if len(withConn) == messagesSent {
				logger.Info("All devices connected",
					zap.String("deviceLabel", deviceLabel),
				)

				rs.PropagateRoutingTable(ctx, currentDevice)

				currentDevice.SetScanningDevices(false)

				return nil
			}
		}
	}

	return nil
}

func (rs deviceService) PropagateRoutingTable(ctx context.Context, device *entities.Device) {
	currDevice := device.GetDeviceLabel()
	device.RemoveFromTableRoutesWith(currDevice)

	for _, deviceConn := range device.GetDevicesWithConn() {
		connDevice := rs.environment.GetDeviceByLabel(deviceConn)
		if connDevice == nil {
			continue
		}

		weight := rs.environment.GetDistanceTo(
			rs.environment.GetDeviceInChart(currDevice).X,
			rs.environment.GetDeviceInChart(currDevice).Y,
			rs.environment.GetDeviceInChart(deviceConn).X,
			rs.environment.GetDeviceInChart(deviceConn).Y,
		)
	
		routingTable := make(map[uuid.UUID]entities.Routing, 0)
		routingTable[uuid.New()] = entities.Routing{
			Source: currDevice,
			Target: deviceConn,
			Weight: weight,
		}
	
		device.AddRouting(routingTable)

		message := entities.NewMessage(
			"update-routing",
			currDevice,
			deviceConn,
			device.GetRoutingTable(),
		)

		rs.SendMessage(device, connDevice, message)
	}
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
			fmt.Sprintf("*/5 * * * * *"),
			true,
		),
		gocron.NewTask(
			rs.ReadMessages,
			context.Background(),
			deviceLabel,
		),
	)
}
