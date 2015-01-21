package sconsify

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

type Playlists struct {
	playlists         map[string]*Playlist
	currentIndexTrack int
	currentPlaylist   string
	playMode          int
	premadeTracks     []*Track
}

const (
	NormalMode = iota
	RandomMode
	AllRandomMode
	SequentialMode
)

func InitPlaylists() *Playlists {
	playlists := &Playlists{
		playlists: make(map[string]*Playlist),
		playMode:  NormalMode,
	}
	return playlists
}

func (playlists *Playlists) Get(name string) *Playlist {
	for _, playlist := range playlists.playlists {
		if playlist.Name() == name {
			return playlist
		}
	}
	return nil
}

func (playlists *Playlists) Playlists() int {
	return len(playlists.playlists)
}

func (playlists *Playlists) AddPlaylist(playlist *Playlist) {
	playlists.checkDuplicatedNames(playlist, playlist.Name(), 1)
	playlists.playlists[playlist.Id()] = playlist
	playlists.buildPlaylistForNewMode()
}

func (playlists *Playlists) checkDuplicatedNames(newPlaylist *Playlist, originalName string, diff int) {
	for _, playlist := range playlists.playlists {
		if newPlaylist.Name() == playlist.Name() {
			newPlaylist.name = originalName + " (" + strconv.Itoa(diff) + ")"
			diff = diff + 1
			playlists.checkDuplicatedNames(newPlaylist, originalName, diff)
			return
		}
	}
}

func (playlists *Playlists) Merge(newPlaylist *Playlists) {
	for key, value := range newPlaylist.playlists {
		playlists.playlists[key] = value
	}
	playlists.buildPlaylistForNewMode()
}

func (playlists *Playlists) Remove(playlistName string) {
	for key, playlist := range playlists.playlists {
		if playlist.Name() == playlistName {
			delete(playlists.playlists, key)
			playlists.buildPlaylistForNewMode()
			return
		}
	}
}

func (playlists *Playlists) Names() []string {
	names := make([]string, playlists.Playlists())
	i := 0
	for _, playlist := range playlists.playlists {
		names[i] = playlist.Name()
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
		if playlist != nil {
			numberOfTracks = playlist.Tracks()
		}
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
	names := playlists.Names()
	sort.Strings(names)

	index := 0
	for _, name := range names {
		playlist := playlists.Get(name)
		for i := 0; i < playlist.Tracks(); i++ {
			playlists.premadeTracks[index] = playlist.Track(i)
			index++
		}
	}
}

func getRandomPermutation(numberOfTracks int) []int {
	return rand.Perm(numberOfTracks)
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

func (playlists *Playlists) SetCurrents(currentPlaylist string, currentIndexTrack int) error {
	if playlist := playlists.Get(currentPlaylist); playlist != nil {
		if playlist.Tracks() > currentIndexTrack {
			playlists.currentPlaylist = currentPlaylist
			playlists.currentIndexTrack = currentIndexTrack
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Invalid index [%v] track or current playlist [%v]", currentIndexTrack, currentPlaylist))
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
	if playlist != nil {
		playlists.currentIndexTrack = playlist.GetNextTrack(playlists.currentIndexTrack)
		return playlist.Track(playlists.currentIndexTrack), repeating
	}

	return nil, false
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
