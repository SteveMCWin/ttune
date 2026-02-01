package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"tuner/tuning"
)

type AppSettings struct {
	AsciiArtFileName string `json:"ascii_art_filename"`
	// ShowAsciiArt bool `json:"show_ascii_art"`
	SelectedTuning string `json:"selected_tuning"`
	BorderStyle string `json:"border_style"`
	Tunings map[string]tuning.Tuning `json:"tunings"`
	ColorThemes map[string]ColorTheme `json:"themes"`
	SelectedTheme string `json:"selected_theme"`
}

func LoadSettings() AppSettings {

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
		log.Println("Local config not found, reading default and writing to $HOME/.config/")
		data, err = os.ReadFile("./config/default.json")
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
