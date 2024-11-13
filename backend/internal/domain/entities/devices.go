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
	Power           int
	Status          bool
	Requests        Requests
	DevicesWithConn map[string]Connection
	ScanningDevices bool
	RoutingTable    Routing

	mu sync.Mutex
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
	d.RoutingTable = make(Routing)
	d.mu.Unlock()
}

// TODO: use this
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

func (d *Device) GetDevicesWithConn() map[string]Connection {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.DevicesWithConn
}

func (d *Device) SetDeviceWithConn(device string, errorRate, latency float64) {
	d.mu.Lock()
	d.DevicesWithConn[device] = Connection{
		ErrorRate: errorRate,
		Latency:   latency,
	}
	d.mu.Unlock()
}

func (d *Device) ResetDeviceConn() {
	d.mu.Lock()
	d.DevicesWithConn = make(map[string]Connection)
	d.mu.Unlock()
}

func (d *Device) GetDeviceLabel() string {
	d.mu.Lock()
	label := d.Label
	d.mu.Unlock()
	return label
}

func (d *Device) AddRouting(routingTable Routing) {
	d.mu.Lock()
	for routeType, route := range routingTable {
		if _, ok := d.RoutingTable[routeType]; !ok {
			d.RoutingTable[routeType] = make(map[string]map[string]float64)
			d.RoutingTable[routeType] = route
			continue
		}

		for sourceLabel, target := range route {
			if _, ok := d.RoutingTable[routeType][sourceLabel]; !ok {
				d.RoutingTable[routeType][sourceLabel] = make(map[string]float64)
				d.RoutingTable[routeType][sourceLabel] = target
				continue
			}

			for targetLabel, weight := range target {
				d.RoutingTable[routeType][sourceLabel][targetLabel] = weight
			}
		}
	}
	d.mu.Unlock()
}

func (d *Device) RemoveRoutings(routeType, source, target string) {
	d.mu.Lock()
	route, exists := d.RoutingTable[routeType][source]
	if exists {
		delete(route, target)
	}
	d.mu.Unlock()
}

func (d *Device) RemoveFromTableRoutesWith(deviceLabel string) {
	d.mu.Lock()
	for _, routes := range d.RoutingTable {
		for sourceLabel, _ := range routes {
			if sourceLabel == deviceLabel {
				delete(d.RoutingTable, sourceLabel)
				continue
			}

			route, exists := d.RoutingTable[sourceLabel][deviceLabel]
			if exists {
				delete(route, deviceLabel)
			}
		}
	}
	d.mu.Unlock()
}

func (d *Device) GetRoutingTable() Routing {
	d.mu.Lock()
	routingTable := make(Routing)
	routingTable = d.RoutingTable
	d.mu.Unlock()
	return routingTable
}

func (d *Device) GetUnreadRequests() map[uuid.UUID]*Request {
	d.mu.Lock()
	unreadRequests := make(map[uuid.UUID]*Request)
	for _, Request := range d.Requests.Received {
		if !Request.IsRead() {
			unreadRequests[Request.ID] = Request
		}
	}
	d.mu.Unlock()
	return unreadRequests
}

func (d *Device) GetReadRequests() map[uuid.UUID]*Request {
	d.mu.Lock()
	readRequests := make(map[uuid.UUID]*Request)
	for _, Request := range d.Requests.Received {
		if Request.IsRead() {
			readRequests[Request.ID] = Request
		}
	}
	d.mu.Unlock()
	return readRequests
}

func (d *Device) GetRequestsSent() map[uuid.UUID]*Request {
	d.mu.Lock()
	RequestsSent := d.Requests.Sent
	d.mu.Unlock()
	return RequestsSent
}

func (d *Device) AddRequestToSent(Request *Request) {
	d.mu.Lock()
	d.Requests.Sent[Request.ID] = Request
	d.mu.Unlock()
}

func (d *Device) AddRequestToReceived(Request *Request) {
	d.mu.Lock()
	d.Requests.Received[Request.ID] = Request
	d.mu.Unlock()
}

func (d *Device) PrintPrettyTable() {
	table := d.GetRoutingTable()
	if len(table) == 0 {
		fmt.Println("No table found")
		return
	}

	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("| %-10s | %-24s | \n", "Device", d.GetDeviceLabel())
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("| %-10s | %-6s | %-6s | %-6s | \n", "Type", "Source", "Target", "Weight")
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	for routingType, routes := range table {
		for sourceLabel, target := range routes {
			for targetLabel, weight := range target  {
				fmt.Printf("| %-10s | %-6s | %-6s | %-6.2f |\n", routingType, sourceLabel, targetLabel, weight)
			}
		}
	}
	fmt.Printf("%s\n", strings.Repeat("-", 40))
}

func (d *Device) DeleteRequest(RequestId uuid.UUID) {
	d.mu.Lock()
	delete(d.Requests.Received, RequestId)
	d.mu.Unlock()
}
