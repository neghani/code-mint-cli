#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
DIST_DIR="${2:-dist}"

if [[ -z "${VERSION}" ]]; then
  echo "Usage: $0 <version> [dist_dir]"
  echo "Example: $0 0.1.0"
  exit 1
fi

if ! command -v gh >/dev/null 2>&1; then
  echo "GitHub CLI (gh) is required. Install it from https://cli.github.com/"
  exit 1
fi

TAG="v${VERSION}"
TITLE="CodeMint CLI ${TAG}"

files=(
  "${DIST_DIR}/codemint_${VERSION}_darwin_arm64.tar.gz"
  "${DIST_DIR}/codemint_${VERSION}_darwin_amd64.tar.gz"
  "${DIST_DIR}/codemint_${VERSION}_linux_arm64.tar.gz"
  "${DIST_DIR}/codemint_${VERSION}_linux_amd64.tar.gz"
  "${DIST_DIR}/codemint_${VERSION}_windows_amd64.zip"
  "${DIST_DIR}/SHA256SUMS"
)

for f in "${files[@]}"; do
  if [[ ! -f "${f}" ]]; then
    echo "Missing release artifact: ${f}"
    exit 1
  fi
done

if ! gh auth status >/dev/null 2>&1; then
  echo "Run 'gh auth login' before creating a release."
  exit 1
fi

echo "Creating GitHub release ${TAG} with artifacts from ${DIST_DIR}"
gh release create "${TAG}" \
  "${files[@]}" \
  --title "${TITLE}" \
  --generate-notes

echo "Release created: ${TAG}"
