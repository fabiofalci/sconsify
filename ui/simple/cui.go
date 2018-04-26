package simple

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/ui"
	"github.com/jroimartin/gocui"
)

var (
	gui                  *Gui
	events               *sconsify.Events
	publisher            *sconsify.Publisher
	queue                *ui.Queue
	playlists            *sconsify.Playlists
	consoleUserInterface sconsify.UserInterface
	player               Player
	loadStateWhenInit    bool
	timeLeftChannels     *TimeLeftChannels
)

const (
	PlaylistsView = "playlists"
	TracksView    = "tracks"
	QueueView     = "queue"
	StatusView    = "status"
	TimeLeftView  = "time_left"
)

type ConsoleUserInterface struct{}

type TimeLeftChannels struct {
	timeLeft   chan time.Duration
	songPaused chan bool
}

type Gui struct {
	g             *gocui.Gui
	playlistsView *gocui.View
	tracksView    *gocui.View
	statusView    *gocui.View
	queueView     *gocui.View
	timeLeftView  *gocui.View

	currentMessage string
	initialised    bool
	PlayingTrack   *sconsify.Track
}

func InitialiseConsoleUserInterface(ev *sconsify.Events, p *sconsify.Publisher, loadState bool) sconsify.UserInterface {
	events = ev
	publisher = p
	gui = &Gui{}
	consoleUserInterface = &ConsoleUserInterface{}
	queue = ui.InitQueue()
	player = &RegularPlayer{}
	loadStateWhenInit = loadState
	return consoleUserInterface
}

func (cui *ConsoleUserInterface) TrackPaused(track *sconsify.Track) {
	gui.setStatus("Paused: " + track.GetFullTitle())
	select {
	case timeLeftChannels.songPaused <- true:
	default:
	}
}

func (cui *ConsoleUserInterface) TrackPlaying(track *sconsify.Track) {
	gui.g.Update(func(g *gocui.Gui) error {
		gui.PlayingTrack = track
		gui.setStatus("Playing: " + track.GetFullTitle())
		gui.updateTracksView()
		select {
		case timeLeftChannels.songPaused <- false:
		default:
		}
		return nil
	})
}

func (cui *ConsoleUserInterface) TrackNotAvailable(track *sconsify.Track) {
	gui.g.Update(func(g *gocui.Gui) error {
		gui.flash("Not available: " + track.GetTitle())
		return nil
	})
}

func (cui *ConsoleUserInterface) Shutdown() {
	persistState()
	publisher.ShutdownEngine()
}

func (cui *ConsoleUserInterface) PlayTokenLost() error {
	gui.setStatus("Play token lost")
	return nil
}

func (cui *ConsoleUserInterface) GetNextToPlay() *sconsify.Track {
	if !queue.IsEmpty() {
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
		gui.g.Update(func(g *gocui.Gui) error {
			playlists.Merge(&newPlaylist)
			gui.updatePlaylistsView()
			gui.updateTracksView()
			return nil
		})
	}
	return nil
}

func (cui *ConsoleUserInterface) ArtistAlbums(folder *sconsify.Playlist) {
	gui.g.Update(func(g *gocui.Gui) error {
		playlists.AddPlaylist(folder)
		gui.updatePlaylistsView()
		gui.updateTracksView()
		return nil
	})
}

func (cui *ConsoleUserInterface) NewTrackLoaded(duration time.Duration) {
	select {
	case timeLeftChannels.timeLeft <- duration:
	default:
	}
}

func (cui *ConsoleUserInterface) TokenExpired() {
	gui.g.Update(func(g *gocui.Gui) error {
		gui.flash("Token has expired. To issue a new one: sconsify -issue-web-api-token")
		return nil
	})
}

func (gui *Gui) countdown() {

	var timeLeft time.Duration
	active := false
	for {
		select {
		case paused := <-timeLeftChannels.songPaused:
			active = !paused
		case duration := <-timeLeftChannels.timeLeft:
			timeLeft = duration
			//active = true
		default:
		}

		if active {
			timeLeftCopy := timeLeft
			gui.g.Update(func(g *gocui.Gui) error {
				gui.clearTimeLeftView()
				fmt.Fprintf(gui.timeLeftView, timeLeftCopy.String())
				return nil
			})
		}

		then := time.Now().Round(time.Second)
		time.Sleep(1 * time.Second)
		diff := time.Now().Round(time.Second).Sub(then)
		if active {
			timeLeft = timeLeft - diff
		}
	}
}

func (gui *Gui) startGui() {
	var err error
	gui.g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.g.Close()

	gui.g.SetManagerFunc(layout)
	if err := keybindings(); err != nil {
		log.Panicln(err)
	}
	gui.g.SelBgColor = gocui.ColorGreen
	gui.g.SelFgColor = gocui.ColorBlack
	gui.g.Cursor = true

	// Time left counter thread
	timeLeftChannels = &TimeLeftChannels{
		timeLeft:   make(chan time.Duration, 1),
		songPaused: make(chan bool, 1),
	}
	go gui.countdown()

	if err := gui.g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (gui *Gui) flash(message string) {
	go func() {
		time.Sleep(4 * time.Second)
		gui.g.Update(func(g *gocui.Gui) error {
			gui.setStatus(gui.currentMessage)
			return nil
		})
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
	gui.g.Update(func(g *gocui.Gui) error {
		gui.clearStatusView()
		fmt.Fprintf(gui.statusView, playlists.GetModeAsString()+"%v\n", message)
		return nil
	})
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
	gui.g.Update(func(g *gocui.Gui) error {
		gui.updateQueueView()
		return nil
	})
	return track
}

func (gui *Gui) playNext() {
	publisher.NextPlay()
}

func (gui *Gui) replay() {
	publisher.Replay()
}

func (gui *Gui) createPlaylistFromQueue(playlistName string) {
	gui.g.Update(func(g *gocui.Gui) error {
		unsavedFolder := playlists.Get("*Unsaved")
		if unsavedFolder == nil {
			unsavedFolder = sconsify.InitFolder("*Unsaved", "*Unsaved", make([]*sconsify.Playlist, 0))
			playlists.AddPlaylist(unsavedFolder)
		}

		playlist := unsavedFolder.GetPlaylist(" " + playlistName)
		if playlist == nil {
			playlist = sconsify.InitPlaylist(playlistName, " "+playlistName, make([]*sconsify.Track, 0))
			unsavedFolder.AddPlaylist(playlist)
		}

		for _, track := range queue.Contents() {
			playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), track.Artist, track.Name, track.Duration))
		}

		gui.clearQueueView()
		gui.updatePlaylistsView()
		gui.updateTracksView()
		return nil
	})

}

func (gui *Gui) getSelectedPlaylistAndTrack() (*sconsify.Playlist, int) {
	if currentPlaylist := gui.getSelectedPlaylist(); currentPlaylist != nil {
		if currentTrack, err := gui.getSelectedTrackName(); err == nil {
			if indexSeparator := strings.Index(currentTrack, "."); indexSeparator > 0 {
				if currentIndexTrack, err := strconv.Atoi(currentTrack[0:indexSeparator]); err == nil {
					return currentPlaylist, currentIndexTrack - 1
				}
			}
			if currentPlaylist.IsOnDemand() {
				return currentPlaylist, -1
			}
		}
	}
	return nil, -1
}

func (gui *Gui) updateTracksView() {
	gui.tracksView.Clear()
	cx, cy := gui.tracksView.Cursor()
	ox, oy := gui.tracksView.Origin()
	gui.tracksView.SetCursor(0, 0)
	gui.tracksView.SetOrigin(0, 0)

	PlayingTrackOnView := false
	if currentPlaylist := gui.getSelectedPlaylist(); currentPlaylist != nil {
		for i := 0; i < currentPlaylist.Tracks(); i++ {
			track := currentPlaylist.Track(i)
			if track == gui.PlayingTrack {
				PlayingTrackOnView = true
				fmt.Fprintf(gui.tracksView, "%v. <<%v>>\n", (i + 1), track.GetTitle())
			} else {
				fmt.Fprintf(gui.tracksView, "%v. %v\n", (i + 1), track.GetTitle())
			}
		}
		if currentPlaylist.IsOnDemand() {
			fmt.Fprintf(gui.tracksView, "Press to load\n")
		}
	}
	if PlayingTrackOnView {
		gui.tracksView.SetCursor(cx, cy)
		gui.tracksView.SetOrigin(ox, oy)
	}
}

func (gui *Gui) updatePlaylistsView() {
	gui.playlistsView.Clear()
	position := 0
	for _, key := range playlists.Names() {
		fmt.Fprintln(gui.playlistsView, key)
		position++

		playlist := playlists.Get(key)
		playlist.Position = position
		if playlist.IsFolder() && playlist.IsFolderOpen() {
			for i := 0; i < playlist.Playlists(); i++ {
				subPlaylist := playlist.Playlist(i)
				fmt.Fprintln(gui.playlistsView, subPlaylist.Name())
				position++

				subPlaylist.Position = position
			}
		}
	}
}

func (gui *Gui) clearQueueView() {
	queue.RemoveAll()
	gui.updateQueueView()
}

func (gui *Gui) updateQueueView() {
	gui.queueView.Clear()
	if !queue.IsEmpty() {
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
	if v, err := g.SetView(PlaylistsView, -1, -1, int(playlistSize), maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		gui.playlistsView = v
		gui.playlistsView.Highlight = true

		gui.updatePlaylistsView()

		if _, err := g.SetCurrentView(PlaylistsView); err != nil {
			return err
		}
	}
	if v, err := g.SetView(TracksView, int(playlistSize), -1, int(playlistSize+trackSize), maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		gui.tracksView = v

		gui.updateTracksView()
	}
	if v, err := g.SetView(QueueView, int(playlistSize+trackSize), -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		gui.queueView = v
	}
	if v, err := g.SetView(StatusView, -1, maxY-2, int(playlistSize+trackSize), maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		gui.statusView = v
	}

	if v, err := g.SetView(TimeLeftView, int(playlistSize+trackSize), maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		gui.timeLeftView = v
	}

	if !gui.initialised && loadStateWhenInit {
		loadInitialState()
	}
	gui.initialised = true
	return nil
}

func loadInitialState() {
	state := loadState()
	loadClosedFoldersFromState(state)
	loadPlaylistFromState(state)
	loadTrackFromState(state)
	loadQueueFromState(state)
}

func loadQueueFromState(state *State) {
	for _, track := range state.Queue {
		addToQueue(track)
	}
}

func addToQueue(track *sconsify.Track) bool {
	if queue.Add(track) == nil {
		return false
	}
	fmt.Fprintf(gui.queueView, "%v\n", track.GetTitle())
	return true
}

func loadPlaylistFromState(state *State) {
	if state.SelectedPlaylist != "" {
		position := 0
		for _, name := range playlists.Names() {
			if name == state.SelectedPlaylist {
				goTo(gui.g, gui.playlistsView, position+1)
				return
			}
			position++
			playlist := playlists.Get(name)
			if playlist.IsFolder() && playlist.IsFolderOpen() {
				for i := 0; i < playlist.Playlists(); i++ {
					subPlaylist := playlist.Playlist(i)
					if subPlaylist.Name() == state.SelectedPlaylist {
						goTo(gui.g, gui.playlistsView, position+1)
						return
					}
					position++
				}
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

func loadClosedFoldersFromState(state *State) {
	if state.ClosedFolders != nil && len(state.ClosedFolders) > 0 {
		updatePlaylist := false
		for _, URI := range state.ClosedFolders {
			if playlist := playlists.GetByURI(URI); playlist != nil {
				playlist.InvertOpenClose()
				updatePlaylist = true
			}
		}
		if updatePlaylist {
			gui.updatePlaylistsView()
		}
	}
}

func (gui *Gui) enableTracksView() error {
	gui.tracksView.Highlight = true
	gui.playlistsView.Highlight = false
	gui.queueView.Highlight = false
	_, err := gui.g.SetCurrentView(TracksView)
	return err
}

func (gui *Gui) enableSideView() error {
	gui.tracksView.Highlight = false
	gui.playlistsView.Highlight = true
	gui.queueView.Highlight = false
	_, err := gui.g.SetCurrentView(PlaylistsView)
	return err
}

func (gui *Gui) enableQueueView() error {
	gui.tracksView.Highlight = false
	gui.playlistsView.Highlight = false
	gui.queueView.Highlight = true
	_, err := gui.g.SetCurrentView(QueueView)
	return err
}

func (gui *Gui) clearStatusView() {
	gui.statusView.Clear()
	gui.statusView.SetCursor(0, 0)
	gui.statusView.SetOrigin(0, 0)
}
