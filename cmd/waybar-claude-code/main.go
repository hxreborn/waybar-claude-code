package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
	"github.com/hxreborn/waybar-claude-code/internal/config"
	"github.com/hxreborn/waybar-claude-code/internal/format"
	"github.com/hxreborn/waybar-claude-code/pkg/waybar"
)

var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("waybar-claude-code %s\n", version)
		return
	}

	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	outputMetrics(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			outputMetrics(ctx)
		}
	}
}

func outputMetrics(ctx context.Context) {
	data, err := ccusage.GetBlocks(ctx)
	if err != nil {
		data = &ccusage.BlocksData{}
	}

	output := waybar.Output{
		Text:    format.FormatText(data),
		Tooltip: format.FormatTooltip(data),
	}

	_ = output.Print()
}
