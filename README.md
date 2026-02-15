# CodeMint CLI

`codemint` is a command line client for CodeMint authentication, item search, and organization listing.

## Build

```bash
make build
```

## Test

```bash
make test
```

## Auth flow

```bash
codemint auth login
codemint auth whoami
codemint auth logout
```

## Domain commands

```bash
codemint items search --q "widget"
codemint org list
```

## Release artifacts

- `codemint_<version>_darwin_arm64.tar.gz`
- `codemint_<version>_darwin_amd64.tar.gz`
- `codemint_<version>_windows_amd64.zip`
- `SHA256SUMS`
