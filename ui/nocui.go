package ui

import (
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/fabiofalci/sconsify/events"
	sp "github.com/op/go-libspotify/spotify"
)

func StartNoUserInterface(events *events.Events, silent *bool) {
	playlists := <-events.WaitForPlaylists()
	go listenForKeyboardEvents(events.NextPlay)

	listenForTermination(events)

	tracks := getTracksInRandomOrder(playlists)
	nextToPlayIndex := 0
	numberOfTracks := len(tracks)

	for {
		track := tracks[nextToPlayIndex]

		events.ToPlay <- track

		if *silent {
			<-events.WaitForStatus()
		} else {
			println(<-events.WaitForStatus())
		}
		<-events.NextPlay

		nextToPlayIndex++
		if nextToPlayIndex >= numberOfTracks {
			nextToPlayIndex = 0
		}
	}
}

func listenForTermination(events *events.Events) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			events.Shutdown()
			<-events.WaitForShutdown()
			os.Exit(0)
		}
	}()
}

func listenForKeyboardEvents(nextPlay chan bool) {
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
			println()
			nextPlay <- true
		}
	}
}

func getTracksInRandomOrder(playlists map[string]*sp.Playlist) []*sp.Track {
	numberOfTracks := 0
	for _, playlist := range playlists {
		playlist.Wait()
		numberOfTracks += playlist.Tracks()
	}

	tracks := make([]*sp.Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)
	permIndex := 0

	for _, playlist := range playlists {
		playlist.Wait()
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i).Track()
			track.Wait()

			tracks[perm[permIndex]] = track
			permIndex++
		}
	}

	return tracks
}

func getRandomPermutation(numberOfTracks int) []int {
	rand.Seed(time.Now().Unix())
	return rand.Perm(numberOfTracks)
}
