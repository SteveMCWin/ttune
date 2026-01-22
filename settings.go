package main

import (
	"encoding/json"
	"log"
	"os"
	"tuner/tuning"
)

type AppSettings struct {
	AsciiArtFileName string `json:"ascii_art_filename"`
	ShowAsciiArt bool `json:"show_ascii_art"`
	SelectedTuning string `json:"selected_tuning"`
	BorderStyle string `json:"border_style"`
	Tunings map[string]tuning.Tuning `json:"tunings"`
	ColorThemes map[string]ColorTheme `json:"themes"`
	SelectedTheme string `json:"selected_theme"`
}

func LoadSettings() AppSettings {
	data, err := os.ReadFile("./config/default.json")
	if err != nil {
		log.Println("Error reading config file!!!!")
		panic(err)
	}

	var settings AppSettings
	err = json.Unmarshal(data, &settings)
	if err != nil {
		log.Println("Error parsing json data")
		panic(err)
	}

	return settings
}

