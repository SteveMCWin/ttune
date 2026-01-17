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

const MIN_FREQUENCY = 70
const MAX_FREQUENCY = 1500

const MIN_AMPLITUDE_THRESHOLD = 0.01
const MIN_CLARITY_RATIO = 2.0

const MIN_BIN = MIN_FREQUENCY * BL / SAMPLE_RATE
const MAX_BIN = MAX_FREQUENCY * BL / SAMPLE_RATE

var frequencyHistory []float64

const HISTORY_SIZE = 3

var AudioStream *portaudio.Stream
var Buffer []float32
var Buffer64 []float64
var Window []float64

func freq_to_octave(freq float64) int {
	var i int
	var f float64
	for i, f = range tuning.OctaveEnds {
		if f > freq {
			return i
		}
	}

	return i
}

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

		err = AudioStream.Start()
		if err != nil {
			log.Println("ERROR starting audio stream")
		} else {
			log.Println("Started audio stream!!")
		}

		return OpenStreamMsg(AudioStream)
	}
}

func closeAudioStream() tea.Cmd {
	return func() tea.Msg {
		if AudioStream == nil {
			log.Println("Tried to close nil stream!!")
			return nil
		}

		AudioStream.Stop()
		err := AudioStream.Close()
		if err != nil {
			log.Println("FAILED TO CLOSE THE AUDIO STREAM")
		}
		log.Println("Closed audio stream successfully")

		return nil
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

func checkSignalStrength() bool {
	var sumSquares float64
	for i := range BL {
		sumSquares += Buffer64[i] * Buffer64[i]
	}
	rms := math.Sqrt(sumSquares / float64(BL))
	return rms > MIN_AMPLITUDE_THRESHOLD
}

func calculateFrequency() (float64, bool) {
	comp := fft.FFTReal(Buffer64)

	// Convert frequency to bin indices
	minBin := MIN_FREQUENCY * BL / SAMPLE_RATE
	maxBin := MAX_FREQUENCY * BL / SAMPLE_RATE

	if maxBin > len(comp)/2 {
		log.Println("Max bin is larger than len(comp)")
		maxBin = len(comp) / 2
	}

	// Find peak magnitude within our frequency range
	var maxMag float64 = 0
	var maxMagIdx int = 0
	var sumMag float64 = 0
	var count int = 0

	for i := minBin; i < maxBin; i++ {
		magnitude := cmplx.Abs(comp[i])
		sumMag += magnitude
		count++

		if magnitude > maxMag {
			maxMag = magnitude
			maxMagIdx = i
		}
	}

	// Check if peak is clear enough (stands out from noise)
	avgMag := sumMag / float64(count)
	clarityRatio := maxMag / avgMag

	if clarityRatio < MIN_CLARITY_RATIO {
		log.Printf("Signal not clear enough: clarity ratio %.2f (need %.2f)\n",
			clarityRatio, MIN_CLARITY_RATIO)
		return 0, false
	}

	// Parabolic interpolation for sub-bin accuracy
	// This significantly improves frequency resolution
	freq := parabolicInterpolation(comp, maxMagIdx)

	log.Printf("Peak at bin %d, interpolated freq: %.2f Hz, clarity: %.2f\n",
		maxMagIdx, freq, clarityRatio)

	return freq, true
}

func parabolicInterpolation(spectrum []complex128, peakBin int) float64 {
	if peakBin <= 0 || peakBin >= len(spectrum)-1 {
		return float64(peakBin) * SAMPLE_RATE / BL
	}

	// Get magnitudes of peak and neighbors
	alpha := cmplx.Abs(spectrum[peakBin-1])
	beta := cmplx.Abs(spectrum[peakBin])
	gamma := cmplx.Abs(spectrum[peakBin+1])

	// Parabolic interpolation formula
	p := 0.5 * (alpha - gamma) / (alpha - 2*beta + gamma)

	// Interpolated bin position
	interpolatedBin := float64(peakBin) + p

	// Convert to frequency
	return interpolatedBin * SAMPLE_RATE / BL
}

func medianFilter(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy and sort it
	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort (fine for small arrays)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Return median
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

func smoothFrequency(freq float64) float64 {
	frequencyHistory = append(frequencyHistory, freq)
	if len(frequencyHistory) > HISTORY_SIZE {
		frequencyHistory = frequencyHistory[1:]
	}

	return medianFilter(frequencyHistory)
}

func FrequencyToNote(freq float64) Note {
	res := Note{}
	if freq < MIN_FREQUENCY {
		return res
	}

	semitone := 12*math.Log2(freq/440.0) + 58.0
	nearestSemitone := math.Round(semitone)

	res.CentsOff = int((semitone - nearestSemitone) * 100)

	noteIndex := int(nearestSemitone-1) % tuning.NUM_SEMITONES
	for noteIndex < 0 {
		noteIndex += tuning.NUM_SEMITONES
	}
	octave := freq_to_octave(freq)

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

		err := AudioStream.Read()
		if err != nil {
			log.Println("Error reading from audio stream:", err)
			return NoteReadingMsg(note)
		}

		buffTo64()

		// Check signal strength before processing
		if !checkSignalStrength() {
			log.Println("Signal too weak, skipping...")
			// Return empty note to indicate no clear signal
			return NoteReadingMsg(note)
		}

		applyWindowToBuffer()
		frequency, isValid := calculateFrequency()

		if !isValid {
			log.Println("Invalid frequency detection")
			return NoteReadingMsg(note)
		}

		// Apply smoothing
		smoothedFreq := smoothFrequency(frequency)
		log.Printf("Raw freq: %.2f Hz, Smoothed: %.2f Hz\n", frequency, smoothedFreq)

		note = FrequencyToNote(smoothedFreq)

		return NoteReadingMsg(note)
	}
}

func prevNote(n Note) Note {
	res := Note{
		Index:  (n.Index - 1 + len(tuning.NoteNames)) % len(tuning.NoteNames),
		Octave: n.Octave,
	}

	if res.Index > n.Index {
		res.Octave -= 1
	}

	return res
}

func nextNote(n Note) Note {
	res := Note{
		Index:  (n.Index + 1) % len(tuning.NoteNames),
		Octave: n.Octave,
	}

	if res.Index < n.Index {
		res.Octave += 1
	}

	return res
}
