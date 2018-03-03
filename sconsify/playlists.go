package sconsify

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

type Playlists struct {
	playlists         map[string]*Playlist
	currentIndexTrack int
	currentPlaylist   string
	playMode          int

	// when shuffle modes or sequential mode we build the tracks here
	premadeTracks *Playlist
}

const (
	NormalMode = iota
	ShuffleMode
	ShuffleAllMode
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
		if playlist.IsFolder() {
			if p := playlist.GetPlaylist(name); p != nil {
				return p
			}
		}
	}
	return nil
}
func (playlists *Playlists) Find(query string, fromIndex int) *Playlist {
	fromIndex--
	names := playlists.Names()
	for i := fromIndex; i < len(names); i++ {
		if strings.Contains(strings.ToUpper(names[i]), strings.ToUpper(query)) {
			return playlists.Get(names[i])
		}
	}
	if fromIndex > 0 {
		for i := 0; i < fromIndex; i++ {
			if strings.Contains(strings.ToUpper(names[i]), strings.ToUpper(query)) {
				return playlists.Get(names[i])
			}
		}
	}

	return nil
}

func (playlists *Playlists) GetByURI(URI string) *Playlist {
	for _, playlist := range playlists.playlists {
		if playlist.URI == URI {
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
	playlists.playlists[playlist.URI] = playlist
	playlists.buildPlaylistForNewMode()
}

func (playlists *Playlists) checkDuplicatedNames(newPlaylist *Playlist, originalName string, diff int) {
	for _, playlist := range playlists.playlists {
		if playlist.HasSameNameIncludingSubPlaylists(newPlaylist) {
			newPlaylist.name = originalName + " (" + strconv.Itoa(diff) + ")"
			diff = diff + 1
			playlists.checkDuplicatedNames(newPlaylist, originalName, diff)
			break
		}
	}
	if newPlaylist.IsFolder() {
		for _, subPlaylist := range newPlaylist.playlists {
			playlists.checkDuplicatedNames(subPlaylist, subPlaylist.Name(), 1)
		}
	}
}

func (playlists *Playlists) Merge(newPlaylists *Playlists) {
	for key, newPlaylist := range newPlaylists.playlists {
		if newPlaylist.IsSearch() {
			searchPlaylist := playlists.GetByURI("Search")
			if searchPlaylist == nil {
				searchPlaylist = InitFolder("Search", "*Search", make([]*Playlist, 0))
				playlists.AddPlaylist(searchPlaylist)
			}

			searchPlaylist.AddPlaylist(newPlaylist)
			searchPlaylist.OpenFolder()
		} else {
			playlists.playlists[key] = newPlaylist
		}
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
		if playlist.IsFolder() {
			if playlist.RemovePlaylist(playlistName) {
				playlists.buildPlaylistForNewMode()
				return
			}
		}
	}
}

func (playlists *Playlists) playlistsAsArray() []Playlist {
	names := make([]Playlist, playlists.Playlists())
	i := 0
	for _, playlist := range playlists.playlists {
		names[i] = *playlist
		i++
	}
	return names
}

func (playlists *Playlists) Names() []string {
	playlistsAsArray := playlists.playlistsAsArray()
	sort.Sort(PlaylistByName(playlistsAsArray))

	namesAsString := make([]string, playlists.Playlists())
	for index, name := range playlistsAsArray {
		namesAsString[index] = name.name
	}
	return namesAsString
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
	if playlists.isShuffleMode() {
		playlist = playlists.Get(playlists.currentPlaylist)
		if playlist != nil {
			numberOfTracks = playlist.Tracks()
		}
	} else {
		// shuffleall and sequential
		numberOfTracks = playlists.Tracks()
	}

	if numberOfTracks == 0 {
		return errors.New("No tracks selected")
	}

	var tracks []*Track
	if playlists.isShuffleMode() {
		tracks = playlists.shufflePlaylist(playlist, numberOfTracks)
	} else if playlists.isShuffleAllMode() {
		tracks = playlists.shuffleAllPlaylists(numberOfTracks)
	} else {
		// sequential
		tracks = playlists.buildSequentialModeTracks()
	}

	playlists.premadeTracks = InitPlaylist("premade", "premade", tracks)
	playlists.currentIndexTrack = -1

	return nil
}

func (playlists *Playlists) shufflePlaylist(playlist *Playlist, numberOfTracks int) []*Track {
	tracks := make([]*Track, numberOfTracks)
	perm := getRandomPermutation(numberOfTracks)

	index := 0
	for i := 0; i < playlist.Tracks(); i++ {
		tracks[perm[index]] = playlist.Track(i)
		index++
	}
	return tracks
}

func (playlists *Playlists) shuffleAllPlaylists(numberOfTracks int) []*Track {
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
	if playlists.playMode == ShuffleMode {
		return "[Shuffled] "
	}
	if playlists.playMode == ShuffleAllMode {
		return "[Playlists Shuffled] "
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

func (playlists *Playlists) isShuffleAllMode() bool {
	return playlists.playMode == ShuffleAllMode
}

func (playlists *Playlists) isShuffleMode() bool {
	return playlists.playMode == ShuffleMode
}

func (playlists *Playlists) isNormalMode() bool {
	return playlists.playMode == NormalMode
}
