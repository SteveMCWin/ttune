package tuning

import (
	"log"
	"math"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"
)


type NoteReading struct {
	Index    int
	Octave   int
	CentsOff int
}

type NoteReadingMsg NoteReading
type OpenStreamMsg *portaudio.Stream

type PitchDetector struct {
	BufferLength int
	SampleRate int

	MinFrequency int
	MaxFrequency int

	MinAmplitudeThreshold float64
	YinCandidateThreshold float64
	YinValidityCeiling float64

	HistorySize int

	AudioStream *portaudio.Stream

    Buffer []float32
    Buffer64 []float64

	frequencyHistory []float64
	minBin float64
	maxBin float64
}

// var BL = 4096 * 2       // NOTE: should be loaded through settings
// var SAMPLE_RATE = 44100 // NOTE: should be loaded through settings
//
// var MIN_FREQUENCY = 70
// var MAX_FREQUENCY = 1500
//
// var MIN_AMPLITUDE_THRESHOLD = 0.01
// var YIN_CANDIDATE_THRESHOLD = 0.10 // Lower = stricter detection, reduces harmonic errors
// var YIN_VALIDITY_CEILING = 0.85 // YIN power threshold - helps filter weak detections
//
// var MIN_BIN = MIN_FREQUENCY * BL / SAMPLE_RATE
// var MAX_BIN = MAX_FREQUENCY * BL / SAMPLE_RATE
//
// var frequencyHistory []float64
//
// var HISTORY_SIZE = 5 // Increased for better smoothing
//
// var AudioStream *portaudio.Stream
// var Buffer []float32
// var Buffer64 []float64

func freq_to_octave(freq float64) int {
	var i int
	var f float64
	for i, f = range OctaveEnds {
		if f > freq {
			return i
		}
	}

	return i
}

func (pd *PitchDetector) InitAudioStream() tea.Cmd {
	return func() tea.Msg {
		if pd.AudioStream == nil {
			log.Println("Audio stream is nil as it should be")
		}
		pd.Buffer = make([]float32, pd.BufferLength)
		pd.Buffer64 = make([]float64, pd.BufferLength)

		var err error
		pd.AudioStream, err = portaudio.OpenDefaultStream(1, 0, float64(pd.SampleRate), pd.BufferLength, pd.Buffer)
		if err != nil {
			log.Println("ERROR opening audio stream")
		}

		err = pd.AudioStream.Start()
		if err != nil {
			log.Println("ERROR starting audio stream")
		}

		return OpenStreamMsg(pd.AudioStream)
	}
}

func (pd *PitchDetector) CloseAudioStream() tea.Cmd {
	return func() tea.Msg {
		if pd.AudioStream == nil {
			log.Println("Tried to close nil stream!!")
			return nil
		}

		err := pd.AudioStream.Stop()
		if err != nil {
			log.Println("ERROR stopping audio stream")
		}
		err = pd.AudioStream.Close()
		if err != nil {
			log.Println("FAILED TO CLOSE THE AUDIO STREAM")
		}

		return nil
	}
}

func (pd *PitchDetector) buffTo64() {
	for i := range pd.BufferLength {
		pd.Buffer64[i] = float64(pd.Buffer[i])
	}
}

func (pd *PitchDetector) checkSignalStrength() bool {
	var sumSquares float64
	for i := range pd.BufferLength {
		sumSquares += pd.Buffer64[i] * pd.Buffer64[i]
	}
	rms := math.Sqrt(sumSquares / float64(pd.BufferLength))
	return rms > pd.MinAmplitudeThreshold
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

func (pd *PitchDetector) yinAbsoluteThreshold(cmndf []float64, threshold float64, tauMin int) int {
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
			if cmndf[tau] < pd.YinValidityCeiling {
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

func (pd *PitchDetector) calculateFrequencyYIN() (float64, bool) {
	// Calculate tau range based on frequency range
	tauMin := int(math.Round(float64(pd.SampleRate) / float64(pd.MaxFrequency)))
	tauMax := int(math.Round(float64(pd.SampleRate) / float64(pd.MinFrequency)))

	tauMax = min(tauMax, len(pd.Buffer64))

	// Step 1: Calculate difference function
	diff := yinDifference(pd.Buffer64, tauMax)

	// Step 2: Cumulative mean normalized difference
	cmndf := yinCumulativeMeanNormalizedDifference(diff)

	// Step 3: Absolute threshold
	tau := pd.yinAbsoluteThreshold(cmndf, pd.YinCandidateThreshold, tauMin)

	// Check if we found a valid period
	if tau == 0 || cmndf[tau] >= 1.0 {
		// log.Printf("YIN: No clear pitch detected (cmndf[%d] = %.3f)\n", tau, cmndf[tau])
		return 0, false
	}

	// Step 4: Parabolic interpolation for better accuracy
	betterTau := yinParabolicInterpolation(cmndf, tau)

	// Convert tau to frequency
	frequency := float64(pd.SampleRate) / betterTau

	// log.Printf("YIN: tau=%d, interpolated=%.2f, freq=%.2f Hz, confidence=%.3f\n",
	// 	tau, betterTau, frequency, 1.0-cmndf[tau])

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

func (pd *PitchDetector) smoothFrequency(freq float64) float64 {
	// Check if frequency is a likely harmonic error
	// If new frequency is very different, clear history
	if len(pd.frequencyHistory) > 0 {
		lastFreq := pd.frequencyHistory[len(pd.frequencyHistory)-1]
		ratio := freq / lastFreq

		// If jump is near a harmonic ratio (2x, 3x, 0.5x, 0.33x), it's suspicious
		// Allow small variations, but reset on large jumps
		if ratio > 1.8 || ratio < 0.55 {
			// Large jump detected, might be harmonic error
			// Only reset if the jump is really significant
			// log.Printf("Frequency jump detected: %.2f -> %.2f (ratio: %.2f)\n", lastFreq, freq, ratio)
			pd.frequencyHistory = []float64{freq}
			return freq
		}
	}

	pd.frequencyHistory = append(pd.frequencyHistory, freq)
	if len(pd.frequencyHistory) > pd.HistorySize {
		pd.frequencyHistory = pd.frequencyHistory[1:]
	}

	return medianFilter(pd.frequencyHistory)
}

func (pd *PitchDetector) FrequencyToNote(freq float64) NoteReading {
	res := NoteReading{}
	if freq < float64(pd.MinFrequency) {
		return res
	}

	semitone := float64(len(NoteNames))*math.Log2(freq/440.0) + 58.0
	nearestSemitone := math.Round(semitone)

	res.CentsOff = int((semitone - nearestSemitone) * 100)

	noteIndex := int(nearestSemitone-1) % NUM_SEMITONES
	for noteIndex < 0 {
		noteIndex += NUM_SEMITONES
	}
	octave := freq_to_octave(freq)

	res.Index = noteIndex
	res.Octave = octave

	return res
}

func SetPitchDetectionParameters()  {
}

func (pd *PitchDetector) CalculateNote() tea.Cmd {
	return func() tea.Msg {
		var note NoteReading
		if pd.AudioStream == nil {
			return NoteReadingMsg(note)
		}

		err := pd.AudioStream.Read()
		if err != nil {
			log.Println("Error reading from audio stream:", err)
			return NoteReadingMsg(note)
		}

		pd.buffTo64()

		// Check signal strength before processing
		if !pd.checkSignalStrength() {
			// log.Println("Signal too weak, skipping...")
			return NoteReadingMsg(note)
		}

		// Use YIN algorithm instead of FFT
		frequency, isValid := pd.calculateFrequencyYIN()

		if !isValid {
			// log.Println("Invalid frequency detection")
			return NoteReadingMsg(note)
		}

		// Apply smoothing
		smoothedFreq := pd.smoothFrequency(frequency)
		// log.Printf("Raw freq: %.2f Hz, Smoothed: %.2f Hz\n", frequency, smoothedFreq)

		note = pd.FrequencyToNote(smoothedFreq)

		return NoteReadingMsg(note)
	}
}

func PrevNote(n NoteReading) NoteReading {
	res := NoteReading{
		Index:  (n.Index - 1 + len(NoteNames)) % len(NoteNames),
		Octave: n.Octave,
	}

	if res.Index > n.Index {
		res.Octave = res.Octave - 1
	}

	return res
}

func NextNote(n NoteReading) NoteReading {
	res := NoteReading{
		Index:  (n.Index + 1) % len(NoteNames),
		Octave: n.Octave,
	}

	if res.Index < n.Index {
		res.Octave += 1
	}

	return res
}
