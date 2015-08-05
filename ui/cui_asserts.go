package ui

func CuiAssertSelectedPlaylist(playlist string) bool {
	selectedPlaylist := gui.getSelectedPlaylist()

	if playlist != selectedPlaylist.Name() {
		return false
	}

	return true
}
