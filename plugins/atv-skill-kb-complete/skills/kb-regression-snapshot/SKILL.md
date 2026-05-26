---
name: kb-regression-snapshot
description: Capture and verify deterministic regression snapshots between KB slices. Use after a slice passes all gates to record machine-checkable state, and before future slices start execution to prove previous slice behavior has not regressed.
---

# KB Regression Snapshot

Remember what earlier slices proved without relying on chat history.

## Modes

| Mode | When | Output |
|---|---|---|
| `capture <slice-id>` | After a slice passes `kb-check`, `kb-functional-test`, and `kb-qa` | `.atv/snapshots/<slice-id>.json` |
| `verify [before-slice-id]` | After `kb-work` Scope Lock and before executing the next slice | pass/fail for all prior snapshots |

Create `.atv/snapshots/` if missing. Never store secrets, cookies, tokens, raw credentials, or private response bodies.

## Snapshot Contract

Each snapshot is JSON:

```json
{
  "slice_id": "slice-001",
  "created_at": "2026-05-26T12:00:00Z",
  "git_ref": "<commit-or-worktree-hash>",
  "checks": [
    {
      "kind": "route|api-schema|file-checksum|command|browser-assertion",
      "target": "<url/path/command/file>",
      "assertion": "<deterministic condition>",
      "artifact": "<log/trace/output path>",
      "status": "pass"
    }
  ]
}
```

Keep snapshots compact. Store pointers to artifacts rather than large output.

## Capture

Capture only deterministic state that proves the slice still works:

- Route responses: status codes and required key DOM elements.
- API endpoint schemas: status, content type, required fields, shape hashes.
- Key file checksums for generated/config/runtime files the slice owns.
- CLI/build/test commands: command, exit code, timestamp, log path.
- Browser assertions: test file or trace path, route, locator/assertion summary.

Prefer existing proof from `kb-check`, `kb-functional-test`, and `kb-qa`. Do not rerun expensive checks just to duplicate proof unless no reusable artifact exists.

## Verify

1. Load all `.atv/snapshots/*.json` older than the current slice.
2. Re-run each check deterministically.
3. If any snapshot fails, STOP before new slice execution.
4. Mark the current slice `🔒 blocked` with:
   - failing snapshot path;
   - failing check target;
   - expected vs observed;
   - command/log/trace path.
5. Do not edit implementation files until the regression is resolved or the user explicitly parks/skips the affected work.

Regression failures are not QA nits. They mean the codebase regressed before the next slice began.

## Output

Capture:

```text
snapshot-capture: PASS slice-001 -> .atv/snapshots/slice-001.json
checks: route=2 api=1 files=3 commands=1 browser=1
```

Verify:

```text
snapshot-verify: PASS 4/4 snapshots
```

or:

```text
snapshot-verify: FAIL .atv/snapshots/slice-001.json
failed: <target> expected <condition> observed <actual>
artifact: <log/trace path>
```

Record the result in the manifest notes. Snapshot verification counts as machine-verifiable proof for `kb-complete`.
