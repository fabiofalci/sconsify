// A lot of pieces copied from the awesome library github.com/op/go-libspotify by Örjan Persson
package spotify

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/fabiofalci/sconsify/events"
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

func newPortAudio() *portAudio {
	return &portAudio{
		buffer: make(chan *audio, 8),
	}
}

var (
	Playlists     = make(map[string]*sp.Playlist)
	currentTrack  *sp.Track
	paused        bool
	cacheLocation = "tmp"
)

func Initialise(username *string, pass *[]byte, allEvents *events.Events) {
	appKey, err := ioutil.ReadFile("spotify_appkey.key")
	if err != nil {
		log.Fatal(err)
	}

	credentials := sp.Credentials{
		Username: *username,
		Password: string(*pass),
	}

	portaudio.Initialize()
	defer portaudio.Terminate()

	pa := newPortAudio()

	deleteCache()
	session, err := sp.NewSession(&sp.Config{
		ApplicationKey:   appKey,
		ApplicationName:  "sconsify",
		CacheLocation:    cacheLocation,
		SettingsLocation: cacheLocation,
		AudioConsumer:    pa,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err = session.Login(credentials, false); err != nil {
		log.Fatal(err)
	}

	err = <-session.LoginUpdates()
	if err != nil {
		log.Fatal(err)
	}

	if checkConnectionState(session) {
		finishInitialisation(session, pa, allEvents)
	} else {
		println("Could not login")
		allEvents.Playlists <- nil
	}
}

func ShutdownSpotify() {
	deleteCache()
}

func deleteCache() {
	os.RemoveAll(cacheLocation)
}

func checkConnectionState(session *sp.Session) bool {
	timeout := make(chan bool)
	go func() {
		time.Sleep(3 * time.Second)
		timeout <- true
	}()
	loggedIn := false
	running := true
	for running {
		select {
		case <-session.ConnectionStateUpdates():
			if session.ConnectionState() == sp.ConnectionStateLoggedIn {
				running = false
				loggedIn = true
			}
		case <-timeout:
			running = false
		}
	}
	return loggedIn
}

func finishInitialisation(session *sp.Session, pa *portAudio, allEvents *events.Events) {
	playlists, _ := session.Playlists()
	playlists.Wait()
	for i := 0; i < playlists.Playlists(); i++ {
		playlist := playlists.Playlist(i)
		playlist.Wait()

		if playlists.PlaylistType(i) == sp.PlaylistTypePlaylist {
			Playlists[playlist.Name()] = playlist
		}
	}

	allEvents.Playlists <- Playlists

	go pa.player()

	go func() {
		for {
			select {
			case <-session.EndOfTrackUpdates():
				allEvents.NextPlay <- true
			}
		}
	}()
	for {
		select {
		case track := <-allEvents.ToPlay:
			Play(session, track, allEvents)
		case <-allEvents.Pause:
			Pause(session, allEvents)
		}
	}
}

func Pause(session *sp.Session, allEvents *events.Events) {
	if currentTrack != nil {
		if paused {
			Play(session, currentTrack, allEvents)
			paused = false
		} else {
			player := session.Player()
			player.Pause()

			artist := currentTrack.Artist(0)
			artist.Wait()
			allEvents.Status <- fmt.Sprintf("Paused: %v - %v [%v]", artist.Name(), currentTrack.Name(), currentTrack.Duration().String())
			paused = true
		}
	}
}

func Play(session *sp.Session, track *sp.Track, allEvents *events.Events) {
	if track.Availability() != sp.TrackAvailabilityAvailable {
		allEvents.Status <- "Not available"
		return
	}
	player := session.Player()
	if err := player.Load(track); err != nil {
		log.Fatal(err)
	}
	player.Play()
	artist := track.Artist(0)
	artist.Wait()
	currentTrack = track
	allEvents.Status <- fmt.Sprintf("Playing: %v - %v [%v]", artist.Name(), currentTrack.Name(), currentTrack.Duration().String())
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
