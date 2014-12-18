package sconsify

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	sp "github.com/op/go-libspotify/spotify"
)

type Playlists struct {
	playlists map[string]*Playlist
}

type Playlist struct {
	tracks []*Track
}

type Track struct {
	Uri      string
	artist   string
	name     string
	duration string
}

func InitPlaylists() *Playlists {
	rand.Seed(time.Now().Unix())

	m := make(map[string]*Playlist)
	playlists := &Playlists{playlists: m}
	return playlists
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

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.artist, track.name, track.duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.artist, track.name)
}

func (playlist *Playlist) GetRandomNextTrack() int {
	return rand.Intn(len(playlist.tracks))
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
