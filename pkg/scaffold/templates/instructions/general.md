# Project Conventions

This project uses the ATV (Agentic Tool & Workflow) Starter Kit.

## Available Workflows

- `/ce-brainstorm` — Explore what to build through collaborative dialogue
- `/ce-plan` — Create a structured implementation plan
- `/ce-work` — Execute the plan with quality checks
- `/ce-review` — Multi-agent code review
- `/ce-compound` — Document solutions for future reference
- `/lfg` — Full autonomous pipeline (plan → work → review)

## Documentation Structure

- `docs/plans/` — Implementation plans (living documents with checkboxes)
- `docs/brainstorms/` — Brainstorm documents (what to build decisions)
- `docs/solutions/` — Documented solutions (institutional knowledge)

## gstack Skills (if installed)

- `/office-hours` — YC-style forcing questions to reframe your product
- `/plan-ceo-review` — Rethink the problem; find the 10-star product
- `/plan-eng-review` — Lock architecture, data flow, edge cases
- `/review` — Staff-level code review; auto-fix obvious issues
- `/qa` — Test app in real browser, find and fix bugs (requires Bun)
- `/ship` — Sync main, run tests, push, open PR
- `/cso` — OWASP Top 10 + STRIDE threat model
- `/careful` — Warn before destructive commands
- `/investigate` — Systematic root-cause debugging
- `/retro` — Weekly retrospective with trends

## ATV Override Rules

When both ATV and gstack provide similar functionality, ATV takes priority:

- **Design docs**: Write to `docs/brainstorms/` (ATV), not `DESIGN.md` (gstack)
- **Solutions**: Document via `/ce-compound` into `docs/solutions/` (ATV), not gstack's `/retro`
- **Plans**: Use `docs/plans/` with ATV naming (`YYYY-MM-DD-NNN-type-name-plan.md`)
- **Reviews**: ATV's `/ce-review` agent selection governs; gstack's `/review` runs alongside
- **Protected artifacts**: Never flag `docs/plans/`, `docs/solutions/`, `docs/brainstorms/`, `compound-engineering.local.md`, or `.github/skills/gstack/` for deletion

## Coding Conventions

- Follow existing patterns in the codebase
- Write tests for new functionality
- Use conventional commit messages (`feat:`, `fix:`, `refactor:`)
