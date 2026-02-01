package main

import (
	"fmt"
	"strings"
	"tuner/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

const TITLE_HEIGHT = 3

func (m Model) View() tea.View {

	view := tea.NewView("Loading...")
	view.AltScreen = true

	switch m.CurrentState {
	case Listening:
		contents := createListeningContents(m)
		view.SetContent(contents)
	case Settings:
		contents := createSettingsContents(m)
		view.SetContent(contents)
	case Help:
		contents := createHelpContents(m)
		view.SetContent(contents)
	}
	
	return view
}

func createListeningContents(m Model) string {

	boxStyle = boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	title_box := boxStyle.Render("tTuner - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalFrameSize()
	tuning_width := whole_widht/4
	if m.AsciiArt == "" {
		tuning_width /= 2
	}
	meter_box_width := whole_widht - tuning_width + 2

	if meter_box_width <= 0 {
		return ""
	}

	boxStyle = boxStyle.Width(tuning_width).Height(m.WindowHeight - title_height+1)
	tuning_name := "Tuning: " + m.SelectedTuning.Name + "\n\n\n"
	var tuning_contents string
	if m.AsciiArt != "" {
		ascii_art := m.AsciiArt
		for i := len(m.SelectedTuning.Notes)-1; i >= 0; i-- {
			ascii_art = strings.Replace(ascii_art, "%%%", m.SelectedTuning.Notes[i], 1)
		}
		tuning_contents = lipgloss.JoinVertical(lipgloss.Center, tuning_name, lipgloss.JoinHorizontal(lipgloss.Top, "", ascii_art))
	} else {
		tuning_notes := "\n"
		for _, t := range m.SelectedTuning.Notes {
			tuning_notes = tuning_notes + t + "\n\n"
		}
		tuning_notes = tuning_notes[:len(tuning_notes)-2]
		tuning_contents = lipgloss.JoinVertical(lipgloss.Center, tuning_name, tuning_notes)
	}
	tuning_box := boxStyle.Render(tuning_contents)


	meter_content_width := (int)(0.8 * (float32)(meter_box_width - boxStyle.GetHorizontalFrameSize()))
	meter_notes_arr := make([]byte, meter_content_width)
	for i := range meter_content_width {
		meter_notes_arr[i] = ' '
	}

	prev_note := prevNote(m.Note)
	next_note := nextNote(m.Note)

	var curr_full_note []byte
	var prev_full_note []byte
	var next_full_note []byte

	if prev_note.Octave < 0 {
		curr_full_note = fmt.Appendf(curr_full_note, "%2s %-2d", tuning.NoteNames[m.Note.Index], m.Note.Octave)
		prev_full_note = []byte("     ")
		next_full_note = []byte("     ")
	} else {
		curr_full_note = fmt.Appendf(curr_full_note, "%2s%-2d", tuning.NoteNames[m.Note.Index], m.Note.Octave)
		prev_full_note = fmt.Appendf(prev_full_note, "%2s%-2d", tuning.NoteNames[prev_note.Index], prev_note.Octave)
		next_full_note = fmt.Appendf(next_full_note, "%2s%-2d", tuning.NoteNames[next_note.Index], next_note.Octave)
	}

	for i := range len(curr_full_note) {
		meter_notes_arr[i+1] = prev_full_note[i]
		meter_notes_arr[meter_content_width/2+i-1] = curr_full_note[i]
		meter_notes_arr[meter_content_width-len(next_full_note)+i] = next_full_note[i]
	}

	meter_bar_arr := make([]rune, meter_content_width)
	acc_indicator_arr := make([]rune, meter_content_width)
	for i := range meter_bar_arr {
		meter_bar_arr[i] = '─'
		acc_indicator_arr[i] = ' '
	}
	meter_bar_arr[0] = '┎'
	meter_bar_arr[len(meter_bar_arr)-1] = '┒'
	meter_bar_arr[len(meter_bar_arr)/2] = '┰'
	
	accuracy_pos_offset := m.Note.CentsOff/10 // NOTE: Perhaps make this a precision setting
	accuracy_pos := max(0, min(len(meter_bar_arr)/2+accuracy_pos_offset, len(meter_bar_arr)-1))
	acc_indicator_arr[accuracy_pos] = '^'

	meter_content := lipgloss.JoinVertical(lipgloss.Center, string(meter_notes_arr), string(meter_bar_arr), string(acc_indicator_arr))

	boxStyle = boxStyle.Width(meter_box_width).Height(m.WindowHeight - title_height+1)
	meter_box := boxStyle.Render(meter_content)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, tuning_box, meter_box)

	instructions_str := "? - help   s - settings   q - quit"
	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color(m.Theme.Secondary)).Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render(instructions_str)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	return all_contents
}

func createSettingsContents(m Model) string {

	boxStyle = boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	title_box := boxStyle.Render("tTuner - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalFrameSize()
	settings_width := whole_widht/4
	if m.AsciiArt == "" {
		settings_width /= 2
	}
	setting_val_box_width := whole_widht - settings_width + 2

	if setting_val_box_width <= 0 {
		return ""
	}

	boxStyle = boxStyle.Width(settings_width).Height(m.WindowHeight - title_height+1)
	setting_names := []string{ "Ascii Art", "Color Scheme", "Selected Tuning", "Border Style"}
	settings_box := boxStyle.Align(lipgloss.Left, lipgloss.Top).Render(lipgloss.JoinVertical(lipgloss.Left, setting_names...))

	boxStyle = boxStyle.Width(setting_val_box_width).Height(m.WindowHeight - title_height+1)
	settings_vals_box := boxStyle.Render("")

	instructions_str := "backspace - back   ↓/j - down   ↑/k - up   ←/h - left   →/l - right   q - quit"
	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color(m.Theme.Secondary)).Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render(instructions_str)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, settings_box, settings_vals_box)

	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	return all_contents
}

func createHelpContents(m Model) string {

	boxStyle = boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	// title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	// title_box := boxStyle.Render("tTuner - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalFrameSize()
	tuning_width := whole_widht/4
	if m.AsciiArt == "" {
		tuning_width /= 2
	}
	meter_box_width := whole_widht - tuning_width + 2

	if meter_box_width <= 0 {
		return ""
	}


	return ""
}
