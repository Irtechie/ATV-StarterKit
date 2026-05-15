---
kanban_id: kb-2026-05-14-compound-loop
slice_id: slice-003
title: "Feed resolved findings into observations"
blockers: [slice-002]
verification: verification-only
hitl: false
expected_files:
  - .github/skills/kanban-work/SKILL.md
status: pending
---

# Slice 3: Feed Resolved Findings into Observations

## What to Build

After P0/P1 findings are resolved (Step 5.5), append structured observation entries to `.atv/observations.jsonl` so that `/learn` can extract instincts from review patterns.

This is the novel integration — no existing system feeds review findings into the learning pipeline.

## Acceptance Criteria

- Step 5.5 includes an instruction to write resolved P0/P1 findings as JSONL entries
- The observation format matches the existing schema: `{ts, hook, tool, args, cwd, result}`
- Each resolved finding maps to one observation entry
- The observation `hook` field identifies these as review-sourced (e.g., `"hook": "ce-review"`)
- The `result` field contains the finding type, severity, and resolution summary
- If `.atv/observations.jsonl` doesn't exist, create it
- P2/P3 findings are NOT written as observations (they're noise, not signal)

## Implementation

Add to the end of Step 5.5 (after the resolution summary), a new sub-section:

```markdown
**Feed learnings to the observation log:**

For each resolved P0/P1 finding, append one line to `.atv/observations.jsonl`:

```json
{"ts":"<ISO-8601>","hook":"ce-review","tool":"kanban-work","args":{"finding_type":"<category>","severity":"P0|P1","resolution":"<what was fixed>"},"cwd":"<repo-root>","result":"resolved"}
```

This connects the review → learn pipeline. Only P0/P1 findings are worth learning from — P2/P3 are style preferences, not systemic patterns.

Create `.atv/observations.jsonl` if it doesn't exist. Append, never overwrite.
```

## Observation Format Rationale

Using the existing schema from `.github/hooks/scripts/observe.js`:
- `ts` — ISO timestamp (when the finding was resolved)
- `hook` — `"ce-review"` (source identification)
- `tool` — `"kanban-work"` (what invoked it)
- `args` — structured data: finding_type, severity, resolution description
- `cwd` — repo root path
- `result` — `"resolved"` (confirms it was actioned)

## Scope Boundary

- Do NOT modify `/learn` — it already reads observations.jsonl
- Do NOT write P2/P3 findings to observations
- Do NOT modify the observation format for existing hook-sourced entries
- Single file edit only

## Test Scenarios

- The JSONL format in the instruction is valid JSON (parseable)
- Only P0/P1 findings generate observations (P2/P3 explicitly excluded)
- The instruction is placed after resolution, not before (can't log what hasn't been fixed)
- The `hook` field uses a distinct value so /learn can filter by source
