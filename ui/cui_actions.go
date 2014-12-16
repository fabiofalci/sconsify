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
	gui.updateStatus(state.currentMessage, false)
	return nil
}

func setAllRandomMode(g *gocui.Gui, v *gocui.View) error {
	state.invertMode(allRandomMode)
	gui.updateStatus(state.currentMessage, false)
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

func removeFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	index := gui.getQeueuSelectedTrackIndex()
	if index > -1 {
		queue.Remove(index)
		gui.updateQueueView()
	}
	return nil
}

func enableSearchInputCommand(g *gocui.Gui, v *gocui.View) error {
	gui.statusView.Clear()
	gui.statusView.SetCursor(0, 0)
	gui.statusView.SetOrigin(0, 0)

	gui.statusView.Editable = true
	gui.g.SetCurrentView("status")

	return nil
}

func searchCommand(g *gocui.Gui, v *gocui.View) error {
	// after the enter the command is at -1
	line, _ := gui.statusView.Line(-1)

	fmt.Fprintf(gui.playlistsView, line)
	gui.g.SetCurrentView("side")
	gui.events.Search(line)

	gui.statusView.Clear()
	gui.statusView.SetCursor(0, 0)
	gui.statusView.SetOrigin(0, 0)

	gui.statusView.Editable = false
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	gui.events.Shutdown()
	<-gui.events.WaitForShutdown()
	return gocui.ErrorQuit
}
