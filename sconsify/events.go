package sconsify

import (
	"time"
)

type Events struct {
	shutdownEngine  chan bool
	shutdownSpotify chan bool

	play   chan *Track
	pause  chan bool
	search chan string
	replay chan bool
	playPauseToggle  chan bool

	getArtistAlbums chan *Artist
	artistAlbums    chan *Playlist

	nextPlay          chan bool
	playTokenLost     chan bool
	playlists         chan Playlists
	trackNotAvailable chan *Track
	trackPlaying      chan *Track
	trackPaused       chan *Track

	newTrackLoaded	  chan time.Duration
}

func InitialiseEvents() *Events {
	return &Events{
		shutdownEngine:  make(chan bool),
		shutdownSpotify: make(chan bool),

		play:   make(chan *Track),
		pause:  make(chan bool),
		search: make(chan string),
		replay: make(chan bool),
		playPauseToggle:  make(chan bool),

		getArtistAlbums: make(chan *Artist),
		artistAlbums:    make(chan *Playlist),

		nextPlay:          make(chan bool),
		playTokenLost:     make(chan bool),
		playlists:         make(chan Playlists),
		trackNotAvailable: make(chan *Track),
		trackPlaying:      make(chan *Track,2),
		trackPaused:       make(chan *Track),

		newTrackLoaded:    make(chan time.Duration,2),
	}
}

func (events *Events) ShutdownEngine() {
	events.shutdownEngine <- true
}

func (events *Events) ShutdownEngineUpdates() <-chan bool {
	return events.shutdownEngine
}

func (events *Events) ShutdownSpotify() {
	events.shutdownSpotify <- true
}

func (events *Events) ShutdownSpotifyUpdates() <-chan bool {
	return events.shutdownSpotify
}

func (events *Events) TrackPlaying(track *Track) {
	select {
	case events.trackPlaying <- track:
	default:
	}
}

func (events *Events) TrackPlayingUpdates() <-chan *Track {
	return events.trackPlaying
}

func (events *Events) TrackPaused(track *Track) {
	select {
	case events.trackPaused <- track:
	default:
	}
}

func (events *Events) TrackPausedUpdates() <-chan *Track {
	return events.trackPaused
}

func (events *Events) Search(query string) {
	events.search <- query
}

func (events *Events) SearchUpdates() <-chan string {
	return events.search
}

func (events *Events) TrackNotAvailable(track *Track) {
	events.trackNotAvailable <- track
}

func (events *Events) TrackNotAvailableUpdates() <-chan *Track {
	return events.trackNotAvailable
}

func (events *Events) NextPlay() {
	events.nextPlay <- true
}

func (events *Events) NextPlayUpdates() <-chan bool {
	return events.nextPlay
}

func (events *Events) Play(track *Track) {
	events.play <- track
}

func (events *Events) PlayUpdates() <-chan *Track {
	return events.play
}

func (events *Events) Replay() {
	events.replay <- true
}

func (events *Events) ReplayUpdates() <-chan bool {
	return events.replay
}

func (events *Events) Pause() {
	events.pause <- true
}

func (events *Events) PauseUpdates() <-chan bool {
	return events.pause
}

func (events *Events) PlayPauseToggle() {
	events.playPauseToggle <- true
}

func (events *Events) PlayPauseToggleUpdates() <-chan bool {
	return events.playPauseToggle
}

func (events *Events) NewPlaylist(playlists *Playlists) {
	events.playlists <- *playlists
}

func (events *Events) PlaylistsUpdates() <-chan Playlists {
	return events.playlists
}

func (events *Events) PlayTokenLost() {
	events.playTokenLost <- true
}

func (events *Events) PlayTokenLostUpdates() <-chan bool {
	return events.playTokenLost
}

func (events *Events) GetArtistAlbums(artist *Artist) {
	events.getArtistAlbums <- artist
}

func (events *Events) GetArtistAlbumsUpdates() <-chan *Artist {
	return events.getArtistAlbums
}

func (events *Events) ArtistAlbums(folder *Playlist) {
	events.artistAlbums <- folder
}

func (events *Events) ArtistAlbumsUpdates() <-chan *Playlist {
	return events.artistAlbums
}

func (events *Events) NewTrackLoaded(duration time.Duration) {
	select {
	case events.newTrackLoaded <- duration:
	default:
	}
}

func (events *Events) NewTrackLoadedUpdate() <-chan time.Duration {
	return events.newTrackLoaded
}
