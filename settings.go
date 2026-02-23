package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"ttune/tuning"

	"embed"
	"github.com/charmbracelet/lipgloss/v2"
)

//go:embed config
var configFS embed.FS

type SettingsData struct {
	Tunings      []tuning.Tuning `json:"tunings"`
	ColorThemes  []ColorTheme    `json:"color_themes"`
	BorderStyles []string        `json:"border_styles"`
	AsciiArt     []AsciiArt      // NOTE: not loaded from json but by looking at the art dir
}

type SettingsOptions struct {
	Name        string
	Description string
	Options     []string
	Previews    []string
	Selected    int
	Apply       func(val int, current AppSettings) AppSettings
}

type AppSettings struct {
	AsciiArt    int `json:"ascii_art_filename"`
	Tuning      int `json:"selected_tuning"`
	BorderStyle int `json:"border_style"`
	ColorTheme  int `json:"selected_theme"`
}

func DefineVisualSettingsOptions(data SettingsData, currentSettings AppSettings) []SettingsOptions {
	ascii_art := SettingsOptions{
		Name:        "Ascii Art",
		Description: "The character art displayed on the left side of the terminal when tuning is in progress. Purely for aesthetical purposes, but I spent a lot of time drawing it :^)",
		Options:     make([]string, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.AsciiArt,
		Apply: func(val int, s AppSettings) AppSettings {
			s.AsciiArt = val
			return s
		},
	}

	for _, v := range data.AsciiArt {
		ascii_art.Options = append(ascii_art.Options, v.FileName)
		// Note: JoinHorizontal is needed for some reason so the rows don't auto align for some reason
		ascii_art.Previews = append(ascii_art.Previews, lipgloss.JoinHorizontal(lipgloss.Center, "", v.FileContents))
	}

	borders := SettingsOptions{
		Name:        "Border Style",
		Description: "The appearance of displayed borders. The double border is my favourite, that's why it's the default one hihi.",
		Options:     make([]string, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.BorderStyle,
		Apply: func(val int, s AppSettings) AppSettings {
			s.BorderStyle = val
			return s
		},
	}

	for _, v := range data.BorderStyles {
		borders.Options = append(borders.Options, v)

		borders.Previews = append(borders.Previews, GetBorderStyleByName(v).Width(12).Height(6).Render(""))
	}

	themes := SettingsOptions{
		Name:        "Color Theme",
		Description: "Colors used for displaying the user interface. Comprised of 3 colors each. Affects only the foreground elements.",
		Options:     make([]string, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.ColorTheme,
		Apply: func(val int, s AppSettings) AppSettings {
			s.ColorTheme = val
			return s
		},
	}

	for _, v := range data.ColorThemes {
		themes.Options = append(themes.Options, v.Name)
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

	tunings := SettingsOptions{
		Name:        "Displayed Tuning",
		Description: "A tuning that will be displayed along the ascii art. Mostly there for aesthetical reasons but also quite handy if you don't have them memorized exactly.",
		Options:     make([]string, 0),
		Previews:    make([]string, 0),
		Selected:    currentSettings.Tuning,
		Apply: func(val int, s AppSettings) AppSettings {
			s.Tuning = val
			return s
		},
	}

	for _, v := range data.Tunings {
		tunings.Options = append(tunings.Options, v.Name)
		var builder strings.Builder
		builder.WriteByte('\n')
		for _, note := range v.Notes {
			builder.WriteString(note)
			builder.WriteByte('\n')
		}
		tunings.Previews = append(tunings.Previews, builder.String())
	}

	return []SettingsOptions{ascii_art, borders, themes, tunings}
}

func LoadSettingsData() SettingsData {
	// Get path of user config dir
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_config_dir_path := filepath.Join(config_dir, "ttune")
	user_config_file_path := filepath.Join(user_config_dir_path, "settings_data.json")

	// if ttune folder doesn't already exist in the config dir, create it
	if _, err := os.Stat(user_config_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_config_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	data, err := os.ReadFile(user_config_file_path)
	if err != nil {
		log.Println("Local settings data not found, reading default and writing to", config_dir)
		data, err = configFS.ReadFile("config/settings_data.json")
		if err != nil {
			log.Println("Error reading config file!!!!")
			panic(err)
		}

		err = os.WriteFile(user_config_file_path, data, 0744)
		if err != nil {
			log.Println("Error saving config file!!!!")
			panic(err)
		}
	}

	var res SettingsData
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Println("Error unmarshaling settings data")
		panic(err)
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

	user_art_dir_path := filepath.Join(filepath.Join(config_dir, "ttune"), "art")

	// if art folder doesn't already exist in the ttune config dir, create it
	if _, err := os.Stat(user_art_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_art_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		default_files, err := configFS.ReadDir("config/art")
		if err != nil {
			log.Fatal("Error reading art dir")
		}
		for _, f := range default_files {
			data, err := configFS.ReadFile("config/art/"+f.Name())
			if err != nil {
				log.Fatal("Error reading ", f.Name(), " :: ", err)
			}
			err = os.WriteFile(filepath.Join(user_art_dir_path, f.Name()), data, 0744)
		}
	}

	files, err := os.ReadDir(user_art_dir_path)
	if err != nil {
		log.Fatal("Error reading art dir")
	}

	ascii_art := make([]AsciiArt, 0)
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(user_art_dir_path, f.Name()))
		if err != nil {
			log.Fatal("Error reading", f.Name())
		}

		ascii_art = append(ascii_art, AsciiArt{FileName: f.Name(), FileContents: string(data)})
	}

	return ascii_art
}

func LoadSettingsSelections() AppSettings {

	// Get path of user config dir
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_config_dir_path := filepath.Join(config_dir, "ttune")
	user_config_file_path := filepath.Join(user_config_dir_path, "config.json")

	// if ttune folder doesn't already exist in the config dir, create it
	if _, err := os.Stat(user_config_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_config_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	data, err := os.ReadFile(user_config_file_path)
	if err != nil {
		log.Println("Local config not found, reading default and writing to", config_dir)
		data, err = configFS.ReadFile("config/default.json")
		if err != nil {
			log.Println("Error reading config file!!!!")
			panic(err)
		}

		err = os.WriteFile(user_config_file_path, data, 0744)
		if err != nil {
			log.Println("Error saving config file!!!!")
			panic(err)
		}
	}

	var settings AppSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Println("Error parsing json data")
		panic(err)
	}

	return settings
}

func StoreSettings(settings AppSettings) {
	config_dir, err := os.UserConfigDir()
	if err != nil {
		log.Println("Error finding user config dir")
		panic(err)
	}

	user_config_dir_path := filepath.Join(config_dir, "ttune")
	user_config_file_path := filepath.Join(user_config_dir_path, "config.json")

	if _, err := os.Stat(user_config_dir_path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(user_config_dir_path, 0744)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	err = os.WriteFile(user_config_file_path, data, 0744)
	if err != nil {
		log.Println("Error saving config file!!!!")
		panic(err)
	}
}
