package sconsify

import (
	"math/rand"

	sp "github.com/op/go-libspotify/spotify"
)

type Playlist struct {
	tracks []*sp.Track
}

func InitPlaylist(tracks []*sp.Track) *Playlist {
	playlist := &Playlist{}
	playlist.tracks = tracks
	return playlist
}

func (playlist *Playlist) GetRandomNextTrack() int {
	return rand.Intn(len(playlist.tracks))
}

func (playlist *Playlist) GetNextTrack(currentIndexTrack int) int {
	if currentIndexTrack >= len(playlist.tracks)-1 {
		return 0
	}
	return currentIndexTrack + 1
}

func (playlist *Playlist) Track(index int) *sp.Track {
	return playlist.tracks[index]
}

func (playlist *Playlist) Tracks() int {
	return len(playlist.tracks)
}
