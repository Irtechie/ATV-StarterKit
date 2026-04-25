---
title: "feat: Include new skills (land, takeoff) in --guided installation"
type: feat
status: completed
date: 2026-04-24
---

# feat: Include new skills (land, takeoff) in --guided installation

## Overview

Two skills landed in `.github/skills/` during the past week — `land` and `takeoff` — but they were never copied into `pkg/scaffold/templates/skills/` (the embedded `//go:embed all:templates` tree the installer ships) and were never registered in `pkg/scaffold/catalog.go` or the guided-mode TUI in `pkg/tui/`. As a result, `npm install ... --guided` does not offer or install them.

This plan wires `land` and `takeoff` into the guided installer so users selecting any of the three presets (Starter / Pro / Full) get them by default, and customizers see them as discrete options in the appropriate category.

It also adds a defensive parity check (test) so future skills added to `.github/skills/` cannot silently miss the installer pipeline again.

## Problem Frame

**Symptom.** User runs `npm install -g @atv/starterkit && atv-starterkit --guided` (or the equivalent), completes the wizard, and the resulting repo has no `.github/skills/land/SKILL.md` or `.github/skills/takeoff/SKILL.md` — even though those skills were merged this week.

**Root cause.** The installer ships skills via `embed.FS` rooted at `pkg/scaffold/templates/skills/`. Files in `.github/skills/` (the source-of-truth dogfooding copy used by this repo's own Copilot configuration) are **not** automatically mirrored. Each new skill requires three coordinated edits:

1. Copy `SKILL.md` (and any other files) into `pkg/scaffold/templates/skills/<name>/`.
2. Register `<name>` in `coreSkillDirectories`, `orchestratorSkillDirectories`, or `easterEggSkillDirectories` in `pkg/scaffold/catalog.go`.
3. Surface the skill in `pkg/tui/categories.go` `atvCategoryMapping` so customize-mode users can see and toggle it.

When a skill is added but only step 1 (or only the dogfooding `.github/skills/` copy) is done, the `--guided` path silently omits it.

**Past-week audit.** The full list of skills added or modified in `.github/skills/` since 2026-04-17:

| Skill | Source copy (`.github/skills/`) | Template (`pkg/scaffold/templates/skills/`) | Catalog wired | TUI wired |
|---|---|---|---|---|
| `karpathy-guidelines` | (not present here, lives only as template) | ✅ present | ✅ in `coreSkillDirectories` | ✅ in `CategoryGuidelines` |
| `meme-iq` | ✅ added 2026-04-24 | ✅ added same commit | ✅ in `easterEggSkillDirectories` | ✅ in `CategoryEasterEgg` |
| `land` | ✅ added 2026-04-24 | ❌ **missing** | ❌ **missing** | ❌ **missing** |
| `takeoff` | ✅ added 2026-04-24 | ❌ **missing** | ❌ **missing** | ❌ **missing** |

Only `land` and `takeoff` are gaps. This plan closes them.

## Requirements Trace

- **R1.** After `--guided` install with the **Starter** preset, the resulting repo contains `.github/skills/land/SKILL.md` and `.github/skills/takeoff/SKILL.md`.
- **R2.** After `--guided` install with **Pro** or **Full** presets, the same files are present (presets share `LayerCoreSkills` / `LayerOrchestrators`, so this follows from R1 once the right layer is chosen).
- **R3.** In customize mode, `land` and `takeoff` appear as toggleable options in a category that makes sense to a user shipping work (Shipping).
- **R4.** The two skills are pre-selected by default in customize mode (matching the behavior of all non-easter-egg ATV skills today).
- **R5.** A regression test verifies that every directory under `pkg/scaffold/templates/skills/` is registered in exactly one of the three skill-directory slices in `catalog.go`. Adding a template without registering it must fail CI.
- **R6.** A second regression test verifies that any skill present in `.github/skills/` is also present under `pkg/scaffold/templates/skills/` (so the dogfooding copy and shipped copy stay in sync), with an explicit allow-list for skills that intentionally only live in one place.

## Scope Boundaries

- **In scope:** Wiring `land` and `takeoff` into the installer, plus parity tests to prevent recurrence.
- **Out of scope:**
  - Editing the content of `land` / `takeoff` SKILL.md files. They were reviewed and merged this week; this plan ships them as-is.
  - Reorganizing the existing preset structure or category taxonomy.
  - Adding new presets.
  - Backfilling other tools (Claude Code, Cursor) — `--guided` here means the ATV-starterkit Copilot installer.
  - Migrating the dogfooding source-of-truth model. Today `.github/skills/` and `pkg/scaffold/templates/skills/` are both maintained by hand; that's a known cost we accept and only guard with a parity test.

## Context & Research

### Relevant Code and Patterns

- `pkg/scaffold/catalog.go` — embeds templates and builds the install component list. Three slices control skill membership:
  - `coreSkillDirectories` (~line 165) — always-on planning/learning/quality skills.
  - `orchestratorSkillDirectories` (~line 187) — opt-in workflow orchestrators (`lfg`, `slfg`, `ralph-loop`, etc.).
  - `easterEggSkillDirectories` (~line 197) — opt-in fun extras (`meme-iq`).
- `pkg/scaffold/templates/skills/<name>/SKILL.md` — embedded via `//go:embed all:templates`. Walked by `skillComponents()` and filtered by the selected directories.
- `pkg/tui/categories.go` `atvCategoryMapping` — maps display categories to ATV skill entries shown in the customize multi-select. The shipping category currently lists `ce-work`, `lfg`, `slfg`, `ce-compound`, `ce-compound-refresh`, `claude-permissions-optimizer`. Land/takeoff fit cleanly here.
- `pkg/tui/presets.go` — Starter/Pro/Full all include `LayerCoreSkills` and `LayerOrchestrators`. Adding `land`/`takeoff` to either layer means all three presets pick them up automatically.

### Decision: which layer for `land` and `takeoff`?

They are session lifecycle workflows that wrap commit/push/PR (land) and backlog briefing (takeoff). The closest neighbors already in the repo are `lfg` and `slfg` in `orchestratorSkillDirectories`. Putting them in `coreSkillDirectories` is also defensible — they're broadly useful, not optional add-ons.

**Decision:** Add to `coreSkillDirectories`. Rationale:

1. They have no runtime prerequisites (unlike `lfg` which orchestrates other skills).
2. They are short, self-contained protocols — closer to `setup` and `ce-plan` in spirit than to `lfg`.
3. Putting them in core means the **Starter** preset gets them, which matches user intent: someone choosing the "lightest install" still expects basic session start/end protocols.
4. Easy to revisit later — moving them to orchestrators is a one-line slice edit.

### Institutional Learnings

No directly applicable `docs/solutions/` entry. The closest prior precedent is the `meme-iq` PR (#22, commit 1f30760), which added a skill end-to-end across all three wiring sites in a single commit. We follow that pattern.

### External References

None needed. This is a wiring change against an existing internal API.

## Key Technical Decisions

- **Copy templates into `pkg/scaffold/templates/skills/`, do not symlink or runtime-load.** The installer uses `//go:embed`, which requires real files at build time. This matches existing convention.
- **Place `land` and `takeoff` in `coreSkillDirectories`** (rationale above).
- **Pre-select by default in customize mode.** `defaultSelectedSkillKeys` in `pkg/tui/wizard.go` already pre-selects every non-easter-egg ATV skill, so the new entries inherit this behavior automatically once added to `atvCategoryMapping`.
- **Add a parity test** rather than a generator. A generator (auto-copy `.github/skills/` → `pkg/scaffold/templates/skills/`) is more invasive and changes the trust boundary between dogfooding and shipped artifacts. A test that fails CI when the two diverge gives the same protection at much lower cost.
- **Keep `.github/skills/<name>/` as the editable source** for skills that exist in both places. The shipped template is a snapshot.

## Open Questions

### Resolved During Planning

- **Q:** Should land/takeoff go in core or orchestrators? **A:** Core (rationale above).
- **Q:** Do they require any new directories in the installer (e.g., docs/changelog)? **A:** No — both are pure SKILL.md, no helper scripts.
- **Q:** Is `karpathy-guidelines` in sync? **A:** Yes — it lives only as a template (no `.github/skills/karpathy-guidelines/` exists in this repo); already in core + Guidelines TUI category.
- **Q:** Is `meme-iq` in sync? **A:** Yes — added in PR #22 across all three wiring sites.

### Deferred to Implementation

- **Q:** Exact wording for the TUI option labels for `land` and `takeoff`. **A:** Implementation will use one-line summaries derived from the SKILL.md descriptions; reviewer can tweak wording in the PR.
- **Q:** Which existing test file to extend vs. create new. **A:** Decide once the implementer reads `pkg/scaffold/` test layout. If a `catalog_test.go` exists it gets extended; otherwise a new `parity_test.go` is added next to `catalog.go`.

## Implementation Units

- [x] **Unit 1: Copy `land` and `takeoff` templates into the embedded skills tree**

**Goal:** Make the SKILL.md files available to `//go:embed` so the installer can write them.

**Requirements:** R1, R2

**Dependencies:** none

**Files:**
- Create: `pkg/scaffold/templates/skills/land/SKILL.md` (copy of `.github/skills/land/SKILL.md`)
- Create: `pkg/scaffold/templates/skills/takeoff/SKILL.md` (copy of `.github/skills/takeoff/SKILL.md`)

**Approach:**
- Verbatim copy of the two SKILL.md files. No content edits.
- After copying, run `go build ./...` to confirm the `embed.FS` picks them up cleanly.

**Patterns to follow:**
- `pkg/scaffold/templates/skills/lfg/SKILL.md` — single-file skill template precedent.
- `pkg/scaffold/templates/skills/meme-iq/SKILL.md` — most recent precedent (PR #22).

**Test scenarios:**
- Test expectation: none for this unit alone — file presence is exercised by Unit 4's parity test and Unit 5's install integration test.

**Verification:**
- `find pkg/scaffold/templates/skills/land pkg/scaffold/templates/skills/takeoff -type f` returns exactly the two SKILL.md files.
- `go build ./...` succeeds.

---

- [x] **Unit 2: Register `land` and `takeoff` in the catalog**

**Goal:** Make the installer treat both skills as part of the core skill layer, so all three presets ship them.

**Requirements:** R1, R2

**Dependencies:** Unit 1

**Files:**
- Modify: `pkg/scaffold/catalog.go`

**Approach:**
- Append `"land"` and `"takeoff"` to `coreSkillDirectories` (the slice currently containing `setup`, `learn`, `karpathy-guidelines`, etc.).
- Place them alphabetically or grouped with similar lifecycle skills — implementer's call. Suggested placement: after `setup` since they bracket a session.
- Do **not** add them to `orchestratorSkillDirectories` or `easterEggSkillDirectories`.

**Patterns to follow:**
- The existing entries in `coreSkillDirectories`. Style matches one-string-per-line.

**Test scenarios:**
- Happy path — Build a filtered catalog with `LayerCoreSkills` selected and assert that `.github/skills/land/SKILL.md` and `.github/skills/takeoff/SKILL.md` appear in the resulting `[]Component` paths.
- Negative — Build a filtered catalog **without** `LayerCoreSkills` and assert neither file appears (proves they aren't accidentally smuggled in via another layer).

**Verification:**
- `go test ./pkg/scaffold/...` passes.
- A small Go scratch program calling `BuildFilteredCatalog(detect.StackGeneral, []string{"core-skills"})` lists both files.

---

- [x] **Unit 3: Surface `land` and `takeoff` in the TUI customize categories**

**Goal:** Customize-mode users see both skills as discrete options, pre-selected by default, in a category that matches their purpose.

**Requirements:** R3, R4

**Dependencies:** Unit 2

**Files:**
- Modify: `pkg/tui/categories.go`

**Approach:**
- Add two entries to the `gstack.CategoryShipping` slice inside `atvCategoryMapping`:
  - `{Label: "Takeoff — backlog briefing at session start", Key: "core-skills:takeoff", Source: "atv"}`
  - `{Label: "Land — commit, push, and open a PR at session end", Key: "core-skills:land", Source: "atv"}`
- Wording is suggestive — implementer may refine after re-reading the SKILL.md files.
- Ordering: place `takeoff` before `ce-work` (start of session) and `land` after `claude-permissions-optimizer` (end of session). Keeps the category readable as a session timeline.
- Default selection is automatic — `defaultSelectedSkillKeys` selects every non-easter-egg ATV entry.

**Patterns to follow:**
- The existing `lfg`/`slfg` entries in `CategoryShipping`.
- The `meme-iq` entry's structure for category placement.

**Test scenarios:**
- Happy path — `BuildCategoryGroups` returns a Shipping group whose `Skills` slice contains `core-skills:takeoff` and `core-skills:land`.
- Default-selection — `defaultSelectedSkillKeys` for the Shipping group (with empty preset gstack set) includes both keys.
- Parse-roundtrip — `ParseSelections([]string{"core-skills:land", "core-skills:takeoff"})` returns `atvLayers` containing `core-skills` exactly once (deduped) and empty `gstackDirs`.

**Verification:**
- `go test ./pkg/tui/...` passes.
- Manual smoke: run `go run ./cmd --guided` in a scratch directory, choose Starter, choose customize, confirm both entries appear and are checked in the Shipping screen.

---

- [x] **Unit 4: Parity regression test — every template directory must be registered**

**Goal:** Prevent future drift where someone copies a skill into `pkg/scaffold/templates/skills/` but forgets to add it to one of the three slices.

**Requirements:** R5

**Dependencies:** Unit 2

**Files:**
- Create: `pkg/scaffold/parity_test.go` (or extend an existing `catalog_test.go` if present)

**Approach:**
- Read the embedded template tree for `templates/skills/`, collect immediate subdirectory names.
- Build the union of `coreSkillDirectories ∪ orchestratorSkillDirectories ∪ easterEggSkillDirectories`.
- Assert the two sets are equal. Failure message must list both "registered but no template" and "template but not registered" entries so the offender knows which side to fix.
- No allow-list — every templated skill should be assigned to a layer.

**Patterns to follow:**
- Standard Go `testing` package idioms.
- `embed.FS.ReadDir("templates/skills")` to enumerate without filesystem assumptions.

**Test scenarios:**
- Happy path — With the codebase post-Unit 2, the test passes.
- Negative (verified manually during development) — Temporarily delete `"land"` from `coreSkillDirectories` and rerun; the test must fail with a clear message naming `land`.
- Negative (verified manually) — Temporarily create `pkg/scaffold/templates/skills/test-fixture-skill/SKILL.md` and rerun; the test must fail naming the unregistered directory.

**Verification:**
- `go test ./pkg/scaffold/ -run TestSkillDirectoryParity` passes.

---

- [x] **Unit 5: Parity regression test — `.github/skills/` and `pkg/scaffold/templates/skills/` stay in sync**

**Goal:** Catch the exact failure this plan is fixing — a skill added to the dogfooding copy but never copied into the shipped templates.

**Requirements:** R6

**Dependencies:** Unit 1

**Files:**
- Create: `pkg/scaffold/dogfood_parity_test.go` (or fold into Unit 4's file)

**Approach:**
- Walk `.github/skills/` at test time using `os` (the dogfooding copy is at the repo root, not embedded).
- Walk the embedded `templates/skills/` directories.
- Define a small allow-list set in the test for skills that are **intentionally** template-only (today: `karpathy-guidelines`) or **intentionally** dogfooding-only (today: none, but reserve the slot).
- Assert that the symmetric difference, minus the allow-list, is empty.
- The test reads the repo via `runtime.Caller` or a relative `../../`-style traversal from the test file. Implementer chooses whichever idiom matches the rest of the test suite.

**Patterns to follow:**
- Look for existing tests that read repo-relative files (e.g., anything that opens `go.mod`); if none exists, use `filepath.Join(filepath.Dir(testFile), "../..")` derived from `runtime.Caller(0)`.

**Test scenarios:**
- Happy path — With Unit 1 done, the symmetric difference equals the allow-list and the test passes.
- Negative (verified manually) — Temporarily delete `pkg/scaffold/templates/skills/land/SKILL.md` and rerun; the test must fail naming `land` as "in `.github/skills/` but not in templates."
- Edge case — A skill present in both but whose SKILL.md content drifts: this test does **not** assert content equality. Document this explicitly in a code comment so future readers know the test is a presence check, not a content check. Content drift is an accepted cost — the dogfooding copy is the editable source and the template is a snapshot.

**Verification:**
- `go test ./pkg/scaffold/ -run TestDogfoodTemplateParity` passes.

---

- [x] **Unit 6: Smoke-test the full guided install end-to-end**

**Goal:** Confirm the user-visible behavior — `--guided` install with the default preset writes both skills.

**Requirements:** R1, R2

**Dependencies:** Units 1–3

**Files:**
- Modify: existing sandbox/integration test if one exercises the guided path. Otherwise, a manual smoke test is sufficient — this unit produces a checklist, not new code, if no integration test exists.
- Reference: `test/sandbox/sandbox_test.go` (already grep'd as a `--guided` consumer) — extend if it covers preset flows.

**Approach:**
- If `test/sandbox/sandbox_test.go` already exercises a preset install, add assertions for the two new SKILL.md files at the expected destination.
- Otherwise, perform a manual smoke run: `go run ./cmd --guided` against a scratch directory, choose Starter without customize, confirm `.github/skills/land/SKILL.md` and `.github/skills/takeoff/SKILL.md` exist in the output.

**Test scenarios:**
- Happy path — Starter preset, no customize: both files written.
- Happy path — Pro preset, no customize: both files written.
- Happy path — Full preset, no customize: both files written.
- Edge case — Starter preset, customize, deselect both options: neither file written. Confirms the toggle works in both directions.

**Verification:**
- All four scenarios above produce the expected outcome.
- If an automated integration test was extended, `go test ./test/...` passes.

## System-Wide Impact

- **Interaction graph:** Three coordinated edits per skill (template copy, catalog slice, TUI category). The new parity tests close this multi-site update into a single enforced contract.
- **Error propagation:** No runtime error paths affected. Only build-time embedding and install-time file writes.
- **State lifecycle risks:** None — both skills are static markdown.
- **API surface parity:** The installer's public-facing behavior changes additively: new files appear in installed repos. No flag changes, no removed behavior.
- **Integration coverage:** Unit 6's end-to-end smoke covers the full pipeline that unit tests cannot — embed + filter + write to disk.
- **Unchanged invariants:** Existing skills, presets, layer keys, and catalog behavior are unchanged. Existing repos that re-run the installer get the new skills additively (existing files are not overwritten, per the standard scaffold contract).

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| `land` or `takeoff` SKILL.md content depends on a helper script or asset that doesn't exist in the embedded tree | Both files are confirmed to be standalone single-file skills (verified via `git log --name-status`). If this changes in the future, Unit 4's parity test won't catch missing assets — but it will catch missing directories. Add to follow-up if asset-bearing skills land. |
| Future skills get added to `.github/skills/` and forgotten again | Unit 5's dogfood parity test fails CI immediately. |
| Choosing core over orchestrators surprises Starter-preset users with two extra files | Acceptable — both skills are short, low-noise, and explicitly about session lifecycle. Preset descriptions don't need updating; "core ATV planning, execution, review, and docs flow" already implies session lifecycle. |
| TUI label wording confuses customize-mode users | Plan calls out wording as deferred to implementation. PR review tightens copy. |

## Documentation / Operational Notes

- No CHANGELOG entry strictly required — this is an installer wiring fix, not a user-facing feature change. However, since `land` and `takeoff` are net-new to `--guided` users, a one-line note under the next release ("Guided installer now ships `land` and `takeoff` skills by default") is worth adding to `CHANGELOG.md` as part of Unit 2 or Unit 3.
- No README updates required — the README already describes the guided installer in general terms.
- No migration or rollout coordination needed.

## Sources & References

- Related code:
  - `pkg/scaffold/catalog.go` (`coreSkillDirectories`, `orchestratorSkillDirectories`, `easterEggSkillDirectories`, `skillComponents`, `BuildFilteredCatalogForPacks`)
  - `pkg/tui/categories.go` (`atvCategoryMapping`)
  - `pkg/tui/presets.go` (Starter/Pro/Full layer membership)
  - `pkg/tui/wizard.go` (`defaultSelectedSkillKeys`)
- Related PRs/commits:
  - `92b49bd` feat(skills): port /land and /takeoff to Copilot skills
  - `dfce627` fix(skills): address verified review findings on land + takeoff
  - `1f30760` feat(installer): add memeIQ Easter Egg scaffolding (#22) — exemplar precedent for end-to-end skill wiring
  - `f47e6e0` feat: add Karpathy Guidelines skill — exemplar of a template-only skill
