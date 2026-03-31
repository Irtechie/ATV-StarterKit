---
title: "refactor: Validate elevated installer launchpad readiness"
type: refactor
status: completed
date: 2026-03-31
origin: docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md
---

# refactor: Validate elevated installer launchpad readiness

## Overview

Assess whether `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md` is ready for implementation as-is or whether the team should take another planning pass before writing code.

Verdict: **the roadmap is ready enough to implement, but only if implementation starts with a narrow Phase 0 kickoff slice that locks a handful of still-implicit contracts.** The team does **not** need another full brainstorm or another full replacement plan. It does need a short readiness-hardening pass so engineers are not forced to make product-shaping decisions mid-refactor.

## Problem Statement / Motivation

The current roadmap is strong at the product and phase level, but the underlying codebase is still built on a single-stack model:

- `pkg/detect/detect.go` returns one `Environment.Stack`
- `pkg/tui/wizard.go` stores one `WizardResult.Stack`
- `pkg/scaffold/catalog.go` composes files from one `stack`
- `cmd/init.go` orchestrates a write-only guided flow with coarse progress and no manifest

That means the plan is **not** ready for engineers to jump directly into multi-stack UI work or launchpad work. However, it **is** ready for a contracts-first implementation kickoff, because the missing pieces are now clear and bounded.

## Proposed Solution

Treat the existing roadmap as the approved umbrella plan and add a readiness gate in front of execution:

1. **Approve the roadmap as the implementation backbone**
2. **Lock five contract decisions up front**
3. **Start coding Phase 0 immediately after those defaults are recorded**
4. **Defer launchpad desktop-companion questions until after manifest/index foundations exist**

This means the right answer is **not** “take another crack at the whole plan.” The right answer is: **start implementation, but start at the contract layer, not the shiny layer.**

## Technical Considerations

- The current product docs are aligned: both brainstorms say the direction is settled and the launchpad should be deterministic-first.
- The current code is not aligned: stack detection, wizard state, scaffold composition, and preset/runtime behavior still assume one primary stack and coarse install outcomes.
- There are no prior `docs/solutions/` writeups to suggest a hidden gotcha beyond what the repo code and March 31 docs already surface.
- External research is not needed for this readiness decision; local repo context is strong enough.

## Readiness Verdict

### Ready to implement now

These areas are implementation-ready:

- **Phase 0 contract and guardrail work** from `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md`
- **Non-guided regression protections**
- **Preset drift fixes** where current code contradicts current product intent
- **Typed install-state/model extraction** for future manifest and recommendation work

### Not ready to implement blindly

These areas should **not** start until the defaults below are recorded and accepted:

- multi-stack selection UI rewrite in `pkg/tui/wizard.go`
- multi-pack composition changes in `pkg/scaffold/catalog.go`
- manifest-backed launchpad behavior
- Electron / PTY / desktop-companion work
- GitHub Copilot SDK assistant behavior

### Bottom-line recommendation

**Proceed with implementation. Do not restart planning.**

But constrain the first coding slice to a readiness-hardening Phase 0. If the team starts with direct UI work, they will hit ambiguity potholes immediately.

## Must-Fix Planning Gaps

These are the only planning gaps significant enough to block clean implementation. Recommended defaults are included so the team can move without another large planning loop.

### 1. `General` pack semantics

**Gap:** The roadmap identifies `General` as part of the new stack-pack model, but current code treats it as a single-stack fallback.

**Recommended default:**
- `General` becomes the **shared/base pack** for non-stack-specific assets
- it remains selectable in guided mode
- it contributes only stack-agnostic assets
- when no stack-specific signals are present, it is the default selected pack
- zero selected packs is invalid

**Why this is enough:** It preserves the brainstorm intent while giving engineers a deterministic composition rule.

### 2. Guided rerun behavior

**Gap:** The roadmap names rerun semantics as a Phase 0 decision, but current scaffold behavior is already mostly additive/skip-based.

**Recommended default:**
- v1 reruns are **additive-only**
- deselection updates intent in the manifest but does **not** delete previously installed files automatically
- cleanup/removal can be a later explicit command if needed

**Why this is enough:** It matches the current safe write behavior and avoids destructive surprises.

### 3. Manifest ownership and location

**Gap:** The roadmap leaves `.atv/` vs `.github/` open.

**Recommended default:**
- write installer state to **`.atv/install-manifest.json`**
- keep it versioned and local to the repo directory
- do not treat it as the durable memory source of truth; it is installer state, not product memory

**Why this is enough:** It separates installer intent from committed repo artifacts and avoids muddying `.github/` with machine-local status.

### 4. Launchpad v1 scope

**Gap:** The roadmap is clear on phasing, but engineers could still over-interpret “launchpad” as “desktop companion now.”

**Recommended default:**
- launchpad v1 means **deterministic local state + reopen command + strong terminal/TUI handoff**
- desktop companion / Electron remains an explicit later option

**Why this is enough:** It keeps launchpad implementation grounded and avoids dragging Windows PTY/Electron risk into the first shipping phases.

### 5. “Untouched” non-guided mode

**Gap:** The roadmap says one-click stays untouched, but that needs a testable definition.

**Recommended default:**
- no new prompts in non-guided mode
- no manifest writing in non-guided mode initially
- same high-level output semantics as current behavior
- regression tests prove this before guided refactors land

**Why this is enough:** It creates a measurable boundary instead of a slogan.

## Best First Implementation Slice

The best first coding slice is **not** the new wizard UI. It is a small contract-and-guardrail slice:

### Slice A — contract layer

Create a small package or module for shared install-state types, for example:

- `pkg/installstate/types.go`
- `pkg/installstate/manifest.go`

Initial types:

- `StackPack`
- `InstallManifest`
- `InstallOutcome`
- `InstallStepStatus`

### Slice B — guardrails

Add regression coverage for non-guided mode and current preset boundaries, likely around:

- `cmd/init.go`
- `pkg/tui/presets.go`
- extracted pure functions where possible

### Slice C — current drift fix

Fix the existing mismatch where:

- `pkg/tui/presets.go` documents Pro as text-only / no browser QA
- but `ProPreset.IncludeAgentBrowser` is currently `true`
- and `pkg/tui/wizard.go` implicitly turns on gstack runtime whenever Bun exists and gstack dirs are selected

This slice would make the current product truer before any multi-stack redesign begins.

## System-Wide Impact

### Interaction graph

Current chain:

`cmd/init.go -> detect.DetectEnvironment -> tui.RunWizard -> scaffold.BuildFilteredCatalog -> scaffold.WriteAll / gstack.Install / agentbrowser.Install -> output.PrintNextSteps`

Near-term implementation chain should become:

`contracts + tests -> guided state extraction -> multi-pack composition -> structured install events -> manifest write -> launchpad/reopen entry point`

### Error propagation

The roadmap is correct to insist that warnings, skips, and failures become first-class. Today `buildInstallSteps()` in `cmd/init.go` collapses scaffolding into a step that ignores write results, and gstack/agent-browser each appear as one coarse step. Engineers need typed outcomes before they build better telemetry.

### State lifecycle risks

The biggest danger is mixing three kinds of state into one fuzzy blob:

- installer intent
- repo memory
- runtime health

This readiness review agrees with the main roadmap: keep those separate from day one.

### API surface parity

Any implementation kickoff must keep these aligned:

- preset definitions in `pkg/tui/presets.go`
- guided wizard state in `pkg/tui/wizard.go`
- stack detection in `pkg/detect/detect.go`
- scaffold composition in `pkg/scaffold/catalog.go`
- post-install summary in `pkg/output/printer.go`

## Acceptance Criteria

- [ ] The team agrees on the five readiness defaults above without reopening the whole roadmap
- [x] A Phase 0 coding slice is identified and limited to contracts, regression guardrails, and current drift fixes
- [ ] Engineers do not need to guess `General` semantics, rerun behavior, manifest location, or launchpad v1 scope while coding
- [x] Non-guided mode has an explicit protected contract before guided-mode refactors begin
- [ ] The roadmap remains the umbrella source of truth, with this readiness plan acting as the kickoff clarifier

## Success Metrics

- Engineers can start coding without needing a new brainstorming round
- The first PRs land in low-regret surfaces: contracts, tests, and current behavior alignment
- The team avoids accidentally starting on Electron/Copilot work before the local state model exists
- The existing roadmap survives as the long-form implementation guide instead of being thrown away and rewritten

## Dependencies & Risks

### Dependencies

- existing roadmap: `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md`
- parent brainstorm: `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md`
- follow-on brainstorm: `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md`

### Risks

- If the team skips the readiness defaults, engineers will make inconsistent contract decisions in UI code
- If the team starts launchpad work before manifest/index contracts, they will build inference-heavy logic that needs rewriting
- If the team keeps Pro/Full drift unresolved, user-facing copy will remain out of sync with actual behavior during the transition

## Sources & References

### Origin

- **Brainstorm document:** `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md` — carried-forward decisions: hybrid guided flow, additive stack-pack direction, telemetry-not-decoration, launchpad handoff

### Internal References

- umbrella roadmap: `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md`
- launchpad brainstorm: `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md`
- launchpad spike: `docs/spikes/architecture-post-install-memory-launchpad-spike.md`
- current installer spine: `cmd/init.go`
- current guided wizard: `pkg/tui/wizard.go`
- current preset definitions: `pkg/tui/presets.go`
- current detection model: `pkg/detect/detect.go`
- current scaffold composition: `pkg/scaffold/catalog.go`

### Research Findings

- No relevant `docs/solutions/` history exists yet for this feature area
- Strong local context made external research unnecessary for this readiness check
- Current repo architecture confirms the roadmap is implementable, but only from a contracts-first starting point
