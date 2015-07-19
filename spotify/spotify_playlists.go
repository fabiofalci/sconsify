package spotify

import (
	"strings"

	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
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
			playlists.AddPlaylist(sconsify.InitFolder(strconv.Itoa(int(folder.Id())), folder.Name(), folderPlaylists))
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

	spotify.events.NewPlaylist(playlists)
	return nil
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
