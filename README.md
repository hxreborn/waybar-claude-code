# waybar-claude-code

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://go.dev/dl/)
[![Zero Deps](https://img.shields.io/badge/deps-zero-success)]()
![Nerd Fonts](https://img.shields.io/badge/nerd%20font-required-orange)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A lightweight Waybar custom module written in Go that displays Claude Code usage metrics to your bar using [ccusage](https://github.com/ryoppippi/ccusage).

## Screenshots

![Waybar module in action](assets/screenshot-1.png)

![Tooltip display](assets/screenshot-2.png)

> [!NOTE]
> Check out the [example configuration](examples/) for CSS and Waybar config samples.

## Features

- **CGO-disabled static binary** - Pure Go compilation, no glibc or external dependencies
- **One-shot execution** - Runs once per poll, exits immediately after output
- **Minimal footprint** - <10MB RSS, <3MB binary with near-zero idle CPU usage
- **Battery-efficient** - No background processes or continuous resource usage
- **Configurable polling** - Refresh interval controlled by Waybar settings

## Requirements

- [npm/npx](https://nodejs.org/) to run ccusage
- [Waybar](https://github.com/Alexays/Waybar) for module integration
- [Nerd Fonts](https://www.nerdfonts.com/) for icon display
- Go 1.21 or later required for source compilation

## Installation

### Precompiled Binary Installation

Download and install to your Waybar modules directory:

```bash
curl -LO https://github.com/hxreborn/waybar-claude-code/releases/latest/download/waybar-claude-code
install -Dm755 waybar-claude-code ~/.config/waybar/modules/waybar-claude-code
```

### Build from Source

Clone and build:

```bash
git clone https://github.com/hxreborn/waybar-claude-code.git /tmp/waybar-claude-code
cd /tmp/waybar-claude-code
make install
```

## Configuration

### Update Waybar Config

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
    "exec": "~/.config/waybar/modules/waybar-claude-code",
    "format": "{text}",
    "interval": 300,
    "restart-interval": 30,
    "tooltip": true,
    "on-click": "kitty -e claude"
  }
}
```

**Config Options:**
- `interval: 300` - Polling interval in seconds
- `restart-interval: 30` - Auto-restart interval in seconds
- `tooltip: true` - Enable hover tooltip with detailed metrics
- `on-click: "kitty -e claude"` - Launch Claude Code in Kitty terminal

### Style Configuration

Add to `~/.config/waybar/style.css`:

```css
#custom-claude-code {
  padding: 0 10px;
  margin: 0 2px;
  color: inherit;
  transition: color 0.2s ease-in-out;
}

#custom-claude-code:hover {
  color: #ff8c00;
}
```

See [examples/](examples/) for basic configuration and styling examples.

## Usage

Hover over the module to see detailed metrics:

```
Requests: 42 | Tokens: 1.2M
Cost: $0.45 | Reset: 2h 15m
```

## License

MIT - see [LICENSE](LICENSE) file for details
