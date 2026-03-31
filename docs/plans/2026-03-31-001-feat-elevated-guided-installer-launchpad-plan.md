---
title: "feat: Elevated Guided Installer and Launchpad Foundations"
type: feat
status: active
date: 2026-03-31
origin: docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md
---

# feat: Elevated Guided Installer and Launchpad Foundations

## Overview

Implement the next generation of the guided ATV installer as a phased program that starts with a richer, multi-stack, confidence-building `--guided` flow and ends with the foundations for a deterministic post-install memory launchpad.

Phase 1 focuses on the guided installer itself: multi-stack pack selection, clearer preset/customize semantics, better capability descriptions, and structured install telemetry. Phase 2 introduces the persistent install and memory state required for a launchpad. Phase 3 delivers a deterministic, dashboard-first launchpad. Phase 4 leaves space for an optional GitHub Copilot SDK memory concierge, but only after the local source of truth is solid.

This plan intentionally treats the installer brainstorm as the parent product direction and the launchpad brainstorm as a focused follow-on experience (see brainstorm: `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md`; see brainstorm: `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md`).

## Problem Statement

Today the repo still implements a mostly single-stack guided flow with coarse install progress and a static post-install summary. That is out of sync with the current product direction:

- guided install should support additive stack packs instead of one primary stack
- presets and customization should be more legible
- install progress should act like telemetry, not decoration
- post-install should culminate in a memory-aware launchpad, not a dead-end summary

The most important architectural gap is that current code still assumes a single stack all the way through detection, wizard results, catalog building, and stack-specific file generation. The second major gap is that there is no canonical install manifest or memory index, which means a future launchpad would otherwise be forced to infer too much from the filesystem after the fact.

## Proposed Solution

Adopt a phased implementation strategy with hard boundaries between guided-installer UX, structured install state, and launchpad UI.

### Key product decisions carried forward

- Keep one-click mode untouched for now (see brainstorm: `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md`)
- Make stack support multi-stack by default in guided mode, with all supported packs selected initially
- Keep the installer hybrid rather than rewriting it as one giant Bubble Tea app
- Treat install progress as structured telemetry with reasons, warnings, skips, and outcomes
- Make the launchpad a follow-on, memory-first dashboard with terminal secondary
- Do not require GitHub Copilot SDK for v1 launchpad usefulness (see brainstorm: `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md`)

## Technical Approach

### Architecture

#### Domain contracts first

Before polishing the guided flow, introduce explicit shared types for:

- `StackPack` — additive selection model for Rails / TypeScript / Python / General
- `InstallCapability` or equivalent — normalized capability/install-plan representation beyond raw layer strings
- `InstallEvent` — machine-readable progress/telemetry events
- `InstallManifest` — what the installer intended, attempted, skipped, failed, and completed
- `RepoState` / `MemoryIndex` — derived local facts used by a later launchpad
- `Recommendation` — deterministic next-action suggestions with rationale

These contracts are the foundation for all later phases and avoid baking singular-stack assumptions into the new UI.

#### Guided installer remains Go-first

Phase 1 should stay inside the current Go CLI architecture:

- `cmd/init.go` remains the orchestration spine
- `pkg/tui/wizard.go` owns guided selection flow
- `pkg/tui/presets.go` owns preset definitions and trade-offs
- `pkg/tui/categories.go` evolves into richer grouped capability metadata
- `pkg/tui/progress.go` evolves from coarse step status into structured telemetry rendering
- `pkg/scaffold/*`, `pkg/gstack/*`, and `pkg/agentbrowser/*` emit structured outcomes rather than forcing the UI to infer them from side effects

#### Launchpad is explicitly later-phase

The launchpad should not be smuggled into the guided-flow work as “just one more screen.” Phase 2 creates its required state model. Phase 3 can then present that state in a dashboard-first experience with a reopen path. The desktop-companion/Electron path remains optional and later.

### Implementation Phases

#### Phase 0: Contracts, decisions, and regression guardrails

**Goal:** Lock the domain model and avoid breaking the existing default path while guided mode evolves.

**Tasks**

- [x] Define `StackPack` semantics, including whether `General` is additive, fallback-only, or base-layer only
- [x] Decide whether zero selected stack packs is invalid or a supported “minimal” path
- [x] Decide guided rerun behavior: additive-only vs manifest-updated-without-removal vs destructive reconciliation
- [x] Add regression tests or golden tests for non-guided mode behavior so “untouched” has a measurable meaning
- [x] Create shared install-state types for manifest, events, and recommendations
- [x] Decide whether manifest storage is repo-local only (`.atv/`) or repo-visible (`.github/`), and version the schema

**Success criteria**

- [x] Multi-stack semantics are explicit and documented
- [x] Non-guided mode contract is testable
- [x] Later phases can build on typed state instead of ad hoc booleans/strings

#### Phase 1: Guided selection flow MVP

**Goal:** Ship the new guided mental model without yet building a full premium details-pane browser.

**Tasks**

- [x] Replace singular stack selection with additive stack-pack selection in `pkg/tui/wizard.go`
- [x] Update `pkg/detect/detect.go` so detection becomes recommendation metadata rather than a hidden decision
- [x] Redefine `WizardResult` to hold multiple stack packs and richer preset/customization state
- [x] Refactor `pkg/scaffold/catalog.go` / `BuildFilteredCatalog(...)` to support multi-pack composition deterministically
- [x] Fix current preset drift in `pkg/tui/presets.go`, especially Pro vs Full browser/runtime behavior
- [x] Improve preset previews so downgrade behavior and prerequisites are obvious
- [x] Upgrade grouped category descriptions in `pkg/tui/categories.go` so users can tell what each capability does before selection
- [x] Preserve progressive disclosure: preset-first path remains the happy path; customization is optional

**Explicitly excluded in this phase**

- full dual-pane capability browser with rich preview surface
- launchpad UI
- Copilot-powered assistant behavior

**Success criteria**

- [x] Guided mode supports additive multi-stack selection with deterministic outcomes regardless of toggle order
- [x] Preset and customization flows clearly explain what will be installed and what prerequisites are needed
- [x] Auto / non-guided mode remains behaviorally unchanged

#### Phase 2: Structured install telemetry and manifest foundation

**Goal:** Turn install progress into structured telemetry and persist the state needed for a later launchpad.

**Tasks**

- [x] Redesign `pkg/tui/progress.go` around `InstallEvent` or equivalent structured event messages
- [ ] Emit parent and substep events from scaffold, gstack, and agent-browser flows
- [ ] Distinguish success, warning, failure, user-skip, prereq-skip, already-installed, and dependency-skipped outcomes
- [ ] Capture durations, reasons, and actionable next steps for each major operation
- [x] Write an atomic, versioned install manifest after guided installs
- [x] Ensure the manifest records requested vs installed vs skipped vs failed, not just intent
- [x] Upgrade the completion summary in `pkg/output/printer.go` to show what was installed, what was skipped, and how to reopen the next experience later

**Explicitly excluded in this phase**

- Electron or desktop companion packaging
- model-driven recommendations

**Success criteria**

- [x] Guided install shows structured telemetry with meaningful warnings and failure reasons
- [x] Manifest writes are atomic, versioned, and survive partial success
- [x] Completion output provides a better handoff even before a full launchpad exists

#### Phase 3: Deterministic launchpad MVP

**Goal:** Deliver a useful, offline-friendly launchpad driven by local facts.

**Tasks**

- [ ] Implement a memory index that classifies durable repo memory, installed intelligence, and runtime health
- [ ] Add deterministic recommendation rules based on repo state, manifest state, and runtime checks
- [ ] Limit primary recommendations to a small set with clear rationale
- [ ] Define a reopen path, likely a new CLI command and/or guided completion affordance
- [ ] Deliver the first launchpad UI as a dashboard-first experience
- [ ] Keep the terminal as a secondary execution surface, not the home screen
- [ ] Ensure graceful fallback in headless/no-GUI environments

**Candidate information architecture**

- Memory library
- Installed intelligence
- Recommended next moves
- Optional activity/log or terminal pane

**Success criteria**

- [ ] Same `RepoState` always produces the same ordered top recommendations
- [ ] Launchpad remains useful without network or Copilot auth
- [ ] Users can reopen the launchpad after install without rerunning the whole setup flow

#### Phase 4: Optional GitHub Copilot SDK concierge

**Goal:** Add an intelligent assistant on top of the deterministic local model without making it the source of truth.

**Tasks**

- [ ] Define typed tools for the assistant (e.g. `getMemorySummary`, `listRecommendations`, `explainRecommendation`, `openArtifact`, `runSuggestedAction`)
- [ ] Expose the local memory index and recommendations through those tools instead of raw file scraping
- [ ] Add clear degraded behavior for no-auth, offline, or slow-response cases
- [ ] Keep the assistant in an explanation/navigation role rather than a truth-owning role
- [ ] Ensure assistant recommendations cannot silently override deterministic ranking without explanation

**Success criteria**

- [ ] Core launchpad value is unchanged when the assistant is disabled
- [ ] The assistant improves explanation and navigation without becoming mandatory

## Alternative Approaches Considered

### 1. Rewrite guided mode as one continuous full-screen Bubble Tea control center

Rejected for the first wave because it adds too much surface-area churn before the underlying multi-stack and install-state contracts are stable. The March 31 brainstorm explicitly recommends elevating the current hybrid model instead of replacing it whole-cloth (see brainstorm: `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md`).

### 2. Build the launchpad desktop companion first

Rejected because the repo currently lacks the install manifest and local memory index that make a launchpad trustworthy. Building the dashboard before the state model would force it to guess, which is a fancy way to ship confusion.

### 3. Make GitHub Copilot SDK mandatory in launchpad v1

Rejected because it adds auth/network/runtime dependencies before the local product is solid. The follow-on brainstorm explicitly recommends deterministic recommendations first and an optional Copilot layer later (see brainstorm: `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md`).

## System-Wide Impact

### Interaction Graph

- Guided branch starts in `cmd/init.go:runInit`
- Guided selection flows through `pkg/tui/wizard.go:RunWizard`
- Current stack/preset selections flow into catalog generation via `pkg/scaffold/catalog.go:BuildCatalog` / `BuildFilteredCatalog`
- Install execution calls scaffold writes plus optional gstack and agent-browser installers
- Completion currently terminates in `pkg/output/printer.go:PrintNextSteps`

Under this plan, the chain becomes:

`guided selection -> normalized install plan -> structured install events -> manifest write -> completion handoff -> later launchpad read/index/recommend`

### Error & Failure Propagation

Current behavior risks flattening partial failures into generic success. The plan should make failures explicit at the event and manifest layer:

- scaffold write failures should not disappear inside a “scaffolding completed” parent step
- prereq-related skips should not look like successful installs
- optional runtime steps must report downgrade reasons clearly
- manifest writing must not silently succeed when the install itself partially failed

### State Lifecycle Risks

Biggest risks:

- stale or corrupted manifest after partial install
- mismatched requested vs actual installed capabilities
- reruns that change intent without cleaning up prior assets
- false inference of runtime capability from repo-local files alone

Mitigations:

- atomic versioned manifest writes
- explicit requested/installed/skipped/failed fields
- rerun semantics defined in Phase 0
- separate installer intent from derived runtime health

### API Surface Parity

Equivalent surfaces that must stay aligned:

- guided UI copy and preset semantics
- catalog filtering logic and file composition rules
- completion summary and any future reopen command
- launchpad recommendations and CLI fallback next-step recommendations

### Integration Test Scenarios

- [x] Guided multi-stack selection produces deterministic file output regardless of selection order
- [ ] Missing prerequisites downgrade or skip runtime-dependent capabilities with clear telemetry and manifest output
- [ ] Partial install success still writes a valid manifest and useful completion summary
- [ ] Re-running guided mode after a prior install preserves or updates state according to the chosen rerun model
- [ ] Launchpad recommendation output matches manifest + repo memory for empty, partial, and mature repos

## Acceptance Criteria

### Functional Requirements

- [x] `init --guided` supports additive stack-pack selection and no longer forces one primary stack
- [x] Detection acts as a recommendation/preselection layer rather than silently choosing for the user
- [x] Preset previews clearly state capability scope, install/runtime expectations, and downgrade behavior
- [x] Customization presents grouped, understandable capability choices
- [x] Install progress includes first-class warnings, skips, failures, and reasons
- [x] Guided install writes a versioned manifest describing requested, installed, skipped, and failed outcomes
- [ ] A deterministic launchpad MVP can read repo memory and manifest state to recommend next actions
- [ ] GitHub Copilot SDK remains optional and non-blocking

### Non-Functional Requirements

- [x] Multi-stack composition is deterministic for the same set of selected packs
- [x] Manifest writes are atomic and recover cleanly from partial install failures
- [x] Non-guided mode remains behaviorally unchanged during Phase 1 and Phase 2
- [ ] Launchpad recommendations are deterministic and offline-friendly
- [ ] Headless/no-GUI environments degrade gracefully to a strong terminal summary and reopen instructions

### Quality Gates

- [x] Add test coverage for guided flow state transitions or extracted planning logic
- [x] Add regression coverage for non-guided mode behavior
- [ ] Document manifest schema and recommendation heuristics
- [x] Resolve preset drift between docs and code before shipping the new guided presets publicly

## Success Metrics

- Users can complete guided setup with a clearer understanding of what ATV installed and why
- Users can see and trust skipped/downgraded runtime steps rather than discovering them later
- The repo gains a reusable install-state foundation for launchpad and future recommendations
- Launchpad v1 provides useful next steps without requiring Copilot auth or model access

## Dependencies & Prerequisites

### Internal dependencies

- Multi-stack-compatible scaffold composition
- Structured install result/event model
- Manifest storage and schema versioning
- Better preset/runtime boundary definitions

### External/runtime dependencies

- Existing Go CLI stack for guided installer work
- Bun / Node / npm / browser runtime prerequisites where optional tools depend on them
- Any future Electron/xterm/node-pty packaging work only after Phase 2 and Phase 3 state foundations exist

## Risk Analysis & Mitigation

| Risk | Severity | Mitigation |
| --- | --- | --- |
| Single-stack architecture resists additive packs | High | Treat multi-stack as a domain-model change in Phase 0/1, not a UI-only task |
| Preset semantics drift from actual runtime behavior | High | Fix Pro/Full boundary early and test downgrade paths |
| Telemetry remains cosmetic | High | Require structured event payloads from installers, not just status strings |
| Launchpad scope balloons into desktop-app work too early | High | Keep manifest/index and deterministic recommendations separate from Electron work |
| Runtime capability inference is inaccurate | Medium | Store explicit verification state and separate repo facts from machine-local health |
| “Untouched” non-guided promise regresses silently | Medium | Add regression tests before broad refactors |

## Future Considerations

- Once the deterministic launchpad exists, consider whether a desktop companion is actually needed or whether a strong terminal/TUI handoff is sufficient
- If Electron is pursued, keep the assistant and terminal layers secondary to the local memory dashboard
- Over time, the recommendation engine can evolve from simple heuristics to richer workflow understanding without changing its local source-of-truth contract

## Documentation Plan

- Update README guided-install sections to match actual multi-stack/preset behavior
- Document install manifest schema and storage location
- Document launchpad reopen workflow once available
- Add a `docs/solutions/` writeup once the first phase lands, since this repo currently has no historical solution doc for this topic

## Sources & References

### Related planning follow-up

- `docs/plans/2026-03-31-002-refactor-elevated-installer-launchpad-readiness-plan.md` — readiness review that confirms the roadmap is implementable, while recording the narrow contract decisions and first coding slice that should happen before UI-heavy work.

### Origin

- **Primary brainstorm:** `docs/brainstorms/2026-03-31-elevated-guided-installer-experience-brainstorm.md` — key decisions carried forward: hybrid guided architecture, additive stack packs, telemetry-not-decoration, launchpad ending
- **Supporting brainstorm:** `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md` — key decisions carried forward: deterministic launchpad first, terminal secondary, optional Copilot SDK later

### Internal References

- Current guided installer spine: `cmd/init.go`
- Current guided wizard: `pkg/tui/wizard.go`
- Current presets: `pkg/tui/presets.go`
- Current grouped capability metadata: `pkg/tui/categories.go`
- Current progress model: `pkg/tui/progress.go`
- Current catalog composition: `pkg/scaffold/catalog.go`
- Current scaffold writer/results: `pkg/scaffold/scaffold.go`
- Current post-install summary: `pkg/output/printer.go`
- Launchpad architecture spike: `docs/spikes/architecture-post-install-memory-launchpad-spike.md`
- Earlier rich-TUI plan for historical context: `docs/plans/2026-03-29-002-feat-rich-tui-experience-plan.md`

### Research Findings Incorporated

- No prior relevant `docs/solutions/` exist yet for this feature area
- Repo research confirms guided work should start in Go CLI/TUI surfaces, while launchpad remains a later net-new surface
- Spec-flow review highlighted unresolved semantics around `General`, reruns, manifest ownership, and launchpad handoff; these are now Phase 0 decisions instead of hidden assumptions
