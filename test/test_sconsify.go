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
	"flag"
)

var output bytes.Buffer

var left = "h"
var right = "l"
var up = "k"
var down = "j"
var openClose = "space"
var search = "slash"
var enter = "Return"
var quit = "q"
var firstLine = "gg"
var lastLine = "G"

func main() {
	runTest := flag.Bool("run-test", false, "Run the test sequence.")
	flag.Parse()

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	events := sconsify.InitialiseEvents()

	infrastructure.InitialiseLogger()
	defer infrastructure.CloseLogger()

	go mock.Initialise(events)

	if *runTest {
		go runTests()
	}

	ui := ui.InitialiseConsoleUserInterface(events, false)
	sconsify.StartMainLoop(events, ui, false)
	println(output.String())
	sleep() // otherwise gocui eventually fails to quit properly
}

func runTests() {
	sleep()
	viNavigation()
	folders()
	searching()
}

func viNavigation() {
	goToFirstPlaylist()

	assert("Bob Marley", "")
	cmdAndAssert(down, "My folder", "")
	cmdAndAssert(down, " Bob Marley and The Wailers", "")
	cmdAndAssert(down, " The Ramones", "")
	cmdAndAssert(down, "Ramones", "")
	cmdAndAssert(down, "Ramones", "")

	cmdAndAssert(up, " The Ramones", "")
	cmdAndAssert(up, " Bob Marley and The Wailers", "")
	cmdAndAssert(up, "My folder", "")
	cmdAndAssert(up, "Bob Marley", "")
	cmdAndAssert(up, "Bob Marley", "")

	cmd(firstLine)
	assert("Bob Marley", "")
	cmd(lastLine)
	assert("Ramones", "")

	cmdAndAssert(right, "Ramones", "I wanna be sedated")

}

func folders() {
	goToFirstPlaylist()

	cmdAndAssert(down, "My folder", "")
	cmdAndAssert(openClose, "[My folder]", "")
	cmdAndAssert(down, "Ramones", "")

	cmdAndAssert(up, "[My folder]", "")
	cmdAndAssert(openClose, "My folder", "")
	cmdAndAssert(down, " Bob Marley and The Wailers", "")
	cmdAndAssert(down, " The Ramones", "")
	cmdAndAssert(down, "Ramones", "")

	cmdAndAssert(up, " The Ramones", "")
	cmdAndAssert(up, " Bob Marley and The Wailers", "")
	cmdAndAssert(up, "My folder", "")
}

func searching() {
	goToFirstPlaylist()

	cmd(search)
	cmds("elvis")
	cmd(enter)

	goToFirstPlaylist()
	assert("*Search", "")
	cmdAndAssert(openClose, "[*Search]", "")
	cmdAndAssert(openClose, "*Search", "")

	cmd(quit)
}

func goToFirstPlaylist() {
	cmd(left)
	cmd(left)
	cmds(firstLine)
}

func cmdAndAssert(key string, expectedPlaylist string, expectedTrack string) {
	cmd(key)
	assert(expectedPlaylist, expectedTrack)
}

func assert(expectedPlaylist string, expectedTrack string) {
	if valid, actualPlaylist := ui.CuiAssertSelectedPlaylist(expectedPlaylist); !valid {
		output.WriteString(fmt.Sprintf("Playlist '%v' not found on position but '%v'", expectedPlaylist, actualPlaylist))
		cmd("q")
		panic("Boom!")
	}
}

func cmds(keys string) {
	for _, key := range keys {
		cmd(string(key))
	}
}

func cmd(key string) {
	exec.Command("xdotool", "key", key).Run()
}

func sleep() {
	time.Sleep(500 * time.Millisecond)
}

