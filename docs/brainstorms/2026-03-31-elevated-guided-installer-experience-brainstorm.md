---
date: 2026-03-31
topic: elevated-guided-installer-experience
---

# Elevated Guided Installer Experience

## What We're Building

A guided installer experience that feels like a cohesive terminal product rather than three disconnected moments. The non-guided path stays fast and unchanged. The guided path becomes a clearer journey: choose the stack packs you want, choose an opinionated setup level, understand what each capability does, watch installation happen with real feedback, and land in a useful post-install launchpad.

The goal is not “more TUI for the sake of TUI.” The goal is confidence. Beginners should feel safely guided, while power users should feel informed and in control.

In practical UX terms, the new guided flow should feel like this:

1. **Stack packs first, all selected by default** — instead of asking for one primary stack, the installer shows TypeScript, Python, Rails, and General as additive packs. Detection only pre-highlights likely matches and explains why; users can unselect anything they do not want.
2. **Preset second, with clearer trade-offs** — Starter / Pro / Full remains the opinionated path, but each preset preview explains install time, network/runtime requirements, and what kinds of capabilities it unlocks.
3. **Capability browser instead of checkbox wall** — if users customize, they should see grouped categories, short descriptions, and a preview/details pane. The goal is to answer “what does this add?” before the user toggles it.
4. **Install telemetry instead of generic progress** — the install screen should show parent steps and important substeps, including skips, prerequisite warnings, and failure reasons. Users should be able to tell whether the installer is scaffolding files, cloning gstack, generating docs, or setting up browser tooling.
5. **Launchpad instead of a dead-end summary** — after install, the user should see selected stack packs, what was skipped, the first slash commands to try, and the highest-value next actions for this repo.

## Why This Approach

We explored three directions:

1. **Elevated hybrid journey** *(recommended)* — keep `huh` for short decision screens, but upgrade the guided flow with stronger structure, richer previews, real progress telemetry, and a polished results screen.
2. **Single full-screen control center** — rewrite the guided flow as one continuous Bubble Tea app with sidebars, panels, and shared state throughout.
3. **Task-centric launcher** — shift from a wizard into a command-palette/list driven interface focused on actions, previews, and reruns.

The recommended direction is the elevated hybrid journey. It compounds what already exists in the repo, respects YAGNI, and still captures the best ideas from Bubble Tea and the strongest TUIs in `awesome-tuis`: k9s-style persistent help and hotkeys, lazydocker-style unified visibility, superfile-style premium visual hierarchy, and Glow-style inline documentation preview.

The most important UX change is that **detection stops being a hidden decision**. Today the product effectively collapses a polyglot repo into one stack. In the new UX, detection becomes a recommendation layer: “we found Rails and TypeScript signals, so those packs are preselected.” That makes the installer feel smarter without taking control away from the user.

The second important change is that **customization becomes legible**. Users should not have to infer the difference between a reviewer, a setup step, a skill pack, and browser tooling from terse labels. The interface should explain the outcome of a choice before install, not after.

## Reference Patterns

- **Bubble Tea + Bubbles**: list browsing, key-help, progress, viewport, tables, pagination, and stateful screen transitions.
- **Lip Gloss**: adaptive colors, panel layouts, borders, spacing, tables, trees, and durable visual hierarchy.
- **k9s**: sticky key-hint footer, breadcrumbs/history, configurable actions, and “everything important is visible” UX.
- **lazydocker**: one-terminal-window thinking — logs, state, and actions in one place.
- **superfile**: modern theming and sectioned layouts that feel intentional instead of utilitarian.
- **Glow**: embedded markdown-style explanation panels for presets, prerequisites, and next steps.

## Key Decisions

- **Keep one-click mode untouched.** The richer experience is for `--guided`, not for every user.
- **Make stack support multi-stack by default.** Replace single-stack selection with stack-pack selection. All supported stack packs start selected, and users can unselect any they do not want.
- **Put stack-pack selection before preset selection.** Users first choose the language/framework coverage they want in the repo, then choose how much process/tooling depth they want.
- **Preserve progressive disclosure.** Presets remain the first-class path; customization becomes more legible, not more complex.
- **Replace the flat customization wall with a capability browser.** Users should see grouped categories, short descriptions, and a details/preview pane before committing.
- **Treat install progress as telemetry, not decoration.** Show prerequisites, nested substeps, skips, failures, and logs clearly enough that users understand what is happening without reading raw terminal noise.
- **End in a launchpad, not a dead-end summary.** After install, show what was installed, what was skipped, first commands, and the most useful docs/actions.
- **Adopt a stronger visual language.** Persistent footer help, better hierarchy, adaptive colors, and calmer spacing should make the installer feel trustworthy and premium.
- **Treat stack-specific assets as additive packs.** Shared root guidance should become stack-agnostic or mergeable, while language-specific reviewers and file instructions install as layered additions.

## Resolved Questions

- **Should we rewrite the wizard from scratch first?** No. The first upgrade wave should elevate the current hybrid architecture, not replace it.
- **Who is the primary user?** First-run users, with a strong secondary path for power users.
- **Where should the “wow” come from?** Clarity, previewability, and feedback — not animation alone.
- **How should polyglot repos work?** Multi-stack is the default. Users start with all stack packs selected and can deselect any they do not want.

## Open Questions

None — this is clear enough to move to planning.

## Next Steps

→ `/ce-plan` for implementation details

Related follow-on brainstorm:

- `docs/brainstorms/2026-03-31-post-install-memory-launchpad-brainstorm.md` — defines the post-install launchpad as a memory-first follow-on experience rather than leaving it implicit inside the installer brainstorm.

Related implementation plan:

- `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md` — phases the work across guided-flow contracts, richer installer UX, manifest/telemetry foundations, deterministic launchpad behavior, and an optional later Copilot SDK concierge.
