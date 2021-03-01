package audio

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/ismtabo/pulseaudio"
)

// DevicePort
type DevicePort struct {
	object      *pulseaudio.Object
	Description string
}

// Device
type Device struct {
	client      *pulseaudio.Client
	object      *pulseaudio.Object
	Description string
	Ports       []*DevicePort
	ActivePort  *DevicePort
}

// NewDevice
func NewDevice(client *pulseaudio.Client, objectPath dbus.ObjectPath) (*Device, error) {
	device := client.Device(objectPath)

	properties, err := device.MapString("PropertyList")
	if err != nil {
		return nil, err
	}

	description, found := properties["device.description"]
	if !found {
		return nil, fmt.Errorf("Description not found for device: %s", objectPath)
	}

	ports, err := getDevicePorts(client, device)
	if err != nil {
		return nil, err
	}

	activePort, err := getDeviceActivePort(client, device)
	if err != nil {
		return nil, err
	}

	return &Device{client, device, description, ports, activePort}, nil
}

func getDeviceActivePort(client *pulseaudio.Client, device *pulseaudio.Object) (*DevicePort, error) {
	portPath, err := device.ObjectPath("ActivePort")
	if err != nil {
		return nil, err
	}

	port := client.DevicePort(portPath)

	description, err := port.String("Description")
	if err != nil {
		return nil, err
	}

	return &DevicePort{port, description}, nil
}

func getDevicePorts(client *pulseaudio.Client, device *pulseaudio.Object) ([]*DevicePort, error) {
	portsPaths, err := device.ListPath("Ports")
	if err != nil {
		return nil, err
	}

	ports := make([]*DevicePort, len(portsPaths))

	for i, portPath := range portsPaths {
		port := client.DevicePort(portPath)

		description, err := port.String("Description")
		if err != nil {
			return nil, err
		}

		ports[i] = &DevicePort{port, description}
	}

	return ports, nil
}

// SetActivePort change the active port of the device with the one given
func (device *Device) SetActivePort(newActivePort *DevicePort) error {
	error := device.object.Set("ActivePort", newActivePort.object.Path())
	if error != nil {
		return error
	}

	device.ActivePort = newActivePort
	return nil
}
