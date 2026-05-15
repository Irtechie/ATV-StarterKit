---
type: kanban-manifest
kanban_id: kb-2026-05-14-compound-loop
brainstorm_path: docs/brainstorms/2026-05-14-kanban-compound-loop-requirements.md
created: 2026-05-14
status: active
slices:
  - id: slice-001
    title: "Wire ce-review as mandatory async"
    path: docs/plans/2026-05-14-001-feat-wire-review-plan.md
    blockers: []
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-002
    title: "P0/P1 resolution gate + compound step"
    path: docs/plans/2026-05-14-002-feat-compound-step-plan.md
    blockers: [slice-001]
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-003
    title: "Feed resolved findings into observations"
    path: docs/plans/2026-05-14-003-feat-observations-feed-plan.md
    blockers: [slice-002]
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-004
    title: "Auto-trigger /learn after compound"
    path: docs/plans/2026-05-14-004-feat-auto-learn-plan.md
    blockers: [slice-003]
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-005
    title: "Add recency decay to /learn"
    path: docs/plans/2026-05-14-005-feat-recency-decay-plan.md
    blockers: []
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-006
    title: "Staleness guard in /evolve"
    path: docs/plans/2026-05-14-006-feat-staleness-guard-plan.md
    blockers: [slice-005]
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
  - id: slice-007
    title: "Auto-evolve on cadence + completion counter"
    path: docs/plans/2026-05-14-007-feat-auto-evolve-plan.md
    blockers: [slice-004, slice-006]
    verification: verification-only
    hitl: false
    status: pending
    notes: ""
---

# Kanban: Compound Engineering Loop

## Origin
Brainstorm: `docs/brainstorms/2026-05-14-kanban-compound-loop-requirements.md`

## Slice Overview
| # | Slice | Blocked By | Verification | HITL | Status |
|---|-------|------------|--------------|------|--------|
| 1 | Wire ce-review as mandatory async | - | verification-only | no | pending |
| 2 | P0/P1 resolution gate + compound step | slice-001 | verification-only | no | pending |
| 3 | Feed resolved findings into observations | slice-002 | verification-only | no | pending |
| 4 | Auto-trigger /learn after compound | slice-003 | verification-only | no | pending |
| 5 | Add recency decay to /learn | - | verification-only | no | pending |
| 6 | Staleness guard in /evolve | slice-005 | verification-only | no | pending |
| 7 | Auto-evolve on cadence + completion counter | slice-004, slice-006 | verification-only | no | pending |

## Dependency DAG

```
Track A (kanban-work flow):      Track B (instinct quality):
  slice-001                        slice-005
      ↓                                ↓
  slice-002                        slice-006
      ↓                                ↓
  slice-003                            │
      ↓                                │
  slice-004 ───────────────────────────┘
      ↓
  slice-007  (converges both tracks)
```

## Existing Infrastructure (already shipped)

- **Step 3.6: Diff-Scope Verification** — hard gate enforcing `expected_files` per slice
- **`expected_files` field in kanban-plan** — slices declare which files they'll touch
- **ce-review** — 14+ agent review with P0-P3 confidence-gated output
- **ce-compound** — 6 subagents, writes YAML-frontmatted docs to `docs/solutions/`
- **/learn** — reads observations.jsonl + git history → instincts
- **/evolve** — promotes instincts at confidence > 0.8, observations > 5
