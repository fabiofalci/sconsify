package events

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Events struct {
	playlists         chan map[string]*sp.Playlist
	play              chan *sp.Track
	nextPlay          chan bool
	pause             chan bool
	trackNotAvailable chan bool
	trackPlaying      chan *sp.Track
	trackPaused       chan *sp.Track
	shutdown          chan bool
	playTokenLost     chan bool
	search            chan string
}

func InitialiseEvents() *Events {
	return &Events{
		playlists:         make(chan map[string]*sp.Playlist),
		play:              make(chan *sp.Track),
		nextPlay:          make(chan bool),
		pause:             make(chan bool),
		trackNotAvailable: make(chan bool),
		trackPlaying:      make(chan *sp.Track),
		trackPaused:       make(chan *sp.Track),
		playTokenLost:     make(chan bool),
		search:            make(chan string),
		shutdown:          make(chan bool)}
}

func (events *Events) TrackPlaying(track *sp.Track) {
	select {
	case events.trackPlaying <- track:
	default:
	}
}

func (events *Events) WaitForTrackPlaying() <-chan *sp.Track {
	return events.trackPlaying
}

func (events *Events) TrackPaused(track *sp.Track) {
	select {
	case events.trackPaused <- track:
	default:
	}
}

func (events *Events) WaitForTrackPaused() <-chan *sp.Track {
	return events.trackPaused
}

func (events *Events) TrackNotAvailable() {
	events.trackNotAvailable <- true
}

func (events *Events) WaitForTrackNotAvailable() <-chan bool {
	return events.trackNotAvailable
}

func (events *Events) NextPlay() {
	events.nextPlay <- true
}

func (events *Events) WaitForNextPlay() <-chan bool {
	return events.nextPlay
}

func (events *Events) Play(track *sp.Track) {
	events.play <- track
}

func (events *Events) WaitPlay() <-chan *sp.Track {
	return events.play
}

func (events *Events) Pause() {
	events.pause <- true
}

func (events *Events) WaitForPlaylists() <-chan map[string]*sp.Playlist {
	return events.playlists
}

func (events *Events) WaitForPause() <-chan bool {
	return events.pause
}

func (events *Events) NewPlaylist(playlists *map[string]*sp.Playlist) {
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
