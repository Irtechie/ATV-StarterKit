---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-002
title: "P0/P1 resolution gate + compound step"
blockers: [slice-001]
verification: verification-only
hitl: false
expected_files:
  - .github/skills/kanban-work/SKILL.md
status: pending
---

# Slice 2: P0/P1 Resolution Gate + Compound Step

## What to Build

Add two new substeps between the current Step 5 (Completion) and Step 6 (Ship It):

1. **Step 5.5: Resolution Gate** — P0/P1 findings from ce-review MUST be resolved before shipping. P2/P3 are logged but don't block.
2. **Step 5.6: Compound** — After resolution passes, invoke `ce-compound` to document what the feature taught the system.

## Acceptance Criteria

- Step 5.5 exists as a hard gate: P0/P1 findings block progression to Step 6
- P2/P3 findings are logged in the manifest `notes` field but do not block
- If no P0/P1 findings exist, the gate passes immediately
- Step 5.6 invokes ce-compound with brief context about what was built
- If nothing novel was learned (pure CRUD/scaffolding), ce-compound is skipped with a one-line note
- The compound step references that micro-learnings from per-slice notes feed into it
- Same voice, formatting, and numbering style as the rest of the document

## Implementation

Insert after Step 5 item 4 (the ce-review invocation from Slice 1) and before Step 6:

```markdown
### Step 5.5: Resolution Gate

Review findings from `ce-review` determine whether shipping is allowed:

| Severity | Action |
|----------|--------|
| P0 (critical) | STOP. Fix before proceeding. Re-run affected tests after fix. |
| P1 (important) | STOP. Fix before proceeding. |
| P2 (suggestion) | Log in manifest `notes`. Do not block. |
| P3 (nit) | Log in manifest `notes`. Do not block. |

This gate is mandatory. The agent MUST NOT proceed to Step 6 while unresolved P0/P1 findings exist.

After resolving all P0/P1s, update the manifest notes with a summary:
`review: P0=0 P1=2(resolved) P2=3(logged) P3=1(logged)`

### Step 5.6: Compound

After the resolution gate passes, document what this feature taught the system:

1. **Invoke `ce-compound`** with context: a one-sentence summary of what was built and any surprising patterns discovered during implementation.
2. ce-compound writes to `docs/solutions/` with YAML frontmatter — let it run without modification.
3. If the implementation was pure boilerplate (no novel patterns, no gotchas, no decisions worth preserving), skip with a manifest note: `compound: skipped — standard implementation, no novel patterns`
4. Per-slice micro-learnings from Step 4 notes feed into the compound context. Reference them when invoking ce-compound.
```

## Scope Boundary

- Do NOT modify Steps 1-4, 3.5, 3.6, 3.7
- Do NOT modify the ce-review invocation added by Slice 1
- Do NOT modify ce-compound or ce-review skills themselves
- Do NOT restructure Step 6 — it remains unchanged

## Test Scenarios

- Diff confirms only new content was added between Step 5 and Step 6
- P0/P1 gate language uses MUST/STOP (enforcement, not suggestion)
- The compound step references per-slice notes from Step 4
- Severity table matches ce-review's actual output levels (P0-P3)
