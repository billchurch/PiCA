package main

import (
	"fmt"
	"os"

	"github.com/billchurch/PiCA/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(ui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running PiCA: %v\n", err)
		os.Exit(1)
	}
}
