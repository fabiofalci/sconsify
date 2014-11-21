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

type NoGui struct {
	silent         *bool
	playlistFilter []string
	tracks         []*sp.Track
}

func StartNoUserInterface(events *events.Events, silent *bool, playlistFilter *string, repeatOn *bool) error {
	nogui := &NoGui{silent: silent}
	nogui.setPlaylistFilter(*playlistFilter)

	playlists := <-events.WaitForPlaylists()
	go listenForKeyboardEvents(events.NextPlay)

	listenForNoCuiTermination(events)

	err := nogui.randomTracks(playlists)
	if err != nil {
		return err
	}

	if !*silent {
		nogui.printPlaylistInfo()
	}

	nextToPlayIndex := 0
	numberOfTracks := len(nogui.tracks)

	for {
		track := nogui.tracks[nextToPlayIndex]

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

func (nogui *NoGui) printPlaylistInfo() {
	fmt.Printf("%v track(s) from ", len(nogui.tracks))
	if nogui.playlistFilter == nil {
		fmt.Printf("all playlist(s)\n")
	} else {
		fmt.Printf("%v playlist(s)\n", len(nogui.playlistFilter))
	}
}

func listenForNoCuiTermination(events *events.Events) {
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
			fmt.Println("")
			nextPlay <- true
		}
	}
}

func (nogui *NoGui) randomTracks(playlists map[string]*sp.Playlist) error {
	numberOfTracks := 0
	for _, playlist := range playlists {
		playlist.Wait()
		if nogui.isOnFilter(playlist.Name()) {
			numberOfTracks += playlist.Tracks()
		}
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	nogui.tracks = make([]*sp.Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)
	permIndex := 0

	for _, playlist := range playlists {
		playlist.Wait()
		if nogui.isOnFilter(playlist.Name()) {
			for i := 0; i < playlist.Tracks(); i++ {
				track := playlist.Track(i).Track()
				track.Wait()

				nogui.tracks[perm[permIndex]] = track
				permIndex++
			}
		}
	}

	return nil
}

func (nogui *NoGui) setPlaylistFilter(playlistFilter string) {
	if playlistFilter == "" {
		return
	}
	nogui.playlistFilter = strings.Split(playlistFilter, ",")
	for i := range nogui.playlistFilter {
		nogui.playlistFilter[i] = strings.Trim(nogui.playlistFilter[i], " ")
	}
}

func (nogui *NoGui) isOnFilter(playlist string) bool {
	if nogui.playlistFilter == nil {
		return true
	}
	for _, filter := range nogui.playlistFilter {
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
