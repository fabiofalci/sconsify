package events

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Events struct {
	Initialised chan bool
	Status      chan string
	ToPlay      chan *sp.Track
	NextPlay    chan bool
	Pause       chan bool
}

func InitialiseEvents() Events {
	return Events{Initialised: make(chan bool, 1), Status: make(chan string), ToPlay: make(chan *sp.Track), NextPlay: make(chan bool), Pause: make(chan bool)}
}
