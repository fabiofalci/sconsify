package sconsify

import (
	"math/rand"

	sp "github.com/op/go-libspotify/spotify"
)

type TrackContainer struct {
	tracks []*sp.Track
}

func InitTrackContainer(tracks []*sp.Track) *TrackContainer {
	trackContainer := &TrackContainer{}
	trackContainer.tracks = tracks
	return trackContainer
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
