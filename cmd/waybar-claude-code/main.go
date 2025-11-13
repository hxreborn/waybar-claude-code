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

const (
	version      = "dev"
	iconStatic   = "󰜡"
	classLoading = "loading"
	classError   = "error"
	execTimeout  = 8 * time.Second
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("waybar-claude-code %s\n", version)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	sigCtx, sigCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer sigCancel()

	printFrame(iconStatic, "Loading Claude Code usage…", classLoading)

	tooltip, err := fetchTooltip(sigCtx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printFrame(iconStatic, "Unable to load stats", classError)
		os.Exit(0)
	}

	printFrame(iconStatic, tooltip, "")
	os.Exit(0)
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

	output := waybar.Output{
		Text:    icon,
		Tooltip: tooltip,
		Class:   class,
	}

	if err := output.PrintTo(writer); err != nil {
		fmt.Fprintf(os.Stderr, "waybar output error: %v\n", err)
		os.Exit(0)
	}

	if err := writer.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "flush error: %v\n", err)
		os.Exit(0)
	}
}
