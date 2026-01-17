package main

import (
	"log"

	"tuner/tuning"

	"github.com/gordonklaus/portaudio"
	tea "charm.land/bubbletea/v2"
)

type OpenStreamMsg *portaudio.Stream
type NoteReadingMsg Note
type State string

const (
	Initializing State = "Initializing"
	Listening State = "Listening"
	Settings State = "Settings"
	Help State = "Help"
)

type Note struct {
	Index int
	Octave int
	CentsOff int
}

type Model struct {
	Theme ColorTheme

	WindowWidth int
	WindowHeight int

	BlockLength int
	Frequency   float64
	Note        Note
	CentsOff    float64

	// AudioStream *portaudio.Stream

	CurrentState State

	SelectedTuning tuning.Tuning
}

func NewModel() Model {

	m := Model {
		BlockLength: BL,
		CurrentState:   Initializing,
		SelectedTuning: tuning.Tunings[tuning.Standard],
		Theme:          DefaultTheme,
	}

	return m
}

func LoadLocalData() tea.Msg {
	return "todo"
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetCurrentTheme(true), // NOTE: hard coded for testing
		initAutioStream(),
	}

	return tea.Batch(cmds...) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
}

// func (m *Model) CloseAudioStream() tea.Cmd {
// 	err := m.AudioStream.Stop()
// 	if err != nil {
// 		log.Println("ERROR STOPPING AUDIO STREAM")
// 	}
//
// 	err = m.AudioStream.Close()
// 	if err != nil {
// 		log.Println("ERROR CLOSING AUDIO STREAM")
// 	}
// 	return nil
// }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch message := msg.(type) {
	case NoteReadingMsg:
		new_note := Note(message)
		if m.Note.Index != new_note.Index || m.Note.CentsOff != new_note.CentsOff {
			log.Printf("Old note: %2s%d %3d   New note: %2s%d %3d\n",
				tuning.NoteNames[m.Note.Index], m.Note.Octave, m.Note.CentsOff,
				tuning.NoteNames[new_note.Index], new_note.Octave, new_note.CentsOff)
		}
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
