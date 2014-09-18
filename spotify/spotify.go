package spotify

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/fabiofalci/sconsify/events"
	"github.com/mitchellh/go-homedir"
	sp "github.com/op/go-libspotify/spotify"
)

type Spotify struct {
	currentTrack  *sp.Track
	paused        bool
	cacheLocation string
	events        *events.Events
	pa            *portAudio
	session       *sp.Session
	appKey        *[]byte
}

func Initialise(username *string, pass *[]byte, events *events.Events) {
	if err := initialiseSpotify(username, pass, events); err != nil {
		fmt.Printf("Error: %v\n", err)
		events.Shutdown()
	}
}

func initialiseSpotify(username *string, pass *[]byte, events *events.Events) error {
	spotify := &Spotify{events: events}
	err := spotify.initKey()
	if err != nil {
		return err
	}
	spotify.initAudio()
	defer portaudio.Terminate()

	err = spotify.initCache()
	if err != nil {
		return err
	}

	spotify.initSession()

	err = spotify.login(username, pass)
	if err != nil {
		return err
	}

	err = spotify.checkIfLoggedIn()
	if err != nil {
		return err
	}

	return nil
}

func (spotify *Spotify) initAudio() {
	portaudio.Initialize()
	spotify.pa = newPortAudio()
}

func (spotify *Spotify) login(username *string, pass *[]byte) error {
	credentials := sp.Credentials{Username: *username, Password: string(*pass)}
	if err := spotify.session.Login(credentials, false); err != nil {
		return err
	}

	err := <-spotify.session.LoginUpdates()
	if err != nil {
		return err
	}

	return nil
}

func (spotify *Spotify) initSession() error {
	var err error
	spotify.session, err = sp.NewSession(&sp.Config{
		ApplicationKey:   *spotify.appKey,
		ApplicationName:  "sconsify",
		CacheLocation:    spotify.cacheLocation,
		SettingsLocation: spotify.cacheLocation,
		AudioConsumer:    spotify.pa,
	})

	if err != nil {
		return err
	}
	return nil
}

func (spotify *Spotify) initKey() error {
	var err error
	spotify.appKey, err = getKey()
	if err != nil {
		return err
	}
	return nil
}

func (spotify *Spotify) initCache() error {
	spotify.initCacheLocation()
	if spotify.cacheLocation == "" {
		return errors.New("Cannot find cache dir")
	}

	spotify.deleteCache()
	return nil
}

func (spotify *Spotify) initCacheLocation() {
	dir, err := homedir.Dir()
	if err == nil {
		dir, err = homedir.Expand(dir)
		if err == nil && dir != "" {
			spotify.cacheLocation = dir + "/.sconsify/cache/"
		}
	}
}

func (spotify *Spotify) shutdownSpotify() {
	spotify.session.Logout()
	spotify.deleteCache()
	spotify.events.Shutdown()
}

func (spotify *Spotify) deleteCache() {
	if strings.HasSuffix(spotify.cacheLocation, "/.sconsify/cache/") {
		os.RemoveAll(spotify.cacheLocation)
	}
}

func (spotify *Spotify) checkIfLoggedIn() error {
	if !spotify.waitForSuccessfulConnectionStateUpdates() {
		return errors.New("Could not login")
	}
	spotify.finishInitialisation()
	return nil
}

func (spotify *Spotify) waitForSuccessfulConnectionStateUpdates() bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(9 * time.Second)
		timeout <- true
	}()
	loggedIn := false
	running := true
	for running {
		select {
		case <-spotify.session.ConnectionStateUpdates():
			if spotify.isLoggedIn() {
				running = false
				loggedIn = true
			}
		case <-timeout:
			running = false
		}
	}
	return loggedIn
}

func (spotify *Spotify) isLoggedIn() bool {
	return spotify.session.ConnectionState() == sp.ConnectionStateLoggedIn
}

func (spotify *Spotify) finishInitialisation() {
	playlists := make(map[string]*sp.Playlist)
	allPlaylists, _ := spotify.session.Playlists()
	allPlaylists.Wait()
	for i := 0; i < allPlaylists.Playlists(); i++ {
		playlist := allPlaylists.Playlist(i)
		playlist.Wait()

		if allPlaylists.PlaylistType(i) == sp.PlaylistTypePlaylist {
			playlists[playlist.Name()] = playlist
		}
	}

	spotify.events.NewPlaylist(&playlists)

	go spotify.pa.player()

	go func() {
		for {
			select {
			case <-spotify.session.EndOfTrackUpdates():
				spotify.events.NextPlay <- true
			}
		}
	}()
	for {
		select {
		case track := <-spotify.events.ToPlay:
			spotify.play(track)
		case <-spotify.events.WaitForPause():
			spotify.pause()
		case <-spotify.events.WaitForShutdown():
			spotify.shutdownSpotify()
		}
	}
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
	spotify.updateStatus("Paused", spotify.currentTrack)
	spotify.paused = true
}

func (spotify *Spotify) isPausedOrPlaying() bool {
	return spotify.currentTrack != nil
}

func (spotify *Spotify) play(track *sp.Track) {
	if !spotify.isTrackAvailable(track) {
		spotify.events.SetStatus("Not available")
		return
	}
	player := spotify.session.Player()
	if err := player.Load(track); err != nil {
		log.Fatal(err)
	}
	player.Play()

	spotify.updateStatus("Playing", track)
}

func (spotify *Spotify) isTrackAvailable(track *sp.Track) bool {
	return track.Availability() == sp.TrackAvailabilityAvailable
}

func (spotify *Spotify) updateStatus(status string, track *sp.Track) {
	spotify.currentTrack = track
	artist := track.Artist(0)
	artist.Wait()
	spotify.events.SetStatus(fmt.Sprintf("%v: %v - %v [%v]", status, artist.Name(), spotify.currentTrack.Name(), spotify.currentTrack.Duration().String()))
}
