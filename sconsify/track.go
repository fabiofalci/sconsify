package sconsify

import (
	"fmt"
)

type Track struct {
	Uri      string
	artist   string
	name     string
	duration string
}

func (track *Track) GetFullTitle() string {
	return fmt.Sprintf("%v - %v [%v]", track.artist, track.name, track.duration)
}

func (track *Track) GetTitle() string {
	return fmt.Sprintf("%v - %v", track.artist, track.name)
}
