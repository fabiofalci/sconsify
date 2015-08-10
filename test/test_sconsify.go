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
)

func main() {
	infrastructure.ProcessSconsifyrc()

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	events := sconsify.InitialiseEvents()

	go mock.Initialise(events)

	var output bytes.Buffer
	go generateKeys(&output)
	ui := ui.InitialiseConsoleUserInterface(events)
	sconsify.StartMainLoop(events, ui, false)
	println(output.String())
}

func generateKeys(output *bytes.Buffer) {
	sleep()

	cmd("h")
	cmd("h")
	cmd("g")
	cmd("g")
	if !ui.CuiAssertSelectedPlaylist("Bob Marley") {
		output.WriteString("Playlist Bob Marley not found on position 1")
		cmd("q")
	}
	sleep()

	cmd("j")
	if !ui.CuiAssertSelectedPlaylist("My folder") {
		output.WriteString("Playlist My folder not found on position 2")
		cmd("q")
	}
	sleep()

	cmd("space")
	if !ui.CuiAssertSelectedPlaylist("[My folder]") {
		output.WriteString("Playlist [My folder] not found on position 2")
		cmd("q")
	}
	sleep()

	cmd("space")
	if !ui.CuiAssertSelectedPlaylist("My folder") {
		output.WriteString("Playlist My folder not found on position 2")
		cmd("q")
	}
	sleep()



	cmd("q")
}


func sleep() {
	time.Sleep(1 * time.Second)
}

func cmd(key string) {
	exec.Command("xdotool", "key", key).Run()
}
