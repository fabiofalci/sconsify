package ui

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fabiofalci/sconsify/sconsify"
)

type State struct {
	SelectedPlaylist string
	SelectedTrack    string
	PlayingTrack     string
}

func loadState() *State {
	if fileLocation := sconsify.GetStateFileLocation(); fileLocation != "" {
		if b, err := ioutil.ReadFile(fileLocation); err == nil {
			var state State
			if err := json.Unmarshal(b, &state); err == nil {
				return &state
			}
		}
	}
	return &State{}
}

func persistState() {
	state := State{}
	selectedPlaylist, index := gui.getSelectedPlaylistAndTrack()

	if selectedPlaylist != nil {
		state.SelectedPlaylist = selectedPlaylist.Name()
		selectedTrack := selectedPlaylist.Track(index)
		if selectedTrack != nil {
			state.SelectedTrack = selectedTrack.Uri
		}
	}

	if playingTrack := playlists.GetPlayingTrack(); playingTrack != nil {
		state.PlayingTrack = playingTrack.Uri
	}

	if b, err := json.Marshal(state); err == nil {
		if fileLocation := sconsify.GetStateFileLocation(); fileLocation != "" {
			sconsify.SaveFile(fileLocation, b)
		}
	}
}
