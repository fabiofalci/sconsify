package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/sconsify"
)

type NoUi struct {
	events *events.Events
	output Printer
}

type Printer interface {
	Print(message string)
}

type SilentPrinter struct{}
type StandardOutputPrinter struct{}

func StartNoUserInterface(events *events.Events, output Printer, repeatOn *bool, random *bool) error {
	if output == nil {
		output = new(StandardOutputPrinter)
	}
	noui := &NoUi{events: events, output: output}

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

	for {
		track, repeating := playlists.GetNext()
		if repeating && !*repeatOn {
			return nil
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
			return nil
		}

		if !goToNext {
			select {
			case <-events.WaitForNextPlay():
			case <-events.WaitForPlayTokenLost():
				noui.output.Print("Play token lost\n")
				return nil
			}
		}
	}

	return nil
}

func (noui *NoUi) waitForPlaylists() *sconsify.Playlists {
	select {
	case playlists := <-noui.events.WaitForPlaylists():
		return &playlists
	case <-noui.events.WaitForShutdown():
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
	noui.events.Shutdown()
	<-noui.events.WaitForShutdown()
	os.Exit(0)
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
			noui.events.NextPlay()
		} else if key == "p" {
			fmt.Println("")
			noui.events.Pause()
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
