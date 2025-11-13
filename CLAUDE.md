# Waybar Minimal Module - Planning & Design Document

## Executive Summary

A ground-up rewrite of a Waybar custom module focusing on simplicity, maintainability, and zero dependencies. The goal is to create a small, reliable module that outputs JSON for Waybar consumption with minimal overhead.

## Research Findings

### Waybar JSON Protocol

From the research, Waybar custom modules with `"return-type": "json"` expect a single line JSON object with these fields:

```json
{
  "text": "display text",           // Required: main display text
  "tooltip": "hover text",           // Optional: tooltip on hover
  "class": "css-class",              // Optional: CSS class(es) - string or array
  "percentage": 50,                  // Optional: 0-100, drives {icon} selection
  "alt": "alternative text"          // Optional: alternative text for {alt} format
}
```

**Key Protocol Rules:**
- One JSON object per line (no pretty printing)
- Output to stdout only
- Errors/logs must go to stderr
- Module can run in two modes:
  - Continuous: loop with sleep, output JSON on each iteration
  - Once: output once and exit (Waybar handles refresh via signals)

### Existing Module Analysis

**waybar-claude-stats (Complex)**
- 270+ lines in main.go alone
- Multiple packages (internal/i18n, state management)
- Complex caching system with ccusage
- 7 different display modes with cycling
- Persistent state across restarts
- Heavy structure but feature-rich

**waybar-lyric (Moderate)**
- Uses external framework (fang)
- Multiple dependencies (cobra-like CLI)
- Clean separation but adds complexity
- Good error handling patterns

**Key Lessons:**
- State persistence adds significant complexity
- External dependencies increase build size
- Simple polling loop is often sufficient
- Error handling should be silent (stderr only)

### Go Best Practices for Minimal Modules

1. **Structure**: Keep it flat for simple tools
   - Single main.go for < 300 lines
   - pkg/ only if truly reusable
   - Avoid internal/ for small projects

2. **Dependencies**: Zero is best
   - Use only stdlib when possible
   - No external CLI frameworks
   - Build with CGO_ENABLED=0 for static binaries

3. **Cross-compilation**: Simple with Go
   ```bash
   GOOS=linux GOARCH=amd64 go build
   GOOS=linux GOARCH=arm64 go build
   ```

## Design Decisions

### Language Choice: Go

**Why Go:**
- Single static binary (no runtime deps)
- Fast startup time
- Low memory footprint
- Easy cross-compilation
- Strong stdlib for JSON, HTTP, system info

**Alternatives Considered:**
- Python: Requires runtime, slower startup
- Rust: Overkill for simple JSON output
- Bash+jq: Fragile, harder to maintain

### Architecture: Minimalist

```
waybar-temp/
├── main.go          # Single file implementation
├── Makefile         # Build automation
├── CLAUDE.md        # This document
└── assets/
    └── waybar.css   # Example styles
```

**Why Single File:**
- Module is ~150-200 lines max
- No need for package separation
- Easier to understand and modify
- Faster compilation

### Update Model: Configurable

```go
// Default: Loop with interval
for {
    output := collectMetrics()
    json.NewEncoder(os.Stdout).Encode(output)
    time.Sleep(interval)
}

// Alternative: Once mode
output := collectMetrics()
json.NewEncoder(os.Stdout).Encode(output)
// Exit, let Waybar handle refresh via signals
```

### Metric Selection

**Primary: System Temperature**
- Universal, useful metric
- Easy to read from /sys/class/thermal/
- Good for demonstrating percentage (0-100°C)
- CSS can color based on heat levels

**Fallback Options:**
- CPU usage (via /proc/stat)
- Memory usage (via /proc/meminfo)
- Network status (ping gateway)

### Configuration Strategy

**CLI Flags (minimal):**
```bash
waybar-temp -interval 60    # Poll interval (0 = once)
waybar-temp -zone thermal_zone0  # Thermal zone to monitor
waybar-temp -debug          # Show debug info to stderr
```

**Environment Variables:**
```bash
WAYBAR_TEMP_ZONE=thermal_zone1
WAYBAR_TEMP_INTERVAL=30
```

**No Config Files:** Keep it simple, use CLI/env only

## Implementation Plan

### Phase 1: Core Functionality
```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "time"
)

type Output struct {
    Text       string      `json:"text"`
    Tooltip    string      `json:"tooltip,omitempty"`
    Class      interface{} `json:"class,omitempty"`
    Percentage int         `json:"percentage,omitempty"`
}

func main() {
    // Parse flags
    // Read metrics
    // Output JSON
    // Loop or exit
}
```

### Phase 2: Metric Collection
```go
func readTemperature(zone string) (float64, error) {
    path := fmt.Sprintf("/sys/class/thermal/%s/temp", zone)
    data, err := os.ReadFile(path)
    if err != nil {
        return 0, err
    }
    // Parse and convert millidegree to degree
    return temp / 1000.0, nil
}
```

### Phase 3: Error Handling
- Silent failures (return default values)
- Log to stderr only if debug flag set
- Never crash, always output valid JSON

### Phase 4: Build System
```makefile
BINARY=waybar-temp
VERSION=$(shell git describe --tags --always --dirty)

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY)

install: build
	install -Dm755 $(BINARY) ~/.config/waybar/modules/

clean:
	rm -f $(BINARY)
```

## Migration from waybar-claude-stats

### What to Keep:
- JSON output structure
- Error handling to stderr
- Version embedding in build

### What to Discard:
- State persistence (unnecessary complexity)
- Mode cycling (over-engineered)
- Caching layer (premature optimization)
- i18n (YAGNI for personal tool)
- Multiple packages (too much structure)

### Simplifications:
- No external commands (ccusage) - read directly from /sys or /proc
- No state files - stateless operation
- No complex CLI framework - just flag package
- No custom types for modes - one simple output

## Waybar Configuration

### Basic Setup
```json
{
  "modules-right": ["custom/temp"],
  "custom/temp": {
    "return-type": "json",
    "format": "{icon} {text}",
    "format-icons": ["", "", "", "", ""],
    "interval": 60,
    "exec": "~/.config/waybar/modules/waybar-temp",
    "tooltip": true
  }
}
```

### Signal-based Updates
```json
{
  "custom/temp": {
    "return-type": "json",
    "format": "{icon} {text}",
    "interval": "once",
    "signal": 8,
    "exec": "~/.config/waybar/modules/waybar-temp -interval 0"
  }
}
// Refresh: pkill -SIGRTMIN+8 waybar
```

## CSS Styling Examples

```css
/* Basic styling */
#custom-temp {
    padding: 0 10px;
    color: #ffffff;
}

/* Temperature-based colors */
#custom-temp.cold { color: #5e81ac; }
#custom-temp.normal { color: #a3be8c; }
#custom-temp.warm { color: #ebcb8b; }
#custom-temp.hot { color: #bf616a; }

/* Icon styling */
#custom-temp.critical {
    animation: blink 1s linear infinite;
}

@keyframes blink {
    50% { opacity: 0; }
}
```

## Testing Strategy

### Manual Testing
```bash
# Test output format
./waybar-temp | jq .

# Test continuous mode
./waybar-temp -interval 2 | head -5

# Test error handling
WAYBAR_TEMP_ZONE=invalid ./waybar-temp

# Memory usage
/usr/bin/time -v ./waybar-temp -interval 0
```

### Integration Testing
1. Copy to Waybar modules directory
2. Add to Waybar config
3. Restart Waybar
4. Verify tooltip, styling, updates

## Performance Targets

- Startup time: < 10ms
- Memory usage: < 5MB RSS
- CPU usage: < 0.1% when polling
- Binary size: < 2MB stripped

## Quick Start Guide

```bash
# Clone and enter directory
cd waybar-temp

# Build
make build

# Test locally
./waybar-temp

# Install to Waybar
make install

# Add to Waybar config (see examples above)
# Restart Waybar
pkill -SIGUSR2 waybar
```

## Maintenance Checklist

- [ ] Keep dependencies at zero
- [ ] Maintain single-file simplicity
- [ ] Test on different thermal zones
- [ ] Verify JSON output format
- [ ] Check binary size stays small
- [ ] Document any deviations

## Future Considerations (YAGNI)

These are explicitly NOT in scope but noted for reference:
- Multiple metrics in one module
- Configuration file support
- Metric history/graphing
- Network API calls
- Database storage
- Plugin system

## Decision Log

1. **2024-11-11**: Chose temperature monitoring as primary metric
   - Reason: Universal, useful, simple to implement
   
2. **2024-11-11**: Single file over package structure
   - Reason: Sub-200 lines doesn't need packages
   
3. **2024-11-11**: No state persistence
   - Reason: Unnecessary complexity for display module
   
4. **2024-11-11**: Go over Python/Rust/Bash
   - Reason: Best balance of simplicity and performance

## References

- [Waybar Wiki - Custom Modules](https://github.com/Alexays/Waybar/wiki/Module:-Custom)
- [Waybar man page](https://man.archlinux.org/man/extra/waybar/waybar-custom.5.en)
- [Go Modules](https://go.dev/ref/mod)
- [Linux Thermal Sysfs](https://www.kernel.org/doc/Documentation/thermal/sysfs-api.txt)
- Existing examples: ~/Documents/GitHub/waybar-things/

## Next Steps

1. Implement basic temperature reader
2. Add JSON output
3. Add interval loop
4. Add CLI flags
5. Create Makefile
6. Write installation docs
7. Test with actual Waybar
8. Add CSS examples
9. Performance optimization
10. Tag v1.0.0

---

*This document represents the planning phase for a minimal Waybar module. Implementation should follow these guidelines to maintain simplicity and reliability.*
