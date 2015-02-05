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

	// when random modes or sequential mode we build the tracks here
	premadeTracks *Playlist
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

func (playlists *Playlists) GetById(id string) *Playlist {
	for _, playlist := range playlists.playlists {
		if playlist.Id() == id {
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
	sort.Strings(names)
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
	if playlists.premadeTracks == nil {
		return 0
	}
	return playlists.premadeTracks.Tracks()
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

	var tracks []*Track
	if playlists.isRandomMode() {
		tracks = playlists.buildRandomModeTracks(playlist, numberOfTracks)
	} else if playlists.isAllRandomMode() {
		tracks = playlists.buildAllRandomModeTracks(numberOfTracks)
	} else {
		// sequential
		tracks = playlists.buildSequentialModeTracks()
	}

	playlists.premadeTracks = InitPlaylist("premade", "premade", tracks)
	playlists.currentIndexTrack = -1

	return nil
}

func (playlists *Playlists) buildRandomModeTracks(playlist *Playlist, numberOfTracks int) []*Track {
	tracks := make([]*Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)

	index := 0
	for i := 0; i < playlist.Tracks(); i++ {
		tracks[perm[index]] = playlist.Track(i)
		index++
	}
	return tracks
}

func (playlists *Playlists) buildAllRandomModeTracks(numberOfTracks int) []*Track {
	tracks := make([]*Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)

	index := 0
	for _, playlist := range playlists.playlists {
		for i := 0; i < playlist.Tracks(); i++ {
			tracks[perm[index]] = playlist.Track(i)
			index++
		}
	}
	return tracks
}

func (playlists *Playlists) buildSequentialModeTracks() []*Track {
	names := playlists.Names()
	sort.Strings(names)
	tracks := make([]*Track, playlists.Tracks())

	index := 0
	for _, name := range names {
		playlist := playlists.Get(name)
		for i := 0; i < playlist.Tracks(); i++ {
			tracks[index] = playlist.Track(i)
			index++
		}
	}
	return tracks
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
	if playingPlaylist := playlists.GetPlayingPlaylist(); playingPlaylist != nil {
		var repeating bool
		playlists.currentIndexTrack, repeating = playingPlaylist.GetNextTrack(playlists.currentIndexTrack)
		return playingPlaylist.Track(playlists.currentIndexTrack), repeating
	}
	return nil, false
}

func (playlists *Playlists) GetPlayingTrack() *Track {
	if playingPlaylist := playlists.GetPlayingPlaylist(); playingPlaylist != nil {
		return playingPlaylist.Track(playlists.currentIndexTrack)
	}
	return nil
}

func (playlists *Playlists) GetPlayingPlaylist() *Playlist {
	if playlists.hasPremadeTracks() {
		return playlists.premadeTracks
	} else if playlist := playlists.getCurrentPlaylist(); playlist != nil {
		return playlist
	}
	return nil
}

func (playlists *Playlists) getCurrentPlaylist() *Playlist {
	return playlists.Get(playlists.currentPlaylist)
}

func (playlists *Playlists) hasPremadeTracks() bool {
	return playlists.premadeTracks != nil
}

func (playlists *Playlists) isCurrentTrackOutOfBounds() bool {
	return playlists.currentIndexTrack >= playlists.PremadeTracks()
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
