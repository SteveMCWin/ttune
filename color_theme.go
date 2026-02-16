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
	FileName string
	FileContents string
}

var (
	boxStyle      = lipgloss.NewStyle().Padding(1).Margin(0, 1).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder())
	// settingsBox   = lipgloss.NewStyle().Border(lipgloss.HiddenBorder())
	// settingsBox   = lipgloss.NewStyle().Border(lipgloss.BlockBorder())
	settingsBox   = lipgloss.NewStyle().Padding(1).Border(lipgloss.DoubleBorder())
	selectedStyle = lipgloss.NewStyle()
	// asciiArtStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	ColorThemes   map[string]ColorTheme
)

func (t ColorTheme) SetToCurrent() {
	boxStyle = boxStyle.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Tertiary))
	settingsBox = settingsBox.BorderForeground(lipgloss.Color(t.Primary)).Foreground(lipgloss.Color(t.Tertiary))
	selectedStyle = selectedStyle.Foreground(lipgloss.Color(t.Secondary))
	// asciiArtStyle = asciiArtStyle.Foreground(lightDark())
}

func SetBorderStyle(style string) {
	switch style {
	case "Normal":
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.NormalBorder())
	case "Rounded":
		boxStyle = boxStyle.BorderStyle(lipgloss.RoundedBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.RoundedBorder())
	case "Double":
		boxStyle = boxStyle.BorderStyle(lipgloss.DoubleBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.DoubleBorder())
	case "Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.BlockBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.BlockBorder())
	case "Inner Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.InnerHalfBlockBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.InnerHalfBlockBorder())
	case "Outer Block":
		boxStyle = boxStyle.BorderStyle(lipgloss.OuterHalfBlockBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.OuterHalfBlockBorder())
	case "Thick":
		boxStyle = boxStyle.BorderStyle(lipgloss.ThickBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.ThickBorder())
	case "None":
		boxStyle = boxStyle.BorderStyle(lipgloss.HiddenBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.HiddenBorder())
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		boxStyle = boxStyle.BorderStyle(lipgloss.NormalBorder())
		settingsBox = settingsBox.BorderStyle(lipgloss.NormalBorder())
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
	default:
		log.Println("! ! ! ! ! Unsupported border style:", style)
		style = style.BorderStyle(lipgloss.NormalBorder())
	}

	return style
}
