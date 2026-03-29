package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ttune/tuning"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"embed"
)

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
}

func (s Setting) SelectedIdx() int {
	return s.GetIdxFromName(s.Selected)
}

// Since an option can be an input field and a multi-choice selection, I made it an interface you interract through these functions
type Option interface {
	GetValue() string
	HanldeSelect() string
	Render() string
}

type MultiChoiceOption string

func (o MultiChoiceOption) GetValue() string {
	return string(o)
}

func (o MultiChoiceOption) HanldeSelect() string {
	return o.GetValue()
}

func (o MultiChoiceOption) Render() string {
	return o.GetValue()
}

type InputFieldOption struct {
	Input textinput.Model
}

func (o InputFieldOption) GetValue() string {
	return o.Input.Value()
}

func (o InputFieldOption) HanldeSelect() string {
	if o.Input.Focused() {
		o.Input.Blur()
	} else {
		o.Input.Focus()
	}

	return o.GetValue()
}

func (o InputFieldOption) Render() string {
	return o.Input.View()
}

func DefineVisualSettingsOptions(data SettingsData, currentSettings SettingsSelections) []Setting {
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
		ascii_art.Options = append(ascii_art.Options, MultiChoiceOption(v.FileName))
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
		borders.Options = append(borders.Options, MultiChoiceOption(v))

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
		themes.Options = append(themes.Options, MultiChoiceOption(v.Name))
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
		tunings.Options = append(tunings.Options, MultiChoiceOption(v.Name))
		var builder strings.Builder
		builder.WriteByte('\n')
		for _, note := range v.Notes {
			builder.WriteString(note)
			builder.WriteByte('\n')
		}
		tunings.Previews = append(tunings.Previews, builder.String())
	}

	// functional := SettingsOptions{
	// 	Name: "Functional Settings",
	// 	Description: "Settings that affect the pitch detection algorithm. Mess around with these to get the optimal tuning for your setup!",
	// 	Options: make([]Option, 0),
	// 	Previews: make([]string, 0),
	// 	Selected: 0,
	// 	Apply: func(val int, s SettingsSelections) SettingsSelections {
	//
	// 	}
	// }
	return []Setting{ascii_art, borders, themes, tunings}
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
	res.BorderStyles = append(res.BorderStyles, defaults.BorderStyles...)
	res.ColorThemes = append(res.ColorThemes, defaults.ColorThemes...)
	res.Tunings = append(res.Tunings, defaults.Tunings...)

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
