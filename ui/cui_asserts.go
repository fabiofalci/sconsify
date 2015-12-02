package ui
import "github.com/fabiofalci/sconsify/sconsify"

func CuiAssertSelectedPlaylist(playlist string) (bool, string) {
	selectedPlaylist := gui.getSelectedPlaylist()

	if playlist != selectedPlaylist.Name() {
		return false, selectedPlaylist.Name()
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

func getTrackName(track *sconsify.Track) string {
	if track != nil {
		return track.Name
	}
	return ""
}
