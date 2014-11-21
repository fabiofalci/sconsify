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
	providedNoUiSilent := flag.Bool("noui-silent", false, "Silent mode when no User Interface is used.")
	providedNoUiPlaylists := flag.String("noui-playlists", "", "Select just some Playlists to play when no User Interface is used. Comma separated list.")
	flag.Parse()

	username, pass := credentials(providedUsername)
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events)

	if *providedUi {
		ui.StartConsoleUserInterface(events)
	} else {
		err := ui.StartNoUserInterface(events, providedNoUiSilent, providedNoUiPlaylists)
		if err != nil {
			fmt.Println(err)
		}
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
	return &username, getPassword()
}

func getPassword() *[]byte {
	passFromEnv := os.Getenv("SCONSIFY_PASSWORD")
	if passFromEnv != "" {
		fmt.Println("Password from environment variable SCONSIFY_PASSWORD.")
		fmt.Println("Unset if you don't want to use it (unset SCONSIFY_PASSWORD).")
		b := []byte(passFromEnv)
		return &b
	}
	fmt.Print("Password: ")
	b := gopass.GetPasswd()
	return &b
}
