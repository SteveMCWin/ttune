package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/mjibson/go-dsp/fft"
)

const BL = 4096 * 2       // NOTE: should be loaded through settings
const SAMPLE_RATE = 44100 // NOTE: should be loaded through settings

var noteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func (m *Model) calcHannWindow() {
	m.Window = make([]float64, m.BlockLength)
	for n := range BL {
		m.Window[n] = 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(n)/float64(BL-1)))
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

func (m *Model) getFrequency() float64 {
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

	frequency := float64(max_mag_idx) * SAMPLE_RATE / BL
	return frequency
}

func (m *Model) GetNote() {
	if m.Frequency < 20 {
		return "---", 0
	}

	semitone := 12*math.Log2(freq/440.0) + 58.0
	nearestSemitone := math.Round(semitone)

	centsOff := (semitone - nearestSemitone) * 100

	noteIndex := int(nearestSemitone-1) % 12
	octave := int(nearestSemitone-1) / 12

	noteName := fmt.Sprintf("%s%d", noteNames[noteIndex], octave)

	return noteName, centsOff
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
