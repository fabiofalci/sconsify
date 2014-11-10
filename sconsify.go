package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/howeyc/gopass"
)

func main() {
	username, pass := credentials()
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)
	ui.Start(events)
}

func credentials() (*string, *[]byte) {
	providedUsername := flag.String("username", "", "Spotify username")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	username := ""
	if *providedUsername == "" {
		fmt.Print("Username: ")
		username, _ = reader.ReadString('\n')
	} else {
		username = *providedUsername
	}
	username = strings.Trim(username, " \n\r")
	fmt.Print("Password: ")
	pass := gopass.GetPasswd()
	return &username, &pass
}
