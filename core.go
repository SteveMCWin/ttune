package main

import (
	"log"
	"math"
	"math/cmplx"
	"tuner/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"

	"github.com/mjibson/go-dsp/fft"
)

const BL = 4096 * 2       // NOTE: should be loaded through settings
const SAMPLE_RATE = 44100 // NOTE: should be loaded through settings

var AudioStream *portaudio.Stream
var Buffer []float32
var Buffer64 []float64
var Window []float64

func initAutioStream() tea.Cmd {
	return func() tea.Msg {
		if AudioStream == nil {
			log.Println("Audio stream is nil as it should be")
		}
		Buffer = make([]float32, BL)
		Buffer64 = make([]float64, BL)
		Window = make([]float64, BL)
		for n := range BL {
			Window[n] = 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(n)/float64(BL-1)))
		}

		var err error
		AudioStream, err = portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BL, Buffer)
		if err != nil {
			log.Println("ERROR opening audio stream")
		} else {
			log.Println("Opened audio stream!!")
		}

		return OpenStreamMsg(AudioStream)
	}
}

func applyWindowToBuffer() {
	for i := range BL {
		Buffer64[i] *= Window[i]
	}
}

func buffTo64() {
	for i := range BL {
		Buffer64[i] = float64(Buffer[i])
	}
}

func calculateFrequency() float64 {
	comp := fft.FFTReal(Buffer64)
	var max_mag float64 = 0
	var max_mag_idx int = 0
	for i := range len(comp) / 2 {
		magnitude := cmplx.Abs(comp[i])
		if magnitude > max_mag {
			max_mag = magnitude
			max_mag_idx = i
		}
	}

	return float64(max_mag_idx) * SAMPLE_RATE / BL
}

func FrequencyToNote(freq float64) Note {
	res := Note{}
	if freq < 20 {
		return res
	}

	semitone := 12*math.Log2(freq/440.0) + 58.0
	nearestSemitone := math.Round(semitone)

	res.CentsOff = int((semitone - nearestSemitone) * 100)

	noteIndex := int(nearestSemitone-1) % tuning.NUM_SEMITONES
	for noteIndex < 0 {
		noteIndex += tuning.NUM_SEMITONES
	}
	octave := int(nearestSemitone-1) / tuning.NUM_SEMITONES

	res.Index = noteIndex
	res.Octave = octave

	return res
}

func CalculateNote() tea.Cmd {
	return func() tea.Msg {
		var note Note
		if AudioStream == nil {
			return NoteReadingMsg(note)
		}
		AudioStream.Read()
		buffTo64()
		applyWindowToBuffer()
		frequency := calculateFrequency()
		note = FrequencyToNote(frequency)

		return NoteReadingMsg(note)
	}
}

func prevNote(n Note) Note {
	res := Note {
		Index: (n.Index-1+len(tuning.NoteNames))%len(tuning.NoteNames),
		Octave: n.Octave,
	}

	if res.Index > n.Index {
		res.Octave -= 1
	}

	return res
}

func nextNote(n Note) Note {
	res := Note {
		Index: (n.Index+1)%len(tuning.NoteNames),
		Octave: n.Octave,
	}

	if res.Index < n.Index {
		res.Octave += 1
	}

	return res
}

// func old_main() {
// 	var buffer = make([]float32, BL)
// 	var window = calcHannWindow()
//
// 	stream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BL, buffer)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer stream.Close()
//
// 	stream.Start()
// 	defer stream.Stop()
//
// 	for {
// 		stream.Read()
// 		buff64 := to64(buffer)
// 		applyWindowToBuffer(buff64, window)
// 		freq := getFrequency(buff64)
// 		note, cents := frequencyToNote(freq)
// 		fmt.Printf("\rFrequency: %.2f Hz | Note: %s | Cents: %+.0f    ", freq, note, cents)
// 	}
// }
