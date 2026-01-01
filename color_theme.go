package main

import (
	"image/color"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ColorTheme struct {
	Primary    color.Color
	Secondary  color.Color
	Accent     color.Color
	TextError  color.Color
	TextTyped  color.Color
	TextUnyped color.Color

	PrimaryLight    color.Color
	SecondaryLight  color.Color
	AccentLight     color.Color
	TextErrorLight  color.Color
	TextTypedLight  color.Color
	TextUnypedLight color.Color
}

var (
	DefaultTheme = ColorTheme{
		Primary:    lipgloss.Color("#1e1e2e"),
		Secondary:  lipgloss.Color("#6c7086"),
		Accent:     lipgloss.Color("#89b4fa"),
		TextError:  lipgloss.Color("#dd8888"),
		TextTyped:  lipgloss.Color("#ffffff"),
		TextUnyped: lipgloss.Color("#aaaaaa"),

		PrimaryLight:    lipgloss.Color("#6c7086"),
		SecondaryLight:  lipgloss.Color("#acb0be"),
		AccentLight:     lipgloss.Color("#1e66f5"),
		TextErrorLight:  lipgloss.Color("#dd8888"),
		TextTypedLight:  lipgloss.Color("#000000"),
		TextUnypedLight: lipgloss.Color("#444444"),
	}
)

var (
	inactiveTabBorder = lipgloss.Border{Bottom: "─", BottomLeft: "─", BottomRight: "─"}
	activeTabBorder   = lipgloss.Border{Top: "─", Bottom: " ", Left: "│", Right: "│", TopLeft: "╭", TopRight: "╮", BottomLeft: "┘", BottomRight: "└"}
	tabGapBorderLeft  = lipgloss.Border{Bottom: "─", BottomLeft: "╭", BottomRight: "─"}
	tabGapBorderRight = lipgloss.Border{Bottom: "─", BottomLeft: "─", BottomRight: "╮"}

	docStyle          = lipgloss.NewStyle().Padding(1, 2).Align(lipgloss.Center)
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	tabGapLeft        = inactiveTabStyle.Border(tabGapBorderLeft, true)
	tabGapRight       = inactiveTabStyle.Border(tabGapBorderRight, true)
	windowStyle       = lipgloss.NewStyle().Padding(2).Align(lipgloss.Left, lipgloss.Center).Border(lipgloss.RoundedBorder()).UnsetBorderTop()
	quoteStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextUnyped)
	typedStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextTyped)
	errorStyle        = lipgloss.NewStyle().Foreground(DefaultTheme.TextError)
	contentStyle      = lipgloss.NewStyle().Padding(0, 8)
)

func (t ColorTheme) SetCurrentTheme(isDark bool) func() tea.Msg {
	var lightDark = lipgloss.LightDark(isDark)
	return func() tea.Msg {
		inactiveTabStyle = inactiveTabStyle.BorderForeground(lightDark(t.Accent, t.AccentLight)).Foreground(lightDark(t.Secondary, t.SecondaryLight))
		activeTabStyle = activeTabStyle.BorderForeground(lightDark(t.Accent, t.AccentLight)).Foreground(lightDark(t.Accent, t.AccentLight))
		tabGapLeft = tabGapLeft.BorderForeground(lightDark(t.Accent, t.AccentLight))
		tabGapRight = tabGapRight.BorderForeground(lightDark(t.Accent, t.AccentLight))
		windowStyle = windowStyle.BorderForeground(lightDark(t.Accent, t.AccentLight))
		return nil
	}
}
