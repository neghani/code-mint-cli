# Quickstart

## Install

Download a package from GitHub Releases and install it.

### Fast install (works even before releases)

```bash
curl -fsSL https://raw.githubusercontent.com/codemint/codemint-cli/main/scripts/install.sh | sh
```

If you publish from your own repo/fork, use your fork's raw URL:

```bash
curl -fsSL https://raw.githubusercontent.com/neghani/code-mint-cli/main/scripts/install.sh | CODEMINT_REPO=neghani/code-mint-cli sh
```

### Linux amd64

```bash
VERSION="0.1.0"
curl -fL -o codemint.tar.gz \
  "https://github.com/codemint/codemint-cli/releases/download/v${VERSION}/codemint_${VERSION}_linux_amd64.tar.gz"
tar -xzf codemint.tar.gz
sudo install -m 0755 codemint /usr/local/bin/codemint
codemint version
```

### Linux arm64

```bash
VERSION="0.1.0"
curl -fL -o codemint.tar.gz \
  "https://github.com/codemint/codemint-cli/releases/download/v${VERSION}/codemint_${VERSION}_linux_arm64.tar.gz"
tar -xzf codemint.tar.gz
sudo install -m 0755 codemint /usr/local/bin/codemint
codemint version
```

### macOS (auto detect arch)

```bash
VERSION="0.1.0"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" && exit 1 ;;
esac
curl -fL -o codemint.tar.gz \
  "https://github.com/codemint/codemint-cli/releases/download/v${VERSION}/codemint_${VERSION}_darwin_${ARCH}.tar.gz"
tar -xzf codemint.tar.gz
sudo install -m 0755 codemint /usr/local/bin/codemint
codemint version
```

### Windows (PowerShell)

```powershell
$Version = "0.1.0"
Invoke-WebRequest `
  -Uri "https://github.com/codemint/codemint-cli/releases/download/v$Version/codemint_${Version}_windows_amd64.zip" `
  -OutFile "codemint.zip"
Expand-Archive -Path "codemint.zip" -DestinationPath ".\codemint-bin" -Force
Move-Item ".\codemint-bin\codemint.exe" "$env:USERPROFILE\bin\codemint.exe" -Force
codemint version
```

## Configure

`codemint --base-url https://app.codemint.app auth login`

Environment overrides:
- `CODEMINT_BASE_URL`
- `CODEMINT_PROFILE`

## Recommended migration flow (right after install)

```bash
# 1) login
codemint auth login

# 2) set your default AI tool once per repo
codemint tool set cursor

# 3) scan + suggest rules for this repository
codemint suggest --path . --type rule

# 4) install one suggested rule (replace slug from output)
codemint add @rule/<slug>
```

If you skip `tool set`, `codemint add` will ask you to pick a tool interactively.

## Core commands

- `codemint auth login`
- `codemint auth whoami`
- `codemint suggest --path . --type rule`
- `codemint suggest --path . --type skill`
- `codemint add @rule/<slug>`
- `codemint items search --q "term"`
- `codemint org list`
