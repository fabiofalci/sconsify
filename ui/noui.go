package ui

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/fabiofalci/sconsify/events"
	sp "github.com/op/go-libspotify/spotify"
)

type NoUi struct {
	silent *bool
	tracks []*sp.Track
	events *events.Events
}

func StartNoUserInterface(events *events.Events, silent *bool, repeatOn *bool) error {
	noui := &NoUi{silent: silent, events: events}

	playlists := noui.waitForPlaylists()
	if playlists == nil {
		return nil
	}

	go noui.listenForKeyboardEvents()

	noui.listenForTermination()

	err := noui.randomTracks(playlists)
	if err != nil {
		return err
	}

	if !*silent {
		noui.printPlaylistInfo()
	}

	nextToPlayIndex := 0
	numberOfTracks := len(noui.tracks)

	for {
		track := noui.tracks[nextToPlayIndex]

		events.ToPlay <- track

		message := <-events.WaitForStatus()
		if !*silent {
			fmt.Println(message)
		}
		select {
		case <-events.NextPlay:
		case <-events.WaitForPlayTokenLost():
			fmt.Printf("Play token lost\n")
			return nil
		}

		nextToPlayIndex++
		if nextToPlayIndex >= numberOfTracks {
			if *repeatOn {
				nextToPlayIndex = 0
			} else {
				return nil
			}
		}
	}

	return nil
}

func (noui *NoUi) waitForPlaylists() map[string]*sp.Playlist {
	select {
	case playlists := <-noui.events.WaitForPlaylists():
		if playlists != nil {
			return playlists
		}
	case <-noui.events.WaitForShutdown():
	}
	return nil
}

func (noui *NoUi) printPlaylistInfo() {
	fmt.Printf("%v track(s)\n", len(noui.tracks))
}

func (noui *NoUi) listenForTermination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			noui.shutdownNogui()
		}
	}()
}

func (noui *NoUi) shutdownNogui() {
	noui.events.Shutdown()
	<-noui.events.WaitForShutdown()
	os.Exit(0)
}

func (noui *NoUi) listenForKeyboardEvents() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	// we could disable echo but I can't enable it back

	// do not display entered characters on the screen
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// defer exec.Command("stty", "-F", "/dev/tty", "echo")

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)

		key := string(b)
		if key == ">" {
			fmt.Println("")
			noui.events.NextPlay <- true
		} else if key == "q" {
			noui.shutdownNogui()
		}
	}
}

func (noui *NoUi) randomTracks(playlists map[string]*sp.Playlist) error {
	numberOfTracks := 0
	for _, playlist := range playlists {
		playlist.Wait()
		numberOfTracks += playlist.Tracks()
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	noui.tracks = make([]*sp.Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)
	permIndex := 0

	for _, playlist := range playlists {
		playlist.Wait()
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i).Track()
			track.Wait()

			noui.tracks[perm[permIndex]] = track
			permIndex++
		}
	}

	return nil
}

func getRandomPermutation(numberOfTracks int) []int {
	rand.Seed(time.Now().Unix())
	return rand.Perm(numberOfTracks)
}
