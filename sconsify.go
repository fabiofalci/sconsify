package main

import (
	"fmt"
	"log"

	"github.com/fabiofalci/sconsify/spotify"
	"github.com/jroimartin/gocui"
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

func getPlaylist(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, l)
		if err := g.SetCurrentView("msg"); err != nil {
			return err
		}
	}
	return nil
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

	if err := g.SetKeybinding("side", gocui.KeyEnter, 0, getPlaylist); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, 0, getPlaylist); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		v.Highlight = true

		fmt.Fprintln(v, "Playlist Dummy")
		if spotify.GetSession() != nil {
			playlists, _ := spotify.GetSession().Playlists()
			// for i := 0; i < playlists.Playlists(); i++ {
			// 	playlist := playlists.Playlist(i)
			// 	playlist.Wait()
			// }
			fmt.Fprintln(v, "Playlist 1 %v", playlists.Playlists())
			fmt.Fprintln(v, "Playlist 2")
			fmt.Fprintln(v, "Playlist 3")
			fmt.Fprintln(v, "Playlist 4")
		}
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}
		fmt.Fprintln(v, "Music A")
		fmt.Fprintln(v, "Music B")
		fmt.Fprintln(v, "Music C")
		fmt.Fprintln(v, "Music D")
		fmt.Fprintln(v, "Music E")
		v.Highlight = true
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var err error

	spotify.Initialise()

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
