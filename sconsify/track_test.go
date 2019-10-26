package sconsify

import (
	"testing"
)

func TestCompletedTrack(t *testing.T) {
	artist0 = InitArtist("artist0", "artist0")
	track := InitTrack("0", artist0, "0", "0")

	if track.IsPartial() {
		t.Error("Track should be completed")
	}
}

func TestPartialTrack(t *testing.T) {
	track := InitPartialTrack("0")

	if !track.IsPartial() {
		t.Error("Track should be partial")
	}
}
