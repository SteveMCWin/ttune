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

	Settings     []Setting
	// FunctionalOptions []Setting // TODO
	SettingsData      SettingsData
	UserSettingsData  SettingsData
	SettingsSelected  SettingsSelections

	SelectedOption      int
	SelectedOptionValue int
	SelectingValues     bool

	HelpItems []HelpItem

	SelectedHelpItem int
}

func NewModel() Model {
	settingsData := LoadSettingsData()
	m := Model{
		BlockLength:      BL,
		CurrentState:     Initializing,
		SettingsData:     settingsData,
		SettingsSelected: LoadSettingsSelections("selections.json", settingsData),
		HelpItems:        InitHelpItems(),
	}

	m.Settings = DefineAllSettingsOptions(m.SettingsData, m.SettingsSelected)
	m.ApplySettings()

	// Force the tui to render the selected preview on startup
	m.SelectedOptionValue = m.Settings[0].SelectedIdx()

	return m
}

func (m *Model) ApplySettings() {
	m.AsciiArt = m.SettingsData.AsciiArt[m.Settings[0].SelectedIdx()].FileContents
	SetBorderStyle(m.SettingsData.BorderStyles[m.Settings[1].SelectedIdx()])
	m.Theme = m.SettingsData.ColorThemes[m.Settings[2].SelectedIdx()]
	m.SelectedTuning = m.SettingsData.Tunings[m.Settings[3].SelectedIdx()]
	m.Theme.SetToCurrent()

	// Store settings to json file
	StoreSettings(m.SettingsSelected, "selections.json")
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
		m.ApplySettings()
		return m, nil

	case tea.KeyMsg:
		// When a text input is focused, route all keys to it instead of normal navigation
		if m.CurrentState == Settings && m.SelectingValues {
			opt := m.Settings[m.SelectedOption].Options[m.SelectedOptionValue]
			if opt.IsFocused() {
				switch message.String() {
				case "esc":
					opt.HanldeSelect()
				case "enter":
					m.SettingsSelected = m.Settings[m.SelectedOption].Apply(opt.GetValue(), m.SettingsSelected)
					m.Settings[m.SelectedOption].Selected = opt.GetValue()
					opt.HanldeSelect()
					cmds = append(cmds, ReRender)
				default:
					if cmd := opt.Update(message); cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
				return m, tea.Batch(cmds...)
			}
		}

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
			switch m.CurrentState {
			case Settings:
				if !m.SelectingValues {
					m.SelectedOption = min(m.SelectedOption+1, len(m.Settings)-1)
					m.SelectedOptionValue = m.Settings[m.SelectedOption].SelectedIdx()
				} else {
					m.SelectedOptionValue = min(
						m.SelectedOptionValue+1,
						len(m.Settings[m.SelectedOption].Options)-1,
					)
				}
			case Help:
				m.SelectedHelpItem = min(m.SelectedHelpItem+1, len(m.HelpItems)-1)
			}

		case "k", "up":
			switch m.CurrentState {
			case Settings:
				if !m.SelectingValues {
					m.SelectedOption = max(m.SelectedOption-1, 0)
					m.SelectedOptionValue = m.Settings[m.SelectedOption].SelectedIdx()
				} else {
					m.SelectedOptionValue = max(m.SelectedOptionValue-1, 0)
				}
			case Help:
				m.SelectedHelpItem = max(m.SelectedHelpItem-1, 0)
			}

		case "l", "right":
			if m.CurrentState == Settings {
				m.SelectingValues = true
			}

		case "enter", "space":
			if m.CurrentState == Settings && m.SelectingValues {
				option_selected := m.Settings[m.SelectedOption].Options[m.SelectedOptionValue]
				option_selected.HanldeSelect()
				if !option_selected.IsFocused() {
					// MultiChoiceOption: not focusable, apply the selection immediately
					m.SettingsSelected = m.Settings[m.SelectedOption].Apply(
						option_selected.GetValue(),
						m.SettingsSelected,
					)
					m.Settings[m.SelectedOption].Selected = option_selected.GetValue()
					cmds = append(cmds, ReRender)
				}
				// InputFieldOption: HanldeSelect just focused it — wait for next Enter to apply
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
