package entities

import (
	"github.com/google/uuid"
	"math/rand"
	"sync"
	"time"
	"math"
	"fmt"
	"strings"
)

const (
	maxPosX = 50
	maxPosY = 50
)


type Device struct {
	ID 			int
	mu 			sync.RWMutex
	Power		int
	PosX		int
	PosY		int
	WalkingSpeed	int
	MessageFreq		int
	RoutingTable	map[uuid.UUID]Routing
}

func (d *Device) Walk() {
	rand.Seed(time.Now().UnixNano())

	d.mu.Lock()
	defer d.mu.Unlock()

	d.PosX += (rand.Intn(3) - 1) * d.WalkingSpeed
	d.PosY += (rand.Intn(3) - 1) * d.WalkingSpeed

	if d.PosX > maxPosX {
		d.PosX = maxPosX
	} else if d.PosX < 0 {
		d.PosX = 0
	}

	if d.PosY > maxPosY {
		d.PosY = maxPosY
	} else if d.PosY < 0 {
		d.PosY = 0
	}
}

func (d *Device) CheckIfIsInTheCoverageArea(posX, posY int) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if (math.Sqrt(math.Pow(float64(d.PosX-posX), 2) + math.Pow(float64(d.PosY-posY), 2)) <= float64(d.Power)) {
		return true
	}

	return false
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

func (d *Device) GetDistanceTo(posX, posY int) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	return math.Sqrt(math.Pow(float64(d.PosX-posX), 2) + math.Pow(float64(d.PosY-posY), 2))
}

func (d *Device) RemoveOwnRoutesFromTable() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for routeUuid, route := range d.RoutingTable {
		if d.ID == route.Source || d.ID == route.Target {
			delete(d.RoutingTable, routeUuid)
		}
	}
}

func (d *Device) UpdateRoutingTable(routes map[uuid.UUID]Routing) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.RoutingTable = routes
}

func (d *Device) GetRoutingTable() map[uuid.UUID]Routing {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.RoutingTable
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
		fmt.Printf("| %-36v | %-6d | %-6d | %-6.2v |\n", routeUuid, route.Source, route.Target, route.Weight)
	}
	fmt.Printf("%s\n", strings.Repeat("-", 67))
}
