package ui

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/events"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

var (
	gui       *Gui
	queue     *Queue
	state     *UiState
	playlists *sconsify.Playlists
)

type Gui struct {
	g             *gocui.Gui
	playlistsView *gocui.View
	tracksView    *gocui.View
	statusView    *gocui.View
	queueView     *gocui.View
	events        *events.Events
	currentTrack  *sconsify.Track
}

func StartConsoleUserInterface(events *events.Events) {
	select {
	case p := <-events.WaitForPlaylists():
		playlists = &p
		if playlists == nil {
			return
		}
	case <-events.WaitForShutdown():
		return
	}

	gui = &Gui{events: events}

	queue = InitQueue()
	state = InitState()

	go func() {
		for {
			select {
			case track := <-gui.events.WaitForTrackPaused():
				gui.trackPaused(track)
			case track := <-gui.events.WaitForTrackPlaying():
				gui.trackPlaying(track)
			case track := <-gui.events.WaitForTrackNotAvailable():
				gui.trackNotAvailable(track)
			case <-gui.events.WaitForPlayTokenLost():
				gui.updateStatus("Play token lost", false)
			case <-gui.events.WaitForNextPlay():
				gui.playNext()
			case newPlaylist := <-events.WaitForPlaylists():
				gui.newPlaylist(&newPlaylist)
			}
		}
	}()

	gui.g = gocui.NewGui()
	if err := gui.g.Init(); err != nil {
		log.Panicln(err)
	}
	defer gui.g.Close()

	gui.g.SetLayout(layout)
	if err := keybindings(); err != nil {
		log.Panicln(err)
	}
	gui.g.SelBgColor = gocui.ColorGreen
	gui.g.SelFgColor = gocui.ColorBlack
	gui.g.ShowCursor = true

	err := gui.g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}

func (gui *Gui) updateStatus(message string, temporary bool) {
	gui.statusView.Clear()
	gui.statusView.SetCursor(0, 0)
	gui.statusView.SetOrigin(0, 0)

	if !temporary {
		state.currentMessage = message
	} else {
		go func() {
			time.Sleep(4 * time.Second)
			gui.updateStatus(state.currentMessage, false)
		}()
	}
	fmt.Fprintf(gui.statusView, state.getModeAsString()+"%v", message)

	// otherwise the update will appear only in the next keyboard move
	gui.g.Flush()
}

func (gui *Gui) trackNotAvailable(track *sconsify.Track) {
	gui.updateStatus("Not available: "+track.GetTitle(), true)
}

func (gui *Gui) trackPlaying(track *sconsify.Track) {
	gui.updateStatus("Playing: "+track.GetFullTitle(), false)
}

func (gui *Gui) trackPaused(track *sconsify.Track) {
	gui.updateStatus("Paused: "+track.GetFullTitle(), false)
}

func (gui *Gui) getSelectedPlaylist() (string, error) {
	return gui.getSelected(gui.playlistsView)
}

func (gui *Gui) getSelectedTrack() (string, error) {
	return gui.getSelected(gui.tracksView)
}

func (gui *Gui) getQeueuSelectedTrackIndex() int {
	_, cy := gui.queueView.Cursor()
	return cy
}

func (gui *Gui) getSelected(v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, nil
}

func (gui *Gui) playNext() error {
	if !queue.isEmpty() {
		gui.playNextFromQueue()
	} else if state.hasPlaylistSelected() {
		gui.playNextFromPlaylist()
	}
	return nil
}

func (gui *Gui) playNextFromPlaylist() {
	playlist := playlists.Get(state.currentPlaylist)
	if state.isAllRandomMode() {
		state.currentPlaylist, state.currentIndexTrack = playlists.GetRandomNextPlaylistAndTrack()
		playlist = playlists.Get(state.currentPlaylist)
	} else if state.isRandomMode() {
		state.currentIndexTrack = playlist.GetRandomNextTrack()
	} else {
		state.currentIndexTrack = playlist.GetNextTrack(state.currentIndexTrack)
	}
	track := playlist.Track(state.currentIndexTrack)

	gui.play(track)
}

func (gui *Gui) playNextFromQueue() {
	gui.play(queue.Pop())
	gui.updateQueueView()
}

func (gui *Gui) play(track *sconsify.Track) {
	gui.currentTrack = track
	gui.events.Play(gui.currentTrack)
}

func getCurrentSelectedTrack() *sconsify.Track {
	var errPlaylist error
	state.currentPlaylist, errPlaylist = gui.getSelectedPlaylist()
	currentTrack, errTrack := gui.getSelectedTrack()
	if errPlaylist == nil && errTrack == nil && playlists != nil {
		playlist := playlists.Get(state.currentPlaylist)

		if playlist != nil {
			currentTrack = currentTrack[0:strings.Index(currentTrack, ".")]
			converted, _ := strconv.Atoi(currentTrack)
			state.currentIndexTrack = converted - 1
			track := playlist.Track(state.currentIndexTrack)
			return track
		}
	}
	return nil
}

func keybindings() error {
	views := []string{"main", "side", "queue"}
	for _, view := range views {
		if err := gui.g.SetKeybinding(view, 'p', 0, pauseCurrentSelectedTrack); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'r', 0, setRandomMode); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'R', 0, setAllRandomMode); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, '>', 0, nextCommand); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, '/', 0, enableSearchInputCommand); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'j', 0, cursorDown); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'k', 0, cursorUp); err != nil {
			return err
		}
		if err := gui.g.SetKeybinding(view, 'q', 0, quit); err != nil {
			return err
		}
	}

	if err := gui.g.SetKeybinding("main", gocui.KeySpace, 0, playCurrentSelectedTrack); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", 'u', 0, queueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", 'd', 0, removeFromQueueCommand); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("status", gocui.KeyEnter, 0, searchCommand); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyHome, 0, cursorHome); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyEnd, 0, cursorEnd); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyPgup, 0, cursorPgup); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyPgdn, 0, cursorPgdn); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyArrowDown, 0, cursorDown); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("", gocui.KeyArrowUp, 0, cursorUp); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", gocui.KeyArrowLeft, 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", gocui.KeyArrowLeft, 0, nextView); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("side", gocui.KeyArrowRight, 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("main", gocui.KeyArrowRight, 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("main", 'h', 0, mainNextViewLeft); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("queue", 'h', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("side", 'l', 0, nextView); err != nil {
		return err
	}
	if err := gui.g.SetKeybinding("main", 'l', 0, mainNextViewRight); err != nil {
		return err
	}

	if err := gui.g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}

	return nil
}

func (gui *Gui) newPlaylist(newPlaylist *sconsify.Playlists) {
	playlists.Merge(newPlaylist)
	gui.updatePlaylistsView()
	gui.updateTracksView()
	gui.g.Flush()
}

func (gui *Gui) updateTracksView() {
	gui.tracksView.Clear()
	gui.tracksView.SetCursor(0, 0)
	gui.tracksView.SetOrigin(0, 0)
	currentPlaylist, err := gui.getSelectedPlaylist()
	if err == nil && playlists != nil {
		playlist := playlists.Get(currentPlaylist)

		if playlist != nil {
			for i := 0; i < playlist.Tracks(); i++ {
				track := playlist.Track(i)
				fmt.Fprintf(gui.tracksView, "%v. %v", (i + 1), track.GetTitle())
			}
		}
	}
}

func (gui *Gui) updatePlaylistsView() {
	gui.playlistsView.Clear()
	keys := playlists.GetNames()
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintln(gui.playlistsView, key)
	}
}

func (gui *Gui) updateQueueView() {
	gui.queueView.Clear()
	if !queue.isEmpty() {
		for _, track := range queue.Contents() {
			fmt.Fprintf(gui.queueView, "%v", track.GetTitle())
		}
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 25, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.playlistsView = v
		gui.playlistsView.Highlight = true

		gui.updatePlaylistsView()

		if err := g.SetCurrentView("side"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("main", 25, -1, maxX-50, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.tracksView = v

		gui.updateTracksView()
	}
	if v, err := g.SetView("queue", maxX-50, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.queueView = v
	}
	if v, err := g.SetView("status", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.statusView = v
	}
	return nil
}
