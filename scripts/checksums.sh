#!/usr/bin/env bash
set -euo pipefail

DIST_DIR="${1:-dist}"
PATTERN="${2:-*}"
cd "$DIST_DIR"

shopt -s nullglob
files=( ${PATTERN}.tar.gz ${PATTERN}.zip )
if [ ${#files[@]} -eq 0 ]; then
  echo "no release archives found for pattern: ${PATTERN}" >&2
  exit 1
fi
sha256sum "${files[@]}" > SHA256SUMS
