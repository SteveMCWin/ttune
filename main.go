package main

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/gordonklaus/portaudio"
)

func logFilePath() string {
    cacheDir, err := os.UserCacheDir()
    if err != nil {
        return "debug.log" // fallback
    }
    appCache := filepath.Join(cacheDir, "ttune")
    os.MkdirAll(appCache, 0755)
    return filepath.Join(appCache, "debug.log")
}

func main() {

	f, err := tea.LogToFile(logFilePath(), "debug")
	if err != nil {
		log.Fatalf("failed setting the debug log file: %v", err)
	}
	defer f.Close()
	// should check for commandline args
	m := NewModel()

	log.Println()
	log.Println("~~~~~~~~~PROGRAM START~~~~~~~~~")
	log.Println()

	// Suppress ALSA warnings during cleanup by redirecting stderr
	oldStderr, _ := syscall.Dup(int(os.Stderr.Fd()))
	devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	syscall.Dup2(int(devNull.Fd()), int(os.Stderr.Fd()))
		
	portaudio.Initialize()
	defer func() {
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
