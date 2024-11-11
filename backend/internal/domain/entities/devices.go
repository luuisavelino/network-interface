package entities

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Devices map[string]*Device

type Device struct {
	Label           string
	mu              sync.Mutex
	Power           int
	Status          bool
	Messages        Messages
	WalkingSpeed    int
	MessageFreq     int
	DevicesWithConn []string
	ScanningDevices bool
	RoutingTable    map[uuid.UUID]Routing
}

func (d *Device) GetStatus() bool {
	d.mu.Lock()
	status := d.Status
	d.mu.Unlock()
	return status
}

func (d *Device) SetStatus(status bool) {
	d.mu.Lock()
	d.Status = status
	d.mu.Unlock()
}

func (d *Device) ResetRoutingTable() {
	d.mu.Lock()
	d.RoutingTable = make(map[uuid.UUID]Routing)
	d.mu.Unlock()
}

func (d *Device) IsScanningDevices() bool {
	d.mu.Lock()
	scanning := d.ScanningDevices
	d.mu.Unlock()
	return scanning
}

func (d *Device) SetScanningDevices(scanning bool) {
	d.mu.Lock()
	d.ScanningDevices = scanning
	d.mu.Unlock()
}

func (d *Device) GetDevicesWithConn() []string {
	d.mu.Lock()
	devicesWithConn := make([]string, len(d.DevicesWithConn))
	copy(devicesWithConn, d.DevicesWithConn)
	d.mu.Unlock()
	return devicesWithConn
}

func (d *Device) SetDeviceWithConn(device string) {
	d.mu.Lock()
	d.DevicesWithConn = append(d.DevicesWithConn, device)
	d.mu.Unlock()
}

func (d *Device) ResetDeviceConn() {
	d.mu.Lock()
	d.DevicesWithConn = []string{}
	d.mu.Unlock()
}

func (d *Device) GetDeviceLabel() string {
	d.mu.Lock()
	label := d.Label
	d.mu.Unlock()
	return label
}

func (d *Device) AddRouting(routingTable map[uuid.UUID]Routing) {
	d.mu.Lock()
	for key, value := range routingTable {
		d.RoutingTable[key] = value
	}
	d.mu.Unlock()
}

func (d *Device) RemoveRoutings(routes []uuid.UUID) {
	d.mu.Lock()
	for _, route := range routes {
		delete(d.RoutingTable, route)
	}
	d.mu.Unlock()
}

func (d *Device) RemoveFromTableRoutesWith(deviceLabel string) {
	d.mu.Lock()
	for routeUuid, route := range d.RoutingTable {
		if deviceLabel == route.Source || deviceLabel == route.Target {
			delete(d.RoutingTable, routeUuid)
		}
	}
	d.mu.Unlock()
}

func (d *Device) GetRoutingTable() map[uuid.UUID]Routing {
	d.mu.Lock()
	routingTable := make(map[uuid.UUID]Routing, len(d.RoutingTable))
	for key, value := range d.RoutingTable {
		routingTable[key] = value
	}
	d.mu.Unlock()
	return routingTable
}

func (d *Device) GetUnreadMessages() map[uuid.UUID]*Message {
	d.mu.Lock()
	unreadMessages := make(map[uuid.UUID]*Message)
	for _, message := range d.Messages.Received {
		if !message.IsRead() {
			unreadMessages[message.ID] = message
		}
	}
	d.mu.Unlock()
	return unreadMessages
}

func (d *Device) GetReadMessages() map[uuid.UUID]*Message {
	d.mu.Lock()
	readMessages := make(map[uuid.UUID]*Message)
	for _, message := range d.Messages.Received {
		if message.IsRead() {
			readMessages[message.ID] = message
		}
	}
	d.mu.Unlock()
	return readMessages
}

func (d *Device) GetMessagesSent() map[uuid.UUID]*Message {
	d.mu.Lock()
	messagesSent := d.Messages.Sent
	d.mu.Unlock()
	return messagesSent
}

func (d *Device) AddMessageToSent(message *Message) {
	d.mu.Lock()
	d.Messages.Sent[message.ID] = message
	d.mu.Unlock()
}

func (d *Device) AddMessageToReceived(message *Message) {
	d.mu.Lock()
	d.Messages.Received[message.ID] = message
	d.mu.Unlock()
}

func (d *Device) PrintPrettyTable() {
	table := d.GetRoutingTable()
	if len(table) == 0 {
		fmt.Println("No table found")
		return
	}

	fmt.Printf("%s\n", strings.Repeat("-", 67))
	fmt.Printf("| %-36s | %-6s | %-6s | %-6s | \n", "Route UUID", "Source", "Target", "Weight")
	fmt.Printf("%s\n", strings.Repeat("-", 67))
	for routeUuid, route := range table {
		fmt.Printf("| %-36v | %-6s | %-6s | %-6.2f |\n", routeUuid, route.Source, route.Target, route.Weight)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 67))
}

func (d *Device) DeleteMessage(messageId uuid.UUID) {
	d.mu.Lock()
	delete(d.Messages.Received, messageId)
	d.mu.Unlock()
}
