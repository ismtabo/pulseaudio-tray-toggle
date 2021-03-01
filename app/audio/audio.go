package audio

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/ismtabo/pulseaudio"
	"github.com/ismtabo/pulseaudio-tray-toggle/pkg/config"
)

// Audio is a adapt type over Pulseaudio module
type Audio struct {
	client *pulseaudio.Client
	Source *Device
	Sink   *Device
}

// NewAudio creates a new instance of audio manager
func NewAudio(cnf *config.Config, client *pulseaudio.Client) (*Audio, error) {
	sink, err := getSink(client, cnf.Pulseaudio.Sink)
	if err != nil {
		return nil, err
	}

	source, err := getSource(client, cnf.Pulseaudio.Source)
	if err != nil {
		return nil, err
	}

	return &Audio{client, source, sink}, nil
}

func getSink(client *pulseaudio.Client, name string) (*Device, error) {
	sinks, err := client.Core().ListPath("Sinks")
	if err != nil {
		return nil, err
	}
	return foundDevice(client, sinks, name)
}

func getSource(client *pulseaudio.Client, name string) (*Device, error) {
	sources, err := client.Core().ListPath("Sources")
	if err != nil {
		return nil, err
	}
	return foundDevice(client, sources, name)
}

func foundDevice(client *pulseaudio.Client, devices []dbus.ObjectPath, targetName string) (*Device, error) {
	for _, path := range devices {
		device := client.Device(path)

		name, err := device.String("Name")
		if err != nil {
			return nil, err
		}

		if name == targetName {
			return NewDevice(client, path)
		}
	}

	return nil, fmt.Errorf("Device not found with name: %s", targetName)
}
