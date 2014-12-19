package sconsify

import (
	"testing"
)

func TestPlaylistTracks(t *testing.T) {
	playlist := createDummyPlaylist()

	if count := playlist.Tracks(); count != 4 {
		t.Errorf("Number of tracks should be 4")
	}
}

func TestPlaylistGetNextTrack(t *testing.T) {
	playlist := createDummyPlaylist()

	if nextIndex := playlist.GetNextTrack(0); nextIndex != 1 {
		t.Errorf("Next track should be track 1")
	}
	if nextIndex := playlist.GetNextTrack(1); nextIndex != 2 {
		t.Errorf("Next track should be track 2")
	}
	if nextIndex := playlist.GetNextTrack(2); nextIndex != 3 {
		t.Errorf("Next track should be track 3")
	}
	if nextIndex := playlist.GetNextTrack(3); nextIndex != 0 {
		t.Errorf("Next track should be track 0")
	}
}

func TestPlaylistTrack(t *testing.T) {
	playlist := createDummyPlaylist()

	if track := playlist.Track(0); track.Uri != "0" {
		t.Errorf("Should be track 0")
	}
	if track := playlist.Track(1); track.Uri != "1" {
		t.Errorf("Should be track 1")
	}
	if track := playlist.Track(2); track.Uri != "2" {
		t.Errorf("Should be track 2")
	}
	if track := playlist.Track(3); track.Uri != "3" {
		t.Errorf("Should be track 3")
	}
}

func createDummyPlaylist() *Playlist {
	tracks := make([]*Track, 4)
	tracks[0] = &Track{Uri: "0"}
	tracks[1] = &Track{Uri: "1"}
	tracks[2] = &Track{Uri: "2"}
	tracks[3] = &Track{Uri: "3"}
	return &Playlist{tracks: tracks}
}
