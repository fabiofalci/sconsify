package spotify

import (
	"errors"
	"fmt"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/infrastructure"
	sp "github.com/op/go-libspotify/spotify"
	webspotify "github.com/zmb3/spotify"
)

type Spotify struct {
	currentTrack   *sconsify.Track
	paused         bool
	events         *sconsify.Events
	pa             *portAudio
	session        *sp.Session
	appKey         []byte
	playlistFilter []string
	client         *webspotify.Client
}

func Initialise(client *webspotify.Client, username string, pass []byte, events *sconsify.Events, playlistFilter *string, preferredBitrate *string) {
	if err := initialiseSpotify(client, username, pass, events, playlistFilter, preferredBitrate); err != nil {
		fmt.Printf("Error: %v\n", err)
		events.ShutdownEngine()
	}
}

func initialiseSpotify(client *webspotify.Client, username string, pass []byte, events *sconsify.Events, playlistFilter *string, preferredBitrate *string) error {
	spotify := &Spotify{events: events}
	spotify.setPlaylistFilter(*playlistFilter)
	spotify.client = client
	if err := spotify.initKey(); err != nil {
		return err
	}
	pa := newPortAudio()

	cacheLocation, err := spotify.initCache()
	if err == nil {
		err = spotify.initSession(pa, cacheLocation, preferredBitrate)
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

func (spotify *Spotify) initSession(pa *portAudio, cacheLocation string, preferredBitrate *string) error {
	var err error
	spotify.session, err = sp.NewSession(&sp.Config{
		ApplicationKey:   spotify.appKey,
		ApplicationName:  "sconsify",
		CacheLocation:    cacheLocation,
		SettingsLocation: cacheLocation,
		AudioConsumer:    pa,
	})

	switch *preferredBitrate {
	case "96k":
		spotify.session.PreferredBitrate(sp.Bitrate96k)
	case "160k":
		spotify.session.PreferredBitrate(sp.Bitrate160k)
	default:
		spotify.session.PreferredBitrate(sp.Bitrate320k)
	}

	return err
}

func (spotify *Spotify) initKey() error {
	var err error
	spotify.appKey, err = getKey()
	return err
}

func (spotify *Spotify) initCache() (string, error) {
	cacheLocation := infrastructure.GetCacheLocation()
	if cacheLocation == "" {
		return "", errors.New("Cannot find cache dir")
	}
	if err := infrastructure.DeleteCache(cacheLocation); err != nil {
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

	sconsify.InitRepository()
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
		case artist := <-spotify.events.GetArtistTopTracksUpdates():
			spotify.artistTopTrack(artist)
		}
	}
}
