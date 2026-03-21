#!/usr/bin/env sh
set -e

REPO="Abhishek-Krishna-A-M/gpad"
INSTALL_DIR="/usr/local/bin"
TMP_DIR="$(mktemp -d)"

# ── cleanup on exit ───────────────────────────────────────────────────────────
cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

# ── check for Go ──────────────────────────────────────────────────────────────
if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed. Install it from https://go.dev/dl/ then re-run this script."
  exit 1
fi

# ── check for git ─────────────────────────────────────────────────────────────
if ! command -v git >/dev/null 2>&1; then
  echo "git is not installed. Install git then re-run this script."
  exit 1
fi

# ── clone and build ───────────────────────────────────────────────────────────
echo "Cloning gpad..."
git clone --depth=1 "https://github.com/${REPO}.git" "$TMP_DIR/gpad" 2>/dev/null

echo "Building gpad..."
cd "$TMP_DIR/gpad"
go build -ldflags="-s -w" -trimpath -o "$TMP_DIR/gpad_bin" ./cmd/gpad/

# ── install ───────────────────────────────────────────────────────────────────
if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP_DIR/gpad_bin" "${INSTALL_DIR}/gpad"
else
  echo "Installing to ${INSTALL_DIR} (needs sudo)..."
  sudo mv "$TMP_DIR/gpad_bin" "${INSTALL_DIR}/gpad"
fi

# ── done ──────────────────────────────────────────────────────────────────────
echo ""
echo "  gpad installed to ${INSTALL_DIR}/gpad"
echo ""
echo "  Get started:"
echo "    gpad today               open today's daily note"
echo "    gpad open my-note.md     create your first note"
echo "    gpad git init <url>      connect git sync (optional)"
echo ""
