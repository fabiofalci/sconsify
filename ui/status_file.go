package ui

import (
	"github.com/fabiofalci/sconsify/sconsify"
	"text/template"
	"io/ioutil"
	"bytes"
)

type StatusTrack struct {
	Action string
	Track  string
	Artist string
}

var fileName string

func toFile(content []byte) {
	ioutil.WriteFile(fileName, content, 0600)
}

func cleanStatusFile() {
	toFile([]byte(""))
}

func ToStatusFile(file string, text string) {
	fileName = file
	toFileEvents := sconsify.InitialiseEvents()

	t := template.Must(template.New("statusTemplate").Parse(text))

	cleanStatusFile()
	for {
		select {
		case track := <-toFileEvents.TrackPausedUpdates():
			var b bytes.Buffer
			t.Execute(&b, StatusTrack{Action: "Paused", Track: track.Name, Artist: track.Artist.Name})
			toFile(b.Bytes())
		case track := <-toFileEvents.TrackPlayingUpdates():
			var b bytes.Buffer
			t.Execute(&b, StatusTrack{Action: "Playing", Track: track.Name, Artist: track.Artist.Name})
			toFile(b.Bytes())
		case <-toFileEvents.ShutdownEngineUpdates():
			cleanStatusFile()
			break
		case <-toFileEvents.TrackNotAvailableUpdates():
		case <-toFileEvents.PlayTokenLostUpdates():
		case <-toFileEvents.NextPlayUpdates():
		case <-toFileEvents.PlaylistsUpdates():
		case <-toFileEvents.ArtistAlbumsUpdates():
		case <-toFileEvents.NewTrackLoadedUpdate():
		case <-toFileEvents.TokenExpiredUpdates():
		case <-toFileEvents.ShutdownSpotifyUpdates():
		case <-toFileEvents.SearchUpdates():
		case <-toFileEvents.PlayUpdates():
		case <-toFileEvents.ReplayUpdates():
		case <-toFileEvents.PauseUpdates():
		case <-toFileEvents.PlayPauseToggleUpdates():
		case <-toFileEvents.GetArtistAlbumsUpdates():
		}
	}
}

