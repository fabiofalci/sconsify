package ui

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/spotify"
	"github.com/jroimartin/gocui"
	sp "github.com/op/go-libspotify/spotify"
)

var g *gocui.Gui
var playlistsView *gocui.View
var tracksView *gocui.View
var statusView *gocui.View
var queueView *gocui.View

var queue *Queue
var state *UiState
var playlists map[string]*sp.Playlist

var playEvents *events.Events

func Start(allEvents *events.Events) {
	playEvents = allEvents

	playlists = <-playEvents.Playlists
	if playlists == nil {
		return
	}

	queue = InitQueue()
	state = InitState()

	go func() {
		for {
			select {
			case message := <-playEvents.Status:
				updateStatus(message)
			case <-playEvents.NextPlay:
				playNext()
			}
		}
	}()

	g = gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.ShowCursor = true

	err := g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}

func updateStatus(message string) {
	statusView.Clear()
	statusView.SetCursor(0, 0)
	statusView.SetOrigin(0, 0)

	state.currentMessage = message
	fmt.Fprintf(statusView, state.getModeAsString()+"%v", state.currentMessage)

	// otherwise the update will appear only in the next keyboard move
	g.Flush()
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	if currentView == nil || currentView.Name() == "side" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("side")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
		if v == playlistsView {
			updateTracksView(g)
		}
	}
	return nil
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
		if v == playlistsView {
			updateTracksView(g)
		}
	}
	return nil
}

func getSelectedPlaylist(g *gocui.Gui) (string, error) {
	return getSelected(g, playlistsView)
}

func getSelectedTrack(g *gocui.Gui) (string, error) {
	return getSelected(g, tracksView)
}

func getSelected(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	spotify.ShutdownSpotify()
	return gocui.ErrorQuit
}

func playNext() error {
	if !queue.isEmpty() {
		playEvents.ToPlay <- queue.Pop()
		updateQueueView()
	} else if state.currentPlaylist != "" {
		playlist := playlists[state.currentPlaylist]
		if !state.randomMode {
			state.currentIndexTrack = getNextTrack(playlist)
		} else {
			state.currentIndexTrack = getRandomNextTrack(playlist)
		}
		playlistTrack := playlist.Track(state.currentIndexTrack)
		track := playlistTrack.Track()
		track.Wait()

		playEvents.ToPlay <- track
	}
	return nil
}

func getNextTrack(playlist *sp.Playlist) int {
	if state.currentIndexTrack >= playlist.Tracks()-1 {
		return 0
	}
	return state.currentIndexTrack + 1
}

func getRandomNextTrack(playlist *sp.Playlist) int {
	return rand.Intn(playlist.Tracks())
}

func playCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	track := getCurrentSelectedTrack()
	if track != nil {
		playEvents.ToPlay <- track
	}
	return nil
}

func pauseCurrentSelectedTrack(g *gocui.Gui, v *gocui.View) error {
	playEvents.Pause <- true
	return nil
}

func setRandomMode(g *gocui.Gui, v *gocui.View) error {
	state.invertMode()
	updateStatus(state.currentMessage)
	return nil
}

func nextCommand(g *gocui.Gui, v *gocui.View) error {
	playNext()
	return nil
}

func queueCommand(g *gocui.Gui, v *gocui.View) error {
	track := getCurrentSelectedTrack()
	if track != nil {
		fmt.Fprintf(queueView, "%v - %v", track.Artist(0).Name(), track.Name())
		queue.Add(track)
	}
	return nil
}

func getCurrentSelectedTrack() *sp.Track {
	var errPlaylist error
	state.currentPlaylist, errPlaylist = getSelectedPlaylist(g)
	currentTrack, errTrack := getSelectedTrack(g)
	if errPlaylist == nil && errTrack == nil && playlists != nil {
		playlist := playlists[state.currentPlaylist]

		if playlist != nil {
			playlist.Wait()
			currentTrack = currentTrack[0:strings.Index(currentTrack, ".")]
			converted, _ := strconv.Atoi(currentTrack)
			state.currentIndexTrack = converted - 1
			playlistTrack := playlist.Track(state.currentIndexTrack)
			track := playlistTrack.Track()
			track.Wait()
			return track
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("main", gocui.KeySpace, 0, playCurrentSelectedTrack); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', 0, pauseCurrentSelectedTrack); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', 0, setRandomMode); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '>', 0, nextCommand); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'u', 0, queueCommand); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, 0, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, 0, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, 0, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowRight, 0, nextView); err != nil {
		return err
	}

	// vi navigation
	if err := g.SetKeybinding("", 'j', 0, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', 0, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'h', 0, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'l', 0, nextView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', 0, quit); err != nil {
		return err
	}

	return nil
}

func updateTracksView(g *gocui.Gui) {
	tracksView.Clear()
	tracksView.SetCursor(0, 0)
	tracksView.SetOrigin(0, 0)
	currentPlaylist, err := getSelectedPlaylist(g)
	if err == nil && playlists != nil {
		playlist := playlists[currentPlaylist]

		if playlist != nil {
			playlist.Wait()
			for i := 0; i < playlist.Tracks(); i++ {
				playlistTrack := playlist.Track(i)
				track := playlistTrack.Track()
				track.Wait()
				fmt.Fprintf(tracksView, "%v. %v - %v", (i + 1), track.Artist(0).Name(), track.Name())
			}
		}
	}
}

func updatePlaylistsView(g *gocui.Gui) {
	playlistsView.Clear()
	if playlists != nil {
		keys := make([]string, len(playlists))
		i := 0
		for k, _ := range playlists {
			keys[i] = k
			i++
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintln(playlistsView, key)
		}
	}
}

func updateQueueView() {
	queueView.Clear()
	if !queue.isEmpty() {
		for _, track := range queue.Contents() {
			fmt.Fprintf(queueView, "%v - %v", track.Artist(0).Name(), track.Name())
		}
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 25, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		playlistsView = v
		playlistsView.Highlight = true

		updatePlaylistsView(g)
	}
	if v, err := g.SetView("main", 25, -1, maxX-50, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		tracksView = v
		tracksView.Highlight = true

		updateTracksView(g)

		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("queue", maxX-50, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		queueView = v
	}
	if v, err := g.SetView("status", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		statusView = v
	}
	return nil
}
