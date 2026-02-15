#!/usr/bin/env bash
set -euo pipefail

DIST_DIR="${1:-dist}"
cd "$DIST_DIR"
sha256sum *.tar.gz *.zip > SHA256SUMS
