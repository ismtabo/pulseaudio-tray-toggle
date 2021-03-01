package main

import (
	"log"
	"os"
	"path"

	"github.com/getlantern/systray"
	"github.com/ismtabo/pulseaudio-tray-toggle/app"
)

//
//--------------------------------------------------------------------[ MAIN ]--

// Create a pulse dbus service with 2 clients, listen to events,
// then use some properties.
//
func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	config := path.Join(dir, "config.yml")
	application, err := app.NewApplication(&config)
	if err != nil {
		log.Fatal(err)
	}
	systray.Run(application.Start, application.Stop)
}
