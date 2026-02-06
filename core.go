package main

import (
	"log"
	"math"
	"ttune/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"
)

const BL = 4096 * 2       // NOTE: should be loaded through settings
const SAMPLE_RATE = 44100 // NOTE: should be loaded through settings

const MIN_FREQUENCY = 70
const MAX_FREQUENCY = 1500

const MIN_AMPLITUDE_THRESHOLD = 0.01
const YIN_THRESHOLD = 0.10 // Lower = stricter detection, reduces harmonic errors

// YIN power threshold - helps filter weak detections
const YIN_POWER_THRESHOLD = 0.85

const MIN_BIN = MIN_FREQUENCY * BL / SAMPLE_RATE
const MAX_BIN = MAX_FREQUENCY * BL / SAMPLE_RATE

var frequencyHistory []float64

const HISTORY_SIZE = 5 // Increased for better smoothing

var AudioStream *portaudio.Stream
var Buffer []float32
var Buffer64 []float64

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

		var err error
		AudioStream, err = portaudio.OpenDefaultStream(1, 0, SAMPLE_RATE, BL, Buffer)
		if err != nil {
			log.Println("ERROR opening audio stream")
		}

		err = AudioStream.Start()
		if err != nil {
			log.Println("ERROR starting audio stream")
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

		return nil
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

// YIN Algorithm Implementation
func yinDifference(buffer []float64, tauMax int) []float64 {
	diff := make([]float64, tauMax)

	for tau := range tauMax {
		for i := 0; i < len(buffer)-tauMax; i++ {
			delta := buffer[i] - buffer[i+tau]
			diff[tau] += delta * delta
		}
	}

	return diff
}

func yinCumulativeMeanNormalizedDifference(diff []float64) []float64 {
	cmndf := make([]float64, len(diff))
	cmndf[0] = 1.0

	runningSum := 0.0
	for tau := 1; tau < len(diff); tau++ {
		runningSum += diff[tau]
		if runningSum == 0 {
			cmndf[tau] = 1.0
		} else {
			cmndf[tau] = diff[tau] / (runningSum / float64(tau))
		}
	}

	return cmndf
}

func yinAbsoluteThreshold(cmndf []float64, threshold float64, tauMin int) int {
	tau := tauMin

	// Find first tau where cmndf drops below threshold
	for tau < len(cmndf) {
		if cmndf[tau] < threshold {
			// Find the minimum in the valley
			for tau+1 < len(cmndf) && cmndf[tau+1] < cmndf[tau] {
				tau++
			}

			// Additional check: verify this is a strong period
			// by checking the power at this tau
			if cmndf[tau] < YIN_POWER_THRESHOLD {
				return tau
			}
		}
		tau++
	}

	// No period found below threshold, return minimum cmndf
	minTau := tauMin
	minVal := cmndf[tauMin]
	for tau := tauMin + 1; tau < len(cmndf); tau++ {
		if cmndf[tau] < minVal {
			minVal = cmndf[tau]
			minTau = tau
		}
	}

	return minTau
}

func yinParabolicInterpolation(cmndf []float64, tau int) float64 {
	if tau < 1 || tau >= len(cmndf)-1 {
		return float64(tau)
	}

	s0 := cmndf[tau-1]
	s1 := cmndf[tau]
	s2 := cmndf[tau+1]

	// Parabolic interpolation
	adjustment := (s2 - s0) / (2 * (2*s1 - s2 - s0))

	return float64(tau) + adjustment
}

func calculateFrequencyYIN() (float64, bool) {
	// Calculate tau range based on frequency range
	tauMin := int(math.Round(float64(SAMPLE_RATE) / MAX_FREQUENCY))
	tauMax := int(math.Round(float64(SAMPLE_RATE) / MIN_FREQUENCY))

	tauMax = min(tauMax, len(Buffer64))

	// Step 1: Calculate difference function
	diff := yinDifference(Buffer64, tauMax)

	// Step 2: Cumulative mean normalized difference
	cmndf := yinCumulativeMeanNormalizedDifference(diff)

	// Step 3: Absolute threshold
	tau := yinAbsoluteThreshold(cmndf, YIN_THRESHOLD, tauMin)

	// Check if we found a valid period
	if tau == 0 || cmndf[tau] >= 1.0 {
		log.Printf("YIN: No clear pitch detected (cmndf[%d] = %.3f)\n", tau, cmndf[tau])
		return 0, false
	}

	// Step 4: Parabolic interpolation for better accuracy
	betterTau := yinParabolicInterpolation(cmndf, tau)

	// Convert tau to frequency
	frequency := float64(SAMPLE_RATE) / betterTau

	log.Printf("YIN: tau=%d, interpolated=%.2f, freq=%.2f Hz, confidence=%.3f\n",
		tau, betterTau, frequency, 1.0-cmndf[tau])

	return frequency, true
}

func medianFilter(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy and sort it
	sorted := make([]float64, len(values))
	copy(sorted, values)

	// Simple bubble sort (fine for small arrays)
	for i := range len(sorted) {
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
	// Check if frequency is a likely harmonic error
	// If new frequency is very different, clear history
	if len(frequencyHistory) > 0 {
		lastFreq := frequencyHistory[len(frequencyHistory)-1]
		ratio := freq / lastFreq

		// If jump is near a harmonic ratio (2x, 3x, 0.5x, 0.33x), it's suspicious
		// Allow small variations, but reset on large jumps
		if ratio > 1.8 || ratio < 0.55 {
			// Large jump detected, might be harmonic error
			// Only reset if the jump is really significant
			log.Printf("Frequency jump detected: %.2f -> %.2f (ratio: %.2f)\n", lastFreq, freq, ratio)
			frequencyHistory = []float64{freq}
			return freq
		}
	}

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

	semitone := float64(len(tuning.NoteNames))*math.Log2(freq/440.0) + 58.0
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
			return NoteReadingMsg(note)
		}

		// Use YIN algorithm instead of FFT
		frequency, isValid := calculateFrequencyYIN()

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
		res.Octave = res.Octave - 1
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
