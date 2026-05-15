---
name: klfg
description: "Full kanban pipeline orchestrator. Chains /kanban-brainstorm → /kanban-plan → /kanban-work → /ce-review → /observe → /learn. Pauses only for the initial brainstorm Q&A in step 1 and any HITL-flagged slices in step 3; otherwise runs autonomously through to DONE. Use when the user says 'klfg', 'kanban lfg', 'run the full kanban', 'go from brainstorm to done', or wants the same hands-off feel as /lfg but for the kanban (vertical-slice) workflow."
argument-hint: "[feature description]"
disable-model-invocation: true
---

CRITICAL: You MUST execute every step below IN ORDER. Do NOT skip any required step. Do NOT jump ahead to coding or implementation. The brainstorm (step 1), plan (step 2), and work (step 3) phases each have a GATE that must verify their output exists before the next step begins. Violating this order produces bad output.

This pipeline is interactive in **two specific places** and autonomous everywhere else:

1. **Step 1 (brainstorm)** stops for product Q&A. That is the design — `kanban-brainstorm` does research first, then asks the user targeted product questions before producing a requirements doc.
2. **Step 3 (work)** stops only on slices the manifest flagged `hitl: true`. Those are intentional human-judgment gates set during planning. `kanban-work` handles them and resumes automatically once the user answers.

Everything else proceeds without prompting. Once the user picks "Proceed to /kanban-plan" at the end of step 1, head off till done.

## Pipeline

1. `/kanban-brainstorm $ARGUMENTS`

   GATE: STOP. Verify the brainstorm produced a requirements document at `docs/brainstorms/*-requirements.md`. If no requirements doc exists, re-run `/kanban-brainstorm $ARGUMENTS` and resume the conversation. Do NOT proceed to step 2 until a written requirements doc exists.

   **Record the requirements doc path.** Refer to it as `<reqs-path>` for the rest of the pipeline.

   Also check the requirements doc for `## Outstanding Questions` → `### Resolve Before Planning`. If that subsection has any unresolved entries, do NOT proceed — return to step 1 and resolve them first. `kanban-brainstorm` is responsible for not handing off until that section is empty, but verify here as a safety check.

2. `/kanban-plan <reqs-path>`

   GATE: STOP. Verify `/kanban-plan` produced a manifest file at `docs/plans/*-kanban-*-manifest.md` and one plan file per slice. If no manifest was created, re-run `/kanban-plan <reqs-path>`. Do NOT proceed to step 3 until both the manifest and per-slice plans exist.

   **Record the manifest path.** Refer to it as `<manifest-path>` for the rest of the pipeline.

3. `/kanban-work <manifest-path>`

   `kanban-work` will execute every pending slice in dependency order. It will pause **only** on slices marked `hitl: true` and will resume automatically once the user answers. Let it run.

   GATE: STOP. After `kanban-work` returns, re-read the manifest. Every slice must be `status: done` or `status: skipped`. If any slice is `pending`, `in_progress`, or `blocked`, re-run `/kanban-work <manifest-path>` to resume. Do NOT proceed to step 4 until the manifest is fully resolved.

   If a slice is genuinely stuck (e.g., `blocked` for an external reason), surface that to the user and stop the pipeline. Do not paper over a blocked slice.

4. `/ce-review mode:autofix plan:<reqs-path>`

   Final review pass over everything that was built across all slices. Pass the **requirements doc path** (not the manifest) so `ce-review` can verify each `R1`, `R2`, … requirement was actually delivered. Autofix mode applies only `safe_auto` fixes silently and emits residual findings as todos for the user to triage later.

   If `ce-review` returns residual P0 findings, surface them to the user before continuing — those are issues the autofix could not address on its own.

5. `/observe`

   Analyze the changed files from this session and capture patterns into the observation log. This feeds the learning pipeline in step 6.

6. `/learn`

   Extract reusable patterns from this session into project instincts (`.atv/instincts/project.yaml`). Cheap, every-session knowledge capture.

7. Output `<promise>DONE</promise>` once steps 1–6 are complete.

## Notes

- **Why no `/unslop`:** intentionally omitted from this pipeline. Run `/unslop` manually if you want a slop-removal pass.
- **Why no `/ce-compound`:** `/ce-compound` is reserved for documenting *notable* problems that were solved during the session (heavyweight case studies in `docs/solutions/`). Not every kanban run produces one. Run it manually after this pipeline if a slice involved real problem-solving worth a case study.
- **Why no `/land`:** committing, pushing, and opening a PR is a separate, deliberate act. Run `/land` after `klfg` finishes when you're ready to ship.
- **Resuming after interruption:** `klfg` is idempotent across restarts because each step's GATE checks for the produced artifact. If the session is interrupted between steps, re-invoke `klfg` with the same arguments and it will pick up at the first failing GATE.

Start with step 1 now. Remember: brainstorm FIRST, plan SECOND, work THIRD. Never skip a phase.
