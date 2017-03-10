package main

import (
	"fmt"
	"time"

	"bytes"
	"flag"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/spotify/mock"
	"github.com/fabiofalci/sconsify/ui/simple"
	"os/exec"
)

var output bytes.Buffer

var left = "h"
var right = "l"
var up = "k"
var down = "j"
var quit = "q"
var firstLine = "gg"
var lastLine = "G"
var queue = "u"
var removeAll = "D"
var remove = "dd"

func main() {
	runTest := flag.Bool("run-test", false, "Run the test sequence.")
	flag.Parse()

	fmt.Println("Sconsify - your awesome Spotify music service in a text-mode interface.")
	events := sconsify.InitialiseEvents()
	publisher := &sconsify.Publisher{}

	infrastructure.InitialiseLogger()
	defer infrastructure.CloseLogger()

	go mock.Initialise(events, publisher)

	if *runTest {
		go runTests()
	}

	ui := simple.InitialiseConsoleUserInterface(events, publisher, false)
	sconsify.StartMainLoop(events, publisher, ui, false)
	println(output.String())
	sleep() // otherwise gocui eventually fails to quit properly
}

func runTests() {
	sleep()

	viNavigation()
	viNavigationJump()
	folders()
	queueing()
	queueingPlaylist()
	removingFromQueue()
	clearingQueue()
	searching()
	removingTrack()
	removingPlaylist()

	cmd(quit)
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

	cmd(lastLine)
	assert("Ramones", "")
	cmd(firstLine)
	assert("Bob Marley", "")

	cmd(lastLine)
	cmdAndAssert(right, "Ramones", "I wanna be sedated")
	cmdAndAssert(down, "Ramones", "Pet semetary")
	cmdAndAssert(down, "Ramones", "Judy is a punk")
	cmdAndAssert(down, "Ramones", "Judy is a punk")

	cmdAndAssert(up, "Ramones", "Pet semetary")
	cmdAndAssert(up, "Ramones", "I wanna be sedated")
	cmdAndAssert(up, "Ramones", "I wanna be sedated")

	cmd(lastLine)
	assert("Ramones", "Judy is a punk")
	cmd(firstLine)
	assert("Ramones", "I wanna be sedated")
}

func viNavigationJump() {
	goToFirstPlaylist()

	assert("Bob Marley", "")
	cmd("3")
	cmd(down)
	assert(" The Ramones", "")

	cmd("2")
	cmd(up)
	assert("My folder", "")
}

func folders() {
	goToFirstPlaylist()

	cmdAndAssert(down, "My folder", "")
	openClose()
	assert("[My folder]", "")
	cmdAndAssert(down, "Ramones", "")

	cmdAndAssert(up, "[My folder]", "")
	openClose()
	assert("My folder", "")
	cmdAndAssert(down, " Bob Marley and The Wailers", "")
	cmdAndAssert(down, " The Ramones", "")
	cmdAndAssert(down, "Ramones", "")

	cmdAndAssert(up, " The Ramones", "")
	cmdAndAssert(up, " Bob Marley and The Wailers", "")
	cmdAndAssert(up, "My folder", "")
}

func searching() {
	goToFirstPlaylist()

	search()
	cmd("elvis")
	enter()

	goToFirstPlaylist()
	assert("*Search", "")
	openClose()
	assert("[*Search]", "")
	openClose()
	assert("*Search", "")
}

func queueing() {
	goToFirstPlaylist()

	cmd(right)
	assert("Bob Marley", "Waiting in vain")
	cmd(queue)
	cmd(queue)
	cmd(down)
	cmd(queue)

	assertNextTrackFromQueue("Waiting in vain")
	assertNextTrackFromQueue("Waiting in vain")
	assertNextTrackFromQueue("Testing")
	assertNextTrackFromQueue("")
}

func queueingPlaylist() {
	goToFirstPlaylist()

	assert("Bob Marley", "")
	cmd(queue)

	assertNextTrackFromQueue("Waiting in vain")
	assertNextTrackFromQueue("Testing")
	assertNextTrackFromQueue("")
}

func removingFromQueue() {
	goToFirstPlaylist()

	assert("Bob Marley", "")
	cmd(queue)

	cmd(right)
	cmd(right)

	cmd(remove)

	assertNextTrackFromQueue("Testing")
	assertNextTrackFromQueue("")
}

func clearingQueue() {
	goToFirstPlaylist()

	assert("Bob Marley", "")
	cmd(queue)

	cmd(right)
	cmd(right)

	cmd(removeAll)

	assertNextTrackFromQueue("")
}

func removingTrack() {
	goToLastPlaylist()

	assert("Ramones", "")
	cmd(right)
	assert("Ramones", "I wanna be sedated")

	cmd(remove)
	assert("Ramones", "Pet semetary")
}

func removingPlaylist() {
	goToLastPlaylist()

	assert("Ramones", "")
	cmd(remove)
	assert("", "")

	cmdAndAssert(up, " The Ramones", "")
}

func goToFirstPlaylist() {
	cmd(left)
	cmd(left)
	cmd(firstLine)
}

func goToLastPlaylist() {
	cmd(left)
	cmd(left)
	cmd(lastLine)
}

func cmdAndAssert(key string, expectedPlaylist string, expectedTrack string) {
	cmd(key)
	assert(expectedPlaylist, expectedTrack)
}

func assert(expectedPlaylist string, expectedTrack string) {
	if valid, actualPlaylist := simple.CuiAssertSelectedPlaylist(expectedPlaylist); !valid {
		output.WriteString(fmt.Sprintf("Playlist '%v' not found on position but '%v'", expectedPlaylist, actualPlaylist))
		cmd("q")
		panic("Boom!")
	}
	if expectedTrack != "" {
		if valid, actualTrack := simple.CuiAssertSelectedTrack(expectedTrack); !valid {
			output.WriteString(fmt.Sprintf("Track '%v' not found but '%v'", expectedTrack, actualTrack))
			cmd("q")
			panic("Boom!")
		}
	}
}

func assertNextTrackFromQueue(expectedNextTrackFromQueue string) {
	if valid, actualTrack := simple.CuiAssertQueueNextTrack(expectedNextTrackFromQueue); !valid {
		output.WriteString(fmt.Sprintf("Track '%v' is not the next in the queue but '%v'", expectedNextTrackFromQueue, actualTrack))
		cmd("q")
		panic("Boom!")
	}
}

func cmd(keys string) {
	for _, key := range keys {
		exec.Command("xdotool", "key", string(key)).Run()
	}
}

func openClose() {
	exec.Command("xdotool", "key", "space").Run()
}

func search() {
	exec.Command("xdotool", "key", "slash").Run()
}

func enter() {
	exec.Command("xdotool", "key", "Return").Run()
}

func sleep() {
	time.Sleep(500 * time.Millisecond)
}
