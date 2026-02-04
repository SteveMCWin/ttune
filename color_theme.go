package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ColorTheme struct {
	Name      string `json:"name"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`
}

var (
	boxStyle      = lipgloss.NewStyle().Padding(1).Margin(0, 1).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	asciiArtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	ColorThemes   map[string]ColorTheme
)

func (t ColorTheme) SetToCurrent(isDark bool) func() tea.Msg {
	return func() tea.Msg {
		boxStyle = boxStyle.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Tertiary))
		// asciiArtStyle = asciiArtStyle.Foreground(lightDark())
		return nil
	}
}

func SetBorderStyle(style string) {
	switch style {
	case "Normal":
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
	case "Rounded":
		boxStyle = boxStyle.BorderStyle(lipgloss.RoundedBorder())
	case "Double":
		boxStyle = boxStyle.BorderStyle(lipgloss.DoubleBorder())
	case "Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.BlockBorder())
	case "Inner Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.InnerHalfBlockBorder())
	case "Outer Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.OuterHalfBlockBorder())
	case "Thick":
		boxStyle = boxStyle.BorderStyle(lipgloss.ThickBorder())
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
	}
}
