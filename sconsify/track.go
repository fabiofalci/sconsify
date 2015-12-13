package sconsify

import (
	"fmt"

	sp "github.com/op/go-libspotify/spotify"
)

type Track struct {
	Uri        string
	Artist     string
	Name       string
	Duration   string
	fromWebApi bool
	loadRetry  int
}

func InitPartialTrack(uri string) *Track {
	return &Track{
		Uri: uri,
	}
}

func InitTrack(uri string, artist string, name string, duration string) *Track {
	return &Track{
		Uri:        uri,
		Artist:     artist,
		Name:       name,
		Duration:   duration,
		fromWebApi: false,
	}
}

func InitWebApiTrack(uri string, artist string, name string, duration string) *Track {
	return &Track{
		Uri:        uri,
		Artist:     artist,
		Name:       name,
		Duration:   duration,
		fromWebApi: true,
		loadRetry:  0,
	}
}

func ToSconsifyTrack(track *sp.Track) *Track {
	artist := track.Artist(0)
	return InitTrack(track.Link().String(), artist.Name(), track.Name(), track.Duration().String())
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.Artist, track.Name, track.Duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.Artist, track.Name)
}

func (track *Track) IsPartial() bool {
	return track.Artist == "" && track.Name == "" && track.Duration == ""
}

func (track *Track) IsFromWebApi() bool {
	return track.fromWebApi
}

func (track *Track) RetryLoading() int {
	track.loadRetry = track.loadRetry + 1
	return track.loadRetry
}
