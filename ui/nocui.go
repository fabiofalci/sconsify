package ui

import (
	"math/rand"

	"github.com/fabiofalci/sconsify/events"
	sp "github.com/op/go-libspotify/spotify"
)

func StartNoUserInterface(events *events.Events, silent *bool) {
	playlists := <-events.WaitForPlaylists()

	allTracks := getAllTracks(playlists).Contents()

	for {
		index := rand.Intn(len(allTracks))
		track := allTracks[index]

		events.ToPlay <- track

		if *silent {
			<-events.WaitForStatus()
		} else {
			println(<-events.WaitForStatus())
		}
		<-events.NextPlay
	}
}

func getAllTracks(playlists map[string]*sp.Playlist) *Queue {
	queue := InitQueue()

	for _, playlist := range playlists {
		playlist.Wait()
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i).Track()
			track.Wait()
			queue.Add(track)
		}
	}

	return queue
}
