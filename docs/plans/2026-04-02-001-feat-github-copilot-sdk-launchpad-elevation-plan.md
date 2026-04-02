---
title: "feat: GitHub Copilot SDK Launchpad Elevation"
type: feat
status: active
date: 2026-04-02
origin: docs/brainstorms/2026-04-02-copilot-sdk-launchpad-elevation-brainstorm.md
---

# feat: GitHub Copilot SDK Launchpad Elevation

## Overview

Transform the ATV installer launchpad from a static post-install dashboard into a realtime monitoring engine powered by the GitHub Copilot SDK. The launchpad becomes a k9s-style live-updating terminal dashboard backed by filesystem watching (`fsnotify`) and an SDK intelligence layer, with a companion VS Code webview panel sharing the same state backend.

This supersedes Phase 4 ("optional concierge") from [docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md](../plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md). The old design treated the SDK as a chat assistant; this design makes it the intelligence engine behind a live monitoring surface (see brainstorm: [docs/brainstorms/2026-04-02-copilot-sdk-launchpad-elevation-brainstorm.md](../brainstorms/2026-04-02-copilot-sdk-launchpad-elevation-brainstorm.md)).

## Problem Statement

The current launchpad (`atv launchpad`) is a polling-based dashboard that refreshes every 3 seconds via `tickRefresh()`. Its recommendations are deterministic heuristics (priority-based, capped at 5) that tell users WHAT to do but never WHY. Users see information across five tabs (Overview, Copilot, CE, Gstack, Moves) but cannot interact meaningfully — the dashboard is read-only.

**Key gaps:**
- **No realtime monitoring** — polling every 3 seconds misses rapid changes and can't deliver proactive alerts
- **No context-aware intelligence** — recommendations are rule-based; they can't explain rationale, adapt to context, or propose nuanced actions
- **No actionability** — users see next steps but must leave the dashboard to execute them
- **No continuous awareness** — closing the TUI means losing awareness of repo state changes
- **Single surface** — terminal-only; no VS Code-native experience

## Proposed Solution

A three-layer architecture:

1. **Local Monitoring Layer** (offline, always-on) — `fsnotify` watches repo files and updates an in-memory state model. No auth required.
2. **SDK Intelligence Layer** (online, additive) — GitHub Copilot SDK Go client queries the state model via typed tools, generates context-aware recommendations with explanations, and proposes executable actions.
3. **View Layer** (dual surface) — Bubble Tea TUI and VS Code webview panel both read from the shared state manager.

(see brainstorm: Architecture Overview diagram)

## Technical Approach

### Architecture

```
┌─────────────────────────────┐  ┌─────────────────────────────┐
│    Bubble Tea TUI (CLI)     │  │   VS Code Webview (Panel)   │
│  Memory|Context|Health|Moves│  │  Memory|Context|Health|Moves│
│  [layers 1,2,4,5 + ~3]     │  │  [all 5 layers, live]       │
└─────────────┬───────────────┘  └──────────────┬──────────────┘
              │ in-process reads                 │ polls state file
              ▼                                  ▼
┌────────────────────────────────────────────────────────────────┐
│                  State Manager (in-memory Go)                  │
│         RepoMemory | Context | Skills | Health | Recs          │
│         ├─ writes → .atv/launchpad-state.json (for VS Code)   │
├────────────────────────┬───────────────────────────────────────┤
│    FS Watcher          │         GitHub Copilot SDK            │
│    (fsnotify)          │         (Go client)                   │
│    - repo files        │         - recommendations             │
│    - manifest          │         - explanations                │
│    - runtime probes    │         - proposed actions             │
│    [OFFLINE]           │         [ONLINE]                      │
├────────────────────────┼───────────────────────────────────────┤
│                        │    VS Code Extension API              │
│                        │    - active skills/agents             │
│                        │    - editor context                   │
│                        │    [VS CODE SURFACE ONLY]             │
└────────────────────────┴───────────────────────────────────────┘
```

**Critical architectural decision: In-process for Phase 1** (see SpecFlow gap #1).

The daemon and TUI run in the same Go process. The TUI reads state directly from the in-memory state manager with zero serialization overhead. The VS Code extension reads state via a serialized `.atv/launchpad-state.json` file that the state manager writes on every change. This avoids the complexity of TCP/socket IPC while still enabling both surfaces.

**Rationale:** Out-of-process daemon adds IPC protocol design, authentication, port management, PID file handling, and cross-platform socket differences — none of which deliver user-facing value vs. in-process + file-sharing. If IPC becomes necessary later (e.g., for headless daemon mode), it can be extracted from the file-based approach.

### Implementation Phases

#### Phase 1: State Manager & FS Watcher Foundation

**Objective:** Replace the 3-second polling loop with event-driven filesystem monitoring. Build the shared state model that both surfaces will consume.

**New package: `pkg/monitor/`**

##### `pkg/monitor/state.go` — State Manager

```go
// LiveState extends LaunchpadSnapshot with realtime monitoring data.
type LiveState struct {
    installstate.LaunchpadSnapshot

    // Layer 1: Repo memory (live-tracked)
    Brainstorms  []ArtifactEntry  // name, path, modTime, size
    Plans        []ArtifactEntry
    Solutions    []ArtifactEntry

    // Layer 2: Context approximation
    ContextEstimate ContextEstimate  // totalInstructionBytes, skillCount, agentCount, estimatedTokens

    // Layer 3: Skills/agents (filesystem approximation for TUI)
    InstalledSkillNames []string
    InstalledAgentNames []string

    // Layer 4: Install health
    DriftEntries []DriftEntry  // path, status (missing|stale|modified), catalogHash, diskHash

    // Layer 5: Runtime health
    RuntimeHealth RuntimeHealthState  // gstack, agentBrowser, mcpServers with status + lastCheck

    // SDK intelligence (nil when offline)
    SDKRecommendations []SDKRecommendation  // title, reason, explanation (WHY), proposedAction
    SDKOnline          bool
    SDKLastQuery       time.Time

    // Metadata
    LastFSEvent  time.Time
    WatchedPaths []string
}

type ArtifactEntry struct {
    Name    string
    Path    string
    ModTime time.Time
    Size    int64
}

type DriftEntry struct {
    Path        string
    Status      DriftStatus  // Missing, Stale, UserModified, Extra
    CatalogHash string
    DiskHash    string
    InstallHash string       // hash at install time, from manifest
}

type ContextEstimate struct {
    TotalInstructionBytes int64
    SkillCount            int
    AgentCount            int
    PromptCount           int
    MCPServerCount        int
    EstimatedTokens       int  // rough approximation: bytes / 4
}

type SDKRecommendation struct {
    ID              string
    Title           string
    Reason          string
    Explanation     string           // WHY — from SDK analysis
    ProposedAction  *ProposedAction  // nil if informational only
    Priority        int
}

type ProposedAction struct {
    Type        ActionType  // SlashCommand, InstallerOp, FileOp, GitOp
    Command     string
    Args        []string
    Description string
    RiskLevel   ActionRisk  // Safe, Caution, Destructive
}
```

##### `pkg/monitor/watcher.go` — FS Watcher

```go
type Watcher struct {
    fsWatcher    *fsnotify.Watcher
    state        *LiveState
    mu           sync.RWMutex
    debounceWin  time.Duration          // 300ms default
    onChange     func(LiveState)         // notify TUI
    ignorePats   []string               // *.swp, *~, .tmp, 4913, .git/
}

func NewWatcher(root string, opts WatcherOptions) (*Watcher, error)
func (w *Watcher) Start(ctx context.Context) error   // begin watching
func (w *Watcher) Stop() error                        // cleanup
func (w *Watcher) State() LiveState                   // threadsafe snapshot read
func (w *Watcher) ForceRefresh() error                // manual re-scan
```

**Watch targets** (only paths that exist at startup; re-scan periodically for new dirs):
- `docs/brainstorms/`, `docs/plans/`, `docs/solutions/` — layer 1
- `.github/copilot-instructions.md`, `.github/skills/`, `.github/agents/`, `.github/copilot-mcp-config.json` — layers 2, 3
- `.atv/install-manifest.json` — layer 4
- `~/.gstack/`, `~/.agent-browser/sessions/` — layer 5 (user-global, watch if exists)

**Debounce strategy:**
- 300ms debounce window (matches typical editor save debounce)
- Filter noise: ignore patterns `*.swp`, `*~`, `*.tmp`, `4913` (Vim), `.DS_Store`, `.git/**`
- On debounce fire: re-scan affected layer only (not full state rebuild)
- Atomic save detection: treat RENAME events as authoritative (ignore preceding CREATE/DELETE of temp files)

**Cross-platform notes:**
- Windows: `ReadDirectoryChangesW` backend — works but doesn't support recursive watching natively; add subdirectories explicitly
- macOS: `kqueue` backend — one FD per file; limit watch targets to minimize FD usage
- Linux: `inotify` — works well; supports recursive via adding subdirectories manually

**Graceful startup without `atv init`:** If watched directories don't exist, skip them silently. Run a parent-directory watcher on `docs/` and `.github/` to detect directory creation after first `atv init`.

##### Tasks

- [ ] Create `pkg/monitor/state.go` with `LiveState`, `ArtifactEntry`, `DriftEntry`, `ContextEstimate`, `SDKRecommendation`, `ProposedAction` types
- [ ] Create `pkg/monitor/watcher.go` with `Watcher` struct, `Start()`, `Stop()`, `State()`, `ForceRefresh()` methods
- [ ] Implement filesystem event filtering (editor temp files, dotfiles, recursive dir addition)
- [ ] Implement 300ms debounce with per-layer targeted re-scan
- [ ] Implement `.atv/launchpad-state.json` serialization on every state change (for VS Code extension)
- [ ] Add `watcher_test.go` — test debounce, filtering, state updates, missing-dir handling
- [ ] Add dependency: `go get github.com/fsnotify/fsnotify`

##### Acceptance Criteria

- [ ] FS watcher detects file creates/modifies/deletes in all layer directories within 300ms
- [ ] Editor temp files (`.swp`, `~`, `.tmp`, `4913`) are filtered out
- [ ] State updates are threadsafe (concurrent reads don't race with writes)
- [ ] Watcher starts gracefully when watched directories don't exist
- [ ] `.atv/launchpad-state.json` is written atomically (temp + rename, matching existing `WriteManifest` pattern)
- [ ] State file is valid JSON and schema-versioned
- [ ] Tests pass on Windows, macOS, Linux (CI matrix)

---

#### Phase 2: Install Drift Detection

**Objective:** Extend the install manifest with per-file checksums and implement a drift detection algorithm that distinguishes catalog updates from user customizations.

##### Manifest Extension: `pkg/installstate/manifest.go`

```go
// Add to InstallManifest
type InstallManifest struct {
    // ... existing fields ...
    CatalogVersion  string            `json:"catalogVersion,omitempty"`  // hash of BuildCatalog output
    FileChecksums   map[string]string `json:"fileChecksums,omitempty"`   // path → SHA256 at install time
}
```

##### Drift Algorithm: `pkg/monitor/drift.go`

```go
func ComputeDrift(manifest InstallManifest, catalog []scaffold.Component, root string) []DriftEntry
```

**Three-way comparison:**
1. **File on disk** matches **install-time checksum** AND matches **catalog checksum** → No drift (installed, up-to-date)
2. **File on disk** matches **install-time checksum** but NOT **catalog checksum** → **Stale** (catalog updated since install)
3. **File on disk** does NOT match **install-time checksum** → **UserModified** (user customized, not drift)
4. **File in catalog** but NOT on disk and NOT in manifest → **Missing** (never installed, or deleted)
5. **File on disk** but NOT in catalog → **Extra** (user-created, informational only)

This prevents false-positive drift alerts on files users intentionally modify (e.g., `copilot-instructions.md`). (see SpecFlow gap #7)

##### `.atv-drift-ignore` File

Users can create `.atv/drift-ignore` to explicitly suppress drift for specific paths:

```
# Paths to ignore in drift detection (glob patterns)
.github/copilot-instructions.md
.github/skills/my-custom-skill/**
```

##### Tasks

- [ ] Extend `InstallManifest` with `CatalogVersion` and `FileChecksums` fields
- [ ] Update `WriteManifest()` to compute and store SHA256 checksums for all scaffolded files
- [ ] Compute catalog version hash from `BuildCatalog()` output
- [ ] Create `pkg/monitor/drift.go` with `ComputeDrift()` implementing three-way comparison
- [ ] Support `.atv/drift-ignore` glob patterns
- [ ] Integrate drift computation into `LiveState` updates (recompute on Layer 4 FS events)
- [ ] Add `drift_test.go` — test all five drift statuses, ignore patterns, missing manifest

##### Acceptance Criteria

- [ ] Manifest stores SHA256 checksums for every scaffolded file
- [ ] User-modified files are classified as `UserModified`, not `Stale`
- [ ] Catalog updates (new version of a template) correctly surface as `Stale` for unmodified files
- [ ] `.atv/drift-ignore` suppresses drift for matched paths
- [ ] Backward compatible: manifests without checksums degrade to simple presence check

---

#### Phase 3: TUI Redesign — Signal-Oriented Panels

**Objective:** Replace the five-tab layout (Overview/Copilot/CE/Gstack/Moves) with four signal-oriented panels backed by the live state manager. Tab redesign is a breaking change from the current UX — the new panels are organized by monitoring signal rather than feature area.

##### New Tab Layout

| Tab | Key | Content |
|-----|-----|---------|
| **Memory** | `1` | Repo memory artifacts: brainstorms, plans, solutions with timestamps and status badges (Draft/Active/Complete). File counts, recent activity. |
| **Context** | `2` | Context estimate: total instruction bytes, skill count, agent count, prompt count, estimated tokens. SDK-powered insights when online. |
| **Health** | `3` | Install drift entries, runtime health (gstack, agent-browser, MCP servers), install outcome summary. Color-coded status indicators. |
| **Moves** | `4` | Recommendations. Local deterministic rules (offline) enriched with SDK explanations (online). Suggest-then-execute UI for approved actions. |

##### TUI Model Update: `pkg/tui/launchpad.go`

```go
type LaunchpadModel struct {
    root     string
    watcher  *monitor.Watcher    // replaces snapshot polling
    tab      LaunchpadTab        // Memory, Context, Health, Moves (4 tabs, not 5)

    // Suggest-then-execute state
    selectedRec   int            // currently highlighted recommendation
    pendingAction *monitor.ProposedAction  // action awaiting approval
    actionResult  *ActionResult  // result of last executed action

    // SDK status
    sdkOnline    bool
    sdkLastQuery time.Time
}
```

**Key changes from current TUI:**
- Replace `snapshot installstate.LaunchpadSnapshot` with `watcher *monitor.Watcher`
- Remove 3-second `tickRefresh()` timer — state updates arrive via watcher `onChange` callback as `tea.Msg`
- Reduce from 5 tabs to 4 tabs (tab keys `1-4` instead of `1-5`)
- Add approve/reject keybindings for suggest-then-execute (`Enter` = approve, `Esc` = dismiss on Moves tab)
- Add `o` key to toggle online/offline SDK status display

**Keyboard Map:**

| Key | Action |
|-----|--------|
| `1-4` | Jump to tab |
| `Tab` / `Shift+Tab` | Next / previous tab |
| `↑/↓` or `j/k` | Navigate items in current panel |
| `r` | Force refresh all layers |
| `Enter` | Approve proposed action (Moves tab) |
| `Esc` | Dismiss proposed action / deselect |
| `o` | Toggle SDK online status display |
| `q` | Quit |

##### Online/Offline Status Bar

```
┌ ATV Launchpad ─────────────────────────── 🟢 SDK Online ─┐
│  [1] Memory  [2] Context  [3] Health  [4] Moves           │
├────────────────────────────────────────────────────────────┤
│                                                            │
│  ... panel content ...                                     │
│                                                            │
├────────────────────────────────────────────────────────────┤
│ Last FS event: 2s ago │ Watched: 42 paths │ q:quit r:refresh│
└────────────────────────────────────────────────────────────┘
```

When SDK is offline:
```
┌ ATV Launchpad ─────────────────── ⚪ Offline (local only) ─┐
```

##### Tasks

- [ ] Refactor `LaunchpadModel` to consume `*monitor.Watcher` instead of polling `LaunchpadSnapshot`
- [ ] Replace `tickRefresh()` with watcher callback → `tea.Msg` bridge
- [ ] Implement four new render methods: `renderMemory()`, `renderContext()`, `renderHealth()`, `renderMoves()`
- [ ] Implement online/offline status bar with `🟢`/`⚪` indicator
- [ ] Implement keyboard navigation for items within panels (`j/k` or `↑/↓`)
- [ ] Implement approve/reject flow for proposed actions on Moves tab
- [ ] Add footer bar with FS event timestamp, watched path count, keybindings
- [ ] Update `cmd/launchpad.go` to initialize `Watcher` and pass to TUI
- [ ] Add `launchpad_test.go` updates — test new tab rendering, watcher integration

##### Acceptance Criteria

- [ ] TUI renders four tabs: Memory, Context, Health, Moves
- [ ] Dashboard updates within 300ms of filesystem changes (no visible polling lag)
- [ ] Online/offline indicator accurately reflects SDK auth state
- [ ] Moves tab shows SDK explanations (WHY) when online, falls back to local recommendations when offline
- [ ] Approve/reject action flow works without crashing on malformed SDK responses
- [ ] Footer shows last FS event timestamp and watched path count

---

#### Phase 4: GitHub Copilot SDK Intelligence Layer

**Objective:** Integrate the GitHub Copilot SDK Go client as the intelligence engine, defining typed tools that expose local state to the model and receiving context-aware recommendations.

##### Prerequisites

- GitHub Copilot CLI installed and authenticated (`copilot --version`)
- Network connectivity for SDK queries

##### SDK Integration: `pkg/sdk/intelligence.go`

```go
package sdk

import (
    copilot "github.com/github/copilot-sdk/go"
    "github.com/github/atv-installer/pkg/monitor"
)

type Intelligence struct {
    client  *copilot.Client
    session *copilot.Session
    state   *monitor.Watcher
    online  bool

    // Rate limiting
    lastQuery    time.Time
    minInterval  time.Duration  // 30s default
    queryBudget  int            // max queries per hour
    queriesUsed  int
}

func NewIntelligence(watcher *monitor.Watcher, opts IntelligenceOptions) *Intelligence
func (i *Intelligence) Start(ctx context.Context) error   // init SDK client, attempt auth
func (i *Intelligence) Stop() error                        // cleanup session + client
func (i *Intelligence) IsOnline() bool
func (i *Intelligence) Query(ctx context.Context) ([]monitor.SDKRecommendation, error)  // rate-limited
func (i *Intelligence) Explain(ctx context.Context, recID string) (string, error)       // on-demand WHY
```

##### Tool Definitions

Seven typed tools exposed to the SDK session so the model can query local state:

```go
// Tool 1: getMemoryIndex
type MemoryIndexParams struct{}
type MemoryIndexResult struct {
    Brainstorms []monitor.ArtifactEntry `json:"brainstorms"`
    Plans       []monitor.ArtifactEntry `json:"plans"`
    Solutions   []monitor.ArtifactEntry `json:"solutions"`
}
var getMemoryIndex = copilot.DefineTool(
    "getMemoryIndex",
    "List all repo memory artifacts (brainstorms, plans, solutions) with names, paths, and modification times",
    func(params MemoryIndexParams, inv copilot.ToolInvocation) (MemoryIndexResult, error) { ... },
)

// Tool 2: getInstallManifest
type ManifestParams struct{}
var getInstallManifest = copilot.DefineTool(
    "getInstallManifest",
    "Get the current install manifest: what was installed, when, which catalog version, and outcomes",
    func(params ManifestParams, inv copilot.ToolInvocation) (installstate.InstallManifest, error) { ... },
)

// Tool 3: getInstallDrift
type DriftParams struct{}
type DriftResult struct {
    Entries []monitor.DriftEntry `json:"entries"`
    Summary string               `json:"summary"` // e.g. "3 stale, 1 missing, 12 up-to-date"
}
var getInstallDrift = copilot.DefineTool(
    "getInstallDrift",
    "Compare installed files against current catalog and return drift entries with status (stale/missing/user-modified)",
    func(params DriftParams, inv copilot.ToolInvocation) (DriftResult, error) { ... },
)

// Tool 4: getRuntimeHealth
type RuntimeHealthParams struct{}
var getRuntimeHealth = copilot.DefineTool(
    "getRuntimeHealth",
    "Probe gstack, agent-browser, and MCP server availability and return health status",
    func(params RuntimeHealthParams, inv copilot.ToolInvocation) (monitor.RuntimeHealthState, error) { ... },
)

// Tool 5: getActiveContext
type ContextParams struct{}
var getActiveContext = copilot.DefineTool(
    "getActiveContext",
    "Get context estimate: total instruction bytes, skill/agent/prompt counts, estimated token usage",
    func(params ContextParams, inv copilot.ToolInvocation) (monitor.ContextEstimate, error) { ... },
)

// Tool 6: listRecommendations
type RecsParams struct{}
var listRecommendations = copilot.DefineTool(
    "listRecommendations",
    "List current deterministic recommendations from local heuristic rules with priorities",
    func(params RecsParams, inv copilot.ToolInvocation) ([]installstate.Recommendation, error) { ... },
)

// Tool 7: proposeAction
type ProposeActionParams struct {
    Title       string `json:"title" jsonschema:"Short action title"`
    Reason      string `json:"reason" jsonschema:"Why this action is recommended"`
    Command     string `json:"command" jsonschema:"The command to execute"`
    RiskLevel   string `json:"riskLevel" jsonschema:"safe, caution, or destructive"`
}
var proposeAction = copilot.DefineTool(
    "proposeAction",
    "Propose an executable action for the user to approve. Used to suggest next steps with explanation.",
    func(params ProposeActionParams, inv copilot.ToolInvocation) (string, error) { ... },
)
```

All read-only tools (1-6) set `SkipPermission = true` to avoid unnecessary permission prompts.

##### SDK Session Configuration

```go
session, err := client.CreateSession(&copilot.SessionConfig{
    OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
    Model:               "gpt-4.1",
    Streaming:           true,
    Tools: []copilot.Tool{
        getMemoryIndex, getInstallManifest, getInstallDrift,
        getRuntimeHealth, getActiveContext, listRecommendations,
        proposeAction,
    },
    SystemMessage: &copilot.SystemMessage{
        Content: `You are the ATV Launchpad intelligence engine. Your role is to analyze
the developer's repo state, installed development tools, and workflow artifacts to provide
context-aware recommendations. Use the available tools to query local state, then synthesize
insights. When recommending actions, always explain WHY — connect the recommendation to
specific evidence from the repo state. Propose concrete, executable actions when appropriate.
Prioritize: 1) fixing broken state, 2) progressing workflow, 3) optimizing setup.`,
    },
})
```

##### Rate Limiting & Cost Control (see SpecFlow gap #5)

- **Minimum interval:** 30 seconds between SDK queries (configurable via `ATV_SDK_INTERVAL`)
- **Query budget:** Max 60 queries per hour (configurable via `ATV_SDK_BUDGET`)
- **Batching:** State changes within the debounce window (300ms) are batched into a single SDK query
- **Trigger heuristic:** Only query SDK when state actually changes (new/modified/deleted artifacts), not on every FS event
- **Backoff on errors:** Exponential backoff starting at 30s, capping at 5 minutes on repeated failures
- **Token cost visibility:** Log estimated token usage per query (based on tool result sizes)

##### Auth Detection & Graceful Degradation (see SpecFlow gap #3)

```go
func (i *Intelligence) Start(ctx context.Context) error {
    client := copilot.NewClient(nil)  // uses logged-in Copilot CLI user by default
    if err := client.Start(); err != nil {
        // Auth not available — run in offline mode
        i.online = false
        log.Info("GitHub Copilot SDK unavailable, running in offline mode")
        return nil  // NOT an error — offline is a valid state
    }
    i.client = client
    i.online = true
    // ... create session
}
```

**Auth expiry mid-session:** Catch SDK errors, classify as auth errors, transition to offline mode, surface "SDK disconnected — re-authenticate with `copilot auth`" in TUI status bar. Periodically retry (every 5 minutes).

**Network flap handling:** Circuit breaker pattern — after 3 consecutive failures within 60 seconds, enter offline mode. Retry after 5 minutes.

##### Tasks

- [ ] Create `pkg/sdk/intelligence.go` with `Intelligence` struct and lifecycle methods
- [ ] Define all 7 typed tools with proper param/result structs
- [ ] Implement rate limiting (30s interval, 60/hour budget, exponential backoff)
- [ ] Implement auth detection with graceful degradation to offline mode
- [ ] Implement circuit breaker for network flap handling
- [ ] Wire SDK recommendations into `LiveState.SDKRecommendations`
- [ ] Add on-demand `Explain()` for individual recommendations
- [ ] Add dependency: `go get github.com/github/copilot-sdk/go`
- [ ] Add `intelligence_test.go` — test offline degradation, rate limiting, tool definitions
- [ ] Bundle Copilot CLI: `go get -tool github.com/github/copilot-sdk/go/cmd/bundler`

##### Acceptance Criteria

- [ ] SDK session initializes with all 7 typed tools registered
- [ ] Read-only tools (1-6) have `SkipPermission = true`
- [ ] `proposeAction` tool produces valid `ProposedAction` structs with risk levels
- [ ] Graceful degradation: SDK auth failure → offline mode (no crash, no error dialog)
- [ ] Rate limiting enforced: no more than 1 query per 30s and 60 per hour
- [ ] Circuit breaker activates after 3 consecutive failures within 60s
- [ ] SDK recommendations include explanation text (WHY)
- [ ] System message guides the model to connect recommendations to evidence

---

#### Phase 5: Suggest-then-Execute Action System

**Objective:** Enable users to approve and execute SDK-proposed actions directly from the launchpad TUI.

##### Action Types & Risk Tiers (see SpecFlow gap #7)

| Risk Level | Actions | UI Treatment |
|------------|---------|-------------|
| **Safe** | Read-only queries, scaffold new files, open documentation | Green indicator. Single `Enter` to approve. |
| **Caution** | Installer re-runs (`atv init`), gstack sync, file modifications | Yellow indicator. Confirmation prompt: "This will modify files. Proceed? [y/n]" |
| **Destructive** | Git operations, file deletion, `--force` flags | Red indicator. Double confirmation: "⚠️ DESTRUCTIVE: This will [action]. Type 'yes' to confirm:" |

##### Action Execution: `pkg/monitor/executor.go`

```go
type Executor struct {
    root   string
    queue  chan ProposedAction
    mu     sync.Mutex
    active bool  // only one action at a time (see SpecFlow gap: concurrent executions)
}

func NewExecutor(root string) *Executor
func (e *Executor) Execute(ctx context.Context, action ProposedAction) (*ActionResult, error)
func (e *Executor) IsActive() bool
```

**Execution constraints:**
- **Sequential only:** One action at a time. If an action is in progress, new approvals are queued (max queue depth: 3).
- **Working directory:** Always the repo root (`root` from `LaunchpadSnapshot`)
- **Environment:** Inherits user's `PATH` and environment
- **Timeout:** 2 minutes per action. Long-running actions (e.g., `atv init`) show a spinner.
- **Output capture:** stdout/stderr captured and displayed in a collapsible panel below the Moves tab
- **Interactive actions:** Actions that require TUI input (e.g., `atv init` wizard) suspend the launchpad, run in foreground, and resume the launchpad after completion.

##### SDK Output Validation (see SpecFlow security gap: hallucinated commands)

```go
// Allowlist of executable command prefixes
var allowedCommands = []string{
    "atv ",           // atv init, atv launchpad
    "gstack ",        // gstack sync, gstack office-hours
    "git ",           // git operations
    "copilot ",       // copilot CLI commands
}

func ValidateAction(action ProposedAction) error {
    // 1. Check command is in allowlist
    // 2. Check risk level matches command type
    // 3. Reject shell metacharacters (;, |, &&, $(), backticks)
    // 4. Reject path traversal (../)
}
```

##### Prompt Injection Mitigation (see SpecFlow security gap)

Since SDK tools read repo files (brainstorms, plans) and feed content to the model, a malicious file could attempt prompt injection. Mitigations:

1. **Human-in-the-loop:** ALL actions require explicit user approval — no autonomous execution
2. **Command allowlist:** Only recognized command prefixes are executable
3. **Shell metacharacter rejection:** No pipes, redirects, subshells, or variable expansion
4. **Risk level enforcement:** Destructive actions require double confirmation regardless of SDK confidence
5. **Output inspection:** Action results are displayed to the user, not fed back to the SDK automatically

##### Tasks

- [ ] Create `pkg/monitor/executor.go` with `Executor`, sequential execution queue, timeout handling
- [ ] Implement command allowlist validation with shell metacharacter rejection
- [ ] Implement risk-tiered approval UI in TUI (single Enter / confirm / double confirm)
- [ ] Implement output capture and collapsible panel display
- [ ] Implement interactive action handling (suspend launchpad, foreground, resume)
- [ ] Add `executor_test.go` — test allowlist, risk tiers, sequential execution, timeout

##### Acceptance Criteria

- [ ] Only allowlisted commands can be executed
- [ ] Shell metacharacters (`;`, `|`, `&&`, `$()`, backticks) are rejected
- [ ] Safe actions require single `Enter`, Caution requires `y/n`, Destructive requires typing `yes`
- [ ] Only one action executes at a time
- [ ] Action timeout at 2 minutes with clean cancellation
- [ ] Interactive actions suspend and resume the launchpad without corruption
- [ ] Action output is captured and displayed in the Moves panel

---

#### Phase 6: VS Code Extension — Webview Panel

**Objective:** Ship a VS Code extension that renders the launchpad as a webview panel, reading state from `.atv/launchpad-state.json` and observing Layer 3 (active skills/agents) via the extension API.

##### Extension Architecture

```
vscode-atv-launchpad/
├── package.json         # Extension manifest, activation events
├── src/
│   ├── extension.ts     # Activation, webview provider registration
│   ├── statePoller.ts   # Polls .atv/launchpad-state.json with FileSystemWatcher
│   ├── layer3.ts        # VS Code API: observe loaded skills, agents, extensions
│   └── webview/
│       ├── index.html   # Dashboard HTML shell
│       ├── main.js      # Panel rendering, message handling
│       └── style.css    # VS Code theme variable integration
├── tsconfig.json
└── .vscodeignore
```

##### State Sync: File-Based Polling

The VS Code extension does **not** communicate with the Go daemon via IPC. It reads `.atv/launchpad-state.json` using a `FileSystemWatcher`:

```typescript
// statePoller.ts
const watcher = vscode.workspace.createFileSystemWatcher(
    new vscode.RelativePattern(workspaceRoot, '.atv/launchpad-state.json')
);
watcher.onDidChange(() => {
    const state = JSON.parse(fs.readFileSync(stateFilePath, 'utf8'));
    panel.webview.postMessage({ type: 'stateUpdate', state });
});
```

**Standalone mode:** If the Go daemon isn't running (no `.atv/launchpad-state.json`), the extension falls back to its own state computation — calling `ScanRepoState()` equivalent logic in TypeScript, or simply showing "Run `atv launchpad` for realtime monitoring."

##### Layer 3: Extension API Observability

```typescript
// layer3.ts — VS Code-only capability
export function getActiveSkillsAndAgents(): Layer3State {
    return {
        loadedExtensions: vscode.extensions.all
            .filter(ext => ext.isActive)
            .map(ext => ext.id),
        // Observe copilot-related extensions specifically
        copilotActive: vscode.extensions.getExtension('github.copilot')?.isActive ?? false,
        copilotChatActive: vscode.extensions.getExtension('github.copilot-chat')?.isActive ?? false,
    };
}
```

##### Webview Rendering

- Use VS Code CSS theme variables (`--vscode-editor-foreground`, etc.) for native look
- Four panels matching TUI: Memory, Context, Health, Moves
- `postMessage` from extension → webview for state updates
- `postMessage` from webview → extension for action approval
- `getState()`/`setState()` for panel state persistence across tab switches

##### Distribution

- Published as a VS Code extension to the Marketplace
- Separate npm package / VSIX from the Go CLI (`atv-installer`)
- Version-locked: extension `package.json` specifies compatible state file schema version

##### Tasks

- [ ] Create `vscode-atv-launchpad/` extension scaffold with TypeScript, webpack build
- [ ] Implement `statePoller.ts` with `FileSystemWatcher` on `.atv/launchpad-state.json`
- [ ] Implement `layer3.ts` for extension API observability
- [ ] Build webview HTML/CSS/JS dashboard with four panels
- [ ] Use VS Code theme CSS variables for native appearance
- [ ] Implement bidirectional `postMessage` for state updates and action approval
- [ ] Implement standalone fallback when daemon isn't running
- [ ] Add `package.json` activation events: `onView`, `workspaceContains:.atv/`
- [ ] Write extension integration tests

##### Acceptance Criteria

- [ ] Extension activates when `.atv/` directory exists in workspace
- [ ] Webview renders four panels matching TUI layout
- [ ] State updates appear within 500ms of `.atv/launchpad-state.json` change
- [ ] Panel respects VS Code light/dark/high-contrast themes
- [ ] Layer 3 shows active Copilot extensions
- [ ] Standalone mode shows meaningful fallback when daemon isn't running
- [ ] Panel state persists across tab switches

---

## Alternative Approaches Considered

1. **On-Demand SDK Analysis** — Scan state at launch, query SDK lazily. Rejected: can't do proactive alerts or realtime monitoring. (see brainstorm: Why This Approach)

2. **VS Code Extension-First** — Build everything as a VS Code extension with webview. Rejected: splits codebase, loses terminal-first users, requires full rewrite. (see brainstorm)

3. **Out-of-Process Daemon with TCP/Socket IPC** — Daemon as separate process, TUI connects via socket. Rejected for Phase 1: adds IPC protocol design, port management, PID files, cross-platform socket differences without delivering user-facing value vs. in-process + file-sharing. Can be extracted later if headless daemon mode is needed.

4. **Chat-based Copilot Assistant** — SDK as a chat pane in the dashboard. Rejected: users want a monitoring dashboard, not a conversation. The SDK should power the recommendation engine invisibly. (see brainstorm: Key Decisions)

## System-Wide Impact

### Interaction Graph

- `atv launchpad` → creates `Watcher` → starts `fsnotify` → watches 5+ directory trees → updates `LiveState` →
  - TUI reads `LiveState` in-process → renders panels
  - State serialized to `.atv/launchpad-state.json` → VS Code extension reads via `FileSystemWatcher`
- `atv launchpad` → creates `Intelligence` → starts Github Copilot SDK client → creates session with 7 tools →
  - On state change + rate limit window: queries SDK → receives recommendations → updates `LiveState.SDKRecommendations`
  - On user action approval: `Executor.Execute()` → runs command → captures output → updates state
- `atv init` → writes files → triggers `fsnotify` events → `Watcher` debounces → state update → TUI re-render + drift recomputation

### Error Propagation

- **fsnotify errors** → logged, watcher continues. Critical errors (e.g., inotify limit reached) → TUI shows warning in footer, falls back to 3-second polling as degraded mode.
- **SDK auth errors** → transition to offline mode, TUI shows "⚪ Offline". Retry every 5 minutes.
- **SDK rate limit errors** → extend backoff interval, reduce query budget for remainder of hour.
- **Executor errors** → display error in Moves panel output area. State is unchanged (no partial state corruption).
- **State serialization errors** → logged. VS Code extension shows stale data until next successful write.

### State Lifecycle Risks

- **Partial FS event processing:** Debounce window ensures batch processing. If crash occurs during state update, next startup re-scans all watched directories.
- **Stale `.atv/launchpad-state.json`:** VS Code extension shows "Last updated: [timestamp]" — stale data is visible but not silent.
- **Concurrent manifest modification:** If `atv init` runs while `atv launchpad` is watching, the watcher detects the manifest write and refreshes state. No lock contention since reads are non-locking.

### API Surface Parity

- `RepoState` struct (existing) → extended by `LiveState` (new). All existing fields preserved — `BuildLaunchpadSnapshot()` continues to work for static mode.
- `Recommendation` struct (existing) → `SDKRecommendation` (new) adds `Explanation` and `ProposedAction`. Local recommendations remain as `[]Recommendation` in the snapshot.
- `cmd/launchpad.go` `--static` flag continues to work — uses existing `BuildLaunchpadSnapshot()` path without watcher or SDK.

### Integration Test Scenarios

1. **File change → TUI update:** Create a file in `docs/brainstorms/`, verify TUI Memory panel updates within 300ms.
2. **Offline → Online transition:** Start launchpad without Copilot auth, authenticate mid-session, verify SDK activates and recommendations gain WHY explanations.
3. **Drift detection accuracy:** Install with manifest, modify one file, update catalog, verify three-way classification (stale vs user-modified).
4. **Action execution end-to-end:** SDK proposes action, user approves in TUI, verify command executes in repo root and output appears in Moves panel.
5. **VS Code state sync:** Run `atv launchpad`, modify a file, verify VS Code extension webview updates within 500ms.

## Acceptance Criteria

### Functional Requirements

- [ ] `atv launchpad` starts an event-driven monitoring dashboard (not polling-based)
- [ ] Five monitoring layers update in realtime from filesystem events
- [ ] GitHub Copilot SDK provides context-aware recommendations with WHY explanations when authenticated
- [ ] Suggest-then-execute enables action approval from the dashboard
- [ ] VS Code webview panel shows the same four panels as the TUI
- [ ] Hybrid offline/online: monitoring works without auth, SDK intelligence requires auth

### Non-Functional Requirements

- [ ] FS event to TUI update latency < 300ms
- [ ] SDK query rate limited to ≤ 1/30s and ≤ 60/hour
- [ ] No crash or hang on auth expiry, network loss, or missing directories
- [ ] Cross-platform: Windows, macOS, Linux (CI matrix)
- [ ] Backward compatible: `atv launchpad --static` continues to work
- [ ] Action execution sandboxed: allowlisted commands only, no shell metacharacters

### Quality Gates

- [ ] Unit tests for all new packages (`monitor`, `sdk`, `executor`)
- [ ] Table-driven tests following existing patterns (see `manifest_test.go`, `recommendations_test.go`)
- [ ] Integration tests for 5 scenarios listed above
- [ ] CI passes on all three platforms with Go 1.26.1
- [ ] golangci-lint clean (install-mode: goinstall)

## Success Metrics

- **Engagement:** Users keep the launchpad running during development sessions (session duration > 30 minutes)
- **Actionability:** > 50% of SDK-proposed actions are approved and executed by users
- **Drift detection value:** Users discover and resolve stale files they wouldn't have noticed manually
- **Dual surface adoption:** Both TUI and VS Code panel see active usage

## Dependencies & Prerequisites

| Dependency | Type | Status |
|-----------|------|--------|
| `github.com/fsnotify/fsnotify` | Go module | New dependency |
| `github.com/github/copilot-sdk/go` | Go module | New dependency |
| GitHub Copilot CLI | Runtime (optional) | User prerequisite for SDK features |
| VS Code Extension API | TypeScript | New codebase |
| Existing `pkg/installstate` types | Internal | Extend, not replace |
| Existing `pkg/tui` TUI | Internal | Refactor |
| Existing `pkg/scaffold` catalog | Internal | Read-only access for drift detection |

## Risk Analysis & Mitigation

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| GitHub Copilot SDK Go API breaks or changes | Blocks Phase 4 | Medium | Pin SDK version in go.mod. SDK is pre-1.0, so expect breaking changes. |
| fsnotify cross-platform inconsistencies | Broken monitoring on some OS | Medium | Prototype on all 3 OS targets before committing. Fall back to polling on failure. |
| SDK token costs exceed user expectations | User frustration, unexpected bills | Low | Rate limiting, query budget, cost visibility in TUI |
| Prompt injection via malicious repo files | Security event | Low | Human-in-the-loop for all actions, command allowlist, shell metacharacter rejection |
| VS Code extension + Go binary distribution complexity | Complicated install | Medium | Separate distribution: Go via npm/binary, VS Code via Marketplace |
| In-process architecture limits future headless daemon mode | Architectural debt | Low | File-based state sharing is extractable to IPC later |

## Future Considerations

- **Headless daemon mode:** Extract state manager into a separate process with IPC for CI/CD and background monitoring use cases
- **WebSocket upgrade:** Replace file-based state sharing with WebSocket for lower-latency VS Code extension updates
- **Telemetry dashboard:** Aggregate SDK query patterns, action approval rates, and session durations for product insights
- **Custom tool registry:** Let users define their own SDK tools beyond the built-in 7
- **Multi-repo monitoring:** Watch multiple repos from a single launchpad instance

## Documentation Plan

- [ ] Update README.md with new `atv launchpad` capabilities
- [ ] Add monitoring layer documentation to `docs/`
- [ ] Document VS Code extension installation and usage
- [ ] Document suggest-then-execute action system and risk tiers
- [ ] Document offline vs. online mode differences

## Sources & References

### Origin

- **Brainstorm document:** [docs/brainstorms/2026-04-02-copilot-sdk-launchpad-elevation-brainstorm.md](../brainstorms/2026-04-02-copilot-sdk-launchpad-elevation-brainstorm.md) — Key decisions carried forward: SDK-powered monitoring daemon, five monitoring layers, dual surface architecture, hybrid offline/online, suggest-then-execute with full scope

### Internal References

- Superseded plan: [docs/plans/2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md](2026-03-31-001-feat-elevated-guided-installer-launchpad-plan.md) Phase 4
- Install manifest solution: [docs/solutions/2026-03-31-guided-install-manifest-and-recommendations.md](../solutions/2026-03-31-guided-install-manifest-and-recommendations.md)
- Architecture spike: [docs/spikes/architecture-post-install-memory-launchpad-spike.md](../spikes/architecture-post-install-memory-launchpad-spike.md)
- Current launchpad: `pkg/tui/launchpad.go`, `pkg/installstate/launchpad.go`
- State types: `pkg/installstate/types.go`
- Recommendations engine: `pkg/installstate/recommendations.go`
- Scaffold catalog (for drift): `pkg/scaffold/catalog.go`

### External References

- GitHub Copilot SDK: `github.com/github/copilot-sdk/go`
- fsnotify: `github.com/fsnotify/fsnotify`
- VS Code Webview API: https://code.visualstudio.com/api/extension-guides/webview
- Charmbracelet Bubble Tea: `github.com/charmbracelet/bubbletea`
