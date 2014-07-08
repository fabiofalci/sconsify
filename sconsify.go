package main

import (
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	sp "github.com/op/go-libspotify/spotify"
)

func main() {
	initialised := make(chan string)
	toPlay := make(chan sp.Track)

	go spotify.Initialise(initialised, toPlay)

	<-initialised

	ui.Start(toPlay)
}
