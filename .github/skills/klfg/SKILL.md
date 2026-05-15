---
name: klfg
description: "Full kanban pipeline orchestrator. Chains /kanban-brainstorm → /kanban-plan → /kanban-work → DONE. kanban-work handles the full gauntlet internally: scope lock, execution, tests, diff-scope, destructive guard, QA (lint + browser), repair loop, Figma sync, ce-review, compound, learn, and evolve. Use when the user says 'klfg', 'kanban lfg', 'run the full kanban', 'go from brainstorm to done', or wants the same hands-off feel as /lfg but for the kanban (vertical-slice) workflow."
argument-hint: "[feature description]"
disable-model-invocation: true
---

CRITICAL: You MUST execute every step below IN ORDER. Do NOT skip any required step. Do NOT jump ahead to coding or implementation. The brainstorm (step 1), plan (step 2), and work (step 3) phases each have a GATE that must verify their output exists before the next step begins. Violating this order produces bad output.

This pipeline is interactive in **two specific places** and autonomous everywhere else:

1. **Step 1 (brainstorm)** stops for product Q&A. That is the design — `kanban-brainstorm` does research first, then asks the user targeted product questions before producing a requirements doc.
2. **Step 3 (work)** stops only on slices the manifest flagged `hitl: true` and when safety gates fire (scope violations, destructive commands, QA failures that exhaust repair). `kanban-work` handles them and resumes automatically once the user answers.

Everything else proceeds without prompting. Once the user picks "Proceed to /kanban-plan" at the end of step 1, hands off till done.

## Pipeline

1. `/kanban-brainstorm $ARGUMENTS`

   GATE: STOP. Verify the brainstorm produced a requirements document at `docs/brainstorms/*-requirements.md`. If no requirements doc exists, re-run `/kanban-brainstorm $ARGUMENTS` and resume the conversation. Do NOT proceed to step 2 until a written requirements doc exists.

   **Record the requirements doc path.** Refer to it as `<reqs-path>` for the rest of the pipeline.

   Also check the requirements doc for `## Outstanding Questions` → `### Resolve Before Planning`. If that subsection has any unresolved entries, do NOT proceed — return to step 1 and resolve them first. `kanban-brainstorm` is responsible for not handing off until that section is empty, but verify here as a safety check.

2. `/kanban-plan <reqs-path>`

   GATE: STOP. Verify `/kanban-plan` produced a manifest file at `docs/plans/*-kanban-*-manifest.md` and one plan file per slice. If no manifest was created, re-run `/kanban-plan <reqs-path>`. Do NOT proceed to step 3 until both the manifest and per-slice plans exist.

   **Record the manifest path.** Refer to it as `<manifest-path>` for the rest of the pipeline.

3. `/kanban-work <manifest-path>`

   `kanban-work` executes every pending slice in dependency order, running the full gauntlet per slice:

   **Per-slice gates (all mandatory):**
   - 3.0 Scope Lock — block writes outside `expected_files`
   - 3 Execute — TDD/integration/verification-only
   - 3.5 System-Wide Test Check — trace side effects
   - 3.6 Diff-Scope Verification — git diff vs declared scope
   - 3.7 Destructive Command Guard — block rm -rf, force push, etc.
   - 3.8 QA — lint (all slices) + browser checks (frontend) → kanban-repair on failure
   - 3.9 Figma Design Sync — UI slices only

   **Post-completion (after all slices):**
   - ce-review — full multi-agent code review
   - Resolution Gate — P0/P1 must be fixed before shipping
   - Compound + Learn + Evolve — document, extract patterns, promote instincts

   HITL pauses: slices flagged `hitl: true`, scope violations, destructive commands, QA failures that exhaust repair (5-iteration cap, stuck detection).

   GATE: STOP. After `kanban-work` returns, re-read the manifest. Every slice must be `status: done` or `status: skipped`. If any slice is `pending`, `in_progress`, or `blocked`, re-run `/kanban-work <manifest-path>` to resume.

   If a slice is genuinely stuck (e.g., `blocked` for an external reason), surface that to the user and stop the pipeline. Do not paper over a blocked slice.

4. Output `<promise>DONE</promise>` once steps 1–3 are complete.

## Notes

- **Why no `/unslop`:** intentionally omitted. Risk of flagging parallel agent WIP as false positives. Run manually if needed.
- **Why no separate `/ce-review`:** kanban-work runs ce-review internally at Step 5.4 with full context (diff, requirements, slice plans). A second pass would be redundant.
- **Why no separate `/learn` or `/observe`:** kanban-work feeds resolved P0/P1 findings to observations.jsonl (Step 5.5), runs `/learn` (Step 5.6), and auto-triggers `/evolve` every 5th kanban completion.
- **Why no separate `/ce-compound`:** kanban-work invokes ce-compound at Step 5.6 for features with novel patterns. Skips automatically for boilerplate.
- **Why no `/land`:** committing, pushing, and opening a PR is a separate, deliberate act. Run `/land` after `klfg` finishes when you're ready to ship.
- **Resuming after interruption:** `klfg` is idempotent across restarts because each step's GATE checks for the produced artifact. If the session is interrupted between steps, re-invoke `klfg` with the same arguments and it will pick up at the first failing GATE.

Start with step 1 now. Remember: brainstorm FIRST, plan SECOND, work THIRD. Never skip a phase.
