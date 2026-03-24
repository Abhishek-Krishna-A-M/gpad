#!/usr/bin/env sh
set -e

REPO="Abhishek-Krishna-A-M/gpad"
BINARY="gpad"
INSTALL_DIR="/usr/local/bin"

# ── detect OS and arch ───────────────────────────────────────────────────────

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux)  OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# ── fetch latest release tag ─────────────────────────────────────────────────

echo "Fetching latest gpad release..."

LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine latest release. Check your internet connection."
  exit 1
fi

echo "Latest release: $LATEST"

# ── download and install ─────────────────────────────────────────────────────

URL="https://github.com/${REPO}/releases/download/${LATEST}/gpad_${OS}_${ARCH}"

TMP="$(mktemp)"
echo "Downloading gpad ${LATEST} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"

# ── pick install location ─────────────────────────────────────────────────────

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP" "${INSTALL_DIR}/${BINARY}"
else
  echo "Installing to ${INSTALL_DIR} (needs sudo)..."
  sudo mv "$TMP" "${INSTALL_DIR}/${BINARY}"
fi

# ── verify ───────────────────────────────────────────────────────────────────

if command -v gpad >/dev/null 2>&1; then
  echo ""
  echo "  gpad $(gpad --version 2>/dev/null | head -1) installed successfully."
  echo ""
  echo "  Get started:"
  echo "    gpad today               open today's daily note"
  echo "    gpad open my-note.md     create your first note"
  echo "    gpad git init <url>      connect git sync (optional)"
  echo ""
else
  echo ""
  echo "  Installed to ${INSTALL_DIR}/gpad"
  echo "  Make sure ${INSTALL_DIR} is in your PATH."
  echo ""
fi
