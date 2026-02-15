#!/bin/sh
set -eu

REPO="${CODEMINT_REPO:-codemint/codemint-cli}"
BINARY="${CODEMINT_BINARY:-codemint}"
INSTALL_DIR="${CODEMINT_INSTALL_DIR:-/usr/local/bin}"
VERSION="${CODEMINT_VERSION:-${1:-latest}}"

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd tar
need_cmd uname
need_cmd mktemp
need_cmd install

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
  darwin|linux) ;;
  *)
    echo "unsupported OS: $OS (supported: darwin, linux)" >&2
    exit 1
    ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *)
    echo "unsupported architecture: $ARCH (supported: amd64, arm64)" >&2
    exit 1
    ;;
esac

if [ "$VERSION" = "latest" ]; then
  TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n 1)"
  if [ -z "$TAG" ]; then
    echo "failed to resolve latest release tag from GitHub API for ${REPO}" >&2
    exit 1
  fi
else
  case "$VERSION" in
    v*) TAG="$VERSION" ;;
    *) TAG="v$VERSION" ;;
  esac
fi

VERSION_NO_V="${TAG#v}"
ASSET="${BINARY}_${VERSION_NO_V}_${OS}_${ARCH}.tar.gz"
BASE_URL="https://github.com/${REPO}/releases/download/${TAG}"

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

ARCHIVE_PATH="${TMP_DIR}/${ASSET}"
SUMS_PATH="${TMP_DIR}/SHA256SUMS"

echo "Downloading ${ASSET} from ${REPO} (${TAG})"
curl -fL "${BASE_URL}/${ASSET}" -o "$ARCHIVE_PATH"

if curl -fsSL "${BASE_URL}/SHA256SUMS" -o "$SUMS_PATH"; then
  EXPECTED="$(grep " ${ASSET}\$" "$SUMS_PATH" | awk '{print $1}' | head -n 1 || true)"
  if [ -n "$EXPECTED" ]; then
    ACTUAL=""
    if command -v sha256sum >/dev/null 2>&1; then
      ACTUAL="$(sha256sum "$ARCHIVE_PATH" | awk '{print $1}')"
    elif command -v shasum >/dev/null 2>&1; then
      ACTUAL="$(shasum -a 256 "$ARCHIVE_PATH" | awk '{print $1}')"
    fi
    if [ -n "$ACTUAL" ] && [ "$ACTUAL" != "$EXPECTED" ]; then
      echo "checksum mismatch for ${ASSET}" >&2
      echo "expected: $EXPECTED" >&2
      echo "actual:   $ACTUAL" >&2
      exit 1
    fi
  fi
fi

tar -xzf "$ARCHIVE_PATH" -C "$TMP_DIR"
if [ ! -f "${TMP_DIR}/${BINARY}" ]; then
  echo "archive did not contain expected binary: ${BINARY}" >&2
  exit 1
fi

DEST="${INSTALL_DIR}/${BINARY}"
if [ -w "$INSTALL_DIR" ]; then
  install -m 0755 "${TMP_DIR}/${BINARY}" "$DEST"
elif command -v sudo >/dev/null 2>&1; then
  sudo install -m 0755 "${TMP_DIR}/${BINARY}" "$DEST"
else
  echo "cannot write to ${INSTALL_DIR}; run as root or set CODEMINT_INSTALL_DIR" >&2
  exit 1
fi

echo "Installed ${BINARY} to ${DEST}"
"$DEST" version || true
