package ui

import (
	"github.com/fabiofalci/sconsify/sconsify"
)

type Player interface {
	Play()
	Pause()
}

type RegularPlayer struct{}

type PersistStatePlayer struct {
	previousPlayingTrack    *sconsify.Track
	previousPlayingPlaylist string
}

func (p *RegularPlayer) Pause() {
	events.Pause()
}

func (p *RegularPlayer) Play() {
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		if trackIndex == -1 {
			if playlist.IsOnDemand() {
				playlist.ExecuteLoad()
				gui.updatePlaylistsView()
				gui.updateTracksView()
			}
		} else {
			track := playlist.Track(trackIndex)
			playlists.SetCurrents(playlist.Name(), trackIndex)
			events.Play(track)
		}
	}
}

func (p *PersistStatePlayer) Pause() {
	if playlist := playlists.GetByURI(p.previousPlayingPlaylist); playlist != nil {
		if currentIndexTrack := playlist.IndexByUri(p.previousPlayingTrack.URI); currentIndexTrack != -1 {
			playlists.SetCurrents(playlist.Name(), currentIndexTrack)
			events.Play(p.previousPlayingTrack)
		}
	}
	player = &RegularPlayer{}
}

func (p *PersistStatePlayer) Play() {
	player = &RegularPlayer{}
	player.Play()
}
