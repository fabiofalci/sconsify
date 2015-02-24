package sconsify

import (
	"strings"
)

type Playlist struct {
	tracks []*Track
	id     string
	name   string
	search bool
}

type PlaylistByName []Playlist

func InitPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks}
}

func InitSearchPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks, search: true}
}

func (playlist *Playlist) GetNextTrack(currentIndexTrack int) (int, bool) {
	if currentIndexTrack >= len(playlist.tracks)-1 {
		return 0, true
	}
	return currentIndexTrack + 1, false
}

func (playlist *Playlist) Track(index int) *Track {
	if index < len(playlist.tracks) {
		return playlist.tracks[index]
	}
	return nil
}

func (playlist *Playlist) IndexByUri(uri string) int {
	for i, track := range playlist.tracks {
		if track.Uri == uri {
			return i
		}
	}
	return -1
}

func (playlist *Playlist) Tracks() int {
	return len(playlist.tracks)
}

func (playlist *Playlist) Name() string {
	return playlist.name
}

func (playlist *Playlist) Id() string {
	return playlist.id
}

func (playlist *Playlist) IsSearch() bool {
	return playlist.search
}

// sort Interface
func (p PlaylistByName) Len() int      { return len(p) }
func (p PlaylistByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PlaylistByName) Less(i, j int) bool {
	return strings.ToLower(p[i].name) < strings.ToLower(p[j].name)
}
