package main

import (
	"log"

	// tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type ColorTheme struct {
	Name      string `json:"name"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`
}

type AsciiArt struct {
	FileName     string
	FileContents string
}

var (
	boxStyle = lipgloss.NewStyle().Padding(1).Margin(0, 1).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	// settingsBoxStyle   = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	// settingsBoxStyle   = lipgloss.NewStyle().Border(lipgloss.BlockBorder())
	settingsBoxStyle     = lipgloss.NewStyle().Padding(1).Border(lipgloss.DoubleBorder())
	selectedStyle        = lipgloss.NewStyle()
	instructionsStyle    = lipgloss.NewStyle().Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0)
	inactive_arrow_style = lipgloss.NewStyle().Faint(true).Faint(true)
	active_arrow_style   = lipgloss.NewStyle().Faint(false)
	// asciiArtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	ColorThemes map[string]ColorTheme
)

func (t ColorTheme) SetToCurrent() {
	boxStyle = boxStyle.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Secondary))
	settingsBoxStyle = settingsBoxStyle.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Secondary))
	selectedStyle = selectedStyle.Foreground(lipgloss.Color(t.Tertiary))
	instructionsStyle = instructionsStyle.Foreground(lipgloss.Color(t.Tertiary))
	inactive_arrow_style = inactive_arrow_style.Foreground(lipgloss.Color(t.Tertiary))
	active_arrow_style = active_arrow_style.Foreground(lipgloss.Color(t.Tertiary))
	// asciiArtStyle = asciiArtStyle.Foreground(lightDark())
}

func SetBorderStyle(style string) {
	switch style {
	case "Normal":
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.NormalBorder())
	case "Rounded":
		boxStyle = boxStyle.BorderStyle(lipgloss.RoundedBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.RoundedBorder())
	case "Double":
		boxStyle = boxStyle.BorderStyle(lipgloss.DoubleBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.DoubleBorder())
	case "Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.BlockBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.BlockBorder())
	case "Inner Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.InnerHalfBlockBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.InnerHalfBlockBorder())
	case "Outer Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.OuterHalfBlockBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.OuterHalfBlockBorder())
	case "Thick":
		boxStyle = boxStyle.BorderStyle(lipgloss.ThickBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.ThickBorder())
	case "Ascii":
		boxStyle = boxStyle.BorderStyle(lipgloss.ASCIIBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.ASCIIBorder())
	case "None":
		boxStyle = boxStyle.BorderStyle(lipgloss.HiddenBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.HiddenBorder())
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
		settingsBoxStyle = settingsBoxStyle.BorderStyle(lipgloss.NormalBorder())
	}
}

func GetBorderStyleByName(style_name string) lipgloss.Style {
	style := lipgloss.NewStyle()
	switch style_name {
	case "Normal":
		style = style.BorderStyle(lipgloss.NormalBorder())
	case "Rounded":
		style = style.BorderStyle(lipgloss.RoundedBorder())
	case "Double":
		style = style.BorderStyle(lipgloss.DoubleBorder())
	case "Block":
		style = style.BorderStyle(lipgloss.BlockBorder())
	case "Inner Block":
		style = style.BorderStyle(lipgloss.InnerHalfBlockBorder())
	case "Outer Block":
		style = style.BorderStyle(lipgloss.OuterHalfBlockBorder())
	case "Thick":
		style = style.BorderStyle(lipgloss.ThickBorder())
	case "None":
		style = style.BorderStyle(lipgloss.HiddenBorder())
	case "Ascii":
		style = style.BorderStyle(lipgloss.ASCIIBorder())
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		style = style.BorderStyle(lipgloss.NormalBorder())
	}

	return style
}
