package ui

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
)

type State struct {
	SelectedPlaylist string
	SelectedTrack    string

	PlayingTrackUri       string
	PlayingTrackFullTitle string
	PlayingPlaylist       string

	ClosedFolders []string
	Queue         []*sconsify.Track
}

func loadState() *State {
	if fileLocation := infrastructure.GetStateFileLocation(); fileLocation != "" {
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
	state := State{
		ClosedFolders: make([]string, 0),
		Queue: make([]*sconsify.Track, 0)}
	selectedPlaylist, index := gui.getSelectedPlaylistAndTrack()

	if selectedPlaylist != nil {
		state.SelectedPlaylist = selectedPlaylist.Name()
		selectedTrack := selectedPlaylist.Track(index)
		if selectedTrack != nil {
			state.SelectedTrack = selectedTrack.Uri
		}
	}

	if playingTrack := playlists.GetPlayingTrack(); playingTrack != nil {
		if playingPlaylist := playlists.GetPlayingPlaylist().Id(); playingPlaylist != "premade" {
			state.PlayingTrackUri = playingTrack.Uri
			state.PlayingTrackFullTitle = playingTrack.GetFullTitle()
			state.PlayingPlaylist = playingPlaylist
		}
	}

	for _, playlistName := range playlists.Names() {
		playlist := playlists.Get(playlistName)
		if playlist.IsFolder() && !playlist.IsFolderOpen() {
			state.ClosedFolders	= append(state.ClosedFolders, playlist.Id())
		}
	}

	for _, track := range queue.Contents() {
		state.Queue = append(state.Queue, track)
	}

	if b, err := json.Marshal(state); err == nil {
		if fileLocation := infrastructure.GetStateFileLocation(); fileLocation != "" {
			infrastructure.SaveFile(fileLocation, b)
		}
	}
}
