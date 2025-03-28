package pages

import (
	"fmt"
	"strings"

	"github.com/billchurch/PiCA/internal/ca"
	"github.com/billchurch/PiCA/internal/ca/commands"
	"github.com/billchurch/PiCA/internal/crypto"
	"github.com/billchurch/PiCA/internal/yubikey"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Action represents the current action being performed
type Action int

const (
	ActionNone Action = iota
	ActionSign
	ActionRevoke
	ActionList
)

// CertItem represents a certificate in the list
type CertItem struct {
	subject      string
	serialNumber string
	notAfter     string
	status       string
}

// FilterValue implements list.Item interface
func (i CertItem) FilterValue() string { return i.subject }

// Title implements list.Item interface
func (i CertItem) Title() string { return i.subject }

// Description implements list.Item interface
func (i CertItem) Description() string {
	return fmt.Sprintf("SN: %s, Expires: %s, Status: %s",
		i.serialNumber, i.notAfter, i.status)
}

// CertManageModel represents the certificate management page
type CertManageModel struct {
	width        int
	height       int
	styles       Styles
	action       Action
	caType       ca.CAType
	inputs       []textinput.Model
	focusIndex   int
	message      string
	certList     list.Model
	certificates []CertItem
}

// NewCertManageModel creates a new CertManageModel
func NewCertManageModel(styles Styles, caType ca.CAType) CertManageModel {
	m := CertManageModel{
		styles:       styles,
		action:       ActionNone,
		caType:       caType,
		inputs:       make([]textinput.Model, 0),
		message:      "",
		certificates: make([]CertItem, 0),
	}

	// Initialize sample cert list (would be loaded from storage in real app)
	items := []list.Item{
		CertItem{
			subject:      "example.com",
			serialNumber: "1234567890",
			notAfter:     "2025-03-26",
			status:       "Valid",
		},
		CertItem{
			subject:      "admin.example.com",
			serialNumber: "1234567891",
			notAfter:     "2025-03-26",
			status:       "Valid",
		},
		CertItem{
			subject:      "revoked.example.com",
			serialNumber: "1234567892",
			notAfter:     "2025-03-26",
			status:       "Revoked",
		},
	}

	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	m.certList = list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	m.certList.Title = "Certificate List"

	return m
}

// setupSignInputs sets up inputs for signing a certificate
func (m *CertManageModel) setupSignInputs() {
	m.inputs = make([]textinput.Model, 4)
	var t textinput.Model

	t = textinput.New()
	t.Placeholder = "Path to CSR file"
	t.Focus()
	t.CharLimit = 100
	t.Width = 50
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "Path to save certificate"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "Profile (e.g., server, client)"
	t.CharLimit = 100
	t.Width = 50
	t.SetValue("server")
	m.inputs[2] = t

	t = textinput.New()
	t.Placeholder = "Path to CA config file"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[3] = t

	m.focusIndex = 0
}

// setupRevokeInputs sets up inputs for revoking a certificate
func (m *CertManageModel) setupRevokeInputs() {
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model

	t = textinput.New()
	t.Placeholder = "Serial number"
	t.Focus()
	t.CharLimit = 100
	t.Width = 50
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "Revocation reason"
	t.CharLimit = 100
	t.Width = 50
	t.SetValue("keyCompromise")
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "Path to CA config file"
	t.CharLimit = 100
	t.Width = 50
	m.inputs[2] = t

	m.focusIndex = 0
}

// Init initializes the model
func (m CertManageModel) Init() tea.Cmd {
	return nil
}

// Update updates the model based on messages
func (m CertManageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			if m.action == ActionNone {
				m.action = ActionSign
				m.setupSignInputs()
				return m, textinput.Blink
			}
		case "r":
			if m.action == ActionNone {
				m.action = ActionRevoke
				m.setupRevokeInputs()
				return m, textinput.Blink
			}
		case "l":
			if m.action == ActionNone {
				m.action = ActionList
				return m, nil
			}
		case "esc":
			m.action = ActionNone
			m.message = ""
			return m, nil

		case "tab", "shift+tab", "up", "down":
			if m.action == ActionSign || m.action == ActionRevoke {
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
			} else if m.action == ActionList {
				var cmd tea.Cmd
				m.certList, cmd = m.certList.Update(msg)
				cmds = append(cmds, cmd)
			}

		case "enter":
			if m.action == ActionSign && m.focusIndex == len(m.inputs)-1 {
				// Process sign form
				m.message = "Signing certificate..."

				// Create CA instance
				ca := ca.NewCA(
					m.caType,
					m.inputs[3].Value(), // CA config
					"",                  // Key file not needed when using YubiKey
					"",                  // Cert file not needed for signing
				)

				// Create sign command
				cmd := commands.NewSignCommand(
					ca,
					m.inputs[0].Value(),                     // CSR
					m.inputs[1].Value(),                     // Output cert
					m.inputs[2].Value(),                     // Profile
					crypto.FromYubiKeySlot(yubikey.SlotCA1), // YubiKey slot
				)

				err := cmd.Execute()
				if err != nil {
					m.message = fmt.Sprintf("Error: %s", err)
				} else {
					m.message = "Certificate signed successfully!"
				}
			} else if m.action == ActionRevoke && m.focusIndex == len(m.inputs)-1 {
				// Process revoke form
				m.message = "Revoking certificate..."

				// Create CA instance
				ca := ca.NewCA(
					m.caType,
					m.inputs[2].Value(), // CA config
					"",                  // Key file not needed when using YubiKey
					"",                  // Cert file not needed for revocation
				)

				// Create revoke command
				cmd := commands.NewRevokeCommand(
					ca,
					m.inputs[0].Value(),                     // Serial number
					m.inputs[1].Value(),                     // Reason
					crypto.FromYubiKeySlot(yubikey.SlotCA1), // YubiKey slot
				)

				err := cmd.Execute()
				if err != nil {
					m.message = fmt.Sprintf("Error: %s", err)
				} else {
					m.message = "Certificate revoked successfully!"
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if m.action == ActionList {
			m.certList.SetSize(msg.Width-4, msg.Height-10)
		}
	}

	// Handle character input for textinputs
	if m.action == ActionSign || m.action == ActionRevoke {
		cmd := m.updateInputs(msg)
		cmds = append(cmds, cmd)
	} else if m.action == ActionList {
		var cmd tea.Cmd
		m.certList, cmd = m.certList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m CertManageModel) View() string {
	var b strings.Builder

	title := "Certificate Management"
	if m.caType == ca.RootCA {
		title += " (Root CA)"
	} else {
		title += " (Sub CA)"
	}

	b.WriteString(m.styles.titleStyle.Render(title))
	b.WriteString("\n\n")

	switch m.action {
	case ActionNone:
		b.WriteString("Select an action:\n\n")
		b.WriteString("[s] Sign a certificate\n")
		b.WriteString("[r] Revoke a certificate\n")
		b.WriteString("[l] List certificates\n")

	case ActionSign:
		b.WriteString("Sign a new certificate:\n\n")
		for i, input := range m.inputs {
			b.WriteString(input.View())
			if i < len(m.inputs)-1 {
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n[esc] Cancel")

	case ActionRevoke:
		b.WriteString("Revoke a certificate:\n\n")
		for i, input := range m.inputs {
			b.WriteString(input.View())
			if i < len(m.inputs)-1 {
				b.WriteString("\n")
			}
		}
		b.WriteString("\n\n[esc] Cancel")

	case ActionList:
		b.WriteString(m.certList.View())
		b.WriteString("\n\n[esc] Back")
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
func (m CertManageModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
