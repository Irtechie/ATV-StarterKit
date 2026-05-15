---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-007
title: "Auto-evolve on cadence + completion counter"
blockers: [slice-004, slice-006]
verification: verification-only
hitl: false
expected_files:
  - .github/skills/kanban-work/SKILL.md
status: pending
---

# Slice 7: Auto-Evolve on Cadence + Completion Counter

## What to Build

Track how many kanban features have been completed and automatically invoke `/evolve` every 5th completion. This creates a natural cadence for instinct promotion without requiring the user to remember to run it.

## Acceptance Criteria

- A completion counter is maintained in `.atv/kanban-completions.txt` (simple integer)
- Counter increments by 1 each time a kanban manifest reaches `status: completed`
- Every 5th completion (5, 10, 15...), invoke `/evolve` automatically
- The evolve invocation is silent unless it actually promotes instincts
- If evolve produces no candidates, log: `evolve: no candidates ready` and continue
- If evolve promotes instincts, log in manifest notes: `evolve: promoted N instincts to skills`
- The counter file is committed alongside the manifest

## Implementation

Add to the end of Step 5.6 (after the /learn invocation from Slice 4):

```markdown
6. **Check evolution cadence:**
   - Read `.atv/kanban-completions.txt` (create with `0` if missing)
   - Increment by 1
   - Write the new value back
   - If the new value is divisible by 5:
     - Invoke `/evolve` to check for promotable instincts
     - Log result in manifest notes: `evolve: promoted N instincts` or `evolve: no candidates ready`
   - If not divisible by 5: skip silently
   - Commit the counter file with the manifest update
```

## Counter Design

A simple plain-text integer file (`.atv/kanban-completions.txt`) rather than YAML or JSON:
- Easy to read/write without parsing
- Easy to inspect (`cat .atv/kanban-completions.txt`)
- Atomic increment (read → +1 → write)
- Committed to repo so counter persists across machines/sessions

## Scope Boundary

- Do NOT modify `/evolve` itself (already updated in Slice 6)
- Do NOT modify `/learn` (already wired in Slice 4)
- Do NOT change the cadence from 5 without brainstorm revision
- Do NOT make evolve blocking — if it fails or finds nothing, shipping continues
- Single file edit (kanban-work SKILL.md) plus new `.atv/kanban-completions.txt` on first run

## Test Scenarios

- Counter increments correctly (0→1, 4→5 triggers evolve, 5→6 does not)
- Counter file is created if missing (cold start)
- Evolve only runs on multiples of 5 (modulo check)
- Evolve failure doesn't block shipping
- Manifest notes record the outcome regardless of result
