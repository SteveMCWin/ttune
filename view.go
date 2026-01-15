package main

import (

	"fmt"
	"tuner/tuning"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

func (m Model) View() tea.View {

	boxStyle = boxStyle.Width(m.WindowWidth - boxStyle.GetHorizontalMargins()).Height(TITLE_HEIGHT)
	title_height := TITLE_HEIGHT + boxStyle.GetVerticalFrameSize()
	title_box := boxStyle.Render("tTuner - " + string(m.CurrentState))

	whole_widht := m.WindowWidth - boxStyle.GetHorizontalFrameSize()
	tuning_width := whole_widht/8
	meter_box_width := whole_widht - tuning_width + 2

	if meter_box_width <= 0 {
		return tea.NewView("")
	}

	boxStyle = boxStyle.Width(tuning_width).Height(m.WindowHeight - title_height+1)
	tuning_contents := "Tuning: " + m.SelectedTuning.Name + "\n\n\n"
	for _, t := range m.SelectedTuning.Notes {
		tuning_contents = tuning_contents + t + "\n\n"
	}
	tuning_contents = tuning_contents[:len(tuning_contents)-2]
	tuning_box := boxStyle.Render(tuning_contents)


	meter_content_width := (int)(0.8 * (float32)(meter_box_width - boxStyle.GetHorizontalFrameSize()))
	meter_notes_arr := make([]byte, meter_content_width)
	for i := range meter_content_width {
		meter_notes_arr[i] = ' '
	}

	prev_note := prevNote(m.Note)
	next_note := nextNote(m.Note)
	curr_full_note := []byte(fmt.Sprintf("%2s %-2d", tuning.NoteNames[m.Note.Index], m.Note.Octave))
	prev_full_note := []byte(fmt.Sprintf("%2s %-2d", tuning.NoteNames[prev_note.Index], prev_note.Octave))
	next_full_note := []byte(fmt.Sprintf("%2s %-2d", tuning.NoteNames[next_note.Index], next_note.Octave))

	for i := range len(curr_full_note) {
		meter_notes_arr[i] = prev_full_note[i]
		meter_notes_arr[meter_content_width/2+i-2] = curr_full_note[i]
		meter_notes_arr[meter_content_width-len(next_full_note)-1+i] = next_full_note[i]
	}

	meter_bar_arr := make([]rune, meter_content_width)
	for i := range meter_bar_arr {
		meter_bar_arr[i] = '─'
	}
	meter_bar_arr[0] = '├'
	meter_bar_arr[len(meter_bar_arr)-1] = '┤'
	meter_bar_arr[len(meter_bar_arr)/2] = '┴'

	meter_content := lipgloss.JoinVertical(lipgloss.Center, string(meter_notes_arr), string(meter_bar_arr))

	boxStyle = boxStyle.Width(meter_box_width).Height(m.WindowHeight - title_height+1)
	meter_box := boxStyle.Render(meter_content)

	main_content := lipgloss.JoinHorizontal(lipgloss.Center, tuning_box, meter_box)

	instructions := lipgloss.NewStyle().Foreground(m.Theme.TextUnyped).Align(lipgloss.Center, lipgloss.Top).Margin(0, 0).Render("q - quit   s - settings   h - help")
	all_contents := lipgloss.JoinVertical(lipgloss.Left, title_box, main_content)
	all_contents = lipgloss.JoinVertical(lipgloss.Center, all_contents, instructions)

	view := tea.NewView(all_contents)
	view.AltScreen = true
	return view
}
