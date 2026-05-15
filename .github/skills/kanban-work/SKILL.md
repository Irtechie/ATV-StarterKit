---
name: kanban-work
description: "Sequential executor for kanban-plan output. Runs all vertical slices in dependency order with fresh context per task, TDD enforcement, and HITL pauses. Use when the user says 'kanban work', 'run the kanban', 'execute all slices', or wants guided execution of a planned feature."
argument-hint: "[path to kanban manifest, or blank to find latest]"
---

# Kanban Work - Sequential Slice Executor

Run all vertical slices from a `kanban-plan` manifest in dependency order. Keep each slice scoped, enforce the requested verification mode, and pause on HITL tasks.

## Quick Start

1. Read the kanban manifest.
2. Validate the dependency DAG and statuses.
3. Confirm execution with the user unless they requested non-interactive execution.
4. Execute ready slices in topological order.
5. Update the manifest after each slice so the workflow is resumable.

## Input

<input> #$ARGUMENTS </input>

**If input is empty:** Scan `docs/plans/` for the most recent `*-kanban-*-manifest.md` file. If found, use it. Otherwise ask: "Which kanban manifest should I execute?"

**If input is a path:** Read the manifest at that path.

## Pre-Flight

1. **Read the manifest** - parse the YAML frontmatter to get the ordered slice list.
2. **Validate DAG** - confirm no cycles in blockers, all referenced slice IDs exist, and all slice files exist.
3. **Check status** - skip any slices already marked `done`. Resume from the first runnable `pending` slice.
4. **Check worktree** - note dirty or untracked files before executing so unrelated user changes are not staged or reverted.
5. **Confirm with user:** "Ready to execute N remaining slices in order. Proceed?"

Treat statuses as:

| Status | Action |
|--------|--------|
| `pending` | Eligible once blockers are `done` or `skipped` |
| `done` | Skip |
| `blocked` | Stop and ask whether to retry, skip, or abort |
| `skipped` | Skip but keep visible in the summary |

## Dependency Ordering

Execute with a topological sort:

1. Build a map of `slice_id -> slice`.
2. For each pending slice, check all `blockers`.
3. Run the first pending slice whose blockers are complete.
4. If pending slices remain but none are runnable, mark the manifest blocked and report the dependency problem.

## Execution Loop

For each slice in dependency order:

### Step 1: Check HITL Flag

If `hitl: true`:

- Present the slice title, description, and the specific question/decision needed.
- Stop and wait for user input.
- Record the user's decision in the slice plan.
- Update manifest status to `done` for this slice only if the decision completes the slice.
- Continue to the next runnable slice.

### Step 2: Deepen If Thin

If the slice plan has fewer than 3 acceptance criteria or no test scenarios:

- Run a lightweight deepening pass on this single slice.
- Add concrete test scenarios and likely file paths.
- Keep the pass bounded; do not re-plan the whole feature.

### Step 3: Execute

Use a fresh sub-agent when the platform supports delegated execution and the user has permitted it. Otherwise execute the slice locally while keeping the scope limited to this slice.

The execution prompt or local checklist must include:

```text
You are executing a single vertical slice. Complete it fully.

Kanban: <kanban_id>
Slice: <slice_id> - <title>
Verification mode: <tdd|integration|verification-only>

Plan contents:
<full slice plan content>

Instructions:
1. Read the plan completely.
2. Set up on the current branch.
3. Apply the verification mode:
   - tdd: write one failing test first, confirm it fails for the right reason, implement minimal code to pass, then refactor.
   - integration: write an integration test proving the path works, then implement the wiring.
   - verification-only: implement the change, then verify builds pass and no regressions.
4. Run relevant tests first, then the full test suite when practical.
5. Stage only files changed for this slice.
6. Commit only if the user asked for commits, with message: "feat(<slice-id>): <title>"

Do not modify other slices' files unless required for this slice.
Do not add scope beyond what the plan specifies.
Do not stage unrelated dirty or untracked files.
```

### Step 3.5: System-Wide Test Check

Before marking a slice done, pause and ask these questions — vertical slices cut through all layers, so side-effects matter:

| Question | What to do |
|----------|------------|
| **What fires when this runs?** Callbacks, middleware, observers, event handlers — trace two levels out from your change. | Read the actual code for callbacks on models you touch, middleware in the request chain, `after_*` hooks. |
| **Do my tests exercise the real chain?** If every dependency is mocked, the test proves logic in isolation — not interaction. | Write at least one integration test that uses real objects through the full callback/middleware chain. |
| **Can failure leave orphaned state?** If your code persists state before calling an external service, what happens when the service fails? | Trace the failure path. If state is created before the risky call, test that failure cleans up or that retry is idempotent. |
| **What other interfaces expose this?** Mixins, DSLs, alternative entry points. | Grep for the method/behavior in related classes. If parity is needed, add it now. |

**When to skip:** Leaf-node changes with no callbacks, no state persistence, no parallel interfaces. Purely additive changes (new helper, new partial) take 10 seconds to confirm "nothing fires, skip."

### Step 3.6: Diff-Scope Verification

After a slice completes, verify that the files actually changed match the slice's declared `expected_files`. This is a hard gate — the agent does not self-report, the actual git diff is checked.

1. **Get the actual diff:**

   ```bash
   git diff --name-only $(git merge-base HEAD main)..HEAD
   ```

   This produces the list of files modified by this slice relative to the branch baseline.

2. **Load the declared scope** from the slice plan's `expected_files` frontmatter field.

3. **Compare and enforce:**

   | Finding | Action |
   |---------|--------|
   | Files changed that are NOT in `expected_files` | **STOP.** Flag each out-of-scope file. Do not proceed to the next slice. Ask the user whether to amend the plan, revert the change, or accept the scope expansion. |
   | Files in `expected_files` that were NOT changed | Flag as potentially incomplete. Ask the user whether the slice is truly done or if work was missed. |
   | Perfect match | Proceed. |

4. **Log results** in the kanban manifest under the slice's `notes` field:

   ```text
   notes: "scope-check: 5/5 expected files changed, 0 out-of-scope"
   ```

5. **If the slice plan has no `expected_files` field**, the gate fails. Stop and require the field to be added before proceeding. Do not infer or guess the expected files — the plan must declare them explicitly.

This gate is mandatory. It cannot be skipped, overridden, or deferred.

### Step 3.7: Figma Design Sync (UI slices only)

If the slice involves UI changes and Figma designs exist:

1. Implement components following design specs
2. Use the **figma-design-sync** agent iteratively to compare
3. Fix visual differences identified
4. Repeat until implementation matches design

Skip this step entirely for non-UI slices.

### Step 4: Verify and Update

After the slice completes:

1. **Check result**
   - If yes: update manifest `status: done` for this slice and update the body table.
   - If no: update manifest `status: blocked`, add failure notes, and ask user how to proceed.

2. **Run verification**
   - Run the relevant test command for the repo.
   - If a full suite is too expensive or unavailable, explain the narrower verification used.

3. **Optional commit**
   - If the user asked for commits, stage only the manifest file for status updates and commit it separately.

4. Continue to the next runnable slice.

### Step 5: Completion

When all slices are `done` or intentionally `skipped`:

1. Update manifest `status: completed`.
2. Run final verification.
3. Report summary:

```text
Kanban <name> complete.
- N slices executed
- S slices skipped
- M tests added
- K files changed
Verification: <command/result>
```

4. Suggest `/ce-review` for a fresh-context review of the full feature.

### Step 6: Ship It

After all slices pass and review is complete:

1. **Quality gate**

   ```bash
   # Run full test suite
   # Run linting (per CLAUDE.md / project conventions)
   ```

2. **Capture screenshots** (if any slice touched UI)

   Use `agent-browser` to screenshot affected routes:
   ```bash
   agent-browser open http://localhost:3000/[route]
   agent-browser snapshot -i
   agent-browser screenshot output.png
   ```

3. **Create PR**

   ```bash
   git push -u origin $(git branch --show-current)

   gh pr create --title "feat: <feature name>" --body "$(cat <<'EOF'
   ## Summary
   - What was built (list slices completed)
   - Origin brainstorm: `<brainstorm_path>`

   ## Slices Executed
   | # | Slice | Verification | Status |
   |---|-------|-------------|--------|
   | 1 | <title> | tdd | ✅ |
   | 2 | <title> | integration | ✅ |

   ## Testing
   - Tests added per slice (TDD/integration as specified)
   - Full suite passing

   ## Post-Deploy Monitoring & Validation
   - **What to monitor:** [logs, metrics, dashboards]
   - **Expected healthy behavior:** [signals]
   - **Failure signals / rollback trigger:** [trigger + action]
   - **Validation window & owner:** [window + owner]

   ## Screenshots
   | Before | After |
   |--------|-------|
   | ![before](URL) | ![after](URL) |

   ---
   🤖 Generated with [Claude Code](https://claude.com/claude-code)

   Co-Authored-By: Claude <noreply@anthropic.com>
   EOF
   )"
   ```

   Omit the Screenshots section if no UI slices were executed. If there is truly no production/runtime impact, include: `No additional operational monitoring required: <reason>`.

4. **Update manifest** — set `status: shipped` and add the PR URL to the manifest frontmatter.

## Failure Handling

| Situation | Action |
|-----------|--------|
| Slice execution fails | Mark blocked, show error, ask user |
| Test suite fails after a slice | Mark blocked, show which tests broke, ask user |
| HITL pause | Present context, wait for user, record decision |
| User says "abort" | Mark remaining slices as `pending`, stop |
| User says "skip" | Mark slice `skipped`, continue to next runnable slice |

## Resume Support

Kanban work is resumable:

- Manifest tracks status per slice.
- Re-running `kanban-work` on the same manifest picks up where it left off.
- Already done or skipped slices are not rerun.

## Success Criteria

- No slice runs before its blockers are complete.
- Manifest frontmatter and body table reflect actual slice status.
- Each completed slice has verification evidence recorded in the final response or failure notes.
- No unrelated files are staged, committed, reverted, or overwritten.

## Integration

- **Input from:** `kanban-plan`
- **Deepening:** Optionally runs `deepen-plan` per slice before execution
- **Execution engine:** Fresh sub-agents when available, local execution otherwise
- **Verification:** Invokes `tdd` skill principles per slice when verification mode requires it
- **Post-completion:** Hand off to `ce-review` for fresh-context review
