package sconsify

type UserInterface interface {
	TrackPaused(track *Track)
	TrackPlaying(track *Track)
	TrackNotAvailable(track *Track)
	PlayTokenLost() error
	GetNextToPlay() *Track
	NewPlaylists(playlists Playlists) error
	ArtistTopTracks(playlist *Playlist)
	Shutdown()
}
