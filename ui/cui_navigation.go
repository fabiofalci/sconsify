package ui

import (
	"github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	if currentView == nil || currentView.Name() == "side" {
		gui.tracksView.Highlight = true
		gui.playlistsView.Highlight = false
		return gui.g.SetCurrentView("main")
	}
	gui.tracksView.Highlight = false
	gui.playlistsView.Highlight = true
	return gui.g.SetCurrentView("side")
}

func cursorEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		selectedPlaylist, err := gui.getSelectedPlaylist()
		if err == nil {
			playlist := playlists[selectedPlaylist]
			if playlist != nil {
				playlist.Wait()
				newIndex := playlist.Tracks() - 1

				ox, _ := v.Origin()
				cx, _ := v.Cursor()
				_, sizeY := v.Size()
				sizeY--

				if newIndex > sizeY {
					v.SetOrigin(ox, newIndex-sizeY)
					v.SetCursor(cx, sizeY)
				} else {
					v.SetCursor(cx, newIndex)
				}
			}
		}
	}
	return nil
}

func cursorHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		v.SetCursor(cx, 0)
		v.SetOrigin(ox, 0)
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if canGoToNewPosition(cy + 1) {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			if v == gui.playlistsView {
				gui.updateTracksView()
			}
		}
	}
	return nil
}

func canGoToNewPosition(newPosition int) bool {
	currentView := gui.g.CurrentView()
	line, err := currentView.Line(newPosition)
	if err != nil || len(line) == 0 {
		return false
	}
	return true
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
		if v == gui.playlistsView {
			gui.updateTracksView()
		}
	}
	return nil
}
