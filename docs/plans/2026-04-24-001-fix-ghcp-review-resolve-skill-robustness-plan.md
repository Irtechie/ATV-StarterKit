---
title: "fix: ghcp-review-resolve skill robustness"
type: fix
status: completed
date: 2026-04-24
---

# fix: ghcp-review-resolve skill robustness

**Target file:** `~/.claude/skills/ghcp-review-resolve/SKILL.md` (global user skill, not in this repo)

## Overview

The `ghcp-review-resolve` skill failed on its first real-world test (ATV-StarterKit PR #9) in four distinct ways that are all symptomatic of the same gap: the skill treats its happy path (Copilot reviewer available, small PR, no prior review, clean merge) as the only path, and bails out when reality deviates.

This plan fixes the skill so it degrades gracefully, detects idempotent state, and surfaces real blockers (merge conflicts) instead of being blocked by imagined ones (a "missing" reviewer that already completed its work weeks ago).

## Problem Frame

Observed on PR #9 run:

1. **Copilot reviewer 422** — `gh pr edit --add-reviewer @copilot` failed because Copilot is not a collaborator on All-The-Vibes/ATV-StarterKit. Skill says this is a "soft failure" but the actual run treated it as a reason to stop the whole pipeline.
2. **Already-resolved findings** — A prior Copilot review from 2026-04-03 had raised 8 findings; all 8 were fixed in commit `fee0c5e` and all 8 threads are resolved. The skill has no way to detect "this PR already went through the loop" and would re-raise or duplicate.
3. **Large PR overflow** — 143 files / +15,044 / −3,262. `gh pr diff` hit the 20k-line API limit. The adjudicator subagent had no reliable way to read the diff.
4. **Merge conflict not surfaced** — `mergeStateStatus=DIRTY`. The skill never checks mergeability, so the user has to discover the real blocker themselves after the skill runs.

Underlying cause: the skill's preflight is too narrow (only checks `gh auth` and "PR exists"), its reviewer-dispatch step has no fallback branch, and it has no concept of "prior run state" or "PR size class".

## Requirements Trace

- R1. Copilot-reviewer assignment failure must not abort the pipeline — skill must continue in single-reviewer mode.
- R2. Skill must detect prior completed Copilot reviews whose findings are already resolved, and skip re-requesting in that case.
- R3. Skill must classify PR size up front and choose a diff-fetch strategy that survives the 20k-line `gh pr diff` cap.
- R4. Skill must surface merge-conflict state as a blocker in its preflight output, with a clear recommended next action, rather than silently proceeding.
- R5. Skill's guardrails (no approve, no merge, no close) must remain intact through every new code path.
- R6. Skill must remain usable on small, clean, first-time PRs — the happy path must not regress.

## Scope Boundaries

In scope:
- Edits to `~/.claude/skills/ghcp-review-resolve/SKILL.md`
- New preflight section with explicit decision tree
- Fallback modes for missing Copilot
- Idempotency check against existing reviews/resolved threads
- Size-class branching for diff fetching

Out of scope:
- Resolving merge conflicts automatically (skill surfaces the blocker and recommends `ce-work`; does not fix)
- Implementing a `pr-review-toolkit`-side chunking strategy (that belongs in that skill)
- Changing how Copilot-review is configured on repos (that's a repo admin task)
- Adding new reviewers beyond Copilot + pr-review-toolkit

## Context & Research

### Relevant files

- `~/.claude/skills/ghcp-review-resolve/SKILL.md` — single file, 256 lines; all edits land here
- Evidence from PR #9 test run in user's feature description

### Related prior art

- `~/.claude/skills/ce-work/` — the skill recommended as a next step when merge conflicts are detected; referenced but not modified
- `pr-review-toolkit:review-pr` — invoked via `Skill()`; contract assumed stable
- GitHub REST endpoints already used: `/pulls/{n}`, `/pulls/{n}/reviews`, `/pulls/{n}/comments`, `/pulls/{n}/files`

## Key Technical Decisions

- **Decision:** Add a formal Step 0 (Preflight) that runs before any reviewer is contacted. **Rationale:** All four observed failures are detectable before any side-effecting call. Putting the checks up front means the skill either proceeds confidently or reports a single actionable blocker, never both.
- **Decision:** Detect "already-reviewed + resolved" via the GraphQL `reviewThreads { isResolved }` API, not by counting `/reviews`. **Rationale:** REST `/reviews` only tells you a review happened, not whether its threads are resolved. Resolution is the real signal.
- **Decision:** Use `gh api .../pulls/{n}/files --paginate` instead of `gh pr diff` for large PRs. **Rationale:** `/files` paginates and gives per-file patch snippets; `gh pr diff` is a single blob with the 20k cap.
- **Decision:** Size thresholds: ≤ 20 files OR ≤ 2,000 lines = "small" (full diff); otherwise "large" (per-file fetch, adjudicator iterates file-by-file). **Rationale:** 20k is GitHub's hard cap; 2k is a soft cap where full-diff reasoning stays useful in a single prompt.
- **Decision:** Copilot fallback is "single-reviewer mode", not abort. **Rationale:** pr-review-toolkit alone still provides value. Aborting wastes a user invocation.
- **Decision:** Merge conflict is a blocker but a **reported** blocker, not a crash. The skill prints the state, recommends `ce-work` to resolve conflicts, and exits cleanly with exit-reason context. **Rationale:** The skill's job is PR review remediation, not merge mediation; but it should tell the user what the real blocker is.

## Open Questions

### Resolved during planning

- **Q:** Should the skill attempt to add `@copilot` anyway and catch the 422? **A:** Yes — probe with the actual API call but treat a 422/403 as "Copilot unavailable" rather than letting it fail loudly. The probe is the most reliable detector.
- **Q:** What counts as "already reviewed and resolved"? **A:** A PR where every Copilot-authored review thread has `isResolved: true` AND the PR head SHA is newer than (or equal to) the SHA at which each of those threads was resolved. If the head SHA moved after resolution, treat threads as potentially stale and re-review.
- **Q:** Should we still run pr-review-toolkit on an "already resolved" PR? **A:** Yes, by default — pr-review-toolkit is independent of Copilot state and may find new things. User can opt out via argument (out of scope for v1).

### Deferred to implementation

- Exact `gh api graphql` query string for resolved-thread detection (verify field names at edit time)
- Whether `gh pr view --json mergeStateStatus,mergeable` covers all the DIRTY states we care about, or if we also need `mergeable_state` via REST

## High-Level Technical Design

> *This illustrates the intended preflight decision flow and is directional guidance for review, not implementation specification. The implementing agent should treat it as context, not code to reproduce.*

```
Step 0 — Preflight
  ├─ gh auth status ──────────────────── fail? stop with clear error
  ├─ detect PR ───────────────────────── none? ask user for #
  ├─ fetch PR metadata (size, merge state, head SHA)
  │    ├─ mergeStateStatus == DIRTY ──── report blocker, recommend ce-work, STOP
  │    ├─ changedFiles > 20 || lines > 2000 ─ mark size_class = "large"
  │    └─ else ──────────────────────── mark size_class = "small"
  ├─ probe Copilot availability
  │    ├─ try: gh pr edit --add-reviewer @copilot
  │    ├─ 422/403 ──────────────────── copilot_available = false
  │    └─ ok ──────────────────────── copilot_available = true
  ├─ check prior review state (GraphQL reviewThreads)
  │    └─ all Copilot threads resolved AND head SHA unchanged since
  │           resolution ──────────── prior_resolved = true
  │
  └─ Decide mode:
       reviewers   = (pr-toolkit) + (copilot if available && !prior_resolved)
       diff_fetch  = full-diff if small else per-file paginated
       if reviewers == {pr-toolkit} and prior_resolved and size_class=="large":
           report "nothing useful to do; existing Copilot review is resolved
                   and PR is too large for cheap re-review" and STOP
       else proceed to Step 1
```

## Implementation Units

- [ ] **Unit 1: Add Step 0 (Preflight) with size, merge, and Copilot detection**

**Goal:** Replace the current "Preconditions" section with an explicit preflight that gathers PR size, merge state, head SHA, and Copilot availability up front.

**Requirements:** R1, R3, R4, R6

**Dependencies:** None

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (Preconditions + Step 1 sections)

**Approach:**
- Rename "Preconditions" → "Step 0 — Preflight"
- Enumerate the exact `gh` / `gh api` calls that populate: `PR_NUMBER`, `PR_HEAD_SHA`, `CHANGED_FILES`, `ADDITIONS+DELETIONS`, `MERGE_STATE_STATUS`, `COPILOT_AVAILABLE`, `PRIOR_RESOLVED`
- Define size class thresholds (20 files / 2,000 lines) as named values so later sections reference them
- Document the preflight report: a compact table the skill emits before doing any work (mirrors the table the user already got on PR #9, but as an intentional artifact instead of a failure postmortem)

**Patterns to follow:**
- Current Step 1's `gh pr view --json ... -q ...` style for metadata fetch

**Test scenarios:**
- Happy path: small clean PR → preflight reports all green, proceeds to Step 1
- Edge case: PR #9 scenario — DIRTY merge state → preflight reports blocker and stops with ce-work recommendation
- Edge case: 15k-line PR → preflight reports `size_class=large` and picks per-file strategy
- Edge case: repo without Copilot → preflight marks `copilot_available=false` and continues in single-reviewer mode
- Integration: user invokes skill with explicit PR number arg → preflight uses that, skips auto-detect

**Verification:**
- SKILL.md parses as valid Markdown with frontmatter intact
- Every preflight output field is referenced by at least one later step
- No later step assumes Copilot is available without checking `COPILOT_AVAILABLE`

- [ ] **Unit 2: Copilot fallback — single-reviewer mode**

**Goal:** Make Copilot assignment optional without reducing the rest of the pipeline's usefulness.

**Requirements:** R1, R5

**Dependencies:** Unit 1

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (Step 2 "Request both reviews in parallel")

**Approach:**
- Reword Step 2a from "fire and hope" to "conditional on `COPILOT_AVAILABLE`"
- When Copilot is unavailable, explicitly log: `Copilot unavailable on this repo (422 on reviewer assignment) — continuing in single-reviewer mode with pr-review-toolkit only.`
- Keep Step 3 (polling) aware of which reviewers are expected — don't wait 10 minutes for a reviewer that was never requested
- Update Step 4 (Collect findings) to gracefully handle an empty Copilot bucket (no overlap detection possible, but unique-finding path still works)

**Patterns to follow:**
- Current soft-failure language for pr-review-toolkit (already tolerates missing results at the 10-minute cap)

**Test scenarios:**
- Happy path: Copilot available → Step 2 fires both, Step 3 waits for both
- Edge case: Copilot unavailable (422) → Step 2 fires only pr-review-toolkit, Step 3 waits only for it
- Error path: Both reviewers fail to produce findings within the cap → Step 3 stops and reports (unchanged from current behavior)
- Integration: overlap-detection logic in Step 4 correctly handles the single-source case (no overlap bonus, but adjudicator still gets the findings)

**Verification:**
- No code path treats a 422 on Copilot as fatal
- Polling loop's expected-reviewer set is driven by Step 0's flags, not hardcoded

- [ ] **Unit 3: Idempotency — detect already-resolved prior reviews**

**Goal:** If Copilot (or any reviewer) has already reviewed this PR and all threads are resolved at the current head SHA, skip re-requesting that reviewer.

**Requirements:** R2, R6

**Dependencies:** Unit 1

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (Step 0 and Step 2)

**Approach:**
- Add a preflight sub-step that uses GraphQL to fetch `reviewThreads { isResolved, comments { author, commit { oid } } }` for the PR
- Classify each thread as resolved-and-fresh (resolved after current HEAD) vs resolved-but-stale (HEAD moved after resolution) vs open
- Derive `PRIOR_RESOLVED` = true when every Copilot thread is resolved-and-fresh
- In Step 2, if `PRIOR_RESOLVED`, skip Copilot re-request and log the reason
- Document the escape hatch: user can force a re-review by passing `--force` (or equivalent argument convention in the skill's arg handling) — deferred to implementation for exact arg syntax

**Patterns to follow:**
- GitHub GraphQL `reviewThreads` field usage, already common in pr-review-toolkit's own implementation

**Test scenarios:**
- Happy path: first-time PR, no prior reviews → `PRIOR_RESOLVED=false`, normal flow
- Edge case: PR #9 scenario — 8 prior findings all resolved, head SHA unchanged → `PRIOR_RESOLVED=true`, skip Copilot re-request
- Edge case: prior resolved but head SHA advanced → `PRIOR_RESOLVED=false`, re-request (threads are stale)
- Edge case: some threads resolved, some still open → `PRIOR_RESOLVED=false`, re-request
- Integration: skill still runs pr-review-toolkit even when Copilot is skipped for idempotency

**Verification:**
- Idempotency check never misclassifies a stale-HEAD resolution as fresh
- When `PRIOR_RESOLVED=true`, skill output clearly explains why Copilot was skipped

- [ ] **Unit 4: Large-PR diff strategy**

**Goal:** Stop using `gh pr diff` for large PRs; use per-file patch iteration instead so the adjudicator has real diff context.

**Requirements:** R3

**Dependencies:** Unit 1

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (Step 5 "Adjudicate findings with a subagent")

**Approach:**
- In Step 5, branch on `size_class`:
  - `small`: current full-diff approach via `gh pr diff`
  - `large`: fetch `gh api repos/{owner}/{repo}/pulls/{n}/files --paginate` to get per-file patch snippets; adjudicator is given only the files referenced by findings, not the whole diff
- Update the adjudicator prompt to say "you have access to this PR's file list and per-file patches; request additional file context via Read if needed" instead of "here is the full diff"
- Cap adjudicator per-file reads at ~30 files per run; if findings span more, report truncation and adjudicate the remainder in a second pass (still cheaper than burning 20k+ lines at once)

**Patterns to follow:**
- Pagination patterns already in Step 3 (`--paginate`)

**Test scenarios:**
- Happy path: small PR (say, 5 files, 200 lines) → uses full-diff path, adjudicator reasons over one blob
- Edge case: 143-file PR (PR #9) → uses per-file path, adjudicator only reads files referenced by findings
- Edge case: finding references a file not in the PR's changed-files list → adjudicator flags it as "finding not grounded in the diff" and drops it (new rejection category)
- Error path: `gh pr diff` returns the 20k-line cap error on a small PR (shouldn't happen but defensive) → skill catches, falls back to per-file path, continues

**Verification:**
- `gh pr diff` is no longer the only diff source
- Adjudicator prompt explicitly names the diff-fetch mode it's operating in so behavior is debuggable

- [ ] **Unit 5: Merge-conflict blocker surfacing**

**Goal:** When the preflight detects `mergeStateStatus=DIRTY`, report it as a blocker with an actionable next step and exit cleanly.

**Requirements:** R4, R5

**Dependencies:** Unit 1

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (Step 0 and Step 8 "Final summary")

**Approach:**
- In Step 0, after fetching merge state, if DIRTY: emit the preflight table, then a single-paragraph "Recommended next action" with two options (`git rebase` manual, or `Skill(skill="compound-engineering:ce-work", args="resolve the merge conflicts")`), then stop the skill with a clean exit
- Make "stop cleanly after reporting blockers" an explicit skill mode, not an error state — no crash, no partial side effects
- In Step 8's "what not to do" list, reinforce: a reported blocker is not a failure; the skill did its job by identifying what's in the way

**Patterns to follow:**
- The preflight table format the test run actually produced on PR #9 — it was already the right artifact; the skill just needs to emit it intentionally

**Test scenarios:**
- Happy path: mergeable PR → no blocker, proceeds
- Edge case: DIRTY state → skill reports table, recommends ce-work, exits with no side effects
- Edge case: `UNKNOWN` merge state (GitHub hasn't computed yet) → skill waits up to 30s for a definitive state, then proceeds or reports
- Guardrail: blocker-exit path must not post any PR comments, submit any reviews, or run any fix loop

**Verification:**
- No `gh api` mutation calls happen on the blocker-exit path
- User-visible output is informational and actionable, not apologetic

- [ ] **Unit 6: Skill frontmatter and example-run update**

**Goal:** Reflect the new behaviors in the skill's description, example run, and guardrails so users and agents discover them.

**Requirements:** R1, R2, R3, R4, R5, R6

**Dependencies:** Units 1–5

**Files:**
- Modify: `~/.claude/skills/ghcp-review-resolve/SKILL.md` (frontmatter description, "Example run shape", Guardrails)

**Approach:**
- Extend the frontmatter `description` to mention "surfaces merge conflicts and prior-review state as explicit preflight output" so the auto-triggering heuristic recognizes those user intents
- Update the example run to show preflight output as a rendered block (the PR #9 shape the user pasted is a good template)
- Add a second example showing the "Copilot unavailable, large PR, pr-toolkit only" path
- Add one line to Guardrails: "If preflight reports a blocker, stop cleanly — do not proceed into Steps 1–8."

**Patterns to follow:**
- Existing example-run block in SKILL.md (already in the right shape, just needs a second scenario)

**Test scenarios:**
- Test expectation: none — this is documentation-only. Manual verification: re-run the skill on PR #9 and confirm output now matches the documented example flow for the merge-conflict + prior-resolved path.

**Verification:**
- Frontmatter YAML still parses
- Description length stays within the ~500-char conventional range
- Example-run section includes both happy and degraded paths

## System-Wide Impact

- **Interaction graph:** Skill invokes `pr-review-toolkit:review-pr` (unchanged contract), `gh api` REST + GraphQL endpoints, and recommends but does not invoke `compound-engineering:ce-work` on the merge-conflict path.
- **Error propagation:** New preflight can exit the skill cleanly before any side effect; all later steps assume preflight passed and its flags are set. If a later step is ever invoked without preflight flags, it must refuse to run rather than assume defaults.
- **State lifecycle risks:** Idempotency check reads PR state that could change between preflight and Step 2. Mitigation: re-check head SHA immediately before any mutation (reviewer request, comment post, commit push) — this is already in the skill's existing Guardrail ("If the PR head SHA changes mid-run, stop").
- **API surface parity:** Only one skill file changes. No other skill depends on this one's internal structure.
- **Integration coverage:** Manual re-run on PR #9 is the end-to-end integration test. Expected new output: preflight table identifies DIRTY + prior-resolved, recommends ce-work, exits cleanly with zero inline comments posted.
- **Unchanged invariants:** Guardrails (no approve, no merge, no close, no acting on rejected findings, no fabricated text) remain intact. No change to the inline-fix loop's per-comment verification discipline.

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| GraphQL `reviewThreads` schema differs from assumption | Verify with a probe query before finalizing Unit 3; keep a REST fallback that reads `/reviews` + `/comments` and heuristically flags "resolved" by reply chain if GraphQL is unavailable |
| Size-class thresholds (20 files / 2000 lines) are wrong for some repos | Values are named and documented in Step 0 so future tuning is a single-point edit; not critical for v1 correctness |
| `mergeStateStatus` can transiently be `UNKNOWN` right after a push | Unit 5 includes a short wait-and-recheck loop before declaring DIRTY |
| Users expect Copilot review even when the repo hasn't enabled it | Unit 2's single-reviewer-mode log message explicitly names the cause (422 / not a collaborator) and recommends enabling Copilot on the repo if they want dual-review |
| Skill file grows past readable length | Total added content estimated < 100 lines; current file is 256 lines, post-change ~350. Still within one-screen-scroll territory |

## Documentation / Operational Notes

- No repo docs in ATV-StarterKit need to change — the skill lives under `~/.claude/skills/` globally.
- If the user maintains a "skills changelog" anywhere, a one-line entry covering "preflight with size/merge/prior-review detection + graceful Copilot fallback" is warranted.

## Sources & References

- PR #9 test run output (in user's feature description for this plan)
- `~/.claude/skills/ghcp-review-resolve/SKILL.md` (current version, 256 lines)
- GitHub REST: [Pulls API — Get a pull request](https://docs.github.com/en/rest/pulls/pulls)
- GitHub GraphQL: `PullRequest.reviewThreads` field
- Referenced but not modified: `compound-engineering:ce-work` skill
