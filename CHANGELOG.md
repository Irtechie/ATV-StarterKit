# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.0.1] — 2026-04-01

Patch release to fix npm distribution. The v2.0.0 npm package was published before the corresponding GitHub release existed, so the postinstall binary download would fail. This release creates the GitHub release with pre-built binaries and publishes a new npm version that resolves correctly.

### Fixed

- **npm install works end-to-end** — v2.0.0 GitHub release now exists with goreleaser binaries for all platforms (macOS, Linux, Windows on amd64/arm64)
- **Version alignment** — VERSION file, npm package, and GitHub release tag are now all in sync at 2.0.1

## [2.0.0] — 2026-03-29

ATV 2.0 — **All The Vibes** becomes a unified installer combining Compound Engineering, gstack, and agent-browser into one agentic coding setup for GitHub Copilot.

### Added

- **gstack integration** — 30 gstack skills from [garrytan/gstack](https://github.com/garrytan/gstack) installable via the guided wizard. Clone, generate docs, copy skills, and create runtime sidecar — all handled by the Go installer.
- **agent-browser integration** — [Vercel agent-browser](https://github.com/vercel-labs/agent-browser) installable via guided mode. Auto-installs the CLI via npm and downloads Chrome for Testing.
- **Preset-based guided wizard** — Three presets (Starter / Pro / Full) replace the flat checkbox list. Beginners pick a preset; power users drill into category-grouped customization.
- **Animated install progress** — Bubbletea-powered progress display with per-step spinners (pending → running → done/failed) replacing the silent line-by-line output.
- **Retro terminal banner** — "ALL THE VIBES 2.0" in block letters with gold gradient and retro boot messages, matching the hero SVG.
- **Runtime sidecar** — `.github/skills/gstack/` directory with binaries (`bin/`), browser runtime (`browse/dist/`), review checklists, and `ETHOS.md` so gstack skills can find their runtime assets.
- **Memory system documentation** — copilot-instructions templates now include ATV Override Rules and conflict resolution between ATV and gstack memory systems.
- **Function-based TUI categories** — Skills organized by what they do (Planning, Review, QA, Security, Shipping, Safety) instead of where they come from.
- **`pkg/gstack/`** — New Go package for gstack installation: prerequisite detection, git clone, skill doc generation, sidecar creation.
- **`pkg/agentbrowser/`** — New Go package for agent-browser: npm install, Chrome download, SKILL.md fetch.
- **`pkg/tui/presets.go`** — Preset definitions (Starter, Pro, Full) with skill mappings.
- **`pkg/tui/progress.go`** — Bubbletea progress model with animated spinners.
- **Sandbox integration tests** — `test/sandbox/` with tests for auto mode, guided mode, gstack install, idempotency, and instructions content.

### Changed

- **README** rewritten as "The Three Pillars" — Compound Engineering (knowledge compounds), gstack (AI sprint process), agent-browser (eyes of the agent). Memory highlighted as a first-class differentiator.
- **Wizard flow** — Stack → Preset → Customize? → (optional) multi-select instead of Stack → flat checkbox wall.
- **copilot-instructions templates** — all 4 stacks (general, typescript, python, rails) now include gstack skill listing, agent-browser section, and ATV Override Rules.
- **`.gitignore`** — added `.gstack/`, `.env`, `*.db`, `.claude/` entries.
- **`cmd/init.go`** — guided mode now runs gstack and agent-browser install steps with progress display; auto mode unchanged.
- **`pkg/output/printer.go`** — new status indicators (🔗 Cloned, 🔨 Built), gstack/agent-browser progress methods, dynamic next-steps with step numbering.

### Fixed

- **gstack skill discovery** — skills placed at `.github/skills/gstack-*/SKILL.md` (one level deep) instead of nested `.github/skills/gstack/{name}/SKILL.md` which Copilot couldn't discover.
- **Windows CRLF** — bash scripts in gstack clone are auto-fixed (`\r\n` → `\n`) before execution.
- **Preset cursor** — default to Starter so all three options are visible in shorter terminals.

## [1.0.0] — 2026-03-11

Initial release of ATV Starter Kit.

### Added

- One-click installer (`atv-installer init`) for GitHub Copilot agentic coding environment
- All 6 Copilot lifecycle hooks: system instructions, setup steps, MCP config, skills, agents, file instructions
- 13 workflow skills (brainstorm, plan, work, review, compound, lfg, slfg, etc.)
- 28 specialized agents (security, performance, architecture, data, design, research)
- Stack detection: TypeScript, Python, Rails, General
- Interactive guided mode (`--guided`) with charmbracelet/huh TUI
- MCP servers: Context7, GitHub, Azure, Terraform
- npm distribution via `npx atv-starterkit`
- Cross-platform binary releases via goreleaser
