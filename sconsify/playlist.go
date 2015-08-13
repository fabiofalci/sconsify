package sconsify

import (
	"strings"
)

type Playlist struct {
	tracks []*Track
	id     string
	name   string
	search bool

	subPlaylist bool
	open        bool
	playlists   []*Playlist
}

type PlaylistByName []Playlist

func InitPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks}
}

func InitSubPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: " " + name, tracks: tracks, subPlaylist: true}
}

func InitSearchPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks, search: true}
}

func InitFolder(id string, name string, playlists []*Playlist) *Playlist {
	folder := &Playlist{id: id, name: name, playlists: playlists, search: false, open: true}

	folder.tracks = make([]*Track, 0)
	for _, subPlaylist := range playlists {
		for _, track := range subPlaylist.tracks {
			folder.tracks = append(folder.tracks, track)
		}
	}

	return folder
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

func (playlist *Playlist) Playlist(index int) *Playlist {
	if index < len(playlist.playlists) {
		return playlist.playlists[index]
	}
	return nil
}

func (playlist *Playlist) AddPlaylist(subPlaylist *Playlist) bool {
	if !playlist.IsFolder() {
		return false
	}
	playlist.playlists = append(playlist.playlists, subPlaylist)
	return true
}

func (playlist *Playlist) RemovePlaylist(playlistName string) bool {
	if !playlist.IsFolder() {
		return false
	}
	for index, p := range playlist.playlists {
		if p.Name() == playlistName {
			playlist.playlists = append(playlist.playlists[:index], playlist.playlists[index+1:]...)
			return true
		}
	}

	return false
}

func (playlist *Playlist) IndexByUri(uri string) int {
	for i, track := range playlist.tracks {
		if track.Uri == uri {
			return i
		}
	}
	return -1
}

func (playlist *Playlist) HasSameNameIncludingSubPlaylists(otherPlaylist *Playlist) bool {
	if playlist.name == otherPlaylist.name {
		return true
	}
	if playlist.IsFolder() {
		for _, subPlaylist := range playlist.playlists {
			if subPlaylist.name == otherPlaylist.name {
				return true
			}
		}
	}
	return false
}



func (playlist *Playlist) Tracks() int {
	return len(playlist.tracks)
}

func (playlist *Playlist) Playlists() int {
	return len(playlist.playlists)
}

func (playlist *Playlist) Name() string {
	return playlist.name
}

func (playlist *Playlist) OriginalName() string {
	if playlist.IsFolder() && !playlist.IsFolderOpen() {
		return playlist.removeClosedFolderName()
	}
	return playlist.name
}

func (playlist *Playlist) removeClosedFolderName() string {
	return playlist.name[1:len(playlist.name) - 1]
}

func (playlist *Playlist) Id() string {
	return playlist.id
}

func (playlist *Playlist) IsSearch() bool {
	return playlist.search
}

func (playlist *Playlist) IsFolder() bool {
	return playlist.playlists != nil
}

func (playlist *Playlist) IsFolderOpen() bool {
	return playlist.open
}

func (playlist *Playlist) OpenFolder() {
	if !playlist.IsFolderOpen() {
		playlist.InvertOpenClose()
	}
}

func (playlist *Playlist) InvertOpenClose() {
	playlist.open = !playlist.open
	if playlist.open {
		playlist.name = playlist.removeClosedFolderName()
	} else {
		playlist.name = "[" + playlist.name + "]"
	}
}


func (playlist *Playlist) RemoveTrack(index int) {
	if len(playlist.tracks) == 0 || index < 0 || index >= len(playlist.tracks) {
		return
	}
	playlist.tracks = append(playlist.tracks[:index], playlist.tracks[index+1:]...)
}

func (playlist *Playlist) RemoveAllTracks() {
	playlist.tracks = make([]*Track, 0)
}

// sort Interface
func (p PlaylistByName) Len() int      { return len(p) }
func (p PlaylistByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PlaylistByName) Less(i, j int) bool {
	return strings.ToLower(p[i].OriginalName()) < strings.ToLower(p[j].OriginalName())
}
