---
date: 2026-05-14
topic: kanban-compound-loop
brainstorm_style: kanban-brainstorm
---

# Close the Compound Engineering Loop in Kanban Skills

## Problem Frame

kanban-work executes vertical slices well (Plan + Work) but stops short of the full compound engineering loop. Review is a suggestion the agent can skip. Compound doesn't exist. Learning never triggers. The system builds features but doesn't get smarter from building them.

The gap: after a kanban completes, nothing captures what was learned, nothing feeds that back into the system, and nothing prevents the same mistakes on the next kanban.

## Research Summary

**Findings that shaped requirements:**

- ExpeL (AAAI 2024) proves batch insight extraction across related tasks produces higher-quality learnings than per-task capture — this validates "compound at end, not per-slice" — Source: arxiv.org/abs/2308.10144
- CrewAI's recency decay (`0.5^(age_days/half_life)`) solves stale-learning accumulation — ATV currently has no staleness mechanism — Source: docs.crewai.com/concepts/memory
- No existing system feeds code review findings back into a learning pipeline — this is a novel integration point — Source: landscape scan of Devin, Claude Code, Windsurf, CrewAI, ExpeL
- ATV's `/learn` → `/instincts` → `/evolve` pipeline already maps to ExpeL's confidence-gated promotion (>0.8 threshold ≈ ExpeL's importance counter surviving DOWNVOTES) — the architecture exists but is never triggered
- Context bloat is the #1 failure mode of accumulated learnings — ATV avoids this by graduating instincts to on-demand skills rather than stuffing a global prompt file

**Confidence:** High — the architecture already exists in ATV, validated by academic research. The work is integration, not invention.

## Requirements

**Review Integration**

- R1. After all kanban slices complete (Step 5), kanban-work MUST invoke `ce-review` on the full diff automatically — not suggest it
- R2. Review runs async/background — it does not block the developer from seeing the completion summary
- R3. P0/P1 review findings MUST be resolved before the PR ships (Step 6 — "Ship It" in kanban-work). P2/P3 are logged in the manifest but do not block
- R4. If ce-review is already running or was already invoked for this branch, skip duplicate invocation

**Compound Step**

- R5. After review completes and P0/P1 findings are resolved (or immediately if none found), kanban-work invokes `ce-compound` automatically
- R6. ce-compound receives context about all slices executed, the review findings (if any), and what was resolved
- R7. If the compound step determines nothing novel was learned (pure CRUD, scaffolding), it logs "No novel patterns — standard implementation" and skips doc creation
- R8. Compound output lands in `docs/solutions/` per existing ce-compound behavior — no new location

**Review → Learning Pipeline**

- R9. Resolved P0/P1 findings from ce-review become observations fed to `/learn` — the fix + the reason it was needed
- R10. Only P0/P1 findings get fed as observations. P2/P3 are too noisy for the learning pipeline
- R11. The observation format must include: what was wrong, why, and what fixed it — not just "fixed bug"

**Automatic Learning Trigger**

- R12. After compound completes, kanban-work auto-runs `/learn` to extract instincts from the session
- R13. `/learn` runs regardless of whether compound produced a doc — review findings alone may contain learnable patterns
- R14. The full loop is hands-off once kanban-work starts: Work → Review (async) → Resolve P0/P1s (interactive if needed) → Compound → Learn → Ship

**Automatic Evolution**

- R15. Every 5th kanban completion (tracked in manifest or `.atv/` state), auto-run `/evolve` silently
- R16. Evolution threshold: confidence ≥0.85 AND ≥5 observations AND passes staleness check
- R17. No user approval gate — the threshold + staleness guard is the quality control. User reviews promoted skills in their next PR diff
- R18. If no instincts qualify, `/evolve` exits silently with no output

**Instinct Recency Decay**

- R19. Instincts that are not re-observed within a configurable window lose confidence over time
- R20. Decay formula: confidence degrades by a factor per period (e.g., CrewAI's `0.5^(age/half_life)` or simpler linear decay)
- R21. Instincts that decay below 0.3 confidence are flagged for review or auto-archived
- R22. `/evolve` checks recency before promoting — an instinct not observed in 90+ days gets a staleness warning and is skipped regardless of confidence score

## Success Criteria

- Running `kanban-work` on a manifest → completing all slices → triggers review + compound + learn automatically with zero manual intervention (except P0/P1 resolution if findings exist)
- After 3-5 kanban runs, `.atv/instincts/project.yaml` exists with real patterns extracted from the work
- Stale instincts decay and don't pollute future `/evolve` runs
- Resolved P0/P1 review findings appear as observations in the learning pipeline
- Developer sees the completion summary immediately; review runs async in background. If P0/P1s are found, developer is prompted to resolve them (interactive gate). Compound + Learn run non-blocking after resolution

## Scope Boundaries

- Do NOT rewrite the core logic of ce-compound, ce-review, /learn, /instincts, or /evolve — but extending their data model (e.g., adding `last_observed` field) and adding pre/post-invocation logic in kanban-work is in scope
- Do NOT restructure kanban-plan — it already has `expected_files` and slice decomposition
- Do NOT add per-slice micro-learnings — batch at end is the chosen pattern
- Do NOT stuff CLAUDE.md or AGENTS.md — learnings go through the instinct pipeline
- Do NOT change how review agents work internally — only wire their output into the loop
- Recency decay applies to the `/learn` skill and instincts format — not to docs/solutions/ (those are permanent institutional knowledge)

## Key Decisions

- **Review weight:** Always full (14+ agents) but async — no blocking. Evidence: ExpeL shows batch quality > per-task quality; running full review in background gives quality without developer tax
- **Compound timing:** End-of-feature only, not per-slice — Evidence: ExpeL proves cross-task batch extraction > per-task capture for related tasks
- **Learning trigger:** Fully automatic after compound — Evidence: manual-only triggers (Devin, CLAUDE.md) create maintenance burden and decay. Automatic + confidence gating (ExpeL) is the proven pattern
- **Evolution threshold:** Raised from existing 0.8 to 0.85 — tighter filter without over-restricting. Research shows observation count (5) is the stronger filter; 0.85 adds a small precision boost without meaningfully reducing recall — user decision
- **Review → learning feedback:** Only P0/P1 — Evidence: CrewAI's auto-extraction from all task output creates noise. Higher threshold = higher signal
- **Staleness:** Recency decay on instincts — Evidence: CrewAI's half-life model. Without decay, ATV would accumulate zombie instincts that block evolve or teach wrong patterns — assumption: specific decay constants need tuning in practice

## Dependencies / Assumptions

- ce-review can be invoked programmatically (not just via slash command) — assumed true based on skill architecture
- ce-compound can receive context about review findings — may need a small interface addition
- `/learn` can accept structured observations (not just "analyze recent work") — needs verification during planning
- `.atv/instincts/project.yaml` format supports a `last_observed` timestamp field — needs verification or addition
- Recency decay can be added to `/learn` without breaking existing instinct format — needs verification

## Alternatives Considered

- **Per-slice micro-learnings (rejected):** Research shows batch extraction across related tasks produces higher-quality insights. Slices in a kanban share domain, making per-slice notes repetitive — ExpeL AAAI 2024
- **Proportional review (rejected):** Complexity of auto-detecting "light" vs "full" scope adds decision logic without clear benefit. Async full review gives quality without blocking — user decision
- **Staleness check only at /evolve (rejected):** By the time instincts reach 0.8 confidence, they may already be stale from early observations. Continuous decay prevents this — CrewAI pattern
- **CLAUDE.md stuffing (rejected):** Context bloat is the #1 failure mode per research. Claude Code caps at 25KB, Windsurf at 6K chars. ATV's on-demand skill loading is architecturally superior — landscape research

## Slice Candidates (advisory for /kanban-plan)

- **Wire ce-review as mandatory async step** — kanban-work invokes ce-review automatically after all slices complete, runs in background, P1s block shipping
- **Add compound step** — after review, invoke ce-compound with full feature context, skip if nothing novel
- **Feed resolved P1s into observations** — format review findings as structured observations for /learn
- **Auto-trigger /learn** — after compound completes, run /learn automatically to extract instincts
- **Add recency decay to instincts** — modify /learn to track last_observed, add decay logic, flag stale instincts
- **Add staleness guard to /evolve** — prevent stale instincts from graduating regardless of confidence score
- **Auto-evolve on cadence** — every 5th kanban completion, silently run /evolve with 0.85 threshold + staleness check
- **Track kanban completion count** — persist a counter in `.atv/` state so evolve cadence survives across sessions

## Outstanding Questions

### Deferred to Planning

- [Affects R9][Technical] What format should observations take when fed to /learn? Does /learn accept structured input or only analyze git history?
- [Affects R5][Technical] Can ce-compound receive arbitrary context (review findings, slice list) or does it only analyze recent git activity?
- [Affects R15][Technical] Does `.atv/instincts/project.yaml` currently have a timestamp field per instinct, or does one need to be added?
- [Affects R16][Needs research] What decay constants work in practice? CrewAI uses half-life in days — what's the right half-life for a coding project?

## Next Steps

→ `/kanban-plan` for vertical-slice decomposition
