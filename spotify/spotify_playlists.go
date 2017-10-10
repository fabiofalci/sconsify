package spotify

import (
	"errors"
	"strings"

	"fmt"
	sp "github.com/fabiofalci/go-libspotify/spotify"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
	webspotify "github.com/zmb3/spotify"
	"strconv"
)

func (spotify *Spotify) initPlaylist() error {
	playlists := sconsify.InitPlaylists()

	if spotify.client != nil {
		if err := spotify.initWebApiPlaylist(playlists); err != nil {
			return err
		}
	} else {
		if err := spotify.initLibspotifyPlaylist(playlists); err != nil {
			return err
		}
	}

	webApiCache := spotify.loadWebApiCacheIfNecessary()
	if spotify.client != nil {
		playlists.AddPlaylist(sconsify.InitOnDemandFolder("Albums", "*Albums", true, func(playlist *sconsify.Playlist) {
			spotify.loadAlbums(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
		playlists.AddPlaylist(sconsify.InitOnDemandPlaylist("Songs", "*Songs", false, func(playlist *sconsify.Playlist) {
			spotify.loadSongs(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
		playlists.AddPlaylist(sconsify.InitOnDemandFolder("New Releases", "*New Releases", true, func(playlist *sconsify.Playlist) {
			spotify.loadNewReleases(playlist, webApiCache)
			spotify.persistWebApiCache(webApiCache)
		}))
	} else {
		if webApiCache.Albums != nil {
			playlist := sconsify.InitOnDemandFolder("Albums", "*Albums", true, func(playlist *sconsify.Playlist) {
				spotify.loadAlbums(playlist, webApiCache)
			})
			playlist.ExecuteLoad()
			playlists.AddPlaylist(playlist)
		}
		if webApiCache.Songs != nil {
			playlist := sconsify.InitOnDemandPlaylist("Songs", "*Songs", true, func(playlist *sconsify.Playlist) {
				spotify.loadSongs(playlist, webApiCache)
			})
			playlist.ExecuteLoad()
			playlists.AddPlaylist(playlist)
		}
		if webApiCache.NewReleases != nil {
			playlist := sconsify.InitOnDemandFolder("New Releases", "*New Releases", true, func(playlist *sconsify.Playlist) {
				spotify.loadNewReleases(playlist, webApiCache)
			})
			playlist.ExecuteLoad()
			playlists.AddPlaylist(playlist)
		}
	}

	spotify.publisher.NewPlaylist(playlists)
	return nil
}

func (spotify *Spotify) initWebApiPlaylist(playlists *sconsify.Playlists) error {
	if privateUser, err := spotify.client.CurrentUser(); err == nil {
		offset := 0
		total := 1
		for offset <= total {
			offset, total = spotify.loadPlaylists(offset, privateUser, playlists)
			if offset > total {
				fmt.Printf("Loaded %v from %v playlists\n", total, total)
			} else {
				fmt.Printf("Loaded %v from %v playlists\n", offset, total)
			}
			if total == 0 {
				return errors.New("No playlist to load")
			}
		}
	} else {
		return err
	}

	return nil
}

func (spotify *Spotify) initLibspotifyPlaylist(playlists *sconsify.Playlists) error {
	fmt.Print("Warning: not using -web-api flag. Sconsify will load playlists using deprecated libspotify API. If not working try -web-api flag.\n")
	allPlaylists, err := spotify.session.Playlists()
	if err != nil {
		return err
	}
	allPlaylists.Wait()
	var folderPlaylists []*sconsify.Playlist
	var folder *sp.PlaylistFolder
	infrastructure.Debugf("# of playlists %v", allPlaylists.Playlists())
	for i := 0; i < allPlaylists.Playlists(); i++ {
		if allPlaylists.PlaylistType(i) == sp.PlaylistTypeStartFolder {
			folder, _ = allPlaylists.Folder(i)
			folderPlaylists = make([]*sconsify.Playlist, 0)
			infrastructure.Debugf("Opening folder '%v' (%v)", folder.Id(), folder.Name())
		} else if allPlaylists.PlaylistType(i) == sp.PlaylistTypeEndFolder {
			if folder != nil {
				playlists.AddPlaylist(sconsify.InitFolder(strconv.Itoa(int(folder.Id())), folder.Name(), folderPlaylists))
				infrastructure.Debugf("Closing folder '%v' (%v)", folder.Id(), folder.Name())
			} else {
				infrastructure.Debug("Closing a null folder, this doesn't look right ")
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
			id := playlist.Link().String()
			infrastructure.Debugf("Playlist '%v' (%v)", id, playlist.Name())
			tracks := make([]*sconsify.Track, playlist.Tracks())
			infrastructure.Debugf("\t# of tracks %v", playlist.Tracks())
			for i := 0; i < playlist.Tracks(); i++ {
				tracks[i] = spotify.initTrack(playlist.Track(i))
			}
			if folderPlaylists == nil {
				playlists.AddPlaylist(sconsify.InitPlaylist(id, playlist.Name(), tracks))
			} else {
				folderPlaylists = append(folderPlaylists, sconsify.InitSubPlaylist(id, playlist.Name(), tracks))
			}
		}
	}
	return nil
}

func (spotify *Spotify) loadPlaylists(offset int, privateUser *webspotify.PrivateUser, playlists *sconsify.Playlists) (int, int) {
	limit := 50
	options := &webspotify.Options{Limit: &limit, Offset: &offset}
	if simplePlaylistPage, err := spotify.client.GetPlaylistsForUserOpt(privateUser.ID, options); err == nil {
		for _, webPlaylist := range simplePlaylistPage.Playlists {
			if spotify.isOnFilter(webPlaylist.Name) {
				spotify.loadPlaylistTracks(&webPlaylist, playlists)
			}
		}

		// simplePlaylistPage.Offset is returning 0 instead of offset
		return offset + limit, simplePlaylistPage.Total
	}

	return 0, 0
}

func (spotify *Spotify) loadPlaylistTracks(webPlaylist *webspotify.SimplePlaylist, playlists *sconsify.Playlists) error {
	limit := 100
	offset := 0
	total := 1
	options := &webspotify.Options{Limit: &limit, Offset: &offset}

	tracks := make([]*sconsify.Track, 0)
	playlist := sconsify.InitPlaylist(string(webPlaylist.URI), webPlaylist.Name, tracks)
	playlists.AddPlaylist(playlist)

	for offset <= total {
		playlistTrackPage, err := spotify.client.GetPlaylistTracksOpt(webPlaylist.Owner.ID, webPlaylist.ID, options, "")
		if err != nil {
			return err
		}

		for _, track := range playlistTrackPage.Tracks {
			if len(track.Track.Artists) > 0 {
				webArtist := track.Track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				track := sconsify.InitWebApiTrack(string(track.Track.URI), artist, track.Track.Name, track.Track.TimeDuration().String())
				playlist.AddTrack(track)
			} else {
				infrastructure.Debugf("%v: track will be ignored as it doesn't have artist\n", track.Track.URI)
			}
		}

		offset = offset + limit
		total = playlistTrackPage.Total

		if offset > total {
			infrastructure.Debugf("%v: loaded %v from %v tracks\n", playlist.Name(), total, total)
		} else {
			infrastructure.Debugf("%v: loaded %v from %v tracks\n", playlist.Name(), offset, total)
		}
	}

	return nil
}

func (spotify *Spotify) loadWebApiCacheIfNecessary() *WebApiCache {
	if spotify.client != nil {
		return &WebApiCache{}
	}
	return spotify.loadWebApiCache()
}

func (spotify *Spotify) loadAlbums(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	if spotify.client != nil {
		if savedAlbumPage, err := spotify.client.CurrentUsersAlbumsOpt(createWebSpotifyOptions(50, playlist.Playlists())); err == nil {
			webApiCache.Albums = savedAlbumPage.Albums
		}
	}

	if webApiCache.Albums != nil {
		for _, album := range webApiCache.Albums {
			tracks := make([]*sconsify.Track, len(album.Tracks.Tracks))
			for j, track := range album.Tracks.Tracks {
				webArtist := track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				tracks[j] = sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String())
			}
			playlist.AddPlaylist(sconsify.InitSubPlaylist(string(album.URI), album.Name, tracks))
		}
		playlist.OpenFolder()
	}

}

func (spotify *Spotify) loadSongs(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	var partialSongs []webspotify.SavedTrack
	if spotify.client != nil {
		if savedTrackPage, err := spotify.client.CurrentUsersTracksOpt(createWebSpotifyOptions(50, playlist.Tracks())); err == nil {
			partialSongs = savedTrackPage.Tracks
			if webApiCache.Songs == nil {
				webApiCache.Songs = make([]webspotify.SavedTrack, 0)
			}
			webApiCache.Songs = append(webApiCache.Songs, partialSongs...)
		}
	}

	if webApiCache.Songs != nil {
		if partialSongs == nil {
			partialSongs = webApiCache.Songs
		}
		for i, track := range partialSongs {
			webApiCache.Songs[i] = track
			webArtist := track.Artists[0]
			artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
			playlist.AddTrack(sconsify.InitWebApiTrack(string(track.URI), artist, track.Name, track.TimeDuration().String()))
		}
	}
}

func (spotify *Spotify) loadNewReleases(playlist *sconsify.Playlist, webApiCache *WebApiCache) {
	if spotify.client != nil {
		if _, simplePlaylistPage, err := spotify.client.FeaturedPlaylistsOpt(&webspotify.PlaylistOptions{Options: *createWebSpotifyOptions(50, playlist.Playlists())}); err == nil {
			webApiCache.NewReleases = make([]webspotify.FullPlaylist, len(simplePlaylistPage.Playlists))
			for i, webPlaylist := range simplePlaylistPage.Playlists {
				if fullPlaylist, err := spotify.client.GetPlaylist(webPlaylist.Owner.ID, webPlaylist.ID); err == nil {
					webApiCache.NewReleases[i] = *fullPlaylist
				}
			}
		}
	}

	if webApiCache.NewReleases != nil {
		for _, fullPlaylist := range webApiCache.NewReleases {
			tracks := make([]*sconsify.Track, len(fullPlaylist.Tracks.Tracks))
			for i, track := range fullPlaylist.Tracks.Tracks {
				webArtist := track.Track.Artists[0]
				artist := sconsify.InitArtist(string(webArtist.URI), webArtist.Name)
				tracks[i] = sconsify.InitWebApiTrack(string(track.Track.URI), artist, track.Track.Name, track.Track.TimeDuration().String())
			}
			playlist.AddPlaylist(sconsify.InitSubPlaylist(string(fullPlaylist.URI), fullPlaylist.Name, tracks))
		}
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
	infrastructure.Debugf("\tTrack '%v' (%v)", track.Link().String(), track.Name())
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
