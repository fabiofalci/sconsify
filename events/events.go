package events

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Events struct {
	Playlists chan map[string]*sp.Playlist
	Status    chan string
	ToPlay    chan *sp.Track
	NextPlay  chan bool
	Pause     chan bool
}

func InitialiseEvents() *Events {
	return &Events{Playlists: make(chan map[string]*sp.Playlist), Status: make(chan string), ToPlay: make(chan *sp.Track), NextPlay: make(chan bool), Pause: make(chan bool)}
}
