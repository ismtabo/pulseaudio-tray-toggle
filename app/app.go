package app

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/ismtabo/pulseaudio"
	"github.com/ismtabo/pulseaudio-tray-toggle/app/audio"
	"github.com/ismtabo/pulseaudio-tray-toggle/pkg/config"
)

// Application implements Service interface
type Application struct {
	pulse        *pulseaudio.Client
	loadedModule bool
	audio        *audio.Audio
	config       *config.Config
}

// NewApplication create a new Application instance for the given config
func NewApplication(configPath *string) (*Application, error) {

	config, err := config.Load(configPath)

	if err != nil {
		return nil, err
	}

	// Load pulseaudio DBus module if needed. This module is mandatory, but it
	// can also be configured in system files. See package doc.
	isLoaded, err := pulseaudio.ModuleIsLoaded()
	if err != nil {
		return nil, err
	}

	if !isLoaded {
		err = pulseaudio.LoadModule()
		if err != nil {
			return nil, err
		}
	}

	// Connect to the pulseaudio dbus service.
	client, err := pulseaudio.New()
	if err != nil {
		return nil, err
	}
	// defer pulse.Close()

	audio, err := audio.NewAudio(config, client)
	if err != nil {
		return nil, err
	}

	application := Application{client, isLoaded, audio, config}

	application.iniTrayIcon()

	return &application, nil
}

// iniTrayIcon function
func (application Application) iniTrayIcon() {
	file, err := os.Open("static/img/headphone-icon.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	systray.SetIcon(bytes)
	systray.SetTitle("Pulseaudio Toggle")
	systray.SetTooltip("Change port with a click")
	systray.AddSeparator()
	initDeviceMenu(application.audio.Sink)
	systray.AddSeparator()
	initDeviceMenu(application.audio.Source)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func initDeviceMenu(device *audio.Device) {
	menuItem := systray.AddMenuItem(device.Description, device.Description)
	menuSubitems := make([]*systray.MenuItem, len(device.Ports))
	for i, port := range device.Ports {
		menuSubitems[i] = menuItem.AddSubMenuItem(port.Description, port.Description)
	}

	for i, subitem := range menuSubitems {
		go func(subitem *systray.MenuItem, port *audio.DevicePort) {
			for range subitem.ClickedCh {
				changeDevicePort(menuSubitems, device, port)
			}
		}(subitem, device.Ports[i])
	}

	updateMenuSubitems(menuSubitems, device)
}

func updateMenuSubitems(menuItems []*systray.MenuItem, device *audio.Device) {
	for i, item := range menuItems {
		if device.Ports[i].Description == device.ActivePort.Description {
			item.Disable()
			item.Check()
		} else {
			item.Enable()
			item.Uncheck()
		}
	}
}

func changeDevicePort(menuItems []*systray.MenuItem, device *audio.Device, port *audio.DevicePort) {
	err := device.SetActivePort(port)
	if err != nil {
		log.Printf("Error - changing device port: %s", err)
		return
	}
	updateMenuSubitems(menuItems, device)
}

// Start initialize application
func (application *Application) Start() {}

// Stop finalize application
func (application Application) Stop() {
	// clean up here
	application.pulse.Close()

	if application.loadedModule {
		pulseaudio.UnloadModule()
	}
}
