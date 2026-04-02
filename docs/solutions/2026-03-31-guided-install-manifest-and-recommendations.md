---
title: "solution: Guided install manifest and deterministic recommendations"
type: solution
status: completed
date: 2026-03-31
origin: docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md
---

# solution: Guided install manifest and deterministic recommendations

## What landed

The guided installer now records a repo-local manifest at `./.atv/install-manifest.json` and fills it with:

- the user's requested stack packs and layers
- coarse structured outcomes for each guided install step
- deterministic next-step recommendations derived from local repo state and install outcomes

This turns the guided flow from a write-only experience into a stateful foundation that later launchpad work can trust.

## Manifest schema

Current manifest fields:

- `version`
  - integer schema version
- `generatedAt`
  - UTC timestamp of the guided install run
- `rerunPolicy`
  - currently `additive-only`
- `requested`
  - `stackPacks`
  - `atvLayers`
  - `gstackDirs`
  - `gstackRuntime`
  - `includeAgentBrowser`
  - `presetName`
- `outcomes`
  - one entry per guided install step
  - includes `step`, `status`, optional `detail`, optional `reason`, and `duration`
- `recommendations`
  - small deterministic list of next moves with `id`, `title`, `reason`, and `priority`

### Current outcome statuses

- `done`
- `warning`
- `failed`
- `skipped`

These statuses are intentionally coarse for now. They are enough to distinguish healthy completion from degraded or partial install paths while leaving room for richer sub-statuses later.

## Recommendation heuristics

The current recommendation engine is local and deterministic.

It combines:

1. **Install outcomes**
   - installer warnings/failures surface a top recommendation to fix prerequisites or degraded setup first
2. **Repo memory state**
   - no brainstorms → recommend `/ce-brainstorm`
   - brainstorms but no plans → recommend `/ce-plan`
   - plan files with unchecked boxes → recommend `/ce-work`
   - completed plans but no solution docs → recommend `/ce-compound`
3. **Installed capabilities**
   - usable gstack install → recommend `/gstack-office-hours`
   - usable agent-browser install → recommend trying browser automation next

### Guardrails

- recommendations are sorted deterministically by priority, then ID
- only the top 3 recommendations are surfaced as primary next moves
- each recommendation must carry a concrete reason so the system can explain *why* it appeared

## Why this matters

This solves two immediate product problems:

1. guided install no longer ends in a dead-end summary
2. later launchpad work no longer has to infer installer intent from scattered filesystem side effects alone

The manifest is not the full memory system. It is the installer's source of truth about what it tried to do and what happened.

## What is still intentionally missing

- a full memory index over repo artifacts
- recommendation rules for review parity and deeper runtime health
- a dedicated reopen command or dashboard UI
- any GitHub Copilot SDK dependency

Those remain follow-on work. The current implementation keeps the core logic trustworthy, offline-friendly, and easy to extend.