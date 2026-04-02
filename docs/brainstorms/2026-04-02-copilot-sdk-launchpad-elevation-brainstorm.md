---
date: 2026-04-02
topic: github-copilot-sdk-launchpad-elevation
---

# GitHub Copilot SDK Launchpad Elevation

## What We're Building

An elevated launchpad experience that transforms the current static post-install dashboard into a live, realtime monitoring engine powered by the GitHub Copilot SDK. The launchpad becomes a k9s-style persistent terminal dashboard that watches five signal layers in realtime:

1. **Repo memory artifacts** — brainstorms, plans, solutions as they're created/updated/deleted
2. **Copilot context & memory state** — context window usage, token consumption, prompt sizes, memory file mutations
3. **Active skill & agent usage** — skills loaded, agents invoked, tools used in the current session. Observable from VS Code extension API on the panel surface; approximated from `.github/` file presence on the TUI surface.
4. **Install health & drift** — what got installed, what's stale, what needs updating, config drift from catalog
5. **Runtime tool health** — gstack availability, agent-browser status, MCP server connectivity

The SDK adds an intelligence layer on top of local monitoring: context-aware recommendations that explain WHY something is suggested, proactive alerts when drift or staleness is detected, and a suggest-then-execute flow where the SDK proposes actions and the user approves them from the dashboard.

This **replaces** the previous "Phase 4 optional concierge" concept from the March 31 plan. The old plan treated the SDK as a chat assistant bolted onto a static dashboard. This new vision makes the SDK the intelligence engine powering a live monitoring surface.

## Why This Approach

We explored three directions:

1. **SDK-Powered Monitoring Daemon** *(selected)* — Go binary runs a persistent background watcher using `fsnotify` for local state changes and the Copilot SDK Go client for intelligence. The Bubble Tea TUI connects to this daemon and renders live-updating panels.

2. **On-Demand SDK Analysis** — Launchpad scans state at launch and on user-triggered refresh. SDK calls happen lazily. Simpler but not truly realtime and can't do proactive alerts.

3. **VS Code Extension-First** — Build the monitoring engine as a VS Code extension with a webview panel. Richer visualization but splits the codebase and loses terminal-first users.

The daemon approach was chosen because:
- It enables genuine realtime monitoring with proactive alerts
- Clean architectural separation: local watcher handles fast filesystem events, SDK provides intelligence on top
- The Go SDK (`github.com/github/copilot-sdk/go`) integrates natively — no subprocess spawning or npm dependency
- Hybrid offline/online model: monitoring works without auth, SDK intelligence activates when available
- Terminal-first philosophy aligns with the existing ATV ethos (k9s, lazydocker, superfile as reference UX)
- Foundation for suggest-and-execute — daemon can queue and run approved actions

## Key Decisions

- **Daemon architecture, not scan-on-demand:** A persistent background process watches filesystem events and maintains live state. The TUI is a view layer that connects to this state. This enables proactive alerts and realtime visualization that scan-on-demand cannot provide.

- **Five monitoring layers as the information architecture:** Repo memory, Copilot context/memory, active skills/agents, install health, and runtime tools. Each layer maps to a launchpad panel or view. This replaces the old five-tab layout (Overview/Copilot/CE/Gstack/Moves) with signal-oriented panels.

- **Hybrid offline/online model:** Core monitoring (filesystem watching, install health, runtime probes) works fully offline with no auth. SDK intelligence (context-aware WHY explanations, proactive staleness alerts, suggest-and-execute) requires Copilot CLI auth and network. Degraded mode is clearly communicated in the TUI.

- **SDK as intelligence engine, not chat assistant:** The SDK does not appear as a chat pane. It powers the recommendation engine — analyzing repo state, computing context-aware suggestions, explaining rationale, and proposing executable actions. The user interacts with the dashboard, not a conversation.

- **Suggest-then-execute flow:** When the SDK recommends an action (e.g., "Your TypeScript plan from March 29 is stale — rerun `/ce-plan`"), the user can approve it from the dashboard and the daemon executes it. No silent execution.

- **Go SDK native integration:** Use `github.com/github/copilot-sdk/go` directly in the Go binary. Define typed tools that expose the local state model to the SDK: `getMemoryIndex`, `listRecommendations`, `explainRecommendation`, `getInstallHealth`, `getRuntimeStatus`, `proposeAction`.

- **fsnotify for local state watching:** Watch `docs/brainstorms/`, `docs/plans/`, `docs/solutions/`, `.github/`, `~/.gstack/`, and manifest files. Debounce events and update an in-memory state model that the TUI renders.

- **Dual surface from day one:** Two view layers ship together over a shared state/intelligence backend:
  - **Terminal TUI** (Bubble Tea) — for Copilot CLI users and terminal workflows. Observes layers 1, 2, 4, 5 fully; approximates layer 3 from filesystem.
  - **VS Code panel** (webview) — for VS Code users. Observes all five layers including live skill/agent usage via the extension API.
  Both surfaces read from the same state manager and SDK intelligence layer.

- **Supersedes old Phase 4:** This design replaces the "optional concierge" concept from `docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md` Phase 4. That plan treated the SDK as a chat assistant bolted onto a static dashboard. This makes it the monitoring intelligence engine.

## Architecture Overview

```
┌─────────────────────────────┐  ┌─────────────────────────────┐
│    Bubble Tea TUI (CLI)     │  │   VS Code Webview (Panel)   │
│  Memory|Context|Health|Moves│  │  Memory|Context|Health|Moves│
│  [layers 1,2,4,5 + ~3]     │  │  [all 5 layers, live]       │
└─────────────┬───────────────┘  └──────────────┬──────────────┘
              │ reads live state                 │
              ▼                                  ▼
┌────────────────────────────────────────────────────────────────┐
│                  State Manager (shared backend)                │
│         RepoMemory | Context | Skills | Health | Recs          │
├────────────────────────┬───────────────────────────────────────┤
│    FS Watcher          │         Copilot SDK                  │
│    (fsnotify)          │         (Go client)                  │
│    - repo files        │         - recommendations            │
│    - manifest          │         - explanations               │
│    - runtime probes    │         - proposed actions            │
│    [OFFLINE]           │         [ONLINE]                     │
├────────────────────────┼───────────────────────────────────────┤
│                        │    VS Code Extension API             │
│                        │    - active skills/agents            │
│                        │    - editor context                  │
│                        │    [VS CODE SURFACE ONLY]            │
└────────────────────────┴───────────────────────────────────────┘
```

## SDK Tools (Typed Go Functions)

The Copilot SDK session will be configured with these custom tools so the model can query local state:

| Tool | Description |
|------|-------------|
| `getMemoryIndex` | Returns list of repo memory artifacts (brainstorms, plans, solutions) with metadata |
| `getInstallManifest` | Returns what was installed, when, which catalog version |
| `getInstallDrift` | Compares installed state against current catalog, returns diffs |
| `getRuntimeHealth` | Probes gstack, agent-browser, MCP servers for availability |
| `getActiveContext` | Returns current Copilot context/memory state if observable |
| `listRecommendations` | Returns deterministic recommendation list from local rules |
| `proposeAction` | SDK proposes a specific action (command, file edit) for user approval |

## Resolved Questions

- **Context/memory observability:** SDK-first with approximate fallback. Try to query Copilot SDK for its own context/memory state first. If the API doesn't expose token counts or context usage, fall back to approximating from local signals (file sizes, instruction counts, skill catalog size, known context limits).

- **Daemon lifecycle:** Explicit launch via `atv launchpad`. No auto-start. The user runs the command, which starts both the daemon and the TUI together. Clean and predictable — no surprise background processes.

- **Action execution scope:** Full execution scope. The dashboard can suggest and execute (with user approval) slash commands, installer operations (re-init, gstack sync), file operations, and git commands. Start with all capabilities available and constrain later if needed, rather than artificially limiting and frustrating users.

## Next Steps

→ `/ce-plan` for phased implementation covering:
- Shared state manager backend (consumed by both surfaces)
- fsnotify integration for the five monitoring layers
- Copilot SDK Go client integration with typed tools
- Bubble Tea TUI redesign for live-updating panels
- VS Code webview panel and extension scaffolding
- VS Code extension API integration for layer 3 (active skills/agents)
- Hybrid offline/online behavior
- Suggest-then-execute action system
