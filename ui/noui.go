package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	e "github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/sconsify"
)

type NoUi struct {
	output           Printer
	internalShutdown chan bool
}

type Printer interface {
	Print(message string)
}

type SilentPrinter struct{}
type StandardOutputPrinter struct{}

var noui *NoUi

func StartNoUserInterface(ev *e.Events, output Printer, repeatOn *bool, random *bool) error {
	events = ev
	if output == nil {
		output = new(StandardOutputPrinter)
	}
	noui = &NoUi{
		output:           output,
		internalShutdown: make(chan bool),
	}

	playlists := noui.waitForPlaylists()
	if playlists == nil || playlists.Tracks() == 0 {
		return errors.New("No track selected")
	}

	go noui.listenForKeyboardEvents()

	noui.listenForTermination()

	if *random {
		playlists.SetMode(sconsify.AllRandomMode)
	} else {
		playlists.SetMode(sconsify.SequentialMode)
	}

	noui.output.Print(fmt.Sprintf("%v track(s)\n", playlists.PremadeTracks()))

	running := true
	for running {
		track, repeating := playlists.GetNext()
		if repeating && !*repeatOn {
			running = false
			break
		}

		events.Play(track)

		goToNext := false
		select {
		case <-events.WaitForTrackNotAvailable():
			goToNext = true
		case track := <-events.WaitForTrackPlaying():
			noui.output.Print(fmt.Sprintf("Playing: %v\n", track.GetFullTitle()))
		case <-events.WaitForPlayTokenLost():
			noui.output.Print("Play token lost\n")
			running = false
			break
		}

		if !goToNext {
			select {
			case <-noui.internalShutdown:
				running = false
			case <-events.WaitForNextPlay():
			case <-events.WaitForPlayTokenLost():
				noui.output.Print("Play token lost\n")
				running = false
			}
		}
	}

	return nil
}

func (noui *NoUi) waitForPlaylists() *sconsify.Playlists {
	select {
	case playlists := <-events.WaitForPlaylists():
		return &playlists
	case <-events.WaitForShutdown():
	}
	return nil
}

func (noui *NoUi) listenForTermination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			noui.shutdownNogui()
		}
	}()
}

func (noui *NoUi) shutdownNogui() {
	events.Shutdown()
	<-events.WaitForShutdown()
	noui.internalShutdown <- true
}

func ShutdownNogui() {
	noui.shutdownNogui()
}

func (noui *NoUi) listenForKeyboardEvents() {
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

	// we could disable echo but I can't enable it back

	// do not display entered characters on the screen
	// exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// defer exec.Command("stty", "-F", "/dev/tty", "echo")

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)

		key := string(b)
		if key == ">" {
			fmt.Println("")
			events.NextPlay()
		} else if key == "p" {
			fmt.Println("")
			events.Pause()
		} else if key == "q" {
			noui.shutdownNogui()
		}
	}
}

func (p *SilentPrinter) Print(message string) {
}

func (p *StandardOutputPrinter) Print(message string) {
	fmt.Print(message)
}
