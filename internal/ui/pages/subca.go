package pages

import (
	"fmt"
	"strings"

	"github.com/billchurch/pica/internal/ca"
	"github.com/billchurch/pica/internal/ca/commands"
	"github.com/billchurch/pica/internal/crypto"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// SubCAModel represents the Sub CA management page
type SubCAModel struct {
	width       int
	height      int
	styles      Styles
	inputs      []textinput.Model
	focusIndex  int
	message     string
	initialized bool
}

// NewSubCAModel creates a new SubCAModel
func NewSubCAModel(styles Styles) SubCAModel {
	m := SubCAModel{
		styles:  styles,
		inputs:  make([]textinput.Model, 6),
		message: "",
	}

	// Configure inputs
	var t textinput.Model

	t = textinput.New()
	t.Placeholder = "Path to Sub CA config file"
	t.Focus()
	t.CharLimit = 100
	t.Width = 50
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "Path to Sub CA CSR file"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "Path to save certificate"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[2] = t

	t = textinput.New()
	t.Placeholder = "Path to Root CA certificate"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[3] = t

	t = textinput.New()
	t.Placeholder = "Path to Root CA config file"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[4] = t

	t = textinput.New()
	t.Placeholder = "Profile (e.g., subca)"
	t.CharLimit = 100
	t.Width = 50
	t.SetValue("subca")
	m.inputs[5] = t

	return m
}

// Init initializes the model
func (m SubCAModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the model based on messages
func (m SubCAModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			// Cycle through inputs
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) - 1
			} else if m.focusIndex >= len(m.inputs) {
				m.focusIndex = 0
			}

			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds = append(cmds, m.inputs[i].Focus())
				} else {
					m.inputs[i].Blur()
				}
			}

		case "enter":
			if m.focusIndex == len(m.inputs)-1 {
				// Process the form
				m.message = "Initializing Sub CA..."

				// Create provider
				provider, err := crypto.CreateDefaultProvider()
				if err != nil {
					m.message = fmt.Sprintf("Error creating provider: %s", err)
					return m, nil
				}

				// Create Init command for Sub CA
				cmd := commands.NewInitCommand(
					ca.SubCA,
					m.inputs[0].Value(), // Config file
					m.inputs[1].Value(), // CSR file
					m.inputs[2].Value(), // Output cert
					crypto.SlotCA2,      // Slot for Sub CA
				)

				// Set the provider
				cmd.Provider = provider

				// Set the additional fields with values from the form
				cmd.RootCACertFile = m.inputs[3].Value()
				cmd.RootCAConfigFile = m.inputs[4].Value()
				cmd.Profile = m.inputs[5].Value()

				// Execute the command
				err = cmd.Execute()
				if err != nil {
					m.message = fmt.Sprintf("Error: %s", err)
				} else {
					m.message = "Sub CA initialized successfully!"
					m.initialized = true
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Handle character input
	cmd := m.updateInputs(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m SubCAModel) View() string {
	var b strings.Builder

	b.WriteString(m.styles.titleStyle.Render("Sub CA Management"))
	b.WriteString("\n\n")

	if !m.initialized {
		b.WriteString("Initialize a new Sub CA:\n\n")

		for i, input := range m.inputs {
			b.WriteString(input.View())
			if i < len(m.inputs)-1 {
				b.WriteString("\n")
			}
		}
	} else {
		b.WriteString(m.styles.successStyle.Render("✓ Sub CA has been initialized\n\n"))
		b.WriteString("You can now use this Sub CA to issue end-entity certificates.\n")
	}

	if m.message != "" {
		b.WriteString("\n\n")
		if strings.HasPrefix(m.message, "Error") {
			b.WriteString(m.styles.errorStyle.Render(m.message))
		} else {
			b.WriteString(m.styles.messageStyle.Render(m.message))
		}
	}

	return b.String()
}

// updateInputs handles updates to the text inputs
func (m SubCAModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}
