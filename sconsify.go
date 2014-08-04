package main

import (
	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
)

func main() {
	events := events.InitialiseEvents()

	go spotify.Initialise(&events)

	if <-events.Initialised {
		ui.Start(&events)
	}
}
