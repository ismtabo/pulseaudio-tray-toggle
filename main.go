package main

import (
	"log"
	"path"
	"path/filepath"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/ismtabo/pulseaudio-tray-toggle/app"
)

//
//--------------------------------------------------------------------[ MAIN ]--

// Create a pulse dbus service with 2 clients, listen to events,
// then use some properties.
//
func main() {
	_, fileName, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		log.Fatal(err)
	}
	config := path.Join(dir, "./config.yml")
	application, err := app.NewApplication(&config)
	if err != nil {
		log.Fatal(err)
	}
	systray.Run(application.Start, application.Stop)
}
