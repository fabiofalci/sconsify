package test

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/howeyc/gopass"
)

func main2() {
	username, pass := credentials()
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)
	playlists := <-events.WaitForPlaylists()

	playlist := playlists["Ramones"]
	playlist.Wait()
	track := playlist.Track(3).Track()
	track.Wait()

	events.ToPlay <- track

	println(track.Name())
	<-events.WaitForStatus()
	<-events.NextPlay
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
