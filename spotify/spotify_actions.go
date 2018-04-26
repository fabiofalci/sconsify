package spotify

import (
	sp "github.com/fabiofalci/go-libspotify/spotify"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
	webspotify "github.com/zmb3/spotify"
	"strings"
	"time"
	"github.com/fabiofalci/sconsify/webapi"
)

func (spotify *Spotify) shutdownSpotify() {
	spotify.session.Logout()
	spotify.session.Close()
	spotify.initCache()
	spotify.publisher.ShutdownEngine()
}

func (spotify *Spotify) play(trackUri *sconsify.Track) {

	player := spotify.session.Player()
	if !spotify.paused || spotify.currentTrack != trackUri {
		link, err := spotify.session.ParseLink(trackUri.URI)
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
						spotify.publisher.Play(trackUri)
					}()
					return
				}
			}
			spotify.publisher.TrackNotAvailable(trackUri)
			return
		}
		if err := player.Load(track); err != nil {
			return
		}
		spotify.publisher.NewTrackLoaded(track.Duration())
	}
	player.Play()

	spotify.publisher.TrackPlaying(trackUri)
	spotify.currentTrack = trackUri
	spotify.paused = false
	return
}

func (spotify *Spotify) isTrackAvailable(track *sp.Track) bool {
	return track.Availability() == sp.TrackAvailabilityAvailable
}

func (spotify *Spotify) search(query string) {
	playlists := sconsify.InitPlaylists()

	query = checkAlias(query)
	name := " " + query

	playlist := sconsify.InitSearchPlaylist(name, name, func(playlist *sconsify.Playlist) {
		options := createWebSpotifyOptions(50, playlist.Tracks())

		if webapi.HasTokenExpired()	{
			webapi.LoadTokenFromFile()
			if webapi.HasTokenExpired()	{
				spotify.publisher.TokenExpired()
			} else {
				client := webapi.NewClient()
				spotify.client = &client
			}
		}

		if searchResult, err := spotify.client.SearchOpt(query, webspotify.SearchTypeTrack, options); err == nil {
			numberOfTracks := len(searchResult.Tracks.Tracks)
			infrastructure.Debugf("Search '%v' returned %v track(s)", query, numberOfTracks)
			for _, track := range searchResult.Tracks.Tracks {
				webArtist := track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String()))
				infrastructure.Debugf("\tTrack '%v' (%v)", track.URI, track.Name)
			}
		} else {
			infrastructure.Debugf("Spotify search returning error: %v", err)
		}
	})
	playlist.ExecuteLoad()
	playlists.AddPlaylist(playlist)

	spotify.publisher.NewPlaylist(playlists)
}

func checkAlias(query string) string {
	if strings.HasPrefix(query, "ar:") {
		return strings.Replace(query, "ar:", "artist:", 1)
	} else if strings.HasPrefix(query, "al:") {
		return strings.Replace(query, "al:", "album:", 1)
	} else if strings.HasPrefix(query, "tr:") {
		return strings.Replace(query, "tr:", "track:", 1)
	}
	return query
}

func (spotify *Spotify) playPauseToggle() {
	if spotify.currentTrack != nil {
		if spotify.paused {
			spotify.play(spotify.currentTrack)
		} else {
			spotify.pauseCurrentTrack()
		}
	}
}

func (spotify *Spotify) pause() {
	if spotify.currentTrack != nil && !spotify.paused {
		spotify.pauseCurrentTrack()
	}
}

func (spotify *Spotify) pauseCurrentTrack() {
	player := spotify.session.Player()
	player.Pause()
	spotify.publisher.TrackPaused(spotify.currentTrack)
	spotify.paused = true
}

func (spotify *Spotify) artistAlbums(artist *sconsify.Artist) {
	if simpleAlbumPage, err := spotify.client.GetArtistAlbums(webspotify.ID(artist.GetSpotifyID())); err == nil {
		folder := sconsify.InitFolder(artist.URI, "*"+artist.Name, make([]*sconsify.Playlist, 0))

		if fullTracks, err := spotify.client.GetArtistsTopTracks(webspotify.ID(artist.GetSpotifyID()), "GB"); err == nil {
			tracks := make([]*sconsify.Track, len(fullTracks))
			for i, track := range fullTracks {
				tracks[i] = sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String())
			}

			folder.AddPlaylist(sconsify.InitPlaylist(artist.URI, " "+artist.Name+" Top Tracks", tracks))
		}

		infrastructure.Debugf("# of albums %v", len(simpleAlbumPage.Albums))
		for _, simpleAlbum := range simpleAlbumPage.Albums {
			infrastructure.Debugf("AlbumsID %v = %v", simpleAlbum.URI, simpleAlbum.Name)
			playlist := sconsify.InitOnDemandPlaylist(string(simpleAlbum.URI), " "+simpleAlbum.Name, true, func(playlist *sconsify.Playlist) {
				infrastructure.Debugf("Album id %v", playlist.ToSpotifyID())
				if simpleTrackPage, err := spotify.client.GetAlbumTracks(webspotify.ID(playlist.ToSpotifyID())); err == nil {
					infrastructure.Debugf("# of tracks %v", len(simpleTrackPage.Tracks))
					for _, track := range simpleTrackPage.Tracks {
						playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String()))
					}
				}
			})
			folder.AddPlaylist(playlist)
		}

		spotify.publisher.ArtistAlbums(folder)
	}
}
