package tuning

type Tuning struct {
	Name string `json:"name"`
	Notes []string `json:"notes"`
}

var OctaveEnds []float64

var NoteNames = [12]string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}

const NUM_SEMITONES = 12

func init() {

	curr_freq := 30.87
	half_note := 1.01395947979

	for range 10 {
		OctaveEnds = append(OctaveEnds, curr_freq*half_note)
		curr_freq *= 2
	}
}
