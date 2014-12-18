package sconsify

import (
	"testing"
)

func TestStateMode(t *testing.T) {
	playlists := InitPlaylists()

	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Empty playlists should be on empty state but it is [%v]", mode)
	}

	playlists.playMode = RandomMode
	if mode := playlists.GetModeAsString(); mode != "[Random] " {
		t.Errorf("Mode was set to true but is returning %v", mode)
	}

	playlists.playMode = AllRandomMode
	if mode := playlists.GetModeAsString(); mode != "[All Random] " {
		t.Errorf("Mode was set to true but is returning %v", mode)
	}

	playlists.playMode = NormalMode
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was set to false but is returning %v", mode)
	}
}

func TestStateInvertMode(t *testing.T) {
	playlists := InitPlaylists()

	playlists.InvertMode(RandomMode)
	if mode := playlists.GetModeAsString(); mode != "[Random] " {
		t.Errorf("Mode was inverted to random but is returning %v", mode)
	}

	playlists.InvertMode(RandomMode)
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to normal but is returning %v", mode)
	}

	playlists.InvertMode(AllRandomMode)
	if mode := playlists.GetModeAsString(); mode != "[All Random] " {
		t.Errorf("Mode was inverted to allRandom but is returning %v", mode)
	}

	playlists.InvertMode(AllRandomMode)
	if mode := playlists.GetModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to random but is returning %v", mode)
	}
}
