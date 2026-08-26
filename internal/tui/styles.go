package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// TruncateName shortens display names by rune count, not byte index, so it
// never splits a multi-byte flag emoji or non-Latin character mid-sequence.
func TruncateName(name string, max int) string {
	runes := []rune(name)
	if len(runes) <= max {
		return name
	}
	return string(runes[:max-3]) + "..."
}

var (
	// Vibrant Colors
	ColorPrimary   = lipgloss.Color("#00F0FF") // Neon Cyan
	ColorSecondary = lipgloss.Color("#7928CA") // Neon Purple
	ColorSuccess   = lipgloss.Color("#00DF89") // Emerald Green
	ColorWarning   = lipgloss.Color("#FFBE0B") // Amber Yellow
	ColorDanger    = lipgloss.Color("#FF0055") // Neon Red
	ColorMuted     = lipgloss.Color("#6C757D") // Gray
	ColorBgDark    = lipgloss.Color("#0D1117") // Deep Background
	ColorCardBg    = lipgloss.Color("#161B22") // Card Surface

	// Typography & Badges
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	BadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorSecondary).
			Padding(0, 1)

	SuccessBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorSuccess).
			Padding(0, 1)

	WarningBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorWarning).
			Padding(0, 1)

	DangerBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorDanger).
			Padding(0, 1)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			Background(ColorCardBg)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorPrimary).
			PaddingBottom(1).
			MarginBottom(1)

	RowSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Background(lipgloss.Color("#21262D"))

	LatencyFast = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	LatencyMid  = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	LatencySlow = lipgloss.NewStyle().Foreground(ColorDanger).Bold(true)
)

// FormatLatency returns colored latency string
func FormatLatency(ms int64) string {
	if ms <= 0 {
		return DangerBadge.Render("TIMEOUT")
	}
	text := fmt.Sprintf("%d ms", ms)
	if ms < 500 {
		return LatencyFast.Render(text)
	} else if ms < 1000 {
		return LatencyMid.Render(text)
	}
	return LatencySlow.Render(text)
}
