package sconsify

import (
	"math/rand"
	"testing"
)

func TestNormalMode(t *testing.T) {
	playlists := InitPlaylists()
	if !playlists.isNormalMode() {
		t.Errorf("Playlists initial state should be Normal")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	playlists.SetCurrents("name", 0)

	if track, repeating := playlists.GetNext(); track.Uri != "1" || repeating {
		t.Error("Next track should be 1 and not repeating: ", track.Uri, repeating)
	}
	if track, repeating := playlists.GetNext(); track.Uri != "2" || repeating {
		t.Error("Next track should be 2 and not repeating: ", track.Uri, repeating)
	}
	if track, repeating := playlists.GetNext(); track.Uri != "3" || repeating {
		t.Error("Next track should be 3 and not repeating: "+track.Uri, repeating)
	}
	if track, repeating := playlists.GetNext(); track.Uri != "0" || !repeating {
		t.Error("Next track should be 0 and repeating : ", track.Uri, repeating)
	}

	if track, repeating := playlists.GetNext(); track.Uri != "1" || repeating {
		t.Error("Next track should be 1 and not repeating: ", track.Uri, repeating)
	}
}

func TestShuffleMode(t *testing.T) {
	rand.Seed(123456789) // repeatable

	playlists := InitPlaylists()
	if !playlists.isNormalMode() {
		t.Errorf("Playlists initial state should be Normal")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	playlists.SetCurrents("name", 0)
	playlists.SetMode(ShuffleMode)

	order := []string{"3", "0", "2", "1"}
	for _, expectedUri := range order {
		if track, repeating := playlists.GetNext(); expectedUri != track.Uri || repeating {
			t.Errorf("Random track should be %v and not repeating but it is %v and isRepeating? %v", track.Uri, repeating)
		}
	}

	// now is repeating
	if track, repeating := playlists.GetNext(); track.Uri != "3" || !repeating {
		t.Errorf("Random track should be 3 and repeating but it is %v and isRepeating? %v", repeating)
	}
}

func TestShuffleAllMode(t *testing.T) {
	rand.Seed(123456789) // repeatable

	playlists := InitPlaylists()
	if !playlists.isNormalMode() {
		t.Errorf("Playlists initial state should be Normal")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	playlists.AddPlaylist(createDummyPlaylist("name1"))
	playlists.SetCurrents("name", 0)
	playlists.SetMode(ShuffleAllMode)

	order := []string{"3", "3", "2", "1", "0", "1", "2", "0"}

	for _, expectedUri := range order {
		if track, repeating := playlists.GetNext(); expectedUri != track.Uri || repeating {
			t.Errorf("Random track should be %v and not repeating but it is %v and isRepeating? %v", expectedUri, track.Uri, repeating)
		}
	}

	// now is repeating
	if track, repeating := playlists.GetNext(); track.Uri != "3" || !repeating {
		t.Errorf("Random track should be 3 and repeating but it is %v and isRepeating? %v", track.Uri, repeating)
	}
}

func TestSequentialShuffleMode(t *testing.T) {
	rand.Seed(123456789) // repeatable

	playlists := InitPlaylists()
	if !playlists.isNormalMode() {
		t.Errorf("Playlists initial state should be Normal")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	playlists.AddPlaylist(createDummyPlaylist("name1"))
	playlists.SetCurrents("name", 0)
	playlists.SetMode(SequentialMode)

	order := []string{"0", "1", "2", "3", "0", "1", "2", "3"}

	for _, expectedUri := range order {
		if track, repeating := playlists.GetNext(); expectedUri != track.Uri || repeating {
			t.Errorf("Random track should be %v and not repeating but it is %v and isRepeating? %v", expectedUri, track.Uri, repeating)
		}
	}

	// now is repeating
	if track, repeating := playlists.GetNext(); track.Uri != "0" || !repeating {
		t.Errorf("Random track should be 0 and repeating but it is %v and isRepeating? %v", track.Uri, repeating)
	}
}

func TestPremadeTracks(t *testing.T) {
	playlists := InitPlaylists()
	if playlists.PremadeTracks() != 0 {
		t.Errorf("PremadeTracks should be empty")
	}

	playlists.SetMode(SequentialMode)
	if playlists.PremadeTracks() != 0 {
		t.Errorf("PremadeTracks should be empty")
	}

	playlists = InitPlaylists()
	playlists.AddPlaylist(createDummyPlaylist("name"))
	playlists.SetMode(SequentialMode)

	if playlists.PremadeTracks() != 4 {
		t.Errorf("PremadeTracks should be 4")
	}

	playlists.AddPlaylist(createDummyPlaylist("name1"))
	if playlists.PremadeTracks() != 8 {
		t.Errorf("PremadeTracks should be 8")
	}
}

func TestSetCurrents(t *testing.T) {
	playlists := InitPlaylists()

	if err := playlists.SetCurrents("not to be found", 10); err == nil {
		t.Errorf("Playlist should not be found")
	}
	if playlists.HasPlaylistSelected() {
		t.Errorf("No playlist should be selected")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))

	if err := playlists.SetCurrents("name", 0); err != nil {
		t.Errorf("Playlist and track should be found")
	}
	if err := playlists.SetCurrents("name", 3); err != nil {
		t.Errorf("Playlist and track should be found")
	}
	if !playlists.HasPlaylistSelected() {
		t.Errorf("It has playlist selected")
	}
}

func TestTracks(t *testing.T) {
	playlists := InitPlaylists()

	if playlists.Tracks() != 0 {
		t.Errorf("Tracks should be empty")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	if playlists.Tracks() != 4 {
		t.Errorf("Tracks should be 4")
	}

	playlists.AddPlaylist(createDummyPlaylist("name1"))
	if playlists.Tracks() != 8 {
		t.Errorf("Tracks should be 8")
	}
}

func TestNames(t *testing.T) {
	playlists := InitPlaylists()

	if len(playlists.Names()) != 0 {
		t.Errorf("Playlists should be empty")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))
	names := playlists.Names()
	if len(names) != 1 {
		t.Errorf("Should have only one name")
	}

	playlists.AddPlaylist(createDummyPlaylist("a list"))
	names = playlists.Names()
	if len(names) != 2 {
		t.Errorf("Should have 2 names")
	}

	playlists.AddPlaylist(createDummyPlaylist("z list"))
	names = playlists.Names()

	for i, name := range []string{"a list", "name", "z list"} {
		if name != names[i] {
			t.Errorf("Names is not in alphabetical order")
		}
	}

	playlists.AddPlaylist(createDummyPlaylist("B list"))
	names = playlists.Names()

	for i, name := range []string{"a list", "B list", "name", "z list"} {
		if name != names[i] {
			t.Errorf("Names is not in alphabetical order")
		}
	}
}

func TestGetNext(t *testing.T) {
	playlists := InitPlaylists()

	if track, _ := playlists.GetNext(); track != nil {
		t.Errorf("Track should not be found")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))

	playlists.SetCurrents("name", 0)

	if track, _ := playlists.GetNext(); track != nil && track.Uri != "1" {
		t.Errorf("Next track should be 1")
	}
	if track, _ := playlists.GetNext(); track != nil && track.Uri != "2" {
		t.Errorf("Next track should be 2")
	}
	if track, _ := playlists.GetNext(); track != nil && track.Uri != "3" {
		t.Errorf("Next track should be 3")
	}
}

func TestPlaylists(t *testing.T) {
	playlists := InitPlaylists()

	if playlists.Playlists() != 0 {
		t.Errorf("Playlist count should be empty")
	}

	playlists.AddPlaylist(createDummyPlaylist("name"))

	if count := playlists.Playlists(); count != 1 {
		t.Errorf("Playlist count should be 1 but it is %v", count)
	}

	playlists.AddPlaylist(createDummyPlaylist("name1"))
	if count := playlists.Playlists(); count != 2 {
		t.Errorf("Playlist count should be 2 but it is %v", count)
	}
}

func TestStateMode(t *testing.T) {
	playlists := InitPlaylists()

	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Empty playlists should be on empty state but it is [%v]", mode)
	}

	playlists.playMode = ShuffleMode
	if mode := playlists.GetModeAsString(); mode != "[Shuffled] " {
		t.Errorf("Mode was set to true but is returning %v", mode)
	}

	playlists.playMode = ShuffleAllMode
	if mode := playlists.GetModeAsString(); mode != "[Playlists Shuffled] " {
		t.Errorf("Mode was set to true but is returning %v", mode)
	}

	playlists.playMode = NormalMode
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was set to false but is returning %v", mode)
	}
}

func TestStateInvertMode(t *testing.T) {
	playlists := InitPlaylists()

	playlists.InvertMode(ShuffleMode)
	if mode := playlists.GetModeAsString(); mode != "[Shuffled] " {
		t.Errorf("Mode was inverted to shuffle but is returning %v", mode)
	}

	playlists.InvertMode(ShuffleMode)
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to normal but is returning %v", mode)
	}

	playlists.InvertMode(ShuffleAllMode)
	if mode := playlists.GetModeAsString(); mode != "[Playlists Shuffled] " {
		t.Errorf("Mode was inverted to shuffle all but is returning %v", mode)
	}

	playlists.InvertMode(ShuffleAllMode)
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to shuffle but is returning %v", mode)
	}
}

func TestRemove(t *testing.T) {
	playlists := InitPlaylists()
	playlists.AddPlaylist(createDummyPlaylist("name0"))
	playlists.AddPlaylist(createDummyPlaylist("name1"))
	playlists.AddPlaylist(createDummyPlaylist("name2"))

	if playlists.Playlists() != 3 {
		t.Error("Number of playlists should be 3")
	}

	playlists.Remove("name1")

	if playlists.Playlists() != 2 {
		t.Error("Number of playlists should be 2")
	}

	names := playlists.Names()
	isCorrectNames := (names[0] == "name0" && names[1] == "name2") || (names[0] == "name2" && names[1] == "name0")
	if !isCorrectNames {
		t.Error("Playlist names look wrong after removing name1 playlist")
	}
}

func TestRemoveSubPlaylist(t *testing.T) {
	playlists := InitPlaylists()
	folder0 := createFolder("0", "name0")
	playlists.AddPlaylist(folder0)
	playlists.AddPlaylist(createDummyPlaylistWithId("1", "name1"))

	if folder0.Playlists() != 4 {
		t.Error("Number of playlists should be 4")
	}

	playlists.Remove(" subPlaylist3")

	if folder0.Playlists() != 3 {
		t.Error("Number of playlists should be 3")
	}
}

func TestDuplicatedNames(t *testing.T) {
	playlists := InitPlaylists()

	playlist := createDummyPlaylistWithId("0", "name")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "name" {
		t.Error("Playlist Name should be 'name': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("1", "name")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "name (1)" {
		t.Error("Playlist Name should be 'name (1)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("2", "name")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "name (2)" {
		t.Error("Playlist Name should be 'name (2)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("3", "name (1)")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "name (1) (1)" {
		t.Error("Playlist Name should be 'name (1) (1)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("4", "name (1)")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "name (1) (2)" {
		t.Error("Playlist Name should be 'name (1) (2)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("5", "testing")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "testing" {
		t.Error("Playlist Name should be 'testing': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("6", "testing (1)")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "testing (1)" {
		t.Error("Playlist Name should be 'testing (1)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("7", "testing")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "testing (2)" {
		t.Error("Playlist Name should be 'testing (2)': ", playlist.Name())
	}

	playlist = createDummyPlaylistWithId("8", "testing")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != "testing (3)" {
		t.Error("Playlist Name should be 'testing (3)': ", playlist.Name())
	}
}

func TestDuplicatedNamesWithFolderAndSubPlaylists(t *testing.T) {
	playlists := InitPlaylists()

	playlist := createDummyPlaylistWithId("0", "name")
	playlists.AddPlaylist(playlist)

	folder := createFolder("1", "name")
	playlists.AddPlaylist(folder)

	if folder.Name() != "name (1)" {
		t.Error("Folder Name should be 'name (1)': ", folder.Name())
	}

	// folder has subplaylist: subPlaylist0, subPlaylist1, subPlaylist2, subPlaylist3
	playlist = createDummyPlaylistWithId("2", " subPlaylist3")
	playlists.AddPlaylist(playlist)

	if playlist.Name() != " subPlaylist3 (1)" {
		t.Error("Folder Name should be ' subPlaylist (1)': ", playlist.Name())
	}

	folder = createFolder("1", "name")
	playlists.AddPlaylist(folder)

	if folder.Name() != "name (2)" {
		t.Error("Folder Name should be 'name (2)': ", folder.Name())
	}
	if folder.Playlist(0).Name() != " subPlaylist0 (1)" {
		t.Error("Folder Name should be ' subPlaylist0 (1)': ", folder.Playlist(0).Name())
	}
	if folder.Playlist(1).Name() != " subPlaylist1 (1)" {
		t.Error("Folder Name should be ' subPlaylist1 (1)': ", folder.Playlist(1).Name())
	}
	if folder.Playlist(2).Name() != " subPlaylist2 (1)" {
		t.Error("Folder Name should be ' subPlaylist2 (1)': ", folder.Playlist(2).Name())
	}
	if folder.Playlist(3).Name() != " subPlaylist3 (2)" {
		t.Error("Folder Name should be ' subPlaylist3 (2)': ", folder.Playlist(3).Name())
	}
}

func TestGet(t *testing.T) {
	playlists := InitPlaylists()

	playlists.AddPlaylist(createDummyPlaylistWithId("0", "name"))
	playlists.AddPlaylist(createDummyPlaylistWithId("1", "any"))

	if playlist := playlists.Get("name"); playlist.Name() != "name" {
		t.Error("Playlist Name should be 'name': ", playlist.Name())
	}
	if playlist := playlists.Get("any"); playlist.Name() != "any" {
		t.Error("Playlist Name should be 'any': ", playlist.Name())
	}

	if playlist := playlists.Get("not found"); playlist != nil {
		t.Error("Playlist should not be found")
	}
}

func TestGetWithSubPlaylists(t *testing.T) {
	playlists := InitPlaylists()

	playlists.AddPlaylist(createDummyPlaylistWithId("0", "name"))
	playlists.AddPlaylist(createFolder("0", "folder"))

	if playlist := playlists.Get("folder"); playlist.Name() != "folder" {
		t.Error("Playlist Name should be 'folder': ", playlist.Name())
	}
	if playlist := playlists.Get(" subPlaylist3"); playlist.Name() != " subPlaylist3" {
		t.Error("Playlist Name should be ' subPlaylist3': ", playlist.Name())
	}
}

func TestGetById(t *testing.T) {
	playlists := InitPlaylists()

	playlists.AddPlaylist(createDummyPlaylistWithId("0", "name"))
	playlists.AddPlaylist(createDummyPlaylistWithId("1", "any"))

	if playlist := playlists.GetById("0"); playlist.Id() != "0" {
		t.Error("Playlist Id should be '0': ", playlist.Id())
	}
	if playlist := playlists.GetById("1"); playlist.Id() != "1" {
		t.Error("Playlist Id should be '1': ", playlist.Id())
	}

	if playlist := playlists.GetById("99"); playlist != nil {
		t.Error("Playlist should not be found")
	}
}
