package main

import (
	"log"
	"net/http"
	"strings"

	"tuner/tuning"

	tea "charm.land/bubbletea/v2"
)

type Model struct {
	// windowWidth  int
	// windowHeight int

	Theme    ColorTheme
	// err      error

	Buffer []float32
	Buffer64 []float64
	BlockLength int
	Window []float64
	Frequency float64

	Tuning tuning.Tuning
}

func NewModel() Model {
	return Model{
		Tuning: tuning.Tunings["standard"],
		Theme: DefaultTheme,
	}
}

func LoadLocalData() tea.Msg {
	return "todo"
}

func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.Theme.SetCurrentTheme(true),                                   // NOTE: hard coded for testing
	}

	return tea.Batch(cmds...) // NOTE: set curr theme should be replaced with a function that loads save data and that handles the theme
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch key := msg.Key(); key.Code {
		case tea.KeyRight, tea.KeyTab:
			if !m.isTyping {
				m.currTab = TabIndex((int(m.currTab) + 1) % len(m.tabs))
			}
		case tea.KeyLeft, tea.KeyLeftShift | tea.KeyTab:
			if !m.isTyping {
				m.currTab = TabIndex((len(m.tabs) + int(m.currTab) - 1) % len(m.tabs))
			}
		case tea.KeyEnter:
			if m.currTab == Home {
				m.isTyping = true
				// cmds = append(cmds, tea.ShowCursor)
			}
		case tea.KeyEsc:
			if m.isTyping {
				// cmds = append(cmds, tea.HideCursor)
				m.isTyping = false
				// stop the test or something
			}
		default:
			switch msg.String() {
			case "ctrl+c":
				seq := tea.Sequence(ChangeFontSize(&m.terminal, 0, true), tea.Quit)
				cmds = append(cmds, seq)
			case "ctrl+r":
				ResetTyingData(&m)
				cmds = append(cmds, GetQuoteFromServer(mod.QUOTE_MEDIUM))
			case "ctrl+up":
				cmds = append(cmds, ChangeFontSize(&m.terminal, 1, true))
			case "ctrl+down":
				cmds = append(cmds, ChangeFontSize(&m.terminal, 1, false))
			default:
				if m.isTyping {
					HandleTyping(&m, msg.String())
				}
			}
		}
	case HttpStatus:
		if int(msg) == http.StatusOK {
			m.isOnline = true
		}
	case HttpError:
		log.Println("ERROR:", msg)
	case tea.WindowSizeMsg:
		log.Println("Terminal width:", msg.Width)
		log.Println("Terminal height:", msg.Height)
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.UpdateCursorStartPos()
		m.CalcCRowEnds()
	case mod.Quote:
		m.cursor.X = m.cStartX
		m.cursor.Y = m.cStartY
		log.Println("startX:", m.cStartX)
		log.Println("startY:", m.cStartY)
		m.quoteLoaded = true
		m.quote = msg
		m.splitQuote = strings.Split(m.quote.Quote, " ")
		m.typedLen = 0
		m.CalcCRowEnds()
	case SupportedTerminals:
		m.terminal = msg
	}

	return m, tea.Batch(cmds...)
}

