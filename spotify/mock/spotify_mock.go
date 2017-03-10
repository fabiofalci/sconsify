package mock

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

var (
	bobMarley    = sconsify.InitArtist("Bob Marley:1", "Bob Marley")
	theRamones   = sconsify.InitArtist("The Ramones:2", "The Ramones")
	elvisPresley = sconsify.InitArtist("Elvis Presley:3", "Elvis Presley")
)

func Initialise(events *sconsify.Events, publisher *sconsify.Publisher) {
	playlists := sconsify.InitPlaylists()

	tracks := make([]*sconsify.Track, 2)
	tracks[0] = sconsify.InitTrack("bobmarley0", bobMarley, "Waiting in vain", "2m3s")
	tracks[1] = sconsify.InitTrack("bobmarley1", bobMarley, "Testing", "5m3s")
	playlists.AddPlaylist(sconsify.InitPlaylist("bobmarleyplaylist0", "Bob Marley", tracks))

	folderPlaylists := make([]*sconsify.Playlist, 2)

	tracks = make([]*sconsify.Track, 2)
	tracks[0] = sconsify.InitTrack("bobmarley2", bobMarley, "Waiting in vain", "2m3s")
	tracks[1] = sconsify.InitTrack("bobmarley3", bobMarley, "Testing", "5m3s")
	folderPlaylists[0] = sconsify.InitPlaylist("bobmarleyplaylist1", " Bob Marley and The Wailers", tracks)

	tracks = make([]*sconsify.Track, 3)
	tracks[0] = sconsify.InitTrack("ramones0", theRamones, "Ramones", "2m3s")
	tracks[1] = sconsify.InitTrack("ramones1", theRamones, "Ramones...", "3m21s")
	tracks[2] = sconsify.InitTrack("ramones2", theRamones, "The Ramones", "1m9s")
	folderPlaylists[1] = sconsify.InitPlaylist("ramonesplaylist0", " The Ramones", tracks)

	playlists.AddPlaylist(sconsify.InitFolder("folder", "My folder", folderPlaylists))

	tracks = make([]*sconsify.Track, 3)
	tracks[0] = sconsify.InitTrack("ramones3", theRamones, "I wanna be sedated", "2m3s")
	tracks[1] = sconsify.InitTrack("ramones4", theRamones, "Pet semetary", "3m21s")
	tracks[2] = sconsify.InitTrack("ramones5", theRamones, "Judy is a punk", "1m9s")
	playlists.AddPlaylist(sconsify.InitPlaylist("ramonesplaylist1", "Ramones", tracks))

	publisher.NewPlaylist(playlists)
	waitForMockEvents(events, publisher)
}

func getSearchedPlaylist() *sconsify.Playlists {
	playlists := sconsify.InitPlaylists()
	playlists.AddPlaylist(sconsify.InitSearchPlaylist("elvispresley1", " Elvis Presley", func(playlist *sconsify.Playlist) {
		playlist.AddTrack(sconsify.InitTrack("elvispreley0", elvisPresley, "Burning Love", "2m3s"))
		playlist.AddTrack(sconsify.InitTrack("elvispreley1", elvisPresley, "Love me tender", "2m4s"))
		playlist.AddTrack(sconsify.InitTrack("elvispreley2", elvisPresley, "It's now or never", "2m5s"))
	}))

	return playlists
}

func waitForMockEvents(events *sconsify.Events, publisher *sconsify.Publisher) {
	for {
		select {
		case <-events.PlayUpdates():
		case <-events.PauseUpdates():
		case <-events.ReplayUpdates():
		case <-events.ShutdownSpotifyUpdates():
			publisher.ShutdownEngine()
		case <-events.SearchUpdates():
			publisher.NewPlaylist(getSearchedPlaylist())
		}
	}
}
