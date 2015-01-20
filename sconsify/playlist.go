package sconsify

type Playlist struct {
	tracks []*Track
	id     string
	name   string
	search bool
}

func InitPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks}
}

func InitSearchPlaylist(id string, name string, tracks []*Track) *Playlist {
	return &Playlist{id: id, name: name, tracks: tracks, search: true}
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

func (playlist *Playlist) Name() string {
	return playlist.name
}

func (playlist *Playlist) Id() string {
	return playlist.id
}

func (playlist *Playlist) IsSearch() bool {
	return playlist.search
}
