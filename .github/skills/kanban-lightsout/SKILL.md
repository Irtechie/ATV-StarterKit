---
name: kanban-lightsout
description: "Autonomous sequential executor for kanban-plan output. Runs all vertical slices in dependency order with fresh context per task, TDD enforcement, and HITL pauses. Use when the user says 'lightsout', 'run the kanban', 'execute all slices', or wants hands-off execution of a planned feature."
argument-hint: "[path to kanban manifest, or blank to find latest]"
---

# Kanban Lightsout — Autonomous Sequential Executor

Run all vertical slices from a `kanban-plan` manifest in dependency order. Fresh context per slice. TDD enforcement. Pause on HITL tasks.

## Input

<input> #$ARGUMENTS </input>

**If input is empty:** Scan `docs/plans/` for the most recent `*-kanban-*-manifest.md` file. If found, use it. Otherwise ask: "Which kanban manifest should I execute?"

**If input is a path:** Read the manifest at that path.

## Pre-Flight

1. **Read the manifest** — parse the YAML frontmatter to get the ordered slice list
2. **Validate DAG** — confirm no cycles in blockers, all referenced slice IDs exist
3. **Check status** — skip any slices already marked `done`. Resume from first `pending` slice.
4. **Confirm with user:** "Ready to execute N remaining slices in order. Proceed?"

## Execution Loop

For each slice in dependency order (respecting `blockers`):

### Step 1: Check HITL Flag

If `hitl: true`:
- Present the slice title, description, and the specific question/decision needed
- **STOP and wait for user input**
- Record the user's decision in the slice plan
- Update manifest status to `done` for this slice
- Continue to next slice

### Step 2: Deepen (Optional, Lightweight)

If the slice plan is thin (fewer than 3 acceptance criteria or no test scenarios):
- Run a lightweight deepening pass on this single slice
- Add concrete test scenarios and file paths
- Keep it under 1 minute — don't over-research

### Step 3: Execute with Fresh Context

Spawn a **fresh sub-agent** (via Task tool, `general-purpose` type) for this slice. The sub-agent prompt MUST include:

```
You are executing a single vertical slice. Complete it fully.

**Kanban:** <kanban_id>
**Slice:** <slice_id> — <title>
**Verification mode:** <tdd|integration|verification-only>

**Plan contents:**
<full slice plan content>

**Instructions:**
1. Read the plan completely
2. Set up on the current branch
3. IF verification mode is "tdd":
   - Write a failing test FIRST (invoke `tdd` skill principles)
   - Confirm it fails for the right reason
   - Implement minimal code to pass
   - Refactor
   IF verification mode is "integration":
   - Write an integration test proving the path works
   - Implement the wiring
   IF verification mode is "verification-only":
   - Implement the change
   - Verify builds pass and no regressions
4. Run the full test suite
5. Commit with message: "feat(<slice-id>): <title>"

Do NOT modify other slices' files unless required for this slice.
Do NOT add scope beyond what the plan specifies.
```

### Step 4: Verify and Update

After the sub-agent completes:

1. **Check result** — did the sub-agent succeed?
   - If YES: update manifest `status: done` for this slice, commit manifest update
   - If NO: update manifest `status: blocked`, add failure notes, ask user how to proceed (retry / skip / abort)

2. **Run full test suite** — ensure no regressions from this slice

3. **Continue** to next slice

### Step 5: Completion

When all slices are `done`:

1. Update manifest `status: completed`
2. Run full test suite one final time
3. Report summary:
   ```
   Kanban <name> complete.
   - N slices executed
   - M tests added
   - K files changed
   All tests passing.
   ```
4. Suggest: "Run `/ce-review` for a fresh-context review of the full feature."

## Failure Handling

| Situation | Action |
|-----------|--------|
| Sub-agent fails on a slice | Mark blocked, show error, ask user |
| Test suite fails after a slice | Mark blocked, show which tests broke, ask user |
| HITL pause | Present context, wait for user, record decision |
| User says "abort" | Mark remaining slices as `pending`, commit manifest, stop |
| User says "skip" | Mark slice `skipped`, continue to next |

## Resume Support

Lightsout is resumable. If interrupted:
- Manifest tracks status per-slice
- Re-running lightsout on the same manifest picks up where it left off
- Already-done slices are skipped

## Integration

- **Input from:** `kanban-plan` (produces the manifest)
- **Deepening:** Optionally runs `deepen-plan` per-slice before execution
- **Execution engine:** Fresh sub-agents (ralph-loop pattern without root-level PRD.md/PROGRESS.md)
- **Verification:** Invokes `tdd` skill principles per-slice when verification mode requires it
- **Post-completion:** Hand off to `ce-review` for fresh-context review
