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

func InitPartialTrack(uri string) *Track {
	return &Track{
		Uri: uri,
	}
}

func InitTrack(uri string, artist string, name string, duration string) *Track {
	return &Track{
		Uri:      uri,
		artist:   artist,
		name:     name,
		duration: duration,
	}
}

func ToSconsifyTrack(track *sp.Track) *Track {
	track.Wait()
	artist := track.Artist(0)
	artist.Wait()
	return InitTrack(track.Link().String(), artist.Name(), track.Name(), track.Duration().String())
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.artist, track.name, track.duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.artist, track.name)
}

func (track *Track) IsPartial() bool {
	return track.artist == "" && track.name == "" && track.duration == ""
}
