package sconsify

import (
	"testing"
)

func TestPlaylistTracks(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	if count := playlist.Tracks(); count != 4 {
		t.Errorf("Number of tracks should be 4")
	}
}

func TestPlaylistGetNextTrack(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	if nextIndex, repeating := playlist.GetNextTrack(0); nextIndex != 1 || repeating {
		t.Error("Next track should be 1 and not repeating: ", nextIndex, repeating)
	}
	if nextIndex, repeating := playlist.GetNextTrack(1); nextIndex != 2 || repeating {
		t.Error("Next track should be 2 and not repeating: ", nextIndex, repeating)
	}
	if nextIndex, repeating := playlist.GetNextTrack(2); nextIndex != 3 || repeating {
		t.Error("Next track should be 3 and not repeating: ", nextIndex, repeating)
	}
	if nextIndex, repeating := playlist.GetNextTrack(3); nextIndex != 0 || !repeating {
		t.Error("Next track should be 0 and repeating: ", nextIndex, repeating)
	}
}

func TestPlaylistTrack(t *testing.T) {
	playlist := createDummyPlaylist("testing")

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

	if track := playlist.Track(4); track != nil {
		t.Errorf("Track should be null")
	}
}

func TestPlaylistTrackByUri(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	if track := playlist.TrackByUri("0"); track.Uri != "0" {
		t.Errorf("Should be track 0")
	}
	if track := playlist.TrackByUri("1"); track.Uri != "1" {
		t.Errorf("Should be track 1")
	}
	if track := playlist.TrackByUri("2"); track.Uri != "2" {
		t.Errorf("Should be track 2")
	}
	if track := playlist.TrackByUri("3"); track.Uri != "3" {
		t.Errorf("Should be track 3")
	}

	if track := playlist.TrackByUri("not found"); track != nil {
		t.Errorf("Track should be null")
	}
}

func TestSearchPlaylist(t *testing.T) {
	tracks := make([]*Track, 1)
	tracks[0] = InitTrack("0", "artist0", "name0", "duration0")
	playlist := InitSearchPlaylist("0", "testing", tracks)

	if !playlist.IsSearch() {
		t.Errorf("Should be a search playlists")
	}

	playlist = InitPlaylist("0", "testing", tracks)

	if playlist.IsSearch() {
		t.Errorf("Should not be a search playlists")
	}
}

func createDummyPlaylist(name string) *Playlist {
	tracks := make([]*Track, 4)
	tracks[0] = InitTrack("0", "artist0", "name0", "duration0")
	tracks[1] = InitTrack("1", "artist1", "name1", "duration1")
	tracks[2] = InitTrack("2", "artist2", "name2", "duration2")
	tracks[3] = InitTrack("3", "artist3", "name3", "duration3")
	return InitPlaylist(name, name, tracks)
}

func createDummyPlaylistWithId(id string, name string) *Playlist {
	tracks := make([]*Track, 4)
	tracks[0] = InitTrack("0", "artist0", "name0", "duration0")
	tracks[1] = InitTrack("1", "artist1", "name1", "duration1")
	tracks[2] = InitTrack("2", "artist2", "name2", "duration2")
	tracks[3] = InitTrack("3", "artist3", "name3", "duration3")
	return InitPlaylist(id, name, tracks)
}
