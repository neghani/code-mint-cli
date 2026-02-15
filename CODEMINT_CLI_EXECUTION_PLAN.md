# codemint CLI: Execution Checklist and Delivery Plan

## Scope
Build and ship a production-ready `codemint` CLI as standalone binaries for:
- macOS: `arm64`, `amd64`
- Windows: `amd64`

Distribution order:
1. GitHub Releases
2. Homebrew
3. winget

## Target Repository Layout

This CLI lives in its **own repo** ([neghani/code-mint-cli](https://github.com/neghani/code-mint-cli)), not under the web app. The app repo is **code-mint-ai** (backend, API, `/cli-auth`, install script served at `/cli/install.sh`). Coordination: see the **code-mint-ai** app repo’s `docs/REPO_COORDINATION.md`.

Structure in this repo:

```text
(code-mint-ali-cli)/
  cmd/
    root.go
    version.go
    auth.go
    auth_login.go
    auth_whoami.go
    auth_logout.go
    items.go
    items_search.go
    org.go
    org_list.go
  internal/
    api/
      client.go
      errors.go
      middleware.go
      models.go
    auth/
      flow.go
      callback_server.go
      browser.go
      storage.go
      storage_darwin.go
      storage_windows.go
      redact.go
    config/
      config.go
      env.go
    output/
      table.go
      json.go
      format.go
  test/
    integration/
      auth_flow_test.go
      whoami_test.go
      items_search_test.go
      org_list_test.go
  .github/
    workflows/
      release.yml
      ci.yml
  docs/
    quickstart.md
    commands.md
    troubleshooting.md
    enterprise.md
  scripts/
    checksums.sh
    release_notes.sh
  go.mod
  go.sum
  main.go
  Makefile
  README.md
```

## Phase-by-Phase Execution Checklist

## Phase A: Foundation (Week 1)

### Backend Alignment (Blocker Check)
- [ ] Confirm `/cli-auth` login redirect preserves `next=...` through successful login.
- [ ] Verify `/api/auth/me`, `/api/items/search`, `/api/org/my` contracts are stable.
- [ ] Freeze auth response schema for CLI (`token`, `expiresAt` if applicable).

### CLI Bootstrap
- [ ] Initialize Go module and pin Go version (`1.22+`).
- [ ] Add Cobra root command and `version` command.
- [ ] Add base config layer:
  - [ ] `base URL`
  - [ ] `profile`
  - [ ] env overrides (`CODEMINT_BASE_URL`, `CODEMINT_PROFILE`)
- [ ] Add shared HTTP client with:
  - [ ] timeout defaults
  - [ ] `User-Agent: codemint/<version> (<os>/<arch>)`
- [ ] Add standard error envelope parsing and exit codes.
- [ ] Add output mode plumbing (`table` default, `--json`).

### Exit Criteria (M1)
- [ ] `codemint help` stable and complete.
- [ ] `codemint version` returns build metadata.
- [ ] CI runs lint + unit tests on pushes/PRs.

## Phase B: Auth Flow (Week 1-2)

### Auth UX and Flow
- [ ] Implement `codemint auth login`:
  - [ ] start local callback server on ephemeral port
  - [ ] open browser to `/cli-auth?port=<port>`
  - [ ] print URL fallback when browser open fails
- [ ] Parse callback payload and capture token safely.
- [ ] Immediately verify token via `GET /api/auth/me`.
- [ ] Persist token in OS secure storage:
  - [ ] macOS Keychain
  - [ ] Windows Credential Manager
- [ ] Implement `codemint auth whoami` using stored token.
- [ ] Implement `codemint auth logout` to delete local secret.

### Security Controls
- [ ] Add token redaction helper in all errors/logging paths.
- [ ] Ensure no token value appears in debug logs or panic output.
- [ ] Return actionable auth failures (expired/revoked/invalid).

### Exit Criteria (M2 partial)
- [ ] `login`, `whoami`, `logout` pass on macOS + Windows.
- [ ] Invalid token path prompts re-login clearly.

## Phase C: Domain Commands (Week 2)

### items/search
- [ ] Implement `codemint items search`:
  - [ ] `--q`
  - [ ] `--type`
  - [ ] `--tags`
  - [ ] pagination flags (`--page`, `--limit`)
- [ ] Map API response to table + JSON output.

### org/list
- [ ] Implement `codemint org list` against `/api/org/my`.
- [ ] Render table and `--json` output consistently.

### Exit Criteria (M2)
- [ ] Search and org commands work against dev and deployed env.
- [ ] Output schema stable for automation consumers.

## Phase D: Hardening (Week 3)

### Reliability
- [ ] Add retry policy for transient network errors.
- [ ] Configure per-request timeout + cancellation handling.
- [ ] Improve API error mapping (401, 403, 422, 429, 5xx).
- [ ] Add proxy/certificate environment support docs and tests.

### Security and Policy
- [ ] Telemetry decision documented (default off).
- [ ] Sensitive header redaction in HTTP debug mode.

### Exit Criteria (M3 partial)
- [ ] Token revocation and network interruption scenarios tested.
- [ ] CLI behavior deterministic under retries/timeouts.

## Phase E: Packaging + Release (Week 3-4)

### Build and Artifacts
- [ ] Add cross-build matrix:
  - [ ] `darwin/arm64`
  - [ ] `darwin/amd64`
  - [ ] `windows/amd64`
- [ ] Produce archives:
  - [ ] `codemint_<version>_darwin_arm64.tar.gz`
  - [ ] `codemint_<version>_darwin_amd64.tar.gz`
  - [ ] `codemint_<version>_windows_amd64.zip`
- [ ] Generate `SHA256SUMS`.

### Signing and Publishing
- [ ] Integrate macOS codesign + notarization in release workflow.
- [ ] Integrate Windows Authenticode signing.
- [ ] Publish GitHub Release assets on tag `v*`.
- [ ] Auto-generate release notes sections:
  - [ ] new commands
  - [ ] breaking changes
  - [ ] known issues
  - [ ] install/upgrade snippets

### Exit Criteria (M3)
- [ ] One-command tag-based release pipeline produces signed artifacts.

## Phase F: Distribution Channels (Week 4)

### Homebrew
- [ ] Create Homebrew tap repo and formula.
- [ ] Verify checksum and binary URL wiring.
- [ ] Validate `brew install` and `brew upgrade`.

### winget
- [ ] Create winget manifest for versioned release.
- [ ] Submit and validate install/upgrade paths.

### Exit Criteria (M4)
- [ ] Install/upgrade works from GitHub, Homebrew, and winget.

## Coordination with code-mint-ai (backend)

- Install script **source of truth** is this repo’s `scripts/install.sh` (and root `install.sh`). The app copies it to `scripts/install-cli.sh` and serves at `/cli/install.sh`; keep in sync when changing install behavior.
- Backend gaps and API contracts: see app repo’s **docs/CLI_GAPS.md** and **docs/cli-integration.md**.
- Feature ideas and roadmap: **FEATURES_ROADMAP.md** (this repo).

## Backend Work Required Before CLI Freeze

## Mandatory
- [ ] Fix/confirm login `next` redirect continuity for `/cli-auth` (in app repo).
- [ ] Keep `POST /api/auth/cli-token` stable.

## Recommended for v1
- [ ] `GET /api/auth/cli-token` list issued CLI tokens.
- [ ] `DELETE /api/auth/cli-token/:id` revoke a token.
- [ ] Add token metadata: `createdAt`, `lastUsedAt`, `revokedAt`, `name`, `userId`.
- [ ] Optional token `scopes` and `expiresAt` support.

## CI/CD Specification

## Trigger
- [ ] On Git tag matching `v*`.

## Pipeline Stages
- [ ] lint + unit tests
- [ ] integration tests (auth + API command smoke)
- [ ] cross-build
- [ ] package archives
- [ ] checksums
- [ ] sign artifacts
- [ ] publish release and notes

## Versioning
- [ ] Follow SemVer (`v1.0.0`, `v1.0.1`, ...).

## Testing Strategy

## Unit Tests
- [ ] command parsing
- [ ] config precedence
- [ ] secure storage abstraction
- [ ] API error mapping

## Integration Tests
- [ ] login callback roundtrip
- [ ] whoami valid/invalid token
- [ ] items search filters/pagination
- [ ] org list path

## Manual Matrix
- [ ] macOS arm64
- [ ] macOS amd64
- [ ] Windows 11 amd64
- [ ] Browser handoff: Chrome, Edge, Safari
- [ ] Corporate proxy scenario

## Regression Gate
- [ ] Do not release unless auth flow passes on macOS + Windows families.

## Documentation Deliverables
- [ ] `docs/quickstart.md`
- [ ] `docs/commands.md`
- [ ] `docs/troubleshooting.md`
- [ ] `docs/enterprise.md`

Troubleshooting must include:
- [ ] browser did not open
- [ ] callback port blocked
- [ ] token revoked/expired

## Day-by-Day Schedule (4 Weeks)

## Week 1
- Day 1: Repo scaffold, Cobra root, config/env loading.
- Day 2: HTTP client, error model, output mode foundation.
- Day 3: Auth callback server + browser open/fallback.
- Day 4: Token verification + secure storage abstraction.
- Day 5: `login`, `whoami`, `logout` E2E on macOS.

## Week 2
- Day 6: Windows secure storage path + auth tests.
- Day 7: `items search` command and API wiring.
- Day 8: `org list` command and output polishing.
- Day 9: JSON mode schema stabilization + integration tests.
- Day 10: Cross-platform bug fixes, milestone review.

## Week 3
- Day 11: retries/timeouts/auth error hardening.
- Day 12: revoked-token UX and redaction audits.
- Day 13: release workflow (`v*`), cross-build artifacts.
- Day 14: checksums + signing integration.
- Day 15: internal alpha release and validation.

## Week 4
- Day 16: private beta via GitHub Releases.
- Day 17: Homebrew tap + formula.
- Day 18: winget manifest submission.
- Day 19: docs finalization and upgrade guides.
- Day 20: GA release readiness + go/no-go checklist.

## Go/No-Go Checklist for GA
- [ ] Auth flow stable on macOS + Windows.
- [ ] Signed binaries and checksums published.
- [ ] Install and upgrade verified on GitHub/Homebrew/winget.
- [ ] Token lifecycle and revoke behavior validated.
- [ ] Critical docs published and reviewed.

## Post-GA Cadence
- [ ] Weekly patch releases initially.
- [ ] Biweekly minor feature releases initially.
- [ ] Track auth-flow and install-friction issues as highest priority.
