package entities

import (
	"github.com/google/uuid"
	"sync"
	"fmt"
	"strings"
)

type Devices map[string]*Device

type Device struct {
	Label 				string
	mu 						sync.Mutex
	Power					int
	Messages			Messages
	WalkingSpeed	int
	MessageFreq		int
	RoutingTable	map[uuid.UUID]Routing
}

func (d *Device) GetDeviceLabel() string {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.Label
}

func (d *Device) AddRouting(routingTable map[uuid.UUID]Routing) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for key, value := range routingTable {
		d.RoutingTable[key] = value
	}
}

func (d *Device) RemoveRoutings(routes []uuid.UUID) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, route := range routes {
		delete(d.RoutingTable, route)
	}	
}

func (d *Device) RemoveFromTableRoutesWith(deviceLabel string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for routeUuid, route := range d.RoutingTable {
		if deviceLabel == route.Source || deviceLabel == route.Target {
			delete(d.RoutingTable, routeUuid)
		}
	}
}

func (d *Device) GetRoutingTable() map[uuid.UUID]Routing {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.RoutingTable
}

func (d *Device) GetUnreadMessages() []*Message {
	d.mu.Lock()
	defer d.mu.Unlock()

	var unreadMessages []*Message
	for _, message := range d.Messages.Received {
		if !message.IsRead() {
			unreadMessages = append(unreadMessages, message)
		}
	}

	return unreadMessages
}

func (d *Device) GetReadMessages() []*Message {
	d.mu.Lock()
	defer d.mu.Unlock()

	var readMessages []*Message
	for _, message := range d.Messages.Received {
		if message.IsRead() {
			readMessages = append(readMessages, message)
		}
	}

	return readMessages
}

func (d *Device) GetMessagesSent() []*Message {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.Messages.Sent
}

func (d *Device) AddMessageToSent(message *Message) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Messages.Sent = append(d.Messages.Sent, message)
}

func (d *Device) AddMessageToReceived(message *Message) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Messages.Received = append(d.Messages.Received, message)
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
		fmt.Printf("| %-36v | %-6s | %-6s | %-6.2v |\n", routeUuid, route.Source, route.Target, route.Weight)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 67))
}
