---
name: cso
description: "OWASP Top 10 + STRIDE threat model for application source code. Scans for SQL injection, XSS, CSRF, auth bypass, and other web application vulnerabilities. Triggers on 'owasp scan', 'threat model', 'cso', 'security review code', 'stride analysis', 'application security'."
argument-hint: "[scope: full (default) | owasp | stride | path/to/scan]"
---

# /cso — Chief Security Officer

Scan application source code for OWASP Top 10 vulnerabilities and produce a STRIDE threat model. Complementary to `/atv-security` which covers agentic config files.

## Arguments

<scope> #$ARGUMENTS </scope>

**Scope detection:**
- Contains "owasp" → OWASP scan only
- Contains "stride" → STRIDE threat model only
- Contains a file/directory path → scope scan to that path
- Otherwise → full scan (OWASP + STRIDE)

## Execution Flow

```
Phase 1: Discover    → Detect stack, find source files
Phase 2: OWASP Scan  → Check Top 10 vulnerability categories
Phase 3: STRIDE      → Threat model the application architecture
Phase 4: Report      → Graded report with findings + threat matrix
```

---

## Phase 1: Discover Application Stack

Use `file_search` and `list_dir` to detect the application stack:

| Signal | Stack | Key files to scan |
|--------|-------|-------------------|
| `package.json`, `*.ts`, `*.js` | Node.js / TypeScript | `src/**`, `routes/**`, `api/**`, `pages/**` |
| `requirements.txt`, `*.py` | Python | `app/**`, `src/**`, `views/**`, `api/**` |
| `Gemfile`, `*.rb` | Ruby / Rails | `app/**`, `config/**`, `db/**` |
| `go.mod`, `*.go` | Go | `**/*.go` |
| `*.cs`, `*.csproj` | .NET | `**/*.cs`, `Controllers/**` |
| `pom.xml`, `*.java` | Java | `src/**/*.java` |

Record: stack detected, total source files found, entry points identified.

If no application source code is found, report: "No application source code detected. `/cso` scans app code — for agentic config security, use `/atv-security`." and stop.

---

## Phase 2: OWASP Top 10 (2021) Scan

For each category, use `grep_search` for Tier 1 patterns and `read_file` + LLM assessment for Tier 2. Scan only application source files, not `.github/` configs.

### A01: Broken Access Control

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `role\s*===?\s*["']admin["']` scoped to route/controller files | Hardcoded role checks instead of RBAC | 🟢 medium |
| `req\.user\s*&&` without authorization middleware | Ad-hoc auth checks bypassing middleware | 🟡 high |

**Tier 2 — LLM assessment:**
- Read route/controller files. Check: Are all state-changing endpoints protected by auth middleware?
- Look for: direct object references without ownership validation (e.g., `/users/:id` without checking `req.user.id === id`)
- Look for: missing authorization on admin/management endpoints
- Severity: 🟡 high per unprotected endpoint

### A02: Cryptographic Failures

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `md5\|sha1\|DES\|RC4` in crypto/hash contexts | Weak/deprecated algorithms | 🟡 high |
| `http://` in API endpoint URLs (not localhost) | Unencrypted data in transit | 🟢 medium |
| `password.*=.*["'][^$]` | Hardcoded passwords | 🔴 critical |

**Tier 2 — LLM assessment:**
- Check: Are passwords hashed with bcrypt/scrypt/argon2, not MD5/SHA1?
- Check: Is sensitive data (PII, tokens, cards) encrypted at rest?
- Check: Are TLS/HTTPS enforced for external communications?

### A03: Injection

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `query\s*\(\s*["'\x60].*\$\{` or `query\s*\(.*\+\s*` | SQL injection via string concatenation | 🔴 critical |
| `eval\s*\(` or `exec\s*\(` or `Function\s*\(` with variable input | Code injection | 🔴 critical |
| `innerHTML\s*=` or `dangerouslySetInnerHTML` or `\| safe` or `\|raw` | XSS via unsafe HTML rendering | 🟡 high |
| `child_process\.\(exec\|spawn\).*\$\{` or `subprocess\.call.*\+` | OS command injection | 🔴 critical |
| `\.find\(\{.*\$` or `\.aggregate\(\[.*\$` in Mongo contexts | NoSQL injection | 🟡 high |

**Tier 2 — LLM assessment:**
- Read files with database queries. Are all queries parameterized?
- Read template files. Is user input escaped before rendering?
- Check for LDAP injection, XPath injection, header injection patterns

### A04: Insecure Design

**Tier 2 — LLM assessment only (no grep patterns):**
- Check: Is there rate limiting on auth endpoints (login, register, password reset)?
- Check: Are there business logic flaws (e.g., negative quantity in cart, price manipulation)?
- Check: Is there account enumeration via different error messages for valid/invalid usernames?
- Severity: 🟡 high per design flaw

### A05: Security Misconfiguration

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `DEBUG\s*=\s*True\|debug:\s*true` | Debug mode enabled | 🟡 high |
| `cors\(\)` without origin restriction, or `origin:\s*["']\*["']` | Unrestricted CORS | 🟡 high |
| `helmet` not imported/used (Node.js) | Missing security headers | 🟢 medium |
| `ALLOWED_HOSTS\s*=\s*\[["']\*["']\]` | Django wildcard hosts | 🟡 high |

**Tier 2 — LLM assessment:**
- Check: Are default credentials changed?
- Check: Are error pages custom (not showing stack traces)?
- Check: Are unnecessary features/endpoints disabled in production config?

### A06: Vulnerable and Outdated Components

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `"dependencies"` in package.json | Check for known-vulnerable versions | 🟢 medium |

**Tier 2 — LLM assessment:**
- Read `package.json`, `requirements.txt`, `Gemfile`, or `go.mod`
- Flag any dependency that hasn't been updated in >1 year (check version patterns)
- Recommend: `npm audit`, `pip-audit`, `bundle audit`, `govulncheck`
- Severity: 🟢 medium (recommend tooling, don't duplicate it)

### A07: Identification and Authentication Failures

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `jwt\.sign.*expiresIn.*["']30d\|["']365d\|["']never` | Excessive token lifetime | 🟡 high |
| `session.*maxAge.*86400000` (>24h in ms) | Long session duration | 🟢 medium |
| `bcrypt.*rounds.*[1-5][^0-9]` or `salt.*rounds.*[1-5][^0-9]` | Weak bcrypt rounds (<6) | 🟡 high |

**Tier 2 — LLM assessment:**
- Check: Is there brute-force protection (account lockout, progressive delays)?
- Check: Is password complexity enforced?
- Check: Are sessions invalidated on logout/password change?

### A08: Software and Data Integrity Failures

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `deserialize\|unserialize\|pickle\.load\|yaml\.load\b` | Unsafe deserialization | 🔴 critical |
| `integrity=` absent in `<script src="https://` | Missing SRI on CDN scripts | 🟢 medium |

**Tier 2 — LLM assessment:**
- Check: Is CI/CD pipeline protected against tampering?
- Check: Are software updates verified with signatures?

### A09: Security Logging and Monitoring Failures

**Tier 2 — LLM assessment only:**
- Check: Are login failures, access denied events, and input validation failures logged?
- Check: Are logs protected against injection (structured logging vs string concat)?
- Check: Is there alerting on suspicious patterns?
- Severity: 🟢 medium per gap

### A10: Server-Side Request Forgery (SSRF)

**Tier 1 — grep patterns:**

| Pattern | What it catches | Severity |
|---------|----------------|----------|
| `fetch\s*\(\s*\w+\|axios\.\w+\(\s*\w+\|requests\.\w+\(\s*\w+` where the URL is a variable | Potential SSRF if URL is user-controlled | 🟡 high |
| `http\.Get\(\s*\w+\|urllib\.request\.urlopen\(\s*\w+` | Same pattern in Go/Python | 🟡 high |

**Tier 2 — LLM assessment:**
- Check: Is user input validated/allowlisted before being used in server-side HTTP requests?
- Check: Are internal network ranges (169.254.x.x, 10.x.x.x, 127.x.x.x) blocked?

---

## Phase 3: STRIDE Threat Model

Read the application's architecture by examining:
- Entry points (routes, API endpoints, webhooks, event handlers)
- Data flows (database queries, external API calls, file I/O)
- Trust boundaries (auth middleware, API gateways, service boundaries)
- Assets (user data, credentials, tokens, business logic)

Produce a threat matrix:

| Threat | Category | Description | Affected Component | Risk | Mitigation |
|--------|----------|-------------|-------------------|------|------------|
| **S**poofing | Identity | Can an attacker impersonate a user/service? | [auth system, API] | [H/M/L] | [existing or missing control] |
| **T**ampering | Data integrity | Can data be modified in transit or at rest? | [database, API payload] | [H/M/L] | [existing or missing control] |
| **R**epudiation | Accountability | Can actions be performed without audit trail? | [logging system] | [H/M/L] | [existing or missing control] |
| **I**nformation Disclosure | Confidentiality | Can sensitive data leak? | [error pages, logs, API responses] | [H/M/L] | [existing or missing control] |
| **D**enial of Service | Availability | Can the service be overwhelmed? | [endpoints without rate limiting] | [H/M/L] | [existing or missing control] |
| **E**levation of Privilege | Authorization | Can a user gain unauthorized access? | [role system, admin endpoints] | [H/M/L] | [existing or missing control] |

For each threat:
- Identify whether a mitigation **already exists** in the codebase
- If missing, provide a concrete recommendation
- Rate risk as High/Medium/Low based on exploitability and impact

---

## Phase 4: Report

### Scoring

**OWASP score** — start at 100, deduct per finding:
- 🔴 critical: −15
- 🟡 high: −10
- 🟢 medium: −5
- 🔵 low: −2

Floor at 0. Grade: A≥90, B≥80, C≥65, D≥50, F<50.

**STRIDE score** — count of unmitigated threats:
- 0 unmitigated: 🟢 Strong posture
- 1–2 unmitigated: 🟡 Moderate risk
- 3+ unmitigated: 🔴 Weak posture

### Report Format

```markdown
## 🔒 CSO Security Report — OWASP Top 10 + STRIDE

**Date:** YYYY-MM-DD
**Stack:** [detected stack]

| Metric | Value |
|--------|-------|
| **OWASP Grade** | [A–F] |
| **OWASP Score** | [0–100]/100 |
| **STRIDE Posture** | [🟢/🟡/🔴] |

### OWASP Findings

#### 🔴 Critical
- **[A03] SQL Injection** in `src/db/queries.ts:42`
  Evidence: `db.query("SELECT * FROM users WHERE id = " + userId)`
  Fix: Use parameterized query: `db.query("SELECT * FROM users WHERE id = $1", [userId])`

#### 🟡 High
- ...

#### 🟢 Medium
- ...

### STRIDE Threat Matrix

| Threat | Risk | Status |
|--------|------|--------|
| Spoofing | [H/M/L] | [✅ Mitigated / ⚠️ Partial / ❌ Unmitigated] |
| Tampering | [H/M/L] | [✅/⚠️/❌] |
| Repudiation | [H/M/L] | [✅/⚠️/❌] |
| Info Disclosure | [H/M/L] | [✅/⚠️/❌] |
| Denial of Service | [H/M/L] | [✅/⚠️/❌] |
| Elevation of Privilege | [H/M/L] | [✅/⚠️/❌] |

### Summary

| Source files scanned | OWASP findings | STRIDE unmitigated | Recommended next |
|---------------------|----------------|-------------------|------------------|
| [N] | [N] | [N] | [action] |

For agentic config security: run `/atv-security`
```

If zero OWASP findings and zero unmitigated STRIDE threats: "Your application code looks secure! No findings detected."

### Phase 5: Persist Report (always, after report is generated)

After printing the report in chat, persist it to disk so it survives the conversation.

**Target file:** `docs/security/YYYY-MM-DD-security-report.md` (today's date, UTC). One shared file per day with separate sections for `/atv-security` and `/cso`.

**Steps:**

1. Ensure `docs/security/` exists. If not, create it (write the file — the directory is created implicitly).
2. Compute today's date as `YYYY-MM-DD` and the run timestamp as ISO-8601 (e.g., `2026-04-24T14:32:10Z`).
3. Try to `read_file` the target path.
   - **If the file does not exist:** `create_file` with this skeleton, then continue at step 4:
     ```markdown
     # Security Report — YYYY-MM-DD

     <!-- atv-security:start -->
     ## /atv-security Scan
     _No scan recorded yet._
     <!-- atv-security:end -->

     <!-- cso:start -->
     ## /cso Scan
     _No scan recorded yet._
     <!-- cso:end -->
     ```
   - **If the file exists:** continue at step 4.
4. Build the new section content (everything between the markers, exclusive):
   ```markdown
   ## /cso Scan — <ISO timestamp>

   - **Scope:** full | owasp | stride | <path>
   - **Stack:** <detected stack>
   - **OWASP Grade:** <A–F> · **Score:** <0–100>/100
   - **STRIDE Posture:** 🟢/🟡/🔴

   <full report markdown from Phase 4, including OWASP findings and STRIDE matrix>
   ```
5. Use `replace_string_in_file` to swap the existing `<!-- cso:start -->` … `<!-- cso:end -->` block with the markers wrapping the new content. The `<!-- atv-security:* -->` block must be left untouched.
6. Confirm in chat: `📄 Report saved to docs/security/YYYY-MM-DD-security-report.md (cso section).`

**Constraints:**
- Never delete or modify the `<!-- atv-security:* -->` section.
- Always keep both marker pairs intact so `/atv-security` can upsert into the same file later.
- If `replace_string_in_file` cannot find the marker block (file was hand-edited), fall back to overwriting the file with a fresh skeleton containing the new `/cso` section and an empty `/atv-security` placeholder, and warn the user that prior `/atv-security` content may have been preserved separately.

---

## What This Skill Does NOT Do

- Scan agentic config files (`.github/`, `.vscode/`) — use `/atv-security` for that
- Run dynamic application security testing (DAST/penetration testing)
- Scan container images or infrastructure-as-code
- Replace dedicated SAST tools (Semgrep, CodeQL, Snyk) — this is a fast triage layer
- Modify source code without explicit user confirmation
