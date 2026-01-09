package main

import (
	"log"
	"github.com/gordonklaus/portaudio"
	tea "charm.land/bubbletea/v2"
)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("failed setting the debug log file: %v", err)
	}
	defer f.Close()
	// should check for commandline args
	m := NewModel()

	log.Println()
	log.Println("~~~~~~~~~PROGRAM START~~~~~~~~~")
	log.Println()

	portaudio.Initialize()
	defer portaudio.Terminate()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal("Unable to run tui:", err)
		m.AudioStream.Stop()
		m.AudioStream.Close()
	}
}
