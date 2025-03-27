package ui

import (
	"github.com/billchurch/pica/internal/ca"
	"github.com/billchurch/pica/internal/ui/pages"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Help  key.Binding
	Quit  key.Binding
	Tab   key.Binding
	Enter key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help, k.Quit, k.Tab, k.Enter},
	}
}

var keys = keyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
}

type page int

const (
	rootCAPage page = iota
	subCAPage
	certManagementPage
)

type Model struct {
	keys           keyMap
	help           help.Model
	styles         pages.Styles
	currentPage    page
	width          int
	height         int
	rootCAModel    pages.RootCAModel
	subCAModel     pages.SubCAModel
	certManageRoot pages.CertManageModel
	certManageSub  pages.CertManageModel
}

func NewModel() Model {
	styles := pages.DefaultStyles()
	
	return Model{
		keys:           keys,
		help:           help.New(),
		styles:         styles,
		currentPage:    rootCAPage,
		rootCAModel:    pages.NewRootCAModel(styles),
		subCAModel:     pages.NewSubCAModel(styles),
		certManageRoot: pages.NewCertManageModel(styles, ca.RootCA),
		certManageSub:  pages.NewCertManageModel(styles, ca.SubCA),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.rootCAModel.Init(),
		m.subCAModel.Init(),
		m.certManageRoot.Init(),
		m.certManageSub.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Tab):
			m.currentPage = (m.currentPage + 1) % 3
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
	}

	// Update active page model
	var cmd tea.Cmd
	switch m.currentPage {
	case rootCAPage:
		var newModel tea.Model
		newModel, cmd = m.rootCAModel.Update(msg)
		m.rootCAModel = newModel.(pages.RootCAModel)
	case subCAPage:
		var newModel tea.Model
		newModel, cmd = m.subCAModel.Update(msg)
		m.subCAModel = newModel.(pages.SubCAModel)
	case certManagementPage:
		// Determine if we're managing Root CA or Sub CA certs
		// For now, let's just use Sub CA cert management
		var newModel tea.Model
		newModel, cmd = m.certManageSub.Update(msg)
		m.certManageSub = newModel.(pages.CertManageModel)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var content string

	switch m.currentPage {
	case rootCAPage:
		content = m.rootCAModel.View()
	case subCAPage:
		content = m.subCAModel.View()
	case certManagementPage:
		content = m.certManageSub.View()
	}

	help := m.help.View(m.keys)
	
	// Add navigation bar
	nav := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA")).
		Render("Press TAB to switch views • ? for help • q to quit")
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			Width(m.width).
			Render("PiCA - Raspberry Pi Certificate Authority"),
		content,
		nav,
		help,
	)
}
