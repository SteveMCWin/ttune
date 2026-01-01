package tuning

import "strings"

type Tuning struct {
	Notes []string
}
var Tunings map[string]Tuning

// NOTE: place this in a file where each line starts with tuning name and then the notes
// this way the user can add their own tunings
const STD_TUNE_STR = "E2;A2;D3;G3;B3;E4"
const SEP_CHAR = ";"

func init() {
	Tunings["standard"] = Tuning{ Notes:  strings.Split(STD_TUNE_STR, SEP_CHAR) }
}
