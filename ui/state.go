package ui

type UiState struct {
	currentMessage string
}

func InitState() *UiState {
	return &UiState{}
}
