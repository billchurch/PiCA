package pages

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles holds the styling definitions for the UI
type Styles struct {
	titleStyle   lipgloss.Style
	pageStyle    lipgloss.Style
	errorStyle   lipgloss.Style
	successStyle lipgloss.Style
	messageStyle lipgloss.Style
	selectStyle  lipgloss.Style
	activeStyle  lipgloss.Style
	infoStyle    lipgloss.Style
}

// DefaultStyles returns the default UI styles
func DefaultStyles() Styles {
	return Styles{
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(1),
		pageStyle: lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingRight(2),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		successStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		messageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		selectStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(1).
			PaddingRight(1),
		activeStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#0088DD")).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),
		infoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Italic(true),
	}
}
