package spotify

import (
	"strings"

	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
	webspotify "github.com/zmb3/spotify"
	"strconv"
)

func (spotify *Spotify) initPlaylist() error {
	playlists := sconsify.InitPlaylists()

	allPlaylists, err := spotify.session.Playlists()
	if err != nil {
		return err
	}
	allPlaylists.Wait()
	var folderPlaylists []*sconsify.Playlist
	var folder *sp.PlaylistFolder
	for i := 0; i < allPlaylists.Playlists(); i++ {
		if allPlaylists.PlaylistType(i) == sp.PlaylistTypeStartFolder {
			folder, _ = allPlaylists.Folder(i)
			folderPlaylists = make([]*sconsify.Playlist, 0)
		} else if allPlaylists.PlaylistType(i) == sp.PlaylistTypeEndFolder {
			if folder != nil {
				playlists.AddPlaylist(sconsify.InitFolder(strconv.Itoa(int(folder.Id())), folder.Name(), folderPlaylists))
			}
			folderPlaylists = nil
			folder = nil
		}

		if allPlaylists.PlaylistType(i) != sp.PlaylistTypePlaylist {
			continue
		}

		playlist := allPlaylists.Playlist(i)
		playlist.Wait()
		if spotify.canAddPlaylist(playlist, allPlaylists.PlaylistType(i)) {
			tracks := make([]*sconsify.Track, playlist.Tracks())
			for i := 0; i < playlist.Tracks(); i++ {
				tracks[i] = spotify.initTrack(playlist.Track(i))
			}
			id := playlist.Link().String()
			if folderPlaylists == nil {
				playlists.AddPlaylist(sconsify.InitPlaylist(id, playlist.Name(), tracks))
			} else {
				folderPlaylists = append(folderPlaylists, sconsify.InitSubPlaylist(id, playlist.Name(), tracks))
			}
		}
	}

	webApiCache := spotify.loadWebApiCache()
	if spotify.client != nil || webApiCache.Albums != nil {
		playlists.AddPlaylist(sconsify.InitOnDemandFolder("Albums", "*Albums", make([]*sconsify.Playlist, 0), func(playlist *sconsify.Playlist) {
			spotify.loadAlbums(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
	}

	if spotify.client != nil || webApiCache.Songs != nil {
		playlists.AddPlaylist(sconsify.InitOnDemandPlaylist("Songs", "*Songs", make([]*sconsify.Track, 0), func(playlist *sconsify.Playlist) {
			spotify.loadSongs(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
	}

	if spotify.client != nil || webApiCache.NewReleases != nil {
		playlists.AddPlaylist(sconsify.InitOnDemandFolder("New Releases", "*New Releases", make([]*sconsify.Playlist, 0), func(playlist *sconsify.Playlist) {
			spotify.loadNewReleases(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
	}

	if spotify.client != nil || webApiCache.Artists != nil {
		playlists.AddPlaylist(sconsify.InitOnDemandFolder("Artists", "*Artists", make([]*sconsify.Playlist, 0), func(playlist *sconsify.Playlist) {
			spotify.loadArtists(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
	}

	spotify.events.NewPlaylist(playlists)
	return nil
}

func (spotify *Spotify) loadAlbums(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	if spotify.client != nil {
		var savedAlbumPage *webspotify.SavedAlbumPage
		var err error
		if savedAlbumPage, err = spotify.client.CurrentUsersAlbumsOpt(createWebSpotifyOptions(50, playlist.Playlists())); err != nil {
			return
		}

		webApiCache.Albums = make([]CachedAlbum, len(savedAlbumPage.Albums))
		for i, album := range savedAlbumPage.Albums {
			tracks := make([]*sconsify.Track, len(album.Tracks.Tracks))
			webApiCache.Albums[i] = CachedAlbum{}
			webApiCache.Albums[i].URI = string(album.URI)
			webApiCache.Albums[i].Name = album.Name
			webApiCache.Albums[i].Tracks = make([]CachedTrack, len(album.Tracks.Tracks))
			for j, track := range album.Tracks.Tracks {
				webArtist := track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				tracks[j] = sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String())
				webApiCache.Albums[i].Tracks[j] = CachedTrack{URI: string(track.URI), Name: track.Name, TimeDuration: track.TimeDuration().String()}
				webApiCache.Albums[i].Tracks[j].ArtistsURI = make([]string, 1)
				webApiCache.Albums[i].Tracks[j].ArtistsURI[0] = string(webArtist.URI)

				cachedArtist := webApiCache.findSharedArtist(string(webArtist.URI))
				if cachedArtist == nil {
					webApiCache.addSharedArtist(CachedArtist{URI: string(webArtist.URI), Name: webArtist.Name})
				}
			}
			playlist.AddPlaylist(sconsify.InitSubPlaylist(string(album.URI), album.Name, tracks))
		}
	} else if webApiCache.Albums != nil {
		for _, album := range webApiCache.Albums {
			tracks := make([]*sconsify.Track, len(album.Tracks))
			for i, track := range album.Tracks {
				webArtist := webApiCache.findSharedArtist(track.ArtistsURI[0])
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				tracks[i] = sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration)
			}
			playlist.AddPlaylist(sconsify.InitSubPlaylist(string(album.URI), album.Name, tracks))
		}
	}

	playlist.OpenFolder()
}

func (spotify *Spotify) loadSongs(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	if spotify.client != nil {
		var savedTrackPage *webspotify.SavedTrackPage
		var err error
		if savedTrackPage, err = spotify.client.CurrentUsersTracksOpt(createWebSpotifyOptions(50, playlist.Tracks())); err != nil {
			return
		}

		if webApiCache.Songs == nil {
			webApiCache.Songs = make([]CachedTrack, len(savedTrackPage.Tracks))
		}
		for _, track := range savedTrackPage.Tracks {
			webArtist := track.Artists[0]
			artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
			playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String()))

			cachedTrack := &CachedTrack{URI: string(track.URI), Name: track.Name, TimeDuration: track.TimeDuration().String()}
			cachedTrack.ArtistsURI = make([]string, 1)
			cachedTrack.ArtistsURI[0] = string(webArtist.URI)

			webApiCache.Songs = append(webApiCache.Songs, *cachedTrack)

			cachedArtist := webApiCache.findSharedArtist(string(webArtist.URI))
			if cachedArtist == nil {
				webApiCache.addSharedArtist(CachedArtist{URI: string(webArtist.URI), Name: webArtist.Name})
			}
		}
	} else if webApiCache.Songs != nil {
		for _, track := range webApiCache.Songs {
			webArtist := webApiCache.findSharedArtist(track.ArtistsURI[0])
			artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
			playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration))
		}
	}
}

func (spotify *Spotify) loadNewReleases(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	var simplePlaylistPage *webspotify.SimplePlaylistPage
	var err error
	if spotify.client != nil {
		if _, simplePlaylistPage, err = spotify.client.FeaturedPlaylistsOpt(&webspotify.PlaylistOptions{Options: *createWebSpotifyOptions(50, playlist.Playlists())}); err != nil {
			return
		}
		webApiCache.NewReleases = simplePlaylistPage.Playlists
	} else if webApiCache.NewReleases == nil {
		return
	}

	for _, album := range webApiCache.NewReleases {
		fullPlaylist, err := spotify.client.GetPlaylist(album.Owner.ID, album.ID)
		if err == nil {
			tracks := make([]*sconsify.Track, len(fullPlaylist.Tracks.Tracks))
			for i, track := range fullPlaylist.Tracks.Tracks {
				webArtist := track.Track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				tracks[i] = sconsify.InitWebApiTrack(string(track.Track.URI), artist, track.Track.Name, track.Track.TimeDuration().String())
			}
			playlist.AddPlaylist(sconsify.InitSubPlaylist(string(album.ID), album.Name, tracks))
		}
		playlist.OpenFolder()
	}
}

func (spotify *Spotify) loadArtists(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	var fullArtistCursorPage *webspotify.FullArtistCursorPage
	var err error
	if spotify.client != nil {
		if fullArtistCursorPage, err = spotify.client.CurrentUsersFollowedArtists(); err != nil {
			return
		}
		webApiCache.Artists = fullArtistCursorPage.Artists
	} else if webApiCache.Artists == nil {
		return
	}

	for _, fullArtist := range webApiCache.Artists {
		tracks := make([]*sconsify.Track, 0)
		playlist.AddPlaylist(sconsify.InitSubPlaylist(string(fullArtist.ID), fullArtist.Name, tracks))
		playlist.OpenFolder()
	}
}

func createWebSpotifyOptions(limit int, offset int) *webspotify.Options {
	return &webspotify.Options{Limit: &limit, Offset: &offset}
}

func (spotify *Spotify) initTrack(playlistTrack *sp.PlaylistTrack) *sconsify.Track {
	track := playlistTrack.Track()
	track.Wait()
	for i := 0; i < track.Artists(); i++ {
		track.Artist(i).Wait()
	}
	return sconsify.ToSconsifyTrack(track)
}

func (spotify *Spotify) canAddPlaylist(playlist *sp.Playlist, playlistType sp.PlaylistType) bool {
	return playlistType == sp.PlaylistTypePlaylist && spotify.isOnFilter(playlist.Name())
}

func (spotify *Spotify) isOnFilter(playlist string) bool {
	if spotify.playlistFilter == nil {
		return true
	}
	for _, filter := range spotify.playlistFilter {
		if filter == playlist {
			return true
		}
	}
	return false
}

func (spotify *Spotify) setPlaylistFilter(playlistFilter string) {
	if playlistFilter == "" {
		return
	}
	spotify.playlistFilter = strings.Split(playlistFilter, ",")
	for i := range spotify.playlistFilter {
		spotify.playlistFilter[i] = strings.Trim(spotify.playlistFilter[i], " ")
	}
}
