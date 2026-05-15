---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-001
title: "Wire ce-review as mandatory async"
blockers: []
verification: verification-only
hitl: false
expected_files:
  - .github/skills/kanban-work/SKILL.md
status: pending
---

# Slice 1: Wire ce-review as Mandatory Async

## What to Build

Replace the current suggestion in kanban-work Step 5 ("Suggest `/ce-review`") with a mandatory automated invocation. After all slices pass, the agent MUST invoke `ce-review` on the full feature diff. This is enforcement — not a suggestion the agent can skip.

## Acceptance Criteria

- Step 5 item 4 changes from "Suggest `/ce-review`" to an imperative instruction
- The agent is instructed to invoke ce-review with the full diff as context
- The review runs asynchronously (agent doesn't block waiting for output inline)
- The review findings (P0-P3) are captured for use by subsequent steps
- No other steps or gates are modified
- Same voice and formatting as existing document

## Implementation

In kanban-work SKILL.md, replace Step 5 item 4:

**Before:**
```
4. Suggest `/ce-review` for a fresh-context review of the full feature.
```

**After:**
```
4. **Invoke `ce-review`** — Run a full multi-agent code review on the feature diff.
   This is mandatory. Do not skip, defer, or make it optional.
   - Pass context: the full `git diff` of the feature branch against baseline
   - Capture the output: each finding has a severity (P0/P1/P2/P3) and confidence score
   - Store findings for the resolution gate (Step 5.5)
```

## Scope Boundary

- Do NOT modify Steps 1-4, 3.5, 3.6, 3.7, or Step 6
- Do NOT modify ce-review itself
- Do NOT add any new files — this is a single-file edit

## Test Scenarios

- Read the updated SKILL.md and confirm the language is imperative (MUST, mandatory)
- Confirm no other steps were touched by running a diff
- Confirm the instruction references passing git diff context and capturing severity levels
