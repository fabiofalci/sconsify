package ui

func CuiAssertSelectedPlaylist(playlist string) (bool, string) {
	selectedPlaylist := gui.getSelectedPlaylist()

	if playlist != selectedPlaylist.Name() {
		return false, selectedPlaylist.Name()
	}

	return true, ""
}
