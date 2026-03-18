package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"hostkit/internal/config"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#7D56F4"))

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))

	categoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#04B575"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Faint(true)
)

// PackageItem represents a package in the list
type PackageItem struct {
	config   *config.Config
	category string
}

func (i PackageItem) FilterValue() string {
	return i.config.Name + " " + i.config.Description
}

// PackageDelegate handles rendering of package items
type PackageDelegate struct{}

func (d PackageDelegate) Height() int                             { return 2 }
func (d PackageDelegate) Spacing() int                            { return 1 }
func (d PackageDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d PackageDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PackageItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s %s\n%s",
		i.config.Icon,
		i.config.Name,
		descriptionStyle.Render(i.config.Description),
	)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

// Model represents the TUI state
type Model struct {
	list        list.Model
	configs     map[string]*config.Config
	selected    *config.Config
	width       int
	height      int
	showDetails bool
}

// NewModel creates a new TUI model
func NewModel(configs map[string]*config.Config) Model {
	items := make([]list.Item, 0, len(configs))

	for _, cfg := range configs {
		items = append(items, PackageItem{
			config:   cfg,
			category: cfg.Category,
		})
	}

	l := list.New(items, PackageDelegate{}, 0, 0)
	l.Title = "HostKit - Server Management Tool"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle

	return Model{
		list:    l,
		configs: configs,
	}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if i, ok := m.list.SelectedItem().(PackageItem); ok {
				m.selected = i.config
				m.showDetails = true
			}

		case "esc":
			if m.showDetails {
				m.showDetails = false
				m.selected = nil
			}

		case "i":
			if m.showDetails && m.selected != nil {
				// Trigger installation
				return m, tea.Printf("Installing %s...", m.selected.Name)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the TUI
func (m Model) View() string {
	if m.showDetails && m.selected != nil {
		return m.renderDetails()
	}
	return m.list.View()
}

func (m Model) renderDetails() string {
	cfg := m.selected

	var b strings.Builder

	b.WriteString(titleStyle.Render(fmt.Sprintf("%s %s", cfg.Icon, cfg.Name)))
	b.WriteString("\n\n")

	b.WriteString(categoryStyle.Render("Category: "))
	b.WriteString(cfg.Category)
	b.WriteString("\n\n")

	b.WriteString(descriptionStyle.Render(cfg.Description))
	b.WriteString("\n\n")

	if cfg.Version.Current != "" {
		b.WriteString(fmt.Sprintf("Version: %s\n", cfg.Version.Current))
	}

	if cfg.Requirements.MinRAM != "" {
		b.WriteString(fmt.Sprintf("Min RAM: %s\n", cfg.Requirements.MinRAM))
	}

	if cfg.Requirements.MinDisk != "" {
		b.WriteString(fmt.Sprintf("Min Disk: %s\n", cfg.Requirements.MinDisk))
	}

	if len(cfg.Requirements.Ports) > 0 {
		b.WriteString(fmt.Sprintf("Ports: %v\n", cfg.Requirements.Ports))
	}

	b.WriteString("\n")
	b.WriteString(itemStyle.Render("Press 'i' to install, 'esc' to go back, 'q' to quit"))

	return b.String()
}
