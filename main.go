package main

import (
	"os"
	"log"
	"syscall"
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

	oldStderr, _ := syscall.Dup(int(os.Stderr.Fd()))
	devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	syscall.Dup2(int(devNull.Fd()), int(os.Stderr.Fd()))
		
	portaudio.Initialize()
	defer func() {
		// Suppress ALSA warnings during cleanup by redirecting stderr
		portaudio.Terminate()
		
		// Restore stderr
		syscall.Dup2(oldStderr, int(os.Stderr.Fd()))
		syscall.Close(oldStderr)
		devNull.Close()
	}()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal("Unable to run tui:", err)
	}
}
