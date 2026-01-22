package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ColorTheme struct {
	Name string `json:"name"`
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Tertiary     string `json:"tertiary"`
}

var (
	boxStyle = lipgloss.NewStyle().Padding(1).Margin(0, 1).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	asciiArtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	ColorThemes map[string]ColorTheme
)

// func init() {
// 	ColorThemes = make(map[string]ColorTheme)
// 	ColorThemes["default"] = ColorTheme{
// 		Primary:    lipgloss.Color("#1e1e2e"),
// 		Secondary:  lipgloss.Color("#6c7086"),
// 		Tertiary:     lipgloss.Color("#89b4fa"),
// 	}
//
// 	ColorThemes["white"] = ColorTheme{
// 		Primary:    lipgloss.Color("#ffffff"),
// 		Secondary:  lipgloss.Color("#ffffff"),
// 		Tertiary:     lipgloss.Color("#ffffff"),
// 	}
// }

func (t ColorTheme) SetToCurrent(isDark bool) func() tea.Msg {
	return func() tea.Msg {
		boxStyle = boxStyle.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Tertiary))
		// asciiArtStyle = asciiArtStyle.Foreground(lightDark())
		return nil
	}
}

func SetBorderStyle(style string) {
	switch style {
	case "normal":
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
	case "rounded":
		boxStyle = boxStyle.BorderStyle(lipgloss.RoundedBorder())
	case "double":
		boxStyle = boxStyle.BorderStyle(lipgloss.DoubleBorder())
	case "block":
		boxStyle = boxStyle.BorderStyle(lipgloss.BlockBorder())
	case "inner_block":
		boxStyle = boxStyle.BorderStyle(lipgloss.InnerHalfBlockBorder())
	case "outer_block":
		boxStyle = boxStyle.BorderStyle(lipgloss.OuterHalfBlockBorder())
	case "thick":
		boxStyle = boxStyle.BorderStyle(lipgloss.ThickBorder())
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
	}
}
