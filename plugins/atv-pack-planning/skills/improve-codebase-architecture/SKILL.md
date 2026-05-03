---
name: improve-codebase-architecture
description: "Find deepening opportunities — refactors that turn shallow modules into deep ones for testability and AI-navigability. Use when the user wants to improve architecture, consolidate coupled modules, make code more testable, or prepare a codebase for AI agents."
argument-hint: "[area of codebase to analyze, or blank for full scan]"
---

# Improve Codebase Architecture

<!-- Adapted from mattpocock/skills — credit: github.com/mattpocock/skills -->

Surface architectural friction and propose **deepening opportunities** — refactors that turn shallow modules into deep ones. The aim is testability and AI-navigability.

## Key Concepts

- **Deep module** — a lot of behavior behind a small interface. High leverage for callers, high locality for maintainers.
- **Shallow module** — interface nearly as complex as the implementation. Low leverage. Often a pass-through.
- **Deletion test** — imagine deleting the module. If complexity vanishes, it was a pass-through. If complexity reappears across N callers, it was earning its keep.
- **Seam** — where an interface lives; a place behavior can be altered without editing in place.

## Why This Matters for AI Agents

A codebase of many tiny, shallow modules is hard for both humans and agents:

- Dependency graph is hard to understand
- Test boundaries are unclear
- Agents modify the wrong layer
- Small changes cause unexpected breakage
- Agent must inspect too many files to understand one behavior

Deep modules fix this: small public interface → agent owns the internals, human owns the boundary.

## Process

### 1. Explore

Use grep/glob/view to walk the codebase organically. Note friction:

- Where does understanding one concept require bouncing between many small modules?
- Where are modules shallow — interface nearly as complex as the implementation?
- Where have pure functions been extracted just for testability, but real bugs hide in how they're called?
- Where do tightly-coupled modules leak across their seams?
- Which parts are untested or hard to test through their current interface?

Apply the **deletion test** to anything you suspect is shallow.

### 2. Present Candidates

Present a numbered list of deepening opportunities. For each:

- **Files** — which files/modules are involved
- **Problem** — why the current architecture causes friction (in terms of depth/locality)
- **Proposed deepening** — what the refactored module would look like (small interface, deep implementation)
- **Benefits** — how tests would improve, how AI agents would navigate it better

Ask: "Which of these would you like to explore?"

### 3. Design the Deepened Module

For the chosen candidate:

- [ ] Define the minimal public interface
- [ ] Identify what moves INSIDE (currently spread across callers)
- [ ] Define the test boundary (tests use the public interface only)
- [ ] Identify what adapters sit at seams
- [ ] Confirm with user before implementing

### 4. Execute (if approved)

Invoke the `tdd` skill to implement the deepened module test-first:

1. Write interface tests (behavior, not implementation)
2. Move/consolidate internal code behind the interface
3. Delete shallow pass-throughs that the deletion test identified
4. Verify all existing tests still pass
5. Remove dead code

## Anti-Patterns to Flag

| Pattern | Problem | Fix |
|---------|---------|-----|
| Many tiny utils | No locality, bugs hide in composition | Consolidate into domain module |
| Service per function | Interface cost exceeds implementation | Merge related services |
| Mock-heavy tests | Testing implementation, not behavior | Test through public interface |
| Config-driven everything | Interface is the config schema | Hardcode sane defaults, expose only what varies |
| God object | Deep but interface is too wide | Extract focused deep modules from it |
