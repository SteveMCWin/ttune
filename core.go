package main

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"

	"github.com/mjibson/go-dsp/fft"
)

const BL = 4096 * 2       // NOTE: should be loaded through settings
const SAMPLE_RATE = 44100 // NOTE: should be loaded through settings

var noteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func (m *Model) calcHannWindow() tea.Cmd {
	return func() tea.Msg {
		m.Window = make([]float64, m.BlockLength)
		for n := range BL {
			m.Window[n] = 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(n)/float64(BL-1)))
		}

		return nil
	}
}

func (m *Model) applyWindowToBuffer() {
	for i := range BL {
		m.Buffer64[i] *= m.Window[i]
	}
}

func (m *Model) buffTo64() {
	for i := range m.BlockLength {
		m.Buffer64[i] = float64(m.Buffer[i])
	}
}

func (m *Model) calculateFrequency() {
	comp := fft.FFTReal(m.Buffer64)
	var max_mag float64 = 0
	var max_mag_idx int = 0
	for i := range len(comp) / 2 {
		magnitude := cmplx.Abs(comp[i])
		if magnitude > max_mag {
			max_mag = magnitude
			max_mag_idx = i
		}
	}

	m.Frequency = float64(max_mag_idx) * SAMPLE_RATE / BL
}

func (m *Model) GetNote() {
	if m.Frequency < 20 {
		m.Note = "---"
		m.CentsOff = 0
	}

	semitone := 12*math.Log2(m.Frequency/440.0) + 58.0
	nearestSemitone := math.Round(semitone)

	m.CentsOff = (semitone - nearestSemitone) * 100

	noteIndex := int(nearestSemitone-1) % 12
	for noteIndex < 0 {
		noteIndex += 12
	}
	octave := int(nearestSemitone-1) / 12

	m.Note = fmt.Sprintf("%s%d", noteNames[noteIndex], octave)
}

func (m *Model) DoTheThing() {
	m.AudioStream.Read()
	m.buffTo64()
	m.applyWindowToBuffer()
	m.calculateFrequency()
	m.GetNote()
}

func (m *Model) openAudioStream() tea.Cmd {
	return func() tea.Msg {
		if m.AudioStream == nil {
			log.Println("Audio stream is nil as it should be")
		}
		stream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BL, m.Buffer)
		if err != nil {
			log.Println("ERROR opening audio stream")
		} else {
			log.Println("Opened audio stream!!")
		}

		return OpenStreamMsg(stream)
	}
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
