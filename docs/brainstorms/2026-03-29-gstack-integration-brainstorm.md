---
date: 2026-03-29
topic: gstack-integration
---

# Integrate gstack into ATV StarterKit

## What We're Building

A new "gstack" layer in the ATV StarterKit installer that bundles [garrytan/gstack](https://github.com/garrytan/gstack) skills alongside our existing ATV skills. Users select gstack components through the existing guided TUI wizard, just like they pick ATV skills, agents, and MCP servers today.

gstack is Garry Tan's open-source "virtual engineering team" — 29 slash-command skills for Claude Code that cover planning (/office-hours), review (/review), QA (/qa with real browser), security (/cso), shipping (/ship), and safety guardrails (/careful, /freeze, /guard). It requires Bun and optionally Playwright for browser-based skills.

## Why This Approach

**Considered approaches:**

### Approach A: Live Git Clone (rejected)
Clone gstack into `.claude/skills/gstack` at install time with `.git` intact.
- **Pros:** Users can `git pull` for updates, matches gstack's own install docs.
- **Cons:** Couples our installer to gstack's git availability, `.git` dir in project, harder to pin versions.

### Approach B: Vendor into ATV Templates (rejected)
Copy all gstack SKILL.md files into `pkg/scaffold/templates/skills/gstack-*/`.
- **Pros:** Fully self-contained, no network dependency at install time beyond our binary.
- **Cons:** 29 skills would bloat our binary, gstack uses a TypeScript build step we can't replicate in Go, browser skills need Bun runtime anyway.

### Approach C: Shallow Clone + Go-Packaged Install (chosen)
Shallow clone gstack at install time via Go's `os/exec` (calling `git clone --depth 1`), strip `.git`, then build the gstack TypeScript binary by shelling out to `bun run build`. The Go installer handles all cross-platform logic — no bash dependency.
- **Pros:** Pinned to latest at install time, no `.git` overhead after clone, full gstack functionality including browser skills, Go installer handles Windows natively.
- **Cons:** Requires git + Bun/Node on user's machine, requires network at install time.

## Key Decisions

1. **Single Go installer**: Keep the current Go binary as the one unified installer. No separate gstack install step — the Go installer handles everything (clone, build, file placement).

2. **ATV overrides gstack for overlapping skills**: Where both provide similar functionality (e.g., ATV `/ce-review` vs gstack `/review`), ATV skills take priority. gstack skills are additive — they fill gaps ATV doesn't cover.

3. **Categorized TUI by function, not source**: The guided wizard groups skills by what they do (Planning, Review, QA, Security, Deploy, Safety), not where they come from. Best-of from both ATV and gstack appear together. Users don't need to know or care which is ATV vs gstack.

4. **Two sub-modes for gstack**:
   - **Markdown-only**: Just copy SKILL.md files. Works without Bun/Playwright. Covers planning, review, and non-browser skills.
   - **Full runtime**: Clone + `bun run build`. Enables /browse, /qa (browser), /connect-chrome, /benchmark. Requires Bun on user's machine.

5. **Installation target**: gstack files go to `.github/skills/gstack/` (VS Code Copilot convention), matching our existing skill placement.

6. **Build at install time**: gstack's TypeScript binary is compiled on the user's machine via `bun run build` when full runtime mode is selected. Go installer shells out to Bun (or Node on Windows as fallback).

## Implementation Sketch

### catalog.go changes
- Reorganize layers by function instead of source:
  - `"planning"` — ATV: ce-brainstorm, ce-plan, deepen-plan | gstack: office-hours, plan-ceo-review, plan-eng-review, plan-design-review, autoplan
  - `"review"` — ATV: ce-review | gstack: review, design-review, design-shotgun
  - `"qa-testing"` — gstack: qa, qa-only, benchmark, browse
  - `"security"` — ATV: security-sentinel agent | gstack: cso
  - `"shipping"` — ATV: ce-work | gstack: ship, land-and-deploy, canary, document-release
  - `"safety"` — gstack: careful, freeze, guard, unfreeze
  - `"debugging"` — gstack: investigate
  - `"retrospective"` — gstack: retro
- ATV skills are default-selected; gstack skills fill gaps
- Mark gstack browser skills as requiring "full runtime"

### wizard.go changes
- Replace source-based grouping with function-based categories
- Each category shows combined ATV + gstack skills
- Add a "Runtime" toggle: "Include browser-based skills (requires Bun)"
- Skills requiring full runtime are grayed out if runtime toggle is off

### scaffold.go changes
- New function: `cloneAndBuildGstack()` — shallow clone, strip `.git`, optionally `bun run build`
- Detect Bun availability before offering full runtime
- Graceful fallback: if Bun not found, warn and install markdown-only skills
- Network error handling: warn but don't fail

### copilot-instructions.md changes
- Single unified skill listing, organized by function
- ATV skills listed first within each category

## Resolved Questions

- **Version pinning**: Shallow clone (`git clone --depth 1`) + strip `.git` directory after. No formal releases needed — we pin to `main` at install time.
- **Windows support**: Package gstack installation as part of the Go installer itself. The Go binary handles cross-platform concerns — no bash dependency for setup.
- **gstack binary build**: gstack's TypeScript binary is built at install time on the user's machine. Requires Bun (or Node on Windows) to be available. The Go installer shells out to `bun run build` (or `node` fallback) after cloning.
- **Tarball source**: Not using tarballs. Shallow git clone from `github.com/garrytan/gstack`, then remove `.git` directory to vendor the snapshot.

## Next Steps

→ `/ce-plan` for implementation details
