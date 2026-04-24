---
title: "feat: Compound Engineering Latest Update"
type: feat
status: active
date: 2026-04-02
origin: docs/brainstorms/2026-04-02-compound-engineering-latest-update-brainstorm.md
---

# feat: Compound Engineering Latest Update

## Overview

Update the ATV installer's embedded compound-engineering (CE) content to match the latest upstream plugin referenced in [ATV-StarterKit PR #7](https://github.com/All-The-Vibes/ATV-StarterKit/pull/7). This adds 20 new agents, rewrites 27 existing agents, introduces 3 new skills, and updates ~20 existing skills. All changes flow through the existing `//go:embed all:templates` pattern — no architectural changes.

## Problem Statement / Motivation

The CE plugin has undergone a major content revision: new specialized reviewers (adversarial, API contract, performance, reliability, testing, etc.), restructured existing agents with JSON-output code-review personas, and expanded skills with multi-file directories (references/, scripts/, assets/). The current ATV installer embeds the prior version. Users who install via `atv init` get outdated CE content.

## Proposed Solution

Cherry-pick PR #7 files into `pkg/scaffold/templates/`, register 3 new skill directories in `catalog.go`, add TUI category entries in `categories.go`, and rebuild. This follows the exact pattern used for the existing 28 agents and 14 skills — manual curation into templates, explicit registration in the catalog. No new layers, no new presets, no walker changes.

## Technical Considerations

- **`skillComponents()` handles multi-file skills already** — uses `fs.WalkDir` recursively, so placing `references/`, `scripts/`, `assets/` subdirectories inside a skill directory works with zero code changes to the walker.
- **All new agents are universal** — none are added to `stackAgents`, so they flow through `LayerUniversalAgents` in all three presets.
- **`agent-browser` special handling preserved** — PR #7 adds references/ and templates/ subdirs to agent-browser, which is a content update to the existing skill. The `IncludeAgentBrowser` / Full-preset-only tracking remains unchanged.
- **Binary size increase** — 20 new agent files + 3 new skill directories + expanded existing files will increase the embedded content. Acceptable for a Go binary.
- **No runtime behavior changes** — this is purely embedded content. `BuildCatalog`, `BuildFilteredCatalogForPacks`, `WriteAll` flow unchanged.
- **Beta content excluded** — `ce-work-beta` and `.bak` backup files from PR #7 are not copied.

## System-Wide Impact

- **Interaction graph**: `BuildCatalog()` → `skills()` → `skillComponents(nil)` walks all embedded skills. `BuildFilteredCatalogForPacks()` → `skillComponents(selectedSkillDirs)` uses the `coreSkillDirectories` and `orchestratorSkillDirectories` slices to filter. New entries in these slices are picked up automatically. Agent flow is identical — `agentsForPacks()` walks `templates/agents/` and filters by `stackAgents` map (unchanged).
- **Error propagation**: `mustRead()` panics on missing embedded files. If a file listed in the brainstorm is not actually placed in templates/, the binary won't compile (embed will succeed but `mustRead` won't be called for skill files — `fs.WalkDir` discovers them). Risk is low.
- **State lifecycle risks**: None — this is compile-time content, not runtime state.
- **API surface parity**: `categories.go` TUI display adds entries for new skills. Wizard flow, preset selections, and gstack integration are unchanged.
- **Integration test scenarios**: (1) `go build ./...` succeeds with new templates embedded. (2) `go test ./...` passes — existing scaffold tests verify the walker and catalog logic. (3) Manual `atv init` produces all new agents and skills in `.github/`.

## Acceptance Criteria

- [ ] 20 new `.agent.md` files in `pkg/scaffold/templates/agents/`
- [ ] 27 existing `.agent.md` files updated with PR #7 content
- [ ] 3 new skill directories in `pkg/scaffold/templates/skills/` (ce-ideate, ce-compound-refresh, claude-permissions-optimizer)
- [ ] ~20 existing skill directories updated (including new subdirectories where PR #7 adds them)
- [ ] `coreSkillDirectories` in `catalog.go` includes `"ce-ideate"`, `"ce-compound-refresh"`
- [ ] `orchestratorSkillDirectories` in `catalog.go` includes `"claude-permissions-optimizer"`
- [ ] `atvCategoryMapping` in `categories.go` includes entries for all 3 new skills
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] No changes to `presets.go`, `wizard.go`, `stackAgents` map, or `skillComponents()` logic
- [ ] `ce-work-beta` and `.bak` files excluded

## Success Metrics

- All existing tests pass (scaffold_test.go, wizard tests, category tests)
- `atv init` generates correct agent and skill files matching PR #7 content
- Binary compiles without errors with expanded embedded templates

## Dependencies & Risks

- **Source dependency**: PR #7 content must be fetched from GitHub (SHA: 7311b001480dfa6f36b16573930404bbf179d70f). PR is draft status but content is stable.
- **Risk: file count accuracy** — the brainstorm lists 20 new + 27 modified agents based on PR diff analysis. Actual count should be verified during implementation.
- **Risk: multi-file skill subdirs** — if any new subdirectory structure doesn't match what `fs.WalkDir` expects (e.g., symlinks, empty dirs), the walker skips it silently. Verify all files are regular files with content.

## Implementation Phases

### Phase 1: Template Content (bulk of work)

**1a. New agent files (20 files)**

Fetch from PR #7 and place in `pkg/scaffold/templates/agents/`:

| File | Source PR path |
|------|---------------|
| adversarial-document-reviewer.agent.md | .github/agents/ |
| adversarial-reviewer.agent.md | .github/agents/ |
| api-contract-reviewer.agent.md | .github/agents/ |
| cli-agent-readiness-reviewer.agent.md | .github/agents/ |
| cli-readiness-reviewer.agent.md | .github/agents/ |
| coherence-reviewer.agent.md | .github/agents/ |
| correctness-reviewer.agent.md | .github/agents/ |
| data-migrations-reviewer.agent.md | .github/agents/ |
| design-lens-reviewer.agent.md | .github/agents/ |
| feasibility-reviewer.agent.md | .github/agents/ |
| issue-intelligence-analyst.agent.md | .github/agents/ |
| maintainability-reviewer.agent.md | .github/agents/ |
| performance-reviewer.agent.md | .github/agents/ |
| previous-comments-reviewer.agent.md | .github/agents/ |
| product-lens-reviewer.agent.md | .github/agents/ |
| project-standards-reviewer.agent.md | .github/agents/ |
| reliability-reviewer.agent.md | .github/agents/ |
| scope-guardian-reviewer.agent.md | .github/agents/ |
| security-lens-reviewer.agent.md | .github/agents/ |
| security-reviewer.agent.md | .github/agents/ |
| testing-reviewer.agent.md | .github/agents/ |

**1b. Update existing agent files (27 files)**

Replace content of all 28 existing agents in `pkg/scaffold/templates/agents/` with PR #7 versions (where changed — verify diff).

**1c. New skill directories (3 dirs)**

Create in `pkg/scaffold/templates/skills/`:
- `ce-ideate/SKILL.md`
- `ce-compound-refresh/SKILL.md` + `references/` + `assets/` subdirs
- `claude-permissions-optimizer/SKILL.md` + `scripts/` subdir

**1d. Update existing skill directories (~20 dirs)**

Replace/update content in existing skill directories. Key updates:
- `ce-brainstorm/` — major SKILL.md rewrite
- `ce-plan/` — major SKILL.md expansion
- `ce-review/` — rewrite + add `references/` subdir
- `ce-work/` — rewrite
- `ce-compound/` — rewrite + add `references/`, `assets/` subdirs
- `document-review/` — rewrite + add `references/` subdir
- `feature-video/` — major update
- `agent-browser/` — major update + add `references/`, `templates/` subdirs (note: agent-browser has special TUI handling but skill content flows normally through `skillComponents()`)

### Phase 2: Registration (small code changes)

**2a. catalog.go** — 3 lines added

```go
var coreSkillDirectories = []string{
    "brainstorming",
    "ce-brainstorm",
    "ce-compound",
    "ce-compound-refresh",  // NEW
    "ce-ideate",            // NEW
    "ce-plan",
    "ce-review",
    "ce-work",
    "deepen-plan",
    "document-review",
    "setup",
}

var orchestratorSkillDirectories = []string{
    "claude-permissions-optimizer",  // NEW
    "feature-video",
    "lfg",
    "resolve_todo_parallel",
    "slfg",
    "test-browser",
}
```

**2b. categories.go** — 3 entries added to `atvCategoryMapping`

```go
gstack.CategoryPlanning: {
    {Label: "Brainstorming — explore what to build", Key: "core-skills:brainstorming", Source: "atv"},
    {Label: "CE Ideate — structured idea exploration", Key: "core-skills:ce-ideate", Source: "atv"},  // NEW
    {Label: "Plan — turn ideas into an implementation plan", Key: "core-skills:ce-plan", Source: "atv"},
    {Label: "Deepen Plan — parallel research to harden the plan", Key: "core-skills:deepen-plan", Source: "atv"},
},
gstack.CategoryShipping: {
    {Label: "CE Work — execute plans with quality checks", Key: "core-skills:ce-work", Source: "atv"},
    {Label: "LFG — full autonomous pipeline", Key: "orchestrators:lfg", Source: "atv"},
    {Label: "SLFG — swarm mode parallel execution", Key: "orchestrators:slfg", Source: "atv"},
    {Label: "CE Compound — document solutions", Key: "core-skills:ce-compound", Source: "atv"},
    {Label: "CE Compound Refresh — refresh documented solutions", Key: "core-skills:ce-compound-refresh", Source: "atv"},  // NEW
},
```

`claude-permissions-optimizer` mapped to an appropriate category (Safety or Shipping):

```go
gstack.CategorySafety: {  // or add new section if CategorySafety doesn't exist yet
    {Label: "Claude Permissions Optimizer — optimize tool permissions", Key: "orchestrators:claude-permissions-optimizer", Source: "atv"},
},
```

### Phase 3: Verification

- `go build ./...`
- `go test ./...`
- Manual spot-check: `atv init` in a temp directory, verify new agents and skills appear

## Files NOT Modified

| File | Reason |
|------|--------|
| `pkg/tui/presets.go` | All presets already include LayerCoreSkills + LayerOrchestrators |
| `pkg/tui/wizard.go` | No new layer constants needed |
| `pkg/scaffold/catalog.go` (skillComponents) | fs.WalkDir already recursive |
| `pkg/scaffold/catalog.go` (stackAgents) | No new stack-specific agents |
| `pkg/scaffold/scaffold.go` | Write logic unchanged |

## Sources & References

- **Origin brainstorm:** [docs/brainstorms/2026-04-02-compound-engineering-latest-update-brainstorm.md](../brainstorms/2026-04-02-compound-engineering-latest-update-brainstorm.md)
- **Source PR:** [All-The-Vibes/ATV-StarterKit#7](https://github.com/All-The-Vibes/ATV-StarterKit/pull/7) (SHA: 7311b001480dfa6f36b16573930404bbf179d70f)
- **Current catalog.go:** [pkg/scaffold/catalog.go](../../pkg/scaffold/catalog.go)
- **Current categories.go:** [pkg/tui/categories.go](../../pkg/tui/categories.go)
