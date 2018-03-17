package simple

import (
	"github.com/jroimartin/gocui"
)

type Player interface {
	Play()
	Pause()
}

type RegularPlayer struct{}

func (p *RegularPlayer) Pause() {
	publisher.PlayPauseToggle()
}

func (p *RegularPlayer) Play() {
	if playlist, trackIndex := gui.getSelectedPlaylistAndTrack(); playlist != nil {
		if trackIndex == -1 {
			if playlist.IsOnDemand() {
				go func() {
					playlist.ExecuteLoad()
					gui.g.Update(func(g *gocui.Gui) error {
						gui.updatePlaylistsView()
						cx, cy := gui.tracksView.Cursor()
						ox, oy := gui.tracksView.Origin()
						gui.updateTracksView()
						gui.tracksView.SetCursor(cx, cy)
						gui.tracksView.SetOrigin(ox, oy)
						return nil
					})
				}()
			}
		} else {
			track := playlist.Track(trackIndex)
			playlists.SetCurrents(playlist.Name(), trackIndex)
			publisher.Play(track)
		}
	}
}
