#!/bin/bash
#
# Hanzo S3 installer — installs the `s3` CLI (S3-compatible object storage).
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/hanzoai/s3/main/install.sh | bash
#   ./install.sh --dir /usr/local/bin --version latest
#
# Methods (auto-detected, in order):
#   1. go install   — if a Go toolchain is present (builds from source)
#   2. docker        — else extracts the prebuilt binary from ghcr.io/hanzoai/s3
#
set -euo pipefail

MODULE="github.com/hanzoai/s3/s3"
IMAGE="ghcr.io/hanzoai/s3"
VERSION="latest"
INSTALL_DIR="/usr/local/bin"

while [ $# -gt 0 ]; do
  case "$1" in
    --dir)     INSTALL_DIR="$2"; shift 2 ;;
    --version) VERSION="$2"; shift 2 ;;
    --help|-h) sed -n '/^# Usage:/,/^$/p' "$0" | sed 's/^# \?//'; exit 0 ;;
    *) echo "unknown option: $1 (try --help)" >&2; exit 1 ;;
  esac
done

info() { echo "[info] $*"; }
ok()   { echo "[ok] $*"; }

place() { # $1 = built binary -> $INSTALL_DIR/s3
  if [ -w "$INSTALL_DIR" ]; then install -m755 "$1" "$INSTALL_DIR/s3"
  else info "elevated permissions needed for $INSTALL_DIR"; sudo install -m755 "$1" "$INSTALL_DIR/s3"; fi
  ok "installed s3 -> $INSTALL_DIR/s3"
  "$INSTALL_DIR/s3" version
}

if command -v go >/dev/null 2>&1; then
  info "Go found — go install ${MODULE}@${VERSION}"
  BIN="$(go env GOBIN)"; [ -z "$BIN" ] && BIN="$(go env GOPATH)/bin"
  go install "${MODULE}@${VERSION}"
  place "$BIN/s3"
elif command -v docker >/dev/null 2>&1; then
  info "no Go — extracting s3 from ${IMAGE}:${VERSION}"
  id="$(docker create "${IMAGE}:${VERSION}")"
  tmp="$(mktemp -d)"; docker cp "$id:/usr/bin/s3" "$tmp/s3"; docker rm "$id" >/dev/null
  place "$tmp/s3"
else
  echo "[error] need either a Go toolchain or Docker to install s3" >&2
  echo "        Go:     go install ${MODULE}@latest" >&2
  echo "        Docker: docker run --rm ${IMAGE}:latest version" >&2
  exit 1
fi
