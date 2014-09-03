package ui

type UiState struct {
	currentIndexTrack int
	currentPlaylist   string
	randomMode        bool
	currentMessage    string
}

func InitState() *UiState {
	return &UiState{randomMode: false}
}

func (state *UiState) getModeAsString() string {
	if state.randomMode {
		return "Random - "
	}
	return ""
}

func (state *UiState) invertMode() bool {
	state.randomMode = !state.randomMode
	return state.randomMode
}
