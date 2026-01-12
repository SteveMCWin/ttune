package main

import (
	"log"

	"tuner/tuning"

	"github.com/gordonklaus/portaudio"
	tea "charm.land/bubbletea/v2"
)

type OpenStreamMsg *portaudio.Stream
type State string

const (
	Initializing State = "Initializing"
	Listening State = "Listening"
	Settings State = "Settings"
	Help State = "Help"
)

type Model struct {
	// windowWidth  int
	// windowHeight int

	Theme ColorTheme
	// err      error

	WindowWidth int
	WindowHeight int

	Buffer      []float32
	Buffer64    []float64
	BlockLength int
	Window      []float64
	Frequency   float64
	Note        string
	CentsOff    float64

	AudioStream *portaudio.Stream

	CurrentState State

	SelectedTuning tuning.Tuning
}

func NewModel() Model {

	m := Model {
		BlockLength: BL,
		CurrentState:   Initializing,
		SelectedTuning: tuning.Tunings["standard"],
		Theme:          DefaultTheme,
		Buffer: make([]float32, BL),
		Buffer64: make([]float64, BL),
		Window: make([]float64, BL),
	}


	return m
}

func LoadLocalData() tea.Msg {
	return "todo"
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetCurrentTheme(true), // NOTE: hard coded for testing
		m.calcHannWindow(),
		m.openAudioStream(),
	}

	return tea.Batch(cmds...) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
}

func (m *Model) CloseAudioStream() tea.Cmd {
	err := m.AudioStream.Stop()
	if err != nil {
		log.Println("ERROR STOPPING AUDIO STREAM")
	}

	err = m.AudioStream.Close()
	if err != nil {
		log.Println("ERROR CLOSING AUDIO STREAM")
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	if m.CurrentState == Listening {
		if m.AudioStream == nil {
			log.Println("Audio stream is nil?!?!?!")
		}
		m.AudioStream.Read()
		m.buffTo64()
		m.applyWindowToBuffer()
		m.calculateFrequency()
		m.GetNote()
	}

	switch message := msg.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c", "q":
			seq := tea.Sequence(m.CloseAudioStream(), tea.Quit)
			cmds = append(cmds, seq)
		case "?", "h":
			m.CurrentState = Help
		case "backspace", "escape":
			m.CurrentState = Listening
		case "s", "tab":
			m.CurrentState = Settings
		}
	case tea.WindowSizeMsg:
		log.Println("Terminal width:", message.Width)
		log.Println("Terminal height:", message.Height)
		m.WindowWidth = message.Width
		m.WindowHeight = message.Height
	case OpenStreamMsg:
		log.Println("Opened stream")
		m.AudioStream = message
		m.CurrentState = Listening
	default:
	}

	return m, tea.Batch(cmds...)
}
