package sconsify

import (
	"math/rand"

	sp "github.com/op/go-libspotify/spotify"
)

type TrackContainer struct {
	Playlist *sp.Playlist
	tracks   []*sp.Track
}

func (trackContainer *TrackContainer) SetTracks(inTracks []*sp.Track) {
	trackContainer.tracks = inTracks
}

func (trackContainer *TrackContainer) GetRandomNextTrack() int {
	return rand.Intn(len(trackContainer.tracks))
}

func (trackContainer *TrackContainer) GetNextTrack(currentIndexTrack int) int {
	if currentIndexTrack >= len(trackContainer.tracks)-1 {
		return 0
	}
	return currentIndexTrack + 1
}

func (trackContainer *TrackContainer) Track(index int) *sp.Track {
	return trackContainer.tracks[index]
}

func (trackContainer *TrackContainer) Tracks() int {
	return len(trackContainer.tracks)
}
