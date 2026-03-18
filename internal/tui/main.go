package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"hostkit/internal/config"
	"hostkit/internal/history"
	"hostkit/internal/status"
)

// MenuState represents the current menu state
type MenuState int

const (
	StateMainMenu MenuState = iota
	StateCategory
	StatePackageList
	StatePackageDetails
	StateDialog
	StateInstalling
	StateHistory
	StateSystemMonitor
	StateServiceManager
	StateLogViewer
	StateBackupManager
	StateFirewall
)

// MainModel is the main TUI model with category-based navigation
type MainModel struct {
	state           MenuState
	categories      []string
	configs         map[string]*config.Config
	categoryList    list.Model
	packageList     *list.Model
	selectedPkg     *config.Config
	currentCategory string
	installModel    InstallModel
	monitorModel    SystemMonitorModel
	dialog          *DialogModel
	statusChecker   *status.Checker
	historyManager  *history.Manager
	statuses        map[string]status.Status
	width           int
	height          int
}

// CategoryItem represents a category in the main menu
type CategoryItem struct {
	name        string
	description string
	icon        string
	count       int
}

func (i CategoryItem) FilterValue() string { return i.name }
func (i CategoryItem) Title() string       { return fmt.Sprintf("%s %s", i.icon, i.name) }
func (i CategoryItem) Description() string {
	return fmt.Sprintf("%s (%d packages)", i.description, i.count)
}

// HistoryItem represents history list item
type HistoryItem struct {
	PackageID   string
	PackageName string
	Version     string
	Status      string
	StartTime   time.Time
	EndTime     time.Time
	Duration    string
	Error       string
}

func (i HistoryItem) FilterValue() string { return i.PackageName }
func (i HistoryItem) Title() string {
	return fmt.Sprintf("%s %s", getStatusIconFromHistory(i.Status), i.PackageName)
}
func (i HistoryItem) Description() string {
	statusText := i.Status
	if i.Status == "success" {
		statusText = "✓ Success"
	} else if i.Status == "failed" {
		statusText = "✗ Failed"
	}
	return fmt.Sprintf("%s - %s - %s", statusText, i.Version, i.StartTime.Format("2006-01-02 15:04"))
}

func getStatusIconFromHistory(status string) string {
	switch status {
	case "success":
		return "✅"
	case "failed":
		return "❌"
	case "installing":
		return "⏳"
	default:
		return "  "
	}
}

// NewMainModel creates a new main model with category menu
func NewMainModel(configs map[string]*config.Config) MainModel {
	// Group configs by category
	categoryMap := make(map[string][]*config.Config)
	for _, cfg := range configs {
		categoryMap[cfg.Category] = append(categoryMap[cfg.Category], cfg)
	}

	// Create category items
	categories := []CategoryItem{
		{name: "Panels", description: "Hosting Control Panels", icon: "🎛️", count: len(categoryMap["panel"])},
		{name: "Databases", description: "Database Servers", icon: "🗄️", count: len(categoryMap["database"])},
		{name: "Web Servers", description: "Web & Application Servers", icon: "🌐", count: len(categoryMap["webserver"])},
		{name: "Security", description: "Security & Firewall Tools", icon: "🔒", count: len(categoryMap["security"])},
		{name: "Services", description: "Container & Service Platforms", icon: "⚙️", count: len(categoryMap["services"])},
		{name: "Monitoring", description: "Monitoring & Analytics", icon: "📊", count: len(categoryMap["monitoring"])},
		{name: "VPN", description: "VPN & Network Tools", icon: "🔐", count: len(categoryMap["vpn"])},
		{name: "System Tools", description: "Services, Monitor, Logs, Backup & Firewall", icon: "🛠️", count: 5},
		{name: "History", description: "Installation History", icon: "📜", count: 0},
	}

	// Filter out empty categories
	var items []list.Item
	var catNames []string
	for _, cat := range categories {
		if cat.count > 0 || cat.name == "History" {
			items = append(items, cat)
			catNames = append(catNames, strings.ToLower(cat.name))
		}
	}

	// Create category list
	catList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	catList.Title = "HostKit - Select Category"
	catList.SetShowStatusBar(false)
	catList.SetFilteringEnabled(false)
	catList.Styles.Title = titleStyle

	// Initialize status checker and history manager
	statusChecker := status.NewChecker()
	statuses := make(map[string]status.Status)

	// Check status for all packages
	for id, cfg := range configs {
		statuses[id] = statusChecker.Check(cfg)
	}

	return MainModel{
		state:          StateMainMenu,
		categories:     catNames,
		configs:        configs,
		categoryList:   catList,
		statusChecker:  statusChecker,
		historyManager: history.NewManager(""),
		statuses:       statuses,
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.categoryList.SetSize(msg.Width, msg.Height-4)
		if m.packageList != nil {
			m.packageList.SetSize(msg.Width, msg.Height-4)
		}

	case tea.KeyMsg:
		// Global keys
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		// State-specific keys
		switch m.state {
		case StateMainMenu:
			return m.handleMainMenuKeys(msg)
		case StatePackageList:
			return m.handlePackageListKeys(msg)
		case StatePackageDetails:
			return m.handlePackageDetailsKeys(msg)
		case StateDialog:
			return m.handleDialogKeys(msg)
		case StateInstalling:
			return m.handleInstallingKeys(msg)
		case StateHistory:
			return m.handleHistoryKeys(msg)
		case StateSystemMonitor:
			monitorModel, cmd := m.monitorModel.Update(msg)
			m.monitorModel = monitorModel
			if msg.String() == "q" || msg.String() == "esc" {
				m.state = StatePackageList
				return m.showSystemTools()
			}
			return m, cmd
		}
	}

	// Update lists
	var cmd tea.Cmd
	switch m.state {
	case StateMainMenu:
		m.categoryList, cmd = m.categoryList.Update(msg)
	case StatePackageList:
		if m.packageList != nil {
			newList, newCmd := m.packageList.Update(msg)
			m.packageList = &newList
			cmd = newCmd
		}
	}

	return m, cmd
}

func (m MainModel) handleMainMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		if item, ok := m.categoryList.SelectedItem().(CategoryItem); ok {
			if item.name == "History" {
				return m.showHistory()
			}
			if item.name == "System Tools" {
				return m.showSystemTools()
			}
			return m.showCategoryPackages(item.name)
		}
	}

	var cmd tea.Cmd
	m.categoryList, cmd = m.categoryList.Update(msg)
	return m, cmd
}

func (m MainModel) showCategoryPackages(category string) (tea.Model, tea.Cmd) {
	m.currentCategory = category
	m.state = StatePackageList

	// Create package list for this category
	var items []list.Item
	for _, cfg := range m.configs {
		if strings.EqualFold(cfg.Category, category) ||
			(category == "Panels" && cfg.Category == "panel") ||
			(category == "Databases" && cfg.Category == "database") ||
			(category == "Web Servers" && cfg.Category == "webserver") ||
			(category == "Security" && cfg.Category == "security") ||
			(category == "Services" && cfg.Category == "services") ||
			(category == "Monitoring" && cfg.Category == "monitoring") ||
			(category == "VPN" && cfg.Category == "vpn") {
			items = append(items, PackageItemWithStatus{
				config: cfg,
				status: m.statuses[cfg.ID],
			})
		}
	}

	newList := list.New(items, PackageDelegateWithStatus{}, m.width, m.height-4)
	newList.Title = fmt.Sprintf("%s - Select Package", category)
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(true)
	newList.Styles.Title = titleStyle
	m.packageList = &newList

	return m, nil
}

// SystemToolItem represents system tools menu item
type SystemToolItem struct {
	name        string
	description string
	icon        string
}

func (i SystemToolItem) FilterValue() string { return i.name }
func (i SystemToolItem) Title() string       { return fmt.Sprintf("%s %s", i.icon, i.name) }
func (i SystemToolItem) Description() string { return i.description }

func (m MainModel) showSystemTools() (tea.Model, tea.Cmd) {
	m.state = StatePackageList
	m.currentCategory = "System Tools"

	tools := []SystemToolItem{
		{name: "Service Manager", description: "Start/Stop/Restart services", icon: "⚙️"},
		{name: "System Monitor", description: "CPU, RAM, Disk, Network usage", icon: "📊"},
		{name: "Log Viewer", description: "View system and service logs", icon: "📜"},
		{name: "Backup Manager", description: "Create and restore backups", icon: "💾"},
		{name: "Firewall", description: "Manage firewall rules", icon: "🔥"},
	}

	var items []list.Item
	for _, tool := range tools {
		items = append(items, tool)
	}

	newList := list.New(items, list.NewDefaultDelegate(), m.width, m.height-4)
	newList.Title = "System Tools"
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(true)
	newList.Styles.Title = titleStyle
	m.packageList = &newList

	return m, nil
}

func (m MainModel) showHistory() (tea.Model, tea.Cmd) {
	m.state = StateHistory

	historyData, err := m.historyManager.GetHistory(50)
	if err != nil {
		historyData = []history.Entry{}
	}

	var items []list.Item
	for _, entry := range historyData {
		items = append(items, HistoryItem{
			PackageID:   entry.PackageID,
			PackageName: entry.PackageName,
			Version:     entry.Version,
			Status:      entry.Status,
			StartTime:   entry.StartTime,
			EndTime:     entry.EndTime,
			Duration:    entry.Duration,
			Error:       entry.Error,
		})
	}

	newList := list.New(items, list.NewDefaultDelegate(), m.width, m.height-4)
	newList.Title = "Installation History"
	newList.SetShowStatusBar(false)
	newList.SetFilteringEnabled(true)
	newList.Styles.Title = titleStyle
	m.packageList = &newList

	return m, nil
}

func (m MainModel) handlePackageListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.currentCategory == "System Tools" {
			if item, ok := m.packageList.SelectedItem().(SystemToolItem); ok {
				switch item.name {
				case "System Monitor":
					m.monitorModel = NewSystemMonitorModel()
					m.state = StateSystemMonitor
					return m, m.monitorModel.Init()
				default:
					return m, nil
				}
			}
		} else {
			if item, ok := m.packageList.SelectedItem().(PackageItemWithStatus); ok {
				m.selectedPkg = item.config
				m.state = StatePackageDetails
			}
		}
		return m, nil

	case "esc", "backspace":
		m.state = StateMainMenu
		m.packageList = nil
		return m, nil
	}

	if m.packageList != nil {
		newList, newCmd := m.packageList.Update(msg)
		m.packageList = &newList
		return m, newCmd
	}

	return m, nil
}

func (m MainModel) handlePackageDetailsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "i":
		if m.selectedPkg != nil {
			// Check if already installed
			if m.statuses[m.selectedPkg.ID] == status.StatusInstalled {
				// Show confirmation dialog for reinstall
				deps := ResolveDependencies(m.selectedPkg, m.configs)
				dialog := NewConfirmDialog(m.selectedPkg, deps)
				m.dialog = &dialog
				m.state = StateDialog
				return m, nil
			}

			// Show confirmation dialog
			deps := ResolveDependencies(m.selectedPkg, m.configs)
			dialog := NewConfirmDialog(m.selectedPkg, deps)
			m.dialog = &dialog
			m.state = StateDialog
			return m, nil
		}

	case "esc", "backspace":
		m.state = StatePackageList
		m.selectedPkg = nil
		return m, nil

	case "q":
		return m, tea.Quit
	}

	return m, nil
}

func (m MainModel) handleDialogKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.dialog == nil {
		return m, nil
	}

	switch msg.String() {
	case "left", "right":
		if m.dialog.selected == 0 {
			m.dialog.selected = 1
		} else {
			m.dialog.selected = 0
		}

	case "enter":
		if m.dialog.IsConfirmed() {
			m.dialog = nil
			m.installModel = NewInstallModel(m.selectedPkg)
			m.state = StateInstalling

			// Record installation start
			entry := history.Entry{
				PackageID:   m.selectedPkg.ID,
				PackageName: m.selectedPkg.Name,
				Version:     m.selectedPkg.Version.Current,
				Status:      "installing",
				StartTime:   time.Now(),
			}
			m.historyManager.AddEntry(entry)

			return m, m.installModel.Init()
		} else {
			m.dialog = nil
			m.state = StatePackageDetails
		}

	case "esc":
		m.dialog = nil
		m.state = StatePackageDetails
	}

	return m, nil
}

func (m MainModel) handleInstallingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	installModel, cmd := m.installModel.Update(msg)
	m.installModel = installModel

	// Check if installation completed
	if m.installModel.state == InstallStateSuccess || m.installModel.state == InstallStateError {
		// Record installation result
		status := "success"
		if m.installModel.state == InstallStateError {
			status = "failed"
		}

		entry := history.Entry{
			PackageID:   m.selectedPkg.ID,
			PackageName: m.selectedPkg.Name,
			Version:     m.selectedPkg.Version.Current,
			Status:      status,
			EndTime:     time.Now(),
			Duration:    time.Since(m.installModel.startTime).String(),
		}
		if m.installModel.state == InstallStateError && m.installModel.error != nil {
			entry.Error = m.installModel.error.Error()
		}
		m.historyManager.AddEntry(entry)

		// Update status
		m.statuses[m.selectedPkg.ID] = m.statusChecker.Check(m.selectedPkg)
	}

	switch msg.String() {
	case "esc":
		if m.installModel.state == InstallStateSuccess || m.installModel.state == InstallStateError {
			m.state = StatePackageDetails
		}

	case "q":
		if m.installModel.state == InstallStateSuccess || m.installModel.state == InstallStateError {
			m.state = StatePackageList
			m.selectedPkg = nil
		}
	}

	return m, cmd
}

func (m MainModel) handleHistoryKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if item, ok := m.packageList.SelectedItem().(HistoryItem); ok {
			// Find the package config
			if cfg, ok := m.configs[item.PackageID]; ok {
				m.selectedPkg = cfg
				m.state = StatePackageDetails
			}
		}
		return m, nil

	case "esc", "backspace":
		m.state = StateMainMenu
		m.packageList = nil
		return m, nil
	}

	if m.packageList != nil {
		newList, newCmd := m.packageList.Update(msg)
		m.packageList = &newList
		return m, newCmd
	}

	return m, nil
}

// View renders the TUI
func (m MainModel) View() string {
	switch m.state {
	case StateMainMenu:
		return m.categoryList.View()
	case StatePackageList, StateHistory:
		if m.packageList != nil {
			return m.packageList.View()
		}
		return "Loading..."
	case StatePackageDetails:
		return m.renderPackageDetails()
	case StateDialog:
		if m.dialog != nil {
			return m.dialog.View()
		}
		return "Loading..."
	case StateInstalling:
		return m.installModel.View()
	case StateSystemMonitor:
		return m.monitorModel.View()
	default:
		return "Unknown state"
	}
}

func (m MainModel) renderPackageDetails() string {
	if m.selectedPkg == nil {
		return "No package selected"
	}

	cfg := m.selectedPkg
	var b strings.Builder

	// Title with status
	statusIcon := status.GetStatusIcon(m.statuses[cfg.ID])
	b.WriteString(titleStyle.Render(fmt.Sprintf("%s %s %s", statusIcon, cfg.Icon, cfg.Name)))
	b.WriteString("\n\n")

	// Status text
	b.WriteString(categoryStyle.Render("Status: "))
	b.WriteString(status.GetStatusText(m.statuses[cfg.ID]))
	b.WriteString("\n\n")

	// Category
	b.WriteString(categoryStyle.Render("Category: "))
	b.WriteString(cfg.Category)
	b.WriteString("\n\n")

	// Description
	b.WriteString(descriptionStyle.Render(cfg.Description))
	b.WriteString("\n\n")

	// Version
	if cfg.Version.Current != "" {
		b.WriteString(itemStyle.Render(fmt.Sprintf("Version: %s\n", cfg.Version.Current)))
	}

	// Requirements
	if cfg.Requirements.MinRAM != "" {
		b.WriteString(itemStyle.Render(fmt.Sprintf("Min RAM: %s\n", cfg.Requirements.MinRAM)))
	}

	if cfg.Requirements.MinDisk != "" {
		b.WriteString(itemStyle.Render(fmt.Sprintf("Min Disk: %s\n", cfg.Requirements.MinDisk)))
	}

	if cfg.Requirements.InstallSize != "" {
		b.WriteString(itemStyle.Render(fmt.Sprintf("Install Size: %s\n", cfg.Requirements.InstallSize)))
	}

	if cfg.Requirements.EstimatedTime != "" {
		b.WriteString(itemStyle.Render(fmt.Sprintf("Estimated Time: %s\n", cfg.Requirements.EstimatedTime)))
	}

	if len(cfg.Requirements.Ports) > 0 {
		ports := make([]string, len(cfg.Requirements.Ports))
		for i, port := range cfg.Requirements.Ports {
			ports[i] = fmt.Sprintf("%d", port)
		}
		b.WriteString(itemStyle.Render(fmt.Sprintf("Ports: %s\n", strings.Join(ports, ", "))))
	}

	// Dependencies
	if len(cfg.Dependencies.Required) > 0 {
		b.WriteString("\n")
		b.WriteString(categoryStyle.Render("Required Dependencies:"))
		b.WriteString("\n")
		for _, dep := range cfg.Dependencies.Required {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  • %s\n", dep)))
		}
	}

	if len(cfg.Dependencies.Optional) > 0 {
		b.WriteString("\n")
		b.WriteString(categoryStyle.Render("Optional Dependencies:"))
		b.WriteString("\n")
		for _, dep := range cfg.Dependencies.Optional {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  • %s\n", dep)))
		}
	}

	// Support links
	if cfg.Support.Docs != "" || cfg.Support.Issues != "" {
		b.WriteString("\n")
		b.WriteString(categoryStyle.Render("Support:"))
		b.WriteString("\n")
		if cfg.Support.Docs != "" {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  Docs: %s\n", cfg.Support.Docs)))
		}
		if cfg.Support.Issues != "" {
			b.WriteString(itemStyle.Render(fmt.Sprintf("  Issues: %s\n", cfg.Support.Issues)))
		}
	}

	// Footer
	b.WriteString("\n")
	action := "Install"
	if m.statuses[cfg.ID] == status.StatusInstalled {
		action = "Reinstall"
	}
	b.WriteString(helpStyle.Render(fmt.Sprintf("i: %s  •  esc: Back  •  q: Quit", action)))

	return b.String()
}
