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

install_binary() {
  src_path="$1"
  dest_path="${INSTALL_DIR}/${BINARY}"
  if [ -w "$INSTALL_DIR" ]; then
    install -m 0755 "$src_path" "$dest_path"
  elif command -v sudo >/dev/null 2>&1; then
    sudo install -m 0755 "$src_path" "$dest_path"
  else
    echo "cannot write to ${INSTALL_DIR}; run as root or set CODEMINT_INSTALL_DIR" >&2
    exit 1
  fi
  echo "Installed ${BINARY} to ${dest_path}"
  "$dest_path" version || true
}

resolve_latest_tag() {
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null \
    | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
    | head -n 1
}

download_release_archive() {
  tag="$1"
  archive_path="$2"
  sums_path="$3"

  version_no_v="${tag#v}"
  asset="${BINARY}_${version_no_v}_${OS}_${ARCH}.tar.gz"
  base_url="https://github.com/${REPO}/releases/download/${tag}"

  echo "Downloading ${asset} from ${REPO} (${tag})"
  if ! curl -fL "${base_url}/${asset}" -o "$archive_path"; then
    return 1
  fi

  if curl -fsSL "${base_url}/SHA256SUMS" -o "$sums_path"; then
    expected="$(grep " ${asset}\$" "$sums_path" | awk '{print $1}' | head -n 1 || true)"
    if [ -n "$expected" ]; then
      actual=""
      if command -v sha256sum >/dev/null 2>&1; then
        actual="$(sha256sum "$archive_path" | awk '{print $1}')"
      elif command -v shasum >/dev/null 2>&1; then
        actual="$(shasum -a 256 "$archive_path" | awk '{print $1}')"
      fi
      if [ -n "$actual" ] && [ "$actual" != "$expected" ]; then
        echo "checksum mismatch for ${asset}" >&2
        echo "expected: $expected" >&2
        echo "actual:   $actual" >&2
        exit 1
      fi
    fi
  fi

  tar -xzf "$archive_path" -C "$TMP_DIR"
  if [ ! -f "${TMP_DIR}/${BINARY}" ]; then
    echo "archive did not contain expected binary: ${BINARY}" >&2
    exit 1
  fi
  install_binary "${TMP_DIR}/${BINARY}"
  return 0
}

build_from_source() {
  ref="$1"
  need_cmd go
  need_cmd find

  if [ "$ref" = "main" ]; then
    src_url="https://codeload.github.com/${REPO}/tar.gz/refs/heads/main"
  else
    src_url="https://codeload.github.com/${REPO}/tar.gz/refs/tags/${ref}"
  fi

  echo "Building ${BINARY} from source (${REPO}@${ref})"
  curl -fL "$src_url" -o "${TMP_DIR}/source.tar.gz"
  tar -xzf "${TMP_DIR}/source.tar.gz" -C "$TMP_DIR"

  src_root="$(find "$TMP_DIR" -mindepth 1 -maxdepth 1 -type d | head -n 1)"
  if [ -z "$src_root" ] || [ ! -f "${src_root}/go.mod" ]; then
    echo "failed to extract source archive for ${REPO}@${ref}" >&2
    exit 1
  fi

  (cd "$src_root" && go build -o "${TMP_DIR}/${BINARY}" .)
  install_binary "${TMP_DIR}/${BINARY}"
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

TMP_DIR="$(mktemp -d)"
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

ARCHIVE_PATH="${TMP_DIR}/release.tar.gz"
SUMS_PATH="${TMP_DIR}/SHA256SUMS"

if [ "$VERSION" = "latest" ]; then
  tag="$(resolve_latest_tag || true)"
  if [ -n "$tag" ]; then
    if download_release_archive "$tag" "$ARCHIVE_PATH" "$SUMS_PATH"; then
      exit 0
    fi
    echo "Latest release assets not available; falling back to source build from main" >&2
  else
    echo "No published release found for ${REPO}; falling back to source build from main" >&2
  fi
  build_from_source "main"
  exit 0
fi

case "$VERSION" in
  v*) TAG="$VERSION" ;;
  *) TAG="v$VERSION" ;;
esac

if download_release_archive "$TAG" "$ARCHIVE_PATH" "$SUMS_PATH"; then
  exit 0
fi

echo "Release asset not found for ${TAG}; attempting source build from tag" >&2
build_from_source "$TAG"
