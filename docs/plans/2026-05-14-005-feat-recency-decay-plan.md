---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-005
title: "Add recency decay to /learn"
blockers: []
verification: verification-only
hitl: false
expected_files:
  - .github/skills/learn/SKILL.md
status: pending
---

# Slice 5: Add Recency Decay to /learn

## What to Build

Add a decay mechanic to `/learn` so that instincts which haven't been observed recently lose confidence over time. This prevents stale patterns from persisting forever. Uses a half-life formula inspired by CrewAI's memory architecture.

## Acceptance Criteria

- `/learn` applies decay before creating/updating instincts (Step 2.5 or similar)
- Formula: `decayed_confidence = confidence × 0.5^(days_since_last_seen / half_life)`
- Default half-life: 90 days (project-level patterns have long shelf life)
- Instincts that decay below 0.3 are archived (moved to `.atv/instincts/archive/`)
- Decay is applied using `last_seen` field (already exists in instinct format)
- Archival includes a note: `archived_reason: confidence decayed below 0.3`
- No change to how new instincts are created or how confidence is increased

## Implementation

Add a new step between existing Step 2 (Analyze Patterns) and Step 3 (Create or Update Instincts):

```markdown
### Step 2.5: Apply Recency Decay

Before creating or updating instincts, apply time-based decay to all existing entries:

1. For each instinct in `.atv/instincts/project.yaml`:
   - Calculate days since `last_seen`
   - Apply decay: `new_confidence = confidence × 0.5^(days / 90)`
   - Update the confidence value in place

2. **Archive stale instincts:**
   - If decayed confidence falls below 0.3, move the instinct to `.atv/instincts/archive/YYYY-MM-DD-archived.yaml`
   - Add `archived_reason: confidence decayed below 0.3 (last seen: <date>)`
   - Remove from `project.yaml`

3. Write updated confidence values back to `project.yaml` before proceeding to Step 3.

**Half-life rationale:** 90 days balances stability (project conventions rarely change weekly) with freshness (patterns unused for 6+ months are likely obsolete). At 90 days, an unobserved instinct at 0.85 decays to:
- 30 days: 0.68
- 90 days: 0.43
- 180 days: 0.21 (archived)
```

## Scope Boundary

- Do NOT modify how confidence is increased on observation (still +0.1)
- Do NOT modify the instinct creation flow
- Do NOT change the instinct YAML schema (uses existing `last_seen` field)
- Do NOT modify `/evolve` — that's Slice 6
- Single file edit

## Test Scenarios

- Decay formula produces correct values (0.85 × 0.5^(90/90) = 0.425)
- Archive threshold is 0.3 (not lower — we want aggressive cleanup)
- Decay runs before updates (so a re-observed instinct doesn't get decayed then boosted in same pass)
- The half-life value (90 days) is stated explicitly, not left as a magic number
