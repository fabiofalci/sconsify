package ui

import (
	"bytes"
	"fmt"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

type KeyMapping struct {
	key  interface{}
	h    gocui.KeybindingHandler
	view string
}

var multipleKeysBuffer bytes.Buffer

func keybindings() error {
	keys := make([]*KeyMapping, 0)

	for _, view := range []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE} {
		addKeyBinding(&keys, newKeyMapping('p', pauseCurrentSelectedTrack, view))
		addKeyBinding(&keys, newKeyMapping('r', setRandomMode, view))
		addKeyBinding(&keys, newKeyMapping('R', setAllRandomMode, view))
		addKeyBinding(&keys, newKeyMapping('>', nextCommand, view))
		addKeyBinding(&keys, newKeyMapping('/', enableSearchInputCommand, view))
		addKeyBinding(&keys, newKeyMapping('j', cursorDown, view))
		addKeyBinding(&keys, newKeyMapping('k', cursorUp, view))
		addKeyBinding(&keys, newKeyMapping('q', quit, view))
	}

	allViews := ""
	addKeyBinding(&keys, newKeyMapping(gocui.KeySpace, playCurrentSelectedTrack, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, playCurrentSelectedTrack, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping('u', queueCommand, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping('d', removeFromPlaylistsCommand, VIEW_PLAYLISTS))
	addKeyBinding(&keys, newKeyMapping('d', removeFromQueueCommand, VIEW_QUEUE))
	addKeyBinding(&keys, newKeyMapping('D', removeAllFromQueueCommand, VIEW_QUEUE))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, searchCommand, VIEW_STATUS))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyHome, cursorHome, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnd, cursorEnd, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyPgup, cursorPgup, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyPgdn, cursorPgdn, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowDown, cursorDown, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowUp, cursorUp, allViews))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowLeft, mainNextViewLeft, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowLeft, nextView, VIEW_QUEUE))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowRight, nextView, VIEW_PLAYLISTS))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowRight, mainNextViewRight, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping('h', mainNextViewLeft, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping('h', nextView, VIEW_QUEUE))
	addKeyBinding(&keys, newKeyMapping('l', nextView, VIEW_PLAYLISTS))
	addKeyBinding(&keys, newKeyMapping('l', mainNextViewRight, VIEW_TRACKS))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyCtrlC, quit, allViews))
	addKeyBinding(&keys, newKeyMapping('G', cursorEnd, allViews))

	for _, key := range keys {
		// it needs to copy the key because closures copy var references and we don't
		// want to execute always the last action
		keyCopy := key
		if err := gui.g.SetKeybinding(key.view, key.key, 0,
			func(g *gocui.Gui, v *gocui.View) error {
				resetMultipleKeys()
				return keyCopy.h(g, v)
			}); err != nil {
			return err
		}
	}

	// multiple keys
	keys = make([]*KeyMapping, 0)

	addKeyBinding(&keys, newKeyMapping('g',
		func(g *gocui.Gui, v *gocui.View) error {
			return multipleKeys(g, v, 'g')
		}, allViews))

	for _, key := range keys {
		keyCopy := key
		if err := gui.g.SetKeybinding(key.view, key.key, 0,
			func(g *gocui.Gui, v *gocui.View) error {
				return keyCopy.h(g, v)
			}); err != nil {
			return err
		}
	}
	return nil
}

func addKeyBinding(keys *[]*KeyMapping, key *KeyMapping) {
	*keys = append(*keys, key)
}

func newKeyMapping(key interface{}, h gocui.KeybindingHandler, view string) *KeyMapping {
	return &KeyMapping{key: key, h: h, view: view}
}

func resetMultipleKeys() {
	multipleKeysBuffer.Reset()
}

func multipleKeys(g *gocui.Gui, v *gocui.View, pressedKey rune) error {
	multipleKeysBuffer.WriteRune(pressedKey)

	switch multipleKeysBuffer.String() {
	case "gg":
		cursorHome(g, v)
		resetMultipleKeys()
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
