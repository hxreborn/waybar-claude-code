package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
		os.Exit(0)
	}

	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if !cfg.Debug {
		log.SetOutput(io.Discard)
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	outputMetrics(ctx, cfg.Debug)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			outputMetrics(ctx, cfg.Debug)
		}
	}
}

func outputMetrics(ctx context.Context, debug bool) {
	data, err := ccusage.GetBlocks(ctx)
	if err != nil {
		if debug {
			log.Printf("Error getting blocks: %v", err)
		}
		data = &ccusage.BlocksData{}
	}

	text := format.FormatText(data)
	tooltip := format.FormatTooltip(data)

	output := waybar.Output{
		Text:    text,
		Tooltip: tooltip,
	}

	output.Print()
}
