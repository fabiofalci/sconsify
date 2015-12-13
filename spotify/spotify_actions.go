package spotify

import (
	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
	"time"
)

func (spotify *Spotify) shutdownSpotify() {
	spotify.session.Logout()
	spotify.initCache()
	spotify.events.ShutdownEngine()
}

func (spotify *Spotify) play(trackUri *sconsify.Track) {
	link, err := spotify.session.ParseLink(trackUri.Uri)
	if err != nil {
		return
	}

	track, err := link.Track()
	if err != nil {
		return
	}

	if trackUri.IsPartial() {
		trackUri = sconsify.ToSconsifyTrack(track)
	}

	if !spotify.isTrackAvailable(track) {
		if trackUri.IsFromWebApi() {
			retry := trackUri.RetryLoading()
			if retry < 4 {
				go func() {
					time.Sleep(100 * time.Millisecond)
					spotify.events.Play(trackUri)
				}()
				return
			}
		}
		spotify.events.TrackNotAvailable(trackUri)
		return
	}
	player := spotify.session.Player()
	if err := player.Load(track); err != nil {
		return
	}
	player.Play()

	spotify.events.TrackPlaying(trackUri)
	spotify.currentTrack = trackUri
}

func (spotify *Spotify) isTrackAvailable(track *sp.Track) bool {
	return track.Availability() == sp.TrackAvailabilityAvailable
}

func (spotify *Spotify) search(query string) {
	searchOptions := &sp.SearchOptions{
		Tracks:    sp.SearchSpec{Offset: 0, Count: 100},
		Albums:    sp.SearchSpec{Offset: 0, Count: 100},
		Artists:   sp.SearchSpec{Offset: 0, Count: 100},
		Playlists: sp.SearchSpec{Offset: 0, Count: 100},
		Type:      sp.SearchStandard,
	}
	search, err := spotify.session.Search(query, searchOptions)
	if err != nil {
		return
	}
	search.Wait()

	numberOfTracks := search.Tracks()
	tracks := make([]*sconsify.Track, numberOfTracks)
	for i := 0; i < numberOfTracks; i++ {
		tracks[i] = sconsify.ToSconsifyTrack(search.Track(i))
	}

	playlists := sconsify.InitPlaylists()
	name := " " + query
	playlists.AddPlaylist(sconsify.InitSearchPlaylist(name, name, tracks))

	spotify.events.NewPlaylist(playlists)
}

func (spotify *Spotify) pause() {
	if spotify.isPausedOrPlaying() {
		if spotify.paused {
			spotify.playCurrentTrack()
		} else {
			spotify.pauseCurrentTrack()
		}
	}
}

func (spotify *Spotify) playCurrentTrack() {
	spotify.play(spotify.currentTrack)
	spotify.paused = false
}

func (spotify *Spotify) pauseCurrentTrack() {
	player := spotify.session.Player()
	player.Pause()
	spotify.events.TrackPaused(spotify.currentTrack)
	spotify.paused = true
}

func (spotify *Spotify) isPausedOrPlaying() bool {
	return spotify.currentTrack != nil
}
