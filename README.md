# waybar-claude-code

> Minimal Waybar custom module for Claude Code usage tracking

A lightweight Waybar module that displays Claude Code API usage metrics from ccusage.

## Features

- Simple polling daemon with configurable interval
- Real-time request count and cost tracking
- Minimal footprint (<5MB memory, <1% CPU)
- Zero dependencies (static Go binary)
- Graceful error handling

## Requirements

- Go 1.21+ (build only)
- npm/npx (runtime - for ccusage)
- Waybar

## Installation

```bash
git clone https://github.com/hxreborn/waybar-claude-code.git
cd waybar-claude-code
make install
```

Installs to `~/.config/waybar/modules/waybar-claude-code`

## Configuration

### Waybar Config

Add to `~/.config/waybar/config.jsonc`:

```jsonc
{
  "modules-right": ["custom/claude-code"],

  "custom/claude-code": {
    "return-type": "json",
    "exec": "~/.config/waybar/modules/waybar-claude-code",
    "interval": 300,
    "tooltip": true
  }
}
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAUDE_INTERVAL_SEC` | `300` | Poll interval (seconds) |
| `CLAUDE_DEBUG` | `false` | Enable debug logging |

Example with custom interval:

```jsonc
"custom/claude-code": {
  "return-type": "json",
  "exec": "env CLAUDE_INTERVAL_SEC=60 ~/.config/waybar/modules/waybar-claude-code",
  "tooltip": true
}
```

## Output Format

**Display:** `Reqs: 42 | $0.45`

**Tooltip:**
```
Requests: 42 | Tokens: 1.2M
Cost: $0.45 | Reset: 2h 15m
```

## How It Works

1. Polls `npx ccusage@latest blocks --active --json --offline` every 5 minutes
2. Extracts request count, token usage, and cost
3. Outputs JSON to stdout for Waybar consumption

## Development

```bash
make build    # Build binary
make test     # Run tests
make fmt      # Format code
make clean    # Remove binary
```

## License

MIT
