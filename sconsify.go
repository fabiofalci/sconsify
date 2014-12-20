package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/events"
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
	flag.Parse()

	username, pass := credentials(providedUsername)
	events := events.InitialiseEvents()

	go spotify.Initialise(username, pass, events, providedPlaylists)

	if *providedUi {
		ui.StartConsoleUserInterface(events)
	} else {
		var output ui.Printer
		if *providedNoUiSilent {
			output = new(ui.SilentPrinter)
		}
		err := ui.StartNoUserInterface(events, output, providedNoUiRepeatOn, providedNoUiRandom)
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
	b := gopass.GetPasswdMasked()
	return &b
}
