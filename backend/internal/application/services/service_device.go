package services

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/RyanCarrier/dijkstra/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
)

func NewDeviceService(environment entities.Environment, scheduler gocron.Scheduler) DeviceService {
	return deviceService{
		environment: environment,
		scheduler:   scheduler,
		jobs: Jobs{
			UpdateRoutingTable: make(map[string]context.CancelFunc),
			Message:            make(map[string]context.CancelFunc),
			Walk:               make(map[string]context.CancelFunc),
		},
	}
}

type deviceService struct {
	environment entities.Environment
	scheduler   gocron.Scheduler
	jobs        Jobs
}

type Jobs struct {
	UpdateRoutingTable map[string]context.CancelFunc
	Message            map[string]context.CancelFunc
	Walk               map[string]context.CancelFunc
}

type DeviceService interface {
	GetDevices(ctx context.Context) (entities.Devices, error)
	InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error)
	UpdateRoutingTable(ctx context.Context, deviceLabel string) error
	GetDevice(ctx context.Context, deviceLabel string) (entities.Device, error)
	GetRoute(tx context.Context, sourceId string, targetId string) ([]entities.Route, error)
	DeleteDevice(ctx context.Context, deviceLabel string) error
	SendUserMessage(ctx context.Context, message entities.Message) error
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
	device.Messages = entities.Messages{
		Sent:     make(map[uuid.UUID]*entities.Message),
		Received: make(map[uuid.UUID]*entities.Message),
	}

	device.SetStatus(true)

	rs.environment.AddDevice(&device)

	label := device.GetDeviceLabel()

	rs.ScheduleReadMessages(label)
	rs.ScheduleUpdateRoutingTable(label)
	rs.ScheduleWalk(label)

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

func (rs deviceService) SendUserMessage(ctx context.Context, message entities.Message) error {
	logger.Info("Init SentMessage service",
		zap.String("journey", "SentMessage"),
		zap.String("current", message.Sender),
		zap.String("target", message.Destination),
		zap.String("message", message.Content.(string)),
	)

	currentDevice := rs.environment.GetDeviceByLabel(message.Sender)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", message.Sender)
	}

	targetDevice := rs.environment.GetDeviceByLabel(message.Destination)
	if targetDevice == nil {
		return fmt.Errorf("device not found: %s", message.Destination)
	}

	newMessage := entities.NewMessage(
		"user-message",
		message.Sender,
		message.Destination,
		message.Content,
	)

	go func() {
		const maxRetries = 3
		const maxTimeout = 5 * time.Second

		for i := 0; i < maxRetries; i++ {
			rs.SendMessage(currentDevice, targetDevice, newMessage)

			select {
			case <-ctx.Done():
				return
			case <-time.After(maxTimeout):
				for _, msg := range targetDevice.GetUnreadMessages() {
					if msg.Topic == "user-message-ack" && msg.Sender == message.Destination {
						// currentDevice.DeleteMessage(msg.ID)
						msg.Read()
						return
					}
				}
			}
		}
		logger.Info("Failed to receive user-message-ack after retries",
			zap.String("current", message.Sender),
			zap.String("target", message.Destination),
		)
	}()

	return nil
}

func (rs deviceService) ReadMessages(ctx context.Context, deviceLabel string) error {
	printAlloc()

	logger.Info("Init ReadMessages service",
		zap.String("journey", "ReadMessages"),
		zap.String("deviceLabel", deviceLabel),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		fmt.Println("device not found: ", deviceLabel)
		return errors.New(fmt.Sprintf("device not found: %s", deviceLabel))
	}

	rs.ProcessAutomaticMessages(ctx, device)

	return nil
}

func (rs deviceService) ProcessAutomaticMessages(ctx context.Context, device *entities.Device) error {
	logger.Info("Init ProcessAutomaticMessages service",
		zap.String("journey", "ProcessAutomaticMessages"),
		zap.String("deviceLabel", device.GetDeviceLabel()),
	)

	if device == nil {
		return errors.New("device not found")
	}

	deviceLabel := device.GetDeviceLabel()

	for _, message := range device.GetUnreadMessages() {
		switch message.Topic {
		case "new-connection":
			rs.NewConnection(ctx, deviceLabel, message.Sender)
			message.Read()
			// device.DeleteMessage(id)
		case "new-connection-ack":
			rs.NewConnectionAck(ctx, deviceLabel, message.Sender)
			message.Read()
			// device.DeleteMessage(id)
		case "confirm-connection":
			rs.ConfirmConnection(ctx, deviceLabel, message.Sender)
			message.Read()
			// device.DeleteMessage(id)
		case "update-routing":
			rs.UpdateRouting(ctx, deviceLabel, message.Sender, message.Content.(map[uuid.UUID]entities.Routing))
		default:
			continue
		}
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

	logger.Info("ConfirmConnection",
		zap.String("current", current),
		zap.String("sender", sender),
	)

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

func (rs deviceService) UserMessage(ctx context.Context, current, sender string) error {
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

	timeout := time.After(15 * time.Second)

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
		return nil, fmt.Errorf("device not found: %s", sourceId)
	}

	graph := dijkstra.NewMappedGraph[string]()
	for device := range rs.environment.GetChart() {
		graph.AddEmptyVertex(device)
	}

	for _, routing := range sourceDevice.GetRoutingTable() {
		graph.AddArc(
			routing.Source,
			routing.Target,
			uint64(routing.Weight*1000),
		)
	}

	best, err := graph.Shortest(sourceId, targetId)
	if err != nil {
		return nil, err
	}

	routes := make([]entities.Route, 0)
	for i := 0; i < len(best.Path)-1; i++ {
		route := entities.Route{Source: best.Path[i], Target: best.Path[i+1]}
		routes = append(routes, route)
	}

	printAlloc()

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
	if message.Topic == "user-message" {
		currentDevice.AddMessageToSent(&message)
	}
}

func (rs *deviceService) scheduleTask(deviceLabel string, interval time.Duration, task func(), jobType string) {
	logger.Info("Init scheduleTask service",
		zap.String("deviceLabel", deviceLabel),
		zap.String("interval", fmt.Sprintf("%v", interval)),
	)

	ctx, cancel := context.WithCancel(context.Background())

	switch jobType {
	case "updateRoutingTable":
		rs.jobs.UpdateRoutingTable[deviceLabel] = cancel
	case "message":
		rs.jobs.Message[deviceLabel] = cancel
	case "walk":
		rs.jobs.Walk[deviceLabel] = cancel
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				task()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (rs *deviceService) ScheduleWalk(deviceLabel string) {
	rs.scheduleTask(deviceLabel, 35*time.Second, func() {
		rs.environment.Walk(deviceLabel)
	}, "walk")
}

func (rs *deviceService) ScheduleUpdateRoutingTable(deviceLabel string) {
	rs.scheduleTask(deviceLabel, 20*time.Second, func() {
		rs.UpdateRoutingTable(context.Background(), deviceLabel)
	}, "updateRoutingTable")
}

func (rs *deviceService) ScheduleReadMessages(deviceLabel string) {
	rs.scheduleTask(deviceLabel, 5*time.Second, func() {
		rs.ReadMessages(context.Background(), deviceLabel)
	}, "message")
}

func (rs *deviceService) CancelJob(jobType, deviceLabel string) {
	switch jobType {
	case "updateRoutingTable":
		if cancel, exists := rs.jobs.UpdateRoutingTable[deviceLabel]; exists {
			cancel()
			delete(rs.jobs.UpdateRoutingTable, deviceLabel)
		}
	case "message":
		if cancel, exists := rs.jobs.Message[deviceLabel]; exists {
			cancel()
			delete(rs.jobs.Message, deviceLabel)
		}
	case "walk":
		if cancel, exists := rs.jobs.Walk[deviceLabel]; exists {
			cancel()
			delete(rs.jobs.Walk, deviceLabel)
		}
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

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("-> %d KB\n", m.Alloc/(1024))
}
