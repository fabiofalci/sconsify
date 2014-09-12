package ui

import (
	"testing"
)

func TestStateMode(t *testing.T) {
	state := InitState()

	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Empty state should be on empty state but it is [%v]", mode)
	}

	state.playMode = randomMode
	if mode := state.getModeAsString(); mode != "Random - " {
		t.Errorf("Mode was set to true but is returning [%v]", mode)
	}

	state.playMode = allRandomMode
	if mode := state.getModeAsString(); mode != "All Random - " {
		t.Errorf("Mode was set to true but is returning [%v]", mode)
	}

	state.playMode = normalMode
	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Mode was set to false but is returning [%v]", mode)
	}
}

func TestStateInvertMode(t *testing.T) {
	state := InitState()

	state.invertMode(randomMode)
	if mode := state.getModeAsString(); mode != "Random - " {
		t.Errorf("Mode was inverted to random but is returning [%v]", mode)
	}

	state.invertMode(randomMode)
	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to normal but is returning [%v]", mode)
	}

	state.invertMode(allRandomMode)
	if mode := state.getModeAsString(); mode != "All Random - " {
		t.Errorf("Mode was inverted to allRandom but is returning [%v]", mode)
	}

	state.invertMode(allRandomMode)
	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to random but is returning [%v]", mode)
	}
}
