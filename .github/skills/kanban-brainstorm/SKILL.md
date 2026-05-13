---
name: kanban-brainstorm
description: 'Research-first brainstorming for vertical-slice work. Runs market and landscape research before asking product questions, so questions are sharper and approaches are grounded in real prior art. Use when the user says ''kanban brainstorm'', ''research-first brainstorm'', ''brainstorm with research'', or ''brainstorm before kanban-plan''. Pick this skill when prior art or competitive landscape is expected to materially change framing, OR when the brainstorm output is intended to feed `/kanban-plan` (vertical slices). Otherwise prefer `/ce-brainstorm`. Keeps the depth of `/ce-brainstorm` (product pressure test, requirements doc, document review) but inverts the order so research informs the conversation.'
argument-hint: "[feature idea or problem to explore]"
---

# Kanban Brainstorm — Research-First Requirements

**Note: The current year is 2026.** Use this when dating requirements documents.

`kanban-brainstorm` answers **WHAT** to build by running market and landscape research **before** asking the user product questions. It is a peer of `/ce-brainstorm` with the same depth, but a different order: research first, then conversation.

This pairs naturally with `/kanban-plan` (vertical-slice decomposition), but the requirements doc it produces is also valid input for `/ce-plan`.

This skill does not implement code. It explores, validates, clarifies, and documents decisions for later planning or execution.

**IMPORTANT: All file references in generated documents must use repo-relative paths (e.g., `src/models/user.rb`), never absolute paths. Absolute paths break portability across machines, worktrees, and teammates.**

## When to Pick `kanban-brainstorm` vs `ce-brainstorm`

| Situation | Pick |
|---|---|
| Design space is well known; conversation is the bottleneck | `ce-brainstorm` |
| Prior art / competitive landscape is **likely to change framing** | `kanban-brainstorm` |
| Output will feed `/kanban-plan` (vertical slices) | `kanban-brainstorm` |
| Output will feed `/ce-plan` and research is unlikely to move you | `ce-brainstorm` |
| Existing brainstorm doc just needs research enrichment | `/deepen-brainstorm` |

`kanban-brainstorm` does **bounded framing research before product decisions**. `/deepen-brainstorm` does **post-doc enrichment and challenge**. Don't pick `kanban-brainstorm` just because research feels good — pick it when research is expected to move framing or you're heading to `/kanban-plan`.

## Core Principles

1. **Research before questions** — Do not interview the user before scanning prior art, competitive landscape, and applicable repo patterns. Questions asked without that context tend to ratify the user's first framing instead of testing it.
2. **Evidence beats intuition** — Every decision in the requirements doc should have either a research citation or an explicit "no evidence — assumption" tag.
3. **Be a thinking partner** — Bring alternatives, challenge assumptions, and surface what-ifs. Don't just extract requirements.
4. **Resolve product decisions here** — User-facing behavior, scope boundaries, and success criteria belong in this workflow. Detailed implementation belongs in planning.
5. **Right-size the artifact** — Match ceremony to scope. Lightweight work gets a compact doc; deep work gets a fuller one. Do not pad sections that add no value.
6. **Apply YAGNI to carrying cost, not coding effort** — Prefer the simplest approach that delivers meaningful value. Avoid speculative complexity, but include low-cost polish that compounds.
7. **Keep implementation out of the requirements doc by default** — Do not include libraries, schemas, endpoints, file layouts, or code-level design unless the brainstorm itself is inherently technical.

## Interaction Rules

1. **Ask one question at a time** — Do not batch several unrelated questions into one message.
2. **Prefer single-select multiple choice** — Use single-select when choosing one direction, one priority, or one next step.
3. **Use multi-select rarely and intentionally** — Only for compatible sets such as goals, constraints, or success criteria that can all coexist. If prioritization matters, follow up by asking which selected item is primary.
4. **Use the platform's question tool when available** — `ask_user` in Copilot CLI, equivalent blocking tools elsewhere. Otherwise present numbered options in chat and wait.
5. **Hold all questions until research is summarized** — The user should see the research brief before the first product question. The only exceptions are clarifying which existing brainstorm to resume (Phase 0) and disambiguating scope (Phase 0.3).

## Output Guidance

- **Keep outputs concise** — Short sections, brief bullets, only enough detail to support the next decision.
- **Use repo-relative paths** — When referencing files, use paths relative to the repo root (e.g., `src/models/user.rb`), never absolute paths.
- **Mark evidence** — When the requirements doc cites research, link or quote the source. When a claim has no source, label it as an assumption.
- **Verify before claiming** — When the brainstorm touches checkable infrastructure (database tables, routes, config files, dependencies, model definitions), read the relevant source files to confirm what actually exists. Any claim that something is absent must be verified or labelled as an unverified assumption.

## Feature Description

<feature_description> #$ARGUMENTS </feature_description>

**If the feature description above is empty, ask the user:** "What would you like to explore? Please describe the feature, problem, or improvement you're thinking about. I'll do market and landscape research before asking product questions."

Do not proceed until you have a feature description from the user.

## Execution Flow

### Phase 0: Resume, Assess, and Route

#### 0.1 Resume Existing Work When Appropriate

If the user references an existing brainstorm topic or document, or there is an obvious recent matching `*-requirements.md` file in `docs/brainstorms/`:

- Read the document.
- Confirm with the user before resuming: "Found an existing requirements doc for [topic]. Should I continue from this, or start fresh?"
- If resuming, summarize the current state briefly, continue from existing decisions and outstanding questions, and update the existing document instead of creating a duplicate. Skip Phase 1 (intake). Decide in Phase 3 whether new external research is needed; the existing doc may already contain enough.

#### 0.2 Assess Whether Brainstorming Is Needed

**Clear-requirements indicators:**

- Specific acceptance criteria provided
- Referenced existing patterns to follow
- Described exact expected behavior
- Constrained, well-defined scope
- No framing risk — the user clearly knows the right shape and the landscape isn't going to change it

**If all of the above are true:** keep the interaction brief. Do a **minimal research pass** (Phase 3.2 prior art and Phase 3.3 applicable skills only — skip 3.1 landscape and 3.4 risks unless something surfaces) and a short research brief, then go to Phase 8 (capture) → Phase 9 (review) → Phase 10 (handoff). Skip Phases 5–7 (pressure test, full Q&A, approaches).

**Do not skip research entirely.** That is what `/ce-brainstorm` is for. The contract of this skill is research-first, even on small scopes.

#### 0.3 Assess Scope

Use the feature description plus a **very light pre-scan** (one or two ripgrep queries at most) to classify the work:

- **Lightweight** — small, well-bounded, low ambiguity
- **Standard** — normal feature or bounded refactor with some decisions to make
- **Deep** — cross-cutting, strategic, or highly ambiguous

Match research depth and Q&A depth to scope. Lightweight scopes get a single research pass with 1–2 questions; deep scopes get full landscape research and a longer dialogue.

If the scope is unclear, ask **one topic-identity question** to disambiguate (e.g., "is this about the API or the UI?"), then proceed. Do not ask scope, user, success-criteria, constraint, or prioritization questions yet — those come after research in Phase 6.

### Phase 1: Topic Intake

Restate the user's feature in your own words in 1–3 sentences and confirm:

- "Here's what I heard: [restated topic]. Did I get the core right?"

If the user corrects you, accept the correction silently and proceed.

**Strict rule for this phase:** Only **topic-identity confirmation** is allowed here. Do **not** ask scope, users, success criteria, constraints, prioritization, or trade-off questions. Those come in Phase 6 after research has run. The point of this phase is to make sure research targets the right thing, not to start product discovery.

### Phase 2: Repo Context Scan

Scan the repo before research. Match depth to scope.

**Lightweight** — Search for the topic, check if something similar already exists, and move on.

**Standard and Deep** — Two passes:

- *Constraint Check* — Check `AGENTS.md` and adjacent project instruction files for workflow, product, or scope constraints that affect the brainstorm. If these add nothing, move on.
- *Topic Scan* — Search for relevant terms. Read the most relevant existing artifact (brainstorm, plan, spec, skill, feature doc). Skim adjacent examples covering similar behavior.

If nothing obvious appears after a short scan, say so and continue.

Two rules govern technical depth during the scan:

1. **Verify before claiming** — When the topic touches checkable infrastructure, read the relevant source files. Claims of absence must be verified or labelled as unverified.
2. **Defer design decisions to planning** — Schemas, migration strategies, endpoint structure, and deployment topology belong in planning unless the brainstorm is itself about a technical decision.

### Phase 3: External Research

Run research in parallel where possible. Time-box it: prefer 3–7 targeted questions over exhaustive landscape scanning.

#### 3.1 Market and Landscape

For the problem space:

- How do other tools or products solve this?
- What is the current state of the art?
- Are there open-source solutions worth studying?
- What user-experience or implementation patterns are considered best practice?
- What scale or complexity thresholds make common approaches break down?

Aim for **3–5 concrete examples** when external research is available. If browsing or network access is unavailable, mark those facts as unverified rather than guessing.

#### 3.2 Prior Art and Learnings

Search for relevant institutional knowledge and similar code:

```bash
rg --files docs/solutions
rg -n "[key terms from the topic]"
```

If `rg` is unavailable, use the platform's native file search.

For each potentially relevant learning:

- Does it apply to this brainstorm?
- What specific insight should carry forward?
- If not applicable, why not?

#### 3.3 Applicable Skills

Check for skills that could provide domain-specific perspective:

```bash
rg --files .github/skills -g "SKILL.md"
```

Also check global / plugin skill roots exposed by the current platform (e.g., `~/.copilot/skills`, `~/.codex/skills`, `~/.claude/skills`).

For each matching skill, apply only the relevant perspective: does it suggest a framing, a constraint, or a known failure mode?

#### 3.4 Risk and Failure-Mode Survey

For the candidate approaches the topic implies, list known failure modes from prior art:

- What goes wrong at scale?
- What goes wrong on day 30 or 90, not day 1?
- What integrations or operational concerns are commonly underestimated?

### Phase 4: Synthesize Research Brief

Before any product question, produce a short **research brief** — the alignment artifact the user needs to answer questions well. This is conversational scaffolding, not part of the requirements doc.

**Distinction:**
- **Research Brief (Phase 4)** — everything notable found in research, used to align the user before Q&A. Lives in chat.
- **Research Summary (Phase 8 doc)** — only the findings that materially affected the requirements or decisions. Lives in the requirements doc.

Do not paste the brief verbatim into the doc later — distill.

```markdown
## Research Brief

**Landscape (3–5 examples):**
- [Tool / approach] — [what they do] — [why notable]

**Established patterns:**
- [Pattern] — [where it shows up] — [when it fits]

**Known failure modes:**
- [Failure] — [conditions]

**Repo prior art:**
- [Existing capability or pattern] (`path/to/file`) — [relevance]

**Applicable learnings:**
- [Learning title from docs/solutions or N/A]

**Open uncertainty:**
- [Question research could not resolve cleanly]
```

Display the brief to the user. Then ask one alignment question:

> "Does any of this change the framing before we go deeper? (a) Yes, here's what shifts; (b) No, my framing still holds; (c) Show me more on [topic]."

If the user selects (c), do another targeted research pass on that topic only and re-show the relevant section.

### Phase 5: Product Pressure Test

Now challenge the request — armed with research, this is sharper than vibes-only critique. Match depth to scope.

**Lightweight:**

- Does the research reveal a simpler off-the-shelf path?
- Is this duplicating something that already covers it?
- Is there a clearly better framing with near-zero extra cost?

**Standard:**

- Is this the right problem, or a proxy for a more important one?
- What user or business outcome actually matters here?
- What happens if we do nothing?
- Given the landscape we just surveyed, is there a nearby framing that compounds value at low extra cost?
- Is the highest-leverage move the request as framed, a reframing, an adjacent addition, a simplification, or doing nothing?

**Deep** — Standard questions plus:

- What durable capability should this create in 6–12 months?
- Does this move the product toward that, or is it only a local patch?
- Which of the failure modes from Phase 3.4 are we accepting?

Use the result to sharpen the conversation, not to bulldoze the user's intent.

### Phase 6: Targeted Q&A

Now run the conversation. The questions are sharper because they reference the research brief. Use the platform's blocking question tool when available.

**Guidelines:**

- Ask questions **one at a time**.
- Prefer **single-select** when choosing one direction, one priority, or one next step.
- Use **multi-select** only for compatible sets that can all coexist; if prioritization matters, ask which selected item is primary.
- Anchor each question to the research where appropriate: "Tools like X do A, others do B — which fits our users?"
- Start broad (problem, users, value) then narrow (constraints, exclusions, edge cases).
- Validate assumptions explicitly: "I'm assuming Y based on research finding Z — is that right?"
- Resolve product decisions here; leave technical implementation choices for planning.
- Make requirements concrete enough that planning will not need to invent behavior.

**Exit condition:** Continue until the idea is clear OR the user explicitly wants to proceed.

### Phase 7: Approaches

If multiple plausible directions remain, propose **2–3 concrete approaches** based on research and conversation. Otherwise state the recommended direction directly.

For each approach, provide:

- Brief description (2–3 sentences)
- Pros and cons
- Key risks or unknowns (from Phase 3.4)
- When it's best suited
- Closest analogue from research (e.g., "this is how X solves it")

When useful, include one deliberately higher-upside alternative — an adjacent reframing or addition that the landscape suggests would compound value, presented as a challenger option, not the default. Omit it when the work is already obviously over-scoped.

Lead with your recommendation and explain why. Prefer simpler solutions when added complexity creates real carrying cost, but do not reject low-cost, high-value polish.

If relevant, call out whether the choice is:

- Reuse an existing pattern
- Extend an existing capability
- Build something net new

### Phase 8: Capture the Requirements

Write or update a requirements document only when the conversation produced durable decisions worth preserving.

This document behaves like a lightweight PRD without PRD ceremony. Include what planning needs to execute well, and skip sections that add no value for the scope. Do **not** include implementation details such as libraries, schemas, endpoints, file layouts, or code structure unless the brainstorm is inherently technical.

**Required content for non-trivial work:**

- Problem frame
- Concrete requirements or intended behavior with stable IDs
- Scope boundaries
- Success criteria
- Research summary (top 3–5 findings with sources)

**Include when materially useful:**

- Key decisions and rationale (with research citations where applicable)
- Dependencies or assumptions
- Outstanding questions
- Alternatives considered (with research citations)
- Slice candidates — when handing off to `/kanban-plan`, list 3–7 candidate **user-visible increments** the research and conversation suggest. Keep these advisory and high-level — describe what each increment delivers, not blockers, ordering, or dependency design. `/kanban-plan` owns sequencing.

**Document structure:** Use this template and omit clearly inapplicable optional sections.

```markdown
---
date: YYYY-MM-DD
topic: <kebab-case-topic>
brainstorm_style: kanban-brainstorm
---

# <Topic Title>

## Problem Frame
[Who is affected, what is changing, and why it matters]

## Research Summary

**Findings that shaped requirements:**
- [Finding] — [which requirements/decisions it affected] — [link or note]

**Confidence:** High / Medium / Low — [one-line justification]

## Requirements

**[Group Header]**
- R1. [Concrete requirement]
- R2. [Concrete requirement]

**[Group Header]**
- R3. [Concrete requirement]

## Success Criteria
- [How we will know this solved the right problem]

## Scope Boundaries
- [Deliberate non-goal or exclusion]

## Key Decisions
- [Decision]: [Rationale] — Evidence: [research citation or "assumption"]

## Dependencies / Assumptions
- [Only include if material]

## Alternatives Considered
- [Approach]: [why not chosen] — [research citation]

## Slice Candidates (advisory for /kanban-plan)
- [Increment title] — [what user-visible behavior it delivers]
- [Increment title] — [what user-visible behavior it delivers]
<!-- Keep advisory. Do not assign blockers, ordering, or dependencies — that's /kanban-plan's job. -->

## Outstanding Questions

### Resolve Before Planning
- [Affects R1][User decision] [Question that must be answered before planning can proceed]

### Deferred to Planning
- [Affects R2][Technical] [Question that should be answered during planning]
- [Affects R2][Needs research] [Question that likely requires deeper research during planning]

## Next Steps
[If `Resolve Before Planning` is empty: `→ /kanban-plan` for vertical-slice decomposition (or `/ce-plan` for phased planning)]
[If `Resolve Before Planning` is not empty: `→ Resume /kanban-brainstorm` to resolve blocking questions before planning]
```

**Visual communication** — Include a visual aid when the requirements would be significantly easier to understand with one. Visual aids are conditional on content patterns, not depth classification.

**When to include:**

| Requirements describe... | Visual aid | Placement |
|---|---|---|
| A multi-step user workflow or process | Mermaid flow diagram or annotated ASCII flow | After Problem Frame, or under its own `## User Flow` heading |
| 3+ behavioral modes, variants, or states | Markdown comparison table | Within the Requirements section |
| 3+ interacting participants (user roles, components, services) | Mermaid or ASCII relationship diagram | After Problem Frame, or under `## Architecture` |
| Multiple competing approaches being compared | Comparison table | Within Phase 7 approach exploration |
| Comparison across landscape examples | Markdown comparison table | Within the Research Summary |

**When to skip:**

- Prose already communicates the concept clearly.
- The diagram would just restate the requirements in visual form.
- The visual describes implementation architecture, schemas, or code structure (that belongs in `/ce-plan` or `/kanban-plan`).
- The brainstorm is simple and linear with no multi-step flows or multi-participant interactions.

**Format selection:**

- **Mermaid** (default) for simple flows — 5–15 nodes, no in-box annotations. Use `TB` direction so diagrams stay narrow.
- **ASCII / box-drawing** for annotated flows that need rich in-box content. 80-column max for code blocks, vertical stacking.
- **Markdown tables** for mode/variant or approach comparisons.
- Place inline at the point of relevance.
- Conceptual level only — user flows, information flows, mode comparisons, component responsibilities.
- Prose is authoritative: when a visual aid and surrounding prose disagree, the prose governs.

After generating a visual aid, verify it accurately represents the prose requirements.

For **Standard** and **Deep** brainstorms, a requirements document is usually warranted.

For **Lightweight** brainstorms, keep the document compact. Skip document creation when only brief alignment is needed and no durable decisions need to be preserved.

For very small requirements docs with only 1–3 simple requirements, plain bullet requirements are acceptable. For **Standard** and **Deep** docs, use stable IDs like `R1`, `R2`, `R3` so planning and review can refer to them unambiguously.

When requirements span multiple distinct concerns, group them under bold topic headers within the Requirements section. Group by logical theme, not discussion order. Requirements keep their original IDs — numbering does not restart per group.

When the work is simple, combine sections rather than padding them. A short requirements document is better than a bloated one.

Before finalizing, check:

- What would `/kanban-plan` or `/ce-plan` still have to invent if this brainstorm ended now?
- Do any requirements depend on something claimed to be out of scope?
- Are any unresolved items actually product decisions rather than planning questions?
- Did implementation details leak in when they shouldn't have?
- Do any requirements claim that infrastructure is absent without verification?
- Is the research summary honest about confidence and gaps?
- Would a visual aid (flow diagram, comparison table, relationship diagram) help a reader grasp the requirements faster than prose alone?

If planning would need to invent product behavior, scope boundaries, or success criteria, the brainstorm is not complete yet.

Ensure `docs/brainstorms/` directory exists before writing.

If the document contains outstanding questions:

- Use `Resolve Before Planning` only for questions that truly block planning.
- If `Resolve Before Planning` is non-empty, keep working those questions during the brainstorm by default.
- If the user explicitly wants to proceed anyway, convert each remaining item into an explicit decision, assumption, or `Deferred to Planning` question first.
- Put technical or research-needing questions under `Deferred to Planning` when they are better answered there.
- Use tags like `[Needs research]` when the planner should likely investigate the question rather than answer from repo context alone.

### Phase 9: Document Review

When a requirements document was created or updated, run the `document-review` skill on it before presenting handoff options. Pass the document path as the argument.

If document-review returns findings that were auto-applied, note them briefly when presenting handoff options. If residual P0/P1 findings were surfaced, mention them so the user can decide whether to address them before proceeding.

When document-review returns "Review complete", proceed to Phase 10.

### Phase 10: Handoff

#### 10.1 Present Next-Step Options

Present next steps using the platform's blocking question tool when available. Otherwise present numbered options in chat and end the turn.

If `Resolve Before Planning` contains any items:

- Ask the blocking questions now, one at a time, by default.
- If the user explicitly wants to proceed anyway, first convert each remaining item into an explicit decision, assumption, or `Deferred to Planning` question.
- If the user chooses to pause instead, present the handoff as paused or blocked rather than complete.
- Do not offer "Proceed to planning" while `Resolve Before Planning` remains non-empty.

**Question when no blocking questions remain:** "Brainstorm complete. What would you like to do next?"

**Question when blocking questions remain and user wants to pause:** "Brainstorm paused. Planning is blocked until the remaining questions are resolved. What would you like to do next?"

Present only the options that apply:

- **Proceed to /kanban-plan (Recommended)** — Vertical-slice decomposition. Default for this skill because the requirements doc includes slice candidates.
- **Proceed to /ce-plan** — Phased implementation plan instead of slices. Use when work is sequential rather than independently slice-able.
- **Proceed directly to /ce-work** — Only offer when scope is lightweight, success criteria are clear, scope boundaries are clear, and no meaningful technical or research questions remain.
- **Run /deepen-brainstorm** — Run another targeted research pass on specific decisions or open questions.
- **Run additional document review** — Offer this only when a requirements document exists. Runs another pass for further refinement.
- **Ask more questions** — Continue clarifying scope, preferences, or edge cases.
- **Share to Proof** — Offer this only when a requirements document exists.
- **Done for now** — Return later.

If the direct-to-work gate is not satisfied, omit that option entirely.

#### 10.2 Handle the Selected Option

**If user selects "Proceed to /kanban-plan (Recommended)":**

Immediately run `/kanban-plan` in the current session, passing the requirements document path. Do not print the closing summary first.

**If user selects "Proceed to /ce-plan":**

Immediately run `/ce-plan` in the current session, passing the requirements document path. Do not print the closing summary first.

**If user selects "Proceed directly to /ce-work":**

Immediately run `/ce-work` in the current session using the finalized brainstorm output as context. If a compact requirements document exists, pass its path. Do not print the closing summary first.

**If user selects "Run /deepen-brainstorm":**

Load `deepen-brainstorm` and apply it to the requirements document for further targeted research. When it returns, present the Phase 10 options again with refreshed state.

**If user selects "Share to Proof":**

```bash
CONTENT=$(cat docs/brainstorms/YYYY-MM-DD-<topic>-requirements.md)
TITLE="Requirements: <topic title>"
RESPONSE=$(curl -s -X POST https://www.proofeditor.ai/share/markdown \
  -H "Content-Type: application/json" \
  -d "$(jq -n --arg title "$TITLE" --arg markdown "$CONTENT" --arg by "ai:compound" '{title: $title, markdown: $markdown, by: $by}')")
PROOF_URL=$(echo "$RESPONSE" | jq -r '.tokenUrl')
```

Display the URL prominently: `View & collaborate in Proof: <PROOF_URL>`

If the curl fails, skip silently. Then return to the Phase 10 options.

**If user selects "Ask more questions":** Return to Phase 6 (Targeted Q&A) and continue asking the user questions one at a time. Probe deeper into edge cases, constraints, preferences, or areas not yet explored.

When the user is satisfied with the additional Q&A, **do not jump straight back to Phase 10**. If the new conversation produced any change to requirements, scope, decisions, or success criteria, re-run Phase 8 (capture / update the requirements doc) → Phase 9 (document review) → Phase 10. Only short-circuit straight back to Phase 10 if the conversation purely confirmed existing decisions and added nothing new to the doc.

**If user selects "Run additional document review":**

Load the `document-review` skill and apply it to the requirements document for another pass. When document-review returns "Review complete", return to the normal Phase 10 options.

#### 10.3 Closing Summary

Use the closing summary only when this run of the workflow is ending or handing off, not when returning to the Phase 10 options.

When complete and ready for planning, display:

```text
Kanban brainstorm complete!

Requirements doc: docs/brainstorms/YYYY-MM-DD-<topic>-requirements.md  # if one was created

Top research findings:
- [Finding 1]
- [Finding 2]

Key decisions:
- [Decision 1]
- [Decision 2]

Slice candidates: [count]
Confidence: [High/Medium/Low]

Recommended next step: `/kanban-plan`
```

If the user pauses with `Resolve Before Planning` still populated, display:

```text
Kanban brainstorm paused.

Requirements doc: docs/brainstorms/YYYY-MM-DD-<topic>-requirements.md  # if one was created

Planning is blocked by:
- [Blocking question 1]
- [Blocking question 2]

Resume with `/kanban-brainstorm` when ready to resolve these before planning.
```

## Quality Checks

- [ ] Research happened **before** the first product question.
- [ ] The research brief was shown to the user before targeted Q&A.
- [ ] Every requirements claim about absent infrastructure was verified or labelled as an assumption.
- [ ] Decisions cite either a research source or are explicitly tagged as assumptions.
- [ ] The Slice Candidates section has 3–7 entries when handing off to `/kanban-plan`, or is omitted with reason.
- [ ] Confidence level in the research summary is honest about gaps.
- [ ] No implementation details leaked into the requirements doc (unless inherently technical).
- [ ] Document-review pass completed.

## Integration with Other Skills

- **Input from:** `/ce-ideate` (idea exploration), or a fresh feature description from the user.
- **Default handoff:** `/kanban-plan` for vertical-slice decomposition.
- **Alternate handoff:** `/ce-plan` for phased planning.
- **Optional follow-up:** `/deepen-brainstorm` for another targeted research pass on the produced doc.
- **Document review:** Always run `document-review` before handoff (Phase 9).
- **Peer skill:** `/ce-brainstorm` — same depth, but research happens after the conversation. Pick that one when conversation is the bottleneck and the design space is already well known.
