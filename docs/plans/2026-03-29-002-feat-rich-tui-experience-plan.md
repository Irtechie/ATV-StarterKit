---
title: "feat: Rich TUI Experience with Presets and Install Progress"
type: feat
status: active
date: 2026-03-29
origin: docs/brainstorms/2026-03-29-rich-tui-experience-brainstorm.md
---

# feat: Rich TUI Experience with Presets and Install Progress

## Overview

Redesign the `--guided` installer TUI with a 3-step wizard (Stack → Preset → Customize) and animated install progress. Uses `charmbracelet/huh` for wizard forms and raw `charmbracelet/bubbletea` + `charmbracelet/bubbles` for the install progress phase.

## Problem Statement

The current TUI dumps 43+ skills into a flat checkbox list with terse labels. Users don't know what they're selecting, there's no beginner path, and the install phase prints text line-by-line with no progress feedback.

## Proposed Solution

Three changes:
1. **Presets** — Starter / Pro / Full with a "Customize?" confirm step
2. **Category descriptions** — each group gets a 1-line explanation
3. **Install progress** — Bubbletea-powered step-by-step status with spinners

## Implementation Phases

### Phase 1: Presets — `pkg/tui/presets.go`

- [ ] Define `Preset` struct: `Name`, `Description`, `ATVLayers []string`, `GstackDirs []string`, `IncludeAgentBrowser bool`
- [ ] Define three presets:
  - **Starter**: Core ATV skills + agents + infra. No gstack, no agent-browser. ~13 skills, instant install
  - **Pro**: Starter + gstack text-only skills (review, ship, careful, investigate, retro, cso, office-hours, plan-*). No browser QA. ~35 skills
  - **Full**: Everything. All gstack skills including browser QA, agent-browser + Chrome. ~45 skills
- [ ] `PresetSkillKeys(preset Preset, prereqs Prerequisites) []string` — returns the skill keys for a preset
- [ ] `PresetDescription(preset Preset) string` — returns a rich description with skill count and install time estimate

### Phase 2: Wizard Redesign — `pkg/tui/wizard.go`

- [ ] **Screen 1: Stack** (unchanged) — Select/confirm detected stack
- [ ] **Screen 2: Preset** — New `huh.Select` with 3 presets:
  ```
  ┃ Choose your setup level
  ┃
  ┃ > ⚡ Starter — Core workflow (13 skills, instant)
  ┃     Plan, build, review, compound. No browser tools.
  ┃
  ┃   🚀 Pro — Full sprint process (35 skills)
  ┃     + gstack review, ship, safety, security, debugging
  ┃
  ┃   🔥 Full — Complete engineering team (45 skills)
  ┃     + browser QA, benchmarks, agent-browser, Chrome
  ┃     Requires: Bun, ~2min install
  ```
- [ ] **Screen 3: Customize?** — `huh.Confirm`: "Want to customize individual skills?"
  - If no → proceed to install with preset selections
  - If yes → show category-grouped multi-select (existing logic, pre-filled with preset's selections)
- [ ] **Screen 4 (optional): Customize** — Existing multi-select but with category headers and descriptions:
  ```
  ┃ 📋 Planning & Design
  ┃   Brainstorm ideas, create plans, research approaches
  ┃   [•] Brainstorming — explore what to build
  ┃   [•] Plan — structured implementation plans
  ┃   [•] [gstack] Office Hours — YC-style forcing questions
  ┃
  ┃ 🔍 Code Review
  ┃   Multi-agent review, security audits, design checks
  ┃   [•] CE Review — multi-agent code review
  ┃   [•] [gstack] Review — staff-level PR review
  ```
- [ ] Update `WizardResult` to include the selected preset name (for progress display)

### Phase 3: Install Progress — `pkg/tui/progress.go`

- [ ] Add `charmbracelet/bubbles` to go.mod (spinners, progress)
- [ ] Define `InstallStep` struct: `Name string`, `Status StepStatus`, `Detail string`
- [ ] Define `StepStatus` enum: `StepPending`, `StepRunning`, `StepDone`, `StepFailed`, `StepSkipped`
- [ ] Create Bubbletea `progressModel` implementing `tea.Model`:
  - `Init()` — start first step
  - `Update()` — handle step completion messages, advance to next step
  - `View()` — render all steps with status indicators:
    ```
    Installing Pro preset for TypeScript...

    ✅ ATV files scaffolded (38 files, 6 dirs)
    ⣾  Cloning gstack...
    ○  Generating gstack skill docs
    ○  Copying skills to .github/skills/
    ○  Installing agent-browser
    ○  Downloading Chrome for Testing
    ```
- [ ] Step status icons:
  - `○` pending (dim)
  - `⣾⣽⣻⢿⡿⣟⣯⣷` spinner (animated, bright)
  - `✅` done (green)
  - `❌` failed (red)
  - `⏭️` skipped (dim)
- [ ] Each step runs as a `tea.Cmd` that returns a `stepCompleteMsg` on finish
- [ ] Wire into `cmd/init.go`: after wizard completes, run `tea.NewProgram(progressModel)` instead of the current sequential print calls

### Phase 4: Wire It All Together — `cmd/init.go`

- [ ] Update `runInit()`:
  - Wizard returns `WizardResult` with preset + optional customizations
  - Build the install plan from result (which steps to run, which to skip)
  - Launch Bubbletea progress program
  - Progress program runs each install step sequentially, updating the view
  - After completion, print the "next steps" summary (outside Bubbletea)
- [ ] Preserve auto mode (`init` without `--guided`) as-is — no TUI, just print

## Acceptance Criteria

### Functional

- [ ] `atv-installer init --guided` shows Stack → Preset → Customize flow
- [ ] Starter preset installs only ATV skills (no network calls)
- [ ] Pro preset installs ATV + gstack text skills
- [ ] Full preset installs everything including agent-browser + Chrome
- [ ] Customize step shows category-grouped skills with descriptions
- [ ] Install phase shows animated spinner for each running step
- [ ] Failed steps show error but don't block subsequent independent steps
- [ ] Auto mode unchanged

### Non-Functional

- [ ] Install progress renders smoothly without flicker
- [ ] Spinner animation runs at ~10fps
- [ ] Each step status updates immediately on completion

## Dependencies

| Dependency | Purpose | Status |
|-----------|---------|--------|
| `charmbracelet/huh` | Wizard forms | ✅ Already in go.mod |
| `charmbracelet/bubbles` | Spinners for progress | ❌ Need to add |
| `charmbracelet/bubbletea` | Progress program runtime | ✅ Already transitive dep |

## Sources & References

### Origin

- **Brainstorm:** [docs/brainstorms/2026-03-29-rich-tui-experience-brainstorm.md](docs/brainstorms/2026-03-29-rich-tui-experience-brainstorm.md) — presets, hybrid huh+bubbletea, step-by-step progress

### Internal References

- Current wizard: `pkg/tui/wizard.go`
- Current categories: `pkg/tui/categories.go`
- Current printer: `pkg/output/printer.go`
- Install flow: `cmd/init.go`
