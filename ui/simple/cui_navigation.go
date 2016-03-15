package simple

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
	if newIndex := getCurrentViewSize(v); newIndex > -1 {
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
	return nil
}

func cursorHome(g *gocui.Gui, v *gocui.View) error {
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	v.SetCursor(cx, 0)
	v.SetOrigin(ox, 0)

	updateTracksView(g, v)
	return nil
}

func cursorPgup(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, pageSizeY := v.Size()
	pageSizeY--

	if newOriginY := oy - pageSizeY; newOriginY > 0 {
		v.SetOrigin(ox, newOriginY)
		v.SetCursor(cx, cy)
	} else {
		v.SetOrigin(ox, 0)
		v.SetCursor(cx, cy)
	}
	updateTracksView(g, v)
	return nil
}

func cursorPgdn(g *gocui.Gui, v *gocui.View) error {
	if maxSize := getCurrentViewSize(v); maxSize > -1 {
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
	return nil
}

func updateTracksView(g *gocui.Gui, v *gocui.View) {
	if v == gui.playlistsView {
		gui.updateTracksView()
	}
}

func getCurrentViewSize(v *gocui.View) int {
	if v == gui.tracksView {
		return getTracksViewSize(v)
	} else if v == gui.playlistsView {
		return getPlaylistsViewSize(v)
	}
	return -1
}

func getTracksViewSize(v *gocui.View) int {
	if selectedPlaylist := gui.getSelectedPlaylist(); selectedPlaylist != nil {
		if selectedPlaylist.IsOnDemand() {
			return selectedPlaylist.Tracks()
		}
		return selectedPlaylist.Tracks() - 1
	}
	return -1
}

func getPlaylistsViewSize(v *gocui.View) int {
	subPlaylists := 0
	for _, key := range playlists.Names() {
		playlist := playlists.Get(key)
		if playlist.IsFolder() && playlist.IsFolderOpen() {
			subPlaylists += playlist.Playlists()
		}
	}
	return playlists.Playlists() + subPlaylists - 1
}

func hasMorePages(newOriginY int, cursorY int, maxSize int) bool {
	return newOriginY+cursorY <= maxSize
}

func isNotInLastPage(originY int, pageSizeY int, maxSize int) bool {
	return originY+pageSizeY <= maxSize
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	offset := getOffsetFromTypedNumbers()
	if cx, cy := v.Cursor(); canGoToNewPosition(cy + offset) {
		if err := v.SetCursor(cx, cy+offset); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+offset); err != nil {
				return err
			}
		}
		if v == gui.playlistsView {
			gui.updateTracksView()
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	offset := getOffsetFromTypedNumbers()
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	if err := v.SetCursor(cx, cy-offset); err != nil && oy > 0 {
		if err := v.SetOrigin(ox, oy-offset); err != nil {
			return err
		}
	}
	if v == gui.playlistsView {
		gui.updateTracksView()
	}
	return nil
}

func getOffsetFromTypedNumbers() int {
	if multipleKeysNumber > 1 {
		return multipleKeysNumber
	}
	return 1
}

func canGoToNewPosition(newPosition int) bool {
	currentView := gui.g.CurrentView()
	line, err := currentView.Line(newPosition)
	if err != nil || len(line) == 0 {
		return false
	}
	return true
}

func canGoToAbsoluteNewPosition(v *gocui.View, newPosition int) bool {
	switch v {
	case gui.playlistsView:
		return newPosition <= playlists.Playlists()
	case gui.tracksView:
		if currentPlaylist := gui.getSelectedPlaylist(); currentPlaylist != nil {
			return newPosition <= currentPlaylist.Tracks()
		}
	case gui.queueView:
	}
	return true
}

func goTo(g *gocui.Gui, v *gocui.View, position int) error {
	if canGoToAbsoluteNewPosition(v, position) {
		position--
		ox, _ := v.Origin()
		cx, _ := v.Cursor()
		v.SetCursor(cx, 0)
		v.SetOrigin(ox, 0)
		if err := v.SetCursor(cx, position); err != nil {
			if err := v.SetOrigin(ox, position); err != nil {
				return err
			}
		}
		if v == gui.playlistsView && gui.tracksView != nil {
			gui.updateTracksView()
		}
	}
	return nil
}

func goToFirstLineCommand(g *gocui.Gui, v *gocui.View) error {
	if multipleKeysNumber <= 0 {
		return cursorHome(g, v)
	}

	return goTo(g, v, multipleKeysNumber)
}

func goToLastLineCommand(g *gocui.Gui, v *gocui.View) error {
	if multipleKeysNumber <= 0 {
		return cursorEnd(g, v)
	}

	return goTo(g, v, multipleKeysNumber)
}
