package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

type KeyMapping struct {
	key  interface{}
	h    gocui.KeybindingHandler
	view string
}

type KeyFunctions struct {
	PauseTrack               string
	ShuffleMode              string
	ShuffleAllMode           string
	NextTrack                string
	ReplayTrack              string
	Search                   string
	Quit                     string
	QueueTrack               string
	RemoveTrackFromPlaylist  string
	RemoveTrackFromQueue     string
	RemoveAllTracksFromQueue string
}

var multipleKeysBuffer bytes.Buffer
var multipleKeysNumber int

func (k *KeyFunctions) defaultValues() {
	if k.PauseTrack == "" {
		k.PauseTrack = "p"
	} 
	if k.ShuffleMode == "" {
		k.ShuffleMode = "s"
	} 
	if k.ShuffleAllMode == "" {
		k.ShuffleAllMode = "S"
	} 
	if k.NextTrack == "" {
		k.NextTrack = ">"
	} 
	if k.ReplayTrack == "" {
		k.ReplayTrack = "<"
	} 
	if k.Search == "" {
		k.Search = "/"
	} 
	if k.Quit == "" {
		k.Quit = "q"
	} 
	if k.QueueTrack == "" {
		k.QueueTrack = "u"
	} 
	if k.RemoveTrackFromPlaylist == "" {
		k.RemoveTrackFromPlaylist = "d"
	} 
	if k.RemoveTrackFromQueue == "" {
		k.RemoveTrackFromQueue = "d"
	} 
	if k.RemoveAllTracksFromQueue == "" {
		k.RemoveAllTracksFromQueue = "D"
	} 
}

func loadKeyFunctions() *KeyFunctions {
	if fileLocation := sconsify.GetKeyFunctionsFileLocation(); fileLocation != "" {
		if b, err := ioutil.ReadFile(fileLocation); err == nil {
			var keyFunctions KeyFunctions
			if err := json.Unmarshal(b, &keyFunctions); err == nil {
				return &keyFunctions
			}
		}
	}
	return &KeyFunctions{}
}

func getFirstRune(value string) rune {
	runes := []rune(value)
	return runes[0]
}

func keybindings() error {
	keyFunctions := loadKeyFunctions()
	keyFunctions.defaultValues()

	keys := make([]*KeyMapping, 0)

	for _, view := range []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE} {
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.PauseTrack), view, pauseTrackCommand))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.ShuffleMode), view, setShuffleMode))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.ShuffleAllMode), view, setShuffleAllMode))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.NextTrack), view, nextTrackCommand))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.ReplayTrack), view, replayTrackCommand))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.Search), view, enableSearchInputCommand))
		addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.Quit), view, quit))
		addKeyBinding(&keys, newKeyMapping('j' , view, cursorDown))
		addKeyBinding(&keys, newKeyMapping('k', view, cursorUp))
	}

	allViews := ""
	addKeyBinding(&keys, newKeyMapping(gocui.KeySpace, VIEW_TRACKS, playCurrentSelectedTrack))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, VIEW_TRACKS, playCurrentSelectedTrack))
	addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.QueueTrack), VIEW_TRACKS, queueTrackCommand))
	addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.RemoveTrackFromPlaylist), VIEW_PLAYLISTS, removeTrackFromPlaylistsCommand))
	addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.RemoveTrackFromQueue), VIEW_QUEUE, removeTrackFromQueueCommand))
	addKeyBinding(&keys, newKeyMapping(getFirstRune(keyFunctions.RemoveAllTracksFromQueue), VIEW_QUEUE, removeAllTracksFromQueueCommand))
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
	player.Play()
	return nil
}

func pauseTrackCommand(g *gocui.Gui, v *gocui.View) error {
	player.Pause()
	return nil
}

func setShuffleMode(g *gocui.Gui, v *gocui.View) error {
	playlists.InvertMode(sconsify.ShuffleMode)
	gui.updateCurrentStatus()
	return nil
}

func setShuffleAllMode(g *gocui.Gui, v *gocui.View) error {
	playlists.InvertMode(sconsify.ShuffleAllMode)
	gui.updateCurrentStatus()
	return nil
}

func nextTrackCommand(g *gocui.Gui, v *gocui.View) error {
	gui.playNext()
	return nil
}

func replayTrackCommand(g *gocui.Gui, v *gocui.View) error {
	gui.replay()
	return nil
}

func queueTrackCommand(g *gocui.Gui, v *gocui.View) error {
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		track := playlist.Track(trackIndex)
		fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
		queue.Add(track)
	}
	return nil
}

func removeAllTracksFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	queue.RemoveAll()
	gui.updateQueueView()
	return gui.enableTracksView()
}

func removeTrackFromQueueCommand(g *gocui.Gui, v *gocui.View) error {
	if index := gui.getQueueSelectedTrackIndex(); index > -1 {
		queue.Remove(index)
		gui.updateQueueView()
	}
	return nil
}

func removeTrackFromPlaylistsCommand(g *gocui.Gui, v *gocui.View) error {
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
	query, _ := gui.statusView.Line(-1)
	query = strings.Trim(query, " ")
	if query != "" {
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
