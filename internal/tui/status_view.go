package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mr-Meshky/vify-cli/internal/core"
	"github.com/Mr-Meshky/vify-cli/internal/model"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// StatusModel renders a live dashboard of the running VPN connection
type StatusModel struct {
	Session   *model.ActiveSession
	Stats     *core.TrafficStats
	StartTime time.Time
	Quitting  bool
}

type tickMsg time.Time

func tickEverySecond() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NewStatusModel creates a live status model
func NewStatusModel(session *model.ActiveSession, stats *core.TrafficStats) StatusModel {
	return StatusModel{
		Session:   session,
		Stats:     stats,
		StartTime: time.Now(),
	}
}

func (m StatusModel) Init() tea.Cmd {
	return tickEverySecond()
}

func (m StatusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		return m, tickEverySecond()
	}
	return m, nil
}

func (m StatusModel) View() string {
	if m.Quitting {
		return SubtitleStyle.Render("Disconnecting Vify VPN...\n")
	}

	node := m.Session.Node
	flag := node.CountryFlag
	if flag == "" {
		flag = "🌐"
	}

	upTotal, downTotal, upSpeed, downSpeed := int64(0), int64(0), int64(0), int64(0)
	if m.Stats != nil {
		upTotal, downTotal, upSpeed, downSpeed = m.Stats.Snapshot()
	}

	uptime := time.Since(m.StartTime).Truncate(time.Second)

	// Build Dashboard Card
	var content strings.Builder

	header := fmt.Sprintf("%s  %s  %s",
		TitleStyle.Render("⚡ VIFY CONNECTED"),
		SuccessBadge.Render("ACTIVE"),
		BadgeStyle.Render(strings.ToUpper(string(m.Session.Mode))),
	)
	content.WriteString(header + "\n\n")

	serverLine := fmt.Sprintf("🌍 Server:   %s %s (%s:%d)",
		flag,
		lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(node.Name),
		node.Server,
		node.Port,
	)
	protoLine := fmt.Sprintf("🔒 Protocol: %s  |  Direct Iran Bypass: %s",
		BadgeStyle.Render(strings.ToUpper(string(node.Protocol))),
		SuccessBadge.Render("ENABLED (0.5x quota)"),
	)
	uptimeLine := fmt.Sprintf("⏱  Uptime:   %s", uptime.String())

	content.WriteString(serverLine + "\n")
	content.WriteString(protoLine + "\n")
	content.WriteString(uptimeLine + "\n\n")

	// Traffic & Speed Gauges
	speedGauge := fmt.Sprintf("↓ Download Speed: %-15s   ↑ Upload Speed: %-15s",
		lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render(core.FormatSpeed(downSpeed)),
		lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(core.FormatSpeed(upSpeed)),
	)
	trafficGauge := fmt.Sprintf("↓ Total Ingest:   %-15s   ↑ Total Egress: %-15s",
		lipgloss.NewStyle().Foreground(ColorSuccess).Render(core.FormatBytes(downTotal)),
		lipgloss.NewStyle().Foreground(ColorPrimary).Render(core.FormatBytes(upTotal)),
	)
	proxyPortLine := fmt.Sprintf("🔌 Local SOCKS5: 127.0.0.1:%-6d        Local HTTP: 127.0.0.1:%-6d",
		m.Session.LocalSocks,
		m.Session.LocalHTTP,
	)

	content.WriteString(speedGauge + "\n")
	content.WriteString(trafficGauge + "\n")
	content.WriteString(proxyPortLine + "\n\n")

	content.WriteString(SubtitleStyle.Render("Press 'q' or Ctrl+C to stop the VPN and cleanly restore system settings."))

	return CardStyle.Render(content.String()) + "\n"
}
