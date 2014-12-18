package sconsify

import (
	"fmt"
	"math/rand"

	sp "github.com/op/go-libspotify/spotify"
)

type Playlist struct {
	tracks []*Track
}

type Track struct {
	Uri      string
	artist   string
	name     string
	duration string
}

func InitPlaylist(tracks []*sp.Track) *Playlist {
	playlist := &Playlist{}

	playlist.tracks = make([]*Track, len(tracks))
	for i, track := range tracks {
		artist := track.Artist(0)
		artist.Wait()
		playlist.tracks[i] = &Track{
			Uri:      track.Link().String(),
			artist:   artist.Name(),
			name:     track.Name(),
			duration: track.Duration().String(),
		}
	}
	return playlist
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.artist, track.name, track.duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.artist, track.name)
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

func (playlist *Playlist) Track(index int) *Track {
	return playlist.tracks[index]
}

func (playlist *Playlist) Tracks() int {
	return len(playlist.tracks)
}
