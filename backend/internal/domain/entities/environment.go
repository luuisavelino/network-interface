package entities

import (
	"sync"
	"math"
	"math/rand"
	"time"
)

const (
	maxPosX = 50
	maxPosY = 50
)

type Environment struct {
	mu sync.Mutex
	Devices Devices
	Chart   Chart
}

func NewEnvironment() Environment {
	return Environment{
		Devices: make(Devices),
		Chart: make(Chart),
	}
}

func (e *Environment) SetDeviceInChart(deviceLabel string, coverageArea CoverageArea) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Chart[deviceLabel] = &coverageArea
}

func (e *Environment) GetDeviceInChart(deviceLabel string) (*CoverageArea) {
	e.mu.Lock()
	defer e.mu.Unlock()

	chart := e.Chart[deviceLabel]

	return chart
}

func (e *Environment) GetChart() Chart {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Chart
}

func (e *Environment) GetDevices() Devices {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Devices
}

func (e *Environment) Walk(deviceLabel string) {
	rand.Seed(time.Now().UnixNano())

	e.mu.Lock()
	defer e.mu.Unlock()

	device, exists := e.Chart[deviceLabel]
	if !exists {
		return
	}

	device.X += (rand.Intn(3) - 1) * 1
	device.Y += (rand.Intn(3) - 1) * 1

	if device.X > maxPosX {
		device.X = maxPosX
	} else if device.X < 0 {
		device.X = 0
	}

	if device.Y > maxPosY {
		device.Y = maxPosY
	} else if device.Y < 0 {
		device.Y = 0
	}
}

func (e *Environment) GetDistanceTo(fromX, fromY, toX, toY int) float64 {
	return math.Sqrt(math.Pow(float64(fromX-toX), 2) + math.Pow(float64(fromY-toY), 2))
}

func (e *Environment) CheckIfIsInTheCoverageArea(distance, r float64) bool {
	if (distance <= r) {
		return true
	}

	return false
}

func (e *Environment) GetEnvironment() Environment {
	e.mu.Lock()
	defer e.mu.Unlock()

	return *e
}

func (e *Environment) GetDeviceByLabel(deviceLabel string) *Device {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Devices[deviceLabel]
}

func (e *Environment) AddDevice(device *Device) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Devices[device.GetDeviceLabel()] = device
}

func (e *Environment) RemoveDevice(deviceLabel string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.Devices, deviceLabel)
}

func (e *Environment) ScanDeviceNearby(deviceLabel string) []*Device {
	e.mu.Lock()
	defer e.mu.Unlock()

	devicesNearby := make([]*Device, 0)
	sourcePosititon, exists := e.Chart[deviceLabel]
	if !exists {
		return devicesNearby
	}

	for _, device := range e.Devices {
		label := device.GetDeviceLabel()
		if label == deviceLabel {
			continue
		}

		targetPosition, exists := e.Chart[label]
		if !exists {
			continue
		}

		distance := e.GetDistanceTo(sourcePosititon.X, sourcePosititon.Y, targetPosition.X, targetPosition.Y)

		if e.CheckIfIsInTheCoverageArea(distance, sourcePosititon.R) {
			devicesNearby = append(devicesNearby, device)
		}
	}

	return devicesNearby
}

func (e *Environment) CheckIfDeviceIsNearby(deviceLabel, target string) bool {
	devicesNearby := e.ScanDeviceNearby(deviceLabel)
	if len(devicesNearby) == 0 {
		return false
	}

	for _, device := range devicesNearby {
		if device.GetDeviceLabel() == target {
			return true
		}
	}

	return false
}
