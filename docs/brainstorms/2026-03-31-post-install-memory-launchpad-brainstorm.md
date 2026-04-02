---
date: 2026-03-31
topic: post-install-memory-launchpad
---

# Post-Install Memory Launchpad Experience

## What We're Building

A post-install launchpad that turns ATV's installed memory into a visible, reusable product surface. It should open naturally at the end of guided install and be easy to reopen later. Instead of dropping users into a dead-end summary, the launchpad answers four questions immediately: what got installed, what memory already exists in this repo, what capabilities are active, and what slash command or workflow should happen next.

The launchpad is not a replacement for the terminal, VS Code, or the installer itself. It is an orientation layer and action hub. The home view should feel like a memory-aware dashboard with a secondary execution surface: a memory library, installed intelligence, recommended next moves, and an optional live terminal/TUI pane for when the user wants to act without context-switching.

## Why This Approach

We explored three directions:

1. **Local rule-based launchpad only** — keep the experience entirely deterministic and file-driven.
2. **Copilot-native memory concierge from day one** — make an embedded agent the main interface for understanding repo memory and deciding next steps.
3. **Phased hybrid launchpad** *(recommended)* — build a deterministic memory-first dashboard first, then add an optional GitHub Copilot SDK assistant as a layer on top.

The phased hybrid direction is the strongest fit. It preserves ATV's ability to be useful without network, auth, or model latency, while still leaving room for a much smarter assistant later. The launchpad should feel trustworthy before it feels magical. That means the first version should compute recommendations from explicit repo state and installed capabilities, then expose a clean seam where a Copilot-powered assistant can explain, summarize, and personalize those recommendations.

## Key Decisions

- **Make the launchpad a follow-on experience, not the main installer.** The guided installer remains the entry point; the launchpad becomes the living home screen after install.
- **Make memory the primary information architecture.** Brainstorms, plans, solutions, instructions, skills, agents, and runtime tools should be surfaced as visible product inventory.
- **Keep the terminal secondary.** The home screen should be a dashboard, with a terminal/TUI pane available for execution, logs, and full-screen tools.
- **Use deterministic recommendations first.** The first version should suggest actions from local facts like file presence, plan state, selected stack packs, and runtime health.
- **Do not require GitHub Copilot SDK for v1.** The launchpad should still be valuable offline or without authentication.
- **Design a seam for a Copilot assistant now.** If added later, the assistant should consume a structured memory index and typed tools, not scrape raw files ad hoc.
- **Use the assistant for explanation and navigation, not truth.** Copilot can help answer “what should I do next?” and “why is this recommended?”, but core state should come from the local indexer.
- **Persist install state explicitly.** The installer should write a small install manifest so the launchpad knows what ATV intended to install, not just what it can infer afterward.

## Resolved Questions

- **Do we need GitHub Copilot SDK in the first release?** No. It should be optional and additive, not a launch requirement.
- **Should the assistant own recommendations?** No. Recommendations should come from deterministic local rules first; the assistant can explain and refine them.
- **What is the “cool” part of the experience?** Seeing repo memory, installed intelligence, and next moves in one place — not just embedding a terminal.
- **Should the launchpad be agent-first or dashboard-first?** Dashboard-first.

## Open Questions

None — this is clear enough to move to planning.

## Next Steps

→ `/ce-plan` for a phased implementation plan covering:

- install manifest + memory index
- dashboard information architecture
- deterministic recommendation engine
- terminal integration strategy
- optional GitHub Copilot SDK phase for a memory concierge

Related implementation plan:

- `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md` — carries this brainstorm forward as Phase 2–4 work after the guided-installer foundations are in place.
