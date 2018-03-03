package sconsify

import (
	"strings"
	"sync"
)

type Playlist struct {
	URI string

	tracks []*Track
	name   string
	search bool

	Position int

	subPlaylist bool
	open        bool
	playlists   []*Playlist

	oneTimeLoad      bool
	loadCallback     func(playlist *Playlist)
	loadCallbackOnce sync.Once
}

type PlaylistByName []Playlist

func InitPlaylist(URI string, name string, tracks []*Track) *Playlist {
	return &Playlist{URI: URI, name: name, tracks: tracks}
}

func InitSubPlaylist(URI string, name string, tracks []*Track) *Playlist {
	return &Playlist{URI: URI, name: " " + name, tracks: tracks, subPlaylist: true}
}

func InitSearchPlaylist(URI string, name string, loadCallback func(playlist *Playlist)) *Playlist {
	return &Playlist{URI: URI, name: name, tracks: make([]*Track, 0), search: true, loadCallback: loadCallback}
}

func InitFolder(URI string, name string, playlists []*Playlist) *Playlist {
	folder := &Playlist{URI: URI, name: name, playlists: playlists, search: false, open: true}
	folder.LoadFolderTracks()
	return folder
}

func InitOnDemandPlaylist(URI string, name string, oneTimeLoad bool, loadCallback func(playlist *Playlist)) *Playlist {
	return &Playlist{URI: URI, name: name, tracks: make([]*Track, 0), oneTimeLoad: oneTimeLoad, loadCallback: loadCallback}
}

func InitOnDemandFolder(URI string, name string, oneTimeLoad bool, loadCallback func(playlist *Playlist)) *Playlist {
	playlist := &Playlist{URI: URI, name: name, playlists: make([]*Playlist, 0), oneTimeLoad: oneTimeLoad, loadCallback: loadCallback, open: true, search: false}
	playlist.InvertOpenClose()
	return playlist
}

func (playlist *Playlist) GetNextTrack(currentIndexTrack int) (int, bool) {
	if currentIndexTrack >= len(playlist.tracks)-1 {
		return 0, true
	}
	return currentIndexTrack + 1, false
}

func (playlist *Playlist) FindTrackIndex(query string, fromIndex int) int {
	for i := fromIndex; i < len(playlist.tracks); i++ {
		if strings.Contains(strings.ToUpper(playlist.tracks[i].Name), strings.ToUpper(query)) ||
			strings.Contains(strings.ToUpper(playlist.tracks[i].Artist.Name), strings.ToUpper(query)) {
			return i
		}
	}
	if fromIndex > 0 {
		for i := 0; i < fromIndex; i++ {
			if strings.Contains(strings.ToUpper(playlist.tracks[i].Name), strings.ToUpper(query)) ||
				strings.Contains(strings.ToUpper(playlist.tracks[i].Artist.Name), strings.ToUpper(query)) {
				return i
			}
		}
	}
	return -1
}

func (playlist *Playlist) Track(index int) *Track {
	if index < len(playlist.tracks) {
		return playlist.tracks[index]
	}
	return nil
}

func (playlist *Playlist) AddTrack(track *Track) {
	playlist.tracks = append(playlist.tracks, track)
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

func (playlist *Playlist) IndexByUri(URI string) int {
	for i, track := range playlist.tracks {
		if track.URI == URI {
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
	if strings.HasPrefix(playlist.name, "[") && strings.HasSuffix(playlist.name, "]") {
		return playlist.name[1 : len(playlist.name)-1]
	}
	return playlist.name
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

func (playlist *Playlist) GetPlaylist(name string) *Playlist {
	for _, subPlaylist := range playlist.playlists {
		if subPlaylist.Name() == name {
			return subPlaylist
		}
	}
	return nil
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

func (playlist *Playlist) IsOnDemand() bool {
	return playlist.loadCallback != nil
}

func (playlist *Playlist) ExecuteLoad() {
	if playlist.oneTimeLoad {
		playlist.loadCallbackOnce.Do(playlist.executeLoad)
	} else {
		playlist.executeLoad()
	}
}

func (playlist *Playlist) executeLoad() {
	playlist.loadCallback(playlist)
	if playlist.oneTimeLoad {
		playlist.loadCallback = nil
	}
	if playlist.playlists != nil {
		playlist.LoadFolderTracks()
	}
}

func (playlist *Playlist) LoadFolderTracks() {
	playlist.tracks = make([]*Track, 0)
	for _, subPlaylist := range playlist.playlists {
		for _, track := range subPlaylist.tracks {
			playlist.tracks = append(playlist.tracks, track)
		}
	}
}

func (playlist *Playlist) ToSpotifyID() string {
	return playlist.URI[strings.LastIndex(playlist.URI, ":")+1:]
}

// sort Interface
func (p PlaylistByName) Len() int      { return len(p) }
func (p PlaylistByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PlaylistByName) Less(i, j int) bool {
	return strings.ToLower(p[i].OriginalName()) < strings.ToLower(p[j].OriginalName())
}
