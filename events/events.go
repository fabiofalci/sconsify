package events

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Events struct {
	playlists chan map[string]*sp.Playlist
	status    chan string
	ToPlay    chan *sp.Track
	NextPlay  chan bool
	Pause     chan bool
	Shutdown  chan bool
}

func InitialiseEvents() *Events {
	return &Events{
		playlists: make(chan map[string]*sp.Playlist),
		status:    make(chan string),
		ToPlay:    make(chan *sp.Track),
		NextPlay:  make(chan bool),
		Pause:     make(chan bool),
		Shutdown:  make(chan bool)}
}

func (events *Events) Play(track *sp.Track) {
	events.ToPlay <- track
}

func (events *Events) WaitForPlaylists() <-chan map[string]*sp.Playlist {
	return events.playlists
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
