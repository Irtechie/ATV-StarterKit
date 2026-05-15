---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-004
title: "Auto-trigger /learn after compound"
blockers: [slice-003]
verification: verification-only
hitl: false
expected_files:
  - .github/skills/kanban-work/SKILL.md
status: pending
---

# Slice 4: Auto-Trigger /learn After Compound

## What to Build

After `ce-compound` runs (Step 5.6), automatically invoke `/learn` to extract instincts from the session. This closes the loop: work → review → compound → learn.

## Acceptance Criteria

- Step 5.6 ends with a mandatory `/learn` invocation
- The learn step runs after compound (not in parallel — it needs the observations written in Step 5.5)
- If no `.atv/instincts/project.yaml` exists, `/learn` creates it (this is already /learn's behavior)
- The instruction notes that observations from Step 5.5 are now available for /learn to consume
- Results are silent (no user-facing output unless new instincts were discovered)
- A manifest note records whether new instincts were created: `learn: 2 new instincts, 1 updated`

## Implementation

Append to the end of Step 5.6 (after the compound invocation):

```markdown
5. **Invoke `/learn`** — Extract instincts from this session's work.
   - Run after compound completes (observations from Step 5.5 are now available)
   - `/learn` reads: observations.jsonl, recent git history, docs/solutions/, existing instincts
   - Record result in manifest notes: `learn: N new instincts, M updated` or `learn: no new patterns`
   - This is automatic — do not ask the user whether to run it
```

## Scope Boundary

- Do NOT modify `/learn` skill itself — invoke it as-is
- Do NOT modify the observation-writing logic from Slice 3
- Do NOT run /learn per-slice (only at end, after compound)
- Single file edit

## Test Scenarios

- /learn invocation is positioned after compound (ordering matters)
- Language is imperative (automatic, not optional)
- Manifest note format is specified for tracking
- No user interaction required (silent unless something was learned)
