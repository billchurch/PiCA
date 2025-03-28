package main

import (
	"fmt"
	"log"
	"os"

	"github.com/billchurch/PiCA/internal/config"
	"github.com/billchurch/PiCA/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.Load(os.Args[1:], "")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create UI model with configuration
	model := ui.NewModelWithConfig(cfg)

	// Run the application
	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running PiCA: %v\n", err)
		os.Exit(1)
	}
}
