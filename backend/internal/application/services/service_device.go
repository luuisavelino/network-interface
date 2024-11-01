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
		jobs: Jobs{
			UpdateRoutingTable: make(map[string]*gocron.Job),
			Message: make(map[string]*gocron.Job),
			Walk: make(map[string]*gocron.Job),
		},
	}
}

type deviceService struct {
	environment entities.Environment
	scheduler gocron.Scheduler
	jobs Jobs
}

type Jobs struct {
	UpdateRoutingTable map[string]*gocron.Job
	Message map[string]*gocron.Job
	Walk map[string]*gocron.Job
}

type DeviceService interface {
	GetDevices(ctx context.Context) (entities.Devices, error)
	InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error)
	UpdateRoutingTable(ctx context.Context, deviceLabel string) error
	GetDevice(ctx context.Context, deviceLabel string) (entities.Device, error)
	GetRoute(tx context.Context, sourceId string, targetId string) ([]entities.Route, error)
	DeleteDevice(ctx context.Context, deviceLabel string) error
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

	device.SetStatus(true)

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

	currentDevice.SetDeviceWithConn(sender)

	return nil
}

func (rs deviceService) ConfirmConnection(ctx context.Context, current, sender string) error {
	logger.Info("Init ConfirmConnection service",
		zap.String("journey", "ConfirmConnection"),
		zap.String("current", current),
		zap.String("sender", sender),
	)

	// currentDevice := rs.environment.GetDeviceByLabel(current)
	// if currentDevice == nil {
	// 	return errors.New(fmt.Sprintf("device not found: %s", current))
	// }

	// senderDevice := rs.environment.GetDeviceByLabel(sender)
	// if senderDevice == nil {
	// 	return errors.New(fmt.Sprintf("device not found: %s", sender))
	// }

	// currentDevice.SetDeviceWithConn(sender)

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

	for _, device := range devicesNearby {
		message := entities.NewMessage(
			"new-connection",
			deviceLabel,
			device.GetDeviceLabel(),
			nil,
		)
		
		rs.SendMessage(currentDevice, device, message)
	}

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-timeout:
			logger.Info("Timeout while waiting for devices to connect",
				zap.String("deviceLabel", deviceLabel),
			)

			rs.PropagateRoutingTable(ctx, currentDevice)

			currentDevice.SetScanningDevices(false)

			return nil
		}
	}

	return nil
}

func (rs deviceService) PropagateRoutingTable(ctx context.Context, device *entities.Device) {
	currDevice := device.GetDeviceLabel()
	device.RemoveFromTableRoutesWith(currDevice)

	fmt.Println("PropagateRoutingTable", device.GetDevicesWithConn())

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

		device.PrintPrettyTable()

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
	logger.Info("Init ScheduleWalk service",
		zap.String("journey", "ScheduleWalk"),
		zap.String("deviceLabel", deviceLabel),
	)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rs.environment.Walk(deviceLabel)
			}
		}
	}()
}

func (rs deviceService) ScheduleUpdateRoutingTable(deviceLabel string) {
	logger.Info("Init ScheduleUpdateRoutingTable service",
		zap.String("journey", "ScheduleUpdateRoutingTable"),
		zap.String("deviceLabel", deviceLabel),
	)

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rs.UpdateRoutingTable(context.Background(), deviceLabel)
			}
		}
	}()
}

func (rs deviceService) ScheduleReadMessages(deviceLabel string) {
	logger.Info("Init ScheduleReadMessages service",
		zap.String("journey", "ScheduleReadMessages"),
		zap.String("deviceLabel", deviceLabel),
	)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rs.ReadMessages(context.Background(), deviceLabel)
			}
		}
	}()
}

func (rs *deviceService) CancelJob(jobType, deviceLabel string) {
	switch jobType {
	case "message":
		fmt.Println("cancel message")
    // if job, exists := rs.jobs.Message[deviceLabel]; exists {
		// 	job.Remove(rs.ReadMessages)
		// 	delete(rs.jobs.Message, deviceLabel)
		// }
	case "walk":
		fmt.Println("cancel walk")
		// if job, exists := rs.jobs.Walk[deviceLabel]; exists {
		// 	job.Remove(rs.environment.Walk)
		// 	delete(rs.jobs.Walk, deviceLabel)
		// }
	case "updateRoutingTable":
		fmt.Println("cancel updateRoutingTable")
		// if job, exists := rs.jobs.UpdateRoutingTable[deviceLabel]; exists {
		// 	job.Remove(rs.ScheduleReadMessages)
		// 	delete(rs.jobs.UpdateRoutingTable, deviceLabel)
		// }
	default:
		return
	}
}

func (rs deviceService) DeleteDevice(ctx context.Context, deviceLabel string) error {
	logger.Info("Init DeleteDevice service",
		zap.String("journey", "DeleteDevice"),
		zap.String("deviceLabel", deviceLabel),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	rs.CancelJob("message", deviceLabel)
	rs.CancelJob("walk", deviceLabel)
	rs.CancelJob("updateRoutingTable", deviceLabel)

	rs.environment.RemoveDevice(deviceLabel)

	device.ResetDeviceConn()
	device.ResetRoutingTable()

	device.SetScanningDevices(false)
	device.SetStatus(false)

	return nil
}
