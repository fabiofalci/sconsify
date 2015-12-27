package sconsify

import (
	"fmt"

	sp "github.com/op/go-libspotify/spotify"
)

type Track struct {
	Uri        string
	Artist     *Artist
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

func InitTrack(uri string, artist *Artist, name string, duration string) *Track {
	return &Track{
		Uri:        uri,
		Artist:     artist,
		Name:       name,
		Duration:   duration,
		fromWebApi: false,
	}
}

func InitWebApiTrack(uri string, artist *Artist, name string, duration string) *Track {
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
	spArtist := track.Artist(0)
	artist := InitArtist(spArtist.Link().String(), spArtist.Link().String(), spArtist.Name())
	return InitTrack(track.Link().String(), artist, track.Name(), track.Duration().String())
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.Artist.Name, track.Name, track.Duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.Artist.Name, track.Name)
}

func (track *Track) IsPartial() bool {
	return track.Artist == nil && track.Name == "" && track.Duration == ""
}

func (track *Track) IsFromWebApi() bool {
	return track.fromWebApi
}

func (track *Track) RetryLoading() int {
	track.loadRetry = track.loadRetry + 1
	return track.loadRetry
}
