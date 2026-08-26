package tui

import (
	"fmt"
	"strings"

	"github.com/Mr-Meshky/vify-cli/internal/model"
	tea "github.com/charmbracelet/bubbletea"
)

// ListModel represents the interactive server picker
type ListModel struct {
	Nodes        []*model.ProxyNode
	Cursor       int
	Selected     *model.ProxyNode
	Quitting     bool
	Filter       string
	FilteredList []*model.ProxyNode
}

// NewListModel creates a list selector model
func NewListModel(nodes []*model.ProxyNode) ListModel {
	m := ListModel{
		Nodes:        nodes,
		FilteredList: nodes,
	}
	return m
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.FilteredList)-1 {
				m.Cursor++
			}
		case "enter", " ":
			if len(m.FilteredList) > 0 && m.Cursor < len(m.FilteredList) {
				m.Selected = m.FilteredList[m.Cursor]
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	if m.Quitting {
		return SubtitleStyle.Render("Selection cancelled.\n")
	}

	var sb strings.Builder
	sb.WriteString(HeaderStyle.Render("⚡ Vify Server Selector (Use ↑/↓ to navigate, Enter to connect, 'q' to quit)"))
	sb.WriteString("\n\n")

	if len(m.FilteredList) == 0 {
		sb.WriteString("  No servers found.\n")
		return sb.String()
	}

	maxShow := 15
	start := 0
	if m.Cursor >= maxShow {
		start = m.Cursor - maxShow + 1
	}
	end := start + maxShow
	if end > len(m.FilteredList) {
		end = len(m.FilteredList)
	}

	for i := start; i < end; i++ {
		node := m.FilteredList[i]
		isCursor := i == m.Cursor

		cursorStr := "  "
		if isCursor {
			cursorStr = "❯ "
		}

		protoBadge := BadgeStyle.Render(strings.ToUpper(string(node.Protocol)))
		flag := node.CountryFlag
		if flag == "" {
			flag = "🌐"
		}

		name := TruncateName(node.Name, 30)

		latStr := "---"
		if node.Latency > 0 {
			latStr = fmt.Sprintf("%dms", node.Latency.Milliseconds())
		}

		line := fmt.Sprintf("%s %s %s  %-30s  %-10s  %s",
			cursorStr,
			flag,
			protoBadge,
			name,
			fmt.Sprintf("%s:%d", node.Server, node.Port),
			LatencyFast.Render(latStr),
		)

		if isCursor {
			line = RowSelectedStyle.Render(line)
		}

		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n" + SubtitleStyle.Render(fmt.Sprintf("Showing %d-%d of %d servers", start+1, end, len(m.FilteredList))))
	sb.WriteString("\n")

	return sb.String()
}
