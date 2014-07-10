package main

import (
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	sp "github.com/op/go-libspotify/spotify"
)

func main() {
	initialised := make(chan bool, 1)
	status := make(chan string)
	toPlay := make(chan sp.Track)
	nextPlay := make(chan string)

	go spotify.Initialise(initialised, toPlay, nextPlay, status)

	if <-initialised {
		ui.Start(toPlay, nextPlay, status)
	}
}
