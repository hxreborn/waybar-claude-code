package format

import (
	"fmt"
	"strings"
	"time"

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

func formatResetTime(t time.Time) string {
	m := t.Minute()

	if m >= 58 {
		t = t.Add(time.Duration(60-m) * time.Minute)
		return t.Format("15h")
	}
	if m <= 2 {
		t = t.Add(-time.Duration(m) * time.Minute)
		return t.Format("15h")
	}

	return t.Format("15h04")
}

func FormatTooltip(data *ccusage.BlocksData) string {
	var sb strings.Builder

	resetTime := time.Now().Add(time.Duration(data.RemainingMinutes) * time.Minute)
	resetTimeStr := formatResetTime(resetTime)

	fmt.Fprintf(&sb, "<b>\uf1b2 Active Block (resets in %s - %s)</b>\n",
		FormatDuration(data.RemainingMinutes),
		resetTimeStr)

	fmt.Fprintf(&sb, "<b>\uf1d8 Requests:</b> %d\n",
		data.Entries)

	fmt.Fprintf(&sb, "<b>\uf145 Tokens:</b> %s (%s in / %s out)\n",
		FormatNumber(data.TotalTokens),
		FormatNumber(data.InputTokens),
		FormatNumber(data.OutputTokens))

	fmt.Fprintf(&sb, "<b>\uf155 Cost:</b> $%.2f @ $%.2f/h",
		data.CostUSD,
		data.CostPerHour)

	return sb.String()
}
