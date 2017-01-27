package sconsify

import (
	"fmt"

	sp "github.com/fabiofalci/go-libspotify/spotify"
)

type Track struct {
	URI string

	Artist     *Artist
	Name       string
	Duration   string
	Album      *Album
	fromWebApi bool
	loadRetry  int
}

func InitPartialTrack(URI string) *Track {
	return &Track{
		URI: URI,
	}
}

func InitTrack(URI string, artist *Artist, name string, duration string) *Track {
	return &Track{
		URI:        URI,
		Artist:     artist,
		Name:       name,
		Duration:   duration,
		fromWebApi: false,
	}
}

func InitWebApiTrack(URI string, artist *Artist, name string, duration string) *Track {
	return &Track{
		URI:        URI,
		Artist:     artist,
		Name:       name,
		Duration:   duration,
		fromWebApi: true,
		loadRetry:  0,
	}
}

func ToSconsifyTrack(track *sp.Track) *Track {
	spArtist := track.Artist(0)
	artist := InitArtist(spArtist.Link().String(), spArtist.Name())
	return InitTrack(track.Link().String(), artist, track.Name(), track.Duration().String())
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.Name, track.Artist.Name, track.Duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.Name, track.Artist.Name)
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
