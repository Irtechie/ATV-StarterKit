---
name: deepen-brainstorm
description: Enhance a brainstorm document with parallel research agents to validate decisions, surface risks, and add market/technical depth before planning
---

## Arguments
[path to brainstorm file]

# Deepen Brainstorm - Research Enhancement Mode

## Introduction

**Note: The current year is 2026.** Use this when searching for recent documentation and best practices.

This command takes an existing brainstorm document (from `/ce-brainstorm`) and enhances it with parallel research agents. Unlike `/deepen-plan` which adds implementation depth, this adds **decision-validation depth**:
- Market research validating the chosen approach
- Competitive landscape for the problem space
- Technical feasibility signals
- Risk factors and mitigation strategies
- Prior art and lessons learned
- Strengthened rationale for key decisions

The result is a brainstorm with battle-tested decisions ready for confident planning.

## Brainstorm File

<brainstorm_path> #$ARGUMENTS </brainstorm_path>

**If the brainstorm path above is empty:**
1. Check for recent brainstorms: `ls -la docs/brainstorms/`
2. Ask the user: "Which brainstorm would you like to deepen? Please provide the path (e.g., `docs/brainstorms/2026-01-15-my-feature-brainstorm.md`)."

Do not proceed until you have a valid brainstorm file path.

## Main Tasks

### 1. Parse and Analyze Brainstorm Structure

<thinking>
Read the brainstorm to identify decisions, approaches chosen, alternatives rejected, and open questions. These are the targets for research validation.
</thinking>

**Read the brainstorm file and extract:**
- [ ] Feature/problem description (the WHAT)
- [ ] Key decisions made and their rationale
- [ ] Chosen approach and why
- [ ] Rejected alternatives and why
- [ ] Open questions (highest priority for research)
- [ ] Assumptions stated or implied
- [ ] Success criteria defined
- [ ] Scope boundaries (in/out)
- [ ] Technologies or patterns mentioned
- [ ] Domain area (UI, API, infrastructure, workflow, etc.)

**Create a research manifest:**
```
Decision 1: [Decision] - Validate: [What to research]
Decision 2: [Decision] - Validate: [What to research]
Open Q 1: [Question] - Research: [What to find]
Assumption 1: [Assumption] - Verify: [How to check]
```

### 2. Launch Parallel Research Agents

<thinking>
For each decision, open question, and assumption, spawn dedicated research agents. The goal is to validate or challenge each one with external evidence.
</thinking>

**For each key decision, spawn a validation agent:**

```
Task explore: "Research whether [chosen approach] is the right choice for [problem context].
Find:
- Who else has solved this problem? What did they choose?
- What are the known failure modes of this approach?
- Are there newer (2024-2026) alternatives we might be missing?
- What scale/complexity thresholds make this approach break down?
Return: Evidence supporting OR challenging this decision."
```

**For each open question, spawn a research agent:**

```
Task explore: "Research: [open question from brainstorm].
Find:
- Industry consensus (if any)
- Tradeoffs between options
- Real-world examples of each option in production
- Data or benchmarks that inform the choice
Return: Concrete recommendation with evidence."
```

**For each assumption, spawn a verification agent:**

```
Task explore: "Verify assumption: [assumption from brainstorm].
Find:
- Is this actually true in 2026?
- What conditions make this assumption false?
- Has anyone documented failures from this assumption?
Return: Confirmed/Challenged with evidence."
```

**Launch ALL agents in PARALLEL.**

### 3. Market & Competitive Research

<thinking>
Understand the broader landscape around this feature/problem. Even for internal tools, someone has likely solved a similar problem.
</thinking>

**Spawn market research agents:**

```
Task explore: "Research the competitive/market landscape for: [problem being solved].
Find:
- How do other tools/products solve this?
- What's the current state of the art (2024-2026)?
- Are there open-source solutions worth studying?
- What user experience patterns are considered best practice?
Return: Landscape summary with 3-5 concrete examples."
```

**Use WebSearch for current context:**

Search for recent articles, blog posts, and documentation related to the brainstorm's problem domain.

### 4. Check Learnings & Prior Art

<thinking>
Check institutional knowledge for relevant past experience.
</thinking>

**Search for relevant learnings:**

```bash
# Project learnings from /ce-compound
find docs/solutions -name "*.md" -type f 2>/dev/null

# Check if similar features exist in the codebase
grep -r "[key terms from brainstorm]" --include="*.py" --include="*.ts" -l
```

**For each potentially relevant learning, spawn a sub-agent:**

```
Task explore: "Read this learning file and determine if it applies to our brainstorm:

Learning: [path]
Brainstorm context: [brief summary of what we're building]

If relevant: What specific insight should we carry forward?
If not: Say 'Not applicable' with brief reason."
```

### 5. Discover and Apply Available Skills

<thinking>
Check for skills that could provide domain-specific insights on the brainstorm's decisions.
</thinking>

**Discover skills from all sources:**

```bash
# Project skills
ls .github/skills/ 2>/dev/null

# User global skills
ls ~/.copilot/skills/ 2>/dev/null

# Plugin skills
find ~/.copilot/plugins/cache -type d -name "skills" 2>/dev/null
```

**For each skill that matches the brainstorm's domain, spawn a sub-agent:**

```
Task general-purpose: "You have the [skill-name] skill at [path].
Read its SKILL.md and apply its perspective to this brainstorm:

[brainstorm content]

Focus on: Does this skill's domain knowledge validate, challenge, or add nuance to any decisions?
Return: Skill-informed insights only (skip if nothing relevant)."
```

### 6. Synthesize and Enhance

<thinking>
Merge all research back into the brainstorm. Preserve original decisions but add evidence layers.
</thinking>

**Collect outputs from all agents and organize by brainstorm section:**

**Enhancement format:**

```markdown
## [Original Section]

[Original content preserved]

### 📊 Research Validation

**Evidence supporting this decision:**
- [Finding 1 with source]
- [Finding 2 with source]

**Risks identified:**
- [Risk 1] — Mitigation: [approach]
- [Risk 2] — Mitigation: [approach]

**Competitive context:**
- [How others solve this]

**Confidence level:** High/Medium/Low — [brief justification]
```

**For Open Questions that research answered:**

```markdown
### [Original Open Question]

**Research recommendation:** [Answer]

**Evidence:**
- [Source 1]
- [Source 2]

**Confidence:** High/Medium/Low
```

**For Assumptions that were verified/challenged:**

```markdown
### Assumptions Audit

| Assumption | Status | Evidence |
|-----------|--------|----------|
| [Assumption 1] | ✅ Confirmed | [brief evidence] |
| [Assumption 2] | ⚠️ Partially true | [conditions] |
| [Assumption 3] | ❌ Challenged | [counter-evidence] |
```

### 7. Add Enhancement Summary

At the top of the brainstorm, add:

```markdown
## Deepening Summary

**Deepened on:** [Date]
**Research agents used:** [Count]
**Decisions validated:** [Count]/[Total]
**Open questions resolved:** [Count]/[Total]
**Assumptions verified:** [Count]/[Total]
**New risks identified:** [Count]

### Confidence Assessment
- **Overall approach confidence:** High/Medium/Low
- **Strongest decision:** [which one and why]
- **Weakest decision:** [which one — consider revisiting]
- **Biggest risk:** [what could go wrong]

### Key Research Findings
1. [Most important finding]
2. [Second most important]
3. [Third most important]
```

### 8. Update Brainstorm File

**Write the enhanced brainstorm:**
- Preserve original filename
- All original content intact
- Research additions clearly marked with `### 📊 Research Validation` headers
- Open questions resolved are moved to a "Resolved Questions" section

## Quality Checks

Before finalizing:
- [ ] All original decisions preserved (never overwrite user choices)
- [ ] Research clearly attributed with sources
- [ ] Confidence levels are honest (don't inflate)
- [ ] Challenged assumptions are flagged prominently
- [ ] No implementation details crept in (that's for /ce-plan)
- [ ] Enhancement summary accurately reflects findings

## Post-Enhancement Options

After writing the enhanced brainstorm, use the **AskUserQuestion tool**:

**Question:** "Brainstorm deepened at `[path]`. What would you like to do next?"

**Options:**
1. **View changes** - Show what research added
2. **Revisit challenged decisions** - Discuss decisions that research challenged
3. **Proceed to planning** - Run `/ce-plan` (will auto-detect this brainstorm)
4. **Deepen further** - Run another research pass on specific decisions
5. **Share to Proof** - Upload to Proof for collaborative review

Based on selection:
- **View changes** → Show summary of additions
- **Revisit challenged decisions** → Re-enter brainstorm dialogue focused on weak decisions
- **Proceed to planning** → Run `/ce-plan` with brainstorm path
- **Deepen further** → Ask which decisions need more research, re-run those agents
- **Share to Proof** →
  ```bash
  CONTENT=$(cat [brainstorm_path])
  TITLE="Brainstorm (Deepened): <topic>"
  RESPONSE=$(curl -s -X POST https://www.proofeditor.ai/share/markdown \
    -H "Content-Type: application/json" \
    -d "$(jq -n --arg title "$TITLE" --arg markdown "$CONTENT" --arg by "ai:compound" '{title: $title, markdown: $markdown, by: $by}')")
  PROOF_URL=$(echo "$RESPONSE" | jq -r '.tokenUrl')
  ```
  Display: `View & collaborate in Proof: <PROOF_URL>`

NEVER CODE! Just research and validate decisions.
