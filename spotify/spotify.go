package spotify

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"code.google.com/p/portaudio-go/portaudio"
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
	Playlists = make(map[string]*sp.Playlist)
)

var statusChannel chan string
var nextPlayChannel chan string

func Initialise(initialised chan string, toPlay chan sp.Track, nextPlay chan string, status chan string) {
	appKey, err := ioutil.ReadFile("spotify_appkey.key")
	if err != nil {
		log.Fatal(err)
	}

	credentials := sp.Credentials{
		Username: "fabiofalci",
		Password: os.Getenv("SPOTIFY_PASSWORD"),
	}

	statusChannel = status
	nextPlayChannel = status

	portaudio.Initialize()
	defer portaudio.Terminate()

	pa := newPortAudio()
	go pa.player()

	session, err := sp.NewSession(&sp.Config{
		ApplicationKey:   appKey,
		ApplicationName:  "testing",
		CacheLocation:    "tmp",
		SettingsLocation: "tmp",
		AudioConsumer:    pa,
	})

	if err != nil {
		log.Fatal(err)
	}

	if err = session.Login(credentials, false); err != nil {
		log.Fatal(err)
	}

	select {
	case err := <-session.LoginUpdates():
		if err != nil {
			log.Fatal(err)
		}
	}

	if session != nil {
		playlists, _ := session.Playlists()
		playlists.Wait()
		for i := 0; i < playlists.Playlists(); i++ {
			playlist := playlists.Playlist(i)
			playlist.Wait()

			if playlists.PlaylistType(i) == sp.PlaylistTypePlaylist {
				Playlists[playlist.Name()] = playlist
			}
		}
	}

	initialised <- ""

	go func() {
		for {
			select {
			case <-session.EndOfTrackUpdates():
				nextPlay <- ""
			}
		}
	}()

	for {
		select {
		case track := <-toPlay:
			Play(session, &track)
		}
	}
}

func Play(session *sp.Session, track *sp.Track) {
	if track.Availability() != sp.TrackAvailabilityAvailable {
		statusChannel <- "Not available"
		return
	}
	player := session.Player()
	if err := player.Load(track); err != nil {
		println("error")
		log.Fatal(err)
	}
	player.Play()
	artist := track.Artist(0)
	artist.Wait()
	statusChannel <- fmt.Sprintf("Playing: %v - %v [%v]", artist.Name(), track.Name(), track.Duration().String())
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
			// for audio := range buffer {
			if len(audio.frames) != 2048*2*2 {
				// panic("unexpected")
				// don't know if it's a panic or track just ended
				// nextPlayChannel <- ""
				break
			}

			j := 0
			for i := 0; i < len(audio.frames); i += 2 {
				out[j] = int16(audio.frames[i]) | int16(audio.frames[i+1])<<8
				j++
			}

			stream.Write()
			// }
		}
		// time.Sleep(1 * time.Second)
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
