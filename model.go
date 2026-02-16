package main

import (
	"log"

	"ttune/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"
)

type OpenStreamMsg *portaudio.Stream
type NoteReadingMsg Note
type State string

type ReRenderMsg struct{}

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

func ReRender() tea.Msg {
	return ReRenderMsg{}
}

type Model struct {
	WindowWidth  int
	WindowHeight int

	Note        Note
	CentsOff    float64
	Frequency   float64
	BlockLength int

	CurrentState State

	// perhaps these should be maps
	// make the map while reading the data so it points the name to the index while reading
	Theme          ColorTheme
	AsciiArt       string
	SelectedTuning tuning.Tuning

	Options          []SettingsOptions
	SettingsData     SettingsData
	SettingsSelected AppSettings

	SelectedOption      int
	SelectedOptionValue int
	SelectingValues     bool
}

func NewModel() Model {
	m := Model{
		BlockLength:      BL,
		CurrentState:     Initializing,
		SettingsSelected: LoadSettingsSelections(),
		SettingsData:     LoadSettingsData(),
	}

	m.ApplySettings()
	m.Options = DefineSettingsOptions(m.SettingsData, m.SettingsSelected)

	return m
}

func (m *Model) ApplySettings() {
	// m.SettingsSelected = LoadSettingsSelections()

	m.AsciiArt = m.SettingsData.AsciiArt[m.SettingsSelected.AsciiArt].FileContents

	SetBorderStyle(m.SettingsData.BorderStyles[m.SettingsSelected.BorderStyle])

	m.SelectedTuning = m.SettingsData.Tunings[m.SettingsSelected.Tuning]

	m.Theme = m.SettingsData.ColorThemes[m.SettingsSelected.ColorTheme]
	m.Theme.SetToCurrent()

	// Store settings to json file
	StoreSettings(m.SettingsSelected)
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetToCurrent(),
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
	case ReRenderMsg:
		m.Options = DefineSettingsOptions(m.SettingsData, m.SettingsSelected)
		m.ApplySettings()
		return m, nil

	case tea.KeyMsg:
		switch message.String() {
		case "ctrl+c", "q":
			seq := tea.Sequence(closeAudioStream(), tea.Quit)
			cmds = append(cmds, seq)

		case "?":
			m.CurrentState = Help

		case "backspace", "esc":
			m.CurrentState = Listening
			m.SelectingValues = false
			cmds = append(cmds, CalculateNote())

		case "s":
			m.CurrentState = Settings

		case "h", "left":
			m.SelectingValues = false

		case "j", "down":
			if !m.SelectingValues {
				m.SelectedOption = (m.SelectedOption + 1) % len(m.Options)
				// UX fix: hover over the currently saved option when switching categories
				m.SelectedOptionValue = m.Options[m.SelectedOption].Selected
			} else {
				m.SelectedOptionValue = (m.SelectedOptionValue + 1) % len(m.Options[m.SelectedOption].Options)
			}

		case "k", "up":
			if !m.SelectingValues {
				m.SelectedOption = (m.SelectedOption - 1 + len(m.Options)) % len(m.Options)
				// UX fix: hover over the currently saved option when switching categories
				m.SelectedOptionValue = m.Options[m.SelectedOption].Selected
			} else {
				m.SelectedOptionValue = (m.SelectedOptionValue - 1 + len(m.Options[m.SelectedOption].Options)) % len(m.Options[m.SelectedOption].Options)
			}

		case "l", "right":
			m.SelectingValues = true

		case "enter", "space":
			if m.CurrentState == Settings && m.SelectingValues {
				m.SettingsSelected = m.Options[m.SelectedOption].Apply(m.SelectedOptionValue, m.SettingsSelected)

				cmds = append(cmds, ReRender)
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
