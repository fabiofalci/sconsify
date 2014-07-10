package main

import (
	"time"

	"github.com/fabiofalci/sconsify/spotify"
	sp "github.com/op/go-libspotify/spotify"
)

func main() {
	initialised := make(chan bool)
	status := make(chan string)
	toPlay := make(chan sp.Track)
	nextPlay := make(chan string)

	go spotify.Initialise(initialised, toPlay, nextPlay, status)

	<-initialised

	playlist := spotify.Playlists["Ramones"]
	playlist.Wait()
	track := playlist.Track(3).Track()
	track.Wait()

	toPlay <- *track

	println(track.Name())

	for {
		time.Sleep(100 * time.Millisecond)
	}
}
