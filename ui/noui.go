package ui

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/events"
	sp "github.com/op/go-libspotify/spotify"
)

type NoUi struct {
	silent         *bool
	playlistFilter []string
	tracks         []*sp.Track
	events         *events.Events
}

func StartNoUserInterface(events *events.Events, silent *bool, playlistFilter *string, repeatOn *bool) error {
	noui := &NoUi{silent: silent, events: events}
	noui.setPlaylistFilter(*playlistFilter)

	playlists := noui.waitForPlaylists(events)
	if playlists == nil {
		return nil
	}

	go noui.listenForKeyboardEvents()

	listenForNoCuiTermination(events)

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

		if *silent {
			<-events.WaitForStatus()
		} else {
			fmt.Println(<-events.WaitForStatus())
		}
		<-events.NextPlay

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

func (noui *NoUi) waitForPlaylists(events *events.Events) map[string]*sp.Playlist {
	select {
	case playlists := <-events.WaitForPlaylists():
		if playlists != nil {
			return playlists
		}
	case <-events.WaitForShutdown():
	}
	return nil
}

func (noui *NoUi) printPlaylistInfo() {
	fmt.Printf("%v track(s) from ", len(noui.tracks))
	if noui.playlistFilter == nil {
		fmt.Printf("all playlist(s)\n")
	} else {
		fmt.Printf("%v playlist(s)\n", len(noui.playlistFilter))
	}
}

func listenForNoCuiTermination(events *events.Events) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			shutdownNogui(events)
		}
	}()
}

func shutdownNogui(events *events.Events) {
	events.Shutdown()
	<-events.WaitForShutdown()
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
		}
	}
}

func (noui *NoUi) randomTracks(playlists map[string]*sp.Playlist) error {
	numberOfTracks := 0
	for _, playlist := range playlists {
		playlist.Wait()
		if noui.isOnFilter(playlist.Name()) {
			numberOfTracks += playlist.Tracks()
		}
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	noui.tracks = make([]*sp.Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)
	permIndex := 0

	for _, playlist := range playlists {
		playlist.Wait()
		if noui.isOnFilter(playlist.Name()) {
			for i := 0; i < playlist.Tracks(); i++ {
				track := playlist.Track(i).Track()
				track.Wait()

				noui.tracks[perm[permIndex]] = track
				permIndex++
			}
		}
	}

	return nil
}

func (noui *NoUi) setPlaylistFilter(playlistFilter string) {
	if playlistFilter == "" {
		return
	}
	noui.playlistFilter = strings.Split(playlistFilter, ",")
	for i := range noui.playlistFilter {
		noui.playlistFilter[i] = strings.Trim(noui.playlistFilter[i], " ")
	}
}

func (noui *NoUi) isOnFilter(playlist string) bool {
	if noui.playlistFilter == nil {
		return true
	}
	for _, filter := range noui.playlistFilter {
		if filter == playlist {
			return true
		}
	}
	return false
}

func getRandomPermutation(numberOfTracks int) []int {
	rand.Seed(time.Now().Unix())
	return rand.Perm(numberOfTracks)
}
