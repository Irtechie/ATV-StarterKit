# Agent Instructions

<!-- KB-WORKFLOW-INSTRUCTIONS:START -->
For KB workflow requests, start with `kb-route`.

Fresh-session preflight:
- If `todo.md` or `docs/context/PROJECT.md` is missing, run `kb-map-bootstrap` immediately.
- If context or handoff folders are partially missing, run `kb-map refresh`.
- Do not ask for confirmation unless a non-empty user file would be overwritten.

Token rule: every token must pay rent. Be concise by default; keep exact paths, commands, errors, decisions, risks, and safety warnings.

Project memory files:
- `todo.md` holds active work, blockers, parked work, and handoff pointers.
- `todo-done.md` holds completed-work summaries.
- `docs/context/PROJECT.md` is the project route map.
- `docs/handoffs/active/`, `docs/handoffs/parked/`, and `docs/handoffs/done/` hold resumable handoffs.

Skills live under `.github/skills/`; these memory files are not skills.
<!-- KB-WORKFLOW-INSTRUCTIONS:END -->

