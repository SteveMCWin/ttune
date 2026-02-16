package main

import (
	"fmt"
	"strings"
	"ttune/tuning"

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

	title_box_style := boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_height := TITLE_HEIGHT + title_box_style.GetVerticalFrameSize()
	title_box := title_box_style.Render("tTuner - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalMargins()
	tuning_width := whole_widht/4
	if m.AsciiArt == "" {
		tuning_width /= 2
	}
	meter_box_width := whole_widht - tuning_width - boxStyle.GetHorizontalMargins()

	if meter_box_width <= 0 {
		return ""
	}

	tuning_box_style := boxStyle.Width(tuning_width).Height(m.WindowHeight - title_height+1)
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
	tuning_box := tuning_box_style.Render(tuning_contents)


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

	meter_box_style := boxStyle.Width(meter_box_width).Height(m.WindowHeight - title_height+1)
	meter_box := meter_box_style.Render(meter_content)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, tuning_box, meter_box)

	instructions_str := "? - help   s - settings   q - quit"
	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color(m.Theme.Secondary)).Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render(instructions_str)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	return all_contents
}

// func createSettingsContents(m Model) string {
//
// 	title_box_style := boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
// 	title_height := TITLE_HEIGHT + title_box_style.GetVerticalFrameSize()
// 	title_box := title_box_style.Render("tTuner - " + string(m.CurrentState))
//
// 	whole_widht := m.WindowWidth - boxStyle.GetHorizontalFrameSize()
// 	settings_width := whole_widht/4
// 	available_values_width := settings_width
// 	if m.AsciiArt == "" {
// 		settings_width /= 2
// 	}
// 	setting_val_box_width := whole_widht - settings_width
//
// 	if setting_val_box_width <= 0 {
// 		return ""
// 	}
//
// 	contents_height := m.WindowHeight - title_height+1
// 	settings_box_style := boxStyle.Width(settings_width).Height(contents_height)
// 	setting_names := make([]string, 0)
// 	for _, s := range m.Options {
// 		setting_names = append(setting_names, s.Name)
// 	}
//
// 	settings_box := settings_box_style.Align(lipgloss.Left, lipgloss.Top).Render(lipgloss.JoinVertical(lipgloss.Left, setting_names...))
//
// 	// TODO: Gotta put these in their own boxes
// 	available_options := lipgloss.JoinVertical(lipgloss.Left, m.Options[m.SelectedOption].Options...)
// 	box_available_options := settingsBox.Align(lipgloss.Left, lipgloss.Top).Width(available_values_width).Height(contents_height*3/4).Render(available_options)
// 	option_description := m.Options[m.SelectedOption].Description
// 	box_option_description := settingsBox.Align(lipgloss.Left, lipgloss.Top).Width(setting_val_box_width-available_values_width).Height(contents_height*3/4).Render(option_description)
// 	option_preview := m.Options[m.SelectedOption].Previews[m.SelectedOptionValue]
// 	box_option_preview := settingsBox.Align(lipgloss.Center, lipgloss.Top).Width(setting_val_box_width-3).Height(contents_height/4).Render(option_preview)
//
// 	val_prev := lipgloss.JoinHorizontal(lipgloss.Top, box_available_options, box_option_preview)
//
// 	boxStyle = boxStyle.Width(setting_val_box_width).Height(m.WindowHeight - title_height+1)
// 	settings_content := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, val_prev, box_option_description))
//
// 	instructions_str := "backspace/esc - back   ↓/j - down   ↑/k - up   ←/h - left   →/l - right   enter/space - select   q - quit"
// 	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color(m.Theme.Secondary)).Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render(instructions_str)
//
// 	main_content := lipgloss.JoinHorizontal(lipgloss.Center, settings_box, settings_content)
//
// 	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
// 	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)
//
// 	return all_contents
// }

func createSettingsContents(m Model) string {
	title_box_style := boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_box := title_box_style.Render("tTuner - " + string(m.CurrentState))

	whole_title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	options_height := m.WindowHeight - whole_title_height + 1

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalMargins()

	if whole_widht <= 0 {
		return ""
	}

	settings_width := whole_widht/4
	// available_values_width := settings_width
	if m.AsciiArt == "" {
		settings_width /= 2
	}
	options_box_width := whole_widht - settings_width - boxStyle.GetHorizontalMargins()

	if options_box_width <= 0 {
		return ""
	}

	settings_box_style := boxStyle.Width(settings_width).Height(options_height)
	setting_names := make([]string, 0)
	for _, s := range m.Options {
		setting_names = append(setting_names, s.Name)
	}

	settings_box := settings_box_style.Align(lipgloss.Left, lipgloss.Top).Render(lipgloss.JoinVertical(lipgloss.Left, setting_names...))

	options_box_style := boxStyle.Width(options_box_width).Height(options_height).Align(lipgloss.Left, lipgloss.Top).Padding(0)

	available_options_widht := (options_box_width - options_box_style.GetHorizontalFrameSize() + options_box_style.GetHorizontalBorderSize())/2
	available_options_height := (options_height - options_box_style.GetVerticalFrameSize() + options_box_style.GetHorizontalBorderSize())/2

	available_options := lipgloss.JoinVertical(lipgloss.Left, m.Options[m.SelectedOption].Options...)
	box_available_options := settingsBox.Align(lipgloss.Left, lipgloss.Top).PaddingLeft(2).Width(available_options_widht).Height(available_options_height).Render(available_options)

	
	options_description_width := available_options_widht
	options_description_height := options_height-available_options_height - settingsBox.GetVerticalFrameSize()

	option_description := m.Options[m.SelectedOption].Description
	box_option_description := settingsBox.Align(lipgloss.Left, lipgloss.Top).Padding(1, 2).Width(options_description_width).Height(options_description_height).Render(option_description)


	preview_width := options_box_width - available_options_widht - options_box_style.GetHorizontalFrameSize()
	preview_height := options_height - settingsBox.GetVerticalFrameSize()

	option_preview := lipgloss.JoinVertical(lipgloss.Center, "Preview", "", "", lipgloss.JoinHorizontal(lipgloss.Center, "", m.Options[m.SelectedOption].Previews[m.SelectedOptionValue]))
	box_option_preview := settingsBox.Align(lipgloss.Center, lipgloss.Top).Width(preview_width).Height(preview_height).Render(option_preview)

	options_box := options_box_style.Render(lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, box_available_options, box_option_description), box_option_preview))

	instructions_str := "backspace/esc - back   ↓/j - down   ↑/k - up   ←/h - left   →/l - right   enter/space - select   q - quit"
	instructions := lipgloss.NewStyle().Foreground(lipgloss.Color(m.Theme.Secondary)).Faint(true).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render(instructions_str)

	settings_options := lipgloss.JoinHorizontal(lipgloss.Top, settings_box, options_box)
	main_contents := lipgloss.JoinVertical(lipgloss.Center, settings_options, instructions)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_contents)

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
