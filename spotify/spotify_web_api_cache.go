package spotify

import (
	webspotify "github.com/zmb3/spotify"
)

type WebApiCache struct {
	Albums      []webspotify.SavedAlbum
	Songs       []webspotify.SavedTrack
	NewReleases []webspotify.SimplePlaylist
	Artists     []webspotify.FullArtist
}
