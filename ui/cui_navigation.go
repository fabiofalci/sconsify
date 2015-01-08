package ui

import (
	"github.com/jroimartin/gocui"
)

// nextView is shared between Playlists and Queue and they all go to Tracks
func nextView(g *gocui.Gui, v *gocui.View) error {
	return gui.enableTracksView()
}

func mainNextViewLeft(g *gocui.Gui, v *gocui.View) error {
	return gui.enableSideView()
}

func mainNextViewRight(g *gocui.Gui, v *gocui.View) error {
	return gui.enableQueueView()
}

func cursorEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		newIndex := getCurrentViewSize(v)
		if newIndex > -1 {
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

			updateTracksView(g, v)
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

		updateTracksView(g, v)
	}
	return nil
}

func cursorPgup(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		_, pageSizeY := v.Size()
		pageSizeY--

		newOriginY := oy - pageSizeY
		if newOriginY > 0 {
			v.SetOrigin(ox, newOriginY)
			v.SetCursor(cx, cy)
		} else {
			v.SetOrigin(ox, 0)
			v.SetCursor(cx, cy)
		}
		updateTracksView(g, v)
	}
	return nil
}

func cursorPgdn(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		maxSize := getCurrentViewSize(v)
		if maxSize > -1 {
			ox, oy := v.Origin()
			cx, cy := v.Cursor()
			_, pageSizeY := v.Size()
			pageSizeY--

			newOriginY := oy + pageSizeY

			if hasMorePages(newOriginY, cy, maxSize) {
				v.SetOrigin(ox, newOriginY)
				v.SetCursor(cx, cy)
			} else if isNotInLastPage(oy, pageSizeY, maxSize) {
				v.SetOrigin(ox, maxSize-pageSizeY)
				v.SetCursor(cx, pageSizeY)
			}
			updateTracksView(g, v)
		}
	}
	return nil
}

func updateTracksView(g *gocui.Gui, v *gocui.View) {
	if v == gui.playlistsView {
		gui.updateTracksView()
	}
}

func getCurrentViewSize(v *gocui.View) int {
	if v == gui.tracksView {
		selectedPlaylist, err := gui.getSelectedPlaylist()
		if err == nil {
			playlist := playlists.Get(selectedPlaylist)
			if playlist != nil {
				return playlist.Tracks() - 1
			}
		}
	} else if v == gui.playlistsView {
		return playlists.Playlists() - 1
	}
	return -1
}

func hasMorePages(newOriginY int, cursorY int, maxSize int) bool {
	return newOriginY+cursorY <= maxSize
}

func isNotInLastPage(originY int, pageSizeY int, maxSize int) bool {
	return originY+pageSizeY <= maxSize
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
