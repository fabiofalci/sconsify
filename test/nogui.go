package test

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	ui "github.com/fabiofalci/sconsify/ui"
	"github.com/howeyc/gopass"
	sp "github.com/op/go-libspotify/spotify"
)

func main2() {
	username, pass := credentials()
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)
	playlists := <-events.WaitForPlaylists()

	allTracks := getAllTracks(playlists).Contents()

	for {
		index := rand.Intn(len(allTracks))
		track := allTracks[index]

		events.ToPlay <- track

		println(<-events.WaitForStatus())
		<-events.NextPlay
	}
}

func getAllTracks(playlists map[string]*sp.Playlist) *ui.Queue {
	queue := ui.InitQueue()

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

func credentials() (*string, *[]byte) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, " \n\r")
	fmt.Print("Password: ")
	pass := gopass.GetPasswd()
	return &username, &pass
}
