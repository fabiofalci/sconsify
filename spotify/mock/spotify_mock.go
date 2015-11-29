package mock

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

func Initialise(events *sconsify.Events) {
	playlists := sconsify.InitPlaylists()

	tracks := make([]*sconsify.Track, 2)
	tracks[0] = sconsify.InitTrack("bobmarley0", "Bob Marley", "Waiting in vain", "2m3s")
	tracks[1] = sconsify.InitTrack("bobmarley1", "Bob Marley", "Testing", "5m3s")
	playlists.AddPlaylist(sconsify.InitPlaylist("bobmarleyplaylist0", "Bob Marley", tracks))

	folderPlaylists := make([]*sconsify.Playlist, 2)

	tracks = make([]*sconsify.Track, 2)
	tracks[0] = sconsify.InitTrack("bobmarley2", "Bob Marley", "Waiting in vain", "2m3s")
	tracks[1] = sconsify.InitTrack("bobmarley3", "Bob Marley", "Testing", "5m3s")
	folderPlaylists[0] = sconsify.InitPlaylist("bobmarleyplaylist1", " Bob Marley and The Wailers", tracks)

	tracks = make([]*sconsify.Track, 3)
	tracks[0] = sconsify.InitTrack("ramones0", "The Ramones", "Ramones", "2m3s")
	tracks[1] = sconsify.InitTrack("ramones1", "The Ramones", "Ramones...", "3m21s")
	tracks[2] = sconsify.InitTrack("ramones2", "The Ramones", "The Ramones", "1m9s")
	folderPlaylists[1] = sconsify.InitPlaylist("ramonesplaylist0", " The Ramones", tracks)

	playlists.AddPlaylist(sconsify.InitFolder("folder", "My folder", folderPlaylists))

	tracks = make([]*sconsify.Track, 3)
	tracks[0] = sconsify.InitTrack("ramones3", "Ramones", "I wanna be sedated", "2m3s")
	tracks[1] = sconsify.InitTrack("ramones4", "Ramones", "Pet semetary", "3m21s")
	tracks[2] = sconsify.InitTrack("ramones5", "Ramones", "Judy is a punk", "1m9s")
	playlists.AddPlaylist(sconsify.InitPlaylist("ramonesplaylist1", "Ramones", tracks))

	events.NewPlaylist(playlists)
	waitForMockEvents(events)
}

func getSearchedPlaylist() *sconsify.Playlists {
	playlists := sconsify.InitPlaylists()
	tracks := make([]*sconsify.Track, 3)
	tracks[0] = sconsify.InitTrack("elvispreley0", "Elvis Presley", "Burning Love", "2m3s")
	tracks[1] = sconsify.InitTrack("elvispreley1", "Elvis Presley", "Love me tender", "2m4s")
	tracks[2] = sconsify.InitTrack("elvispreley2", "Elvis Presley", "It's now or never", "2m5s")
	playlists.AddPlaylist(sconsify.InitSearchPlaylist("elvispresley1", " Elvis Presley", tracks))

	return playlists
}

func waitForMockEvents(events *sconsify.Events) {
	for {
		select {
		case <-events.PlayUpdates():
		case <-events.PauseUpdates():
		case <-events.ReplayUpdates():
		case <-events.ShutdownSpotifyUpdates():
			events.ShutdownEngine()
		case <-events.SearchUpdates():
			events.NewPlaylist(getSearchedPlaylist())
		}
	}
}

