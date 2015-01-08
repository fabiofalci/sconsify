package sconsify

type Playlist struct {
	tracks []*Track
	name   string
}

func InitPlaylist(name string, tracks []*Track) *Playlist {
	return &Playlist{name: name, tracks: tracks}
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
