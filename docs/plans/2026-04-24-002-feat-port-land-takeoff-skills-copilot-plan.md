---
title: "feat: Port /land and /takeoff skills to ATV-starterkit for GitHub Copilot"
type: feat
status: active
date: 2026-04-24
---

# feat: Port /land and /takeoff skills to ATV-starterkit for GitHub Copilot

## Overview

Port two Claude Code skills — `/land` (session completion: commit → push → PR → handoff) and `/takeoff` (session kickoff: prioritized backlog briefing) — into the ATV-starterkit repo so they are usable via GitHub Copilot. The ATV-starterkit already hosts a large library of Copilot-compatible skills under `.github/skills/<skill-name>/SKILL.md`, so this port follows the established pattern rather than introducing anything new structurally.

The primary adaptation work is:

1. Remove Claude-Code-specific surface area (e.g., `ExitWorktree` tool calls, `mcp__backlog__*` MCP tool references) and replace with Copilot-equivalent behavior (plain shell / `gh` CLI / prose guidance).
2. Adjust skill frontmatter and triggers to match Copilot's skill invocation model as used by the sibling skills already in `.github/skills/`.
3. Align both skills with the ATV-starterkit's documented conventions (conventional commits, `docs/plans/` naming, optional Backlog.md support, gstack interop) as described in `.github/copilot-instructions.md`.

## Problem Frame

The user has two battle-tested session bookends — `/takeoff` to start a session and `/land` to finish one — defined as Claude Code user-scope skills at `~/.claude/skills/{land,takeoff}/SKILL.md`. They want these same flows available inside ATV-starterkit when driving work via GitHub Copilot, which in this repo means the Copilot Skills mechanism backed by `.github/skills/<name>/SKILL.md` files (see sibling skills like `git-commit/`, `git-commit-push-pr/`, `ce-work/`, etc.).

Directly copying the SKILL.md files is insufficient: the Claude versions assume tools and conventions that Copilot does not have (e.g., `ExitWorktree`, `mcp__backlog__task_create`, Claude's keyword-trigger "this is not an exact match" language). They also make assumptions about stack detection that should be tightened for this specific repo (Go primary, with Node tooling under `npm/`).

## Requirements Trace

- R1. Provide a Copilot-invokable `land` skill that runs commit → push → PR → handoff, matching the intent of `~/.claude/skills/land/SKILL.md`.
- R2. Provide a Copilot-invokable `takeoff` skill that produces a prioritized backlog briefing, matching the intent of `~/.claude/skills/takeoff/SKILL.md`.
- R3. Both skills must live under `.github/skills/<name>/SKILL.md` and follow the same frontmatter/structure shape as existing ATV-starterkit skills (`git-commit`, `git-commit-push-pr`).
- R4. Skills must use only Copilot-available tooling — shell, `gh`, `git`, optional `backlog` CLI. No Claude-Code-only tool references (`ExitWorktree`, `mcp__backlog__*`, `AskUserQuestion`, etc.).
- R5. Both skills must preserve the signature UX affordances that make the originals useful: the terminal banners (`🛬 PLANE LANDED — NICE WORK`, `✈️ TAKE OFF — NOW AT 30,000 FEET`), the "never merge the PR" rule, the "never stop before pushing" rule, and the bullet-group takeoff shape with emoji headers.
- R6. Skills must interoperate with this repo's conventions: conventional commits, `docs/plans/` as the plan home, graceful fallback when no `backlog/` directory exists, and deference to `.github/copilot-instructions.md`.
- R7. The pair must be registered/discoverable the same way other skills in this repo are (mere presence under `.github/skills/<name>/SKILL.md` appears to be the registration mechanism — verify during implementation).

## Scope Boundaries

**In scope:**
- Creating `.github/skills/land/SKILL.md` and `.github/skills/takeoff/SKILL.md`.
- Tailoring language, triggers, and tool references for GitHub Copilot.
- Detecting ATV-starterkit's actual stack (Go + Node) in the `land` quality-gate step.
- A short mention in `.github/copilot-instructions.md` under "Available Workflows" so users know the skills exist.

**Out of scope / non-goals:**
- Not porting to other coding agents (Cursor, Cline, etc.) — Copilot only.
- Not modifying the original Claude skills at `~/.claude/skills/{land,takeoff}/`.
- Not auto-invoking these skills via hooks; invocation remains user-initiated.
- Not building a Backlog.md integration for ATV-starterkit itself — the repo currently has no `backlog/` directory, and the skill should fall back gracefully, matching the original's edge-case handling.
- Not merging PRs as part of `/land`. That remains an explicit human action.

## Context & Research

### Relevant Code and Patterns

- `.github/skills/git-commit/SKILL.md` — canonical shape for a git-oriented Copilot skill in this repo. Uses YAML frontmatter with `name` and `description`, then numbered `### Step N:` sections. No `argument-hint` in the Copilot variant.
- `.github/skills/git-commit-push-pr/SKILL.md` — the closest functional cousin to `/land`. Shows the repo's pattern for `gh pr view` / `gh pr create`, exit-code capture, and PR description construction. `/land` should not duplicate this skill's PR-body authoring logic — it should either delegate ("follow the conventions in `git-commit-push-pr`") or stay terse on PR body, focusing on the commit-push-PR *sequence* rather than the body craft.
- `.github/skills/ce-work/` and similar — confirm the `.github/skills/<name>/SKILL.md` structure is the registration mechanism for Copilot skills in this repo.
- `.github/copilot-instructions.md` — lists available workflows (`/ce-brainstorm`, `/ce-plan`, etc.). New skills should be added to this index so Copilot users discover them.
- `~/.claude/skills/land/SKILL.md` and `~/.claude/skills/takeoff/SKILL.md` — source documents. Carry over the 10-step / 6-step structures, the banners, and the "never merge / never skip push" rules verbatim. Strip Claude-specific tool calls.

### Institutional Learnings

- Sibling Copilot skills in this repo drop Claude-specific frontmatter keys (`argument-hint`) — follow their shape for portability.
- Stack detection in `/land`'s quality-gate step should be *adaptive* — check repo root for `go.mod`, `package.json`, `pnpm-workspace.yaml`, `Cargo.toml`, `pyproject.toml` and run the matching commands. ATV-starterkit has both Go and Node surfaces (`go.mod` at root, `npm/` subproject), so the skill should run Go checks at root and note the Node subproject.
- Backlog CLI support is optional. ATV-starterkit does not currently have `backlog/` at the root; the skill should fall back as the original does.

### External References

- None required. This is an internal port with existing reference implementations.

## Key Technical Decisions

- **Copy-then-adapt, not rewrite.** Start from the original SKILL.md files and make surgical edits. Wholesale rewrites would lose hard-won UX details (banners, bullet-group shapes, the "never merge the PR" language).
- **Drop Claude-only tool references.** Replace `ExitWorktree(action: "keep")` with a prose note ("if working in a git worktree, leave it in place for PR review; remove manually when merged"). Replace `mcp__backlog__task_create` references with CLI-only (`backlog task create ...`) guarded behind a `command -v backlog` check.
- **Drop the `argument-hint` frontmatter key** to match sibling Copilot skills in this repo (`git-commit`, `git-commit-push-pr`). Argument handling instructions stay in the skill body.
- **Keep both banners exactly as-is.** They are part of the skill's identity and user muscle memory.
- **Keep "never merge the PR" as a critical rule in `/land`.** Copilot is more likely than Claude to interpret "finish the work" as "merge" — the rule is more important here, not less.
- **Delegate PR-body craft to `git-commit-push-pr` where possible.** `/land` should call out that PR body construction follows the conventions of the existing `git-commit-push-pr` skill rather than reimplementing it, keeping the two skills in sync.
- **Tighten stack detection for ATV-starterkit.** Detect stack from repo-root signal files; if `go.mod` is present, run `go build ./... && go vet ./...`; if `npm/package.json` exists, optionally run `(cd npm && npm run build)` when npm-subproject files changed in this session. This avoids the "run everything" sprawl the generic skill could fall into.
- **Mention the pair in `.github/copilot-instructions.md` under a new "Session bookends" heading** so users can discover them alongside the `/ce-*` workflows.

## Open Questions

### Resolved During Planning

- **How are Copilot skills invoked in this repo?** — By user slash command (e.g., `/land`, `/takeoff`) matching the skill directory name under `.github/skills/`. Registration is implicit via presence; no manifest update needed. Verified by inspecting sibling skills.
- **Should we keep `argument-hint`?** — No. Sibling skills don't use it. Argument handling stays in the skill body prose.
- **Which skill owns PR body craft?** — `git-commit-push-pr`. `/land` references it rather than duplicating.

### Deferred to Implementation

- **Exact wording of the Copilot-specific worktree fallback** — decide once during implementation; short prose note is sufficient.
- **Whether to add a one-line cross-reference from `git-commit-push-pr` back to `/land`** — evaluate during implementation; if natural, add it; if forced, skip.

## Implementation Units

- [ ] **Unit 1: Port `/land` skill to `.github/skills/land/SKILL.md`**

**Goal:** Create a Copilot-compatible `land` skill that executes the full session-completion checklist.

**Requirements:** R1, R3, R4, R5, R6

**Dependencies:** None

**Files:**
- Create: `.github/skills/land/SKILL.md`

**Approach:**
- Copy `~/.claude/skills/land/SKILL.md` as the starting point.
- Strip the `argument-hint` frontmatter key.
- Rewrite the description field to match the voice of sibling Copilot skills (concise, starts with action verbs, lists trigger phrases).
- Replace `ExitWorktree(action: ...)` calls in Step 7 with a short prose block: "If you are working inside a git worktree, leave it in place while the PR is open; remove it manually with `git worktree remove` only after the PR is merged or abandoned."
- Replace `mcp__backlog__task_create` references in Step 1 with guarded CLI usage: `if command -v backlog >/dev/null 2>&1 && [ -d backlog ]; then backlog task create ...; fi`, and fall back to the handoff list otherwise.
- In Step 2 (quality gates), replace the generic multi-stack block with an ordered, adaptive detection sequence tuned for this repo (check `go.mod` → run `go build ./... && go vet ./...`; check `npm/package.json` → run Node build only if npm files changed; add generic fallbacks for Python/Rust/pnpm for portability when this skill is copied to other repos).
- Preserve Steps 1–10 numbering and all critical rules verbatim. Keep the `🛬 PLANE LANDED — NICE WORK` banner rule unchanged.
- Add a one-line pointer in Step 6 that PR body construction should follow the conventions in `.github/skills/git-commit-push-pr/SKILL.md`.
- Update the "Project-specific considerations" footer to reference `.github/copilot-instructions.md` and `AGENTS.md` (both exist in this repo) instead of the Claude-world `CLAUDE.md`/`AGENTS.md` pair.

**Patterns to follow:**
- `.github/skills/git-commit/SKILL.md` for frontmatter shape and step headers.
- `.github/skills/git-commit-push-pr/SKILL.md` for `gh pr view` exit-code-capture idiom — reuse that exact pattern in Step 6.

**Test scenarios:**
- Happy path: Dirty working tree with Go changes → skill stages specific files, commits with conventional message, pushes, creates PR, emits banner.
- Happy path (nothing to commit): Clean tree, already pushed → skill skips Steps 4–5, still emits banner.
- Edge case: Branch has no upstream → skill runs `git push -u origin "$branch"` without attempting a pre-push rebase.
- Edge case: `backlog/` directory absent → Step 1 falls back to the handoff list without error.
- Edge case: Running inside a git worktree → skill emits the prose worktree note, does not attempt `ExitWorktree`.
- Error path: `git push` fails due to protected branch or hook rejection → skill surfaces the failure, does *not* emit the banner, loops back per the critical rules.
- Error path: Quality gate (`go build` or `go vet`) fails → skill halts before commit, does *not* emit the banner.
- Integration: After `/land` succeeds, `git log origin/<branch>..HEAD` returns empty and `git status` reports a clean tree.

**Verification:**
- Skill file parses as valid Markdown + YAML frontmatter (same linter posture as sibling skills).
- No string matches for Claude-only tool names (`ExitWorktree`, `mcp__backlog__`, `AskUserQuestion`) remain in the file.
- Invoking `/land` in a dummy branch with a single file change produces a commit, push, PR, handoff, and final banner line.

- [ ] **Unit 2: Port `/takeoff` skill to `.github/skills/takeoff/SKILL.md`**

**Goal:** Create a Copilot-compatible `takeoff` skill that produces a prioritized backlog briefing at session start.

**Requirements:** R2, R3, R4, R5, R6

**Dependencies:** None (parallelizable with Unit 1)

**Files:**
- Create: `.github/skills/takeoff/SKILL.md`

**Approach:**
- Copy `~/.claude/skills/takeoff/SKILL.md` as the starting point.
- Strip the `argument-hint` frontmatter key; move argument docs into a body section as sibling skills do.
- Rewrite the description field in the sibling-skill voice.
- Replace `mcp__backlog__task_list` references with CLI-only: prefer `backlog task list --plain --sort priority` when the CLI is available and `backlog/` exists; otherwise fall back to reading `docs/plans/*.md` frontmatter for active plan titles.
- Preserve the bullet-group output shape, emoji headers (🛫, 🔵, 🟢, ⚪, 🔴), and the `✈️ TAKE OFF — NOW AT 30,000 FEET` final banner verbatim.
- Preserve the edge-case list (no backlog, everything blocked, empty list, tasks without priority) and the "ask before starting work" recommendation pattern at Step 5.
- Update the fallback path: when `backlog/` is absent, scan `docs/plans/` for files with `status: active` frontmatter and render a bullet list of those plan titles + filenames as the actionable group. Mention the fallback in the output so the user knows why the shape differs.

**Patterns to follow:**
- `.github/skills/git-commit/SKILL.md` for frontmatter shape.
- Existing `docs/plans/YYYY-MM-DD-NNN-*-plan.md` filename convention — the fallback should parse these filenames to extract type and slug.

**Test scenarios:**
- Happy path: `backlog/` directory exists with a mix of To Do / In Progress / Done tasks → skill renders 🛫, 🟢, ⚪ groups correctly, emits banner.
- Happy path (this repo today): No `backlog/` directory, `docs/plans/` has several `status: active` plans → skill renders the plan-based fallback list, states that it's falling back, emits banner.
- Edge case: Task with unresolved dependency → skill annotates the line with `(blocked by <ID>)` rather than hiding it.
- Edge case: Empty backlog and empty plans → skill congratulates the user, suggests `/ce-ideate` or `/ce-plan`, still emits banner.
- Edge case: `--top 3` argument → skill truncates the top-priority group to 3 items.
- Edge case: Task title contains a `|` pipe character → skill escapes it as `\|`.
- Error path: `backlog` CLI is present but returns a non-zero exit code → skill reports the error honestly, falls back to `docs/plans/`, emits banner.

**Verification:**
- Skill file parses as valid Markdown + YAML frontmatter.
- No string matches for `mcp__backlog__` or other Claude-only tools remain in the file.
- Invoking `/takeoff` in ATV-starterkit today (no backlog present) returns the `docs/plans/` fallback list with the banner.

- [ ] **Unit 3: Index the new skills in `.github/copilot-instructions.md`**

**Goal:** Make the skills discoverable via the repo's Copilot onboarding file.

**Requirements:** R7

**Dependencies:** Units 1 and 2

**Files:**
- Modify: `.github/copilot-instructions.md`

**Approach:**
- Add a new subsection under the top-level workflow list (alongside `/ce-brainstorm` etc.) titled "Session bookends" with two bullets:
  - `/takeoff` — Prioritized backlog briefing to start a session
  - `/land` — Commit → push → PR → handoff to finish a session
- Keep the phrasing terse to match the rest of the file.
- Do not restructure or renumber existing sections.

**Patterns to follow:**
- Current bullet shape in `.github/copilot-instructions.md` (single-line, em-dash-separated).

**Test scenarios:**
- Test expectation: none — documentation-only change, verified by visual inspection and markdown lint.

**Verification:**
- Diff is limited to the new section; no unrelated lines changed.
- File still renders cleanly on GitHub.

## System-Wide Impact

- **Interaction graph:** Both skills shell out to `git` and `gh`; `/land` also optionally to `backlog`. No runtime callbacks or observers.
- **Error propagation:** Failures in `/land` quality gates or `git push` must halt the routine without emitting the success banner. `/takeoff` failures should prefer graceful degradation (fallback to `docs/plans/`) over halting.
- **State lifecycle risks:** `/land` must not leave partial state — if push fails after commit, it must surface that cleanly so the user can resolve and re-run. The "never `git add -A`" rule protects against accidental secret staging.
- **API surface parity:** Other Copilot skills in this repo invoke `gh pr view` with exit-code capture (see `git-commit-push-pr`). `/land` should use the same idiom rather than inventing a new one.
- **Integration coverage:** The `/land` → `git-commit-push-pr` reference relationship should be tested end-to-end at least once; drift between the two skills is a likely future bug source.
- **Unchanged invariants:** Existing skills (`git-commit`, `git-commit-push-pr`, `ce-*`) are untouched. Existing `docs/plans/` naming convention is respected by `/takeoff`'s fallback. `AGENTS.md` and `.github/copilot-instructions.md` conventions remain authoritative.

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| Copilot's slash-command invocation semantics differ subtly from Claude's and cause a skill to not trigger | Mirror the frontmatter + directory-name shape of existing working skills (`git-commit`, `git-commit-push-pr`); smoke-test by invoking `/land` and `/takeoff` once after creation. |
| `/land` runs `go build ./...` against a large module graph and stalls the session | Note the command's cost in the skill body; if it proves slow in practice, follow up with a targeted-test variant (e.g., `go build ./cmd/... ./pkg/...`) rather than blanket disabling the gate. |
| `/land` is interpreted as "finish and merge" by a Copilot user or future automation | Keep the "NEVER merge the PR unless the user explicitly says 'merge this PR'" line in the Critical Rules section, verbatim, and repeat it in the description where users scan first. |
| Drift between `/land`'s PR creation and `git-commit-push-pr`'s richer PR-body craft | `/land` references `git-commit-push-pr` as the source of truth for PR body conventions rather than restating them. Drift becomes a one-line skill update instead of two full rewrites. |
| `/takeoff` fallback to `docs/plans/` surfaces too many stale plans and overwhelms the briefing | Filter to `status: active` frontmatter only; cap at the same default `--top 5` as the backlog path. |

## Documentation / Operational Notes

- Update `.github/copilot-instructions.md` (Unit 3) is the only documentation change required.
- No rollout, feature-flag, or migration work. Skills take effect on the branch that introduces them.
- Rollback is a `git revert` of the introducing commit — skills are additive and have no persistent state.

## Sources & References

- Source skill: `~/.claude/skills/land/SKILL.md`
- Source skill: `~/.claude/skills/takeoff/SKILL.md`
- Sibling pattern: `.github/skills/git-commit/SKILL.md`
- Sibling pattern: `.github/skills/git-commit-push-pr/SKILL.md`
- Repo conventions: `.github/copilot-instructions.md`, `AGENTS.md`
