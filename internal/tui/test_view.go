package tui

import (
	"fmt"
	"strings"

	"github.com/Mr-Meshky/vify-cli/internal/model"
)

// RenderLeaderboard formats the tested nodes into a sleek terminal table
func RenderLeaderboard(nodes []*model.ProxyNode) string {
	if len(nodes) == 0 {
		return DangerBadge.Render(" No healthy nodes discovered. Check your internet or clean IP. ") + "\n"
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(HeaderStyle.Render("🏆 Vify Latency Leaderboard (HTTP/204 Verified)"))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf(" %-4s %-4s %-8s %-32s %-22s %-12s\n",
		"#", "FLAG", "PROTO", "NAME", "ENDPOINT", "LATENCY"))
	sb.WriteString(strings.Repeat("─", 88) + "\n")

	for i, node := range nodes {
		rank := fmt.Sprintf("#%d", i+1)
		flag := node.CountryFlag
		if flag == "" {
			flag = "🌐"
		}
		proto := strings.ToUpper(string(node.Protocol))

		name := TruncateName(node.Name, 30)

		endpoint := fmt.Sprintf("%s:%d", node.Server, node.Port)
		if len(endpoint) > 20 {
			endpoint = endpoint[:17] + "..."
		}

		latStr := FormatLatency(node.Latency.Milliseconds())

		sb.WriteString(fmt.Sprintf(" %-4s %-4s %-8s %-32s %-22s %s\n",
			rank,
			flag,
			BadgeStyle.Render(proto),
			name,
			endpoint,
			latStr,
		))
	}

	sb.WriteString(strings.Repeat("─", 88) + "\n\n")
	return sb.String()
}
