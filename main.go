package main

import (
	"fmt"
	"math"
	"math/cmplx"

	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

const BL = 4096*2
const SAMPLE_RATE = 44100

var noteNames = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}

func calcHannWindow() [BL]float64 {
	var res [BL]float64
	for n := range BL {
		res[n] = 0.5*(1.0-math.Cos(2.0*math.Pi*float64(n)/float64(BL-1)))
	}

	return res
}

func applyWindowToBuffer(buff []float64, window [BL]float64) {
	for i := range BL {
		buff[i] *= window[i]
	}
}

func to64(buff []float32) []float64 {
	res := make([]float64, len(buff))
	for i := range len(buff) {
		res[i] = float64(buff[i])
	}

	return res
}

func getFrequency(buff []float64) float64 {
	comp := fft.FFTReal(buff)
	var max_mag float64 = 0
	var max_mag_idx int = 0
	for i := range len(comp)/2 {
		magnitude := cmplx.Abs(comp[i])
		if magnitude > max_mag {
			max_mag = magnitude
			max_mag_idx = i
		}
	}

	frequency := float64(max_mag_idx) * SAMPLE_RATE / BL
	return frequency
}

func frequencyToNote(freq float64) (string, float64) {
    if freq < 20 {
        return "---", 0
    }
    
    // Calculate semitone number
    semitone := 12*math.Log2(freq/440.0) + 58.0
    nearestSemitone := math.Round(semitone)
    
    // Calculate cents off (100 cents = 1 semitone)
    centsOff := (semitone - nearestSemitone) * 100
    
    // Get note name and octave
    noteIndex := int(nearestSemitone-1) % 12
    octave := int(nearestSemitone-1) / 12
    
    noteName := fmt.Sprintf("%s%d", noteNames[noteIndex], octave)
    
    return noteName, centsOff
}

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	var buffer = make([]float32, BL)
	var window = calcHannWindow()

	stream, err := portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BL, buffer)
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	stream.Start()
	defer stream.Stop()

	for {
		stream.Read()
		buff64 := to64(buffer)
		applyWindowToBuffer(buff64, window)
		freq := getFrequency(buff64)
		note, cents := frequencyToNote(freq)
		fmt.Printf("\rFrequency: %.2f Hz | Note: %s | Cents: %+.0f    ", freq, note, cents)
	}
}
