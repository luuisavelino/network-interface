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
	e.Chart[deviceLabel] = &coverageArea
	e.mu.Unlock()
}

func (e *Environment) GetDeviceInChart(deviceLabel string) *CoverageArea {
	e.mu.Lock()
	chart := e.Chart[deviceLabel]
	e.mu.Unlock()
	return chart
}

func (e *Environment) GetChart() Chart {
	e.mu.Lock()
	chart := e.Chart
	e.mu.Unlock()
	return chart
}

func (e *Environment) GetDevices() Devices {
	e.mu.Lock()
	devices := e.Devices
	e.mu.Unlock()
	return devices
}

func (e *Environment) Walk(deviceLabel string) {
	rand.Seed(time.Now().UnixNano())

	e.mu.Lock()
	device, exists := e.Chart[deviceLabel]
	if !exists {
		e.mu.Unlock()
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
	e.mu.Unlock()
}

func (e *Environment) GetDistanceTo(fromX, fromY, toX, toY int) float64 {
	return math.Sqrt(math.Pow(float64(fromX-toX), 2) + math.Pow(float64(fromY-toY), 2))
}

func (e *Environment) CheckIfIsInTheCoverageArea(distance, r float64) bool {
	return distance <= r
}

func (e *Environment) GetEnvironment() Environment {
	e.mu.Lock()
	env := *e
	e.mu.Unlock()
	return env
}

func (e *Environment) GetDeviceByLabel(deviceLabel string) *Device {
	e.mu.Lock()
	device := e.Devices[deviceLabel]
	e.mu.Unlock()
	return device
}

func (e *Environment) AddDevice(device *Device) {
	e.mu.Lock()
	e.Devices[device.GetDeviceLabel()] = device
	e.mu.Unlock()
}

func (e *Environment) RemoveDevice(deviceLabel string) {
	e.mu.Lock()
	delete(e.Devices, deviceLabel)
	delete(e.Chart, deviceLabel)
	e.mu.Unlock()
}

func (e *Environment) ScanDeviceNearby(deviceLabel string) []*Device {
	e.mu.Lock()
	sourcePosititon, exists := e.Chart[deviceLabel]
	if !exists {
		e.mu.Unlock()
		return nil
	}

	devicesNearby := make([]*Device, 0)
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
	e.mu.Unlock()
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
