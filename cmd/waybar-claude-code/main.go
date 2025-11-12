package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
	"github.com/hxreborn/waybar-claude-code/internal/format"
	"github.com/hxreborn/waybar-claude-code/pkg/waybar"
)

var (
	version    = "dev"
	iconStatic = "󰜡"
	writer     = bufio.NewWriter(os.Stdout)
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("waybar-claude-code %s\n", version)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	printStatus(iconStatic, "Loading Claude Code usage…", "loading")
	outputCycle(ctx)
}

func outputCycle(ctx context.Context) {
	icon := iconStatic

	// Give ccusage more time to fetch data
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	data, err := ccusage.GetBlocks(ctx)
	if err != nil {
		printStatus(icon, "Unable to load stats", "error")
		return
	}

	tooltip := format.FormatTooltip(data)

	printStatus(icon, tooltip, "")
}

func printStatus(icon, tooltip, class string) {
	output := waybar.Output{
		Text:       icon,
		Tooltip:    tooltip,
		Percentage: 0,
	}

	if class != "" {
		output.Class = class
	}

	_ = output.PrintTo(writer)
	writer.Flush()
}
