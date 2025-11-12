package format

import (
	"fmt"
	"strings"

	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
)

// FormatNumber formats large numbers with K/M suffix
// Examples: 1500 -> "1.5K", 2500000 -> "2.5M", 500 -> "500"
func FormatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

// FormatDuration formats minutes into human-readable duration
// Examples: 125 -> "2h 5m", 45 -> "45m", 0 -> "0m"
func FormatDuration(minutes int) string {
	if minutes <= 0 {
		return "0m"
	}

	hours := minutes / 60
	mins := minutes % 60

	if hours > 0 && mins > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", mins)
}

// FormatTooltip generates a multi-line tooltip from BlocksData
// Format:
//   Requests: 42 | Tokens: 1.2M
//   Cost: $0.45 | Reset: 2h 15m
func FormatTooltip(data *ccusage.BlocksData) string {
	var lines []string

	// Line 1: Requests and Tokens
	line1 := fmt.Sprintf("Requests: %d | Tokens: %s",
		data.Entries,
		FormatNumber(data.TotalTokens))
	lines = append(lines, line1)

	// Line 2: Cost and Reset time
	line2 := fmt.Sprintf("Cost: $%.2f | Reset: %s",
		data.CostUSD,
		FormatDuration(data.RemainingMinutes))
	lines = append(lines, line2)

	return strings.Join(lines, "\n")
}

// FormatText generates compact text for Waybar display
// Format: "Reqs: X | $Y"
func FormatText(data *ccusage.BlocksData) string {
	return fmt.Sprintf("Reqs: %d | $%.2f", data.Entries, data.CostUSD)
}
