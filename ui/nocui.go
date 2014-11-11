package ui

import (
	"math/rand"
	"os"
	"os/exec"

	"github.com/fabiofalci/sconsify/events"
	sp "github.com/op/go-libspotify/spotify"
)

func StartNoUserInterface(events *events.Events, silent *bool) {
	playlists := <-events.WaitForPlaylists()
	go listenForKeyboardEvents(events.NextPlay)

	allTracks := getAllTracks(playlists).Contents()

	for {
		index := rand.Intn(len(allTracks))
		track := allTracks[index]

		events.ToPlay <- track

		if *silent {
			<-events.WaitForStatus()
		} else {
			println(<-events.WaitForStatus())
		}
		<-events.NextPlay
	}
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

func getAllTracks(playlists map[string]*sp.Playlist) *Queue {
	queue := InitQueue()

	for _, playlist := range playlists {
		playlist.Wait()
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i).Track()
			track.Wait()
			queue.Add(track)
		}
	}

	return queue
}
