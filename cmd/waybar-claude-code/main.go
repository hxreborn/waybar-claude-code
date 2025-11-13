package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
	"github.com/hxreborn/waybar-claude-code/internal/format"
	"github.com/hxreborn/waybar-claude-code/pkg/waybar"
)

const (
	version      = "dev"
	iconStatic   = "󰜡"
	classLoading = "loading"
	classError   = "error"
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

	printFrame(iconStatic, "Loading Claude Code usage…", classLoading)

	tooltip, err := fetchTooltip(ctx)
	if err != nil {
		printFrame(iconStatic, "Unable to load stats", classError)
		return
	}

	printFrame(iconStatic, tooltip, "")
}

func fetchTooltip(ctx context.Context) (string, error) {
	data, err := ccusage.GetBlocks(ctx)
	if err != nil {
		return "", err
	}
	return format.FormatTooltip(data), nil
}

func printFrame(icon, tooltip, class string) {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	output := waybar.Output{
		Text:    icon,
		Tooltip: tooltip,
		Class:   class,
	}

	if err := output.PrintTo(writer); err != nil {
		fmt.Fprintf(os.Stderr, "waybar output error: %v\n", err)
	}
}
