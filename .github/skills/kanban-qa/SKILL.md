---
name: kanban-qa
description: "Browser-based verification for frontend slices. Connects via CDP, Playwright, or Agent Browser based on context. Hard gate — the browser reports what rendered, the model does not self-report. Use when a kanban slice touches frontend files (.tsx, .jsx, .html, .css, .vue, .svelte) and needs visual/interactive verification."
argument-hint: "[slice plan path, or blank to verify the current slice]"
---

# Kanban QA — Browser Verification Gate

Verify frontend slices by looking at what the app actually renders. The agent tests as a user, not a developer. Source code is irrelevant during QA — only what the browser shows matters.

## When to Run

Called from `kanban-work` after Step 3.6 (Diff-Scope Verification) passes, **only** for slices whose `expected_files` include frontend file extensions: `.tsx`, `.jsx`, `.html`, `.css`, `.scss`, `.vue`, `.svelte`, `.ejs`, `.hbs`.

If the slice is backend-only, skip entirely with a one-line note: `qa: skipped — no frontend files in expected_files`

## Input

<input> #$ARGUMENTS </input>

**If input is empty:** Read the current slice context from the active kanban manifest.

**If input is a path:** Read that slice plan directly.

## Step 0: Transport Selection

Pick the transport based on what the slice needs, not a fixed priority list.

### Decision Logic

1. **Is this an internal/corporate site?** (SSO, Conditional Access, company-owned domains, session cookies from a real login)

   - **YES** → CDP required. Connect to the user's existing browser session via `ws://localhost:9222` (or `CDP_ENDPOINT` env var). The real browser already has cookies, tokens, and session state from the user's login. No way to fake this.
   - If CDP is unavailable → **STOP.** Do not attempt Playwright or Agent Browser on internal sites — they cannot pass SSO/Conditional Access. Log in kanban.md: `qa: skipped — internal site, no CDP session. Start browser with --remote-debugging-port=9222`

2. **Is this a regular site or local dev server?** (localhost, public URLs, no corporate auth)

   - Agent Browser if installed (`agent-browser` on PATH) — structured element targeting, fast
   - Playwright if available — headless, clean viewport control, good for responsive testing
   - CDP as fallback

3. **Does the slice need responsive/viewport testing?** (deep tier, or slice touches layout/grid/responsive components)

   - Playwright preferred for multi-viewport. Headless, spawns 375px/768px/1440px cleanly.
   - Fallback: CDP with device emulation.

4. **None available** → Log: `qa: skipped — no browser transport available (checked: CDP, Playwright, Agent Browser)`. Not fatal, but flag it in kanban.md.

## Step 1: Connect and Navigate

1. Connect via the selected transport.
2. Navigate to the page(s) relevant to the slice. Determine URLs from:
   - Slice plan's acceptance criteria
   - Route patterns in `expected_files` (e.g., `pages/dashboard.tsx` → `/dashboard`)
   - Dev server URL (default `http://localhost:3000`, respect `DEV_SERVER_URL` env var)
3. Wait for the page to reach a stable state (network idle or DOM content loaded).

## Step 2: Capture Evidence

1. **Screenshot** the rendered page. Save to `.atv/qa-screenshots/<slice-id>-<timestamp>.png`.
2. **Capture console output** — errors, warnings, and failed network requests (4xx/5xx).
3. **Check for render failures** — blank pages, loading spinners that never resolve, error boundaries.

Create `.atv/qa-screenshots/` if it doesn't exist.

## Step 3: Verify Against Slice Requirements

Read the slice plan's acceptance criteria. For each criterion that describes visible behavior:

| Check | Method |
|-------|--------|
| Element exists | Query the DOM for the expected element |
| Text content matches | Read rendered text, compare to expected |
| Layout is correct | Screenshot comparison against description |
| No regressions | Console has no new errors or warnings |

**Never read source code during QA.** Test as a user — if you can't verify it from the browser, flag it as needing manual verification.

## Step 4: Interaction Checks (standard + deep tiers)

If the tier is `standard` or `deep` and the slice added interactive elements:

1. Click buttons, links, and interactive elements added by the slice.
2. Fill form fields if the slice added forms.
3. After each interaction, check for:
   - New console errors
   - Failed network requests
   - Unexpected navigation
   - UI elements disappearing or breaking
4. Screenshot after each significant interaction.

## Step 5: Responsive Checks (deep tier only)

If the tier is `deep`:

1. Resize viewport to **375px** (mobile) — screenshot + console check.
2. Resize viewport to **768px** (tablet) — screenshot + console check.
3. Resize viewport to **1440px** (desktop) — screenshot + console check.
4. Flag any layout breakage, overflow, or elements hidden at specific breakpoints.

## Step 6: Report

**On pass:**

```text
qa: PASS — <transport used>
  checks: N/N passed
  console: clean (0 errors, 0 warnings)
  screenshots: .atv/qa-screenshots/<slice-id>-*.png
  tier: quick|standard|deep
```

**On fail:**

For each failure:
1. Verify it's reproducible — retry once before reporting.
2. Screenshot the failure state (mandatory).
3. Log the specific check that failed and what was expected vs. observed.

```text
qa: FAIL — <transport used>
  failed: "<check description>"
  expected: "<what the slice said>"
  observed: "<what the browser showed>"
  screenshot: .atv/qa-screenshots/<slice-id>-fail-<n>.png
  console_errors: [list if any]
```

**FAIL on any critical check is a hard gate.** The agent MUST NOT proceed to the next slice. Same enforcement as Step 3.6 (Diff-Scope Verification).

Log all results in `docs/kanban.md` under the slice's status or notes. Also update the manifest `notes` field.

## Tiers

Set per-feature or per-slice in kanban.md or the manifest. Default is `quick`.

| Tier | What it does |
|------|-------------|
| `quick` | Screenshot + console check only |
| `standard` | + interaction checks on new/modified elements |
| `deep` | + responsive checks at 3 breakpoints (375px, 768px, 1440px) |

## Transport Reference

Actions are transport-agnostic above. Here's the mapping:

| Action | CDP | Agent Browser | Playwright |
|--------|-----|---------------|------------|
| Connect | `ws://localhost:9222/json` → get `webSocketDebuggerUrl` | `agent-browser open <url>` | `playwright.chromium.launch()` |
| Navigate | `Page.navigate` | `agent-browser open <url>` | `page.goto(url)` |
| Screenshot | `Page.captureScreenshot` | `agent-browser screenshot <file>` | `page.screenshot()` |
| Console errors | `Runtime.consoleAPICalled` event listener | `agent-browser snapshot -i` (inspect output) | `page.on('console')` |
| Click element | `DOM.querySelector` + `Input.dispatchMouseEvent` | `agent-browser click @<ref>` | `page.click(selector)` |
| Get text | `DOM.querySelector` + `DOM.getOuterHTML` | `agent-browser snapshot -i` + parse | `page.textContent(selector)` |
| Resize viewport | `Emulation.setDeviceMetricsOverride` | Not supported — use Playwright | `page.setViewportSize()` |
| Network errors | `Network.responseReceived` event listener | Not directly — check console | `page.on('response')` |

## Principles

- Never read source code during QA. Test as a user, not a developer.
- Verify before documenting — retry once to confirm reproducible.
- Screenshot evidence is mandatory for any failure.
- Write incrementally, don't batch findings.
- Auth is sacred — never attempt authenticated routes without a real session.
- If a check can't be verified from the browser, flag it for manual review rather than guessing.

## Integration

- **Called from:** `kanban-work` (after Step 3.6, frontend slices only)
- **Results feed into:** `ce-review` (Step 5.4) as additional context
- **Screenshots persist:** `.atv/qa-screenshots/` (gitignored, ephemeral)
- **Logs persist:** kanban.md notes + manifest notes
