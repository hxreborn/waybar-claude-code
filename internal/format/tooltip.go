package format

import (
	"fmt"
	"strings"

	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
)

func FormatNumber(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

func FormatDuration(minutes int) string {
	if minutes <= 0 {
		return "0m"
	}

	hours, mins := minutes/60, minutes%60

	switch {
	case hours > 0 && mins > 0:
		return fmt.Sprintf("%dh %dm", hours, mins)
	case hours > 0:
		return fmt.Sprintf("%dh", hours)
	default:
		return fmt.Sprintf("%dm", mins)
	}
}

func FormatTooltip(data *ccusage.BlocksData) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Requests: %d | Tokens: %s\n",
		data.Entries,
		FormatNumber(data.TotalTokens))

	fmt.Fprintf(&b, "Cost: $%.2f | Reset: %s",
		data.CostUSD,
		FormatDuration(data.RemainingMinutes))

	return b.String()
}

func FormatText(data *ccusage.BlocksData) string {
	return fmt.Sprintf("Reqs: %d | $%.2f", data.Entries, data.CostUSD)
}
