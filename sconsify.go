package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/rpc"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/fabiofalci/sconsify/ui/noui"
	"github.com/fabiofalci/sconsify/ui/simple"
	"github.com/howeyc/gopass"
	"runtime"
	"strconv"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/fabiofalci/sconsify/webapi"
)

var version string
var commit string
var buildDate string
var spotifyClientId string
var authRedirectUrl string

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	infrastructure.ProcessSconsifyrc()

	providedUsername := flag.String("username", "", "Spotify username.")
	providedWebApi := flag.Bool("web-api", true, "Use Spotify WEB API for more features. It requires web authorization.")
	providedOpenBrowser := flag.String("open-browser-cmd", "", "Open browser command to complete the web authorization.")
	providedStatusFile := flag.String("status-file", "", "File that sconsify will output status such as track being played.")
	providedStatusFileTemplate := flag.String("status-file-template", "", "Status file template.")
	providedUi := flag.Bool("ui", true, "Run Sconsify with Console User Interface. If false then no User Interface will be presented and it'll shuffle tracks.")
	providedPlaylists := flag.String("playlists", "", "Select just some Playlists to play. Comma separated list.")
	providedPreferredBitrate := flag.String("preferred-bitrate", "320k", "Preferred bitrate: 96k, 160k, 320k.")
	providedNoUiSilent := flag.Bool("noui-silent", false, "Silent mode when no UI is used.")
	providedNoUiRepeatOn := flag.Bool("noui-repeat-on", true, "Play your playlist and repeat it after the last track.")
	providedNoUiShuffle := flag.Bool("noui-shuffle", true, "Shuffle tracks or follow playlist order.")
	providedWebApiCacheToken := flag.Bool("web-api-cache-token", true, "Cache the web-api token as plain text in ~/.sconsify until its expiration.")
	providedIssueWebApiToken := flag.Bool("issue-web-api-token", false, "Issue a new web-api token.")
	providedWebApiCacheContent := flag.Bool("web-api-cache-content", true, "Cache some of the web-api content as plain text in ~/.sconsify.")
	providedDebug := flag.Bool("debug", false, "Enable debug mode.")
	askingVersion := flag.Bool("version", false, "Print version.")
	providedCommand := flag.String("command", "", "Execute a command in the server: replay, play_pause, next")
	providedServer := flag.Bool("server", true, "Start a background server to accept commands.")
	flag.Parse()

	if *askingVersion {
		fmt.Println("Version: " + version)
		fmt.Println("Git commit: " + commit)
		if i, err := strconv.ParseInt(buildDate, 10, 64); err == nil {
			fmt.Println("Build date: " + time.Unix(i, 0).UTC().String())
		}
		fmt.Println("Go version: " + runtime.Version())
		os.Exit(0)
	}

	if *providedDebug {
		infrastructure.InitialiseLogger()
		defer infrastructure.CloseLogger()
	}

	if *providedCommand != "" {
		rpc.Client(*providedCommand)
		return
	}

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")

	initConf := &spotify.SpotifyInitConf{
		WebApiAuth:         *providedWebApi,
		PlaylistFilter:     *providedPlaylists,
		PreferredBitrate:   *providedPreferredBitrate,
		CacheWebApiToken:   *providedWebApiCacheToken,
		CacheWebApiContent: *providedWebApiCacheContent,
		SpotifyClientId:    spotifyClientId,
		AuthRedirectUrl:    authRedirectUrl,
		OpenBrowserCommand: *providedOpenBrowser,
	}

	if *providedIssueWebApiToken {
		webapi.Auth(initConf.SpotifyClientId, initConf.AuthRedirectUrl, initConf.CacheWebApiToken, initConf.OpenBrowserCommand)
		return
	}

	username, pass := credentials(providedUsername)
	events := sconsify.InitialiseEvents()
	publisher := &sconsify.Publisher{}

	if *providedStatusFile != "" {
		statusFileTemplate := "{{.Action}}: {{.Track}} - {{.Artist}}\n"
		if *providedStatusFileTemplate!= "" {
			statusFileTemplate = *providedStatusFileTemplate
		}
		go ui.ToStatusFile(*providedStatusFile, statusFileTemplate)
	}

	go spotify.Initialise(initConf, username, pass, events, publisher)

	if *providedServer {
		go rpc.StartServer(publisher)
	}

	if *providedUi {
		ui := simple.InitialiseConsoleUserInterface(events, publisher, true)
		sconsify.StartMainLoop(events, publisher, ui, false)
	} else {
		var output noui.Printer
		if *providedNoUiSilent {
			output = new(noui.SilentPrinter)
		}
		ui := noui.InitialiseNoUserInterface(events, publisher, output, providedNoUiRepeatOn, providedNoUiShuffle)
		sconsify.StartMainLoop(events, publisher, ui, true)
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
	b, _ := gopass.GetPasswdMasked()
	return b
}
