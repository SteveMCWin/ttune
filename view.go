package main

import (
	"fmt"
	"log"

	// "os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

func (m Model) View() tea.View {

	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isActive := i == int(m.currTab)
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()
		style = style.Border(border)

		if m.isTyping {
			style = style.Faint(true)
		}

		renderedTabs = append(renderedTabs, style.Render(t.TabName))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	gap_size := max(0, m.windowWidth-lipgloss.Width(row)-12)
	gap_l := tabGapLeft.Render(strings.Repeat(" ", gap_size/2))
	gap_r := tabGapRight.Render(strings.Repeat(" ", gap_size/2+gap_size%2))
	row = lipgloss.JoinHorizontal(lipgloss.Top, gap_l, row, gap_r)
	doc.WriteString(row)
	doc.WriteString("\n")

	var contents string

	switch m.currTab {
	case About:
		contents = m.tabs[m.currTab].Contents
	case Settings:
		contents = m.tabs[m.currTab].Contents
	case Home:
		contents = GetHomeContents(&m)
	case Leaderboard:
		contents = m.tabs[m.currTab].Contents
	case ProfileView:
		contents = m.tabs[m.currTab].Contents
	default:
	}

	// log.Println("windowStyle width: ", windowStyle.GetHorizontalFrameSize())
	// log.Println("docStyle width: ", docStyle.GetHorizontalFrameSize())
	// log.Println("contentStyle width: ", contentStyle.GetHorizontalFrameSize())
	// log.Println()
	// log.Println("windowStyle height: ", windowStyle.GetVerticalFrameSize())
	// log.Println("docStyle height: ", docStyle.GetVerticalFrameSize())
	// log.Println("contentStyle height: ", contentStyle.GetVerticalFrameSize())
	// log.Println("tabStyle height: ", activeTabStyle.GetVerticalFrameSize())
	// log.Println()

	_, err := doc.WriteString(windowStyle.Width(m.windowWidth - windowStyle.GetHorizontalFrameSize()+2).Height(m.windowHeight-windowStyle.GetVerticalFrameSize()).Render(contents))
	if err != nil {
		log.Println("Error displaying window and contents:", err)
	}

	whole_string := docStyle.Render(doc.String())

	view := tea.NewView(whole_string)
	view.Cursor = m.cursor
	view.AltScreen = true
	return view
}

func GetHomeContents(m *Model) string {

	// There are two problems: 
	//  - We have to determine the starting position of the text
	//  - We have to determine when the cursor is supposed to go to the next row of text
	// 
	// Ok so the idea is to split the quote into rows, where each row ends with a newline character and then display it as that.
	// Since the text is aligned to the left, the X coordinate of the start of each row is the same
	// So once the startX + currX of the cursor are greater than the length of the text in currRow, set the cursors currX to startX and increment currY
	// I've tried searching the result string that gets rendered to the terminal to get the index of the first word and position the cursor there,
	// but the result string has a bunch of characters that aren't visible and are for style, so it's way off. The only other solution I could
	// think of is to just calculate the position based on the style margins and paddings, which does work for the startX, but the startY is currenlty
	// variable because of the alignment of the content style. This means I will have to handle the alignment myself by calculating the vertical padding
	// for the content based on the window height.

	// Ok so to get the startX, we just divide all style.GetHorizontalFrameSize with 2 and add them together
	// For the startY however, it's a different story. Calculating the offest from the top we sum the following:
	// (windowStyle-1)/2   we subtract 1 because the top part of the border is unset, so it's not counted in the GetVerticalFrameSize
	// tabStyle+1   we add 1 because the getVerticalFrameSize returns only the space taken up by padding, margins and borders, not the content
	// docStyle/2   the only 'normal' part of the sum
	//
	// To get the bottom offset we add the following:
	// (windowStyle+1)/2   now we add 1 because there is a border on the bottom
	// docStyle/2   same as the top part
	// 
	// Note that we don't add the tab offset, because the tabs are at the top of the screen
	// Now we need to subtract both of the results from the terminal height and we get the height our working area
	// The last part is probably to just divide that by 2, but I'll trial and error that since I don't know how the alignment works exactly

	contents := ""

	if m.quoteLoaded && !m.quoteCompleted {

		curr_word := m.splitQuote[m.wordsTyped]

		lines := ""
		log.Println("Already typed rows: ")
		for i := m.currLine-m.linesToShow/2; i < m.currLine; i++ {
			log.Println("Line", i, ":", m.quoteLines[i])
			lines += m.quoteLines[i] + "\n"
		}
		contents = typedStyle.Render(lines) // Already typed lines

		correct_chars := len(m.typedWord)
		incorrect_chars := len(m.typedErr)
		typed_chars := correct_chars + incorrect_chars

		correctly_typed := curr_word[:min(correct_chars, len(curr_word))]
		incorrectly_typed := curr_word[min(correct_chars, len(curr_word)):min(typed_chars, len(curr_word))]
		yet_to_type := curr_word[min(typed_chars, len(curr_word)):]
		overtyped := m.typedErr[min(len(curr_word)-correct_chars, incorrect_chars):]

		contents += typedStyle.Render(m.quoteLines[m.currLine][:m.typedLen]) // Current line - typed correctly up to current word
		contents += typedStyle.Render(correctly_typed) // Current word - typed correctly
		contents += errorStyle.Render(incorrectly_typed) // Current word - typed incorrectly
		contents += quoteStyle.Render(yet_to_type) // Current word - untyped
		contents += errorStyle.Render(overtyped) // Current word - overtyped
		
		curr_word_lens_combined := len(correctly_typed + incorrectly_typed + yet_to_type + overtyped)

		contents += quoteStyle.Render(m.quoteLines[m.currLine][m.typedLen + curr_word_lens_combined :]) // Rest of the current line

		log.Println("Currently typed row: ", m.quoteLines[m.currLine][:m.typedLen] + correctly_typed + incorrectly_typed + yet_to_type + overtyped + m.quoteLines[m.currLine][m.typedLen + curr_word_lens_combined :])

		lines = ""
		for i := m.currLine+1; i <= m.currLine+m.linesToShow/2 && i < len(m.quoteLines); i++ {
			lines += m.quoteLines[i] + "\n"
			log.Println(i)
		}
		log.Println("Lines:", lines)
		contents += quoteStyle.Render(lines)

	} else if m.quoteCompleted {
		contents = typedStyle.Render("Completed test!!! :D")

		stats_str := fmt.Sprintf("\nWPM: %f\nCPM: %f\nACC: %f\n", m.stats.Wpm, m.stats.Cpm, m.stats.Acc)
		contents += typedStyle.Render(stats_str)
	}

	contents = contentStyle.Width(m.windowWidth-windowStyle.GetHorizontalFrameSize()).Render(contents)

	return contents
}
