package ui

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/sconsify"
)

type NoUi struct {
	silent *bool
	tracks []*sconsify.Track
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
				fmt.Println("Playing: " + track.GetFullTitle())
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

func (noui *NoUi) waitForPlaylists() *sconsify.Playlists {
	select {
	case playlists := <-noui.events.WaitForPlaylists():
		return &playlists
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

func (noui *NoUi) setTracks(playlists *sconsify.Playlists, random *bool) error {
	var err error
	noui.tracks, err = playlists.GetTracks(random)
	if err != nil {
		return err
	}

	return nil
}
