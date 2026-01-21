package main

import (
	"encoding/json"
	"log"
	"os"
)

type AppSettings struct {
	AsciiArtFileName string `json:"ascii_art_filename"`
	ShowAsciiArt bool `json:"show_ascii_art"`
	TuningName string `json:"tuning_name"`
	ColorThemeName string `json:"color_theme_name"`
	BorderStyle string `json:"border_style"`
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

