package ui

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/jroimartin/gocui"
)

var (
	gui       *Gui
	events    *sconsify.Events
	queue     *Queue
	playlists *sconsify.Playlists
)

type Gui struct {
	g              *gocui.Gui
	playlistsView  *gocui.View
	tracksView     *gocui.View
	statusView     *gocui.View
	queueView      *gocui.View
	currentTrack   *sconsify.Track
	currentMessage string
}

func InitialiseConsoleUserInterface(ev *sconsify.Events) sconsify.UserInterface {
	events = ev
	gui = &Gui{}
	queue = InitQueue()
	return gui
}

func (gui *Gui) TrackPaused(track *sconsify.Track) {
	gui.updateStatus("Paused: "+track.GetFullTitle(), false)
}

func (gui *Gui) TrackPlaying(track *sconsify.Track) {
	gui.updateStatus("Playing: "+track.GetFullTitle(), false)
}

func (gui *Gui) TrackNotAvailable(track *sconsify.Track) {
	gui.updateStatus("Not available: "+track.GetTitle(), true)
}

func (gui *Gui) Shutdown() {}

func (gui *Gui) PlayTokenLost() error {
	gui.updateStatus("Play token lost", false)
	return nil
}

func (gui *Gui) GetNextToPlay() *sconsify.Track {
	if !queue.isEmpty() {
		return gui.getNextFromQueue()
	} else if playlists.HasPlaylistSelected() {
		return gui.getNextFromPlaylist()
	}
	return nil
}

func (gui *Gui) NewPlaylists(newPlaylist sconsify.Playlists) error {
	if playlists == nil {
		playlists = &newPlaylist
		go gui.initGui()
	} else {
		playlists.Merge(&newPlaylist)
		go func() {
			gui.updatePlaylistsView()
			gui.updateTracksView()
			gui.g.Flush()
		}()
	}
	return nil
}

func (gui *Gui) initGui() {
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
		gui.currentMessage = message
	} else {
		go func() {
			time.Sleep(4 * time.Second)
			gui.updateStatus(gui.currentMessage, false)
		}()
	}
	fmt.Fprintf(gui.statusView, playlists.GetModeAsString()+"%v", message)

	// otherwise the update will appear only in the next keyboard move
	gui.g.Flush()
}

func (gui *Gui) getSelectedPlaylist() (string, error) {
	return gui.getSelected(gui.playlistsView)
}

func (gui *Gui) getSelectedTrack() (string, error) {
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

	return l, nil
}

func (gui *Gui) playNextFromPlaylist() {
	track, _ := playlists.GetNext()
	gui.play(track)
}

func (gui *Gui) getNextFromPlaylist() *sconsify.Track {
	gui.currentTrack, _ = playlists.GetNext()
	return gui.currentTrack
}

func (gui *Gui) playNextFromQueue() {
	gui.play(queue.Pop())
	gui.updateQueueView()
}

func (gui *Gui) getNextFromQueue() *sconsify.Track {
	gui.currentTrack = queue.Pop()
	go gui.updateQueueView()
	return gui.currentTrack
}

func (gui *Gui) play(track *sconsify.Track) {
	gui.currentTrack = track
	events.Play(gui.currentTrack)
}

func (gui *Gui) playNext() {
	if !queue.isEmpty() {
		gui.playNextFromQueue()
	} else if playlists.HasPlaylistSelected() {
		gui.playNextFromPlaylist()
	}
}

func getCurrentSelectedTrack() *sconsify.Track {
	currentPlaylist, errPlaylist := gui.getSelectedPlaylist()
	currentTrack, errTrack := gui.getSelectedTrack()
	if errPlaylist == nil && errTrack == nil {
		playlist := playlists.Get(currentPlaylist)

		if playlist != nil {
			currentTrack = currentTrack[0:strings.Index(currentTrack, ".")]
			currentIndexTrack, _ := strconv.Atoi(currentTrack)
			currentIndexTrack = currentIndexTrack - 1
			track := playlist.Track(currentIndexTrack)
			playlists.SetCurrents(currentPlaylist, currentIndexTrack)
			return track
		}
	}
	return nil
}

func (gui *Gui) updateTracksView() {
	gui.tracksView.Clear()
	gui.tracksView.SetCursor(0, 0)
	gui.tracksView.SetOrigin(0, 0)
	currentPlaylist, err := gui.getSelectedPlaylist()
	if err == nil {
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
