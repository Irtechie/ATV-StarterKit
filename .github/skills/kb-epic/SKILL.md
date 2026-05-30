---
name: kb-epic
description: Large-initiative coordinator for KB workflows. Use for app migrations, framework rewrites, major architecture changes, multi-subsystem initiatives, multiple brainstorms/manifests, long backlogs, or dark-factory execution across related workstreams.
argument-hint: "[initiative description or epic path]"
---

# KB Epic

Coordinate large work without turning it into one huge plan. This is the lane for "migrate this app to a different framework", "rewrite the architecture", or "run a whole program of related work."

## Goal

Create an epic map that points to multiple brainstorms, plans, manifests, handoffs, research notes, and subsystem docs. Keep each execution unit small enough for `kb-plan` and `kb-work`.

## When To Use

Use `kb-epic` when:

- The request is bigger than one brainstorm or one manifest.
- Example: language/runtime migration, app-shell rewrite, auth overhaul, full interaction architecture replacement.
- One plan would be too large to execute or review.
- Work spans multiple subsystems.
- Architecture direction affects many later slices.
- Several brainstorms should feed one release.
- The user wants a dark-factory queue.

## Epic Location

Create:

```text
docs/context/epics/<epic-name>.md
```

If `docs/context/epics/` does not exist, create it and link it from `docs/context/PROJECT.md`.

## Epic Template

```markdown
# <Epic Name>

Status: draft|active|parked|complete
Created: YYYY-MM-DD
Last refreshed: YYYY-MM-DD

## Intent

## Success Criteria

## Architecture Decisions

## Research

## Workstreams

| Workstream | Brainstorm | Manifest | Status | Notes |
|---|---|---|---|---|

## Dependency Map

## Dark Factory Queue

## Human Checkpoints

## Parked / Blocked

## Completion Criteria
```

## Scheduling

Default to serial execution when workstreams share files, schemas, prompts, auth, routing, generated artifacts, or architecture decisions.

Use read-only swarms for research/review. Use coding swarms only when file ownership is declared and non-overlapping.

## Workstream Routing

An epic may have one umbrella brainstorm, but execution should not depend on one
giant brainstorm. Use the smallest useful artifact per workstream:

| Workstream State | Route | Output |
|---|---|---|
| Goal, user value, constraints, or success criteria are unclear | `kb-brainstorm` for that workstream | `docs/brainstorms/<date>-<workstream>-requirements.md` |
| Behavior is clear enough to slice, but execution shape is not | `kb-plan` for that workstream | `docs/plans/<date>-000-kb-<workstream>-manifest.md` plus slice plans |
| A valid manifest already exists and is current | `kb-work <manifest>` | executed slices |
| Workstream is blocked by architecture/research decision | `kb-research` or a small decision note first | research/decision linked from epic |

Brainstorm granularity:

- Use one umbrella brainstorm when the epic has one coherent decision surface.
  Its slice candidates may seed the Workstreams table.
- Use one brainstorm per workstream when the workstream has its own product
  tradeoffs, unclear acceptance criteria, or human checkpoint.
- Multiple brainstorms are right when the workstreams could make different
  product or architecture decisions without invalidating each other.
- Do not run a brainstorm for a workstream whose behavior is already clear; send
  it straight to `kb-plan`.
- Before sending a brainstormed workstream to `kb-plan`, resolve or explicitly
  park any `Resolve Before Planning` items. If those items affect scope,
  acceptance criteria, architecture direction, data contracts, or human risk,
  planning is premature.

Manifest convention:

- A manifest is the `kb-plan` output file under `docs/plans/` whose frontmatter
  has `type: kb-manifest` and whose filename normally looks like
  `<date>-000-kb-<topic>-manifest.md`.
- Record that manifest path in the epic Workstreams table and in `todo.md`.
- Queue only runnable or deliberately parked manifests. Do not queue a raw
  brainstorm as runnable work.

Todo target:

- Default target is the active project root's `todo.md` loaded by `kb-map`.
- If the environment has a separate session/private todo store, do not use it as
  the durable queue unless repo instructions explicitly say so. The epic must be
  resumable from repo-local memory.

`todo.md` queue convention:

```markdown
| <Workstream name> | ⬜ pending | P0|P1|P2 | `docs/plans/<manifest>.md` |
```

Use `🔧 in_progress` only while an agent is actively executing that manifest.
Use `🔒 blocked` or `🛑 human-required` when the next action cannot run without a
specific decision. Keep completed manifests out of active rows; move summaries
to `todo-done.md`.

## Flow

1. Read `todo.md`, `docs/context/PROJECT.md`, and relevant subsystem docs.
2. Create or refresh the epic map.
3. Split the initiative into workstreams small enough for normal KB flow.
4. Fill the Workstreams table as you route:
   - `Brainstorm`: blank, `skipped-clear`, or the brainstorm path.
   - `Manifest`: blank until `kb-plan` creates it, then the manifest path.
   - `Status`: `draft`, `planned`, `queued`, `running`, `blocked`, `done`, or
     `parked`.
   - `Notes`: blocker, dependency, human checkpoint, or proof pointer.
5. Resolve or park any brainstorm `Resolve Before Planning` items before
   routing that workstream to `kb-plan`.
6. Route each workstream using the Workstream Routing table above.
7. Queue each created manifest in `todo.md` using this repo's queue convention.
8. Use `kb-work` for runnable manifests.
9. Use `kb-complete` after each manifest or at epic milestones.
10. Use `kb-map refresh` after durable architecture changes.
11. Use `kb-ship` when release readiness matters.

Refresh cold epic work older than 72 hours before execution.
