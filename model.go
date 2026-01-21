package main

import (
	"encoding/json"
	"log"
	"os"

	"tuner/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"
)

type OpenStreamMsg *portaudio.Stream
type NoteReadingMsg Note
type State string

const (
	Initializing State = "Initializing"
	Listening    State = "Listening"
	Settings     State = "Settings"
	Help         State = "Help"
)

type Note struct {
	Index    int
	Octave   int
	CentsOff int
}

type Model struct {
	Theme ColorTheme

	WindowWidth  int
	WindowHeight int

	BlockLength int
	Frequency   float64
	Note        Note
	CentsOff    float64

	CurrentState State

	SelectedTuning tuning.Tuning

	Settings AppSettings
	AsciiArt string
}

type AppSettings struct {
	AsciiArtFileName string `json:"ascii_art_filename"`
	ShowAsciiArt bool `json:"show_ascii_art"`
}

func NewModel() Model {
	m := Model{
		BlockLength:    BL,
		CurrentState:   Initializing,
		SelectedTuning: tuning.Tunings[tuning.Standard],
		Theme:          WhiteTheme,
	}

	m.LoadSettings()

	return m
}

func (m *Model) LoadSettings() {
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

	m.Settings = settings

	data, err = os.ReadFile(settings.AsciiArtFileName)
	if err != nil {
		log.Println("Error reading ascii art file name")
		panic(err)
	}
	m.AsciiArt = string(data)
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetCurrentTheme(true), // NOTE: hard coded for testing
		initAutioStream(),
	}

	return tea.Batch(cmds...) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch message := msg.(type) {
	case NoteReadingMsg:
		new_note := Note(message)
		m.Note = new_note
		cmds = append(cmds, CalculateNote())
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c", "q":
			seq := tea.Sequence(closeAudioStream(), tea.Quit)
			cmds = append(cmds, seq)
		case "?", "h":
			m.CurrentState = Help
		case "backspace", "escape":
			m.CurrentState = Listening
		case "s", "tab":
			m.CurrentState = Settings
		}
	case tea.WindowSizeMsg:
		// log.Println("Terminal width:", message.Width)
		// log.Println("Terminal height:", message.Height)
		m.WindowWidth = message.Width
		m.WindowHeight = message.Height
	case OpenStreamMsg:
		log.Println("Opened stream")
		m.CurrentState = Listening
		cmds = append(cmds, CalculateNote())
	default:
	}

	return m, tea.Batch(cmds...)
}
