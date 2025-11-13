# Repository Guidelines

## Project Structure & Module Organization
Source for the Waybar binary lives in `cmd/waybar-claude-code`, while reusable helpers (JSON marshalling, tooltip builders) stay under `pkg/waybar` such as `output.go`. `assets/` contains screenshots referenced in the README, and `examples/waybar/` carries config + CSS snippets. House research notes inside `DEVELOPMENT.md` or `ANALYSIS_SUMMARY.md` so the Go tree remains small and cache-friendly. If you add new ccusage adapters or caches, stage them inside `internal/` packages to keep the CLI thin.

## Build, Test, and Development Commands
- `go build ./cmd/waybar-claude-code` builds a local binary; add `-o ~/.config/waybar/modules/waybar-claude-code` when installing manually.
- `make install` (target described in README) wraps the optimized CGO-disabled build and copies it into the Waybar modules path.
- `CLAUDE_INTERVAL_SEC=5 go run ./cmd/waybar-claude-code` shortens the ticker for development and mirrors Waybar’s stdout contract.
- `go vet ./...`, `staticcheck ./...`, and `golangci-lint run` catch API regressions before you open a PR.

## Coding Style & Naming Conventions
Stick to Go 1.21+, tabs, and whatever `gofmt` + `goimports` emit; never hand-align JSON tags. Keep package names short and nouns (`ccusage`, `tooltip`), export only what Waybar needs, and document non-obvious state machines with single-line comments. Follow the design mantra from `DEVELOPMENT.md`: explicit loops over clever abstractions.

## Testing Guidelines
Place `*_test.go` beside the code under test and name functions `TestThing_Scenario`. Run `go test ./...` before every push, upgrade to `go test -race -cover ./...` when touching concurrency, and lean on env vars for edge cases (`CLAUDE_TIMEOUT_SEC=0 go test ./...` forces timeout paths). Snapshot tooltip text under `testdata/` to guard against accidental formatting drift.

## Commit & Pull Request Guidelines
Project history (see the upstream GitHub mirror) uses scope-prefixed Conventional Commits such as `fix: clamp tooltip width` or `docs: refresh README`; keep subjects imperative and under 72 characters. Reference the motivating issue in the body, list the commands you ran, and attach refreshed screenshots whenever UI output changes. PRs should mention any env vars required for reviewers to reproduce your test run.

## Security & Configuration Tips
This module shells out to `npx ccusage@latest`, so avoid hard-coding paths or tokens; rely on env vars (`CLAUDE_CACHE_TTL_SEC`, `CLAUDE_TIMEOUT_SEC`) instead. Restart Waybar via `pkill -SIGUSR2 waybar` after installing new binaries, and keep Nerd Fonts available so the `󰜡` glyph renders correctly. Document any new configuration knobs in README and `examples/waybar/config.jsonc` immediately.
