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
	playlists.AddPlaylist("name", sconsify.CreateDummyPlaylist())
	events.NewPlaylist(playlists)

	message := <-output.message
	if message != "4 track(s)" {
		t.Errorf("Should be playing 4 tracks")
	}

	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist0 - name0 [duration0]" {
		t.Errorf("Not showing right track, instead [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist1 - name1 [duration1]" {
		t.Errorf("Not showing right track, instead [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist2 - name2 [duration2]" {
		t.Errorf("Not showing right track, instead [%v]", message)
	}

	events.NextPlay()
	events.TrackPlaying(<-events.WaitPlay())
	message = <-output.message
	if message != "Playing: artist3 - name3 [duration3]" {
		t.Errorf("Not showing right track, instead [%v]", message)
	}

	events.NextPlay()

	<-finished
}
