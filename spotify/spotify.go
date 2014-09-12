// A lot of pieces copied from the awesome library github.com/op/go-libspotify by Örjan Persson
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

type audio struct {
	format sp.AudioFormat
	frames []byte
}

type audio2 struct {
	format sp.AudioFormat
	frames []int16
}

type portAudio struct {
	buffer chan *audio
}

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
		events.Shutdown <- true
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

func newPortAudio() *portAudio {
	return &portAudio{buffer: make(chan *audio, 8)}
}

func (spotify *Spotify) initAudio() {
	portaudio.Initialize()

	spotify.pa = newPortAudio()
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
	spotify.events.Shutdown <- true
}

func (spotify *Spotify) deleteCache() {
	if strings.HasSuffix(spotify.cacheLocation, "/.sconsify/cache/") {
		os.RemoveAll(spotify.cacheLocation)
	}
}

func (spotify *Spotify) checkIfLoggedIn() error {
	if spotify.checkConnectionState() {
		spotify.finishInitialisation()
	} else {
		spotify.events.NewPlaylist(nil)
		return errors.New("Could not login")
	}
	return nil
}

func (spotify *Spotify) checkConnectionState() bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		timeout <- true
	}()
	loggedIn := false
	running := true
	for running {
		select {
		case <-spotify.session.ConnectionStateUpdates():
			if spotify.session.ConnectionState() == sp.ConnectionStateLoggedIn {
				running = false
				loggedIn = true
			}
		case <-timeout:
			running = false
		}
	}
	return loggedIn
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
			spotify.Play(track)
		case <-spotify.events.Pause:
			spotify.Pause()
		case <-spotify.events.Shutdown:
			spotify.shutdownSpotify()
		}
	}
}

func (spotify *Spotify) Pause() {
	if spotify.currentTrack != nil {
		if spotify.paused {
			spotify.Play(spotify.currentTrack)
			spotify.paused = false
		} else {
			player := spotify.session.Player()
			player.Pause()

			artist := spotify.currentTrack.Artist(0)
			artist.Wait()
			spotify.events.SetStatus(fmt.Sprintf("Paused: %v - %v [%v]", artist.Name(), spotify.currentTrack.Name(), spotify.currentTrack.Duration().String()))
			spotify.paused = true
		}
	}
}

func (spotify *Spotify) Play(track *sp.Track) {
	if track.Availability() != sp.TrackAvailabilityAvailable {
		spotify.events.SetStatus("Not available")
		return
	}
	player := spotify.session.Player()
	if err := player.Load(track); err != nil {
		log.Fatal(err)
	}
	player.Play()
	artist := track.Artist(0)
	artist.Wait()
	spotify.currentTrack = track
	spotify.events.SetStatus(fmt.Sprintf("Playing: %v - %v [%v]", artist.Name(), spotify.currentTrack.Name(), spotify.currentTrack.Duration().String()))
}

func (pa *portAudio) player() {
	out := make([]int16, 2048*2)

	stream, err := portaudio.OpenDefaultStream(
		0,
		2,     // audio.format.Channels,
		44100, // float64(audio.format.SampleRate),
		len(out),
		&out,
	)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	for {
		// Decode the incoming data which is expected to be 2 channels and
		// delivered as int16 in []byte, hence we need to convert it.

		select {
		case audio := <-pa.buffer:
			if len(audio.frames) != 2048*2*2 {
				// panic("unexpected")
				// don't know if it's a panic or track just ended
				break
			}

			j := 0
			for i := 0; i < len(audio.frames); i += 2 {
				out[j] = int16(audio.frames[i]) | int16(audio.frames[i+1])<<8
				j++
			}

			stream.Write()
		}
	}
}

func (pa *portAudio) WriteAudio(format sp.AudioFormat, frames []byte) int {
	audio := &audio{format, frames}

	if len(frames) == 0 {
		return 0
	}

	select {
	case pa.buffer <- audio:
		return len(frames)
	default:
		return 0
	}
}
