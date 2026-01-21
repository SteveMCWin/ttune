package main

import (
	"image/color"
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const TITLE_HEIGHT = 3

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
	boxStyle = lipgloss.NewStyle().Padding(1).Margin(0, 1).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	ColorThemes map[string]ColorTheme
)

func init() {
	ColorThemes = make(map[string]ColorTheme)
	ColorThemes["default"] = ColorTheme{
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

	ColorThemes["white"] = ColorTheme{
		Primary:    lipgloss.Color("#ffffff"),
		Secondary:  lipgloss.Color("#ffffff"),
		Accent:     lipgloss.Color("#ffffff"),
		TextError:  lipgloss.Color("#ffffff"),
		TextTyped:  lipgloss.Color("#ffffff"),
		TextUnyped: lipgloss.Color("#ffffff"),

		PrimaryLight:    lipgloss.Color("#ffffff"),
		SecondaryLight:  lipgloss.Color("#ffffff"),
		AccentLight:     lipgloss.Color("#ffffff"),
		TextErrorLight:  lipgloss.Color("#ffffff"),
		TextTypedLight:  lipgloss.Color("#ffffff"),
		TextUnypedLight: lipgloss.Color("#ffffff"),
	}
}

func (t ColorTheme) SetToCurrent(isDark bool) func() tea.Msg {
	var lightDark = lipgloss.LightDark(isDark)
	return func() tea.Msg {
		boxStyle = boxStyle.BorderForeground(lightDark(t.Accent, t.AccentLight)).Foreground(lightDark(t.Secondary, t.SecondaryLight))
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
