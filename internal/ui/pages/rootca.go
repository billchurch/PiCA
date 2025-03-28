package pages

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/ca/commands"
	"github.com/billchurch/PiCA/internal/config"
	"github.com/billchurch/PiCA/internal/crypto"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// RootCAModel represents the Root CA management page
type RootCAModel struct {
	width       int
	height      int
	styles      Styles
	inputs      []textinput.Model
	focusIndex  int
	message     string
	initialized bool
	config      *config.Config // Add configuration
}

// NewRootCAModel creates a new RootCAModel without configuration
func NewRootCAModel(styles Styles) RootCAModel {
	return NewRootCAModelWithConfig(styles, nil)
}

// NewRootCAModelWithConfig creates a new RootCAModel with configuration
func NewRootCAModelWithConfig(styles Styles, cfg *config.Config) RootCAModel {
	// Use default config if none provided
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	m := RootCAModel{
		styles:  styles,
		inputs:  make([]textinput.Model, 3),
		message: "",
		config:  cfg,
	}

	// Configure inputs
	var t textinput.Model

	t = textinput.New()
	t.Placeholder = "Path to Root CA config file"
	t.Focus()
	t.CharLimit = 100
	t.Width = 50
	// Pre-fill from configuration if available
	if cfg.CAConfigFile != "" {
		t.SetValue(cfg.CAConfigFile)
	} else {
		t.SetValue(fmt.Sprintf("%s/cfssl/root-ca-config.json", cfg.ConfigDir))
	}
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "Path to Root CA CSR file"
	t.CharLimit = 100
	t.Width = 50
	// Pre-fill with default CSR file if not specified
	if cfg.CAConfigFile != "" {
		t.SetValue(strings.Replace(cfg.CAConfigFile, "config", "csr", 1))
	} else {
		t.SetValue(fmt.Sprintf("%s/cfssl/root-ca-csr.json", cfg.ConfigDir))
	}
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "Path to save certificate"
	t.CharLimit = 100
	t.Width = 50
	// Pre-fill with default certificate path
	if cfg.CACertFile != "" {
		t.SetValue(cfg.CACertFile)
	} else {
		t.SetValue(fmt.Sprintf("%s/root-ca.pem", cfg.CertDir))
	}
	m.inputs[2] = t

	return m
}

// Init initializes the model
func (m RootCAModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the model based on messages
func (m RootCAModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.message = "Initializing Root CA..."

				// Create provider using configuration or environment
				var provider crypto.Provider
				var err error
				if m.config.ProviderType != "" {
					// Force specific provider type
					os.Setenv("PICA_PROVIDER", m.config.ProviderType)
				}
				
				provider, err = crypto.CreateDefaultProvider()
				if err != nil {
					m.message = fmt.Sprintf("Error creating provider: %s", err)
					return m, nil
				}
				defer provider.Close()

				// Determine key slot
				keySlot := crypto.SlotCA1
				if m.config.KeySlot != "" {
					// Convert hex string to slot
					slotVal, err := strconv.ParseInt(m.config.KeySlot, 16, 64)
					if err == nil {
						keySlot = crypto.Slot(slotVal)
					}
				}

				// Initialize the Root CA
				cmd := commands.NewInitCommand(
					ca.RootCA,
					m.inputs[0].Value(),
					m.inputs[1].Value(),
					m.inputs[2].Value(),
					keySlot,
				)
				
				// Set provider
				cmd.Provider = provider

				err = cmd.Execute()
				if err != nil {
					m.message = fmt.Sprintf("Error: %s", err)
				} else {
					m.message = "Root CA initialized successfully!"
					m.initialized = true
					
					// Update configuration with new values
					m.config.CAConfigFile = m.inputs[0].Value()
					m.config.CACertFile = m.inputs[2].Value()
					m.config.CAType = "root"
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
func (m RootCAModel) View() string {
	var b strings.Builder

	b.WriteString(m.styles.titleStyle.Render("Root CA Management"))
	b.WriteString("\n\n")

	if !m.initialized {
		b.WriteString("Initialize a new Root CA:\n\n")

		for i, input := range m.inputs {
			b.WriteString(input.View())
			if i < len(m.inputs)-1 {
				b.WriteString("\n")
			}
		}
	} else {
		b.WriteString(m.styles.successStyle.Render("✓ Root CA has been initialized\n\n"))
		b.WriteString("You can now use this Root CA to sign Sub CA certificates.\n")
		b.WriteString("Remember to keep this Root CA offline for security.\n")
		b.WriteString(fmt.Sprintf("\nRoot CA Certificate: %s\n", m.config.CACertFile))
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
func (m RootCAModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}
