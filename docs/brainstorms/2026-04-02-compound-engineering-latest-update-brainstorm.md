---
date: 2026-04-02
topic: compound-engineering-latest-update
---

# Integrating Latest Compound Engineering Plugin into ATV Installer

## What We're Building

Update the ATV installer's embedded compound-engineering (CE) content to match the latest upstream plugin (referenced in [ATV-StarterKit PR #7](https://github.com/All-The-Vibes/ATV-StarterKit/pull/7)). The PR adds 20 new agents, rewrites ~27 existing agents, introduces 4 new skills (some with multi-file directories), and updates ~20 existing skills with major content improvements.

CE content is embedded into the Go binary at compile time via `//go:embed all:templates` in `pkg/scaffold/catalog.go`. This is fundamentally different from gstack, which clones and builds at user install time. For CE, the update workflow is: copy updated files into `pkg/scaffold/templates/`, register new entries in catalog.go and the TUI category system, then rebuild the binary. No runtime network calls are involved.

## Why This Approach

We explored three directions:

### Approach A: Manual Cherry-Pick into Templates (chosen)
Copy PR #7 files into `pkg/scaffold/templates/`, update `coreSkillDirectories` and `atvCategoryMapping` registrations, rebuild binary.
- **Pros:** Follows the existing embedded template pattern exactly. No new tooling, no build complexity. All CE content is already managed this way (28 agents, 14 skills in templates/ today). Full curation control — skip beta/experimental content per directory.
- **Cons:** Manual effort per sync. Future CE updates require the same copy-and-register process.

### Approach B: Git Subtree/Submodule (rejected)
Point `pkg/scaffold/templates/.github/` at the upstream CE plugin repo as a git subtree or submodule.
- **Pros:** Automates upstream tracking, `git subtree pull` for future syncs.
- **Cons:** Complicates the build (contributors must init submodules), all-or-nothing from upstream, conflicts when ATV customizes CE files. Breaks the existing pattern where each template file is explicitly registered.

### Approach C: Sync Script (rejected)
Build a script to clone the CE repo, diff against current templates, and generate catalog.go updates.
- **Pros:** Reproducible, automates the manual copy step.
- **Cons:** New tooling to maintain. YAGNI — this is a one-time sync until the next major CE release.

Manual cherry-pick was chosen because it's exactly how CE content is already managed. The existing 28 agents and 14 skills in `templates/` were placed there manually and registered in `catalog.go`. This update follows the same pattern at a larger scale.

## Key Decisions

- **Multi-file skill directories work with no walker changes**: `skillComponents()` in `catalog.go` already uses `fs.WalkDir` to recursively traverse all files under each skill directory — not just SKILL.md. Placing subdirectories (references/, scripts/, assets/) inside `templates/skills/<skill-name>/` is sufficient. The walker processes every non-directory entry and maps it to `.github/skills/<skill-name>/...`. No code change needed to `skillComponents()`.

- **New agents → templates/agents/ (universal)**: All 20 new agents go into `templates/agents/`. They are not added to the `stackAgents` map in `catalog.go`, so they are classified as universal agents. They install when the `universal-agents` layer is selected, which is included in all three presets (Starter, Pro, Full).

- **New CE skills → `coreSkillDirectories`**: `ce-ideate` and `ce-compound-refresh` are added to `coreSkillDirectories` in `catalog.go`. This makes them eligible when the `core-skills` layer is filtered in `BuildFilteredCatalogForPacks()`. They also need entries in `atvCategoryMapping` in `categories.go` so they appear in the guided TUI wizard — `ce-ideate` maps to the Planning category, `ce-compound-refresh` maps to the Shipping category (alongside existing `ce-compound`).

- **`claude-permissions-optimizer` → `orchestratorSkillDirectories`**: Rather than introducing a new `utility-skills` layer and a new `LayerUtilitySkills` constant (which would require changes to all three presets in `presets.go`, `BuildFilteredCatalogForPacks()` in `catalog.go`, and `AllLayers()` in `wizard.go`), this skill fits into `orchestratorSkillDirectories` — the existing bucket for non-core workflow tools like `feature-video`, `lfg`, `slfg`, `resolve_todo_parallel`, and `test-browser`. It needs an `atvCategoryMapping` entry under a suitable gstack category.

- **agent-browser has special handling**: `agent-browser` is not a normal skill in the TUI. It has dedicated tracking: `IncludeAgentBrowser` field on `WizardResult`, separate logic in `presets.go` (only included in Full preset), and its own key `"agent-browser"` in `categories.go`. PR #7's agent-browser updates (new references/, templates/) are a content update to the existing skill, not a new registration. The special handling remains unchanged.

- **Preset impact — content only, no structural changes**: All three presets (Starter, Pro, Full) already include `LayerCoreSkills`, `LayerOrchestrators`, `LayerUniversalAgents`, and `LayerStackAgents`. New agents and skills flow through existing layers with no preset structure changes needed. The presets' `ATVLayers` slices remain identical.

- **Beta skills excluded**: `ce-work-beta` is not copied into templates/ until stable.

- **Frontmatter updates applied**: All 27 modified agents get their content replaced with the PR #7 version. This includes `user-invocable: true` frontmatter additions and major content rewrites to agents like `dhh-rails-reviewer` (JSON-output persona format), `kieran-*-reviewer`, `pr-comment-resolver`, and `repo-research-analyst`.

- **`stackAgents` map preserved**: The existing 10 entries in `stackAgents` (mapping Rails, Python, TypeScript agents to their stacks) are unchanged. No new agents are stack-specific.

## Scope of Changes

### New Agents to Add (20 files → templates/agents/)
| Agent | Lines | Purpose |
|-------|-------|---------|
| adversarial-document-reviewer | 86 | Challenges plan premises and assumptions |
| adversarial-reviewer | 103 | Constructs failure scenarios for code |
| api-contract-reviewer | 44 | Reviews for breaking API contract changes |
| cli-agent-readiness-reviewer | 441 | CLI agent-readiness review (standalone) |
| cli-readiness-reviewer | 66 | CLI readiness (code-review persona) |
| coherence-reviewer | 36 | Document internal consistency |
| correctness-reviewer | 44 | Logic errors and edge cases |
| data-migrations-reviewer | 48 | Migration safety and data integrity |
| design-lens-reviewer | 43 | Missing design decisions in plans |
| feasibility-reviewer | 39 | Technical feasibility of plans |
| issue-intelligence-analyst | 229 | GitHub issue theme analysis |
| maintainability-reviewer | 44 | Code maintainability |
| performance-reviewer | 46 | Performance patterns in code |
| previous-comments-reviewer | 60 | PR comment history context |
| product-lens-reviewer | 67 | Product/business lens on plans |
| project-standards-reviewer | 76 | Project-specific convention enforcement |
| reliability-reviewer | 44 | Error handling and resilience |
| scope-guardian-reviewer | 51 | Scope-goal alignment in plans |
| security-lens-reviewer | 35 | Security implications in plans |
| security-reviewer | 46 | Security patterns in code |
| testing-reviewer | 44 | Test coverage and quality |

### Modified Agents to Update (27 files)
Most get `user-invocable: true` added. Major rewrites:
- `agent-native-reviewer` — restructured with triage step, stack-specific search strategies, noun test
- `dhh-rails-reviewer` — converted to JSON-output code-review persona
- `julik-frontend-races-reviewer` — condensed from 219 to 45 lines, converted to JSON output
- `kieran-*-reviewer` (Python, Rails, TypeScript) — all converted to JSON-output personas
- `pr-comment-resolver` — major expansion with structured resolution workflow
- `repo-research-analyst` — major expansion with structured research methodology
- `schema-drift-detector` — enhanced with cross-reference checking
- `learnings-researcher` — refactored search approach
- `spec-flow-analyzer` — significant enhancements

### New Skills to Add
| Skill | Directory Slice | Multi-file? | Notes |
|-------|----------------|-------------|-------|
| ce-ideate | coreSkillDirectories | SKILL.md only | New CE pipeline step for ideation |
| ce-compound-refresh | coreSkillDirectories | Yes (references/, assets/) | Refreshed compound learning |
| claude-permissions-optimizer | orchestratorSkillDirectories | Yes (scripts/) | Specialized workflow tooling |

### Modified Skills to Update (~20 files)
Major rewrites:
- `ce-brainstorm` — 294 additions, substantially reworked
- `ce-plan` — 776 additions, major expansion
- `ce-review` — 500 additions + new references/ directory
- `ce-work` — 174 additions, reworked
- `ce-compound` — 221 additions + new references/ and assets/ directories
- `document-review` — 264 additions + new references/ directory
- `agent-browser` — 556 additions + new references/ and templates/ directories
- `frontend-design` — 243 additions
- `feature-video` — 209 additions

### Implementation Changes Required

**catalog.go:**
1. Add `"ce-ideate"`, `"ce-compound-refresh"` to `coreSkillDirectories` slice
2. Add `"claude-permissions-optimizer"` to `orchestratorSkillDirectories` slice
3. No changes to `stackAgents` map — all new agents are universal
4. No changes to `skillComponents()` — the `fs.WalkDir` walker already handles recursive subdirectories

**categories.go:**
5. Add `ce-ideate` to `atvCategoryMapping[gstack.CategoryPlanning]` skills list
6. Add `ce-compound-refresh` to `atvCategoryMapping[gstack.CategoryShipping]` skills list (or CategoryRetrospective)
7. Add `claude-permissions-optimizer` to an appropriate `atvCategoryMapping` category
8. No changes to `InfraLayers` — agents flow through existing `LayerUniversalAgents`

**presets.go:**
9. No structural changes — all three presets already include `LayerCoreSkills`, `LayerOrchestrators`, `LayerUniversalAgents`

**wizard.go:**
10. No changes — no new layer constants needed

**templates/agents/:**
11. Copy 20 new `.agent.md` files
12. Replace content of 27 existing `.agent.md` files

**templates/skills/:**
13. Copy new skill directories (`ce-ideate/`, `ce-compound-refresh/`, `claude-permissions-optimizer/`)
14. Replace/update content of ~20 existing skill directories (including adding references/, scripts/, assets/ subdirectories where PR #7 introduces them)

### Files to Skip
- `ce-work-beta` — beta, not stable
- `.github/copilot-mcp-config.json.bak.*` — backup file from PR author's local state

## Resolved Questions

- **Include all 48 agents or curate?** Include all. New review-persona agents are universal (no `stackAgents` entries), so they flow through the existing `universal-agents` layer in all three presets.
- **How to handle multi-file skill directories?** Just place the files. `skillComponents()` already uses `fs.WalkDir` to traverse all files recursively — no code change needed. `//go:embed all:templates` embeds the full tree.
- **Include beta skills?** No. `ce-work-beta` is excluded until stable.
- **Where does claude-permissions-optimizer go?** Into `orchestratorSkillDirectories` — the existing bucket for non-core workflow tools. Adding a new `utility-skills` layer would require changes across `catalog.go`, `wizard.go`, and all three presets in `presets.go`, which violates YAGNI for a single skill.
- **Do presets need changes?** No. All three presets (Starter, Pro, Full) already include `LayerCoreSkills`, `LayerOrchestrators`, `LayerUniversalAgents`, and `LayerStackAgents`. New content flows through existing layers.
- **Does agent-browser registration change?** No. `agent-browser` has special handling in `WizardResult.IncludeAgentBrowser`, preset-level tracking (Full only), and its own key in `atvCategoryMapping`. PR #7's changes are content-only updates to the existing skill directory.
- **Merge or replace?** PR #7 is used as reference only (not merged). Files are cherry-picked into templates/ with curation.

## Next Steps

→ `/ce-plan` for implementation details
