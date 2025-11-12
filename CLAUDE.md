# waybar-claude-code

Waybar module displaying Claude Code usage statistics. Zero external dependencies, single binary output.

## Architecture

**Single execution model**: Module runs once per invocation, outputs JSON, exits. Waybar handles polling via `interval` config.

**Why one-shot execution:**
- Simpler process lifecycle (no ticker management)
- Waybar already provides interval mechanism
- Lower memory footprint
- Easier debugging

## Waybar JSON Protocol

Custom modules with `"return-type": "json"` expect single-line JSON on stdout:

```json
{"text": "display", "tooltip": "hover text", "class": "css-class", "percentage": 50}
```

**Protocol rules:**
- Single line output (no pretty printing)
- stdout for JSON only
- stderr for errors/logs
- Exit immediately after output

## Configuration

**Waybar config:**
```jsonc
{
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

## CSS Styling

**GTK CSS limitations:**
- No `transform` property (no scale/rotate/translate)
- No `font-size` animation support
- Working properties: `color`, `background-color`, `box-shadow`, `opacity`
- Supported pseudo-classes: `:hover`, `:active`, `:focus` (must be last element in selector)

**Hover effect:**
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

## HyDE Integration

HyDE manages waybar via Python watcher and systemd service.

**Key paths:**
- Layouts: `~/.local/share/waybar/layouts/hyprdots/`
- Styles: `~/.local/share/waybar/styles/`
- State: `~/.local/state/hyde/staterc`
- Config: `~/.config/waybar/config.jsonc` (auto-generated, do not edit)

**Reload after layout changes:**
```bash
~/.local/lib/hyde/waybar.py --set hyprdots/04-claude-stats
```

**NEVER use `pkill waybar && waybar &`** - creates duplicate bars.

## Build

```bash
make build    # Build binary
make install  # Install to ~/.config/waybar/modules/
```

**Build flags:**
- `CGO_ENABLED=0` - static binary
- `-ldflags="-s -w"` - strip debug symbols
- `-X main.version` - embed version from git

## Module Structure

```
cmd/waybar-claude-code/main.go  # Entry point
internal/ccusage/               # Claude Code usage API client
internal/config/                # Config loading
internal/format/                # Tooltip formatting
pkg/waybar/                     # Waybar JSON output
```

**Kept multiple packages because:**
- ccusage logic is complex (API client, caching, timeout handling)
- Format logic has specific business rules
- pkg/waybar is reusable across waybar modules

## Testing

```bash
# Output validation
./waybar-claude-code | jq .

# Live test in waybar
~/.config/waybar/modules/waybar-claude-code
```

## File Locations

**Personal system files (for testing changes):**
- Config: `~/.config/waybar/modules/custom-claude-code.jsonc`
- CSS: `~/.config/waybar/user-style.css`
- Binary: `~/.config/waybar/modules/waybar-claude-code`

**Repository example files (update after testing):**
- Config: `examples/waybar-config.jsonc`
- CSS: `examples/style.css`

**Workflow:**
1. Modify personal system files to test changes
2. Rebuild with `make build && make install`
3. Reload waybar: `systemctl --user restart hyde-Hyprland-bar.service` (preferred for HyDE) or `pkill -SIGUSR2 waybar`
4. If working, copy changes to example files

## Known Issues

**GTK CSS property support:**
- No `transform`: scale/rotate/translate not supported
- No `font-size` animation: causes rendering issues
- Working properties: `background-color`, `box-shadow`, `opacity`, `color`
- Waybar CSS parser uses GTK3 CSS spec, not standard CSS

## Performance

- Startup: ~50ms (includes API call with 8s timeout)
- Memory: <10MB RSS
- Binary: <3MB stripped
- CPU: negligible (one-shot execution)

## References

- [Waybar custom modules](https://github.com/Alexays/Waybar/wiki/Module:-Custom)
- [GTK CSS properties](https://docs.gtk.org/gtk3/css-properties.html)
- [HyDE documentation](https://github.com/HyDE-Project/HyDE)
