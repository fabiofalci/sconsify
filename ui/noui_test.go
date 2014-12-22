package ui

import (
	"fmt"
	"math/rand"
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
	repeatOn := true
	random := false
	events := events.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}

	finished := make(chan bool)
	go func() {
		err := StartNoUserInterface(events, output, &repeatOn, &random)
		finished <- err == nil
	}()

	sendNewPlaylist(events)

	assertPrintFourTracks(t, events, output)

	assertFirstTrack(t, events, output)
	assertNextThreeTracks(t, events, output)
	assertRepeatingAllFourTracks(t, events, output)

	assertShutdown(t, events, finished)
}

func TestNoUiSequentialAndNotRepeating(t *testing.T) {
	repeatOn := false
	random := false
	events := events.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}

	finished := make(chan bool)
	go func() {
		StartNoUserInterface(events, output, &repeatOn, &random)
		finished <- true
	}()

	sendNewPlaylist(events)

	assertPrintFourTracks(t, events, output)

	assertFirstTrack(t, events, output)
	assertNextThreeTracks(t, events, output)
	assertNoNextTrack(events, finished)
}

func TestNoUiRandomAndRepeating(t *testing.T) {
	rand.Seed(123456789) // repeatable

	repeatOn := true
	random := true
	events := events.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}

	finished := make(chan bool)
	go func() {
		err := StartNoUserInterface(events, output, &repeatOn, &random)
		finished <- err == nil
	}()

	sendNewPlaylist(events)

	assertPrintFourTracks(t, events, output)

	assertRandomFirstTrack(t, events, output)
	assertRandomNextThreeTracks(t, events, output)
	assertRandomRepeatingAllFourTracks(t, events, output)

	assertShutdown(t, events, finished)
}

func TestNoUiRandomAndNotRepeating(t *testing.T) {
	rand.Seed(123456789) // repeatable

	repeatOn := false
	random := true
	events := events.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}

	finished := make(chan bool)
	go func() {
		StartNoUserInterface(events, output, &repeatOn, &random)
		finished <- true
	}()

	sendNewPlaylist(events)

	assertPrintFourTracks(t, events, output)

	assertRandomFirstTrack(t, events, output)
	assertRandomNextThreeTracks(t, events, output)
	assertNoNextTrack(events, finished)
}

func sendNewPlaylist(events *events.Events) {
	playlists := sconsify.InitPlaylists()
	playlists.AddPlaylist("name", createDummyPlaylist())
	events.NewPlaylist(playlists)
}

func assertShutdown(t *testing.T, events *events.Events, finished chan bool) {
	go ShutdownNogui()

	<-events.WaitForShutdown()
	events.Shutdown()

	if !<-finished {
		t.Errorf("Not properly finished")
	}
}

func assertPrintFourTracks(t *testing.T, events *events.Events, output *TestPrinter) {
	message := <-output.message
	if message != "4 track(s)" {
		t.Errorf("Should be playing 4 tracks")
	}
}

func assertNoNextTrack(events *events.Events, finished chan bool) {
	events.NextPlay()
	<-finished
}

func assertFirstTrack(t *testing.T, events *events.Events, output *TestPrinter) {
	events.TrackPlaying(<-events.WaitPlay())
	message := <-output.message
	if message != "Playing: artist0 - name0 [duration0]" {
		t.Errorf("Should be showing track0 instead is showing [%v]", message)
	}
}

func assertRandomFirstTrack(t *testing.T, events *events.Events, output *TestPrinter) {
	events.TrackPlaying(<-events.WaitPlay())
	message := <-output.message
	if message != "Playing: artist3 - name3 [duration3]" {
		t.Errorf("Should be showing track3 instead is showing [%v]", message)
	}
}

func assertNextThreeTracks(t *testing.T, events *events.Events, output *TestPrinter) {
	playNext(t, events, output, []string{"1", "2", "3"})
}

func assertRandomNextThreeTracks(t *testing.T, events *events.Events, output *TestPrinter) {
	playNext(t, events, output, []string{"0", "2", "1"})
}

func assertRepeatingAllFourTracks(t *testing.T, events *events.Events, output *TestPrinter) {
	playNext(t, events, output, []string{"0", "1", "2", "3"})
}

func assertRandomRepeatingAllFourTracks(t *testing.T, events *events.Events, output *TestPrinter) {
	playNext(t, events, output, []string{"3", "0", "2", "1"})
}

func playNext(t *testing.T, events *events.Events, output *TestPrinter, tracks []string) {
	for _, track := range tracks {
		events.NextPlay()
		events.TrackPlaying(<-events.WaitPlay())
		message := <-output.message
		expectedMessage := fmt.Sprintf("Playing: artist%v - name%v [duration%v]", track, track, track)
		if message != expectedMessage {
			t.Errorf("Should be showing track%v instead is showing [%v]", track, message)
		}
	}
}

func createDummyPlaylist() *sconsify.Playlist {
	tracks := make([]*sconsify.Track, 4)
	tracks[0] = sconsify.InitTrack("0", "artist0", "name0", "duration0")
	tracks[1] = sconsify.InitTrack("1", "artist1", "name1", "duration1")
	tracks[2] = sconsify.InitTrack("2", "artist2", "name2", "duration2")
	tracks[3] = sconsify.InitTrack("3", "artist3", "name3", "duration3")
	return sconsify.InitPlaylist(tracks)
}
