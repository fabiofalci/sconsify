package ui

import (
	"fmt"

	sp "github.com/op/go-libspotify/spotify"
)

func formatTrack(status string, track *sp.Track) string {
	artist := track.Artist(0)
	artist.Wait()
	return fmt.Sprintf("%v: %v - %v [%v]", status, artist.Name(), track.Name(), track.Duration().String())
}
