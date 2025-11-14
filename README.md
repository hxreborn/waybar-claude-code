# waybar-claude-code

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg?logo=go)](https://go.dev/dl/)
[![Zero Deps](https://img.shields.io/badge/deps-zero-success)]()
![Nerd Fonts](https://img.shields.io/badge/nerd%20font-required-orange)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A lightweight Waybar custom module written in Go that displays Claude Code usage metrics to your bar using [ccusage](https://github.com/ryoppippi/ccusage).

## Demo

<p align="center">
  <img src="assets/demo.gif" width="800" alt="Demo">
</p>

> [!NOTE]
> If you're using waybar for the first time checkout the [example configuration](examples).

## Features

- **CGO-disabled static binary** - Single 3MB Go executable with zero runtime dependencies. Runs anywhere without glibc or library conflicts.
- **Live usage without daemons** - Executes ccusage on each poll to show current Claude Code stats without background processes hogging resources.
- **Waybar-native JSON output** - Built-in tooltip support and CSS class integration. Works seamlessly with your existing bar configuration and theme.
- **Microscopic overhead** - Spawns per interval, fetches stats and exits. ~10MB memory footprint with effectively zero CPU load.
- **Configurable refresh rate** - Polls every 5 minutes by default. Change one config value to update faster or slower.

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
"custom/claude-code": {
  "return-type": "json",
  "exec": "~/.config/waybar/modules/waybar-claude-code",
  "format": "{text}",
  "interval": 300,
  "tooltip": true,
  "on-click": "${TERMINAL:-kitty} -e claude"
}
```

### Style Configuration

Add to `~/.config/waybar/style.css`:

```css
#custom-claude-code {
  padding: 0 10px;
  margin: 0 2px;
  color: inherit;
}

#custom-claude-code:hover {
  color: #ff8c00;
}

tooltip {
  font-family: "MesloLGSDZ Nerd Font", monospace;
}
```

**Note:** Tooltip icons require a [Nerd Font](https://www.nerdfonts.com/). See [examples/](examples/) for complete styling with active states and tooltip customization.

## Usage

Hover over the 󰜡 icon to see detailed metrics:

```
 Active Block (resets in 3h 45m - 18h30)
 Requests: 110
 Tokens: 2.8M (1.5M in / 1.3M out)
 Cost: $1.47 @ $0.38/h
```

Reset time uses 24-hour format and rounds to the nearest hour if within 2 minutes (e.g., `17:59` → `18h`).

## License

MIT - see [LICENSE](LICENSE) file for details
