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
	providedUsername := flag.String("username", "", "Spotify username.")
	providedUi := flag.Bool("ui", true, "Run Sconsify with Console User Interface. If false then no User Interface will be presented and it'll only random between Playlists.")
	flag.Parse()

	username, pass := credentials(providedUsername)
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)

	if *providedUi {
		ui.StartConsoleUserInterface(events)
	} else {
		ui.StartNoUserInterface(events)
	}
}

func credentials(providedUsername *string) (*string, *[]byte) {
	username := ""
	if *providedUsername == "" {
		fmt.Print("Username: ")
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
	} else {
		username = *providedUsername
	}
	username = strings.Trim(username, " \n\r")
	fmt.Print("Password: ")
	pass := gopass.GetPasswd()
	return &username, &pass
}
