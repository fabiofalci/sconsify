package sconsify

import (
	"time"
)

type Publisher struct {

}

type Events struct {
	shutdownEngine  chan bool
	shutdownSpotify chan bool

	play            chan *Track
	pause           chan bool
	search          chan string
	replay          chan bool
	playPauseToggle chan bool

	getArtistAlbums chan *Artist
	artistAlbums    chan *Playlist

	nextPlay          chan bool
	playTokenLost     chan bool
	playlists         chan Playlists
	trackNotAvailable chan *Track
	trackPlaying      chan *Track
	trackPaused       chan *Track

	newTrackLoaded chan time.Duration
	tokenExpired   chan bool
}

var (
	subscribers []*Events
)

func init() {
	subscribers = make([]*Events, 0)
}

func InitialiseEvents() *Events {
	events := &Events{
		shutdownEngine:  make(chan bool),
		shutdownSpotify: make(chan bool),

		play:            make(chan *Track),
		pause:           make(chan bool),
		search:          make(chan string),
		replay:          make(chan bool),
		playPauseToggle: make(chan bool),

		getArtistAlbums: make(chan *Artist),
		artistAlbums:    make(chan *Playlist),

		nextPlay:          make(chan bool),
		playTokenLost:     make(chan bool),
		playlists:         make(chan Playlists),
		trackNotAvailable: make(chan *Track),
		trackPlaying:      make(chan *Track, 2),
		trackPaused:       make(chan *Track),

		newTrackLoaded: make(chan time.Duration, 2),
		tokenExpired:   make(chan bool),
	}

	subscribers = append(subscribers, events)
	return events
}

func (publisher *Publisher) ShutdownEngine() {
	for _, subscriber := range subscribers {
		subscriber.shutdownEngine <- true
	}
}

func (events *Events) ShutdownEngineUpdates() <-chan bool {
	return events.shutdownEngine
}

func (publisher *Publisher) ShutdownSpotify() {
	for _, subscriber := range subscribers {
		subscriber.shutdownSpotify <- true
	}
}

func (events *Events) ShutdownSpotifyUpdates() <-chan bool {
	return events.shutdownSpotify
}

func (publisher *Publisher) TrackPlaying(track *Track) {
	for _, subscriber := range subscribers {
		subscriber.trackPlaying <- track
	}
}

func (events *Events) TrackPlayingUpdates() <-chan *Track {
	return events.trackPlaying
}

func (publisher *Publisher) TrackPaused(track *Track) {
	for _, subscriber := range subscribers {
		subscriber.trackPaused <- track
	}
}

func (events *Events) TrackPausedUpdates() <-chan *Track {
	return events.trackPaused
}

func (publisher *Publisher) Search(query string) {
	for _, subscriber := range subscribers {
		subscriber.search <- query
	}
}

func (events *Events) SearchUpdates() <-chan string {
	return events.search
}

func (publisher *Publisher) TrackNotAvailable(track *Track) {
	for _, subscriber := range subscribers {
		subscriber.trackNotAvailable <- track
	}
}

func (events *Events) TrackNotAvailableUpdates() <-chan *Track {
	return events.trackNotAvailable
}

func (publisher *Publisher) NextPlay() {
	for _, subscriber := range subscribers {
		subscriber.nextPlay <- true
	}
}

func (events *Events) NextPlayUpdates() <-chan bool {
	return events.nextPlay
}

func (publisher *Publisher) Play(track *Track) {
	for _, subscriber := range subscribers {
		subscriber.play <- track
	}
}

func (events *Events) PlayUpdates() <-chan *Track {
	return events.play
}

func (publisher *Publisher) Replay() {
	for _, subscriber := range subscribers {
		subscriber.replay <- true
	}
}

func (events *Events) ReplayUpdates() <-chan bool {
	return events.replay
}

func (publisher *Publisher) Pause() {
	for _, subscriber := range subscribers {
		subscriber.pause <- true
	}
}

func (events *Events) PauseUpdates() <-chan bool {
	return events.pause
}

func (publisher *Publisher) PlayPauseToggle() {
	for _, subscriber := range subscribers {
		subscriber.playPauseToggle <- true
	}
}

func (events *Events) PlayPauseToggleUpdates() <-chan bool {
	return events.playPauseToggle
}

func (publisher *Publisher) NewPlaylist(playlists *Playlists) {
	for _, subscriber := range subscribers {
		subscriber.playlists <- *playlists
	}
}

func (events *Events) PlaylistsUpdates() <-chan Playlists {
	return events.playlists
}

func (publisher *Publisher) PlayTokenLost() {
	for _, subscriber := range subscribers {
		subscriber.playTokenLost <- true
	}
}

func (events *Events) PlayTokenLostUpdates() <-chan bool {
	return events.playTokenLost
}

func (publisher *Publisher) GetArtistAlbums(artist *Artist) {
	for _, subscriber := range subscribers {
		subscriber.getArtistAlbums <- artist
	}
}

func (events *Events) GetArtistAlbumsUpdates() <-chan *Artist {
	return events.getArtistAlbums
}

func (publisher *Publisher) ArtistAlbums(folder *Playlist) {
	for _, subscriber := range subscribers {
		subscriber.artistAlbums <- folder
	}
}

func (events *Events) ArtistAlbumsUpdates() <-chan *Playlist {
	return events.artistAlbums
}

func (publisher *Publisher) NewTrackLoaded(duration time.Duration) {
	for _, subscriber := range subscribers {
		select {
		case subscriber.newTrackLoaded <- duration:
		default:
		}
	}
}

func (events *Events) NewTrackLoadedUpdate() <-chan time.Duration {
	return events.newTrackLoaded
}

func (publisher *Publisher) TokenExpired() {
	for _, subscriber := range subscribers {
		select {
		case subscriber.tokenExpired <- true:
		default:
		}
	}
}

func (events *Events) TokenExpiredUpdates() <-chan bool {
	return events.tokenExpired
}
