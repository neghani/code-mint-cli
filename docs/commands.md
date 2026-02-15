# Commands

## Recommended first command for migration

- `codemint suggest --path . --type rule`
- `codemint suggest --path . --type skill`

## Global flags

- `--json`
- `--base-url`
- `--profile`
- `--config`
- `--debug`

## Auth

- `codemint auth login`
- `codemint auth whoami`
- `codemint auth logout`

## Items

- `codemint items search --q <query> [--type] [--tags] [--page] [--limit]`

## Org

- `codemint org list`

## Migration and install lifecycle

- `codemint scan [path]`
- `codemint suggest [--path <dir>] [--type rule|skill]`
- `codemint tool set <name>`
- `codemint add @rule/<slug>|@skill/<slug> [--tool <name>] [--dry-run]`
- `codemint list [--installed]`
- `codemint remove @rule/<slug>|@skill/<slug>`
- `codemint sync [--dry-run]`
