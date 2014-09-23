package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

func playCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	track := getCurrentSelectedTrack()
	if track != nil {
		gui.events.Play(track)
	}
	return nil
}

func pauseCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	gui.events.Pause()
	return nil
}

func setRandomMode(g *gocui.Gui, v *gocui.View) error {
	state.invertMode(randomMode)
	gui.updateStatus(state.currentMessage)
	return nil
}

func setAllRandomMode(g *gocui.Gui, v *gocui.View) error {
	state.invertMode(allRandomMode)
	gui.updateStatus(state.currentMessage)
	return nil
}

func nextCommand(g *gocui.Gui, v *gocui.View) error {
	gui.playNext()
	return nil
}

func queueCommand(g *gocui.Gui, v *gocui.View) error {
	track := getCurrentSelectedTrack()
	if track != nil {
		fmt.Fprintf(gui.queueView, "%v - %v", track.Artist(0).Name(), track.Name())
		queue.Add(track)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	gui.events.Shutdown()
	<-gui.events.WaitForShutdown()
	return gocui.ErrorQuit
}
