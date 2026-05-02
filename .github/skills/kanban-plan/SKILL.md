---
name: kanban-plan
description: "Break a brainstorm or feature into vertical-slice task plans with dependency DAG, verification strategy, and HITL flags. Alternative to ce-plan that enforces end-to-end slicing over horizontal phases. Use when the user says 'kanban plan', 'slice this', 'break into vertical slices', 'kanban', or wants independently-grabbable tasks."
argument-hint: "[brainstorm path, feature description, or PRD]"
---

# Kanban Plan — Vertical Slice Decomposition

<!-- Inspired by mattpocock/skills to-issues — credit: github.com/mattpocock/skills -->

Break work into independently-executable **vertical slices** (tracer bullets). Each slice cuts through ALL layers end-to-end — never horizontal phases.

## Interaction Method

Use the platform's blocking question tool when available (`ask_user` in Copilot CLI, `AskUserQuestion` in Claude Code). Ask one question at a time. Prefer concise single-select choices.

## Input

<input> #$ARGUMENTS </input>

**If input is empty:** Check `docs/brainstorms/` for recent brainstorm documents. If found, ask which one to use. Otherwise ask: "What feature or work would you like to decompose into vertical slices?"

**If input is a brainstorm path:** Read it thoroughly — this is the source of truth for WHAT to build. Carry forward all decisions, scope boundaries, and requirements.

**If input is a feature description:** Proceed directly to decomposition.

## Core Rules

### Vertical Slices Only

Each slice must deliver a narrow but COMPLETE path through every relevant layer (schema, service, API, UI, tests). A completed slice is demoable or verifiable on its own.

```
WRONG (horizontal phases):
  Task 1: Create database schema
  Task 2: Build service layer
  Task 3: Add API routes
  Task 4: Build frontend

RIGHT (vertical slices):
  Task 1: Award points on lesson completion + show on dashboard
  Task 2: Track streaks (builds on task 1)
  Task 3: Add level progression display
```

### Enabling Slices Are Acceptable (With Constraints)

Some work IS legitimately enabling infrastructure (migrations, auth plumbing, shared config). Allow enabling slices ONLY when:

- They unlock a named downstream slice
- They are the smallest viable prerequisite
- The slice names its immediate consumer(s)

### Every Slice Has a Verification Strategy

| Mode | When | Gate |
|------|------|------|
| `tdd` | Behavior changes, business logic | Failing test → implement → passes |
| `integration` | Wiring, cross-boundary, API contracts | Integration test proves path works |
| `verification-only` | Config, scaffolding, ops | Builds pass, no regression |
| `hitl` | UX taste, design judgment | Human confirms acceptable |

## Process

### 1. Understand the Source Material

Read the brainstorm/PRD/description. Extract:

- What behaviors need to exist
- What the user-visible outcomes are
- What constraints/dependencies exist
- What's explicitly out of scope

### 2. Draft Vertical Slices

Break the work into thin end-to-end slices. For each slice, determine:

- **Title** — short descriptive name
- **What it delivers** — end-to-end behavior description
- **Verification mode** — tdd / integration / verification-only / hitl
- **Blocked by** — which other slices must complete first (or "none")
- **HITL flag** — does this need human judgment? (most should be `false` if brainstorm was thorough)

### 3. Present and Quiz the User

Show the proposed breakdown as a numbered list. Ask:

- Does the granularity feel right? (too coarse / too fine)
- Are dependency relationships correct?
- Should any slices be merged or split?
- Are verification modes correct?
- Any HITL flags wrong?

Iterate until approved.

### 4. Generate Plan Files

Create a **manifest** and **individual slice plans**.

#### Manifest: `docs/plans/YYYY-MM-DD-000-kanban-<name>-manifest.md`

```yaml
---
type: kanban-manifest
kanban_id: kb-YYYY-MM-DD-<name>
brainstorm_path: docs/brainstorms/<source-file>.md
created: YYYY-MM-DD
status: active
slices:
  - id: slice-001
    title: "<title>"
    path: docs/plans/YYYY-MM-DD-001-<type>-<name>-plan.md
    blockers: []
    verification: tdd
    hitl: false
    status: pending
  - id: slice-002
    title: "<title>"
    path: docs/plans/YYYY-MM-DD-002-<type>-<name>-plan.md
    blockers: [slice-001]
    verification: tdd
    hitl: false
    status: pending
---

# Kanban: <Feature Name>

## Origin
Brainstorm: `<brainstorm_path>`

## Slice Overview
| # | Slice | Blocked By | Verification | HITL | Status |
|---|-------|-----------|--------------|------|--------|
| 1 | <title> | — | tdd | no | pending |
| 2 | <title> | slice-001 | tdd | no | pending |
| 3 | <title> | — | integration | no | pending |
```

#### Individual Slice Plans: `docs/plans/YYYY-MM-DD-NNN-<type>-<name>-plan.md`

Each slice plan uses standard ATV plan format with additional frontmatter:

```yaml
---
kanban_id: kb-YYYY-MM-DD-<name>
slice_id: slice-NNN
title: "<title>"
blockers: []
verification: tdd
hitl: false
status: pending
---
```

The plan body should include:
- What to build (end-to-end behavior, not layer-by-layer)
- Acceptance criteria (testable)
- Files likely involved
- Test scenarios (specific enough for TDD)
- Scope boundary (what this slice does NOT include)

### 5. Commit

```bash
git add docs/plans/
git commit -m "kanban-plan: decompose <feature> into N vertical slices"
```

## Integration with Other Skills

- **Input from:** `ce-brainstorm`, `deepen-brainstorm`
- **Deepening:** Run `deepen-plan` on individual slices (one at a time to keep context small)
- **Execution:** `kanban-lightsout` runs all slices in order, OR `ce-work` picks up one slice at a time
- **Verification:** Each slice uses `tdd` skill when verification mode is `tdd`
