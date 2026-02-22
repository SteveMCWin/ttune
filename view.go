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
	title_box := title_box_style.Render("tTune - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalMargins()
	tuning_width := whole_widht / 4
	if m.AsciiArt == "" {
		tuning_width = tuning_width*2/3
	}
	meter_box_width := whole_widht - tuning_width - boxStyle.GetHorizontalMargins()

	if meter_box_width <= 0 {
		return ""
	}

	tuning_box_style := boxStyle.Width(tuning_width).Height(m.WindowHeight - title_height + 1)
	tuning_name := "Tuning: " + m.SelectedTuning.Name + "\n\n\n"
	var tuning_contents string
	if m.AsciiArt != "" {
		ascii_art := m.AsciiArt
		for i := len(m.SelectedTuning.Notes) - 1; i >= 0; i-- {
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

	meter_content_width := (int)(0.8 * (float32)(meter_box_width-boxStyle.GetHorizontalFrameSize()))
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

	accuracy_pos_offset := m.Note.CentsOff / 10 // NOTE: Perhaps make this a precision setting
	accuracy_pos := max(0, min(len(meter_bar_arr)/2+accuracy_pos_offset, len(meter_bar_arr)-1))
	acc_indicator_arr[accuracy_pos] = '^'

	meter_content := lipgloss.JoinVertical(lipgloss.Center, string(meter_notes_arr), string(meter_bar_arr), string(acc_indicator_arr))

	meter_box_style := boxStyle.Width(meter_box_width).Height(m.WindowHeight - title_height + 1)
	meter_box := meter_box_style.Render(meter_content)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, tuning_box, meter_box)

	instructions_str := "? - help   s - settings   q - quit"
	instructions := instructionsStyle.Render(instructions_str)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	return all_contents
}

func createSettingsContents(m Model) string {
	title_box_style := boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_box := title_box_style.Render("tTune - " + string(m.CurrentState))

	whole_title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	settings_height := m.WindowHeight - whole_title_height + 1

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalMargins()

	if whole_widht <= 0 {
		return ""
	}

	settings_width := whole_widht / 4

	options_box_width := whole_widht - settings_width

	if options_box_width <= 0 {
		return ""
	}

	settings_box_style := boxStyle.Width(settings_width).Height(settings_height)
	setting_names := []string{"Settings", ""}
	for i, o := range m.VisualOptions {
		var line string
		if m.SelectedOption != i {
			selection_box := "[ ] "
			line = selection_box + o.Name
		} else {
			selection_box := "[o] "
			line = selection_box + o.Name
			if !m.SelectingValues {
				line = selectedStyle.Render(line)
			}
		}
		setting_names = append(setting_names, line)
	}

	settings_box := settings_box_style.Align(lipgloss.Left, lipgloss.Top).PaddingLeft(2).Render(lipgloss.JoinVertical(lipgloss.Left, setting_names...))

	options_box_style := boxStyle.UnsetBorderStyle().Width(options_box_width).Height(settings_height).Align(lipgloss.Left, lipgloss.Top).Padding(0)

	available_options_widht := (options_box_width - options_box_style.GetHorizontalFrameSize()) / 3
	available_options_height := (settings_height - options_box_style.GetVerticalFrameSize()) / 2

	options_names := []string{"Options", ""}
	for i, o := range m.VisualOptions[m.SelectedOption].Options {

		prefix := "[ ] "
		if m.VisualOptions[m.SelectedOption].Selected == i {
			prefix = "[o] "
		}
		line := prefix + o

		// Is the user currently hovering over this setting?
		if m.SelectedOptionValue == i && m.SelectingValues {
			line = selectedStyle.Render(line)
		}

		options_names = append(options_names, line)
	}
	available_options := lipgloss.JoinVertical(lipgloss.Left, options_names...)
	box_available_options := settingsBox.Align(lipgloss.Left, lipgloss.Top).PaddingLeft(2).Width(available_options_widht).Height(available_options_height).Render(available_options)

	options_description_width := available_options_widht
	options_description_height := settings_height - available_options_height

	option_description := lipgloss.JoinVertical(lipgloss.Left, "Description", "", m.VisualOptions[m.SelectedOption].Description)
	box_option_description := settingsBox.Align(lipgloss.Left, lipgloss.Top).Padding(1, 2).Width(options_description_width).Height(options_description_height).Render(option_description)

	preview_width := options_box_width - available_options_widht - options_box_style.GetHorizontalFrameSize()
	preview_height := settings_height

	option_preview := lipgloss.JoinVertical(lipgloss.Center, "Preview", "", "", m.VisualOptions[m.SelectedOption].Previews[m.SelectedOptionValue])
	box_option_preview := settingsBox.Align(lipgloss.Center, lipgloss.Center).MarginLeft(2).Width(preview_width - 2).Height(preview_height).Render(option_preview)

	options_box := options_box_style.UnsetBorderStyle().Render(lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.JoinVertical(lipgloss.Left, box_available_options, box_option_description), box_option_preview))

	instructions_str := "backspace/esc - back   ↓/j - down   ↑/k - up   ←/h - left   →/l - right   enter/space - select   ? - help   q - quit"
	instructions := instructionsStyle.Render(instructions_str)

	settings_options := lipgloss.JoinHorizontal(lipgloss.Top, settings_box, options_box)
	main_contents := lipgloss.JoinVertical(lipgloss.Center, settings_options, instructions)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_contents)

	return all_contents
}

func createHelpContents(m Model) string {

	title_box_style := boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_height := TITLE_HEIGHT + title_box_style.GetVerticalFrameSize()
	title_box := title_box_style.Render("tTune - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalMargins()
	help_list_width := whole_widht / 4
	help_contents_width := whole_widht - help_list_width - boxStyle.GetHorizontalMargins()

	if help_contents_width <= 0 {
		return ""
	}

	help_list_box_style := boxStyle.Width(help_list_width).Height(m.WindowHeight - title_height + 1).Align(lipgloss.Left, lipgloss.Top).Padding(1, 2)
	help_list_contents := []string{"Help", ""}
	for i, o := range m.HelpItems {
		var line string
		if m.SelectedHelpItem != i {
			selection_box := "[ ] "
			line = selection_box + o.Name
		} else {
			selection_box := "[o] "
			line = selection_box + o.Name
			if !m.SelectingValues {
				line = selectedStyle.Render(line)
			}
		}
		help_list_contents = append(help_list_contents, line)
	}

	help_list_box := help_list_box_style.Render(lipgloss.JoinVertical(lipgloss.Left, help_list_contents...))

	help_contents_box_style := boxStyle.Width(help_contents_width).Height(m.WindowHeight - title_height + 1).Padding(0, 2).Align(lipgloss.Left)
	help_contents_box := help_contents_box_style.Render(m.HelpItems[m.SelectedHelpItem].Contents)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, help_list_box, help_contents_box)

	instructions_str := "backspace/esc - back   s - settings   q - quit"
	instructions := instructionsStyle.Render(instructions_str)
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	return all_contents
}
