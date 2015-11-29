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
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/howeyc/gopass"
)

var version string
var commit string
var buildDate string

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	infrastructure.ProcessSconsifyrc()

	providedUsername := flag.String("username", "", "Spotify username.")
	providedUi := flag.Bool("ui", true, "Run Sconsify with Console User Interface. If false then no User Interface will be presented and it'll shuffle tracks.")
	providedPlaylists := flag.String("playlists", "", "Select just some Playlists to play. Comma separated list.")
	providedNoUiSilent := flag.Bool("noui-silent", false, "Silent mode when no UI is used.")
	providedNoUiRepeatOn := flag.Bool("noui-repeat-on", true, "Play your playlist and repeat it after the last track.")
	providedNoUiShuffle := flag.Bool("noui-shuffle", true, "Shuffle tracks or follow playlist order.")
	providedDebug := flag.Bool("debug", false, "Enable debug mode.")
	askingVersion := flag.Bool("version", false, "Print version.")
	flag.Parse()

	if *askingVersion {
		fmt.Println("Version: " + version)
		fmt.Println("Git commit: " + commit)
		fmt.Println("Build date: " + buildDate)
		os.Exit(0)
	}

	if *providedDebug {
		infrastructure.InitialiseLogger()
		defer infrastructure.CloseLogger()
	}

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	username, pass := credentials(providedUsername)
	events := sconsify.InitialiseEvents()

	go spotify.Initialise(username, pass, events, providedPlaylists)

	if *providedUi {
		ui := ui.InitialiseConsoleUserInterface(events, true)
		sconsify.StartMainLoop(events, ui, false)
	} else {
		var output ui.Printer
		if *providedNoUiSilent {
			output = new(ui.SilentPrinter)
		}
		ui := ui.InitialiseNoUserInterface(events, output, providedNoUiRepeatOn, providedNoUiShuffle)
		sconsify.StartMainLoop(events, ui, true)
	}
}

func credentials(providedUsername *string) (string, []byte) {
	username := ""
	if *providedUsername == "" {
		fmt.Print("Premium account username: ")
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
