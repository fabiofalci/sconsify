package main

import (
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	sp "github.com/op/go-libspotify/spotify"
)

func main() {
	initialised := make(chan string)
	status := make(chan string)
	toPlay := make(chan sp.Track)
	nextPlay := make(chan string)

	go spotify.Initialise(initialised, toPlay, nextPlay, status)

	<-initialised

	ui.Start(toPlay, nextPlay, status)
}
