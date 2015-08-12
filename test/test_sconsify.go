package main

import (
	"fmt"
	"time"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/fabiofalci/sconsify/spotify/mock"
	"os/exec"
	"bytes"
	"strconv"
)

var output bytes.Buffer

func main() {
	infrastructure.ProcessSconsifyrc()

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	events := sconsify.InitialiseEvents()

	infrastructure.InitialiseLogger()
	defer infrastructure.CloseLogger()

	go mock.Initialise(events)

	go testSequence()
	ui := ui.InitialiseConsoleUserInterface(events)
	sconsify.StartMainLoop(events, ui, false)
	println(output.String())
}

func testSequence() {
	sleep()

	cmd("h")
	cmd("h")
	cmd("g")
	cmdAndAssert("g", "Bob Marley", "", 1)
	cmdAndAssert("j", "My folder", "", 2)
	cmdAndAssert("space", "[My folder]", "", 2)
	cmdAndAssert("space", "My folder", "", 3)

	cmd("q")
}

func cmdAndAssert(key string, selectedPlaylist string, selectedTrack string, playlistPosition int) {
	cmd(key)
	if !ui.CuiAssertSelectedPlaylist(selectedPlaylist) {
		output.WriteString("Playlist '" + selectedPlaylist + "' not found on position " + strconv.Itoa(playlistPosition))
		cmd("q")
	}
	sleep()
}

func cmd(key string) {
	exec.Command("xdotool", "key", key).Run()
}

func sleep() {
	time.Sleep(1 * time.Second)
}

