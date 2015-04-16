package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/howeyc/gopass"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	sconsify.ProcessSconsifyrc()

	providedUsername := flag.String("username", "", "Spotify username.")
	providedUi := flag.Bool("ui", true, "Run Sconsify with Console User Interface. If false then no User Interface will be presented and it'll only random between Playlists.")
	providedPlaylists := flag.String("playlists", "", "Select just some Playlists to play. Comma separated list.")
	providedNoUiSilent := flag.Bool("noui-silent", false, "Silent mode when no UI is used.")
	providedNoUiRepeatOn := flag.Bool("noui-repeat-on", true, "Play your playlist and repeat it after the last track.")
	providedNoUiRandom := flag.Bool("noui-random", true, "Random between tracks or follow playlist order.")
	providedDebug := flag.Bool("debug", false, "Enable debug mode.")
	flag.Parse()

	if *providedDebug {
		sconsify.InitialiseLogger()
		defer sconsify.CloseLogger()
	}

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	username, pass := credentials(providedUsername)
	events := sconsify.InitialiseEvents()

	go spotify.Initialise(username, pass, events, providedPlaylists)

	if *providedUi {
		ui := ui.InitialiseConsoleUserInterface(events)
		sconsify.StartMainLoop(events, ui, false)
	} else {
		var output ui.Printer
		if *providedNoUiSilent {
			output = new(ui.SilentPrinter)
		}
		ui := ui.InitialiseNoUserInterface(events, output, providedNoUiRepeatOn, providedNoUiRandom)
		sconsify.StartMainLoop(events, ui, true)
	}
}

func credentials(providedUsername *string) (string, []byte) {
	username := ""
	if *providedUsername == "" {
		fmt.Print("Username: ")
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
	} else {
		username = *providedUsername
		fmt.Println("Provided username: " + username)
	}
	return strings.Trim(username, " \n\r"), getPassword()
}

func getPassword() []byte {
	passFromEnv := os.Getenv("SCONSIFY_PASSWORD")
	if passFromEnv != "" {
		fmt.Println("Reading password from environment variable SCONSIFY_PASSWORD.")
		return []byte(passFromEnv)
	}
	fmt.Print("Password: ")
	return gopass.GetPasswdMasked()
}
