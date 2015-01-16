package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/fabiofalci/sconsify/sconsify"
)

type NoUi struct {
	output    Printer
	random    bool
	repeatOn  bool
	playlists *sconsify.Playlists
	events    *sconsify.Events
}

type Printer interface {
	Print(message string)
}

type SilentPrinter struct{}
type StandardOutputPrinter struct{}

func InitialiseNoUserInterface(events *sconsify.Events, output Printer, repeatOn *bool, random *bool) sconsify.UserInterface {
	if output == nil {
		output = new(StandardOutputPrinter)
	}
	noui := &NoUi{
		output:   output,
		random:   *random,
		repeatOn: *repeatOn,
		events:   events,
	}

	go noui.listenForTermination()
	go noui.listenForKeyboardEvents()
	return noui
}

func (noui *NoUi) TrackPaused(track *sconsify.Track) {
	noui.output.Print(fmt.Sprintf("Paused: %v\n", track.GetFullTitle()))
}

func (noui *NoUi) TrackNotAvailable(track *sconsify.Track) {
	go noui.events.NextPlay()
}

func (noui *NoUi) PlayTokenLost() error {
	noui.output.Print("Play token lost\n")
	return errors.New("Play token lost")
}

func (noui *NoUi) TrackPlaying(track *sconsify.Track) {
	noui.output.Print(fmt.Sprintf("Playing: %v\n", track.GetFullTitle()))
}

func (noui *NoUi) GetNextToPlay() *sconsify.Track {
	if noui.playlists != nil {
		track, repeating := noui.playlists.GetNext()
		if repeating && !noui.repeatOn {
			go noui.Shutdown()
		} else {
			return track
		}
	}
	return nil
}

func (noui *NoUi) NewPlaylists(playlists sconsify.Playlists) error {
	if playlists.Tracks() == 0 {
		noui.output.Print("No track selected\n")
		return errors.New("No track selected")
	}
	if noui.random {
		playlists.SetMode(sconsify.AllRandomMode)
	} else {
		playlists.SetMode(sconsify.SequentialMode)
	}

	noui.output.Print(fmt.Sprintf("%v track(s)\n", playlists.PremadeTracks()))
	noui.playlists = &playlists
	return nil
}

func (noui *NoUi) listenForTermination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		noui.Shutdown()
	}
}

func (noui *NoUi) Shutdown() {
	noui.events.ShutdownEngine()
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
			noui.Shutdown()
		}
	}
}

func (p *SilentPrinter) Print(message string) {
}

func (p *StandardOutputPrinter) Print(message string) {
	fmt.Print(message)
}
