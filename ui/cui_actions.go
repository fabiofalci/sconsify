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
var multipleKeysNumber int

func keybindings() error {
	keys := make([]*KeyMapping, 0)

	for _, view := range []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE} {
		addKeyBinding(&keys, newKeyMapping('p', view, pauseCurrentSelectedTrack))
		addKeyBinding(&keys, newKeyMapping('r', view, setRandomMode))
		addKeyBinding(&keys, newKeyMapping('R', view, setAllRandomMode))
		addKeyBinding(&keys, newKeyMapping('>', view, nextCommand))
		addKeyBinding(&keys, newKeyMapping('/', view, enableSearchInputCommand))
		addKeyBinding(&keys, newKeyMapping('j', view, cursorDown))
		addKeyBinding(&keys, newKeyMapping('k', view, cursorUp))
		addKeyBinding(&keys, newKeyMapping('q', view, quit))
	}

	allViews := ""
	addKeyBinding(&keys, newKeyMapping(gocui.KeySpace, VIEW_TRACKS, playCurrentSelectedTrack))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, VIEW_TRACKS, playCurrentSelectedTrack))
	addKeyBinding(&keys, newKeyMapping('u', VIEW_TRACKS, queueCommand))
	addKeyBinding(&keys, newKeyMapping('d', VIEW_PLAYLISTS, removeFromPlaylistsCommand))
	addKeyBinding(&keys, newKeyMapping('d', VIEW_QUEUE, removeFromQueueCommand))
	addKeyBinding(&keys, newKeyMapping('D', VIEW_QUEUE, removeAllFromQueueCommand))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, VIEW_STATUS, searchCommand))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyHome, allViews, cursorHome))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnd, allViews, cursorEnd))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyPgup, allViews, cursorPgup))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyPgdn, allViews, cursorPgdn))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowDown, allViews, cursorDown))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowUp, allViews, cursorUp))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowLeft, VIEW_TRACKS, mainNextViewLeft))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowLeft, VIEW_QUEUE, nextView))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowRight, VIEW_PLAYLISTS, nextView))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyArrowRight, VIEW_TRACKS, mainNextViewRight))
	addKeyBinding(&keys, newKeyMapping('h', VIEW_TRACKS, mainNextViewLeft))
	addKeyBinding(&keys, newKeyMapping('h', VIEW_QUEUE, nextView))
	addKeyBinding(&keys, newKeyMapping('l', VIEW_PLAYLISTS, nextView))
	addKeyBinding(&keys, newKeyMapping('l', VIEW_TRACKS, mainNextViewRight))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyCtrlC, allViews, quit))
	addKeyBinding(&keys, newKeyMapping('G', allViews, uppergCommand))

	for _, key := range keys {
		// it needs to copy the key because closures copy var references and we don't
		// want to execute always the last action
		keyCopy := key
		if err := gui.g.SetKeybinding(key.view, key.key, 0,
			func(g *gocui.Gui, v *gocui.View) error {
				err := keyCopy.h(g, v)
				resetMultipleKeys()
				return err
			}); err != nil {
			return err
		}
	}

	// multiple keys
	keys = make([]*KeyMapping, 0)

	addKeyBinding(&keys, newKeyMapping('g', allViews,
		func(g *gocui.Gui, v *gocui.View) error {
			return multipleKeysPressed(g, v, 'g')
		}))

	// numbers
	for i := 0; i < 10; i++ {
		numberCopy := i
		addKeyBinding(&keys, newKeyMapping(rune(i+48), allViews,
			func(g *gocui.Gui, v *gocui.View) error {
				return multipleKeysNumberPressed(numberCopy)
			}))
	}

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

func newKeyMapping(key interface{}, view string, h gocui.KeybindingHandler) *KeyMapping {
	return &KeyMapping{key: key, h: h, view: view}
}

func resetMultipleKeys() {
	multipleKeysBuffer.Reset()
	multipleKeysNumber = 0
}

func multipleKeysNumberPressed(pressedNumber int) error {
	if multipleKeysNumber == 0 {
		multipleKeysNumber = pressedNumber
	} else {
		multipleKeysNumber = multipleKeysNumber*10 + pressedNumber
	}
	return nil
}

func multipleKeysPressed(g *gocui.Gui, v *gocui.View, pressedKey rune) error {
	multipleKeysBuffer.WriteRune(pressedKey)

	switch multipleKeysBuffer.String() {
	case "gg":
		ggCommand(g, v)
		resetMultipleKeys()
	}

	return nil
}

func playCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		track := playlist.Track(trackIndex)
		playlists.SetCurrents(playlist.Name(), trackIndex)
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
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		track := playlist.Track(trackIndex)
		fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
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
	return gocui.Quit
}
