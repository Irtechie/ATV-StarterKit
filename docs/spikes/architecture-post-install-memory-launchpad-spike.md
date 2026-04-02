---
title: "Post-Install Memory Launchpad Research"
category: "Architecture & Design"
status: "🟢 Complete"
priority: "High"
timebox: "1 week"
created: 2026-03-31
updated: 2026-03-31
owner: "GitHub Copilot"
tags: ["technical-spike", "architecture", "electron", "memory", "post-install"]
---

# Post-Install Memory Launchpad Research

## Summary

**Spike Objective:** Determine how to build a live post-install Electron experience that surfaces ATV's memory system, shows what is installed, and recommends the right slash commands based on repo state and available capabilities.

**Why This Matters:** ATV already installs a differentiated memory system, but today users mainly discover it through docs and terminal output. A launchpad could turn memory from a hidden superpower into a visible product surface.

**Timebox:** 1 week

**Decision Deadline:** Before planning any desktop companion or post-install launcher work.

## Research Question(s)

**Primary Question:** What is the right architecture and UX for a post-install Electron companion that reads ATV memory artifacts locally, presents them clearly, and recommends next actions?

**Secondary Questions:**

- How should the app combine dashboard UI with terminal/TUI execution?
- Which installed artifacts are reliable sources of memory and capability state?
- What should be persisted by the installer to make the launcher accurate and cheap to build?
- How should recommendations avoid overwhelming users?

## Investigation Plan

### Research Tasks

- [x] Inventory installed memory-like artifacts in the repository and installer outputs
- [x] Research Electron architecture and security guidance
- [x] Research xterm.js and node-pty integration patterns
- [x] Research recommendation and launchpad UX patterns
- [x] Synthesize an architecture recommendation
- [x] Document findings and follow-up actions

### Success Criteria

**This spike is complete when:**

- [x] The ATV memory surfaces are classified into durable, runtime, and derived layers
- [x] A desktop architecture recommendation is documented
- [x] A recommendation engine approach is proposed
- [x] A concrete MVP scope is identified

## Technical Context

**Related Components:** `cmd/init.go`, `pkg/tui/*`, `pkg/output/printer.go`, `pkg/gstack/*`, `pkg/agentbrowser/*`, `pkg/scaffold/*`, `README.md`, `docs/brainstorms/*`, `docs/plans/*`, `docs/solutions/*`

**Dependencies:** Any future implementation depends on deciding whether the launcher is optional or bundled, and whether the installer writes a manifest/index.

**Constraints:**

- The current installer does **not** persist a canonical install manifest
- Much of the valuable state is file-based, but some runtime state lives in `~/.gstack/` and `~/.agent-browser/sessions/`
- Electron renderer security must remain strict: narrow preload API, no broad Node exposure
- Embedding a terminal requires native-module/PTTY complexity on Windows

## Research Findings

### Investigation Results

ATV already installs or enables several strong memory surfaces:

1. **Durable repo memory**
   - `docs/brainstorms/*.md`
   - `docs/plans/*.md`
   - `docs/solutions/**/*.md`
   - `compound-engineering.local.md`

2. **Installed capability memory**
   - `.github/copilot-instructions.md`
   - `.github/copilot-setup-steps.yml`
   - `.github/copilot-mcp-config.json`
   - `.github/skills/**`
   - `.github/agents/**`
   - stack/file instruction assets

3. **Runtime memory**
   - `.gstack/**` in project
   - `~/.gstack/` user-global session/config state
   - `~/.agent-browser/sessions/` user-global browser state

4. **Conceptual but not currently persisted state**
   - selected install preset
   - chosen layers/stack packs
   - install health summary
   - runtime verification status

This means the launchpad should **not** be sold as a generic shell wrapper. Its real value is: “show me what ATV knows, what ATV installed, and what I should do next.”

### UX Recommendation

The best product shape is a **dashboard-first Electron companion with an embedded terminal/TUI as a secondary pane**.

The home screen should not be a terminal. It should be a memory-aware launchpad with four primary regions:

1. **Memory Library**
   - Brainstorms
   - Plans
   - Solutions
   - Recent/active items
   - Status badges like Draft / Active / Complete / Empty

2. **Installed Intelligence**
   - Stack packs installed
   - Skills, agents, MCP servers
   - gstack runtime presence/health
   - agent-browser presence/session health

3. **Recommended Next Moves**
   - Ranked slash command suggestions
   - Each suggestion includes a rationale, e.g. “2 brainstorms have no plan” or “active plan has unchecked tasks”

4. **Terminal / Live TUI Pane**
   - Optional drawer/tab/right pane using xterm.js + node-pty
   - Used to run `atv-installer`, embedded Go TUIs, shell commands, and logs

This keeps discovery structured and lets the terminal stay a power surface instead of the main information architecture.

### Recommendation Engine Model

Recommendations should be based on **installed capabilities + memory state + runtime health**.

Suggested heuristic layers:

- **State layer**
  - no brainstorms → recommend `/ce-brainstorm`
  - brainstorm without plan → recommend `/ce-plan`
  - active plan with unchecked items → recommend `/ce-work`
  - completed work but no review → recommend `/ce-review`
  - solved work but no solution doc → recommend `/ce-compound`

- **Capability layer**
  - if gstack installed and runtime healthy → unlock `/gstack-review`, `/gstack-ship`, `/gstack-qa`
  - if agent-browser installed and browser session exists → unlock browser QA and debug actions
  - if stack packs installed → surface matching reviewers/instructions

- **Confidence layer**
  - show why a recommendation exists
  - avoid more than 3 primary recommendations at once
  - keep lower-priority suggestions collapsed under “More ideas”

### Prototype/Testing Notes

Electron architecture research strongly supports:

- `contextIsolation: true`
- `nodeIntegration: false`
- narrow `contextBridge` APIs only
- validating URLs before `shell.openExternal`
- using `protocol.handle()` for custom local protocols if needed
- using `utilityProcess` or a background worker for indexing
- using `node-pty` specifically for interactive terminal hosting
- using xterm.js addons like Fit, Search, WebLinks, Clipboard, and optional WebGL

Windows-specific risk is concentrated in:

- `node-pty` native rebuilds
- terminal lifecycle cleanup
- alternate-screen/full-screen TUI behavior
- shell/path quirks

### External Resources

- Electron official docs: context isolation, security, IPC, sandbox, protocol handling
- xterm.js docs: FitAddon, SearchAddon, WebLinksAddon, ClipboardAddon, custom link providers
- node-pty docs: Electron integration and Windows PTY behavior
- ATV README / release notes / installer code

## Decision

### Recommendation

Build the launcher as an **optional post-install Electron companion** with a **memory-first dashboard** and an **embedded PTY terminal pane**.

Do **not** make the terminal the primary UI.
Do **not** rewrite the Go installer/TUI logic in JavaScript.
Do **not** rely on runtime inference alone.

Instead:

1. Read durable repo memory and installed capability files
2. Add a small persisted install manifest/index
3. Compute recommendations locally
4. Use the terminal only for execution and live TUI hosting

### Rationale

This approach best fits ATV’s differentiator: the installer already gives users a richer memory system than a normal scaffolder. The launcher should make that visible.

A terminal-first Electron app would undersell the memory system. A dashboard-only app would lose the power of live execution. The hybrid model gets both.

### Implementation Notes

Recommended architecture:

- **Electron main**: IPC broker, file access, PTY/session manager, safe external actions
- **Background indexer**: scans repo memory and runtime state, builds derived summary
- **Preload**: typed `contextBridge` API only
- **Renderer**: dashboard UI + xterm.js pane

Recommended new persisted artifact from installer:

- `.atv/install-manifest.json` or `.github/atv-install-manifest.json`

Suggested manifest fields:

- installed preset / stack packs
- enabled layers
- gstack selected/runtime mode
- agent-browser selected
- install timestamp/version
- runtime verification results

Suggested launchpad sections:

- Memory
- Installed
- Recommended
- Terminal
- Activity/health

### Follow-up Actions

- [ ] Decide whether the Electron companion is bundled, separately installed, or Full-preset-only
- [ ] Define the install manifest schema
- [ ] Define the memory index schema and scoring logic for recommendations
- [ ] Create a brainstorm or plan specifically for the launcher UX
- [ ] Decide whether the first MVP embeds a PTY or opens an external terminal

## Status History

| Date       | Status         | Notes |
| ---------- | -------------- | ----- |
| 2026-03-31 | 🔴 Not Started | Spike created and scoped |
| 2026-03-31 | 🟡 In Progress | Research completed across repo artifacts, Electron, xterm.js, and node-pty |
| 2026-03-31 | 🟢 Complete    | Recommended a memory-first Electron launchpad with optional embedded terminal |

---

_Last updated: 2026-03-31 by GitHub Copilot_
