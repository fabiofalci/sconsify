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
var multipleKeysHandlers map[string]gocui.KeybindingHandler

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
	return getAsRuneArray(value)[0]
}

func getAsRuneArray(value string) []rune {
	return []rune(value)
}

func isMultipleKey(value string) bool {
	return len(getAsRuneArray(value)) > 1
}

func createKeyMapping(handler gocui.KeybindingHandler, command string, view string) (*KeyMapping, bool) {
	if isMultipleKey(command) {
		keyRune := getAsRuneArray(command)
		multipleKeysHandlers[command] = handler
		return newKeyMapping(keyRune[0], view,
			func(g *gocui.Gui, v *gocui.View) error {
				return multipleKeysPressed(g, v, keyRune[0])
			}), true
	}
	return newKeyMapping(getFirstRune(command), view, handler), false
}

func addToKeys(isMultiple bool, keyMapping *KeyMapping, keys *[]*KeyMapping, multipleKeys *[]*KeyMapping) {
	if isMultiple {
		addKeyBinding(multipleKeys, keyMapping)
	} else {
		addKeyBinding(keys, keyMapping)
	}
}

func keybindings() error {
	keyFunctions := loadKeyFunctions()
	keyFunctions.defaultValues()

	keys := make([]*KeyMapping, 0)
	multipleKeys := make([]*KeyMapping, 0)
	multipleKeysHandlers = make(map[string]gocui.KeybindingHandler)

	var keyMapping *KeyMapping
	var isMultiple bool

	for _, view := range []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE} {
		keyMapping, isMultiple = createKeyMapping(pauseTrackCommand, keyFunctions.PauseTrack, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(setShuffleMode, keyFunctions.ShuffleMode, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(setShuffleAllMode, keyFunctions.ShuffleAllMode, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(nextTrackCommand, keyFunctions.NextTrack, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(replayTrackCommand, keyFunctions.ReplayTrack, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(enableSearchInputCommand, keyFunctions.Search, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		keyMapping, isMultiple = createKeyMapping(quit, keyFunctions.Quit, view)
		addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

		addKeyBinding(&keys, newKeyMapping('j', view, cursorDown))
		addKeyBinding(&keys, newKeyMapping('k', view, cursorUp))
	}

	allViews := ""
	keyMapping, isMultiple = createKeyMapping(queueTrackCommand, keyFunctions.QueueTrack, allViews)
	addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

	keyMapping, isMultiple = createKeyMapping(removeTrackFromPlaylistsCommand, keyFunctions.RemoveTrackFromPlaylist, allViews)
	addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

	keyMapping, isMultiple = createKeyMapping(removeTrackFromQueueCommand, keyFunctions.RemoveTrackFromQueue, allViews)
	addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

	keyMapping, isMultiple = createKeyMapping(removeAllTracksFromQueueCommand, keyFunctions.RemoveAllTracksFromQueue, allViews)
	addToKeys(isMultiple, keyMapping, &keys, &multipleKeys)

	addKeyBinding(&keys, newKeyMapping(gocui.KeySpace, VIEW_TRACKS, playCurrentSelectedTrack))
	addKeyBinding(&keys, newKeyMapping(gocui.KeyEnter, VIEW_TRACKS, playCurrentSelectedTrack))
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

	addKeyBinding(&multipleKeys, newKeyMapping('g', allViews,
		func(g *gocui.Gui, v *gocui.View) error {
			return multipleKeysPressed(g, v, 'g')
		}))

	// numbers
	for i := 0; i < 10; i++ {
		numberCopy := i
		addKeyBinding(&multipleKeys, newKeyMapping(rune(i+48), allViews,
			func(g *gocui.Gui, v *gocui.View) error {
				return multipleKeysNumberPressed(numberCopy)
			}))
	}

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

	for _, key := range multipleKeys {
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

	handler := multipleKeysHandlers[multipleKeysBuffer.String()]
	if handler != nil {
		handler(g, v)
		resetMultipleKeys()
	} else {
		switch multipleKeysBuffer.String() {
		case "gg":
			ggCommand(g, v)
			resetMultipleKeys()
		}
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
