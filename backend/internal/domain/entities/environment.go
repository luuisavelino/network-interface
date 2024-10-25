package entities

import (
	"sync"
)

type Environment struct {
	mu sync.Mutex
	Devices map[int]*Device
}

func NewEnvironment() Environment {
	return Environment{
		Devices: make(map[int]*Device),
	}
}

func (e *Environment) GetDevices() map[int]*Device {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Devices
}

func (e *Environment) GetEnvironment() Environment {
	e.mu.Lock()
	defer e.mu.Unlock()

	return *e
}

func (e *Environment) GetDeviceById(deviceId int) *Device {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Devices[deviceId]
}

func (e *Environment) AddDevice(device *Device) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.Devices[device.GetDeviceID()] = device
}

func (e *Environment) RemoveDevice(deviceId int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.Devices, deviceId)
}

func (e *Environment) ScanDevicesWithCommunication(deviceId int) []*Device {
	e.mu.Lock()
	defer e.mu.Unlock()

	currentDevice := e.Devices[deviceId]

	devicesWithCommunication := make([]*Device, 0)

	for _, device := range e.Devices {
		if device.GetDeviceID() == currentDevice.GetDeviceID() {
			continue
		}

		if currentDevice.CheckIfIsInTheCoverageArea(device.PosX, device.PosY) && 
				device.CheckIfIsInTheCoverageArea(currentDevice.PosX, currentDevice.PosY) {
			devicesWithCommunication = append(devicesWithCommunication, device)
		}
	}

	return devicesWithCommunication
}
