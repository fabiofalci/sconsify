package sconsify

import (
	sp "github.com/op/go-libspotify/spotify"
)

type Playlist struct {
	tracks []*Track
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

// FIXME temporary to be used in tests
func CreateDummyPlaylist() *Playlist {
	tracks := make([]*Track, 4)
	tracks[0] = &Track{Uri: "0", artist: "artist0", name: "name0", duration: "duration0"}
	tracks[1] = &Track{Uri: "1", artist: "artist1", name: "name1", duration: "duration1"}
	tracks[2] = &Track{Uri: "2", artist: "artist2", name: "name2", duration: "duration2"}
	tracks[3] = &Track{Uri: "3", artist: "artist3", name: "name3", duration: "duration3"}
	return &Playlist{tracks: tracks}
}
