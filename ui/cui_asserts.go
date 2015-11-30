package ui

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
