package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/luuisavelino/network-interface/internal/domain/entities"
	"github.com/luuisavelino/network-interface/internal/domain/entities/dijkstra"
	"github.com/luuisavelino/network-interface/pkg/logger"
	"go.uber.org/zap"
)

func NewDeviceService(environment *entities.Environment, scheduler gocron.Scheduler) DeviceService {
	return deviceService{
		environment: environment,
		scheduler:   scheduler,
		jobs: Jobs{
			UpdateRoutingTable: make(map[string]context.CancelFunc),
			Request:            make(map[string]context.CancelFunc),
			Walk:               make(map[string]context.CancelFunc),
		},
	}
}

type deviceService struct {
	environment *entities.Environment
	scheduler   gocron.Scheduler
	jobs        Jobs
}

type Jobs struct {
	UpdateRoutingTable map[string]context.CancelFunc
	Request            map[string]context.CancelFunc
	Walk               map[string]context.CancelFunc
}

type DeviceService interface {
	GetDevices(ctx context.Context) (entities.Devices, error)
	InsertDevice(ctx context.Context, device entities.Device) (entities.Device, error)
	UpdateRoutingTable(ctx context.Context, deviceLabel string) error
	GetDevice(ctx context.Context, deviceLabel string) (entities.Device, error)
	GetRoute(tx context.Context, sourceId string, targetId string, routingType string) ([]entities.Route, error)
	DeleteDevice(ctx context.Context, deviceLabel string) error
	SendUserMessage(ctx context.Context, request entities.Request) error
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

	device.RoutingTable = make(entities.Routing, 0)
	device.Requests = entities.Requests{
		Sent:     make(map[uuid.UUID]*entities.Request),
		Received: make(map[uuid.UUID]*entities.Request),
	}

	device.SetStatus(true)

	rs.environment.AddDevice(&device)

	label := device.GetDeviceLabel()

	rs.ScheduleReadRequests(label)
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
		return entities.Device{}, fmt.Errorf("device not found: %s", deviceLabel)
	}

	return *device, nil
}

func (rs deviceService) SendUserMessage(ctx context.Context, request entities.Request) error {
	logger.Info("Init SendUserMessage service",
		zap.String("journey", "SendUserMessage"),
		zap.String("current", request.Header.Sender),
		zap.String("target", request.Header.Destination),
		zap.String("request", request.Body.(string)),
	)

	currentDevice := rs.environment.GetDeviceByLabel(request.Header.Sender)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", request.Header.Sender)
	}

	targetDevice := rs.environment.GetDeviceByLabel(request.Header.Destination)
	if targetDevice == nil {
		return fmt.Errorf("device not found: %s", request.Header.Destination)
	}

	var routingType string
	switch request.Header.ContentType {
	case "text":
		routingType = "distance"
	case "audio":
		routingType = "latency"
	case "file":
		routingType = "error-rate"
	default:
		routingType = "distance"
	}

	routes, err := rs.GetRoute(ctx, request.Header.Sender, request.Header.Destination, routingType)
	if err != nil {
		return err
	}

	if len(routes) == 0 {
		return fmt.Errorf("no route to send message")
	}

	rs.PropagateRequest(ctx, routes, request)

	return nil
}

func (rs deviceService) PropagateRequest(ctx context.Context, routes []entities.Route, request entities.Request) {
	sender := routes[0].Source
	target := routes[0].Target

	currentDevice := rs.environment.GetDeviceByLabel(sender)
	targetDevice := rs.environment.GetDeviceByLabel(target)

	routes = routes[1:]

	newRequest := entities.NewRequest(
		"user-message",
		sender,
		target,
		routes,
		request.Body,
	)

	rs.SendRequest(currentDevice, targetDevice, newRequest)

	// go func() {
	// 	const maxRetries = 1
	// 	const maxTimeout = 5 * time.Second

	// 	for i := 0; i < maxRetries; i++ {
	// 		rs.SendRequest(currentDevice, targetDevice, newRequest)

	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case <-time.After(maxTimeout):
	// 			for _, msg := range targetDevice.GetUnreadRequests() {
	// 				if msg.Header.Topic == "user-message-ack" && msg.Header.Sender == target {
	// 					msg.Read()
	// 					return
	// 				}
	// 			}
	// 		}
	// 	}
	// 	logger.Info("Failed to receive user-message-ack after retries",
	// 		zap.String("current", request.Header.Sender),
	// 		zap.String("target", request.Header.Destination),
	// 	)
	// }()
}

func (rs deviceService) ReadRequests(ctx context.Context, deviceLabel string) error {
	logger.Info("Init ReadRequests service",
		zap.String("journey", "ReadRequests"),
		zap.String("deviceLabel", deviceLabel),
	)

	device := rs.environment.GetDeviceByLabel(deviceLabel)
	if device == nil {
		fmt.Println("device not found: ", deviceLabel)
		return fmt.Errorf("device not found: %s", deviceLabel)
	}

	rs.ProcessAutomaticRequests(ctx, device)

	return nil
}

func (rs deviceService) ProcessAutomaticRequests(ctx context.Context, device *entities.Device) error {
	logger.Info("Init ProcessAutomaticRequests service",
		zap.String("journey", "ProcessAutomaticRequests"),
		zap.String("deviceLabel", device.GetDeviceLabel()),
	)

	if device == nil {
		return errors.New("device not found")
	}

	deviceLabel := device.GetDeviceLabel()

	for id, request := range device.GetUnreadRequests() {
		switch request.Header.Topic {
		case "new-connection":
			rs.NewConnection(ctx, deviceLabel, request.Header.Sender)
			request.Read()
			device.DeleteRequest(id)

		case "new-connection-ack":
			rs.NewConnectionAck(ctx, deviceLabel, request.Header.Sender)
			request.Read()
			device.DeleteRequest(id)
		case "confirm-connection":
			rs.ConfirmConnection(ctx, deviceLabel, request.Header.Sender)
			request.Read()
			device.DeleteRequest(id)
		case "update-routing":
			rs.UpdateRouting(ctx, deviceLabel, request.Header.Sender, request.Body.(entities.Routing))
			device.DeleteRequest(id)
		case "user-message":
			rs.UserMessage(ctx, deviceLabel, request.Header.Sender, *request)
			request.Read()
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
		return fmt.Errorf("device not found: %s", current)
	}

	isNearby := rs.environment.CheckIfDeviceIsNearby(current, sender)
	if !isNearby {
		fmt.Println("device not nearby: ", sender)
		return nil
	}

	senderDevice := rs.environment.GetDeviceByLabel(sender)

	request := entities.NewRequest(
		"new-connection-ack",
		current,
		sender,
		nil,
		nil,
	)

	rs.SendRequest(currentDevice, senderDevice, request)

	return nil
}

func (rs deviceService) NewConnectionAck(ctx context.Context, current, sender string) error {
	logger.Info("Init NewConnectionAck service",
		zap.String("journey", "NewConnectionAck"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", current)
	}

	senderDevice := rs.environment.GetDeviceByLabel(sender)
	if senderDevice == nil {
		return fmt.Errorf("device not found: %s", sender)
	}

	request := entities.NewRequest(
		"confirm-connection",
		current,
		sender,
		nil,
		nil,
	)

	rs.SendRequest(currentDevice, senderDevice, request)

	currentDevice.SetDeviceWithConn(
		sender,
		rand.Float64() * 10,
		rand.Float64() * 10,
	)

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

func (rs deviceService) UpdateRouting(ctx context.Context, current, sender string, routingTable entities.Routing) error {
	logger.Info("Init UpdateRouting service",
		zap.String("journey", "UpdateRouting"),
		zap.String("current", current),
	)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", current)
	}

	currentDevice.RemoveFromTableRoutesWith(sender)
	currentDevice.AddRouting(routingTable)

	currentDevice.PrintPrettyTable()

	return nil
}

func (rs deviceService) UserMessage(ctx context.Context, current, sender string, request entities.Request) error {
	logger.Info("Init UserMessage service",
		zap.String("journey", "UserMessage"),
		zap.String("current", current),
	)

	fmt.Println("UserMessage service:", current)

	currentDevice := rs.environment.GetDeviceByLabel(current)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", current)
	}

	senderDevice := rs.environment.GetDeviceByLabel(request.Header.Sender)
	if senderDevice == nil {
		return fmt.Errorf("device not found: %s", sender)
	}

	userMessageAck := entities.NewRequest(
		"user-message-ack",
		current,
		sender,
		nil,
		nil,
	)

	rs.SendRequest(currentDevice, senderDevice, userMessageAck)

	if len(request.Header.Path) == 0 {
		return nil
	}

	rs.PropagateRequest(ctx, request.Header.Path, request)

	return nil
}

func (rs deviceService) UpdateRoutingTable(ctx context.Context, deviceLabel string) error {
	logger.Info("Init UpdateRoutingTable service",
		zap.String("journey", "UpdateRoutingTable"),
		zap.String("deviceLabel", deviceLabel),
	)

	currentDevice := rs.environment.GetDeviceByLabel(deviceLabel)
	if currentDevice == nil {
		return fmt.Errorf("device not found: %s", deviceLabel)
	}

	currentDevice.SetScanningDevices(true)
	currentDevice.ResetDeviceConn()

	devicePosition := rs.environment.GetDeviceInChart(deviceLabel)
	if devicePosition == nil {
		return fmt.Errorf("device not found: %s", deviceLabel)
	}

	devicesNearby := rs.environment.ScanDeviceNearby(deviceLabel)
	if len(devicesNearby) == 0 {
		currentDevice.RemoveFromTableRoutesWith(deviceLabel)
		return nil
	}

	for _, device := range devicesNearby {
		request := entities.NewRequest(
			"new-connection",
			deviceLabel,
			device.GetDeviceLabel(),
			nil,
			nil,
		)

		rs.SendRequest(currentDevice, device, request)
	}

	time.Sleep(15 * time.Second)

	rs.PropagateRoutingTable(ctx, currentDevice)

	currentDevice.SetScanningDevices(false)

	return nil
}

func (rs deviceService) BuildRoutingTableRow(device *entities.Device, deviceConn string) entities.Routing {
	currDevice := device.GetDeviceLabel()

	weight := rs.environment.GetDistanceTo(
		rs.environment.GetDeviceInChart(currDevice).X,
		rs.environment.GetDeviceInChart(currDevice).Y,
		rs.environment.GetDeviceInChart(deviceConn).X,
		rs.environment.GetDeviceInChart(deviceConn).Y,
	)

	routingTable := make(entities.Routing, 0)
	
	routingTable["distance"] = make(map[string]map[string]float64)
	routingTable["distance"][currDevice] = make(map[string]float64)
	routingTable["distance"][currDevice][deviceConn] = weight
	
	routingTable["error-rate"] = make(map[string]map[string]float64)
	routingTable["error-rate"][currDevice] = make(map[string]float64)
	routingTable["error-rate"][currDevice][deviceConn] = device.DevicesWithConn[deviceConn].GetErrorRate()
	
	routingTable["latency"] = make(map[string]map[string]float64)
	routingTable["latency"][currDevice] = make(map[string]float64)
	routingTable["latency"][currDevice][deviceConn] = device.DevicesWithConn[deviceConn].GetLatency()

	return routingTable
}

func (rs deviceService) PropagateRoutingTable(ctx context.Context, device *entities.Device) {
	currDevice := device.GetDeviceLabel()
	device.RemoveFromTableRoutesWith(currDevice)

	for deviceConn := range device.GetDevicesWithConn() {
		connDevice := rs.environment.GetDeviceByLabel(deviceConn)
		if connDevice == nil {
			continue
		}

		routingTable := rs.BuildRoutingTableRow(device, deviceConn)
		device.AddRouting(routingTable)

		request := entities.NewRequest(
			"update-routing",
			currDevice,
			deviceConn,
			nil,
			device.GetRoutingTable(),
		)

		rs.SendRequest(device, connDevice, request)
	}
}

func (rs deviceService) GetRoute(tx context.Context, sourceId, targetId, routingType string) ([]entities.Route, error) {
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

	for sourceLabel, target := range sourceDevice.GetRoutingTable()[routingType] {
		for targetLabel, weight := range target {
			graph.AddArc(
				sourceLabel,
				targetLabel,
				uint64(weight*1000),
			)
		}
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

func (rs deviceService) SendRequest(currentDevice *entities.Device, targetDevice *entities.Device, request entities.Request) {
	logger.Info("Init SendRequest service",
		zap.String("journey", "SendRequest"),
		zap.String("source", currentDevice.GetDeviceLabel()),
		zap.String("target", targetDevice.GetDeviceLabel()),
		zap.String("request", request.Header.Topic),
	)

	targetDevice.AddRequestToReceived(&request)
	if request.Header.Topic == "user-message" || request.Header.Topic == "user-message-ack" {
		currentDevice.AddRequestToSent(&request)
	}
	// currentDevice.AddRequestToSent(&request)
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
	case "request":
		rs.jobs.Request[deviceLabel] = cancel
	case "walk":
		rs.jobs.Walk[deviceLabel] = cancel
	default:
		cancel()
		return
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

func (rs *deviceService) ScheduleReadRequests(deviceLabel string) {
	rs.scheduleTask(deviceLabel, 5*time.Second, func() {
		rs.ReadRequests(context.Background(), deviceLabel)
	}, "request")
}

func (rs *deviceService) CancelJob(jobType, deviceLabel string) {
	switch jobType {
	case "updateRoutingTable":
		if cancel, exists := rs.jobs.UpdateRoutingTable[deviceLabel]; exists {
			cancel()
			delete(rs.jobs.UpdateRoutingTable, deviceLabel)
		}
	case "request":
		if cancel, exists := rs.jobs.Request[deviceLabel]; exists {
			cancel()
			delete(rs.jobs.Request, deviceLabel)
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
		return fmt.Errorf("device not found: %s", deviceLabel)
	}

	rs.CancelJob("request", deviceLabel)
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
