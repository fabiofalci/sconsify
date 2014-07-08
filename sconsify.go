package main

import (
	"fmt"
	"log"

	"github.com/fabiofalci/sconsify/spotify"
	"github.com/jroimartin/gocui"
	sp "github.com/op/go-libspotify/spotify"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView := g.CurrentView()
	if currentView == nil || currentView.Name() == "side" {
		return g.SetCurrentView("main")
	}
	return g.SetCurrentView("side")
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getPlaylist(g *gocui.Gui, v *gocui.View) (string, error) {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	return l, nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrorQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyArrowDown, 0, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, 0, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowLeft, 0, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowRight, 0, nextView); err != nil {
		return err
	}

	// vi navigation
	if err := g.SetKeybinding("", 'j', 0, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', 0, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'h', 0, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", 'l', 0, nextView); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, 0, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', 0, quit); err != nil {
		return err
	}

	// if err := g.SetKeybinding("side", gocui.KeyEnter, 0, getPlaylist); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding("main", gocui.KeyEnter, 0, getPlaylist); err != nil {
	// 	return err
	// }

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		sideView = v
		sideView.Highlight = true

		if spotify.GetSession() != nil {
			playlists, _ := spotify.GetSession().Playlists()
			for i := 0; i < playlists.Playlists(); i++ {
				playlist := playlists.Playlist(i)
				playlist.Wait()
				fmt.Fprintln(v, playlist.Name())
			}
		}
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}

		currentPlaylist, err := getPlaylist(g, sideView)
		if err == nil && playlistsMap != nil {
			playlist := playlistsMap[currentPlaylist]

			if playlist != nil {
				playlist.Wait()
				for i := 0; i < playlist.Tracks(); i++ {
					playlistTrack := playlist.Track(i)
					track := playlistTrack.Track()
					track.Wait()
					fmt.Fprintf(v, "%v", track.Name())
					// track.Wait()
					// fmt.Fprintf(v, "%v", track.Name())
				}
			}

		}
		v.Highlight = true
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

var sideView *gocui.View

var (
	playlistsMap = make(map[string]*sp.Playlist)
)

func main() {
	var err error

	spotify.Initialise()

	if spotify.GetSession() != nil {
		playlists, _ := spotify.GetSession().Playlists()
		playlists.Wait()
		for i := 0; i < playlists.Playlists(); i++ {
			playlist := playlists.Playlist(i)
			playlist.Wait()

			if playlists.PlaylistType(i) == sp.PlaylistTypePlaylist {
				playlistsMap[playlist.Name()] = playlist
			}
		}
	}

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.ShowCursor = true

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}
