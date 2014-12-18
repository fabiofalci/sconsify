package events

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

type Events struct {
	playlists         chan map[string]*sconsify.Playlist
	play              chan *sconsify.Track
	nextPlay          chan bool
	pause             chan bool
	trackNotAvailable chan *sconsify.Track
	trackPlaying      chan *sconsify.Track
	trackPaused       chan *sconsify.Track
	shutdown          chan bool
	playTokenLost     chan bool
	search            chan string
}

func InitialiseEvents() *Events {
	return &Events{
		playlists:         make(chan map[string]*sconsify.Playlist),
		play:              make(chan *sconsify.Track),
		nextPlay:          make(chan bool),
		pause:             make(chan bool),
		trackNotAvailable: make(chan *sconsify.Track),
		trackPlaying:      make(chan *sconsify.Track),
		trackPaused:       make(chan *sconsify.Track),
		playTokenLost:     make(chan bool),
		search:            make(chan string),
		shutdown:          make(chan bool)}
}

func (events *Events) TrackPlaying(track *sconsify.Track) {
	select {
	case events.trackPlaying <- track:
	default:
	}
}

func (events *Events) WaitForTrackPlaying() <-chan *sconsify.Track {
	return events.trackPlaying
}

func (events *Events) TrackPaused(track *sconsify.Track) {
	select {
	case events.trackPaused <- track:
	default:
	}
}

func (events *Events) WaitForTrackPaused() <-chan *sconsify.Track {
	return events.trackPaused
}

func (events *Events) TrackNotAvailable(track *sconsify.Track) {
	events.trackNotAvailable <- track
}

func (events *Events) WaitForTrackNotAvailable() <-chan *sconsify.Track {
	return events.trackNotAvailable
}

func (events *Events) NextPlay() {
	events.nextPlay <- true
}

func (events *Events) WaitForNextPlay() <-chan bool {
	return events.nextPlay
}

func (events *Events) Play(track *sconsify.Track) {
	events.play <- track
}

func (events *Events) WaitPlay() <-chan *sconsify.Track {
	return events.play
}

func (events *Events) Pause() {
	events.pause <- true
}

func (events *Events) WaitForPlaylists() <-chan map[string]*sconsify.Playlist {
	return events.playlists
}

func (events *Events) WaitForPause() <-chan bool {
	return events.pause
}

func (events *Events) NewPlaylist(playlists *map[string]*sconsify.Playlist) {
	events.playlists <- *playlists
}

func (events *Events) WaitForShutdown() <-chan bool {
	return events.shutdown
}

func (events *Events) Shutdown() {
	events.shutdown <- true
}

func (events *Events) PlayTokenLost() {
	events.playTokenLost <- true
}

func (events *Events) WaitForPlayTokenLost() <-chan bool {
	return events.playTokenLost
}

func (events *Events) WaitForSearch() <-chan string {
	return events.search
}

func (events *Events) Search(query string) {
	events.search <- query
}
