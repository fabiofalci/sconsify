package ui

import (
	"fmt"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

func keybindings() error {
	views := []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE}
	allViews := ""
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

	if err := gui.g.SetKeybinding(VIEW_TRACKS, gocui.KeySpace, 0, playCurrentSelectedTrack); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_TRACKS, gocui.KeyEnter, 0, playCurrentSelectedTrack); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(VIEW_TRACKS, 'u', 0, queueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_PLAYLISTS, 'd', 0, removeFromPlaylistsCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_QUEUE, 'd', 0, removeFromQueueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_QUEUE, 'D', 0, removeAllFromQueueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_STATUS, gocui.KeyEnter, 0, searchCommand); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(allViews, gocui.KeyHome, 0, cursorHome); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(allViews, gocui.KeyEnd, 0, cursorEnd); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(allViews, gocui.KeyPgup, 0, cursorPgup); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(allViews, gocui.KeyPgdn, 0, cursorPgdn); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(allViews, gocui.KeyArrowDown, 0, cursorDown); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(allViews, gocui.KeyArrowUp, 0, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(VIEW_TRACKS, gocui.KeyArrowLeft, 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_QUEUE, gocui.KeyArrowLeft, 0, nextView); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(VIEW_PLAYLISTS, gocui.KeyArrowRight, 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_TRACKS, gocui.KeyArrowRight, 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(VIEW_TRACKS, 'h', 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_QUEUE, 'h', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_PLAYLISTS, 'l', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding(VIEW_TRACKS, 'l', 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding(allViews, gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}

	return nil
}

func playCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	if track := gui.getCurrentSelectedTrack(); track != nil {
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
	gui.updateCurrentStatus()
	return nil
}

func setAllRandomMode(g *gocui.Gui, v *gocui.View) error {
	playlists.InvertMode(sconsify.AllRandomMode)
	gui.updateCurrentStatus()
	return nil
}

func nextCommand(g *gocui.Gui, v *gocui.View) error {
	gui.playNext()
	return nil
}

func queueCommand(g *gocui.Gui, v *gocui.View) error {
	if track := gui.getCurrentSelectedTrack(); track != nil {
		fmt.Fprintf(gui.queueView, "%v", track.GetTitle())
		queue.Add(track)
	}
	return nil
}

func removeAllFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	queue.RemoveAll()
	gui.updateQueueView()
	return gui.enableTracksView()
}

func removeFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	if index := gui.getQueueSelectedTrackIndex(); index > -1 {
		queue.Remove(index)
		gui.updateQueueView()
	}
	return nil
}

func removeFromPlaylistsCommand(g *gocui.Gui, v *gocui.View) error {
	if playlist := gui.getSelectedPlaylist(); playlist != nil && playlist.IsSearch() {
		playlists.Remove(playlist.Name())
		gui.updatePlaylistsView()
		gui.updateTracksView()
	}
	return nil
}

func enableSearchInputCommand(g *gocui.Gui, v *gocui.View) error {
	gui.clearStatusView()
	gui.statusView.Editable = true
	gui.g.SetCurrentView(VIEW_STATUS)
	return nil
}

func searchCommand(g *gocui.Gui, v *gocui.View) error {
	// after user hit Enter, the typed command is at position -1
	if query, _ := gui.statusView.Line(-1); query != "" {
		events.Search(query)
	}
	gui.enableSideView()
	gui.clearStatusView()
	gui.statusView.Editable = false
	gui.updateCurrentStatus()
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	consoleUserInterface.Shutdown()
	// TODO wait for shutdown
	// <-events.ShutdownUpdates()
	return gocui.ErrorQuit
}
