package ui

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

var (
	gui                  *Gui
	events               *sconsify.Events
	queue                *Queue
	playlists            *sconsify.Playlists
	consoleUserInterface sconsify.UserInterface
	player               Player
)

const (
	VIEW_PLAYLISTS = "playlists"
	VIEW_TRACKS    = "tracks"
	VIEW_QUEUE     = "queue"
	VIEW_STATUS    = "status"
)

type ConsoleUserInterface struct{}

type Gui struct {
	g              *gocui.Gui
	playlistsView  *gocui.View
	tracksView     *gocui.View
	statusView     *gocui.View
	queueView      *gocui.View
	currentMessage string
	initialised    bool
	PlayingTrack   *sconsify.Track
}

func InitialiseConsoleUserInterface(ev *sconsify.Events) sconsify.UserInterface {
	events = ev
	gui = &Gui{}
	consoleUserInterface = &ConsoleUserInterface{}
	queue = InitQueue()
	player = &RegularPlayer{}
	return consoleUserInterface
}

func (cui *ConsoleUserInterface) TrackPaused(track *sconsify.Track) {
	gui.setStatus("Paused: " + track.GetFullTitle())
}

func (cui *ConsoleUserInterface) TrackPlaying(track *sconsify.Track) {
	gui.PlayingTrack = track
	gui.setStatus("Playing: " + track.GetFullTitle())
}

func (cui *ConsoleUserInterface) TrackNotAvailable(track *sconsify.Track) {
	gui.flash("Not available: " + track.GetTitle())
}

func (cui *ConsoleUserInterface) Shutdown() {
	persistState()
	events.ShutdownEngine()
}

func (cui *ConsoleUserInterface) PlayTokenLost() error {
	gui.setStatus("Play token lost")
	return nil
}

func (cui *ConsoleUserInterface) GetNextToPlay() *sconsify.Track {
	if !queue.isEmpty() {
		return gui.getNextFromQueue()
	} else if playlists.HasPlaylistSelected() {
		return gui.getNextFromPlaylist()
	}
	return nil
}

func (cui *ConsoleUserInterface) NewPlaylists(newPlaylist sconsify.Playlists) error {
	if playlists == nil {
		playlists = &newPlaylist
		go gui.startGui()
	} else {
		playlists.Merge(&newPlaylist)
		gui.updatePlaylistsView()
		gui.updateTracksView()
		gui.g.Flush()
	}
	return nil
}

func (gui *Gui) startGui() {
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

	if err := gui.g.MainLoop(); err != nil && err != gocui.Quit {
		log.Panicln(err)
	}
}

func (gui *Gui) flash(message string) {
	go func() {
		time.Sleep(4 * time.Second)
		gui.setStatus(gui.currentMessage)
	}()
	gui.updateStatus(message)
}

func (gui *Gui) setStatus(message string) {
	gui.currentMessage = message
	gui.updateCurrentStatus()
}

func (gui *Gui) updateCurrentStatus() {
	gui.updateStatus(gui.currentMessage)
}

func (gui *Gui) updateStatus(message string) {
	gui.clearStatusView()
	fmt.Fprintf(gui.statusView, playlists.GetModeAsString()+"%v\n", message)
	// otherwise the update will appear only in the next keyboard move
	gui.g.Flush()
}

func (gui *Gui) getSelectedPlaylist() *sconsify.Playlist {
	if playlistName, _ := gui.getSelected(gui.playlistsView); playlistName != "" {
		return playlists.Get(playlistName)
	}
	return nil
}

func (gui *Gui) getSelectedTrackName() (string, error) {
	return gui.getSelected(gui.tracksView)
}

func (gui *Gui) getQueueSelectedTrackIndex() int {
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

	return l, err
}

func (gui *Gui) getNextFromPlaylist() *sconsify.Track {
	track, _ := playlists.GetNext()
	return track
}

func (gui *Gui) getNextFromQueue() *sconsify.Track {
	track := queue.Pop()
	go gui.updateQueueView()
	return track
}

func (gui *Gui) playNext() {
	events.NextPlay()
}

func (gui *Gui) replay() {
	events.Replay()
}

func (gui *Gui) getSelectedPlaylistAndTrack() (*sconsify.Playlist, int) {
	currentPlaylist := gui.getSelectedPlaylist()
	currentTrack, errTrack := gui.getSelectedTrackName()
	if currentPlaylist != nil && errTrack == nil {
		if currentPlaylist != nil {
			currentTrack = currentTrack[0:strings.Index(currentTrack, ".")]
			currentIndexTrack, _ := strconv.Atoi(currentTrack)
			currentIndexTrack = currentIndexTrack - 1
			return currentPlaylist, currentIndexTrack
		}
	}
	return nil, -1
}

func (gui *Gui) updateTracksView() {
	gui.tracksView.Clear()
	gui.tracksView.SetCursor(0, 0)
	gui.tracksView.SetOrigin(0, 0)

	if currentPlaylist := gui.getSelectedPlaylist(); currentPlaylist != nil {
		for i := 0; i < currentPlaylist.Tracks(); i++ {
			track := currentPlaylist.Track(i)
			fmt.Fprintf(gui.tracksView, "%v. %v\n", (i + 1), track.GetTitle())
		}
	}
}

func (gui *Gui) updatePlaylistsView() {
	gui.playlistsView.Clear()
	for _, key := range playlists.Names() {
		fmt.Fprintln(gui.playlistsView, key)
		playlist := playlists.Get(key)
		if playlist.IsFolder() {
			for i := 0; i < playlist.Playlists(); i++ {
				subPlaylist := playlist.Playlist(i)
				fmt.Fprintln(gui.playlistsView, subPlaylist.Name())
			}
		}
	}
}

func (gui *Gui) updateQueueView() {
	gui.queueView.Clear()
	if !queue.isEmpty() {
		for _, track := range queue.Contents() {
			fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
		}
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	max := float32(maxX)
	playlistSize := max * 0.15
	trackSize := max * 0.60
	if v, err := g.SetView(VIEW_PLAYLISTS, -1, -1, int(playlistSize), maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.playlistsView = v
		gui.playlistsView.Highlight = true

		gui.updatePlaylistsView()

		if err := g.SetCurrentView(VIEW_PLAYLISTS); err != nil {
			return err
		}
	}
	if v, err := g.SetView(VIEW_TRACKS, int(playlistSize), -1, int(playlistSize+trackSize), maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.tracksView = v

		gui.updateTracksView()
	}
	if v, err := g.SetView(VIEW_QUEUE, int(playlistSize+trackSize), -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.queueView = v
	}
	if v, err := g.SetView(VIEW_STATUS, -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		gui.statusView = v
	}

	if !gui.initialised {
		loadInitialState()
	}
	gui.initialised = true
	return nil
}

func loadInitialState() {
	state := loadState()
	loadPlaylistFromState(state)
	loadTrackFromState(state)
	loadPlayingTrackFromState(state)
}

func loadPlayingTrackFromState(state *State) {
	if state.PlayingTrackUri != "" {
		enablePreviousPlayingTrackToBePlayed(state)
	}
}

func loadPlaylistFromState(state *State) {
	if state.SelectedPlaylist != "" {
		for index, name := range playlists.Names() {
			if name == state.SelectedPlaylist {
				goTo(gui.g, gui.playlistsView, index+1)
				break
			}
		}
	}
}

func loadTrackFromState(state *State) {
	if state.SelectedPlaylist != "" && state.SelectedTrack != "" {
		if playlist := playlists.Get(state.SelectedPlaylist); playlist != nil {
			if index := playlist.IndexByUri(state.SelectedTrack); index != -1 {
				goTo(gui.g, gui.tracksView, index+1)
				gui.enableTracksView()
			}
		}
	}
}

func enablePreviousPlayingTrackToBePlayed(state *State) {
	player = &PersistStatePlayer{
		previousPlayingTrack:    sconsify.InitPartialTrack(state.PlayingTrackUri),
		previousPlayingPlaylist: state.PlayingPlaylist,
	}
	fmt.Fprintf(gui.statusView, "Paused: %v\n", state.PlayingTrackFullTitle)
}

func (gui *Gui) enableTracksView() error {
	gui.tracksView.Highlight = true
	gui.playlistsView.Highlight = false
	gui.queueView.Highlight = false
	return gui.g.SetCurrentView(VIEW_TRACKS)
}

func (gui *Gui) enableSideView() error {
	gui.tracksView.Highlight = false
	gui.playlistsView.Highlight = true
	gui.queueView.Highlight = false
	return gui.g.SetCurrentView(VIEW_PLAYLISTS)
}

func (gui *Gui) enableQueueView() error {
	gui.tracksView.Highlight = false
	gui.playlistsView.Highlight = false
	gui.queueView.Highlight = true
	return gui.g.SetCurrentView(VIEW_QUEUE)
}

func (gui *Gui) clearStatusView() {
	gui.statusView.Clear()
	gui.statusView.SetCursor(0, 0)
	gui.statusView.SetOrigin(0, 0)
}
