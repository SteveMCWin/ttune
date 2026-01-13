package tuning

import (
	"log"
	"strings"
)

const (
	Standard string = "Standard"
)

type Tuning struct {
	Name string
	Notes []string
}

var Tunings map[string]Tuning

var NoteNames = []string{"C ", "C#", "D ", "D#", "E ", "F ", "F#", "G ", "G#", "A ", "A#", "B "}

// NOTE: place this in a file where each line starts with tuning name and then the notes
// this way the user can add their own tunings
const STD_TUNE_STR = "E2;A2;D3;G3;B3;E4"
const SEP_CHAR = ";"

func init() {
	Tunings = make(map[string]Tuning)
	Tunings[Standard] = Tuning{ Name: Standard, Notes:  strings.Split(STD_TUNE_STR, SEP_CHAR) }
	log.Println(Tunings["standard"])
}
