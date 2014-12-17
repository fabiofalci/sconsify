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
	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
)

type NoUi struct {
	silent *bool
	tracks []*sp.Track
	events *events.Events
}

func StartNoUserInterface(events *events.Events, silent *bool, repeatOn *bool, random *bool) error {
	noui := &NoUi{silent: silent, events: events}

	playlists := noui.waitForPlaylists()
	if playlists == nil {
		return nil
	}

	go noui.listenForKeyboardEvents()

	noui.listenForTermination()

	if err := noui.setTracks(playlists, random); err != nil {
		return err
	}

	if !*silent {
		noui.printPlaylistInfo()
	}

	nextToPlayIndex := 0
	numberOfTracks := len(noui.tracks)

	for {
		track := noui.tracks[nextToPlayIndex]

		events.Play(track)

		goToNext := false
		if !*silent {
			select {
			case <-events.WaitForTrackNotAvailable():
				goToNext = true
			case track := <-events.WaitForTrackPlaying():
				fmt.Println(formatTrack("Playing", track))
			}
		}

		if !goToNext {
			select {
			case <-events.WaitForTrackNotAvailable():
			case <-events.WaitForNextPlay():
			case <-events.WaitForPlayTokenLost():
				fmt.Printf("Play token lost\n")
				return nil
			}
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

func (noui *NoUi) waitForPlaylists() map[string]*sconsify.Playlist {
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
			noui.events.NextPlay()
		} else if key == "p" {
			fmt.Println("")
			noui.events.Pause()
		} else if key == "q" {
			noui.shutdownNogui()
		}
	}
}

func (noui *NoUi) setTracks(playlists map[string]*sconsify.Playlist, random *bool) error {
	numberOfTracks := 0
	for _, playlist := range playlists {
		numberOfTracks += playlist.Tracks()
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	noui.tracks = make([]*sp.Track, numberOfTracks)

	var perm []int
	if *random {
		perm = getRandomPermutation(numberOfTracks)
	}
	index := 0

	for _, playlist := range playlists {
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i)

			if *random {
				noui.tracks[perm[index]] = track
			} else {
				noui.tracks[index] = track
			}
			index++
		}
	}

	return nil
}

func getRandomPermutation(numberOfTracks int) []int {
	rand.Seed(time.Now().Unix())
	return rand.Perm(numberOfTracks)
}
