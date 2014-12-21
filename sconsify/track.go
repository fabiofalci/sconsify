package sconsify

import (
	"fmt"

	sp "github.com/op/go-libspotify/spotify"
)

type Track struct {
	Uri      string
	artist   string
	name     string
	duration string
}

func InitTrack(uri string, artist string, name string, duration string) *Track {
	return &Track{
		Uri:      uri,
		artist:   artist,
		name:     name,
		duration: duration,
	}
}

func ToSconsifyTrack(spotifyTracks []*sp.Track) []*Track {
	tracks := make([]*Track, len(spotifyTracks))
	for i, track := range spotifyTracks {
		artist := track.Artist(0)
		artist.Wait()
		tracks[i] = InitTrack(track.Link().String(), artist.Name(), track.Name(), track.Duration().String())
	}
	return tracks
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.artist, track.name, track.duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.artist, track.name)
}
