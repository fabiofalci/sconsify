package spotify

import (
	"errors"
	"fmt"
	"time"

	sp "github.com/fabiofalci/go-libspotify/spotify"
	"github.com/fabiofalci/sconsify/infrastructure"
	"github.com/fabiofalci/sconsify/sconsify"
	"github.com/fabiofalci/sconsify/webapi"
	"github.com/gordonklaus/portaudio"
	webspotify "github.com/zmb3/spotify"
)

type Spotify struct {
	currentTrack       *sconsify.Track
	paused             bool
	events             *sconsify.Events
	publisher          *sconsify.Publisher
	pa                 *portAudio
	session            *sp.Session
	appKey             []byte
	playlistFilter     []string
	client             *webspotify.Client
	cacheWebApiContent bool
}

type SpotifyInitConf struct {
	WebApiAuth         bool
	PlaylistFilter     string
	PreferredBitrate   string
	CacheWebApiToken   bool
	CacheWebApiContent bool
	SpotifyClientId    string
	AuthRedirectUrl    string
	OpenBrowserCommand string
}

func Initialise(initConf *SpotifyInitConf, username string, pass []byte, events *sconsify.Events, publisher *sconsify.Publisher) {
	if err := initialiseSpotify(initConf, username, pass, events, publisher); err != nil {
		fmt.Printf("Error: %v\n", err)
		publisher.ShutdownEngine()
	}
}

func initialiseSpotify(initConf *SpotifyInitConf, username string, pass []byte, events *sconsify.Events, publisher *sconsify.Publisher) error {
	spotify := &Spotify{events: events, publisher: publisher}
	spotify.setPlaylistFilter(initConf.PlaylistFilter)
	spotify.cacheWebApiContent = initConf.CacheWebApiContent
	if err := spotify.initKey(); err != nil {
		return err
	}
	pa := newPortAudio()

	cacheLocation, err := spotify.initCache()
	if err == nil {
		err = spotify.initSession(pa, cacheLocation, initConf.PreferredBitrate)
		if err == nil {
			err = spotify.login(username, pass)
			if err == nil {
				err = spotify.checkIfLoggedIn(initConf, pa)
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

func (spotify *Spotify) initSession(pa *portAudio, cacheLocation string, preferredBitrate string) error {
	var err error
	spotify.session, err = sp.NewSession(&sp.Config{
		ApplicationKey:   spotify.appKey,
		ApplicationName:  "sconsify",
		CacheLocation:    cacheLocation,
		SettingsLocation: cacheLocation,
		AudioConsumer:    pa,
	})

	switch preferredBitrate {
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

func (spotify *Spotify) checkIfLoggedIn(initConf *SpotifyInitConf, pa *portAudio) error {
	if !spotify.waitForSuccessfulConnectionStateUpdates() {
		return errors.New("Could not login")
	}
	return spotify.finishInitialisation(initConf, pa)
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

func (spotify *Spotify) finishInitialisation(initConf *SpotifyInitConf, pa *portAudio) error {
	if initConf.WebApiAuth {
		var err error
		spotify.client, err = webapi.Auth(initConf.SpotifyClientId, initConf.AuthRedirectUrl, initConf.CacheWebApiToken, initConf.OpenBrowserCommand)
		if err != nil {
			return err
		}
		if spotify.client!= nil {
			if privateUser, err := spotify.client.CurrentUser(); err == nil {
				if privateUser.ID != spotify.session.LoginUsername() {
					return errors.New("Username doesn't match with web-api authorization")
				}
			} else {
				spotify.client = nil
			}
		}
	}
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
			spotify.publisher.NextPlay()
		case <-spotify.session.PlayTokenLostUpdates():
			spotify.publisher.PlayTokenLost()
		case track := <-spotify.events.PlayUpdates():
			spotify.play(track)
		case <-spotify.events.PauseUpdates():
			spotify.pause()
		case <-spotify.events.PlayPauseToggleUpdates():
			spotify.playPauseToggle()
		case <-spotify.events.ReplayUpdates():
			spotify.play(spotify.currentTrack)
		case <-spotify.events.ShutdownSpotifyUpdates():
			spotify.shutdownSpotify()
		case query := <-spotify.events.SearchUpdates():
			spotify.search(query)
		case artist := <-spotify.events.GetArtistAlbumsUpdates():
			spotify.artistAlbums(artist)
		}
	}
}
