---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-006
title: "Staleness guard in /evolve"
blockers: [slice-005]
verification: verification-only
hitl: false
expected_files:
  - .github/skills/evolve/SKILL.md
status: pending
---

# Slice 6: Staleness Guard in /evolve

## What to Build

Tighten `/evolve`'s promotion criteria: raise the confidence threshold from 0.8 to 0.85, and add a staleness check — instincts not seen in 90+ days cannot be promoted regardless of their historical confidence.

## Acceptance Criteria

- Confidence threshold raised from 0.8 to 0.85
- New filter: `last_seen` must be within the last 90 days
- Both conditions must be met (AND, not OR)
- The "no candidates found" message reflects the updated criteria
- Observation count remains at > 5 (no change)
- Existing archive check remains (not already evolved)
- Same formatting and voice as the rest of the document

## Implementation

Modify Step 1 (Identify Candidates) in evolve/SKILL.md:

**Before:**
```markdown
Read `.atv/instincts/project.yaml` and filter for:
- Confidence > 0.8
- Observations > 5
- Not already evolved (check `.atv/instincts/archive/`)
```

**After:**
```markdown
Read `.atv/instincts/project.yaml` and filter for:
- Confidence > 0.85
- Observations > 5
- `last_seen` within the last 90 days (rejects stale instincts even if confidence is high)
- Not already evolved (check `.atv/instincts/archive/`)
```

Also update the "no candidates" message:

**Before:**
```
Instincts need confidence > 0.8 and 5+ observations to evolve.
```

**After:**
```
Instincts need confidence > 0.85, 5+ observations, and activity within 90 days to evolve.
```

## Scope Boundary

- Do NOT modify Step 2 (clustering) or Step 3+ (skill generation)
- Do NOT modify the archive structure
- Do NOT change observation count threshold
- Do NOT add recency decay here (that's /learn's job in Slice 5)
- Single file edit

## Test Scenarios

- Threshold is 0.85, not 0.8 (matches brainstorm requirement)
- 90-day staleness check uses `last_seen` field (already exists)
- The "no candidates" message matches the actual criteria
- All three conditions (confidence + observations + recency) are AND'd
