package spotify

import (
	"strings"

	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
)

func (spotify *Spotify) initPlaylist() error {
	playlists := sconsify.InitPlaylists()

	allPlaylists, err := spotify.session.Playlists()
	if err != nil {
		return err
	}
	allPlaylists.Wait()
	for i := 0; i < allPlaylists.Playlists(); i++ {
		if allPlaylists.PlaylistType(i) != sp.PlaylistTypePlaylist {
			continue
		}
		playlist := allPlaylists.Playlist(i)
		playlist.Wait()

		if spotify.canAddPlaylist(playlist, allPlaylists.PlaylistType(i)) {
			tracks := make([]*sconsify.Track, playlist.Tracks())
			for i := 0; i < playlist.Tracks(); i++ {
				playlistTrack := playlist.Track(i)
				playlistTrack.Track().Wait()
				playlistTrack.Track().Artist(0).Wait()
				tracks[i] = sconsify.ToSconsifyTrack(playlistTrack.Track())
			}
			id := playlist.Link().String()
			playlists.AddPlaylist(sconsify.InitPlaylist(id, playlist.Name(), tracks))
		}
	}

	spotify.events.NewPlaylist(playlists)
	return nil
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
