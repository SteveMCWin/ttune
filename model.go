package main

import (
	"log"

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

	Data SettingsData
	SettingsSelected AppSettings
	AsciiArt string
}

func NewModel() Model {
	m := Model{
		BlockLength:    BL,
		CurrentState:   Initializing,
		SettingsSelected: LoadSettingsSelections(),
		Data: LoadSettingsData(),
	}

	log.Println(m.Data)

	m.ApplySettings()

	return m
}

func (m *Model) ApplySettings() {
	m.SettingsSelected = LoadSettingsSelections()

	m.AsciiArt = m.Data.AsciiArt[m.SettingsSelected.AsciiArtFileName]

	SetBorderStyle(m.Data.BorderStyles[m.SettingsSelected.BorderStyle])

	m.SelectedTuning = m.Data.Tunings[m.SettingsSelected.SelectedTuning]

	m.Theme = m.Data.ColorThemes[m.SettingsSelected.SelectedTheme]

	// Store settings to json file
	StoreSettings(m.SettingsSelected)
}


func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetToCurrent(true),
		initAutioStream(),
	}

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch message := msg.(type) {
	case NoteReadingMsg:
		new_note := Note(message)
		m.Note = new_note
		if m.CurrentState == Listening {
			cmds = append(cmds, CalculateNote())
		}
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c", "q":
			seq := tea.Sequence(closeAudioStream(), tea.Quit)
			cmds = append(cmds, seq)
		case "?":
			m.CurrentState = Help
		case "backspace":
			m.CurrentState = Listening
			cmds = append(cmds, CalculateNote())
		case "s":
			m.CurrentState = Settings
		case "h", "left":
			if m.CurrentState == Settings || m.CurrentState == Help {
				log.Println("LEFT")
			}
		case "j", "down":
			if m.CurrentState == Settings || m.CurrentState == Help {
				log.Println("DOWN")
			}
		case "k", "up":
			if m.CurrentState == Settings || m.CurrentState == Help {
				log.Println("UP")
			}
		case "l", "right":
			if m.CurrentState == Settings || m.CurrentState == Help {
				log.Println("RIGHT")
			}
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
