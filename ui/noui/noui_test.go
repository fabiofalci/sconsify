package noui

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

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
	shuffle := true
	events := sconsify.InitialiseEvents()
	publisher := &sconsify.Publisher{}

	go func() {
		playlists := sconsify.InitPlaylists()
		publisher.NewPlaylist(playlists)
	}()

	ui := InitialiseNoUserInterface(events, publisher, nil, &repeatOn, &shuffle)
	err := sconsify.StartMainLoop(events, publisher, ui, true)
	if err == nil {
		t.Errorf("No track selected should return an error")
	}
}

func TestNoUiSequentialAndRepeating(t *testing.T) {
	repeatOn := true
	shuffle := false
	events := sconsify.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}
	publisher := &sconsify.Publisher{}
	ui := InitialiseNoUserInterface(events, publisher, output, &repeatOn, &shuffle)

	finished := make(chan bool)
	go func() {
		err := sconsify.StartMainLoop(events, publisher, ui, true)
		finished <- err == nil
	}()

	sendNewPlaylist(publisher)

	assertPrintFourTracks(t, events, output)

	assertFirstTrack(t, events, publisher, output)
	assertNextThreeTracks(t, events, publisher, output)
	assertRepeatingAllFourTracks(t, events, publisher, output)

	assertShutdown(t, ui, events, publisher, finished)
}

func TestNoUiSequentialAndNotRepeating(t *testing.T) {
	repeatOn := false
	shuffle := false
	events := sconsify.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}
	publisher := &sconsify.Publisher{}
	ui := InitialiseNoUserInterface(events, publisher, output, &repeatOn, &shuffle)

	finished := make(chan bool)
	go func() {
		sconsify.StartMainLoop(events, publisher, ui, true)
		finished <- true
	}()

	sendNewPlaylist(publisher)

	assertPrintFourTracks(t, events, output)

	assertFirstTrack(t, events, publisher, output)
	assertNextThreeTracks(t, events, publisher, output)
	assertNoNextTrack(events, publisher, finished)
}

func TestNoUiShuffleAndRepeating(t *testing.T) {
	rand.Seed(123456789) // repeatable

	repeatOn := true
	shuffle := true
	events := sconsify.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}
	publisher := &sconsify.Publisher{}
	ui := InitialiseNoUserInterface(events, publisher, output, &repeatOn, &shuffle)

	finished := make(chan bool)
	go func() {
		err := sconsify.StartMainLoop(events, publisher, ui, true)
		finished <- err == nil
	}()

	sendNewPlaylist(publisher)

	assertPrintFourTracks(t, events, output)

	assertShuffleFirstTrack(t, events, publisher, output)
	assertShuffleNextThreeTracks(t, events, publisher, output)
	assertShuffleRepeatingAllFourTracks(t, events, publisher, output)

	assertShutdown(t, ui, events, publisher, finished)
}

func TestNoUiShuffleAndNotRepeating(t *testing.T) {
	rand.Seed(123456789) // repeatable

	repeatOn := false
	shuffle := true
	events := sconsify.InitialiseEvents()
	output := &TestPrinter{message: make(chan string)}
	publisher := &sconsify.Publisher{}
	ui := InitialiseNoUserInterface(events, publisher, output, &repeatOn, &shuffle)

	finished := make(chan bool)
	go func() {
		sconsify.StartMainLoop(events, publisher, ui, true)
		finished <- true
	}()

	sendNewPlaylist(publisher)

	assertPrintFourTracks(t, events, output)

	assertShuffleFirstTrack(t, events, publisher, output)
	assertShuffleNextThreeTracks(t, events, publisher, output)
	assertNoNextTrack(events, publisher, finished)
}

func sendNewPlaylist(publisher *sconsify.Publisher) {
	playlists := sconsify.InitPlaylists()
	playlists.AddPlaylist(createDummyPlaylist())
	publisher.NewPlaylist(playlists)
}

func assertShutdown(t *testing.T, ui sconsify.UserInterface, events *sconsify.Events, publisher *sconsify.Publisher, finished chan bool) {
	go ui.Shutdown()

	// playing spotify shutdown here
	<-events.ShutdownSpotifyUpdates()
	publisher.ShutdownEngine()

	if !<-finished {
		t.Errorf("Not properly finished")
	}
}

func assertPrintFourTracks(t *testing.T, events *sconsify.Events, output *TestPrinter) {
	message := <-output.message
	if message != "4 track(s)" {
		t.Errorf("Should be playing 4 tracks")
	}
}

func assertNoNextTrack(events *sconsify.Events, publisher *sconsify.Publisher, finished chan bool) {
	publisher.NextPlay()

	// playing spotify shutdown here
	<-events.ShutdownSpotifyUpdates()
	publisher.ShutdownEngine()

	<-finished
}

func assertFirstTrack(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	publisher.TrackPlaying(<-events.PlayUpdates())
	message := <-output.message
	if message != "Playing: artist0 - name0 [duration0]" {
		t.Errorf("Should be showing track0 instead is showing [%v]", message)
	}
}

func assertShuffleFirstTrack(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	publisher.TrackPlaying(<-events.PlayUpdates())
	message := <-output.message
	if message != "Playing: artist3 - name3 [duration3]" {
		t.Errorf("Should be showing track3 instead is showing [%v]", message)
	}
}

func assertNextThreeTracks(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	playNext(t, events, publisher, output, []string{"1", "2", "3"})
}

func assertShuffleNextThreeTracks(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	playNext(t, events, publisher, output, []string{"0", "2", "1"})
}

func assertRepeatingAllFourTracks(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	playNext(t, events, publisher, output, []string{"0", "1", "2", "3"})
}

func assertShuffleRepeatingAllFourTracks(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter) {
	playNext(t, events, publisher, output, []string{"3", "0", "2", "1"})
}

func playNext(t *testing.T, events *sconsify.Events, publisher *sconsify.Publisher, output *TestPrinter, tracks []string) {
	for _, track := range tracks {
		publisher.NextPlay()
		publisher.TrackPlaying(<-events.PlayUpdates())
		message := <-output.message
		expectedMessage := fmt.Sprintf("Playing: artist%v - name%v [duration%v]", track, track, track)
		if message != expectedMessage {
			t.Errorf("Should be showing track%v instead is showing [%v]", track, message)
		}
	}
}

func createDummyPlaylist() *sconsify.Playlist {
	tracks := make([]*sconsify.Track, 4)
	artist0 := sconsify.InitArtist("artist0", "artist0")
	tracks[0] = sconsify.InitTrack("0", artist0, "name0", "duration0")
	artist1 := sconsify.InitArtist("artist1", "artist1")
	tracks[1] = sconsify.InitTrack("1", artist1, "name1", "duration1")
	artist2 := sconsify.InitArtist("artist2", "artist2")
	tracks[2] = sconsify.InitTrack("2", artist2, "name2", "duration2")
	artist3 := sconsify.InitArtist("artist3", "artist3")
	tracks[3] = sconsify.InitTrack("3", artist3, "name3", "duration3")
	return sconsify.InitPlaylist("0", "test", tracks)
}
