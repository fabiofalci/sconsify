package sconsify

import (
	"strconv"
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

func TestPlaylistIndexByUri(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	if index := playlist.IndexByUri("0"); index != 0 {
		t.Errorf("Should be track 0")
	}
	if index := playlist.IndexByUri("1"); index != 1 {
		t.Errorf("Should be track 1")
	}
	if index := playlist.IndexByUri("2"); index != 2 {
		t.Errorf("Should be track 2")
	}
	if index := playlist.IndexByUri("3"); index != 3 {
		t.Errorf("Should be track 3")
	}

	if index := playlist.IndexByUri("not found"); index != -1 {
		t.Errorf("Track index should be -1")
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

func TestSubPlaylist(t *testing.T) {
	subPlaylist := createSubPlaylist("1", "testing0")

	if subPlaylist.Name() != " testing0" {
		t.Errorf("Should not be a search playlists")
	}
}

func TestFolderPlaylist(t *testing.T) {
	folder := createFolder("0", "folder0")

	if folder.IsSearch() {
		t.Errorf("Should not be a search playlists")
	}

	if !folder.IsFolder() {
		t.Errorf("Should be a folder")
	}

	if folder.Playlists() != 4 {
		t.Errorf("Folder playlists should be 4")
	}

	tracks := 0
	for i := 0; i < folder.Playlists(); i++ {
		if folder.Playlist(i).Name() != " subPlaylist"+strconv.Itoa(i) {
			t.Errorf("Playlist order is not correct: %v", folder.Playlist(i).Name())
		}
		tracks += folder.Playlist(i).Tracks()
	}

	if tracks != folder.Tracks() {
		t.Errorf("Folder should contain all subplaylist tracks, %v != %v", tracks, folder.Tracks())
	}
}

func TestPlaylistHasSameNameIncludingSubPlaylists(t *testing.T) {
	playlist0 := createDummyPlaylistWithId("0", "testing0")
	playlist1 := createDummyPlaylistWithId("1", "testing1")

	if playlist0.HasSameNameIncludingSubPlaylists(playlist1) || playlist1.HasSameNameIncludingSubPlaylists(playlist0) {
		t.Errorf("Playlists don't have same name")
	}

	otherPlaylist0 := createDummyPlaylistWithId("2", "testing0")
	if !playlist0.HasSameNameIncludingSubPlaylists(otherPlaylist0) || !otherPlaylist0.HasSameNameIncludingSubPlaylists(playlist0) {
		t.Errorf("Playlists have same name")
	}
}

func TestSubPlaylistHasSameNameIncludingSubPlaylists(t *testing.T) {
	playlist0 := createDummyPlaylistWithId("0", "playlist0")
	folder0 := createFolder("0", "folder0")

	if playlist0.HasSameNameIncludingSubPlaylists(folder0) || folder0.HasSameNameIncludingSubPlaylists(playlist0) {
		t.Errorf("Playlists don't have same name")
	}

	otherPlaylist0 := createDummyPlaylistWithId("2", " subPlaylist3")
	if !folder0.HasSameNameIncludingSubPlaylists(otherPlaylist0) {
		t.Errorf("Playlists have same name")
	}
	if otherPlaylist0.HasSameNameIncludingSubPlaylists(folder0) {
		t.Errorf("Playlists don't have same name")
	}

	anotherFolder0 := createFolder("1", "folder0")
	if !anotherFolder0.HasSameNameIncludingSubPlaylists(folder0) || !folder0.HasSameNameIncludingSubPlaylists(anotherFolder0) {
		t.Errorf("Playlists have same name")
	}
}

func TestOpenClose(t *testing.T) {
	folder := createFolder("0", "testing")

	if !folder.IsFolderOpen() {
		t.Errorf("Folder should initialise as open")
	}
	if folder.Name() != folder.OriginalName() || folder.Name() != "testing" {
		t.Errorf("Folder name should not be closed. Name %v, original %v", folder.Name(), folder.OriginalName())
	}

	folder.InvertOpenClose()

	if folder.IsFolderOpen() {
		t.Errorf("Folder should be closed")
	}
	if folder.Name() != "[testing]" && folder.OriginalName() != "testing" {
		t.Errorf("Folder name should be closed. Name %v, original %v", folder.Name(), folder.OriginalName())
	}

	folder.InvertOpenClose()

	if !folder.IsFolderOpen() {
		t.Errorf("Folder should be opened")
	}
	if folder.Name() != folder.OriginalName() || folder.Name() != "testing" {
		t.Errorf("Folder name should not be closed. Name %v, original %v", folder.Name(), folder.OriginalName())
	}
}

func TestOpenFolder(t *testing.T) {
	folder := createFolder("0", "testing")

	if !folder.IsFolderOpen() {
		t.Errorf("Folder should initialise as open")
	}

	folder.InvertOpenClose()

	if folder.IsFolderOpen() {
		t.Errorf("Folder should be closed")
	}

	folder.OpenFolder()

	if !folder.IsFolderOpen() {
		t.Errorf("Folder should be open")
	}

	folder.OpenFolder()

	if !folder.IsFolderOpen() {
		t.Errorf("Folder should be open")
	}
}

func TestAddPlaylist(t *testing.T) {
	folder := createFolder("0", "testing")
	playlist0 := createDummyPlaylistWithId("0", "playlist0")
	playlist1 := createDummyPlaylistWithId("1", "playlist1")

	if !folder.AddPlaylist(playlist0) {
		t.Errorf("Folder should accept new playlist")
	}

	if playlist1.AddPlaylist(playlist0) {
		t.Errorf("Playlist should not accept new playlist")
	}
}

func TestRemovePlaylist(t *testing.T) {
	folder := createFolder("0", "testing")
	playlist0 := createDummyPlaylistWithId("0", "playlist0")

	if !folder.RemovePlaylist(" subPlaylist3") {
		t.Errorf("Folder should remove existing playlist")
	}

	if folder.RemovePlaylist(" subPlaylist99") {
		t.Errorf("Folder should not remove playlist")
	}

	if playlist0.RemovePlaylist(" subPlaylist3") {
		t.Errorf("Playlist should not remove playlist")
	}
}

func TestRemoveTrackFromPlaylist(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	playlist.RemoveTrack(0)

	if count := playlist.Tracks(); count != 3 {
		t.Errorf("Number of tracks should be 3")
	}

	playlist.RemoveTrack(2)

	if count := playlist.Tracks(); count != 2 {
		t.Errorf("Number of tracks should be 2")
	}

	if playlist.Track(0).Name != "name1" {
		t.Errorf("Track name should be name1")
	}

	if playlist.Track(1).Name != "name2" {
		t.Errorf("Track name should be name2")
	}
}

func TestRemoveAllTracksFromPlaylist(t *testing.T) {
	playlist := createDummyPlaylist("testing")

	playlist.RemoveAllTracks()

	if count := playlist.Tracks(); count != 0 {
		t.Errorf("After a remove all track number of tracks should be 0")
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

func createSubPlaylist(id string, name string) *Playlist {
	tracks := make([]*Track, 4)
	tracks[0] = InitTrack("0", "artist0", "name0", "duration0")
	tracks[1] = InitTrack("1", "artist1", "name1", "duration1")
	tracks[2] = InitTrack("2", "artist2", "name2", "duration2")
	tracks[3] = InitTrack("3", "artist3", "name3", "duration3")
	return InitSubPlaylist(id, name, tracks)
}

func createFolder(id string, name string) *Playlist {
	playlists := make([]*Playlist, 4)
	playlists[0] = createSubPlaylist("0", "subPlaylist0")
	playlists[1] = createSubPlaylist("1", "subPlaylist1")
	playlists[2] = createSubPlaylist("2", "subPlaylist2")
	playlists[3] = createSubPlaylist("3", "subPlaylist3")
	return InitFolder(id, name, playlists)
}
