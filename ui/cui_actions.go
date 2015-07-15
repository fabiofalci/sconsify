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

type Keyboard struct {
	ConfiguredKeys map[string][]string
	UsedFunctions  map[string]bool

	Keys []*KeyMapping
	MultipleKeys []*KeyMapping
}

type KeyEntry struct {
	Key     string
	Command string
}

const (
	PauseTrack string = "PauseTrack"
	ShuffleMode string = "ShuffleMode"
	ShuffleAllMode string = "ShuffleAllMode"
	NextTrack string = "NextTrack"
	ReplayTrack string = "ReplayTrack"
	Search string = "Search"
	Quit string = "Quit"
	QueueTrack string = "QueueTrack"
	QueuePlaylist string = "QueuePlaylist"
	RepeatPlayingTrack string = "RepeatPlayingTrack"
	RemoveSearchFromPlaylists string = "RemoveSearchFromPlaylists"
	RemoveTrackFromQueue string = "RemoveTrackFromQueue"
	RemoveAllTracksFromQueue string = "RemoveAllTracksFromQueue"
	GoToFirstLine string = "GoToFirstLine"
	GoToLastLine string = "GoToLastLine"
	PlaySelectedTrack string = "PlaySelectedTrack"
)

var multipleKeysBuffer bytes.Buffer
var multipleKeysNumber int
var multipleKeysHandlers map[string]gocui.KeybindingHandler

func (keyboard *Keyboard) defaultValues() {
	if !keyboard.UsedFunctions[PauseTrack] {
		keyboard.addKey("p", PauseTrack)
	}
	if !keyboard.UsedFunctions[ShuffleMode] {
		keyboard.addKey("s", ShuffleMode)
	}
	if !keyboard.UsedFunctions[ShuffleAllMode] {
		keyboard.addKey("S", ShuffleAllMode)
	}
	if !keyboard.UsedFunctions[NextTrack] {
		keyboard.addKey(">", NextTrack)
	}
	if !keyboard.UsedFunctions[ReplayTrack] {
		keyboard.addKey("<", ReplayTrack)
	}
	if !keyboard.UsedFunctions[Search] {
		keyboard.addKey("/", Search)
	}
	if !keyboard.UsedFunctions[Quit] {
		keyboard.addKey("q", Quit)
	}
	if !keyboard.UsedFunctions[QueueTrack] {
		keyboard.addKey("u", QueueTrack)
	}
	if !keyboard.UsedFunctions[QueuePlaylist] {
		keyboard.addKey("u", QueuePlaylist)
	}
	if !keyboard.UsedFunctions[RepeatPlayingTrack] {
		keyboard.addKey("r", RepeatPlayingTrack)
	}
	if !keyboard.UsedFunctions[RemoveSearchFromPlaylists] {
		keyboard.addKey("d", RemoveSearchFromPlaylists)
	}
	if !keyboard.UsedFunctions[RemoveTrackFromQueue] {
		keyboard.addKey("d", RemoveTrackFromQueue)
	}
	if !keyboard.UsedFunctions[RemoveAllTracksFromQueue] {
		keyboard.addKey("D", RemoveAllTracksFromQueue)
	}
	if !keyboard.UsedFunctions[GoToFirstLine] {
		keyboard.addKey("gg", GoToFirstLine)
	}
	if !keyboard.UsedFunctions[GoToLastLine] {
		keyboard.addKey("G", GoToLastLine)
	}
	if !keyboard.UsedFunctions[PlaySelectedTrack] {
		keyboard.addKey("<space>", PlaySelectedTrack)
		keyboard.addKey("<enter>", PlaySelectedTrack)
	}
}

func (keyboard *Keyboard) loadKeyFunctions() {
	if fileLocation := sconsify.GetKeyFunctionsFileLocation(); fileLocation != "" {
		if b, err := ioutil.ReadFile(fileLocation); err == nil {
			fileContent := make([]KeyEntry, 0)
			if err := json.Unmarshal(b, &fileContent); err == nil {
				for _, keyEntry := range fileContent {
					keyboard.addKey(keyEntry.Key, keyEntry.Command)
				}
			}
		}
	}
}

func (keyboard *Keyboard) addKey(key string, command string) {
	if keyboard.ConfiguredKeys[key] == nil {
		keyboard.ConfiguredKeys[key] = make([]string, 0)
	}
	keyboard.ConfiguredKeys[key] = append(keyboard.ConfiguredKeys[key], command)
	keyboard.UsedFunctions[command] = true
}

func (keyboard *Keyboard) configureKey(handler gocui.KeybindingHandler, command string, view string) {
	for key, commands := range keyboard.ConfiguredKeys {
		for _, c := range commands {
			if c == command {
				keyMapping, isMultiple := createKeyMapping(handler, key, view)
				keyboard.addToKeys(isMultiple, keyMapping)
			}
		}
	}
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

func createKeyMapping(handler gocui.KeybindingHandler, key string, view string) (*KeyMapping, bool) {
	switch key {
	case "<enter>":
		return newKeyMapping(gocui.KeyEnter, view, handler), false
	case "<space>":
		return newKeyMapping(gocui.KeySpace, view, handler), false
	}
	if isMultipleKey(key) {
		keyRune := getAsRuneArray(key)
		multipleKeysHandlers[key] = handler
		return newKeyMapping(keyRune[0], view,
			func(g *gocui.Gui, v *gocui.View) error {
				return multipleKeysPressed(g, v, keyRune[0])
			}), true
	}
	return newKeyMapping(getFirstRune(key), view, handler), false
}

func (keyboard *Keyboard) addToKeys(isMultiple bool, keyMapping *KeyMapping) {
	if isMultiple {
		addKeyBinding(&keyboard.MultipleKeys, keyMapping)
	} else {
		addKeyBinding(&keyboard.Keys, keyMapping)
	}
}

func keybindings() error {
	keyboard := &Keyboard{
		ConfiguredKeys: make(map[string][]string),
		UsedFunctions: make(map[string]bool),
		Keys: make([]*KeyMapping, 0),
		MultipleKeys: make([]*KeyMapping, 0)}

	keyboard.loadKeyFunctions()
	keyboard.defaultValues()

	multipleKeysHandlers = make(map[string]gocui.KeybindingHandler)

	for _, view := range []string{VIEW_TRACKS, VIEW_PLAYLISTS, VIEW_QUEUE} {
		keyboard.configureKey(pauseTrackCommand, PauseTrack, view)
		keyboard.configureKey(setShuffleMode, ShuffleMode, view)
		keyboard.configureKey(setShuffleAllMode, ShuffleAllMode, view)
		keyboard.configureKey(nextTrackCommand, NextTrack, view)
		keyboard.configureKey(replayTrackCommand, ReplayTrack, view)
		keyboard.configureKey(enableSearchInputCommand, Search, view)
		keyboard.configureKey(repeatPlayingTrackCommand, RepeatPlayingTrack, view)
		keyboard.configureKey(quit, Quit, view)
		addKeyBinding(&keyboard.Keys, newKeyMapping('j', view, cursorDown))
		addKeyBinding(&keyboard.Keys, newKeyMapping('k', view, cursorUp))
	}

	allViews := ""
	keyboard.configureKey(queueTrackCommand, QueueTrack, VIEW_TRACKS)
	keyboard.configureKey(queuePlaylistCommand, QueuePlaylist, VIEW_PLAYLISTS)
	keyboard.configureKey(removeSearchPlaylistsCommand, RemoveSearchFromPlaylists, VIEW_PLAYLISTS)
	keyboard.configureKey(removeTrackFromQueueCommand, RemoveTrackFromQueue, VIEW_QUEUE)
	keyboard.configureKey(removeAllTracksFromQueueCommand, RemoveAllTracksFromQueue, VIEW_QUEUE)
	keyboard.configureKey(goToFirstLineCommand, GoToFirstLine, allViews)
	keyboard.configureKey(goToLastLineCommand, GoToLastLine, allViews)
	keyboard.configureKey(playSelectedTrack, PlaySelectedTrack, VIEW_TRACKS)

	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyEnter, VIEW_STATUS, searchCommand))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyHome, allViews, cursorHome))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyEnd, allViews, cursorEnd))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyPgup, allViews, cursorPgup))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyPgdn, allViews, cursorPgdn))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowDown, allViews, cursorDown))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowUp, allViews, cursorUp))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowLeft, VIEW_TRACKS, mainNextViewLeft))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowLeft, VIEW_QUEUE, nextView))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowRight, VIEW_PLAYLISTS, nextView))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyArrowRight, VIEW_TRACKS, mainNextViewRight))
	addKeyBinding(&keyboard.Keys, newKeyMapping('h', VIEW_TRACKS, mainNextViewLeft))
	addKeyBinding(&keyboard.Keys, newKeyMapping('h', VIEW_QUEUE, nextView))
	addKeyBinding(&keyboard.Keys, newKeyMapping('l', VIEW_PLAYLISTS, nextView))
	addKeyBinding(&keyboard.Keys, newKeyMapping('l', VIEW_TRACKS, mainNextViewRight))
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyCtrlC, allViews, quit))

	// numbers
	for i := 0; i < 10; i++ {
		numberCopy := i
		addKeyBinding(&keyboard.MultipleKeys, newKeyMapping(rune(i+48), allViews,
			func(g *gocui.Gui, v *gocui.View) error {
				return multipleKeysNumberPressed(numberCopy)
			}))
	}

	for _, key := range keyboard.Keys {
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

	for _, key := range keyboard.MultipleKeys {
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
	}

	return nil
}

func playSelectedTrack(g *gocui.Gui, v *gocui.View) error {
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

func repeatPlayingTrackCommand(g *gocui.Gui, v *gocui.View) error {
	if gui.PlayingTrack != nil {
		queue.Insert(gui.PlayingTrack)
		gui.updateQueueView()
	}
	return nil
}

func queuePlaylistCommand(g *gocui.Gui, v *gocui.View) error {
	if playlist, _ := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i)
			fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
			if queue.Add(track) == nil {
				return nil
			}
		}
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

func removeSearchPlaylistsCommand(g *gocui.Gui, v *gocui.View) error {
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
