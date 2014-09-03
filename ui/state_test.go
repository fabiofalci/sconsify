package ui

import (
	"testing"
)

func TestStateMode(t *testing.T) {
	state := InitState()

	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Empty state should be on empty state but it is [%v]", mode)
	}

	state.randomMode = true
	if mode := state.getModeAsString(); mode != "Random - " {
		t.Errorf("Mode was set to true but is returning [%v]", mode)
	}

	state.randomMode = false
	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Mode was set to false but is returning [%v]", mode)
	}
}

func TestStateInvertMode(t *testing.T) {
	state := InitState()

	state.invertMode()
	if mode := state.getModeAsString(); mode != "Random - " {
		t.Errorf("Mode was inverted to true but is returning [%v]", mode)
	}

	state.invertMode()
	if mode := state.getModeAsString(); mode != "" {
		t.Errorf("Mode was inverted to false but is returning [%v]", mode)
	}
}
