package events

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Events struct {
	playlists     chan map[string]*sp.Playlist
	status        chan string
	ToPlay        chan *sp.Track
	NextPlay      chan bool
	pause         chan bool
	shutdown      chan bool
	playTokenLost chan bool
}

func InitialiseEvents() *Events {
	return &Events{
		playlists:     make(chan map[string]*sp.Playlist),
		status:        make(chan string),
		ToPlay:        make(chan *sp.Track),
		NextPlay:      make(chan bool),
		pause:         make(chan bool),
		playTokenLost: make(chan bool),
		shutdown:      make(chan bool)}
}

func (events *Events) Play(track *sp.Track) {
	events.ToPlay <- track
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

func (events *Events) WaitForStatus() <-chan string {
	return events.status
}

func (events *Events) SetStatus(message string) {
	events.status <- message
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
