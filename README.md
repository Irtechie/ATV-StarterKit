<p align="center">
       <img src="./assets/hero-retro.svg" alt="ATV — All The Vibes 2.0 Starter Kit" width="100%" />
</p>

<h1 align="center">ATV — All The Vibes 2.0 Starter Kit</h1>

<p align="center"><strong>One command. Full agentic coding setup. Maximum tasteful chaos.</strong></p>

<p align="center">
       <a href="https://go.dev"><img alt="Go 1.26+" src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white"></a>
       <a href="https://opensource.org/licenses/MIT"><img alt="MIT License" src="https://img.shields.io/badge/License-MIT-ffd700?style=flat-square"></a>
       <a href="https://github.com/features/copilot"><img alt="GitHub Copilot Ready" src="https://img.shields.io/badge/GitHub%20Copilot-Ready-8957e5?style=flat-square&logo=github"></a>
       <a href="#the-full-sprint"><img alt="45 skills" src="https://img.shields.io/badge/Skills-45-ff8c00?style=flat-square"></a>
       <a href="#the-agent-roster"><img alt="29 agents" src="https://img.shields.io/badge/Agents-29-f97316?style=flat-square"></a>
</p>

<p align="center">
       <a href="#quick-start">Quick start</a> ·
       <a href="#the-three-pillars">Three pillars</a> ·
       <a href="#the-guided-experience">Guided experience</a> ·
       <a href="#the-launchpad">Launchpad</a> ·
       <a href="#the-full-sprint">Full sprint</a> ·
       <a href="#development">Development</a>
</p>

---

## What is ATV 2.0?

ATV 2.0 is a one-command installer that wires together three open-source systems into a single coherent agentic coding environment for GitHub Copilot:

- **Compound Engineering** — the planning-to-knowledge pipeline
- **gstack** — the sprint execution engine
- **agent-browser** — the browser automation layer

Each brings a distinct philosophy. Together they cover the full software lifecycle — from \"what should I build?\" through \"is it healthy in production?\" — with 45 skills, 29 agents, a memory-aware launchpad, and a knowledge system that makes your repo smarter with every PR.

---

## The Three Pillars

ATV 2.0 isn't a thing we built from scratch. It's the integration point for three independent projects, each with a philosophy worth understanding.

### Compound Engineering — knowledge compounds

**Origin:** [compound-engineering](https://github.com/EveryInc/compound-engineering-plugin) by Every

**Philosophy:** The first time you solve a problem takes hours of research. If you document it, the second time takes minutes. If you wire that documentation into your planning system, the third time is automatic. *Knowledge compounds.*

**What it provides:**
- `/ce-brainstorm` → `/ce-plan` → `/ce-work` → `/ce-review` → `/ce-compound` — a gated pipeline where each step produces an artifact the next step consumes
- `docs/solutions/` — structured solution documents with YAML frontmatter, searchable by the `learnings-researcher` agent during future planning sessions
- `docs/plans/` and `docs/brainstorms/` — living documents that track decisions, not just code
- `compound-engineering.local.md` — per-project configuration for which review agents fire

**The key insight:** Most AI coding tools treat every session as a blank slate. Compound Engineering treats every session as an investment. The `/ce-compound` skill writes what you learned into `docs/solutions/`, and the next time `/ce-plan` runs, the `learnings-researcher` agent searches those files first. Your repo accumulates institutional knowledge that prevents repeated mistakes.

### gstack — the AI sprint process

**Origin:** [gstack](https://github.com/garrytan/gstack) by Garry Tan (Y Combinator)

**Philosophy:** A single person with the right AI tooling can ship like a team of twenty. The difference isn't raw code generation speed — it's having a *process*. Think → Plan → Build → Review → Test → Ship → Reflect. Each step feeds the next. Nothing falls through the cracks because every skill knows what came before it.

**What it provides:**
- 30 slash-command skills covering office hours, engineering review, browser QA, shipping, deploy verification, security audits, safety guardrails, and weekly retros
- A real Chromium browser that the agent controls — not a mock, not a headless puppeteer script, but a persistent daemon with sub-second commands and cookie state
- Safety guardrails (`/gstack-careful`, `/gstack-freeze`, `/gstack-guard`) that prevent destructive commands before they execute
- Session tracking and per-project learning via `~/.gstack/`

**The key insight:** gstack doesn't just give the AI more tools. It gives the AI a *role*. `/gstack-review` acts as a staff engineer. `/gstack-cso` acts as a chief security officer. `/gstack-office-hours` acts as a YC partner challenging your premises. The skills aren't prompts — they're opinionated engineering processes encoded as markdown.

### agent-browser — the eyes of the agent

**Origin:** [agent-browser](https://github.com/vercel-labs/agent-browser) by Vercel

**Philosophy:** AI agents need to see the web the same way users do. Not through DOM dumps or HTML parsing, but through an accessibility tree with stable element references that survive page changes. Give the agent fast, reliable browser control and it can QA, debug, scrape, and test like a human — except at 100ms per command.

**What it provides:**
- A native Rust CLI that controls Chrome via CDP with ~100ms latency per command
- Snapshot refs (`@e1`, `@e2`) — deterministic element selection that works like screen coordinates but for the DOM
- Sessions, profiles, authentication vault, cookie persistence — the plumbing needed for real-world browser automation
- Security controls: domain allowlists, action policies, content boundaries, output limits

**The key insight:** Most browser automation tools are built for test suites. `agent-browser` is built for AI agents. The snapshot-ref workflow (`open → snapshot → interact → re-snapshot`) is designed to fit cleanly into an LLM's tool-calling loop. The agent doesn't need to write CSS selectors or XPath — it looks at the accessibility tree, picks a ref, and acts.

---

## Why Memory Matters

Most agentic coding setups are stateless. You install some skills, run some commands, and every session starts fresh. ATV 2.0 is different because memory is a first-class feature, not an afterthought.

### How memory works across the three pillars

| Layer | What remembers | Where it lives | Who reads it |
|---|---|---|---|
| **Institutional knowledge** | Solved problems, gotchas, patterns | `docs/solutions/*.md` (git-tracked) | `learnings-researcher` agent during `/ce-plan` and `/ce-review` |
| **Design decisions** | Why we chose approach A over B | `docs/brainstorms/*.md` (git-tracked) | `/ce-plan` auto-discovers recent brainstorms |
| **Implementation plans** | What to build, acceptance criteria, checkboxes | `docs/plans/*.md` (git-tracked) | `/ce-work` reads and checks off items as it implements |
| **Install manifest** | What the installer intended, attempted, skipped, failed | `.atv/install-manifest.json` (repo-local) | `atv-installer launchpad`, concierge tools |
| **Project config** | Which review agents to run, stack settings | `compound-engineering.local.md` | `/ce-review`, `/ce-work` |
| **gstack session state** | Active sessions, user preferences, prefix choice | `~/.gstack/` (user-global) | Every gstack skill preamble |
| **gstack project learning** | Per-project self-learning data | `.gstack/` (gitignored) | `/gstack-learn` |
| **Browser state** | Cookies, localStorage, login sessions | `~/.agent-browser/sessions/` | `agent-browser` session persistence |

The compound engineering memory loop is the most powerful:

```text
solve problem → /ce-compound documents it → docs/solutions/
                                                    ↓
future /ce-plan → learnings-researcher searches docs/solutions/ → avoids past mistakes
```

**Every PR makes your repo smarter.** Solutions are git-tracked, so they travel with the codebase. New team members get the benefit of every mistake the team already made and solved. This is the opposite of how most AI tools work — instead of losing context at the end of each session, you're building a searchable knowledge base that future sessions mine automatically.

---

## Quick Start

### 1. Install

```bash
cd your-project
npx atv-starterkit init
```

Auto-detects your stack. Installs 13 core ATV skills, 29 agents, MCP servers, and docs structure. Done in seconds.

Want to choose your preset and stack packs? Use `npx atv-starterkit init --guided` for the interactive TUI with multi-stack selection.

### 2. Use

Open **Copilot Chat** in VS Code (⌃⌘I / Ctrl+Shift+I) and run skills as slash commands:

```text
/ce-brainstorm   →  Explore the problem, produce a design doc
/ce-plan         →  Generate an implementation plan with acceptance criteria
/ce-work         →  Build against the plan with incremental commits
/ce-review       →  Multi-agent code review (security, architecture, performance)
/ce-compound     →  Document what you learned for future sessions
```

Or skip the steps and run the full pipeline in one shot:

```text
/lfg             →  Plan → deepen → build → review → test → compound
```

### 3. Reopen the Launchpad

```bash
atv-installer launchpad
```

Shows your memory dashboard: installed intelligence, repo memory snapshot, and deterministic next-step recommendations. Reopenable any time — no reinstall needed.

### 4. Compound

Every time you run `/ce-compound`, solved problems get saved to `docs/solutions/`. Next time `/ce-plan` runs, the `learnings-researcher` agent searches those files first — so your repo gets smarter with every PR.

---

## The Guided Experience

The guided installer walks you through four screens:

### Screen 1: Stack Packs

```text
┃ Which stack packs should be included?
┃ [✓] General
┃ [✓] TypeScript    (tsconfig.json detected)
┃ [ ] Python
┃ [ ] Rails
```

Multi-select — choose as many stacks as your project uses. Auto-detected packs are pre-selected. Stack packs are additive: selecting both TypeScript and Rails installs agents and file instructions for both.

### Screen 2: Preset

```text
┃ Choose your setup level
┃
┃ > ⚡ Starter — Core workflow (13 skills, instant)
┃     Plan, build, review, compound. No browser tools.
┃
┃   🚀 Pro — Full sprint process (35+ skills)
┃     + gstack review, ship, safety, security, debugging
┃
┃   🔥 Full — Complete engineering team (45+ skills)
┃     + browser QA, benchmarks, agent-browser, Chrome
┃     Requires: Bun, ~2min install
```

**Starter** is pure Compound Engineering — no network calls, instant install. **Pro** adds the gstack sprint skills (text-only, no browser). **Full** is everything: all 45 skills, gstack browser runtime, agent-browser CLI, and Chrome for Testing.

### Screen 3: Customize?

```text
┃ Want to customize individual skills?
┃   Yes, let me pick / No, install preset as-is
```

Power users can drill into category-grouped multi-select. Beginners skip straight to install.

### Screen 4: Install Progress

```text
  Installing Pro preset for typescript...

  ✅ Scaffolding ATV files (24 files created, 8 directories) · 340ms
  ⚠️  Syncing gstack skills — setup failed, fell back to docs · 2.1s
  ✅ Installing agent-browser (CLI ready, skill copied) · 1.8s
```

Real-time animated spinners. Each step shows pending → running → done/warned/failed with structured telemetry: durations, reasons, skip explanations. Substep events track individual file writes, git clone vs build stages, and npm operations.

### Screen 5: Summary + Recommendations

```text
  Guided install summary

  ✅ Scaffolding ATV files (24 files created) · 340ms
  ⚠️  Syncing gstack skills — fell back to markdown-only · 2.1s
  ✅ Installing agent-browser (CLI ready, skill copied) · 1.8s

  Recommended next moves

    1. Fix installer warnings before relying on every capability
       fell back to markdown-only sync
    2. Start with /ce-brainstorm to shape the first feature
       No brainstorms were found in docs/brainstorms yet.

  🎉 ATV Starter Kit ready!
  Install state saved to .atv/install-manifest.json
  Reopen later with: atv-installer launchpad
```

The installer writes a versioned manifest to `.atv/install-manifest.json` recording requested vs installed vs skipped vs failed outcomes. Deterministic recommendations derive from local repo state — no network required.

---

## The Launchpad

After install, run `atv-installer launchpad` to see your repo's memory dashboard:

```text
  ATV Launchpad
  Local memory + install intelligence for this repo

  Installed intelligence
  ✅ Guided manifest found at .atv/install-manifest.json
  • Last guided run: 2026-04-01 14:30 UTC
  • Preset: Pro
  • Stack packs: General, TypeScript
  • gstack: 2 skill dirs requested (markdown-only)
  • agent-browser: requested in the last guided run
  • Outcomes: 2 done, 1 warnings, 0 failed, 0 skipped

  Repo memory snapshot
  • brainstorms: 3
  • plans: 2
  • solutions: 1
  • active plan state: unchecked work remains
  • agents: 18, skills: 12
  • copilot-instructions.md: present
  • gstack staging: present (.gstack/)
  • agent-browser skill: installed

  Recommended next moves
  1. Continue the active plan with /ce-work
     At least one plan still has unchecked items.
  2. Use /gstack-office-hours for a deeper sprint kickoff
     gstack skills were requested and synced successfully enough to use.
  3. Open the app in a real browser with agent-browser
     Browser automation tooling was installed or partially prepared.
```

The launchpad is **deterministic** — same repo state always produces the same recommendations. It works fully offline, requires no auth, and reads only local filesystem state.

### Concierge Agent

The `atv-concierge` agent (scaffolded to `.github/agents/atv-concierge.agent.md`) lets Copilot Chat navigate the launchpad data:

| Command | Purpose |
|---|---|
| `atv-installer concierge memory-summary` | Full repo memory + install intelligence as JSON |
| `atv-installer concierge list-recommendations` | Deterministic next-step recommendations |
| `atv-installer concierge explain-recommendation [id]` | Why a recommendation exists + suggested command |
| `atv-installer concierge open-artifact [name]` | Resolve "plans" / "manifest" / "skills" to filesystem paths |
| `atv-installer concierge run-suggested-action [id]` | Describe the command (never auto-executes) |

The concierge is **secondary** — it explains and navigates the deterministic local model but never overrides it. All tools work without network or Copilot auth.

---

## The Full Sprint

ATV covers the complete software lifecycle:

<table>
       <tr>
              <td width="25%" valign="top">
                     <strong>💭 Think</strong><br />
                     <sub>Frame the problem</sub><br /><br />
                     <code>/ce-brainstorm</code><br />
                     <code>/gstack-office-hours</code>
              </td>
              <td width="25%" valign="top">
                     <strong>📋 Plan</strong><br />
                     <sub>Pressure-test the approach</sub><br /><br />
                     <code>/ce-plan</code><br />
                     <code>/gstack-plan-ceo-review</code><br />
                     <code>/gstack-plan-eng-review</code><br />
                     <code>/gstack-plan-design-review</code><br />
                     <code>/gstack-autoplan</code>
              </td>
              <td width="25%" valign="top">
                     <strong>🔨 Build</strong><br />
                     <sub>Execute with momentum</sub><br /><br />
                     <code>/ce-work</code><br />
                     <code>/lfg</code><br />
                     <code>/slfg</code>
              </td>
              <td width="25%" valign="top">
                     <strong>👀 Review</strong><br />
                     <sub>Find what you missed</sub><br /><br />
                     <code>/ce-review</code><br />
                     <code>/gstack-review</code><br />
                     <code>/gstack-design-review</code><br />
                     <code>/gstack-cso</code><br />
                     <code>/gstack-codex</code>
              </td>
       </tr>
       <tr>
              <td width="33.33%" valign="top">
                     <strong>🧪 Test</strong><br />
                     <sub>Use real browser eyes</sub><br /><br />
                     <code>agent-browser</code><br />
                     <code>/gstack-qa</code><br />
                     <code>/gstack-qa-only</code><br />
                     <code>/gstack-benchmark</code><br />
                     <code>/gstack-browse</code>
              </td>
              <td width="33.33%" valign="top">
                     <strong>🚀 Ship</strong><br />
                     <sub>Land without chaos</sub><br /><br />
                     <code>/gstack-ship</code><br />
                     <code>/gstack-land-and-deploy</code><br />
                     <code>/gstack-canary</code><br />
                     <code>/gstack-document-release</code>
              </td>
              <td width="33.33%" valign="top">
                     <strong>📊 Reflect</strong><br />
                     <sub>Compound what you learned</sub><br /><br />
                     <code>/ce-compound</code><br />
                     <code>/gstack-retro</code><br />
                     <code>/gstack-learn</code>
              </td>
       </tr>
</table>

> 🛡️ Safety guardrails apply across the whole sprint: `/gstack-careful`, `/gstack-freeze`, `/gstack-guard`, and `/gstack-investigate`.

### Think

| Skill | What it does |
|---|---|
| `/ce-brainstorm` | Interactive dialogue to clarify requirements; produces design docs in `docs/brainstorms/` |
| `/gstack-office-hours` | YC-style forcing questions that challenge your framing before you write code |
| `/gstack-plan-ceo-review` | CEO-level review: find the 10-star product hiding in the request |

### Plan

| Skill | What it does |
|---|---|
| `/ce-plan` | Parallel research agents scan codebase + external docs; auto-discovers brainstorms; outputs plans with acceptance criteria |
| `/deepen-plan` | Enriches each plan section with best practices and performance guidance |
| `/gstack-plan-eng-review` | Forces hidden assumptions into the open: architecture, data flow, edge cases |
| `/gstack-plan-design-review` | Scores design quality 0-10 per dimension; rewrites plan to hit 10 |
| `/gstack-autoplan` | Runs CEO → design → eng review in one command |

### Build

| Skill | What it does |
|---|---|
| `/ce-work` | Implements against the plan with incremental commits and system-wide sanity checks |
| `/lfg` | Full pipeline: plan → deepen → work → review → test → video → compound |
| `/slfg` | Parallelized version via swarm agents |

### Review

| Skill | What it does |
|---|---|
| `/ce-review` | Parallel review agents: security, performance, architecture, language-specific |
| `/gstack-review` | Staff-level code review with auto-fix and completeness checks |
| `/gstack-design-review` | Design audit with atomic fix commits |
| `/gstack-cso` | OWASP Top 10 + STRIDE threat model |
| `/gstack-codex` | Cross-model review via OpenAI Codex CLI |

### Test

| Skill | What it does |
|---|---|
| `agent-browser` | Direct browser automation: open, snapshot, click, fill, screenshot, inspect |
| `/gstack-qa` | Full QA loop: find bugs in real browser, fix them, write regressions, re-verify |
| `/gstack-qa-only` | Report-only QA |
| `/gstack-benchmark` | Page load baselines, Core Web Vitals, resource sizes |
| `/gstack-browse` | Persistent browser runtime for deeper sessions |

### Ship

| Skill | What it does |
|---|---|
| `/gstack-ship` | Sync main, run tests, audit coverage, push, open PR |
| `/gstack-land-and-deploy` | Merge → CI → deploy → verify production |
| `/gstack-canary` | Post-deploy monitoring for errors and regressions |
| `/gstack-document-release` | Auto-update project docs to match what shipped |

### Reflect

| Skill | What it does |
|---|---|
| `/ce-compound` | Documents solved problems in `docs/solutions/` — compounds knowledge for future sessions |
| `/gstack-retro` | Team-aware weekly retro with per-person breakdowns |
| `/gstack-learn` | Per-project self-learning infrastructure |

### Safety Guardrails

| Skill | What it does |
|---|---|
| `/gstack-careful` | Warns before `rm -rf`, `DROP TABLE`, force-push |
| `/gstack-freeze` | Restricts edits to one directory while debugging |
| `/gstack-guard` | Careful + Freeze combined |
| `/gstack-investigate` | No fixes without systematic investigation first |

---

## The Agent Roster

28 specialized agents in `.github/agents/`, invoked by skills during review, planning, and debugging:

| Category | Agents |
|---|---|
| **Code Review** | `kieran-rails-reviewer`, `kieran-python-reviewer`, `kieran-typescript-reviewer`, `dhh-rails-reviewer`, `code-simplicity-reviewer`, `julik-frontend-races-reviewer` |
| **Security** | `security-sentinel` |
| **Architecture** | `architecture-strategist` |
| **Performance** | `performance-oracle` |
| **Data** | `data-integrity-guardian`, `data-migration-expert`, `schema-drift-detector`, `deployment-verification-agent` |
| **Design** | `design-implementation-reviewer`, `design-iterator`, `figma-design-sync` |
| **Research** | `repo-research-analyst`, `best-practices-researcher`, `framework-docs-researcher`, `learnings-researcher`, `git-history-analyzer` |
| **Process** | `pr-comment-resolver`, `spec-flow-analyzer`, `bug-reproduction-validator`, `pattern-recognition-specialist` |
| **Meta** | `agent-native-reviewer`, `ankane-readme-writer` |
| **Launchpad** | `atv-concierge` |
| **Ops** | `lint` |

---

## What Gets Installed

### All 6 Copilot Lifecycle Hooks

| # | Hook | File | When it fires |
|---|---|---|---|
| 1 | **System Instructions** | `.github/copilot-instructions.md` | Every Copilot chat |
| 2 | **Setup Steps** | `.github/copilot-setup-steps.yml` | Coding Agent initialization |
| 3 | **MCP Servers** | `.github/copilot-mcp-config.json` | Copilot startup |
| 4 | **Skills** | `.github/skills/*/SKILL.md` | When description matches request |
| 5 | **Agents** | `.github/agents/*.agent.md` | Subagent orchestration |
| 6 | **File Instructions** | `.github/*.instructions.md` | `applyTo` glob matches |

### Supported Stacks

| Stack | Detection | Additions |
|---|---|---|
| **TypeScript** | `tsconfig.json` | TypeScript reviewer, TS file instructions |
| **Python** | `pyproject.toml` / `requirements.txt` | Python reviewer, Python file instructions |
| **Rails** | `Gemfile` + `config/routes.rb` | 8 Rails-specific agents, Ruby file instructions |
| **General** | fallback | Universal agents and skills |

### MCP Servers

| Server | Type | Package |
|---|---|---|
| **Context7** | SSE | `mcp.context7.com` |
| **GitHub** | stdio | `@modelcontextprotocol/server-github` |
| **Azure** | stdio | `@azure/mcp` |
| **Terraform** | stdio | `terraform-mcp-server` |

---

## How It Works Under the Hood

```text
atv-installer init --guided
        │
        ▼
 Detect stack + prerequisites (git, bun, node)
        │
        ▼
 Screen 1: Stack Packs (multi-select) → Screen 2: Preset → Screen 3: Customize?
        │
        ▼
 Install with structured telemetry:
        │
        ├── ATV scaffold ──► Embedded templates → .github/skills/*/SKILL.md
        │                    └── Substep events per file (created/skipped/merged)
        │
        ├── gstack ──► git clone → .gstack/ (staging)
        │               ├── gen:skill-docs → .agents/skills/gstack-*/
        │               ├── Copy SKILL.md → .github/skills/gstack-*/
        │               └── Substeps: clone → build/doc-gen → copy skills
        │
        └── agent-browser ──► npm install -g → agent-browser install (Chrome)
                              ├── .github/skills/agent-browser/SKILL.md
                              └── Substeps: npm install → copy SKILL.md
        │
        ▼
 Write manifest to .atv/install-manifest.json
        │
        ├── Requested state (packs, layers, preset)
        ├── Outcomes with substeps + skip reasons
        └── Deterministic recommendations
        │
        ▼
 atv-installer launchpad  ──► Reopenable memory dashboard
 atv-installer concierge  ──► JSON API for Copilot SDK assistant
```

- `.gstack/` is gitignored — staging area with the full repo and runtime
- `.github/skills/gstack-*/SKILL.md` are lightweight copies Copilot discovers
- `.github/skills/gstack/` is the runtime sidecar (binaries, checklists, ETHOS.md)
- All skills at one level deep in `.github/skills/` — Copilot's discovery convention
- Idempotent: re-running skips existing files, merges JSON configs

---

## Prerequisites

**Required:** Git, Node.js 16+

**Optional:**
- **Bun** — for gstack browser skills (`/gstack-qa`, `/gstack-browse`, `/gstack-benchmark`)
- **GitHub PAT** — for GitHub MCP server
- **Azure CLI** — for Azure MCP server

Without Bun, text-based gstack skills still work. `agent-browser` works independently of Bun.

## Installation

### npm (recommended)

```bash
npx atv-starterkit init              # quick run
npm install -g atv-starterkit        # global install
```

### Binary

Download from [GitHub Releases](https://github.com/All-The-Vibes/ATV-StarterKit/releases).

### From source

```bash
git clone https://github.com/All-The-Vibes/ATV-StarterKit.git
cd ATV-StarterKit && go build -o atv-installer .
```

## Development

```bash
go build -o atv-installer .             # build
go test ./...                            # all tests
go test ./pkg/concierge/ -v              # concierge tool tests
go test ./pkg/installstate/ -v           # manifest + recommendations tests
go test ./test/sandbox/ -v               # integration tests (22 E2E scenarios)
go test ./test/sandbox/ -v -run E2E      # comprehensive lifecycle tests only

# sandbox test
mkdir /tmp/test && cd /tmp/test
echo '{}' > tsconfig.json && git init
/path/to/atv-installer init --guided

# verify launchpad + concierge
/path/to/atv-installer launchpad
/path/to/atv-installer concierge memory-summary
/path/to/atv-installer concierge list-recommendations
```

## Limitations

- **Bun required for browser skills** — `/gstack-qa`, `/gstack-browse`, `/gstack-benchmark`
- **Network required for gstack** — clones ~22MB at install time
- **gstack setup on Windows** — falls back to `bun run gen:skill-docs` (bash path issues)
- **Token-heavy pipelines** — long multi-agent sessions can hit context limits

---

<div align="center">

MIT — Built by [All The Vibes](https://github.com/All-The-Vibes)

Powered by [Compound Engineering](https://github.com/EveryInc/compound-engineering-plugin) · [gstack](https://github.com/garrytan/gstack) · [agent-browser](https://github.com/vercel-labs/agent-browser)

</div>
