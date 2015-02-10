package spotify

import (
	"errors"
	"fmt"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/fabiofalci/sconsify/sconsify"
	sp "github.com/op/go-libspotify/spotify"
)

type Spotify struct {
	currentTrack   *sconsify.Track
	paused         bool
	events         *sconsify.Events
	pa             *portAudio
	session        *sp.Session
	appKey         []byte
	playlistFilter []string
}

func Initialise(username string, pass []byte, events *sconsify.Events, playlistFilter *string) {
	if err := initialiseSpotify(username, pass, events, playlistFilter); err != nil {
		fmt.Printf("Error: %v\n", err)
		events.ShutdownEngine()
	}
}

func initialiseSpotify(username string, pass []byte, events *sconsify.Events, playlistFilter *string) error {
	spotify := &Spotify{events: events}
	spotify.setPlaylistFilter(*playlistFilter)
	if err := spotify.initKey(); err != nil {
		return err
	}
	pa := newPortAudio()

	cacheLocation, err := spotify.initCache()
	if err == nil {
		err = spotify.initSession(pa, cacheLocation)
		if err == nil {
			err = spotify.login(username, pass)
			if err == nil {
				err = spotify.checkIfLoggedIn(pa)
			}
		}
	}

	return err
}

func (spotify *Spotify) login(username string, pass []byte) error {
	credentials := sp.Credentials{Username: username, Password: string(pass)}
	if err := spotify.session.Login(credentials, false); err != nil {
		return err
	}

	return <-spotify.session.LoggedInUpdates()
}

func (spotify *Spotify) initSession(pa *portAudio, cacheLocation string) error {
	var err error
	spotify.session, err = sp.NewSession(&sp.Config{
		ApplicationKey:   spotify.appKey,
		ApplicationName:  "sconsify",
		CacheLocation:    cacheLocation,
		SettingsLocation: cacheLocation,
		AudioConsumer:    pa,
	})

	return err
}

func (spotify *Spotify) initKey() error {
	var err error
	spotify.appKey, err = getKey()
	return err
}

func (spotify *Spotify) initCache() (string, error) {
	cacheLocation := sconsify.GetCacheLocation()
	if cacheLocation == "" {
		return "", errors.New("Cannot find cache dir")
	}
	if err := sconsify.DeleteCache(cacheLocation); err != nil {
		return "", err
	}
	return cacheLocation, nil
}

func (spotify *Spotify) checkIfLoggedIn(pa *portAudio) error {
	if !spotify.waitForSuccessfulConnectionStateUpdates() {
		return errors.New("Could not login")
	}
	return spotify.finishInitialisation(pa)
}

func (spotify *Spotify) waitForSuccessfulConnectionStateUpdates() bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(9 * time.Second)
		timeout <- true
	}()
	for {
		select {
		case <-spotify.session.ConnectionStateUpdates():
			return spotify.isLoggedIn()
		case <-timeout:
			return false
		}
	}
	return false
}

func (spotify *Spotify) isLoggedIn() bool {
	return spotify.session.ConnectionState() == sp.ConnectionStateLoggedIn
}

func (spotify *Spotify) finishInitialisation(pa *portAudio) error {
	// init audio could happen after initPlaylist but this logs to output therefore
	// the screen isn't built properly
	portaudio.Initialize()
	go pa.player()
	defer portaudio.Terminate()

	if err := spotify.initPlaylist(); err != nil {
		return err
	}

	spotify.waitForEvents()
	return nil
}

func (spotify *Spotify) waitForEvents() {
	for {
		select {
		case <-spotify.session.EndOfTrackUpdates():
			spotify.events.NextPlay()
		case <-spotify.session.PlayTokenLostUpdates():
			spotify.events.PlayTokenLost()
		case track := <-spotify.events.PlayUpdates():
			spotify.play(track)
		case <-spotify.events.PauseUpdates():
			spotify.pause()
		case <-spotify.events.ReplayUpdates():
			spotify.playCurrentTrack()
		case <-spotify.events.ShutdownSpotifyUpdates():
			spotify.shutdownSpotify()
		case query := <-spotify.events.SearchUpdates():
			spotify.search(query)
		}
	}
}
