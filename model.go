package main

import (
	"log"
	"net/http"
	"strings"

	"tuner/tuning"

	"github.com/gordonklaus/portaudio"
	tea "charm.land/bubbletea/v2"
)

type State int

const (
	Tuning State = iota
	Settings
	Help
)

type Model struct {
	// windowWidth  int
	// windowHeight int

	Theme ColorTheme
	// err      error

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
		CurrentState:   Tuning,
		SelectedTuning: tuning.Tunings["standard"],
		Theme:          DefaultTheme,
		Buffer: make([]float32, BL),
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

	switch m.CurrentState {
	case Tuning:
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
			m.CurrentState = Tuning
		case "s", "tab":
			m.CurrentState = Settings
		}
	default:
	}

	return m, tea.Batch(cmds...)
}
