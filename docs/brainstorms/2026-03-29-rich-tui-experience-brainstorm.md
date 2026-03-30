---
date: 2026-03-29
topic: rich-tui-experience
---

# Rich TUI Experience for ATV Installer

## What We're Building

A redesigned guided installer experience with three key improvements:

1. **Multi-step wizard with presets** — Three screens: Stack → Preset (Starter/Pro/Full) → Customize. Beginners pick a preset and go; power users tweak individual skills in a customize step.

2. **Rich install progress** — Step-by-step progress indicators during the install phase (ATV scaffold, gstack clone, bun install, doc generation, agent-browser, Chrome download), each with clear status (pending → running → done/failed).

3. **Hybrid huh + Bubbletea** — Keep `charmbracelet/huh` for the wizard form screens (it's already built on Bubbletea and handles forms well). Use raw Bubbletea for the install progress phase with animated spinners and real-time status updates.

## Why This Approach

The current TUI dumps 43+ skills into a flat checkbox list. Users don't understand what they're selecting, there's no visual hierarchy, and there's no guidance for beginners. The install phase prints text line-by-line with no progress feedback.

**Multi-step with presets** solves the beginner problem — pick "Starter" and you're done. **Customize** solves the power user problem — drill into categories and toggle individual skills. **Progress indicators** solve the install feedback problem — users know what's happening during the 30+ second gstack clone.

## Key Decisions

1. **Three presets**: Starter (core ATV only, 13 skills, instant), Pro (core + gstack sprint skills, no browser QA), Full (everything: all 45 skills, gstack runtime, agent-browser + Chrome)
2. **Preset → customize flow**: After selecting a preset, show a "Customize?" confirm. If yes, show the category-grouped multi-select pre-filled with the preset's selections. If no, proceed to install.
3. **Huh for forms, Bubbletea for progress**: Don't rewrite forms from scratch. Use raw Bubbletea `tea.Program` only for the install phase where we need animated spinners and concurrent status updates.
4. **Category descriptions in the customize step**: Each category group gets a 1-line description so users understand what they're toggling.

## Open Questions

None — all resolved during brainstorm.

## Next Steps

→ `/ce-plan` for implementation details
