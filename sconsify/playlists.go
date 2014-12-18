package sconsify

import (
	"errors"
	"math/rand"
	"time"
)

type Playlists struct {
	playlists         map[string]*Playlist
	currentIndexTrack int
	currentPlaylist   string
	playMode          int
}

const (
	NormalMode    = iota
	RandomMode    = iota
	AllRandomMode = iota
)

func InitPlaylists() *Playlists {
	rand.Seed(time.Now().Unix())

	playlists := &Playlists{
		playlists: make(map[string]*Playlist),
		playMode:  NormalMode,
	}
	return playlists
}

func (playlists *Playlists) Get(name string) *Playlist {
	return playlists.playlists[name]
}

func (playlists *Playlists) Playlists() int {
	return len(playlists.playlists)
}

func (playlists *Playlists) AddPlaylist(name string, playlist *Playlist) {
	playlists.playlists[name] = playlist
}

func (playlists *Playlists) Merge(newPlaylist *Playlists) {
	for key, value := range newPlaylist.playlists {
		playlists.playlists[key] = value
	}
}

func (playlists *Playlists) GetNames() []string {
	names := make([]string, playlists.Playlists())
	i := 0
	for name, _ := range playlists.playlists {
		names[i] = name
		i++
	}
	return names
}

func (playlists *Playlists) Tracks() int {
	numberOfTracks := 0
	for _, playlist := range playlists.playlists {
		numberOfTracks += playlist.Tracks()
	}
	return numberOfTracks
}

func (playlists *Playlists) GetTracks(random *bool) ([]*Track, error) {
	numberOfTracks := playlists.Tracks()
	if numberOfTracks == 0 {
		return nil, errors.New("No tracks selected")
	}

	tracks := make([]*Track, numberOfTracks)

	var perm []int
	if *random {
		perm = getRandomPermutation(numberOfTracks)
	}

	index := 0
	for _, playlist := range playlists.playlists {
		for i := 0; i < playlist.Tracks(); i++ {
			track := playlist.Track(i)

			if *random {
				tracks[perm[index]] = track
			} else {
				tracks[index] = track
			}
			index++
		}
	}

	return tracks, nil
}

func getRandomPermutation(numberOfTracks int) []int {
	return rand.Perm(numberOfTracks)
}

func (playlists *Playlists) GetRandomNextPlaylistAndTrack() (string, int) {
	index := rand.Intn(playlists.Playlists())
	count := 0
	var playlist *Playlist
	var newPlaylistName string
	for key, value := range playlists.playlists {
		if index == count {
			newPlaylistName = key
			playlist = value
			break
		}
		count++
	}
	return newPlaylistName, playlist.GetRandomNextTrack()
}

func (playlists *Playlists) GetModeAsString() string {
	if playlists.playMode == RandomMode {
		return "[Random] "
	}
	if playlists.playMode == AllRandomMode {
		return "[All Random] "
	}
	return ""
}

func (playlists *Playlists) SetCurrents(currentPlaylist string, currentIndexTrack int) {
	playlists.currentPlaylist = currentPlaylist
	playlists.currentIndexTrack = currentIndexTrack
}

func (playlists *Playlists) GetNext() *Track {
	playlist := playlists.Get(playlists.currentPlaylist)
	if playlists.isAllRandomMode() {
		playlists.currentPlaylist, playlists.currentIndexTrack = playlists.GetRandomNextPlaylistAndTrack()
		playlist = playlists.Get(playlists.currentPlaylist)
	} else if playlists.isRandomMode() {
		playlists.currentIndexTrack = playlist.GetRandomNextTrack()
	} else {
		playlists.currentIndexTrack = playlist.GetNextTrack(playlists.currentIndexTrack)
	}

	return playlist.Track(playlists.currentIndexTrack)
}

func (playlists *Playlists) InvertMode(mode int) int {
	if mode == playlists.playMode {
		playlists.playMode = NormalMode
	} else {
		playlists.playMode = mode
	}
	return playlists.playMode
}

func (playlists *Playlists) HasPlaylistSelected() bool {
	return playlists.currentPlaylist != ""
}

func (playlists *Playlists) isAllRandomMode() bool {
	return playlists.playMode == AllRandomMode
}

func (playlists *Playlists) isRandomMode() bool {
	return playlists.playMode == RandomMode
}
