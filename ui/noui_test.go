package ui

import (
	"strings"
	"testing"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/sconsify"
)

type TestPrinter struct {
	message chan string
}

func (p *TestPrinter) Print(message string) {
	p.message <- strings.Trim(message, " \n\r")
}

func TestNoUiEmptyPlaylists(t *testing.T) {
	repeatOn := true
	random := true
	events := events.InitialiseEvents()
	go func() {
		playlists := sconsify.InitPlaylists()
		events.NewPlaylist(playlists)
	}()

	err := StartNoUserInterface(events, nil, &repeatOn, &random)
	if err == nil {
		t.Errorf("No track selected should return an error")
	}
}

func TestNoUiSequentialAndRepeating(t *testing.T) {
	repeatOn := false
	random := false
	events := events.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}

	finished := make(chan bool)
	go func() {
		StartNoUserInterface(events, output, &repeatOn, &random)
		finished <- true
	}()

	playlists := sconsify.InitPlaylists()
	playlists.AddPlaylist("name", createDummyPlaylist())
	events.NewPlaylist(playlists)

	message := <-output.message
	if message != "4 track(s)" {
		t.Errorf("Should be playing 4 tracks")
	}

	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist0 - name0 [duration0]" {
		t.Errorf("Should be showing track0 instead is showing [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist1 - name1 [duration1]" {
		t.Errorf("Should be showing track1 instead is showing [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist2 - name2 [duration2]" {
		t.Errorf("Should be showing track2 instead is showing [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist3 - name3 [duration3]" {
		t.Errorf("Should be showing track3 instead is showing [%v]", message)
	}

	events.NextPlay()

	<-finished
}

func createDummyPlaylist() *sconsify.Playlist {
	tracks := make([]*sconsify.Track, 4)
	tracks[0] = sconsify.InitTrack("0", "artist0", "name0", "duration0")
	tracks[1] = sconsify.InitTrack("1", "artist1", "name1", "duration1")
	tracks[2] = sconsify.InitTrack("2", "artist2", "name2", "duration2")
	tracks[3] = sconsify.InitTrack("3", "artist3", "name3", "duration3")
	return sconsify.InitPlaylist(tracks)
}
