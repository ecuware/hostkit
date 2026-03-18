package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"hostkit/internal/config"
	"hostkit/internal/status"
)

// PackageItemWithStatus represents a package with installation status
type PackageItemWithStatus struct {
	config *config.Config
	status status.Status
}

func (i PackageItemWithStatus) FilterValue() string {
	return i.config.Name + " " + i.config.Description
}

// PackageDelegateWithStatus handles rendering with status icons
type PackageDelegateWithStatus struct{}

func (d PackageDelegateWithStatus) Height() int  { return 2 }
func (d PackageDelegateWithStatus) Spacing() int { return 1 }
func (d PackageDelegateWithStatus) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d PackageDelegateWithStatus) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(PackageItemWithStatus)
	if !ok {
		return
	}

	statusIcon := status.GetStatusIcon(i.status)

	str := fmt.Sprintf("%s %s %s\n%s",
		statusIcon,
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
