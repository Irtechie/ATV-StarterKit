<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=0:FF8C00,50:FFD700,100:FFA500&height=250&section=header&text=ATV%20STARTER%20KIT&fontSize=55&fontColor=ffffff&animation=fadeIn&fontAlignY=35&desc=One%20command.%20Instant%20agentic%20coding.&descAlignY=55&descSize=18&descColor=ffffff" width="100%"/>

**A**gentic **T**ool & **V**ibes — a one-click installer that scaffolds a complete AI-powered engineering environment into any project. Compound Engineering workflows for planning and review, plus Garry Tan's [gstack](https://github.com/garrytan/gstack) sprint process for QA, shipping, safety, and browser-based testing. **43 skills. 28 agents. One command.**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![GitHub Copilot](https://img.shields.io/badge/GitHub%20Copilot-Powered-8957e5?style=flat-square&logo=github)](https://github.com/features/copilot)

</div>

## What is ATV?

ATV gives you a virtual engineering team inside GitHub Copilot. It combines two complementary systems into a single installer:

- **Compound Engineering** — a multi-agent pipeline for brainstorming, planning, code review, and institutional knowledge. Think → Plan → Build → Review → Document.
- **gstack** — Garry Tan's open-source sprint process with 30 slash-command skills for QA, shipping, safety guardrails, security audits, and browser-based testing.

Together they cover the full software lifecycle: from "what should I build?" through "is it deployed and healthy?" — with 43 skills and 28 specialized agents, all discoverable by GitHub Copilot.

## Installation

### Option 1: npm / npx (Recommended)

Requires [Node.js 16+](https://nodejs.org/):

```bash
# Quick run (no global install)
npx atv-starterkit init

# Or install globally
npm install -g atv-starterkit
atv-starterkit init
```

### Option 2: Download Binary

> **Zero dependencies** — single static binary, works immediately.

Download from [GitHub Releases](https://github.com/All-The-Vibes/ATV-StarterKit/releases):

| Platform | Download |
|----------|----------|
| **Windows** | `atv-installer_windows_amd64.zip` |
| **macOS (Intel)** | `atv-installer_darwin_amd64.tar.gz` |
| **macOS (Apple Silicon)** | `atv-installer_darwin_arm64.tar.gz` |
| **Linux** | `atv-installer_linux_amd64.tar.gz` |

```bash
# macOS/Linux
tar xzf atv-installer_*.tar.gz
sudo mv atv-installer /usr/local/bin/
```

### Option 3: Build from Source

```bash
git clone https://github.com/All-The-Vibes/ATV-StarterKit.git
cd ATV-StarterKit
go build -o atv-installer .
sudo mv atv-installer /usr/local/bin/
```

## Quick Start

### One-Click Mode (Default)

```bash
cd your-project
atv-installer init
```

Auto-detects your stack, installs all 6 Copilot lifecycle hooks and 13 core skills. Done in seconds.

### Guided Mode — Choose Your Skills

```bash
atv-installer init --guided
```

The interactive TUI presents skills organized by function — not source. Pick what you need:

```
┃ Workflow skills (ATV + gstack)
┃ Prerequisites: git 2.52.0, bun 1.3.10, node v22.18.0
┃
┃   [•] Brainstorming — explore what to build
┃   [•] Plan — structured implementation plans
┃   [•] [gstack] Office Hours — YC-style forcing questions
┃   [•] [gstack] Plan CEO Review — find the 10-star product
┃   [•] CE Review — multi-agent code review
┃   [•] [gstack] Review — staff-level PR review
┃   [•] [gstack] QA — test app in real browser ⚠️ (requires Bun)
┃   [•] [gstack] Ship — sync, test, push, open PR
┃   [•] [gstack] Careful — warn before destructive commands
┃   ...
```

When you select gstack skills, the installer:
1. Clones gstack to `.gstack/` (a gitignored staging area)
2. Runs gstack's own setup to generate optimized skill docs
3. Copies each skill to `.github/skills/gstack-*/SKILL.md` — flat, one level deep, auto-discovered by Copilot
4. Creates a runtime sidecar at `.github/skills/gstack/` with binaries and assets

**Result:** 43 skills at `.github/skills/*/SKILL.md`, all discoverable by GitHub Copilot.

---

## The Full Sprint

ATV covers the complete software lifecycle. Here's how the skills work together:

```
  Think          Plan          Build         Review        Test          Ship          Reflect
   💭             📋            🔨            👀            🧪            🚀            📊
  brainstorm   ce-plan       ce-work       ce-review      qa           ship          retro
  office-hrs   plan-ceo      lfg/slfg      review         qa-only      land-deploy   compound
               plan-eng                    design-review  benchmark    canary        document-rel
               plan-design                 cso            browse       careful       learn
               autoplan                    codex                       freeze/guard
```

### Phase 1: Think — What should we build?

| Skill | What it does |
|-------|-------------|
| `/ce-brainstorm` | Interactive dialogue to clarify requirements; runs repo research; produces a design doc in `docs/brainstorms/` |
| `/gstack-office-hours` | YC-style forcing questions that reframe your product before you write code. Pushes back on your framing, challenges premises, generates alternatives |
| `/gstack-plan-ceo-review` | CEO-level review: find the 10-star product hiding in the request. Four modes: Expansion, Selective Expansion, Hold Scope, Reduction |

### Phase 2: Plan — How do we build it?

| Skill | What it does |
|-------|-------------|
| `/ce-plan` | Parallel research agents scan your codebase + external docs. Risk-based decision: security topics always get deep research. Outputs a plan in `docs/plans/` with acceptance criteria |
| `/deepen-plan` | Enriches each plan section with parallel research agents — best practices, performance, UI patterns |
| `/gstack-plan-eng-review` | Lock architecture, data flow, ASCII diagrams, edge cases, and tests. Forces hidden assumptions into the open |
| `/gstack-plan-design-review` | Rates each design dimension 0-10, explains what a 10 looks like, edits the plan to get there. AI slop detection |
| `/gstack-autoplan` | One command, fully reviewed plan. Runs CEO → design → eng review automatically |

### Phase 3: Build — Execute the plan

| Skill | What it does |
|-------|-------------|
| `/ce-work` | Reads the plan, breaks it into tasks, implements with incremental commits. System-Wide Test Check at each step (callbacks, orphaned state, error alignment). Checks off plan items as completed |
| `/lfg` | Full autonomous pipeline: plan → deepen → work → review → test → video → compound. Sequential gates ensure nothing is skipped |
| `/slfg` | Same as `/lfg` but parallelizes review + testing via swarm agents |

### Phase 4: Review — Find bugs before they ship

| Skill | What it does |
|-------|-------------|
| `/ce-review` | Launches parallel review agents (security, performance, architecture, language-specific). Configurable via `compound-engineering.local.md`. Ultra-thinking deep dive |
| `/gstack-review` | Staff-level code review. Auto-fixes obvious issues, flags completeness gaps |
| `/gstack-design-review` | Design audit then fix loop with atomic commits and before/after screenshots |
| `/gstack-cso` | OWASP Top 10 + STRIDE threat model. Zero-noise: 17 false positive exclusions, 8/10+ confidence gate |
| `/gstack-codex` | Independent code review from OpenAI Codex CLI — cross-model analysis |

### Phase 5: Test — Verify it works

| Skill | What it does |
|-------|-------------|
| `/gstack-qa` | Opens a real Chromium browser, clicks through flows, finds bugs, fixes them, generates regression tests, re-verifies. *Requires Bun* |
| `/gstack-qa-only` | Same QA methodology, but report only — no code changes |
| `/gstack-benchmark` | Baseline page load times, Core Web Vitals, resource sizes. Compare before/after on every PR |
| `/gstack-browse` | Give the agent eyes. Real Chromium, real clicks, ~100ms per command |

### Phase 6: Ship — Get it to production

| Skill | What it does |
|-------|-------------|
| `/gstack-ship` | Sync main, run tests, audit coverage, push, open PR. Bootstraps test frameworks if you don't have one |
| `/gstack-land-and-deploy` | Merge PR, wait for CI and deploy, verify production health. One command from "approved" to "verified in production" |
| `/gstack-canary` | Post-deploy monitoring loop. Watches for console errors and performance regressions |
| `/gstack-document-release` | Updates all project docs to match what you shipped. Catches stale READMEs |

### Phase 7: Reflect — Learn and protect

| Skill | What it does |
|-------|-------------|
| `/ce-compound` | Documents solved problems in `docs/solutions/` with YAML frontmatter. Future sessions search this via `learnings-researcher` — knowledge compounds |
| `/gstack-retro` | Team-aware weekly retro with per-person breakdowns, shipping streaks, test health trends |
| `/gstack-learn` | Per-project self-learning infrastructure |

### Safety Guardrails

| Skill | What it does |
|-------|-------------|
| `/gstack-careful` | Warns before destructive commands (rm -rf, DROP TABLE, force-push). Say "be careful" to activate |
| `/gstack-freeze` | Restrict file edits to one directory while debugging. Hard block, not just a warning |
| `/gstack-guard` | Careful + Freeze combined. Maximum safety for production work |
| `/gstack-investigate` | Systematic root-cause debugging. Iron Law: no fixes without investigation first |

---

## Knowledge Compounding Loop

This is what makes ATV more than a collection of skills. Knowledge flows through the system:

```
solve problem → /ce-compound documents it → docs/solutions/
                                                    ↓
future /ce-plan → learnings-researcher searches docs/solutions/ → avoids past mistakes
```

Plans live in `docs/plans/` (created by `/ce-plan`, consumed by `/ce-work`). Brainstorms live in `docs/brainstorms/` (created by `/ce-brainstorm`, auto-discovered by `/ce-plan`). Solutions live in `docs/solutions/` (created by `/ce-compound`, searched by every future planning session). Everything is file-based and git-tracked — your team's knowledge compounds with every PR.

---

## The Agent Roster (28 Specialized Agents)

Every agent in `.github/agents/` is a specialist that can be invoked by skills during review, planning, or debugging:

| Category | Agents |
|----------|--------|
| **Code Review** | `kieran-rails-reviewer`, `kieran-python-reviewer`, `kieran-typescript-reviewer`, `dhh-rails-reviewer`, `code-simplicity-reviewer`, `julik-frontend-races-reviewer` |
| **Security** | `security-sentinel` (OWASP, input validation, secrets) |
| **Architecture** | `architecture-strategist` (SOLID, coupling, boundaries) |
| **Performance** | `performance-oracle` (Big-O, DB queries, memory, scaling) |
| **Data** | `data-integrity-guardian`, `data-migration-expert`, `schema-drift-detector`, `deployment-verification-agent` |
| **Design** | `design-implementation-reviewer`, `design-iterator`, `figma-design-sync` |
| **Research** | `repo-research-analyst`, `best-practices-researcher`, `framework-docs-researcher`, `learnings-researcher`, `git-history-analyzer` |
| **Process** | `pr-comment-resolver`, `spec-flow-analyzer`, `bug-reproduction-validator`, `pattern-recognition-specialist` |
| **Meta** | `agent-native-reviewer`, `ankane-readme-writer` |
| **Ops** | `lint` |

---

## All 6 Copilot Lifecycle Hooks

| # | Hook | File | When It Fires |
|---|------|------|---------------|
| 1 | **System Instructions** | `.github/copilot-instructions.md` | Every Copilot chat — injected as system context |
| 2 | **Setup Steps** | `.github/copilot-setup-steps.yml` | When Copilot Coding Agent initializes |
| 3 | **MCP Servers** | `.github/copilot-mcp-config.json` | When Copilot starts — registers tool servers |
| 4 | **Skills** | `.github/skills/*/SKILL.md` | When skill description matches user's request |
| 5 | **Agents** | `.github/agents/*.agent.md` | When invoked by subagent orchestration |
| 6 | **File Instructions** | `.github/*.instructions.md` | Auto-loaded by `applyTo` glob when editing files |

## Supported Stacks

| Stack | Detection | Additional Content |
|-------|-----------|-------------------|
| **TypeScript** | `tsconfig.json` | `kieran-typescript-reviewer` agent, TS file instructions |
| **Python** | `pyproject.toml` / `requirements.txt` | `kieran-python-reviewer` agent, Python file instructions |
| **Rails** | `Gemfile` + `config/routes.rb` | 8 additional agents (DHH, data integrity, schema drift, ...), Ruby file instructions |
| **General** | fallback | Universal agents and skills only |

## MCP Servers

Pre-configured in `.github/copilot-mcp-config.json`:

| Server | Type | Package |
|--------|------|---------|
| **Context7** | SSE (remote) | `mcp.context7.com` |
| **GitHub** | stdio (npx) | `@modelcontextprotocol/server-github` (needs PAT) |
| **Azure** | stdio (npx) | `@azure/mcp` |
| **Terraform** | stdio (npx) | `terraform-mcp-server` |

---

## Prerequisites

### Required

- **Git** — for gstack clone and general usage
- **Node.js 16+** — for npx-based MCP servers

### Optional (for gstack browser skills)

- **Bun** — required for `/gstack-qa`, `/gstack-browse`, `/gstack-benchmark`. Install: [bun.sh](https://bun.sh)
- **GitHub PAT** — needed for GitHub MCP server (`repo`, `read:org` scopes)
- **Azure CLI** — `az login` for Azure MCP authentication

Without Bun, text-based gstack skills (`/gstack-review`, `/gstack-ship`, `/gstack-careful`, etc.) work fine. Only browser-based skills are disabled.

## How It Works Under the Hood

```
atv-installer init --guided
        │
        ▼
 Detect stack (TS/Python/Rails/General) + detect git/bun/node
        │
        ▼
 TUI: Pick skills by function (planning, review, QA, security, ...)
        │
        ├── ATV skills ──► Embedded templates → .github/skills/*/SKILL.md
        │
        └── gstack skills ──► git clone → .gstack/ (staging)
                                    │
                                    ├── bun run gen:skill-docs --host codex
                                    ├── Copy gstack-*/SKILL.md → .github/skills/
                                    └── Create sidecar: .github/skills/gstack/
                                         (bin/, browse/, ETHOS.md, review assets)
```

**Key details:**
- `.gstack/` is a gitignored staging area with the full gstack repo and runtime
- `.github/skills/gstack-*/SKILL.md` are the lightweight copies Copilot discovers
- `.github/skills/gstack/` is the runtime sidecar (binaries, utilities, checklists)
- All skills are at one level deep in `.github/skills/` — exactly what Copilot expects
- ATV skills have no prefix; gstack skills have `gstack-` prefix — no naming collisions
- Idempotent: running again skips existing files, merges JSON configs

---

## Development

```bash
# Build
go build -o atv-installer .

# Run locally
./atv-installer init
./atv-installer init --guided

# Test
go test ./...                              # all tests
go test ./pkg/gstack/ -v                   # gstack unit tests
go test ./test/sandbox/ -v                 # integration tests (sandbox)
go test ./test/sandbox/ -v -run TestGstack # network-dependent gstack tests (skip with -short)

# Test in a sandbox
mkdir /tmp/test-project && cd /tmp/test-project
echo '{}' > tsconfig.json && git init
/path/to/atv-installer init --guided
```

## Limitations

- **Bun required for browser skills** — `/gstack-qa`, `/gstack-browse`, `/gstack-benchmark` need Bun installed
- **Network required for gstack** — guided mode clones gstack at install time (~22MB)
- **Token-heavy pipeline** — running 5+ parallel agents in a long session can hit context limits
- **TUI requires proper terminal** — `--guided` mode needs Windows Terminal / iTerm (not VS Code integrated terminal)
- **gstack `./setup` on Windows** — falls back to `bun run gen:skill-docs` (bash path issues with Git Bash)

---

<div align="center">

## License

MIT

Built with ❤️ by [All The Vibes](https://github.com/All-The-Vibes) — powered by [gstack](https://github.com/garrytan/gstack)

</div>
