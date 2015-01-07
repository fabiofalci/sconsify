package ui

import (
	"fmt"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

func keybindings() error {
	views := []string{"main", "side", "queue"}
	for _, view := range views {
		if err := gui.g.SetKeybinding(view, 'p', 0, pauseCurrentSelectedTrack); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'r', 0, setRandomMode); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'R', 0, setAllRandomMode); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, '>', 0, nextCommand); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, '/', 0, enableSearchInputCommand); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'j', 0, cursorDown); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'k', 0, cursorUp); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'q', 0, quit); err != nil {
			return err
		}
	}

	if err := gui.g.SetKeybinding("main", gocui.KeySpace, 0, playCurrentSelectedTrack); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("main", gocui.KeyEnter, 0, playCurrentSelectedTrack); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", 'u', 0, queueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", 'd', 0, removeFromQueueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", 'D', 0, removeAllFromQueueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("status", gocui.KeyEnter, 0, searchCommand); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyHome, 0, cursorHome); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyEnd, 0, cursorEnd); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyPgup, 0, cursorPgup); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyPgdn, 0, cursorPgdn); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyArrowDown, 0, cursorDown); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyArrowUp, 0, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", gocui.KeyArrowLeft, 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", gocui.KeyArrowLeft, 0, nextView); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("side", gocui.KeyArrowRight, 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("main", gocui.KeyArrowRight, 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", 'h', 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", 'h', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("side", 'l', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("main", 'l', 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}

	return nil
}

func playCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	track := gui.getCurrentSelectedTrack()
	if track != nil {
		events.Play(track)
	}
	return nil
}

func pauseCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	events.Pause()
	return nil
}

func setRandomMode(g *gocui.Gui, v *gocui.View) error {
	playlists.InvertMode(sconsify.RandomMode)
	gui.updateStatus(gui.currentMessage)
	return nil
}

func setAllRandomMode(g *gocui.Gui, v *gocui.View) error {
	playlists.InvertMode(sconsify.AllRandomMode)
	gui.updateStatus(gui.currentMessage)
	return nil
}

func nextCommand(g *gocui.Gui, v *gocui.View) error {
	gui.playNext()
	return nil
}

func queueCommand(g *gocui.Gui, v *gocui.View) error {
	track := gui.getCurrentSelectedTrack()
	if track != nil {
		fmt.Fprintf(gui.queueView, "%v", track.GetTitle())
		queue.Add(track)
	}
	return nil
}

func removeAllFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	queue.RemoveAll()
	gui.updateQueueView()
	return gui.enableMainView()
}

func removeFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	index := gui.getQueueSelectedTrackIndex()
	if index > -1 {
		queue.Remove(index)
		gui.updateQueueView()
	}
	return nil
}

func enableSearchInputCommand(g *gocui.Gui, v *gocui.View) error {
	gui.clearStatusView()
	gui.statusView.Editable = true
	gui.g.SetCurrentView("status")
	return nil
}

func searchCommand(g *gocui.Gui, v *gocui.View) error {
	// after the enter the command is at -1
	line, _ := gui.statusView.Line(-1)

	gui.enableSideView()
	events.Search(line)
	gui.clearStatusView()
	gui.statusView.Editable = false
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	gui.Shutdown()
	// TODO wait for shutdown
	// <-events.ShutdownUpdates()
	return gocui.ErrorQuit
}
