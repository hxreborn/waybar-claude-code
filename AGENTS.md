# Repository Guidelines

## Project Structure & Modules
Core Go entrypoint lives in `cmd/waybar-claude-code`, while reusable logic is grouped under `internal/ccusage` (ccusage client), `internal/config` (env parsing), and `internal/format` (render helpers). Shared Waybar glue belongs in `pkg/waybar`. Assets for docs and demos live in `assets/`, and ready-to-copy Waybar + CSS samples are in `examples/`. Keep binaries and build artefacts out of Git; `waybar-claude-code` at the root is the compiled output only.

## Build, Test & Run
Use the Makefile for repeatable workflows:
- `make build` – produces a static `waybar-claude-code` binary with version metadata.
- `make install` – builds then drops the binary into `~/.config/waybar/modules/`.
- `make run RUN_INTERVAL=15` – runs the module with a given refresh cadence for local debugging.
- `make fmt`, `make lint`, `make test` – wrap `go fmt`, `go vet`+`staticcheck`, and `go test -race -cover ./...` respectively. Run them before any PR.

## Coding Style & Naming
Rely on `go fmt` defaults (tabs, single blank lines between top-level declarations). Exported APIs use UpperCamelCase, internal helpers stay lowerCamelCase. CLI flags and environment knobs follow the `CLAUDE_*` pattern; Waybar config keys mirror the JSON schema shown in `examples/waybar-config.jsonc`. Keep files small and cohesive: config parsing stays in `internal/config`, formatting logic in `internal/format`, and Waybar structs in `pkg/waybar`.

## Testing Expectations
Place `_test.go` beside the code under test, using `TestXxx` naming so `go test ./...` picks them up. Prefer table-driven tests for formatters and JSON serializers, and include integration coverage for the Waybar payload (see existing tests in `internal/...`). Aim to keep race detector clean and maintain at least the current coverage level (the `-cover` flag enforces regressions).

## File Locations for Development
When testing changes, modify the personal system config files first, then sync to examples:

**Personal system (live testing):**
- `~/.config/waybar/modules/custom-claude-code.jsonc` - module config
- `~/.config/waybar/user-style.css` - CSS styles
- `~/.config/waybar/modules/waybar-claude-code` - compiled binary (from `make install`)

**Repository examples (update after testing):**
- `examples/waybar-config.jsonc` - reference module config
- `examples/style.css` - reference CSS

Test workflow: edit personal files → `make build && make install` → `systemctl --user restart hyde-Hyprland-bar.service` (preferred for HyDE) or `pkill -SIGUSR2 waybar` → verify → copy to examples.

## Styling Approach
- Uses GTK3 CSS `:hover` pseudo-class for visual feedback
- Simple color transition on hover, no flag files or state management
- GTK3 supports `:hover`, `:active`, `:focus` (must be last element in selector)
- Working properties: `color`, `background-color`, `box-shadow`, `opacity`
- No `transform` or `font-size` animation support in GTK3 CSS

## Commit & PR Guidelines
Commits follow Conventional Commit prefixes (`feat`, `fix`, `docs`, `chore`, `test`, etc.) as seen in recent history. Keep messages imperative and scoped to a single concern. PRs should describe the user-facing impact, list test commands you ran, link any tracked issues, and include updated screenshots when UI output (icons, tooltip content) changes. Mention any new env vars or Waybar config knobs so downstream configs can be updated quickly.

## Future Enhancements
- **Adaptive polling:** convert the binary into a long-running process with a dual-interval ticker (5s inactive, CLAUDE_INTERVAL_SEC when active). Add `/proc/*/comm` scanning for instant-on detection, emit inactive vs active classes, and switch Waybar configs to `interval: "once"` + `restart-interval: 30`. Update user CSS/README to explain the new states.
