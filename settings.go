package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ttune/tuning"

	"embed"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const MIN_BUFFER_LEN = 4096
const MAX_BUFFER_LEN = 4096 * 8

//go:embed config
var configFS embed.FS

type SettingsSelections struct {
	AsciiArt    string `json:"ascii_art_filename"`
	Tuning      string `json:"selected_tuning"`
	BorderStyle string `json:"border_style"`
	ColorTheme  string `json:"selected_theme"`

	BufferLength   int     `json:"buffer_length"`
	SampleRate     int     `json:"sample_rate"`
	MinFrequency   int     `json:"min_frequency"`
	MaxFrequency   int     `json:"max_frequency"`
	AmplTreshold   float32 `json:"amplitude_treshold"`
	YinMinTreshold float32 `json:"yin_min_treshold"`
	YinMaxTreshold float32 `json:"yin_max_treshold"`
	HistorySize    int     `json:"history_size"`
}

const (
	AsciiArtSetting = int(iota)
	BorderStyleSetting
	ColorThemeSetting
	TuningSetting

	BufferLengthSetting
	SampleRateSetting
	MinFrequencySetting
	MaxFrequencySetting
	AmplTresholdSetting
	YinMinTresholdSetting
	YinMaxTresholdSetting
	HistorySizeSetting

	NumSettings
)

// return []Setting{ascii_art, borders, themes, tunings, buffer_len, sample_rate, min_frequency, max_frequency, ampl_threshold, yin_min_threshold, yin_max_threshold, history_size}

// Data that is meant to be configured in json files
type SettingsData struct {
	Tunings      []tuning.Tuning `json:"tunings"`
	ColorThemes  []ColorTheme    `json:"color_themes"`
	BorderStyles []string        `json:"border_styles"`
	AsciiArt     []AsciiArt      // NOTE: not loaded from json but by looking at the art dir
}

type Setting struct {
	Name           string
	Description    string
	Options        []Option
	Previews       []string
	Selected       string
	Apply          func(selection string, current SettingsSelections) SettingsSelections
	GetIdxFromName func(selection string) int
	Clamp          func(string) string
}

func validateInt(s string) error {
	for _, r := range s {
		if r < '0' || r > '9' {
			return errors.New("only digits allowed")
		}
	}
	return nil
}

func validateFloat(s string) error {
	dots := 0
	for _, r := range s {
		if r == '.' {
			dots++
			if dots > 1 {
				return errors.New("only one decimal point allowed")
			}
		} else if r < '0' || r > '9' {
			return errors.New("only digits and one decimal point allowed")
		}
	}
	return nil
}

func (s Setting) SelectedIdx() int {
	return s.GetIdxFromName(s.Selected)
}

// Since an option can be an input field and a multi-choice selection, I made it an interface you interract through these functions
type Option interface {
	GetValue() string
	HanldeSelect() tea.Cmd
	Render(selected bool) string
	Update(tea.Msg) tea.Cmd
	IsFocused() bool
	SetTheme(t ColorTheme)
}

type MultiChoiceOption struct {
	value string
	theme ColorTheme
}

func (o *MultiChoiceOption) GetValue() string {
	return string(o.value)
}

func (o *MultiChoiceOption) HanldeSelect() tea.Cmd {
	return nil
}

func (o *MultiChoiceOption) Render(selected bool) string {
	if selected {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(o.theme.Tertiary)).Render(o.value)
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(o.theme.Secondary)).Render(o.value)
}

func (o *MultiChoiceOption) SetTheme(t ColorTheme) {
	o.theme = t
}

func (o *MultiChoiceOption) Update(_ tea.Msg) tea.Cmd { return nil }
func (o *MultiChoiceOption) IsFocused() bool          { return false }

type InputFieldOption struct {
	Input         textinput.Model
	hoveredStyles textinput.Styles
}

func (o *InputFieldOption) GetValue() string {
	return o.Input.Value()
}

func (o *InputFieldOption) HanldeSelect() tea.Cmd {
	if o.Input.Focused() {
		o.Input.Blur()
		return nil
	}
	return o.Input.Focus()
}

func (o *InputFieldOption) Render(selected bool) string {
	prompt := "- "
	if o.IsFocused() {
		prompt = "> "
	}

	o.Input.Prompt = prompt

	if selected && !o.Input.Focused() {
		tmp := o.Input
		tmp.SetStyles(o.hoveredStyles)
		return tmp.View()
	}
	return o.Input.View()
}

func (o *InputFieldOption) SetTheme(t ColorTheme) {
	// pri_style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Primary))
	sec_style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Secondary))
	ter_style := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Tertiary))

	// primaryState := textinput.StyleState{
	// 	Placeholder: pri_style,
	// 	Suggestion:  pri_style,
	// 	Prompt:      pri_style,
	// 	Text:        pri_style,
	// }

	secondaryState := textinput.StyleState{
		Placeholder: sec_style,
		Suggestion:  sec_style,
		Prompt:      sec_style,
		Text:        sec_style,
	}

	tertiaryState := textinput.StyleState{
		Placeholder: ter_style,
		Suggestion:  ter_style,
		Prompt:      ter_style,
		Text:        ter_style,
	}

	cursor := textinput.CursorStyle{
		Color:      lipgloss.Color(t.Tertiary),
		Shape:      tea.CursorBlock,
		Blink:      true,
		BlinkSpeed: time.Millisecond * 500,
	}

	textinput_style = textinput.Styles{
		Blurred: secondaryState,
		Focused: tertiaryState,
		Cursor:  cursor,
	}

	o.hoveredStyles = textinput.Styles{
		Blurred: tertiaryState,
		Focused: tertiaryState,
		Cursor:  cursor,
	}

	o.Input.SetStyles(textinput_style)
}

func (o *InputFieldOption) Update(msg tea.Msg) tea.Cmd {
	prev := o.Input.Value()
	var cmd tea.Cmd
	o.Input, cmd = o.Input.Update(msg)
	if o.Input.Err != nil {
		o.Input.SetValue(prev)
		o.Input.Err = nil
	}
	return cmd
}

func (o *InputFieldOption) IsFocused() bool {
	return o.Input.Focused()
}

func DefineAllSettingsOptions(data SettingsData, currentSettings SettingsSelections) []Setting {

	makeIntInput := func(current int) textinput.Model {
		input := textinput.New()
		input.SetValue(strconv.Itoa(current))
		input.Validate = validateInt
		return input
	}

	makeFloatInput := func(current float32) textinput.Model {
		input := textinput.New()
		input.SetValue(strconv.FormatFloat(float64(current), 'f', -1, 32))
		input.Validate = validateFloat
		return input
	}

	buffer_len := Setting{
		Name:        "Buffer Length",
		Description: "Number of audio samples processed per block. Larger values improve low-frequency detection accuracy but increase latency. Must be between 4096 and 32768.",
		Options:     []Option{&InputFieldOption{Input: makeIntInput(currentSettings.BufferLength)}},
		Previews:    []string{""},
		Selected:    strconv.Itoa(currentSettings.BufferLength),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Println("buffer_len: invalid value:", val)
				return s
			}
			s.BufferLength = n
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp: func(val string) string {
			if val == "" {
				return strconv.Itoa(MIN_BUFFER_LEN)
			}
			n, err := strconv.Atoi(val)
			if err != nil {
				return strconv.Itoa(MIN_BUFFER_LEN)
			}
			return strconv.Itoa(min(MAX_BUFFER_LEN, max(MIN_BUFFER_LEN, n)))
		},
	}

	sample_rate := Setting{
		Name:        "Sample Rate",
		Description: "Audio sample rate in Hz. Higher values capture more detail but require more processing power. 44100 Hz is the standard CD-quality rate.",
		Options:     []Option{&InputFieldOption{Input: makeIntInput(currentSettings.SampleRate)}},
		Previews:    []string{""},
		Selected:    strconv.Itoa(currentSettings.SampleRate),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Println("sample_rate: invalid value:", val)
				return s
			}
			s.SampleRate = n
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp: func(val string) string {
			if val == "" {
				return "44100"
			}
			n, err := strconv.Atoi(val)
			if err != nil {
				return "44100"
			}
			return strconv.Itoa(min(192000, max(8000, n)))
		},
	}

	min_frequency := Setting{
		Name:        "Min Frequency",
		Description: "Lowest frequency in Hz the tuner will attempt to detect. Lower values let you tune bass instruments but may introduce false readings. Default is 70 Hz.",
		Options:     []Option{&InputFieldOption{Input: makeIntInput(currentSettings.MinFrequency)}},
		Previews:    []string{""},
		Selected:    strconv.Itoa(currentSettings.MinFrequency),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Println("min_frequency: invalid value:", val)
				return s
			}
			s.MinFrequency = n
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp: func(val string) string {
			if val == "" {
				return "70"
			}
			n, err := strconv.Atoi(val)
			if err != nil {
				return "70"
			}
			return strconv.Itoa(min(5000, max(20, n)))
		},
	}

	max_frequency := Setting{
		Name:        "Max Frequency",
		Description: "Highest frequency in Hz the tuner will attempt to detect. Higher values cover more of the upper register. Default is 1500 Hz.",
		Options:     []Option{&InputFieldOption{Input: makeIntInput(currentSettings.MaxFrequency)}},
		Previews:    []string{""},
		Selected:    strconv.Itoa(currentSettings.MaxFrequency),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Println("max_frequency: invalid value:", val)
				return s
			}
			s.MaxFrequency = n
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp: func(val string) string {
			if val == "" {
				return "1500"
			}
			n, err := strconv.Atoi(val)
			if err != nil {
				return "1500"
			}
			return strconv.Itoa(min(20000, max(100, n)))
		},
	}

	clampFloat01 := func(val, fallback string) string {
		if val == "" || val == "." {
			return fallback
		}
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return fallback
		}
		if f < 0 {
			f = 0
		} else if f > 1 {
			f = 1
		}
		return strconv.FormatFloat(f, 'f', -1, 32)
	}

	ampl_threshold := Setting{
		Name:        "Amplitude Threshold",
		Description: "Minimum RMS signal level required before pitch detection runs. Raise this if background noise is triggering false readings. Default is 0.01.",
		Options:     []Option{&InputFieldOption{Input: makeFloatInput(currentSettings.AmplTreshold)}},
		Previews:    []string{""},
		Selected:    strconv.FormatFloat(float64(currentSettings.AmplTreshold), 'f', -1, 32),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				log.Println("ampl_threshold: invalid value:", val)
				return s
			}
			s.AmplTreshold = float32(f)
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp:          func(val string) string { return clampFloat01(val, "0.01") },
	}

	yin_min_threshold := Setting{
		Name:        "YIN Min Threshold",
		Description: "YIN candidate threshold for pitch detection. Lower values are stricter and reduce harmonic errors but may miss weak signals. Default is 0.10.",
		Options:     []Option{&InputFieldOption{Input: makeFloatInput(currentSettings.YinMinTreshold)}},
		Previews:    []string{""},
		Selected:    strconv.FormatFloat(float64(currentSettings.YinMinTreshold), 'f', -1, 32),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				log.Println("yin_min_threshold: invalid value:", val)
				return s
			}
			s.YinMinTreshold = float32(f)
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp:          func(val string) string { return clampFloat01(val, "0.1") },
	}

	yin_max_threshold := Setting{
		Name:        "YIN Max Threshold",
		Description: "YIN validity ceiling — readings above this power threshold are discarded as weak detections. Raise to accept more readings, lower to filter noise. Default is 0.85.",
		Options:     []Option{&InputFieldOption{Input: makeFloatInput(currentSettings.YinMaxTreshold)}},
		Previews:    []string{""},
		Selected:    strconv.FormatFloat(float64(currentSettings.YinMaxTreshold), 'f', -1, 32),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				log.Println("yin_max_threshold: invalid value:", val)
				return s
			}
			s.YinMaxTreshold = float32(f)
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp:          func(val string) string { return clampFloat01(val, "0.85") },
	}

	history_size := Setting{
		Name:        "History Size",
		Description: "Number of recent frequency readings used by the median filter for smoothing. Larger values stabilise the display but slow response to pitch changes. Default is 5.",
		Options:     []Option{&InputFieldOption{Input: makeIntInput(currentSettings.HistorySize)}},
		Previews:    []string{""},
		Selected:    strconv.Itoa(currentSettings.HistorySize),
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			n, err := strconv.Atoi(val)
			if err != nil {
				log.Println("history_size: invalid value:", val)
				return s
			}
			s.HistorySize = n
			return s
		},
		GetIdxFromName: func(selection string) int { return 0 },
		Clamp: func(val string) string {
			if val == "" {
				return "5"
			}
			n, err := strconv.Atoi(val)
			if err != nil {
				return "5"
			}
			return strconv.Itoa(min(50, max(1, n)))
		},
	}

	ascii_art := Setting{
		Name:        "Ascii Art",
		Description: "The character art displayed on the left side of the terminal when tuning is in progress. Purely for aesthetical purposes, but I spent a lot of time drawing it :^)",
		Options:     make([]Option, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.AsciiArt,
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			s.AsciiArt = val
			return s
		},
		GetIdxFromName: func(selection string) int {
			for i, v := range data.AsciiArt {
				if v.FileName == selection {
					return i
				}
			}
			return 0
		},
	}

	for _, v := range data.AsciiArt {
		ascii_art.Options = append(ascii_art.Options, &MultiChoiceOption{value: v.FileName})
		// Note: JoinHorizontal is needed for some reason so the rows don't auto align for some reason
		ascii_art.Previews = append(
			ascii_art.Previews,
			lipgloss.JoinHorizontal(lipgloss.Center, "", v.FileContents),
		)
	}

	borders := Setting{
		Name:        "Border Style",
		Description: "The appearance of displayed borders. The double border is my favourite, that's why it's the default one hihi.",
		Options:     make([]Option, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.BorderStyle,
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			s.BorderStyle = val
			return s
		},
		GetIdxFromName: func(selection string) int {
			for i, v := range data.BorderStyles {
				if v == selection {
					return i
				}
			}

			return 0
		},
	}

	for _, v := range data.BorderStyles {
		borders.Options = append(borders.Options, &MultiChoiceOption{value: v})

		borders.Previews = append(
			borders.Previews,
			GetBorderStyleByName(v).Width(12).Height(6).Render(""),
		)
	}

	themes := Setting{
		Name:        "Color Theme",
		Description: "Colors used for displaying the user interface. Comprised of 3 colors each. Affects only the foreground elements.",
		Options:     make([]Option, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.ColorTheme,
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			s.ColorTheme = val
			return s
		},
		GetIdxFromName: func(selection string) int {
			for i, v := range data.ColorThemes {
				if v.Name == selection {
					return i
				}
			}

			return 0
		},
	}

	for _, v := range data.ColorThemes {
		themes.Options = append(themes.Options, &MultiChoiceOption{value: v.Name})
		blocks := `
███████
███████
███████
`
		tmp_style := lipgloss.NewStyle()
		preview := tmp_style.Foreground(lipgloss.Color(v.Primary)).Render(blocks)
		preview += tmp_style.Foreground(lipgloss.Color(v.Secondary)).Render(blocks)
		preview += tmp_style.Foreground(lipgloss.Color(v.Tertiary)).Render(blocks)
		themes.Previews = append(themes.Previews, preview)
	}

	tunings := Setting{
		Name:        "Displayed Tuning",
		Description: "A tuning that will be displayed along the ascii art. Mostly there for aesthetical reasons but also quite handy if you don't have them memorized exactly.",
		Options:     make([]Option, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.Tuning,
		Apply: func(val string, s SettingsSelections) SettingsSelections {
			s.Tuning = val
			return s
		},
		GetIdxFromName: func(selection string) int {
			for i, v := range data.Tunings {
				if v.Name == selection {
					return i
				}
			}

			return 0
		},
	}

	for _, v := range data.Tunings {
		tunings.Options = append(tunings.Options, &MultiChoiceOption{value: v.Name})
		var builder strings.Builder
		builder.WriteByte('\n')
		for _, note := range v.Notes {
			builder.WriteString(note)
			builder.WriteByte('\n')
		}
		tunings.Previews = append(tunings.Previews, builder.String())
	}

	res := make([]Setting, NumSettings)

	res[AsciiArtSetting] = ascii_art
	res[BorderStyleSetting] = borders
	res[ColorThemeSetting] = themes
	res[TuningSetting] = tunings
	res[BufferLengthSetting] = buffer_len
	res[SampleRateSetting] = sample_rate
	res[SampleRateSetting] = sample_rate
	res[MinFrequencySetting] = min_frequency
	res[MaxFrequencySetting] = max_frequency
	res[AmplTresholdSetting] = ampl_threshold
	res[YinMinTresholdSetting] = yin_min_threshold
	res[YinMaxTresholdSetting] = yin_max_threshold
	res[HistorySizeSetting] = history_size

	return res
}

func LoadOrWriteConfigFile(config_file_name string) ([]byte, error) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_config_dir_path := filepath.Join(config_dir, "ttune")
	user_config_file_path := filepath.Join(user_config_dir_path, config_file_name)

	// if ttune folder doesn't already exist in the config dir, create it
	if _, err := os.Stat(user_config_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_config_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Println("Error while making config dir", err)
			return []byte{}, err
		}
	}

	data, err := os.ReadFile(user_config_file_path)
	if err != nil {
		log.Println(config_file_name, "not found, reading default and writing to", config_dir)
		data, err = configFS.ReadFile("config/" + config_file_name)
		if err != nil {
			log.Println("Error reading config file!!!!")
			return []byte{}, err
		}

		err = os.WriteFile(user_config_file_path, data, 0744)
		if err != nil {
			log.Println("Error saving config file!!!!")
			return []byte{}, err
		}
	}

	return data, nil
}

func LoadSettingsData() SettingsData {
	var res SettingsData

	// Load users options
	user_options_filename := "custom_options.json"
	data, err := LoadOrWriteConfigFile(user_options_filename)
	if err != nil {
		return SettingsData{}
	}
	var formatted SettingsData
	if err := json.Unmarshal(data, &formatted); err != nil {
		log.Println("Error unmarshaling settings data")
		panic(err)
	}
	res.BorderStyles = append(res.BorderStyles, formatted.BorderStyles...)
	res.ColorThemes = append(res.ColorThemes, formatted.ColorThemes...)
	res.Tunings = append(res.Tunings, formatted.Tunings...)

	// Load default options
	defaultData, err := configFS.ReadFile("config/default_options.json")
	if err != nil {
		log.Println("Error reading embedded default_options.json")
		panic(err)
	}
	var defaults SettingsData
	if err := json.Unmarshal(defaultData, &defaults); err != nil {
		log.Println("Error unmarshaling default settings data")
		panic(err)
	}
	existingThemes := make(map[string]bool, len(res.ColorThemes))
	for _, t := range res.ColorThemes {
		existingThemes[t.Name] = true
	}
	for _, t := range defaults.ColorThemes {
		if !existingThemes[t.Name] {
			res.ColorThemes = append(res.ColorThemes, t)
		}
	}

	existingTunings := make(map[string]bool, len(res.Tunings))
	for _, t := range res.Tunings {
		existingTunings[t.Name] = true
	}
	for _, t := range defaults.Tunings {
		if !existingTunings[t.Name] {
			res.Tunings = append(res.Tunings, t)
		}
	}

	existingBorders := make(map[string]bool, len(res.BorderStyles))
	for _, b := range res.BorderStyles {
		existingBorders[b] = true
	}
	for _, b := range defaults.BorderStyles {
		if !existingBorders[b] {
			res.BorderStyles = append(res.BorderStyles, b)
		}
	}

	res.AsciiArt = LoadAsciiArt()
	return res
}

func LoadAsciiArt() []AsciiArt {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_art_dir_path := filepath.Join(config_dir, "ttune", "art")

	// Create art dir if it doesn't exist
	if err := os.MkdirAll(user_art_dir_path, 0744); err != nil {
		log.Fatal("Error creating art dir: ", err)
	}

	// Sync embedded art files
	default_files, err := configFS.ReadDir("config/art")
	if err != nil {
		log.Fatal("Error reading embedded art dir")
	}

	for _, f := range default_files {
		dest := filepath.Join(user_art_dir_path, f.Name())
		if _, err := os.Stat(dest); errors.Is(err, os.ErrNotExist) {
			data, err := configFS.ReadFile("config/art/" + f.Name())
			if err != nil {
				log.Fatal("Error reading embedded file ", f.Name(), " :: ", err)
			}
			if err := os.WriteFile(dest, data, 0744); err != nil {
				log.Fatal("Error writing art file ", f.Name(), " :: ", err)
			}
		}
	}

	// Load all art from the user's dir (includes both defaults and any user-added files)
	files, err := os.ReadDir(user_art_dir_path)
	if err != nil {
		log.Fatal("Error reading user art dir")
	}

	ascii_art := make([]AsciiArt, 0)
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(user_art_dir_path, f.Name()))
		if err != nil {
			log.Fatal("Error reading ", f.Name())
		}
		ascii_art = append(ascii_art, AsciiArt{FileName: f.Name(), FileContents: string(data)})
	}

	return ascii_art
}

func LoadSettingsSelections(config_file_name string, data SettingsData) SettingsSelections {
	raw, err := LoadOrWriteConfigFile(config_file_name)
	if err != nil {
		return SettingsSelections{}
	}

	var settings SettingsSelections
	if err := json.Unmarshal(raw, &settings); err == nil {
		return settings
	}

	// Attempt to migrate from legacy format where selections were stored as indices
	var legacy struct {
		AsciiArt    int `json:"ascii_art_filename"`
		Tuning      int `json:"selected_tuning"`
		BorderStyle int `json:"border_style"`
		ColorTheme  int `json:"selected_theme"`
	}
	if err := json.Unmarshal(raw, &legacy); err != nil {
		log.Println("Error parsing settings json, resetting to defaults")
		return SettingsSelections{}
	}
	log.Println("Migrating settings from legacy index format to string format")
	if legacy.AsciiArt < len(data.AsciiArt) {
		settings.AsciiArt = data.AsciiArt[legacy.AsciiArt].FileName
	}
	if legacy.Tuning < len(data.Tunings) {
		settings.Tuning = data.Tunings[legacy.Tuning].Name
	}
	if legacy.BorderStyle < len(data.BorderStyles) {
		settings.BorderStyle = data.BorderStyles[legacy.BorderStyle]
	}
	if legacy.ColorTheme < len(data.ColorThemes) {
		settings.ColorTheme = data.ColorThemes[legacy.ColorTheme].Name
	}
	StoreSettings(settings, config_file_name)
	return settings
}

func StoreSettings(settings SettingsSelections, config_file_name string) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_config_dir_path := filepath.Join(config_dir, "ttune")
	user_config_file_path := filepath.Join(user_config_dir_path, config_file_name)

	if _, err := os.Stat(user_config_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_config_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		log.Println("Error marshaling config settings")
		panic(err)
	}

	err = os.WriteFile(user_config_file_path, data, 0744)
	if err != nil {
		log.Println("Error saving config file!!!!")
		panic(err)
	}
}
