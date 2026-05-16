# Kanban Skills — Vertical-Slice Agent Pipeline

> **Fork of [All-The-Vibes/ATV-StarterKit](https://github.com/All-The-Vibes/ATV-StarterKit)**
> Credit to the [All The Vibes](https://github.com/All-The-Vibes) community for the foundational ATV framework this builds upon.

The kanban skills are an **enforcement-first agent pipeline** that decomposes work into vertical slices and executes them through mandatory safety gates. Unlike the upstream `/lfg` pipeline (which trusts the agent to self-report), the kanban pipeline verifies every claim against ground truth: git diffs, lint output, browser renders, and test results.

---

## Architecture

```text
/klfg — the orchestrator (hands-off, one command)
   │
   ├─ /kanban-brainstorm    ← research-first requirements
   │
   ├─ /kanban-plan          ← vertical-slice decomposition with DAG
   │
   └─ /kanban-work          ← sequential execution engine
         │
         ├─ 3.0  Scope Lock             (proactive — blocks writes before they happen)
         ├─ 3    Execute                 (TDD / integration / verification-only)
         ├─ 3.5  System-Wide Tests       (trace side effects 2 levels out)
         ├─ 3.6  Diff-Scope Verification (reactive — checks git diff after the fact)
         ├─ 3.7  Destructive Guard       (blocks rm -rf, force push, DROP TABLE)
         ├─ 3.8  QA Gate                 (lint + browser) → /kanban-repair on failure
         ├─ 3.9  Figma Sync             (UI slices only)
         ├─ 4    Verify & Update
         └─ 5    Completion              (persists scope context for kanban-complete)
   │
   └─ /kanban-complete       ← post-work quality & learning
         │
         ├─ 1    ce-review              (multi-agent code review, scope pre-loaded)
         ├─ 2    Resolution Gate         (P0/P1 must be fixed)
         ├─ 3    Compound + Learn + Evolve
         └─ 4    Cleanup
```

---

## Skills

### `/klfg` — Full Pipeline Orchestrator

**What it does:** Chains brainstorm → plan → work → complete → DONE in one command. Interactive at three points: brainstorm Q&A, safety gate pauses during work, and a "continue to review?" prompt after all slices complete.

**When to use:** You want hands-off execution from idea to reviewed, documented code. One command, walk away.

**How it's different from `/lfg`:** The upstream `/lfg` runs a horizontal pipeline (plan → deepen → work → review → compound). `/klfg` enforces vertical slicing — every slice cuts through all layers end-to-end — and splits the pipeline: kanban-work owns slice execution + gates, kanban-complete owns quality review + learning.

```
/klfg "add user streak tracking"
```

---

### `/kanban-brainstorm` — Research-First Requirements

**What it does:** Runs market/landscape research *before* asking product questions. Questions are sharper because they're grounded in real prior art.

**Produces:** `docs/brainstorms/*-requirements.md`

**Key difference from `/ce-brainstorm`:** Inverts the order. ce-brainstorm asks questions first, then validates. kanban-brainstorm researches first, then asks — so questions are "given X exists, should we do Y?" instead of "what do you want?"

**When to use:**
- Prior art or competitive landscape materially changes framing
- Output is intended to feed `/kanban-plan` (vertical slices)
- You need research to inform the conversation before committing to an approach

---

### `/kanban-plan` — Vertical Slice Decomposition

**What it does:** Breaks a brainstorm/PRD/feature description into independently-executable vertical slices with a dependency DAG, verification strategy, and HITL flags.

**Produces:**
- `docs/plans/YYYY-MM-DD-000-kanban-<name>-manifest.md` — the DAG
- `docs/plans/YYYY-MM-DD-NNN-<type>-<name>-plan.md` — one per slice

**Key innovation — `expected_files`:** Every slice must declare which files it will create or modify. This isn't guidance — it's the contract that `kanban-work` enforces at Steps 3.0 and 3.6. If a slice doesn't declare its files, execution cannot begin.

**Slice format:**
```yaml
---
kanban_id: kb-2026-05-14-streaks
slice_id: slice-001
title: "Award points on lesson completion"
blockers: []
verification: tdd
hitl: false
expected_files:
  - path: src/services/points.py
    op: create
  - path: src/models/user.py
    op: edit
    scope: "add points_balance field"
  - path: tests/test_points.py
    op: create
status: pending
---
```

**Verification modes:**

| Mode | When | Gate |
|------|------|------|
| `tdd` | Business logic, behavior changes | Failing test → implement → passes |
| `integration` | Cross-boundary wiring, API contracts | Integration test proves path works |
| `verification-only` | Config, scaffolding, ops | Builds pass, no regression |
| `hitl` | UX taste, design judgment | Human confirms acceptable |

---

### `/kanban-work` — Sequential Slice Executor

**What it does:** Executes all slices from a kanban manifest in dependency order, running every safety gate per slice. After all slices complete, persists scope context for `kanban-complete` to pick up.

**The Gauntlet (per-slice):**

| Gate | Type | What It Checks |
|------|------|----------------|
| **3.0 Scope Lock** | Proactive | Blocks writes to files not in `expected_files` BEFORE they happen |
| **3.5 System Tests** | Analytical | Traces side effects 2 levels out — callbacks, middleware, observers |
| **3.6 Diff-Scope** | Reactive | `git diff --name-only` vs declared scope AFTER execution |
| **3.7 Destructive Guard** | Preventive | Blocks `rm -rf`, `git push --force`, `DROP TABLE`, etc. |
| **3.8 QA** | Observational | Lint (all slices) + browser verification (frontend slices) |
| **3.9 Figma Sync** | Visual | Compares rendered UI against Figma designs |

**Convention-matched test files:** Test files corresponding to an `expected_files` entry are automatically in scope. `src/foo.py` allows `tests/test_foo.py` without explicit declaration. Test files with no matching source are still flagged.

**Multi-agent board sync:** Uses `docs/kanban.md` as a shared handoff file. Agents claim slices before working and release after completing. Board wins over manifest if they diverge.

**Resume support:** Re-running `kanban-work` on the same manifest picks up where it left off. Already-done slices are not re-run.

---

### `/kanban-qa` — Quality Assurance Gate

**What it does:** Hard quality verification. The browser reports what rendered. The linter reports what's dirty. The model does not self-report.

**Runs:**
- **Lint** on `expected_files` — every slice, no exceptions
- **Browser verification** against acceptance criteria — frontend slices only

**Key features:**
- Continuous console monitoring (catches runtime errors during interaction)
- Diff-aware page scoping (only verifies pages the slice actually touches)
- Atomic commits per fix (one commit = one revert if it regresses)
- Enhanced stuck detection (reverts, multi-file spirals, same-file circles)

**On failure:** Invokes `kanban-repair`. Does NOT stop immediately — gives repair a chance to fix it surgically.

---

### `/kanban-repair` — Surgical Fix Loop

**What it does:** When QA finds failures, repair attempts targeted fixes without losing context. No handoff to a new agent — the executing agent keeps its full session.

**Progress-based, not count-based:**
- Each iteration must make measurable progress toward fixing the failure
- "Measurable" = fewer lint errors, more tests passing, closer to acceptance criteria
- If an iteration produces no measurable improvement, it counts toward stuck detection

**Stuck detection signals:**
| Signal | Interpretation |
|--------|---------------|
| Same failure message appears twice in a row | Stuck — stop |
| Agent reverts its own fix (2x = stuck) | Oscillating — stop |
| Fix touches 3+ files in one iteration | Scope creep — stop and ask |
| Same file edited → reverted → re-edited | Circular — stop |

**Hard ceiling:** 5 iterations maximum, regardless of progress. After that, STOP and surface to the user. This is not negotiable.

---

### `/kanban-complete` — Post-Work Quality & Learning

**What it does:** After `kanban-work` finishes all slices, runs the quality review and knowledge capture pipeline. Separated from kanban-work so the user gets a natural pause point before investing in review and documentation.

**Pipeline:**
1. **ce-review** — multi-agent code review with scope-verified file list pre-loaded from kanban-work's gates (skips redundant scope discovery)
2. **Resolution Gate** — P0/P1 findings must be fixed before proceeding. P2/P3 logged but don't block.
3. **Compound + Learn + Evolve** — document patterns (ce-compound), extract instincts (/learn), and promote mature ones (/evolve every 5th completion)
4. **Cleanup** — prune QA screenshots, trim observations log to 90 days

**When to use:** After `kanban-work` reports all slices complete. `klfg` prompts automatically. Standalone users invoke it manually with the manifest path.

**Standalone invocation:**
```
/kanban-complete docs/plans/2025-05-16-001-kanban-feature-manifest.md
```

---

## The Safety Philosophy

The kanban pipeline exists because **agents lie by omission.** They don't intend to — but when you ask "did you stay in scope?" they check their own memory, not the filesystem. When you ask "does it pass lint?" they recall what they think they wrote, not what the linter actually says.

Every gate in this pipeline checks ground truth:

| What we verify | How we verify it | Why the agent can't self-report |
|----------------|-----------------|-------------------------------|
| File scope | `git diff --name-only` | Agent may edit files without noticing |
| Code quality | Actual linter output | Agent may "fix" issues that still fail |
| Visual correctness | Browser screenshot | Agent can't render DOM in its head |
| Test results | Test runner exit code | Agent may misread stack traces |
| Destructive commands | Pattern matching on shell input | Agent may not realize a command is destructive |

---

## Integration with ATV

The kanban skills build on the ATV foundation:

| ATV Component | How Kanban Uses It |
|---------------|-------------------|
| `ce-brainstorm` | kanban-brainstorm extends it with research-first ordering |
| `ce-review` | Called by kanban-complete with scope-verified file list (skips redundant discovery) |
| `ce-compound` | Called by kanban-complete for novel patterns |
| `/learn` + `/evolve` | Called by kanban-complete after every completion |
| `docs/solutions/` | Consumed by `learnings-researcher` during planning |
| `docs/brainstorms/` | Produced by kanban-brainstorm, consumed by kanban-plan |
| `docs/plans/` | Produced by kanban-plan, consumed by kanban-work |
| `agent-browser` | Used by kanban-qa for browser verification |
| `kanban-repair` | Called by kanban-qa when checks fail |

---

## Quick Reference

| Command | Does What | Produces |
|---------|-----------|----------|
| `/klfg "feature"` | Full pipeline, one command | PR with everything |
| `/kanban-brainstorm "idea"` | Research → requirements | `docs/brainstorms/*-requirements.md` |
| `/kanban-plan path/to/reqs.md` | Slice decomposition | Manifest + per-slice plans |
| `/kanban-work path/to/manifest.md` | Execute all slices | Working code, scope context |
| `/kanban-complete path/to/manifest.md` | Review + learn | Reviewed code, documentation |
| `/kanban-qa` | Lint + browser checks | Pass/fail with repair attempt |
| `/kanban-repair` | Fix QA failures | Atomic fix commits |

---

## Comparison: `/lfg` vs `/klfg`

| Dimension | `/lfg` (upstream) | `/klfg` (kanban) |
|-----------|-------------------|------------------|
| Decomposition | Horizontal phases | Vertical slices |
| Scope enforcement | None — trusts the agent | Hard gates at Steps 3.0 + 3.6 |
| QA | Post-hoc (gstack-qa) | Per-slice (kanban-qa + repair) |
| Review | Separate step | kanban-complete (scope pre-loaded from work) |
| Learning | Separate step | kanban-complete (auto-cadence) |
| Resumability | Limited | Full — manifest tracks per-slice status |
| Multi-agent | Not designed for it | Board sync protocol in kanban.md |
| Destructive safety | `/gstack-careful` (optional) | Step 3.7 (mandatory, cannot be skipped) |

---

## File Layout

```
.github/skills/
├── klfg/SKILL.md              # Pipeline orchestrator
├── kanban-brainstorm/SKILL.md # Research-first requirements
├── kanban-plan/SKILL.md       # Vertical slice decomposition
├── kanban-work/SKILL.md       # Sequential executor + all gates
├── kanban-complete/SKILL.md   # Post-work review + learning
├── kanban-qa/SKILL.md         # Quality assurance (lint + browser)
└── kanban-repair/SKILL.md     # Surgical fix loop

docs/
├── brainstorms/               # Requirements docs from kanban-brainstorm
├── plans/                     # Manifests + slice plans from kanban-plan
├── solutions/                 # Institutional knowledge from ce-compound
├── kanban.md                  # Live board (multi-agent handoff)
└── kanban-done.md             # Archived completed features

.atv/
├── observations.jsonl         # Tool use log (auto-trimmed to 90 days)
├── instincts/project.yaml    # Learned patterns with confidence scores
├── qa-screenshots/            # Browser captures (pruned after PR)
└── kanban-completions.txt     # Counter for evolve cadence (every 5th run)
```

---

## Credits

- **[All-The-Vibes/ATV-StarterKit](https://github.com/All-The-Vibes/ATV-StarterKit)** — The foundational framework. Compound Engineering pipeline, learning system, observer hooks, agent architecture, and 45+ skills that this builds on.
- **[Compound Engineering](https://github.com/EveryInc/compound-engineering-plugin)** — Planning-to-knowledge pipeline by Every, Inc.
- **[gstack](https://github.com/garrytan/gstack)** — Sprint execution engine by Garry Tan / Y Combinator. Research from gstack's `/qa` informed kanban-qa's continuous console monitoring and enhanced stuck detection.
- **[agent-browser](https://github.com/vercel-labs/agent-browser)** — Browser automation layer by Vercel Labs. Powers kanban-qa's browser verification.
- **[Andrej Karpathy](https://x.com/karpathy/status/2015883857489522876)** — Behavioral guardrails philosophy. The kanban pipeline's enforcement-over-suggestion approach is a direct response to Karpathy's observation that "models make wrong assumptions on your behalf and just run along with them."
- **[mattpocock/skills](https://github.com/mattpocock/skills)** — Inspiration for kanban-plan's slice-to-issues pattern.

---

## License

Same as the parent repository. See [LICENSE](../LICENSE).
