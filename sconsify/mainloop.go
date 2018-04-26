package sconsify

func StartMainLoop(events *Events, publisher *Publisher, ui UserInterface, askForFirstTrack bool) error {
	select {
	case playlists := <-events.PlaylistsUpdates():
		err := ui.NewPlaylists(playlists)
		if err != nil {
			return err
		}
	case <-events.ShutdownEngineUpdates():
		// TODO it is an error
		return nil
	}

	defer func() {
		publisher.ShutdownSpotify()
		// wait for spotify shutdown
		<-events.ShutdownEngineUpdates()
	}()

	getNextToPlay := func() {
		if track := ui.GetNextToPlay(); track != nil {
			publisher.Play(track)
		}
	}

	if askForFirstTrack {
		getNextToPlay()
	}

	for {
		select {
		case track := <-events.TrackPausedUpdates():
			ui.TrackPaused(track)
		case track := <-events.TrackPlayingUpdates():
			ui.TrackPlaying(track)
		case track := <-events.TrackNotAvailableUpdates():
			ui.TrackNotAvailable(track)
		case <-events.PlayTokenLostUpdates():
			if err := ui.PlayTokenLost(); err != nil {
				return nil
			}
		case <-events.NextPlayUpdates():
			getNextToPlay()
		case newPlaylist := <-events.PlaylistsUpdates():
			ui.NewPlaylists(newPlaylist)
		case playlist := <-events.ArtistAlbumsUpdates():
			ui.ArtistAlbums(playlist)
		case <-events.ShutdownEngineUpdates():
			return nil
		case duration := <-events.NewTrackLoadedUpdate():
			ui.NewTrackLoaded(duration)
		case <-events.TokenExpiredUpdates():
			ui.TokenExpired()
		}
	}

	return nil
}
