package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"tuner/tuning"
)

type SettingsData struct {
	Tunings      map[string]tuning.Tuning `json:"tunings"`
	ColorThemes  map[string]ColorTheme    `json:"color_themes"`
	BorderStyles map[string]string        `json:"border_styles"`
	AsciiArt     map[string]string        `json:"ascii_art"`
}

type SettingsOptions struct {
	Description string   `json:"description"`
	Options     []string `json:"options"`
}

type AppSettings struct {
	AsciiArtFileName string `json:"ascii_art_filename"`
	SelectedTuning   string `json:"selected_tuning"`
	BorderStyle      string `json:"border_style"`
	SelectedTheme    string `json:"selected_theme"`
	// Tunings map[string]tuning.Tuning `json:"tunings"`
	// ColorThemes map[string]ColorTheme `json:"themes"`
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
		log.Println("Local settings data not found, reading default and writing to ", config_dir)
		data, err = os.ReadFile("./config/settings_data.json")
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

func LoadAsciiArt() map[string]string {
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

		cmd_str := "cp ./config/art/* " + user_art_dir_path
		cmd := exec.Command(cmd_str)
		if err := cmd.Run(); err != nil {
			log.Fatal("Error copying art files to user config dir")
		}
	}

	files, err := os.ReadDir(user_art_dir_path)
	if err != nil {
		log.Fatal("Error reading art dir")
	}

	ascii_art := make(map[string]string)
	for _, f := range files {
		data, err := os.ReadFile(f.Name())
		if err != nil {
			log.Fatal("Error reading", f.Name())
		}

		ascii_art[f.Name()] = string(data)
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
