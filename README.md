# CodeMint CLI

`codemint` is the official command-line interface for CodeMint.
It supports authentication, catalog search, repository-aware recommendations, and local installation/synchronization of rules and skills.

## Highlights

- Browser-based login with secure local token storage
- Search and discovery commands for items, org data, and catalog recommendations
- Repository scanning (`scan`) plus guided suggestions (`suggest`)
- Local install lifecycle for rules/skills (`add`, `list`, `remove`, `sync`)
- Tool-aware installs with per-repo default AI tool settings
- Scriptable output with `--json`

## Installation

### Option 1: Download and install prebuilt packages

Release packages:

- `codemint_<version>_darwin_arm64.tar.gz`
- `codemint_<version>_darwin_amd64.tar.gz`
- `codemint_<version>_windows_amd64.zip`
- `SHA256SUMS`

Download and install on macOS/Linux:

```bash
# set the version you want
VERSION="0.1.0"

# choose ONE package that matches your machine
curl -fL -o codemint.tar.gz \
  "https://github.com/codemint/codemint-cli/releases/download/v${VERSION}/codemint_${VERSION}_darwin_arm64.tar.gz"

# extract and install
tar -xzf codemint.tar.gz
install -m 0755 codemint /usr/local/bin/codemint

# verify
codemint version
```

Download and install on Windows (PowerShell):

```powershell
$Version = "0.1.0"
Invoke-WebRequest `
  -Uri "https://github.com/codemint/codemint-cli/releases/download/v$Version/codemint_${Version}_windows_amd64.zip" `
  -OutFile "codemint.zip"
Expand-Archive -Path "codemint.zip" -DestinationPath ".\codemint-bin" -Force
Move-Item ".\codemint-bin\codemint.exe" "$env:USERPROFILE\bin\codemint.exe" -Force

# Add $env:USERPROFILE\bin to PATH if needed, then verify:
codemint version
```

### Option 2: Build from source

Requirements:
- Go `1.22+`

```bash
# install Go module dependencies
go mod download

# build
make build
./bin/codemint version
```

## Quickstart

```bash
# 1) authenticate
codemint auth login

# 2) verify identity
codemint auth whoami

# 3) search items
codemint items search --q "widget"

# 4) list organizations
codemint org list
```

For non-default environments:

```bash
codemint --base-url https://app.codemint.app auth login
```

## Run Locally

Build once, then run commands from the local binary:

```bash
make build
./bin/codemint version
./bin/codemint auth whoami
```

You can also use `go run` during development:

```bash
go run . version
go run . items search --q "widget"
```

## Test Locally

Run the full test suite:

```bash
make test
```

Run tests directly with Go:

```bash
go test ./...
```

Run tests for one package:

```bash
go test ./internal/scan -v
```

## Command Overview

| Area | Commands |
|---|---|
| Auth | `auth login`, `auth whoami`, `auth logout` |
| Search | `items search`, `org list` |
| Repo analysis | `scan [path]`, `suggest [--path <dir>] [--type rule\|skill]` |
| Install lifecycle | `add @rule/<slug>\|@skill/<slug>`, `list`, `remove <ref>`, `sync` |
| Tool settings | `tool list`, `tool current`, `tool set <name>` |
| Diagnostics | `doctor`, `version` |

Run `codemint <command> --help` for full usage and flags.

## Global Flags

- `--json` output JSON for automation
- `--base-url` override API base URL
- `--profile` choose credential/profile namespace
- `--config` custom config file path
- `--debug` enable debug logging

## Configuration

Default config file:

- `~/.config/codemint/config.json`

Environment variables:

- `CODEMINT_BASE_URL`
- `CODEMINT_PROFILE`

Override precedence (highest to lowest):

1. CLI flags (`--base-url`, `--profile`)
2. Environment variables
3. Config file
4. Built-in defaults

## Change Base URL (Useful for Multi-Platform/Test Environments)

Use one of these methods:

1. Per command (best for quick testing):

```bash
codemint --base-url https://staging.codemint.app auth login
```

2. Environment variable (best for local session testing):

macOS/Linux:

```bash
export CODEMINT_BASE_URL=https://staging.codemint.app
codemint auth whoami
```

Windows PowerShell:

```powershell
$env:CODEMINT_BASE_URL = "https://staging.codemint.app"
codemint auth whoami
```

3. Config file (best for stable local setup):

- File: `~/.config/codemint/config.json`
- Example:

```json
{
  "base_url": "https://staging.codemint.app",
  "profile": "default"
}
```

For testing multiple platforms/environments, pair base URLs with profiles:

```bash
codemint --base-url https://staging.codemint.app --profile staging auth login
codemint --base-url https://app.codemint.app --profile prod auth login

codemint --profile staging auth whoami
codemint --profile prod auth whoami
```

## Development

```bash
# format source
make fmt

# run tests
make test

# build local binary
make build
```

Generate local release artifacts:

```bash
make release-dry-run VERSION=0.1.0
```

## Documentation

- Quickstart: `docs/quickstart.md`
- Commands: `docs/commands.md`
- Troubleshooting: `docs/troubleshooting.md`

## Troubleshooting

- Browser did not open: copy the printed login URL and open it manually.
- Callback port blocked: allow loopback traffic to `127.0.0.1`, then retry.
- Token revoked/expired: run `codemint auth login` again.
