package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hxreborn/waybar-claude-code/internal/ccusage"
	"github.com/hxreborn/waybar-claude-code/internal/config"
	"github.com/hxreborn/waybar-claude-code/internal/format"
	"github.com/hxreborn/waybar-claude-code/pkg/waybar"
)

var (
	version = "dev"
	frames  = spinner.CharSets[14]
	frame   int
)

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

	outputMetrics(ctx, cfg.Animate)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			outputMetrics(ctx, cfg.Animate)
		}
	}
}

func outputMetrics(ctx context.Context, animate bool) {
	data, err := ccusage.GetBlocks(ctx)
	if err != nil {
		data = &ccusage.BlocksData{}
	}

	text := format.FormatText(data)
	if animate {
		text = frames[frame%len(frames)]
		frame++
	}

	output := waybar.Output{
		Text:    text,
		Tooltip: format.FormatTooltip(data),
	}

	_ = output.Print()
}
