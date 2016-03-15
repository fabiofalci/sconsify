package simple

import "github.com/fabiofalci/sconsify/sconsify"

func CuiAssertSelectedPlaylist(playlist string) (bool, string) {
	selectedPlaylistName := getPlaylistName(gui.getSelectedPlaylist())

	if playlist != selectedPlaylistName {
		return false, selectedPlaylistName
	}

	return true, ""
}

func CuiAssertSelectedTrack(track string) (bool, string) {
	selectedPlaylist, index := gui.getSelectedPlaylistAndTrack()
	selectedTrack := selectedPlaylist.Track(index)

	if track != selectedTrack.Name {
		return false, selectedTrack.Name
	}

	return true, ""
}

func CuiAssertQueueNextTrack(trackName string) (bool, string) {
	nextTrackName := getTrackName(gui.getNextFromQueue())

	if trackName == nextTrackName {
		return true, ""
	}

	return false, nextTrackName
}

func getPlaylistName(playlist *sconsify.Playlist) string {
	if playlist != nil {
		return playlist.Name()
	}
	return ""
}

func getTrackName(track *sconsify.Track) string {
	if track != nil {
		return track.Name
	}
	return ""
}
