package sconsify

import (
	"errors"
	"math/rand"
	"sort"
)

type Playlists struct {
	playlists         map[string]*Playlist
	currentIndexTrack int
	currentPlaylist   string
	playMode          int
	premadeTracks     []*Track
}

const (
	NormalMode     = iota
	RandomMode     = iota
	AllRandomMode  = iota
	SequentialMode = iota
)

func InitPlaylists() *Playlists {
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

func (playlists *Playlists) PremadeTracks() int {
	return len(playlists.premadeTracks)
}

func (playlists *Playlists) buildPlaylistForNewMode() error {
	if playlists.isNormalMode() {
		playlists.premadeTracks = nil
		return nil
	}

	var numberOfTracks int
	var playlist *Playlist
	if playlists.isRandomMode() {
		playlist = playlists.Get(playlists.currentPlaylist)
		numberOfTracks = playlist.Tracks()
	} else {
		// all random and sequential
		numberOfTracks = playlists.Tracks()
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	playlists.premadeTracks = make([]*Track, numberOfTracks)
	if playlists.isRandomMode() {
		playlists.buildRandomModeTracks(playlist, numberOfTracks)
	} else if playlists.isAllRandomMode() {
		playlists.buildAllRandomModeTracks(numberOfTracks)
	} else {
		// sequential
		playlists.buildSequentialModeTracks()
	}

	playlists.currentIndexTrack = -1

	return nil
}

func (playlists *Playlists) buildRandomModeTracks(playlist *Playlist, numberOfTracks int) {
	perm := getRandomPermutation(numberOfTracks)

	index := 0
	for i := 0; i < playlist.Tracks(); i++ {
		playlists.premadeTracks[perm[index]] = playlist.Track(i)
		index++
	}
}

func (playlists *Playlists) buildAllRandomModeTracks(numberOfTracks int) {
	perm := getRandomPermutation(numberOfTracks)

	index := 0
	for _, playlist := range playlists.playlists {
		for i := 0; i < playlist.Tracks(); i++ {
			playlists.premadeTracks[perm[index]] = playlist.Track(i)
			index++
		}
	}
}

func (playlists *Playlists) buildSequentialModeTracks() {
	names := playlists.GetNames()
	sort.Strings(names)

	index := 0
	for _, name := range names {
		playlist := playlists.playlists[name]
		for i := 0; i < playlist.Tracks(); i++ {
			playlists.premadeTracks[index] = playlist.Track(i)
			index++
		}
	}
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

func (playlists *Playlists) GetNext() (*Track, bool) {
	repeating := false
	if playlists.premadeTracks != nil {
		playlists.currentIndexTrack++
		if playlists.currentIndexTrack >= len(playlists.premadeTracks) {
			playlists.currentIndexTrack = 0
			repeating = true
		}
		return playlists.premadeTracks[playlists.currentIndexTrack], repeating
	}

	playlist := playlists.Get(playlists.currentPlaylist)
	playlists.currentIndexTrack = playlist.GetNextTrack(playlists.currentIndexTrack)

	return playlist.Track(playlists.currentIndexTrack), repeating
}

func (playlists *Playlists) SetMode(mode int) {
	playlists.playMode = mode
	playlists.buildPlaylistForNewMode()
}

func (playlists *Playlists) InvertMode(mode int) int {
	if mode == playlists.playMode {
		playlists.SetMode(NormalMode)
	} else {
		playlists.SetMode(mode)
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

func (playlists *Playlists) isNormalMode() bool {
	return playlists.playMode == NormalMode
}
