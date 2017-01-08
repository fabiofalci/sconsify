// A lot of pieces copied from the awesome library github.com/op/go-libspotify by Örjan Persson
package spotify

import (
	"github.com/gordonklaus/portaudio"
	sp "github.com/op/go-libspotify/spotify"
	"time"
	"github.com/fabiofalci/sconsify/infrastructure"
)

type audio struct {
	format sp.AudioFormat
	frames []byte
}

type portAudio struct {
	buffer chan audio
}

type portAudioStream struct {
	device *portaudio.DeviceInfo
	stream *portaudio.Stream

	channels   int
	sampleRate int
}

func newPortAudio() *portAudio {
	return &portAudio{buffer: make(chan audio, 8)}
}

func (pa *portAudio) player(stream *portAudioStream) {
	buffer := make([]int16, 8192)
	output := buffer[:]

	previous := time.Now().Second()
	for {
		current := time.Now().Second()
		if previous != current {
			infrastructure.Debugf("Getting data %v\n", current)
			previous = current
		}
		var input audio
		select {
		case input = <-pa.buffer:
			// Initialize the audio stream based on the specification of the input format.
			err := stream.Stream(&output, input.format.Channels, input.format.SampleRate)
			if err != nil {
				panic(err)
			}

			// Decode the incoming data which is expected to be 2 channels and
			// delivered as int16 in []byte, hence we need to convert it.
			i := 0
			for i < len(input.frames) {
				j := 0
				for j < len(buffer) && i < len(input.frames) {
					buffer[j] = int16(input.frames[i]) | int16(input.frames[i+1])<<8
					j += 1
					i += 2
				}

				output = buffer[:j]
				stream.Write()
			}
		//default:
		//	infrastructure.Debugf("no message received")
		}

	}
	infrastructure.Debugf("Out!!")
}

func (s *portAudioStream) Stream(buffer *[]int16, channels int, sampleRate int) error {
	if s.stream == nil || s.channels != channels || s.sampleRate != sampleRate {
		if err := s.reset(); err != nil {
			return err
		}

		params := portaudio.HighLatencyParameters(nil, s.device)
		params.Output.Channels = channels
		params.SampleRate = float64(sampleRate)
		params.FramesPerBuffer = len(*buffer)

		stream, err := portaudio.OpenStream(params, buffer)
		if err != nil {
			return err
		}
		if err := stream.Start(); err != nil {
			stream.Close()
			return err
		}

		s.stream = stream
		s.channels = channels
		s.sampleRate = sampleRate
	}
	return nil
}

func (s *portAudioStream) reset() error {
	if s.stream != nil {
		if err := s.stream.Stop(); err != nil {
			return err
		}
		if err := s.stream.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Write pushes the data in the buffer through to PortAudio.
func (s *portAudioStream) Write() error {
	return s.stream.Write()
}

func (pa *portAudio) WriteAudio(format sp.AudioFormat, frames []byte) int {
	select {
	case pa.buffer <- audio{format, frames}:
		return len(frames)
	default:
		return 0
	}
}
