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

	VisualOptions     []SettingsOptions
	FunctionalOptions []SettingsOptions // TODO
	SettingsData      SettingsData
	SettingsSelected  AppSettings

	SelectedOption      int
	SelectedOptionValue int
	SelectingValues     bool

	HelpItems []HelpItem

	SelectedHelpItem int
}

func NewModel() Model {
	m := Model{
		BlockLength:      BL,
		CurrentState:     Initializing,
		SettingsSelected: LoadSettingsSelections(),
		SettingsData:     LoadSettingsData(),
		HelpItems:        InitHelpItems(),
	}

	m.VisualOptions = DefineVisualSettingsOptions(m.SettingsData, m.SettingsSelected)
	// m.FunctionalOptions = DefineFunctionalSettingsOptions(m.SettingsData, m.SettingsSelected) // TODO
	m.ApplySettings()

	// Force the tui to render the selected preview on startup
	m.SelectedOptionValue = m.VisualOptions[0].Selected % len(m.VisualOptions[0].Options)

	return m
}

func (m *Model) ApplySettings() {
	ascii_selected := m.SettingsSelected.AsciiArt % len(m.VisualOptions[0].Options)
	m.AsciiArt = m.SettingsData.AsciiArt[ascii_selected].FileContents

	SetBorderStyle(m.SettingsData.BorderStyles[m.SettingsSelected.BorderStyle])

	m.SelectedTuning = m.SettingsData.Tunings[m.SettingsSelected.Tuning]

	m.Theme = m.SettingsData.ColorThemes[m.SettingsSelected.ColorTheme]
	m.Theme.SetToCurrent()

	// Store settings to json file
	StoreSettings(m.SettingsSelected)
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		initAutioStream(),
		ReRender,
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
		m.VisualOptions = DefineVisualSettingsOptions(m.SettingsData, m.SettingsSelected)
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
			// if statement needed to prevent race condition
			if m.CurrentState != Listening {
				m.CurrentState = Listening
				m.SelectingValues = false
				cmds = append(cmds, CalculateNote())
			}

		case "s":
			m.CurrentState = Settings

		case "h", "left":
			if m.CurrentState == Settings {
				m.SelectingValues = false
			}

		case "j", "down":
			if m.CurrentState == Settings {
				if !m.SelectingValues {
					m.SelectedOption = (m.SelectedOption + 1) % len(m.VisualOptions)
					m.SelectedOptionValue = m.VisualOptions[m.SelectedOption].Selected
				} else {
					m.SelectedOptionValue = (m.SelectedOptionValue + 1) % len(m.VisualOptions[m.SelectedOption].Options)
				}
			} else if m.CurrentState == Help {
				m.SelectedHelpItem = (m.SelectedHelpItem + 1) % len(m.HelpItems)
			}

		case "k", "up":
			if m.CurrentState == Settings {
				if !m.SelectingValues {
					m.SelectedOption = (m.SelectedOption - 1 + len(m.VisualOptions)) % len(m.VisualOptions)
					m.SelectedOptionValue = m.VisualOptions[m.SelectedOption].Selected
				} else {
					m.SelectedOptionValue = (m.SelectedOptionValue - 1 + len(m.VisualOptions[m.SelectedOption].Options)) % len(m.VisualOptions[m.SelectedOption].Options)
				}
			} else if m.CurrentState == Help {
				m.SelectedHelpItem = (m.SelectedHelpItem - 1 + len(m.HelpItems)) % len(m.HelpItems)
			}

		case "l", "right":
			if m.CurrentState == Settings {
				m.SelectingValues = true
			}

		case "enter", "space":
			if m.CurrentState == Settings && m.SelectingValues {
				m.SettingsSelected = m.VisualOptions[m.SelectedOption].Apply(m.SelectedOptionValue, m.SettingsSelected)

				cmds = append(cmds, ReRender)
			}
		}

	case tea.WindowSizeMsg:
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
