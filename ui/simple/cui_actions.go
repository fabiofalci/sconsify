package simple

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

type keyHandler func(*gocui.Gui, *gocui.View) error

type KeyMapping struct {
	key  interface{}
	mod  gocui.Modifier
	h    keyHandler
	view string
}

type Keyboard struct {
	ConfiguredKeys map[string][]string
	UsedFunctions  map[string]bool

	SequentialKeys map[string]keyHandler

	Keys []*KeyMapping
}

type KeyEntry struct {
	Key     string
	Command string
}

const (
	PauseTrack         = "PauseTrack"
	ShuffleMode        = "ShuffleMode"
	ShuffleAllMode     = "ShuffleAllMode"
	NextTrack          = "NextTrack"
	ReplayTrack        = "ReplayTrack"
	Search             = "Search"
	SearchView         = "SearchView"
	RepeatSearchView   = "RepeatSearchView"
	Quit               = "Quit"
	QueueTrack         = "QueueTrack"
	QueuePlaylist      = "QueuePlaylist"
	RepeatPlayingTrack = "RepeatPlayingTrack"
	RemoveTrack        = "RemoveTrack"
	RemoveAllTracks    = "RemoveAllTracks"
	GoToFirstLine      = "GoToFirstLine"
	GoToLastLine       = "GoToLastLine"
	PlaySelectedTrack  = "PlaySelectedTrack"
	Up                 = "Up"
	Down               = "Down"
	Left               = "Left"
	Right              = "Right"
	OpenCloseFolder    = "OpenCloseFolder"
	ArtistAlbums       = "ArtistAlbums"
	CreatePlaylist     = "CreatePlaylist"
)

var multipleKeysBuffer []rune
var multipleKeysNumber int
var keyboard *Keyboard
var actionBeingExecuted string
var currentView string
var lastQuery string

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
	if !keyboard.UsedFunctions[SearchView] {
		keyboard.addKey("\\", SearchView)
	}
	if !keyboard.UsedFunctions[RepeatSearchView] {
		keyboard.addKey("n", RepeatSearchView)
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
	if !keyboard.UsedFunctions[RemoveTrack] {
		keyboard.addKey("dd", RemoveTrack)
	}
	if !keyboard.UsedFunctions[RemoveAllTracks] {
		keyboard.addKey("D", RemoveAllTracks)
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
	if !keyboard.UsedFunctions[Up] {
		keyboard.addKey("<up>", Up)
		keyboard.addKey("k", Up)
	}
	if !keyboard.UsedFunctions[Down] {
		keyboard.addKey("<down>", Down)
		keyboard.addKey("j", Down)
	}
	if !keyboard.UsedFunctions[Left] {
		keyboard.addKey("<left>", Left)
		keyboard.addKey("h", Left)
	}
	if !keyboard.UsedFunctions[Right] {
		keyboard.addKey("<right>", Right)
		keyboard.addKey("l", Right)
	}
	if !keyboard.UsedFunctions[OpenCloseFolder] {
		keyboard.addKey("<space>", OpenCloseFolder)
	}
	if !keyboard.UsedFunctions[ArtistAlbums] {
		keyboard.addKey("i", ArtistAlbums)
	}
	if !keyboard.UsedFunctions[CreatePlaylist] {
		keyboard.addKey("c", CreatePlaylist)
	}
}

func (keyboard *Keyboard) loadKeyFunctions() {
	if fileLocation := infrastructure.GetKeyFunctionsFileLocation(); fileLocation != "" {
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

func (keyboard *Keyboard) configureKey(handler keyHandler, command string, view string) {
	for key, commands := range keyboard.ConfiguredKeys {
		switch key {
		case "<enter>":
			key = string(gocui.KeyEnter)
		case "<space>":
			key = string(gocui.KeySpace)
		case "<up>":
			key = string(gocui.KeyArrowUp)
		case "<down>":
			key = string(gocui.KeyArrowDown)
		case "<left>":
			key = string(gocui.KeyArrowLeft)
		case "<right>":
			key = string(gocui.KeyArrowRight)
		}
		for _, c := range commands {
			if c == command {
				if view == "" {
					keyboard.SequentialKeys[PlaylistsView+" "+key] = handler
					keyboard.SequentialKeys[QueueView+" "+key] = handler
					keyboard.SequentialKeys[TracksView+" "+key] = handler
				} else {
					keyboard.SequentialKeys[view+" "+key] = handler
				}
			}
		}
	}
}

func keybindings() error {
	keyboard = &Keyboard{
		ConfiguredKeys: make(map[string][]string),
		UsedFunctions:  make(map[string]bool),
		Keys:           make([]*KeyMapping, 0),
		SequentialKeys: make(map[string]keyHandler)}

	multipleKeysBuffer = make([]rune, 0, 0)
	keyboard.loadKeyFunctions()
	keyboard.defaultValues()

	for _, view := range []string{TracksView, PlaylistsView, QueueView} {
		for i := 'a'; i <= 'z'; i++ {
			key := i
			addKeyBinding(&keyboard.Keys, newKeyMapping(i, view, func(g *gocui.Gui, v *gocui.View) error {
				return keyPressed(key, g, v)
			}))
		}

		for i := 'A'; i <= 'Z'; i++ {
			key := i
			addKeyBinding(&keyboard.Keys, newKeyMapping(i, view, func(g *gocui.Gui, v *gocui.View) error {
				return keyPressed(key, g, v)
			}))
		}

		for _, value := range []rune{'>', '<', '/', '\\'} {
			key := value
			addKeyBinding(&keyboard.Keys, newKeyMapping(key, view, func(g *gocui.Gui, v *gocui.View) error {
				return keyPressed(key, g, v)
			}))
		}
		var specialKeys = []gocui.Key{
			gocui.KeySpace, gocui.KeyArrowUp,
			gocui.KeyArrowDown, gocui.KeyArrowLeft,
			gocui.KeyArrowRight, gocui.KeyEnter}

		for _, value := range specialKeys {
			key := value
			addKeyBinding(&keyboard.Keys, newKeyMapping(key, view, func(g *gocui.Gui, v *gocui.View) error {
				return keyPressed(rune(key), g, v)
			}))
		}

		keyboard.configureKey(pauseTrackCommand, PauseTrack, view)
		keyboard.configureKey(setShuffleMode, ShuffleMode, view)
		keyboard.configureKey(setShuffleAllMode, ShuffleAllMode, view)
		keyboard.configureKey(nextTrackCommand, NextTrack, view)
		keyboard.configureKey(replayTrackCommand, ReplayTrack, view)
		keyboard.configureKey(enableSearchInputCommand, Search, view)
		keyboard.configureKey(enableSearchViewInputCommand, SearchView, view)
		keyboard.configureKey(enableRepeatSearchViewInputCommand, RepeatSearchView, view)
		keyboard.configureKey(repeatPlayingTrackCommand, RepeatPlayingTrack, view)
		keyboard.configureKey(quit, Quit, view)
		keyboard.configureKey(goToFirstLineCommand, GoToFirstLine, view)
		keyboard.configureKey(goToLastLineCommand, GoToLastLine, view)
		addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyHome, view, cursorHome))
		addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyEnd, view, cursorEnd))
		addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyPgup, view, cursorPgup))
		addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyPgdn, view, cursorPgdn))
		keyboard.configureKey(cursorUp, Up, view)
		keyboard.configureKey(cursorDown, Down, view)
		keyboard.configureKey(removeTrackCommand, RemoveTrack, view)
		keyboard.configureKey(removeAllTracksCommand, RemoveAllTracks, view)

		// numbers
		for i := 0; i < 10; i++ {
			numberCopy := i
			addKeyBinding(&keyboard.Keys, newKeyMapping(rune(i+48), view,
				func(g *gocui.Gui, v *gocui.View) error {
					return multipleKeysNumberPressed(numberCopy)
				}))
		}
	}

	keyboard.configureKey(queueTrackCommand, QueueTrack, TracksView)
	keyboard.configureKey(queuePlaylistCommand, QueuePlaylist, PlaylistsView)
	keyboard.configureKey(playSelectedTrack, PlaySelectedTrack, TracksView)

	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyEnter, StatusView, executeAction))
	keyboard.configureKey(mainNextViewLeft, Left, TracksView)
	keyboard.configureKey(nextView, Left, QueueView)
	keyboard.configureKey(nextView, Right, PlaylistsView)
	keyboard.configureKey(mainNextViewRight, Right, TracksView)
	keyboard.configureKey(openCloseFolderCommand, OpenCloseFolder, PlaylistsView)
	keyboard.configureKey(artistAlbums, ArtistAlbums, TracksView)
	addKeyBinding(&keyboard.Keys, newKeyMapping(gocui.KeyCtrlC, "", quit))
	keyboard.configureKey(enableCreatePlaylistCommand, CreatePlaylist, QueueView)

	for _, key := range keyboard.Keys {
		// it needs to copy the key because closures copy var references and we don't
		// want to execute always the last action
		keyCopy := key
		if err := gui.g.SetKeybinding(key.view, key.key, key.mod,
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

func newKeyMapping(key interface{}, view string, h  keyHandler) *KeyMapping {
	return newModifiedKeyMapping(gocui.ModNone, key, view, h)
}

func newModifiedKeyMapping(mod gocui.Modifier, key interface{}, view string, h keyHandler) *KeyMapping {
	return &KeyMapping{mod: mod, key: key, h: h, view: view}
}

func keyPressed(key rune, g *gocui.Gui, v *gocui.View) error {
	multipleKeysBuffer = append(multipleKeysBuffer, key)
	var keyCombination string
	if len(multipleKeysBuffer) == 1 {
		keyCombination = string(multipleKeysBuffer[0])
	} else {
		keyCombination = string(multipleKeysBuffer[0]) + string(multipleKeysBuffer[1])
	}

	if handler := keyboard.SequentialKeys[v.Name()+" "+keyCombination]; handler != nil {
		multipleKeysBuffer = make([]rune, 0, 0)
		err := handler(g, v)
		multipleKeysNumber = 0
		return err
	}

	if len(multipleKeysBuffer) >= 2 {
		key1 := multipleKeysBuffer[1]
		multipleKeysBuffer = make([]rune, 0, 0)
		return keyPressed(rune(key1), g, v)
	}
	return nil
}

func multipleKeysNumberPressed(pressedNumber int) error {
	if multipleKeysNumber == 0 {
		multipleKeysNumber = pressedNumber
	} else {
		multipleKeysNumber = multipleKeysNumber*10 + pressedNumber
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
		for i := 1; i <= getOffsetFromTypedNumbers(); i++ {
			track := playlist.Track(trackIndex)
			if queue.Add(track) != nil {
				fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
			}
		}
	}
	return nil
}

func openCloseFolderCommand(g *gocui.Gui, v *gocui.View) error {
	if playlist := gui.getSelectedPlaylist(); playlist != nil {
		if playlist.IsFolder() {
			playlist.InvertOpenClose()
			gui.updatePlaylistsView()
		}
	}
	return nil
}

func artistAlbums(g *gocui.Gui, v *gocui.View) error {
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		track := playlist.Track(trackIndex)
		publisher.GetArtistAlbums(track.Artist)
	}
	return nil
}

func repeatPlayingTrackCommand(g *gocui.Gui, v *gocui.View) error {
	if gui.PlayingTrack != nil {
		for i := 1; i <= getOffsetFromTypedNumbers(); i++ {
			queue.Insert(gui.PlayingTrack)
			gui.updateQueueView()
		}
	}
	return nil
}

func queuePlaylistCommand(g *gocui.Gui, v *gocui.View) error {
	if playlist, _ := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		for i := 1; i <= getOffsetFromTypedNumbers(); i++ {
			for i := 0; i < playlist.Tracks(); i++ {
				track := playlist.Track(i)
				if !addToQueue(track) {
					return nil
				}
			}
		}
	}
	return nil
}

func removeAllTracksCommand(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case PlaylistsView:
	case TracksView:
		if playlist, index := gui.getSelectedPlaylistAndTrack(); index > -1 {
			playlist.RemoveAllTracks()
			gui.updateTracksView()
			return gui.enableSideView()
		}
	case QueueView:
		gui.clearQueueView()
		return gui.enableTracksView()
	}
	return nil
}

func removeTrackCommand(g *gocui.Gui, v *gocui.View) error {
	switch v.Name() {
	case PlaylistsView:
		if playlist := gui.getSelectedPlaylist(); playlist != nil {
			playlists.Remove(playlist.Name())
			gui.updatePlaylistsView()
			gui.updateTracksView()
		}
	case TracksView:
		if playlist, index := gui.getSelectedPlaylistAndTrack(); index > -1 {
			for i := 1; i <= getOffsetFromTypedNumbers(); i++ {
				playlist.RemoveTrack(index)
			}
			gui.updateTracksView()
			goTo(g, v, index+1)
		}
	case QueueView:
		if index := gui.getQueueSelectedTrackIndex(); index > -1 {
			for i := 1; i <= getOffsetFromTypedNumbers(); i++ {
				if queue.Remove(index) != nil {
					continue
				}
			}
			gui.updateQueueView()
		}
	}
	return nil
}

func enableSearchViewInputCommand(g *gocui.Gui, v *gocui.View) error {
	gui.clearStatusView()
	gui.statusView.Editable = true
	currentView = v.Name()
	gui.g.SetCurrentView(StatusView)
	actionBeingExecuted = SearchView
	return nil
}

func enableRepeatSearchViewInputCommand(g *gocui.Gui, v *gocui.View) error {
	executeSearchView(lastQuery)
	return nil
}

func searchViewCommand(g *gocui.Gui, v *gocui.View) error {
	if query := getTypedCommand(); query != "" {
		executeSearchView(query)
		lastQuery = query
	}
	gui.clearStatusView()
	gui.statusView.Editable = false
	gui.updateCurrentStatus()
	return nil
}
func executeSearchView(query string) {
	if currentView == PlaylistsView {
		gui.enableSideView()
		currentPosition := 0
		selectedPlaylist := gui.getSelectedPlaylist()
		if selectedPlaylist != nil {
			currentPosition = selectedPlaylist.Position
		}
		playlist := playlists.Find(query, currentPosition)
		if playlist != nil {
			goTo(gui.g, gui.playlistsView, playlist.Position)
		}
	} else if currentView == QueueView {
		gui.enableQueueView()
	} else if currentView == TracksView {
		gui.enableTracksView()
		playlist, trackIndex := gui.getSelectedPlaylistAndTrack()
		trackIndex++
		if playlist != nil {
			trackIndex := playlist.FindTrackIndex(query, trackIndex)
			if trackIndex > -1 {
				goTo(gui.g, gui.tracksView, trackIndex + 1)
			}
		}
	}
}

func enableSearchInputCommand(g *gocui.Gui, v *gocui.View) error {
	gui.clearStatusView()
	gui.statusView.Editable = true
	gui.g.SetCurrentView(StatusView)
	actionBeingExecuted = Search
	return nil
}

func searchCommand(g *gocui.Gui, v *gocui.View) error {
	if query := getTypedCommand(); query != "" {
		publisher.Search(query)
	}
	gui.enableSideView()
	gui.clearStatusView()
	gui.statusView.Editable = false
	gui.updateCurrentStatus()
	return nil
}

func enableCreatePlaylistCommand(g *gocui.Gui, v *gocui.View) error {
	gui.clearStatusView()
	gui.statusView.Editable = true
	gui.g.SetCurrentView(StatusView)
	actionBeingExecuted = CreatePlaylist
	return nil
}

func createPlaylistCommand(g *gocui.Gui, v *gocui.View) error {
	if playlistName := getTypedCommand(); playlistName != "" {
		gui.createPlaylistFromQueue(playlistName)
	}
	gui.enableSideView()
	gui.clearStatusView()
	gui.statusView.Editable = false
	gui.updateCurrentStatus()
	return nil
}

func getTypedCommand() string {
	typed, _ := gui.statusView.Line(0)
	return strings.Trim(typed, " \x00")
}

func executeAction(g *gocui.Gui, v *gocui.View) error {
	if actionBeingExecuted == Search {
		return searchCommand(g, v)
	} else if actionBeingExecuted == SearchView {
		return searchViewCommand(g, v)
	} else if actionBeingExecuted == CreatePlaylist {
		return createPlaylistCommand(g, v)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	consoleUserInterface.Shutdown()
	// TODO wait for shutdown
	// <-events.ShutdownUpdates()
	return gocui.ErrQuit
}

func (gui *Gui) clearTimeLeftView() {
	gui.timeLeftView.Clear()
	gui.timeLeftView.SetCursor(0, 0)
	gui.timeLeftView.SetOrigin(0, 0)
}
