// A lot of pieces copied from the awesome library github.com/op/go-libspotify by Örjan Persson
package spotify

import (
	sp "github.com/fabiofalci/go-libspotify/spotify"
	"github.com/gordonklaus/portaudio"
)

type audio struct {
	format sp.AudioFormat
	frames []byte
}

type portAudio struct {
	buffer chan *audio
}

func newPortAudio() *portAudio {
	return &portAudio{buffer: make(chan *audio, 8)}
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
