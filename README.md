<p align="center">
       <img src="./assets/hero-retro.svg" alt="ATV — All The Vibes 2.0 Starter Kit" width="100%" />
</p>

<h1 align="center">ATV — All The Vibes 2.0 Starter Kit</h1>

<p align="center"><strong>One command. Full agentic coding setup. Maximum tasteful chaos.</strong></p>

<p align="center">
       <a href="https://www.npmjs.com/package/atv-starterkit"><img alt="npm version" src="https://img.shields.io/npm/v/atv-starterkit?style=flat-square&logo=npm&logoColor=white&color=cb3837"></a>
       <a href="https://go.dev"><img alt="Go 1.26+" src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white"></a>
       <a href="https://opensource.org/licenses/MIT"><img alt="MIT License" src="https://img.shields.io/badge/License-MIT-ffd700?style=flat-square"></a>
       <a href="https://github.com/features/copilot"><img alt="GitHub Copilot Ready" src="https://img.shields.io/badge/GitHub%20Copilot-Ready-8957e5?style=flat-square&logo=github"></a>
       <a href="#the-full-sprint"><img alt="45 skills" src="https://img.shields.io/badge/Skills-45-ff8c00?style=flat-square"></a>
       <a href="#the-agent-roster"><img alt="29 agents" src="https://img.shields.io/badge/Agents-29-f97316?style=flat-square"></a>
</p>

<p align="center">
       <a href="#quick-start">Quick start</a> ·
       <a href="#installation">Installation</a> ·
       <a href="#the-three-pillars">Three pillars</a> ·
       <a href="#continuous-learning--your-repo-gets-smarter-every-session">Learning</a> ·
       <a href="#the-full-sprint">Full sprint</a> ·
       <a href="#development">Development</a>
</p>

<p align="center">
       <a href="./assets/demo.mp4">
              <img src="https://img.shields.io/badge/▶%20Watch%20Demo-91s-f59e0b?style=for-the-badge&logo=youtube&logoColor=white" alt="Watch Demo Video" />
       </a>
</p>

<video src="https://github.com/All-The-Vibes/ATV-StarterKit/releases/download/v2.4.0/demo.mp4" width="100%" autoplay loop muted playsinline controls></video>

<details>
<summary><strong>📽️ Demo Video (91 seconds)</strong> — click to expand scene list</summary>

<br />

**What you'll see:**
1. **The Problem** — AI that never remembers your patterns
2. **One Command** — `npx atv-starterkit@latest init` installs everything
3. **What You Get** — 45 skills, 29 agents, infinite memory
4. **The Full Sprint** — Think → Plan → Build → Review → Test → Ship → Reflect
5. **/unslop** ⭐ — Three parallel passes strip AI slop: Code Slop, Comment Rot, Design Slop
6. **The Learning Loop** — Observer hooks → /learn → instincts → /evolve → permanent skills
7. **Memory Cycles** — Three reinforcing loops that make your repo smarter
8. **Get Started** — One command away

</details>

---

## What is ATV 2.0?

ATV 2.0 is a one-command installer that wires together three open-source systems into a single coherent agentic coding environment for GitHub Copilot:

- **Compound Engineering** — the planning-to-knowledge pipeline
- **gstack** — the sprint execution engine
- **agent-browser** — the browser automation layer

Together they cover the full software lifecycle — from "what should I build?" through "is it healthy in production?" — with 45 skills, 29 agents, a continuous learning pipeline, and a knowledge system that makes your repo smarter with every session.

---

## Quick Start

### 1. Install

```bash
cd your-project
npx atv-starterkit@latest init
```

Auto-detects your stack. Installs 16 core skills, 29 agents, observer hooks, MCP servers, and docs structure. Done in seconds.

For the interactive TUI with multi-stack selection:

```bash
npx atv-starterkit@latest init --guided
```

**What `init` does:**

```
  ⚡ All The Vibes 2.0 ⚡
  One command. Full agentic coding setup.

  Auto-detected primary: typescript project (tsconfig.json found, existing git repo)
  Likely stack packs: TypeScript

  📁 .github/skills
  📁 .github/agents
  📁 .github/hooks/scripts
  📁 .vscode
  📁 docs/plans
  📁 docs/brainstorms
  📁 docs/solutions
  📁 .atv/instincts
  ✅ .github/copilot-instructions.md
  ✅ .github/copilot-setup-steps.yml
  ✅ .github/copilot-mcp-config.json
  ✅ .github/skills/atv-learn/SKILL.md
  ✅ .github/skills/atv-instincts/SKILL.md
  ✅ .github/skills/atv-evolve/SKILL.md
  ✅ .github/skills/atv-observe/SKILL.md
  ✅ .github/skills/atv-unslop/SKILL.md
  ✅ .github/skills/ce-brainstorm/SKILL.md
  ✅ .github/skills/ce-plan/SKILL.md
  ✅ .github/skills/ce-work/SKILL.md
  ✅ .github/skills/ce-review/SKILL.md
  ✅ .github/skills/ce-compound/SKILL.md
  ✅ .github/hooks/copilot-hooks.json
  ✅ .github/hooks/scripts/observe.js
  ...
```

### 2. Use

Open **Copilot Chat** in VS Code (⌃⌘I / Ctrl+Shift+I) and run skills as slash commands:

```text
/ce-brainstorm   →  Explore the problem, produce a design doc
/ce-plan         →  Generate an implementation plan with acceptance criteria
/ce-work         →  Build against the plan with incremental commits
/ce-review       →  Multi-agent code review (security, architecture, performance)
/ce-compound     →  Document what you learned for future sessions
```

Or run the full pipeline in one shot:

```text
/lfg             →  Plan → deepen → build → review → test → compound
```

### 3. Learn

ATV includes a **continuous learning pipeline** that captures your coding patterns and evolves them into project-specific skills:

```text
/learn           →  Extract patterns from recent work into instincts
/instincts       →  View learned patterns with confidence scores
/evolve          →  Promote mature instincts into discoverable skills
/observe         →  Run focused analysis on a specific domain
```

---

## Installation

### npm (recommended)

```bash
npx atv-starterkit@latest init       # quick run — downloads binary automatically
npm install -g atv-starterkit        # global install
atv-starterkit init                  # then run from anywhere
```

The npm package downloads the correct platform binary from [GitHub Releases](https://github.com/All-The-Vibes/ATV-StarterKit/releases) during install — no Go toolchain needed.

### Binary (direct download)

Grab a pre-built binary from [GitHub Releases](https://github.com/All-The-Vibes/ATV-StarterKit/releases/latest) for your platform (macOS, Linux, Windows — amd64/arm64).

### From source

```bash
git clone https://github.com/All-The-Vibes/ATV-StarterKit.git
cd ATV-StarterKit && go build -o atv-installer .
```

### Prerequisites

**Required:** Git, Node.js 16+

**Optional:**
- **Bun** — for gstack browser skills (`/gstack-qa`, `/gstack-browse`, `/gstack-benchmark`)
- **GitHub PAT** — for GitHub MCP server
- **Azure CLI** — for Azure MCP server

Without Bun, text-based gstack skills still work. `agent-browser` works independently of Bun.

---

## The Guided Experience

The guided installer (`--guided`) walks you through:

### Screen 1: Stack Packs

```text
┃ Which stack packs should be included?
┃ [✓] General
┃ [✓] TypeScript    (tsconfig.json detected)
┃ [ ] Python
┃ [ ] Rails
```

Multi-select — auto-detected packs are pre-selected. Stack packs are additive.

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

**Starter** is pure Compound Engineering — no network calls, instant install. **Pro** adds gstack sprint skills. **Full** is everything: all 45 skills, gstack browser runtime, agent-browser CLI, and Chrome for Testing.

### Screen 3: Customize?

Power users can drill into category-grouped multi-select. Beginners skip straight to install.

### Screen 4: Install Progress

```text
  Installing Pro preset for typescript...

  ✅ Scaffolding ATV files (24 files created, 8 directories) · 340ms
  ⚠️  Syncing gstack skills — setup failed, fell back to docs · 2.1s
  ✅ Installing agent-browser (CLI ready, skill copied) · 1.8s
```

Real-time animated spinners with structured telemetry: durations, skip reasons, substep events.

### Screen 5: Summary + Recommendations

```text
  Guided install summary

  ✅ Scaffolding ATV files (24 files created) · 340ms
  ⚠️  Syncing gstack skills — fell back to markdown-only · 2.1s
  ✅ Installing agent-browser (CLI ready, skill copied) · 1.8s

  Recommended next moves

    1. Fix installer warnings before relying on every capability
    2. Start with /ce-brainstorm to shape the first feature

  🎉 ATV Starter Kit ready!
  Install state saved to .atv/install-manifest.json
```

---

## The Three Pillars

### Compound Engineering — knowledge compounds

**Origin:** [compound-engineering](https://github.com/EveryInc/compound-engineering-plugin) by Every

A gated pipeline where each step produces an artifact the next step consumes:

- `/ce-brainstorm` → `/ce-plan` → `/ce-work` → `/ce-review` → `/ce-compound`
- `docs/solutions/` — structured solution docs, searchable by the `learnings-researcher` agent during future planning
- `docs/plans/` and `docs/brainstorms/` — living documents that track decisions, not just code

**The key insight:** Every time you run `/ce-compound`, solved problems get saved to `docs/solutions/`. Next time `/ce-plan` runs, the `learnings-researcher` agent searches those files first. Your repo gets smarter with every PR.

### gstack — the AI sprint process

**Origin:** [gstack](https://github.com/garrytan/gstack) by Garry Tan (Y Combinator)

- 30 slash-command skills covering office hours, engineering review, browser QA, shipping, deploy verification, security audits, safety guardrails, and weekly retros
- A real Chromium browser the agent controls with sub-second commands and cookie state
- Safety guardrails (`/gstack-careful`, `/gstack-freeze`, `/gstack-guard`) that prevent destructive commands

**The key insight:** gstack doesn't just give the AI more tools — it gives the AI a *role*. `/gstack-review` acts as a staff engineer. `/gstack-cso` acts as a chief security officer. The skills are opinionated engineering processes encoded as markdown.

### agent-browser — the eyes of the agent

**Origin:** [agent-browser](https://github.com/vercel-labs/agent-browser) by Vercel

- A native Rust CLI that controls Chrome via CDP with ~100ms latency per command
- Snapshot refs (`@e1`, `@e2`) — deterministic element selection for AI tool-calling loops
- Sessions, profiles, authentication vault, cookie persistence

**The key insight:** The snapshot-ref workflow (`open → snapshot → interact → re-snapshot`) fits cleanly into an LLM's tool-calling loop. No CSS selectors or XPath needed.

---

## Continuous Learning — Your Repo Gets Smarter Every Session

> **This is the headline feature of ATV 1.3.** Most AI coding tools treat every session as day one. ATV remembers.

### The Problem

Every time you start a Copilot session, the AI has no memory of how *your team* writes code. It doesn't know you always wrap errors with `%w`, that you prefer table-driven tests, or that your services use constructor injection. You end up correcting the same patterns session after session.

### The Solution

ATV installs a **continuous learning pipeline** that observes how you code, extracts reusable patterns ("instincts"), and graduates proven ones into discoverable Copilot skills — all automatically.

```text
┌─────────────────────────────────────────────────────────────────────┐
│                    THE LEARNING LOOP                                │
│                                                                     │
│   You code normally                                                 │
│        ↓                                                            │
│   Observer hooks silently capture tool use → .atv/observations.jsonl│
│        ↓                                                            │
│   /learn analyzes observations + git history                        │
│        ↓                                                            │
│   Instincts created with confidence scores → .atv/instincts/       │
│        ↓                                                            │
│   Confidence grows with each session (0.5 → 0.6 → 0.7 → 0.8)      │
│        ↓                                                            │
│   /evolve promotes mature instincts into discoverable skills        │
│        ↓                                                            │
│   .github/skills/learned-*/SKILL.md ← Copilot auto-discovers these │
│        ↓                                                            │
│   Next session: Copilot already knows your patterns                 │
└─────────────────────────────────────────────────────────────────────┘
```

### How It Works

**1. Observer Hooks** — Silent data capture via Copilot CLI's hook system

ATV installs hooks for all 6 Copilot lifecycle events: `sessionStart`, `sessionEnd`, `preToolUse`, `postToolUse`, `userPromptSubmitted`, and `errorOccurred`. A lightweight Node.js observer captures every tool interaction to `.atv/observations.jsonl` — silently, with zero impact on your workflow.

**2. `/learn`** — Pattern extraction from evidence

Analyzes git history, diffs, observations, and existing solutions to find recurring patterns. Each pattern becomes an "instinct" — a small learned behavior with a confidence score that grows through repeated observation.

```yaml
# .atv/instincts/project.yaml
instincts:
  - id: always-wrap-errors
    trigger: "when returning an error from a function"
    behavior: "wrap with fmt.Errorf using %w"
    confidence: 0.85
    observations: 12
```

**3. `/instincts`** — The pattern dashboard

```text
Project Instincts — my-api

  Error Handling (2 instincts)
    ★ always-wrap-errors        0.9  "wrap errors with fmt.Errorf %w"    15 obs
    ● sentinel-errors           0.6  "use sentinel errors for expected"   5 obs

  Testing (1 instinct)
    ★ table-driven-tests        0.85 "use table-driven test pattern"     12 obs

  Legend: ★ ready to evolve (>0.8)  ● active  ○ tentative (<0.5)
```

**4. `/evolve`** — Graduate patterns into permanent skills

When an instinct reaches >0.8 confidence through repeated observation, `/evolve` promotes it into a full SKILL.md at `.github/skills/learned-*/`. These generated skills are auto-discovered by Copilot — your AI assistant now *permanently knows* your team's conventions.

### Why This Matters for Developers

| Without ATV Learning | With ATV Learning |
|---|---|
| Copilot suggests generic patterns every session | Copilot knows *your* patterns from day one |
| You manually correct the same style issues repeatedly | Conventions are enforced automatically via discovered skills |
| New team members have no way to learn "how we do things here" | Instincts are committed to git — the whole team inherits them |
| Tribal knowledge lives in people's heads | Patterns are documented with evidence (commit hashes, observation counts) |
| AI assistance plateaus — it never gets better at *your* project | Every session makes the AI smarter about *your* codebase |

### The Instinct Lifecycle

```text
Session 1:  Pattern first observed               → confidence 0.5  ○ tentative
Session 3:  Pattern seen again, reinforced        → confidence 0.7  ● active
Session 5:  Contradictory evidence found          → confidence 0.55 ● (adjusted)
Session 8:  Pattern consistently confirmed        → confidence 0.85 ★ ready to evolve
            /evolve → .github/skills/learned-*/   → permanent Copilot skill
```

### Key Design Decisions

- **Instincts are committed to git** — the whole team benefits, not just one developer
- **Observations are gitignored** — raw tool use data is ephemeral, instincts are permanent
- **Generated skills use `learned-` prefix** — visually distinct from hand-written skills
- **Cross-platform observer** — Node.js script works on macOS, Linux, and Windows
- **Silent and fast** — observer hooks have 5-second timeout, never block the agent
- **Confidence scoring prevents noise** — only well-established patterns get promoted

---

## De-Slop — The Quality Gate No Other AI Coding Tool Has

> **Every AI coding tool helps you write code faster. ATV is the only one that makes sure the output doesn't *look* like AI wrote it.**

### The Problem

AI coding assistants have a "tell." They over-abstract things that should be simple. They write comments that restate the obvious. They default to purple-to-blue gradients and generic card grids. Experienced developers spot it instantly — and it erodes trust in your codebase.

Every AI tool on the market focuses on generating more code faster. **None of them have a quality gate for detecting and removing the patterns that make AI output look artificial.** The result: codebases fill up with over-engineered abstractions, `// This function handles the logic for...` comments, and template-looking UI. Code review catches bugs — but nobody catches *slop*.

ATV is the first tool to ship a built-in de-slop pass.

### The Solution

`/unslop` runs **three parallel analysis passes** on your recent changes and produces a unified report:

```text
/unslop                          →  Report slop in changed files
/unslop src/components/          →  Scope to a directory
/unslop fix                      →  Auto-apply safe fixes
```

| Pass | What it catches | Example |
|------|----------------|---------|
| **Code Slop** | Over-abstraction, YAGNI violations, nested ternaries, commented-out code | Interface used once → inline it |
| **Comment Rot** | Obvious restatements, AI filler phrases, stale TODOs, inaccurate docs | `// This function handles auth` → delete it |
| **Design Slop** | Generic gradients, template layouts, missing hover states, stock UI | Purple-to-blue default → use brand palette |

### How It Works

```text
┌─────────────────────────────────────────────────────────────────────┐
│                     /unslop PIPELINE                                │
│                                                                     │
│   Determine scope (git diff or explicit path)                       │
│        ↓                                                            │
│   Classify files → which passes to run                              │
│        ↓                                                            │
│   ┌──────────────┬──────────────────┬──────────────────┐            │
│   │  Code Slop   │   Comment Rot    │   Design Slop    │  PARALLEL  │
│   │  Detector    │   Detector       │   Detector       │            │
│   └──────┬───────┴────────┬─────────┴────────┬─────────┘            │
│          └────────────────┼──────────────────┘                      │
│                           ↓                                         │
│   Merge, deduplicate, sort by severity                              │
│        ↓                                                            │
│   De-slop Report (table format)                                     │
│        ↓                                                            │
│   Optional: /unslop fix → auto-apply safe fixes                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Pipeline Integration

`/unslop` is wired into both autonomous pipelines:

- **`/lfg`** — runs `/unslop fix` after `ce-review` autofix, before todo-resolve
- **`/slfg`** — runs `/unslop` (report-only) in the **parallel phase** for free, then `/unslop fix` in the sequential autofix phase

In `/slfg`, the report pass adds **zero wall-clock time** because it runs alongside `ce-review` and browser testing:

```text
/slfg parallel phase:
  ┌─────────────────────────────┐
  │ ce-review (report-only)     │
  │ test-browser                │  ← all three run simultaneously
  │ /unslop (report-only)       │
  └─────────────────────────────┘
```

### AI Filler Phrase Blocklist

`/unslop` maintains a hard-ban list of AI-generated filler that gets flagged immediately:

| Phrase | Why it's slop |
|--------|--------------|
| "This function is responsible for handling..." | Restates the function name |
| "robust and scalable" | Meaningless without evidence |
| "leverages" / "utilizes" | Just say "uses" |
| "seamlessly integrates" | Nothing integrates seamlessly |
| "comprehensive solution" | What solution isn't? |
| "In today's rapidly evolving..." | Marketing copy in a codebase |

### What Makes This Different from Code Review

| `/ce-review` | `/unslop` |
|---|---|
| Asks: "Is this **correct**?" | Asks: "Does this **look human-written**?" |
| Catches bugs, security issues, architecture problems | Catches over-engineering, filler comments, template UI |
| Uses P0-P3 severity (can block merge) | Uses High/Medium/Low (never blocks — slop isn't a bug) |
| 17 specialized reviewer personas | 3 focused passes (code, comments, design) |
| Part of the Review phase | Part of the Reflect phase |

They're complementary — run both. `/ce-review` ensures your code works. `/unslop` ensures it doesn't look like a robot wrote it.

---

## Why Memory Matters

ATV builds **seven layers of memory** — each serving a different timescale and audience:

| Layer | What remembers | Where it lives | Who reads it | Timescale |
|---|---|---|---|---|
| **Learned instincts** | Coding patterns, conventions, style | `.atv/instincts/project.yaml` | `/learn`, `/instincts`, `/evolve` | Grows every session |
| **Evolved skills** | Graduated instincts as full skills | `.github/skills/learned-*/` | Copilot (auto-discovered) | Permanent |
| **Observations** | Raw tool use, commands, errors | `.atv/observations.jsonl` | Observer hooks, `/learn` | Per-session (gitignored) |
| **Institutional knowledge** | Solved problems, gotchas, patterns | `docs/solutions/*.md` | `learnings-researcher` agent during `/ce-plan` | Permanent |
| **Design decisions** | Why we chose approach A over B | `docs/brainstorms/*.md` | `/ce-plan` auto-discovers recent brainstorms | Permanent |
| **Implementation plans** | What to build, acceptance criteria | `docs/plans/*.md` | `/ce-work` reads and checks off items | Per-feature |
| **Install manifest** | What was installed, attempted, skipped | `.atv/install-manifest.json` | `atv-starterkit init` | Per-install |

The **compound memory loop** — three reinforcing cycles:

```text
Cycle 1: Knowledge Compounding (per-PR)
  solve problem → /ce-compound → docs/solutions/
                                       ↓
  future /ce-plan → learnings-researcher → avoids past mistakes

Cycle 2: Pattern Learning (per-session)
  code normally → observer hooks → .atv/observations.jsonl
                                         ↓
  /learn → instincts (confidence scoring) → /evolve → permanent skills

Cycle 3: Team Propagation (per-commit)
  .atv/instincts/project.yaml committed to git
                    ↓
  entire team inherits learned patterns → consistency without style guides
```

> **The combined effect:** Cycle 1 captures *what* you solved. Cycle 2 captures *how* you solve. Cycle 3 shares both with the team. Over weeks, your repo develops a memory that makes every Copilot session more effective than the last.

---

## The Full Sprint

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
                     <code>/learn</code><br />
                     <code>/instincts</code><br />
                     <code>/evolve</code><br />
                     <code>/unslop</code><br />
                     <code>/gstack-retro</code><br />
                     <code>/gstack-learn</code>
              </td>
       </tr>
</table>

> 🛡️ Safety guardrails apply across the whole sprint: `/gstack-careful`, `/gstack-freeze`, `/gstack-guard`, and `/gstack-investigate`.

<details>
<summary><strong>Skill reference by phase</strong></summary>

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
| `/learn` | Extracts coding patterns from recent work into instincts with confidence scoring |
| `/instincts` | Dashboard showing all learned patterns grouped by domain with evolution readiness |
| `/evolve` | Promotes mature instincts (confidence >0.8) into permanent, auto-discovered Copilot skills |
| `/observe` | Focused pattern analysis on a specific domain, file pattern, or question |
| `/unslop` | Unified de-slop pass: code simplification + comment rot detection + design slop check |
| `/gstack-retro` | Team-aware weekly retro with per-person breakdowns |
| `/gstack-learn` | Per-project self-learning infrastructure |

### Safety Guardrails

| Skill | What it does |
|---|---|
| `/gstack-careful` | Warns before `rm -rf`, `DROP TABLE`, force-push |
| `/gstack-freeze` | Restricts edits to one directory while debugging |
| `/gstack-guard` | Careful + Freeze combined |
| `/gstack-investigate` | No fixes without systematic investigation first |

</details>

---

## The Agent Roster

29 specialized agents in `.github/agents/`, invoked by skills during review, planning, learning, and debugging:

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
| **Learning** | `pattern-observer` — analyzes tool use patterns and feeds the learning pipeline |
| **Meta** | `agent-native-reviewer`, `ankane-readme-writer` |
| **Ops** | `lint` |

---

## What Gets Installed

### All 6 Copilot Lifecycle Hooks + Observer System

| # | Hook | File | When it fires |
|---|---|---|---|
| 1 | **System Instructions** | `.github/copilot-instructions.md` | Every Copilot chat |
| 2 | **Setup Steps** | `.github/copilot-setup-steps.yml` | Coding Agent initialization |
| 3 | **MCP Servers** | `.github/copilot-mcp-config.json` | Copilot startup |
| 4 | **Skills** | `.github/skills/*/SKILL.md` | When description matches request |
| 5 | **Agents** | `.github/agents/*.agent.md` | Subagent orchestration |
| 6 | **File Instructions** | `.github/*.instructions.md` | `applyTo` glob matches |
| + | **Observer Hooks** | `.github/hooks/copilot-hooks.json` | Every tool use (silent) |

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
 Screen 1: Stack Packs → Screen 2: Preset → Screen 3: Customize?
        │
        ▼
 Install with structured telemetry:
        │
        ├── ATV scaffold ──► Embedded templates → .github/skills/*/SKILL.md
        │
        ├── Learning pipeline ──► Observer hooks → copilot-hooks.json
        │                         Learning skills → atv-learn, atv-instincts, atv-evolve, atv-observe
        │                         Observer script → .github/hooks/scripts/observe.js
        │                         Instinct storage → .atv/instincts/
        │
        ├── gstack ──► git clone → .gstack/ (staging)
        │               ├── gen:skill-docs → .agents/skills/gstack-*/
        │               └── Copy SKILL.md → .github/skills/gstack-*/
        │
        └── agent-browser ──► npm install -g → agent-browser install (Chrome)
                              └── .github/skills/agent-browser/SKILL.md
        │
        ▼
 Write manifest to .atv/install-manifest.json
```

- `.gstack/` is gitignored — staging area with the full repo and runtime
- `.github/skills/gstack-*/SKILL.md` are lightweight copies Copilot discovers
- All skills at one level deep in `.github/skills/` — Copilot's discovery convention
- Idempotent: re-running skips existing files, merges JSON configs

---

## Development

```bash
go build -o atv-installer .             # build
go test ./...                            # all tests
go test ./pkg/installstate/ -v           # manifest + recommendations tests
go test ./pkg/monitor/ -v                # watcher + drift detection tests
go test ./test/sandbox/ -v               # integration tests (E2E scenarios)

# sandbox test
mkdir /tmp/test && cd /tmp/test
echo '{}' > tsconfig.json && git init
npx atv-starterkit init --guided
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
