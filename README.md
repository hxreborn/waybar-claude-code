# waybar-claude-code

> Waybar custom module for monitoring Claude Code usage metrics in real-time

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://go.dev/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/hxreborn/waybar-claude-code)](https://goreportcard.com/report/github.com/hxreborn/waybar-claude-code)
[![Release](https://img.shields.io/github/v/release/hxreborn/waybar-claude-code)](https://github.com/hxreborn/waybar-claude-code/releases)
[![Stars](https://img.shields.io/github/stars/hxreborn/waybar-claude-code?style=social)](https://github.com/hxreborn/waybar-claude-code)
[![Made with Love](https://img.shields.io/badge/made%20with-%E2%9D%A4-red.svg)](https://github.com/hxreborn/waybar-claude-code)

A lightweight Waybar module that displays Claude Code usage metrics for Linux systems.

## Demo

![Waybar module in action](assets/screenshot-1.png)

![Tooltip display](assets/screenshot-2.png)

> [!NOTE]
> If you're using Waybar for the first time, check out the [example configuration](examples/waybar/).

## Features

- **Icon-only display** - Clean, minimal status bar presence with detailed tooltip
- **Real-time metrics** - Token usage, cache efficiency, and live request counts
- **Accurate pricing** - Uses `npx ccusage@latest` for current pricing data
- **Long-running daemon** - Efficient internal ticker, no process restart overhead
- **Configurable polling** - Tune `CLAUDE_INTERVAL_SEC` for slower/faster refresh cycles
- **Request + cost tracking** - Surfaces live request counts and block spend in the tooltip
- **Tiny footprint** - <5 MB RSS and idle CPU under 1%, even on low-power systems
- **Zero external deps** - Single static Go binary; only `npx ccusage@latest` is invoked at runtime
- **Automatic recovery** - Waybar 0.9.24+ auto-restarts on crash
- **Zero configuration** - Works out of the box; env vars handle power users

## Installation

**Requirements:**
- npm/npx - Required to run `ccusage@latest`
- Waybar 0.9.24+ - For `restart-interval` and `hide-empty-text` features
- Nerd Fonts - For icon display (`󰜡` glyph)
- Go 1.21+ - For building from source

**Option 1: Prebuilt binaries** (recommended)

Download the latest release:

```bash
# Download latest release (replace VERSION with actual version)
curl -LO https://github.com/hxreborn/waybar-claude-code/releases/latest/download/waybar-claude-code-linux-amd64

# Install
install -Dm755 waybar-claude-code-linux-amd64 ~/.config/waybar/modules/waybar-claude-code
```

**Option 2: Build from source**

```bash
git clone https://github.com/hxreborn/waybar-claude-code.git
cd waybar-claude-code
make install
```

Installs to `~/.config/waybar/modules/waybar-claude-code`

**Option 3: Manual build**

```bash
CGO_ENABLED=0 go build \
  -trimpath \
  -ldflags "-s -w" \
  -o waybar-claude-code \
  ./cmd/waybar-claude-code

install -Dm755 waybar-claude-code ~/.config/waybar/modules/waybar-claude-code
```

**After installation, add to your Waybar config and restart:**

```bash
pkill -SIGUSR2 waybar
```

## Configuration

### Waybar Config

Add to `~/.config/waybar/config.jsonc`:

```jsonc
{
  "modules-right": [
    "custom/claude-code",
    "pulseaudio",
    "network",
    "clock"
  ],

  "custom/claude-code": {
    "return-type": "json",
    "format": "{icon}",
    "format-icons": ["󰜡"],
    "interval": "once",
    "restart-interval": 30,
    "exec": "~/.config/waybar/modules/waybar-claude-code",
    "tooltip": true,
    "escape": false,
    "hide-empty-text": true,
    "on-click": "sh -c 'claude'",
    "on-click-right": "xdg-open https://console.anthropic.com/settings/limits"
  }
}
```

### CSS Styling

Add to `~/.config/waybar/style.css`:

```css
#custom-claude-code {
  font-size: 0.95em;
  margin: 0 6px;
  padding: 2px 10px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.06);
  transition: background-color 0.2s ease;
}

#custom-claude-code:hover {
  background: rgba(0, 0, 0, 0.12);
}

#custom-claude-code.claude.error {
  background: rgba(220, 53, 69, 0.18);
  color: #dc3545;
}
```

### Environment Variables

Customize via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAUDE_INTERVAL_SEC` | `300` | Poll interval in seconds (set `0` for dev loops) |
| `CLAUDE_DEBUG` | `false` | Emit verbose stderr logs while running |

## Usage

### Tooltip Content

Hover over the module to see detailed metrics:

```
CLAUDE CODE · BLOCK METRICS
⟐ Tokens     1.12M / 2.00M  (56%)
💲 Cost       $0.57 this block
✉ Messages   51

EFFICIENCY & PERFORMANCE
⟡ Cache Hit  94%
⚡ Burn Rate  $0.33/hr · 11k tok/min

TOKEN BREAKDOWN
↧ Input      45k
↥ Output     20k
⇣ Cache Read 980k
⇡ Cache Write 60k

RATE LIMIT WINDOW
Progress   [██████████░░░░░░░░]  43%
Remaining  2h 52m of 5h
Resets     19:00

Updated 6s ago
```

## Troubleshooting

### Module doesn't appear

**Check:** Is npm installed?
```bash
which npm npx
```

**Check:** Is ccusage accessible?
```bash
npx ccusage@latest blocks --active --json --offline
```

**Check:** Are Waybar logs showing errors?
```bash
journalctl -f -t waybar
```

### Module shows error state (red background)

**Possible causes:**
- ccusage not found in PATH
- ccusage timeout (increase `CLAUDE_TIMEOUT_SEC`)
- Network issues preventing `npx` from downloading package

**Debug:**
```bash
~/.config/waybar/modules/waybar-claude-code 2>error.log
cat error.log
```

### Icon shows placeholder (□ or ?)

**Issue:** Nerd Fonts not installed or not configured in Waybar

**Fix:**
```bash
# Install Nerd Fonts (Arch example)
yay -S ttf-nerd-fonts-symbols-mono

# Update Waybar style.css
* {
  font-family: "Your Font", "Symbols Nerd Font Mono";
}
```

### Tooltip doesn't update

**Check:** Is module actually running?
```bash
pgrep -a waybar-claude-code
```

**Fix:** Restart Waybar
```bash
pkill waybar && waybar &
```

### High memory usage

**Expected:** <5 MB RSS for long-running process

**Check:**
```bash
ps aux | grep waybar-claude-code
```

If significantly higher, file an issue with details.

## Data

This module uses `ccusage` from [ryoppippi/ccusage](https://github.com/ryoppippi/ccusage) to fetch usage data.

**How it works:**

```
┌─────────┐   JSON    ┌────────────────────┐   JSON   ┌────────┐
│ ccusage │──stdout──→│ waybar-claude-code │──stdout─→│ Waybar │
│ @latest │           │                    │          └────────┘
└─────────┘           │  Parse & Format    │               ↓
                      │  Build Tooltip     │         User sees
                      │  Emit JSON         │         metrics
                      └────────────────────┘
```

1. Internal ticker wakes every 15 seconds
2. Executes `npx ccusage@latest blocks --active --json --offline`
3. Executes `npx ccusage@latest today --json --offline`
4. Parses JSON responses, extracts metrics
5. Formats multi-section tooltip with ASCII art
6. Emits single-line JSON to stdout
7. Waybar reads JSON, updates module display

**Why npx @latest?**
- Fresh bundled pricing data (30x more accurate than stale global install)
- Automatic updates without manual intervention
- Only 745ms overhead once per minute due to caching

**Why long-running process?**
- No restart overhead (Go runtime + npx initialization happens once)
- Persistent in-memory cache (60s TTL)
- ~0% CPU when sleeping between updates

## License

MIT License - see [LICENSE](LICENSE) file for details.
