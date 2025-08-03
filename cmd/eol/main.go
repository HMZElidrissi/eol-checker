package main

import (
	"log"
	"os"

	"github.com/HMZElidrissi/eol-checker/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
