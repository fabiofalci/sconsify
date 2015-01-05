package events

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

type Events struct {
	shutdown chan bool

	play   chan *sconsify.Track
	pause  chan bool
	search chan string

	nextPlay          chan bool
	playTokenLost     chan bool
	playlists         chan sconsify.Playlists
	trackNotAvailable chan *sconsify.Track
	trackPlaying      chan *sconsify.Track
	trackPaused       chan *sconsify.Track
}

func InitialiseEvents() *Events {
	return &Events{
		shutdown: make(chan bool),

		play:   make(chan *sconsify.Track),
		pause:  make(chan bool),
		search: make(chan string),

		nextPlay:          make(chan bool),
		playTokenLost:     make(chan bool),
		playlists:         make(chan sconsify.Playlists),
		trackNotAvailable: make(chan *sconsify.Track),
		trackPlaying:      make(chan *sconsify.Track),
		trackPaused:       make(chan *sconsify.Track)}
}

func (events *Events) Shutdown() {
	events.shutdown <- true
}

func (events *Events) ShutdownUpdates() <-chan bool {
	return events.shutdown
}

func (events *Events) TrackPlaying(track *sconsify.Track) {
	select {
	case events.trackPlaying <- track:
	default:
	}
}

func (events *Events) TrackPlayingUpdates() <-chan *sconsify.Track {
	return events.trackPlaying
}

func (events *Events) TrackPaused(track *sconsify.Track) {
	select {
	case events.trackPaused <- track:
	default:
	}
}

func (events *Events) TrackPausedUpdates() <-chan *sconsify.Track {
	return events.trackPaused
}

func (events *Events) Search(query string) {
	events.search <- query
}

func (events *Events) SearchUpdates() <-chan string {
	return events.search
}

func (events *Events) TrackNotAvailable(track *sconsify.Track) {
	events.trackNotAvailable <- track
}

func (events *Events) TrackNotAvailableUpdates() <-chan *sconsify.Track {
	return events.trackNotAvailable
}

func (events *Events) NextPlay() {
	events.nextPlay <- true
}

func (events *Events) NextPlayUpdates() <-chan bool {
	return events.nextPlay
}

func (events *Events) Play(track *sconsify.Track) {
	events.play <- track
}

func (events *Events) PlayUpdates() <-chan *sconsify.Track {
	return events.play
}

func (events *Events) Pause() {
	events.pause <- true
}

func (events *Events) PauseUpdates() <-chan bool {
	return events.pause
}

func (events *Events) NewPlaylist(playlists *sconsify.Playlists) {
	events.playlists <- *playlists
}

func (events *Events) PlaylistsUpdates() <-chan sconsify.Playlists {
	return events.playlists
}

func (events *Events) PlayTokenLost() {
	events.playTokenLost <- true
}

func (events *Events) PlayTokenLostUpdates() <-chan bool {
	return events.playTokenLost
}
