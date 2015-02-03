package ui

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fabiofalci/sconsify/sconsify"
)

type State struct {
	SelectedPlaylist string
	SelectedTrack    string
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
	selectedPlaylist := gui.getSelectedPlaylist()
	selectedTrack := gui.getCurrentSelectedTrack()
	if selectedPlaylist != nil && selectedTrack != nil {
		state := State{SelectedPlaylist: selectedPlaylist.Name(), SelectedTrack: selectedTrack.Uri}
		if b, err := json.Marshal(state); err == nil {
			if fileLocation := sconsify.GetStateFileLocation(); fileLocation != "" {
				sconsify.SaveFile(fileLocation, b)
			}
		}
	}
}
