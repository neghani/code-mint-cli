# CodeMint CLI — Feature roadmap and missed ideas

This doc lists useful features that are not yet implemented, plus ideas to improve adoption and UX. Priorities: **P0** = critical, **P1** = high value, **P2** = nice to have.

## Already tracked in code-mint-ai (backend) docs

- **CLI_GAPS.md** — Windows token read, org list shape (fixed), Cursor skills path (fixed), Windsurf/Continue/Claude skills dirs, checksum on install, manifest `lastSyncAt`, file locking, `--debug`, Linux in releases (fixed), `--force` on add, `--limit` on suggest, doctor API check, token expiry, remove empty skill dir, etc.
- **TODO_CLI_INTEGRATION.md** — Backend and CLI contract checklist (many items done).

## High-value features not yet in gaps doc

### P1 — Adoption and DX

| Feature | Description |
|--------|-------------|
| **Shell completion** | `codemint completion bash|zsh|fish` and install instructions. Speeds up discovery and reduces typos. |
| **Self-upgrade** | `codemint upgrade` or `codemint self-update`: fetch latest release for current OS/arch, replace binary, print version. Optional `--check` to only report if upgrade available. |
| **Batch add** | `codemint add @rule/a @rule/b @skill/c` — accept multiple refs in one command. |
| **Pin version on add** | `codemint add @rule/slug@1.2.0` — install a specific catalog version; manifest stores it for reproducible installs. |
| **Config init** | `codemint init` or `codemint config init`: create `.codemint/config.json` or ensure `.codemint/` and manifest with defaults (e.g. tool from auto-detect). |
| **Token list/revoke from CLI** | `codemint auth tokens` (list), `codemint auth revoke <id>` — call `GET /api/auth/cli-token` and `DELETE /api/auth/cli-token/:id` so users can manage tokens without the browser. |

### P1 — Reliability and polish

| Feature | Description |
|--------|-------------|
| **Checksum on install** | After writing a rule/skill file, verify SHA256 against `item.Checksum`; on mismatch, rollback and error (see CLI_GAPS #8). |
| **Manifest lastSyncAt** | Set and persist `lastSyncAt` after successful sync (CLI_GAPS #9). |
| **File locking** | Use `.codemint/.lock` (flock) around manifest load/save and install/remove to avoid races with extension or multiple terminals (CLI_GAPS #10). |
| **Remove empty skill dir** | After removing `.cursor/skills/<slug>/SKILL.md`, remove the empty `<slug>` directory (CLI_GAPS #24). |
| **Windsurf/Continue/Claude skills dirs** | Use `.windsurf/skills`, `.continue/skills`, `.claude/skills` for skills instead of putting them in rules dir (CLI_GAPS #7). |

### P2 — Discovery and telemetry

| Feature | Description |
|--------|-------------|
| **Suggest reason** | When tag overlap is low, show a clearer reason string (e.g. “No tag match; try browsing by tech”) and optionally surface matched tags (CLI_GAPS #21). |
| **Doctor API check** | `codemint doctor` calls `GET /api/auth/me` and reports “API: OK” or “API: unreachable” (CLI_GAPS #22). |
| **Opt-in telemetry** | Document decision (default off); if added later, limit to anonymous usage counts and errors. |
| **Export/backup** | `codemint export` or `codemint backup` — tar/zip of `.codemint/` and tool-specific installed paths for the repo, for restore or sharing. |

### P2 — Distribution

| Feature | Description |
|--------|-------------|
| **Homebrew** | Formula in a tap; `brew install neghani/code-mint-cli/codemint` (execution plan Phase F). |
| **winget** | Manifest and publish (execution plan Phase F). |
| **Signing** | macOS notarization, Windows Authenticode (execution plan Phase E). |

## Backend (code-mint-ai) items that unblock CLI

- Login redirect: preserve `next` so after sign-in user returns to `/cli-auth?port=...`.
- Token lifecycle: keep `GET/POST/DELETE /api/auth/cli-token` stable; optional `expiresAt` and scopes for tokens.
- Catalog: resolve and sync endpoints are in place; ensure pagination and error shapes stay consistent.

## Reference

- Execution plan: **CODEMINT_CLI_EXECUTION_PLAN.md**
- Gaps (and many of these items): app repo **docs/CLI_GAPS.md**
- Coordination: app repo **docs/REPO_COORDINATION.md**
