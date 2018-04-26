package sconsify

import "time"

type UserInterface interface {
	TrackPaused(track *Track)
	TrackPlaying(track *Track)
	TrackNotAvailable(track *Track)
	PlayTokenLost() error
	GetNextToPlay() *Track
	NewPlaylists(playlists Playlists) error
	ArtistAlbums(folder *Playlist)
	Shutdown()
	NewTrackLoaded(duration time.Duration)
	TokenExpired()
}
