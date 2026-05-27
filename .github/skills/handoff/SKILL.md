---
name: handoff
description: Create a repo-local handoff restart packet for another agent or fresh session to pick up. Prefer kb-handoff for KB workflows.
argument-hint: "What will the next session be used for?"
---

Write a compact handoff document so a fresh agent can continue the work without relying on chat history.

Save inside the active project root:

```text
docs/handoffs/active/YYYY-MM-DD-<short-topic>.md
```

Resolve the project root first with `git rev-parse --show-toplevel`. If no valid project root exists, ask the user to change into the project directory or provide the project path. Do not write handoffs to `C:\Users\marowe\.copilot\handoffs`, home folders, drive roots, or global skill folders.

Credit: Based on mattpocock/skills (MIT). Adapted for ATV/KB repo-local workflows.

## What to include

1. **Context** — repo root, branch, and what the user was doing.
2. **Decisions made** — key architectural or design decisions, with reasoning.
3. **Work completed** — files changed, commits made, what's done.
4. **Work remaining** — what's left, blockers, next steps.
5. **Key files** — repo-local paths the next session should read first.
6. **Suggested route** — `kb-start`, `kb-plan`, `kb-work`, `kb-fix`, or another specific skill.

## Rules

- Do not duplicate content already in plans, brainstorms, PRDs, committed docs, or `todo.md`; reference by path.
- Redact API keys, passwords, tokens, PII, and private credentials.
- Keep it under 1200 words unless the user asks for more.
- If the user passed arguments, treat them as what the next session should focus on.
- Update `todo.md` with a compact handoff pointer when the handoff represents active or blocked work.
