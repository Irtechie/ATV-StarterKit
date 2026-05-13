---
name: deepen-brainstorm
description: Enhance a brainstorm document with focused research to validate decisions, surface risks, and add market or technical depth before planning. Use when the user says "deepen brainstorm", "research this brainstorm", or wants a brainstorm hardened before planning.
argument-hint: "[path to brainstorm file]"
---

# Deepen Brainstorm - Research Enhancement Mode

## Quick Start

1. Read the requested brainstorm completely.
2. Extract decisions, assumptions, open questions, and scope boundaries.
3. Research only the items that can change the quality of the later plan.
4. Patch the brainstorm in place with clearly marked research additions.
5. Report the strongest validation, challenged assumptions, and remaining uncertainty.

Use the current session date when dating output. For current external facts, verify against current sources; if browsing or network access is unavailable, mark those facts as unverified instead of guessing.

Unlike `deepen-plan`, which adds implementation detail, this skill adds decision-validation depth:

- Market research validating the chosen approach
- Competitive landscape for the problem space
- Technical feasibility signals
- Risk factors and mitigation strategies
- Prior art and lessons learned
- Strengthened rationale for key decisions

## Input

<brainstorm_path> #$ARGUMENTS </brainstorm_path>

**If the brainstorm path above is empty:**
1. Check for recent brainstorms in `docs/brainstorms/`.
2. Ask the user: "Which brainstorm would you like to deepen? Please provide the path (e.g., `docs/brainstorms/2026-01-15-my-feature-brainstorm.md`)."

Do not proceed until you have a valid brainstorm file path.

## Operating Rules

- Do not write code. This is research and decision validation only.
- Preserve the original brainstorm content. Add evidence layers rather than rewriting the user's decisions.
- Keep research bounded. Prefer 3-7 targeted research questions over exhaustive landscape scanning.
- Use parallel sub-agents only when the platform supports them and the invocation permits delegated research. Otherwise perform the same research locally in batches.
- Use the platform's blocking question tool when available; otherwise ask concise direct questions.
- Do not upload to Proof, post externally, or run network side effects without explicit user confirmation.
- Stage or commit only when the user explicitly asked for a commit.

## Workflow

### 1. Parse and Analyze Brainstorm Structure

Read the brainstorm to identify decisions, approaches chosen, alternatives rejected, and open questions. These are the targets for research validation.

Extract:

- Feature/problem description
- Key decisions and rationale
- Chosen approach and rejected alternatives
- Open questions
- Stated and implied assumptions
- Success criteria
- Scope boundaries
- Technologies or patterns mentioned
- Domain area such as UI, API, infrastructure, workflow, data, or security

Create a research manifest:

```text
Decision 1: [Decision] - Validate: [What to research]
Decision 2: [Decision] - Validate: [What to research]
Open Q 1: [Question] - Research: [What to find]
Assumption 1: [Assumption] - Verify: [How to check]
```

### 2. Run Decision Research

For each key decision, research:

- Who else has solved this problem? What did they choose?
- What are the known failure modes of this approach?
- Are there newer alternatives we might be missing?
- What scale/complexity thresholds make this approach break down?
- What evidence supports or challenges the decision?

For each open question, research:

- Industry consensus, if any
- Tradeoffs between options
- Real-world examples of each option in production
- Data or benchmarks that inform the choice
- A concrete recommendation with confidence level

For each assumption, verify:

- Is this actually true now?
- What conditions make this assumption false?
- Has anyone documented failures from this assumption?
- Is the assumption confirmed, partially true, or challenged?

### 3. Market and Competitive Research

Research the landscape around the problem:

- How do other tools/products solve this?
- What is the current state of the art?
- Are there open-source solutions worth studying?
- What user experience patterns are considered best practice?

Return a landscape summary with 3-5 concrete examples when external research is available.

### 4. Check Learnings and Prior Art

Search for relevant learnings and similar code:

```bash
rg --files docs/solutions
rg -n "[key terms from brainstorm]"
```

If `rg` is unavailable, use the platform's native file search. For each potentially relevant learning, determine:

- Does this learning apply to the brainstorm?
- What specific insight should carry forward?
- If not applicable, why not?

### 5. Discover and Apply Available Skills

Check for skills that could provide domain-specific insights. Search project and user skill locations that exist in the current environment:

```bash
rg --files .github/skills -g "SKILL.md"
```

Also check available global/plugin skill roots exposed by the current platform, such as `~/.copilot/skills`, `~/.codex/skills`, or `~/.claude/skills`.

For each matching skill, apply only the relevant perspective:

- Does the skill validate any decisions?
- Does it challenge any assumptions?
- Does it add a nuance that should affect planning?

### 6. Synthesize and Enhance

Merge research back into the brainstorm. Preserve original decisions but add evidence layers.

Use this format under relevant original sections:

```markdown
### Research Validation

**Evidence supporting this decision:**
- [Finding 1 with source]
- [Finding 2 with source]

**Risks identified:**
- [Risk 1] - Mitigation: [approach]
- [Risk 2] - Mitigation: [approach]

**Competitive context:**
- [How others solve this]

**Confidence level:** High/Medium/Low - [brief justification]
```

For open questions that research answered:

```markdown
### [Original Open Question]

**Research recommendation:** [Answer]

**Evidence:**
- [Source 1]
- [Source 2]

**Confidence:** High/Medium/Low
```

For assumptions:

```markdown
### Assumptions Audit

| Assumption | Status | Evidence |
|-----------|--------|----------|
| [Assumption 1] | Confirmed | [brief evidence] |
| [Assumption 2] | Partially true | [conditions] |
| [Assumption 3] | Challenged | [counter-evidence] |
```

### 7. Add Enhancement Summary

At the top of the brainstorm, add:

```markdown
## Deepening Summary

**Deepened on:** [Date]
**Research questions:** [Count]
**Decisions validated:** [Count]/[Total]
**Open questions resolved:** [Count]/[Total]
**Assumptions verified:** [Count]/[Total]
**New risks identified:** [Count]

### Confidence Assessment
- **Overall approach confidence:** High/Medium/Low
- **Strongest decision:** [which one and why]
- **Weakest decision:** [which one - consider revisiting]
- **Biggest risk:** [what could go wrong]

### Key Research Findings
1. [Most important finding]
2. [Second most important]
3. [Third most important]
```

### 8. Update Brainstorm File

Write the enhanced brainstorm:

- Preserve original filename
- Keep all original content intact
- Mark research additions with `### Research Validation`
- Move resolved open questions to a "Resolved Questions" section

## Quality Checks

- [ ] All original decisions preserved
- [ ] Research clearly attributed with sources
- [ ] Confidence levels are honest
- [ ] Challenged assumptions are flagged prominently
- [ ] No implementation details crept in
- [ ] Enhancement summary accurately reflects findings

## Success Criteria

- The brainstorm still contains every original decision and scope boundary.
- Research additions distinguish verified facts from inference.
- The top summary names the strongest decision, weakest decision, biggest risk, and remaining open questions.
- The document is ready to feed into `ce-plan` or `kanban-plan`.

## Post-Enhancement Options

After writing the enhanced brainstorm, ask what to do next:

1. **View changes** - Show what research added
2. **Revisit challenged decisions** - Discuss decisions that research challenged
3. **Proceed to planning** - Run `/ce-plan` or `/kanban-plan` with the brainstorm path
4. **Deepen further** - Run another research pass on specific decisions
5. **Share to Proof** - Upload to Proof for collaborative review, after explicit confirmation
