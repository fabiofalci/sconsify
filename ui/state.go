package ui

type UiState struct {
	currentIndexTrack int
	currentPlaylist   string
	playMode          int
	currentMessage    string
}

const (
	normalMode    = iota
	randomMode    = iota
	allRandomMode = iota
)

func InitState() *UiState {
	return &UiState{playMode: normalMode}
}

func (state *UiState) getModeAsString() string {
	if state.playMode == randomMode {
		return "Random - "
	}
	if state.playMode == allRandomMode {
		return "All Random - "
	}
	return ""
}

func (state *UiState) invertMode(mode int) int {
	if mode == state.playMode {
		state.playMode = normalMode
	} else {
		state.playMode = mode
	}
	return state.playMode
}
