#!/usr/bin/env sh

set -e

REPO="Abhishek-Krishna-A-M/gpad"
BIN_NAME="gpad"
INSTALL_DIR="/usr/local/bin"

echo ">>> Detecting OS and Architecture..."

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)   OS="linux" ;;
    Darwin*)  OS="macos" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

echo ">>> OS: $OS"
echo ">>> ARCH: $ARCH"

echo ">>> Fetching latest version tag..."
TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep -oP '"tag_name": "\K(.*)(?=")')

if [ -z "$TAG" ]; then
    echo "Failed to fetch latest release tag"
    exit 1
fi

echo ">>> Latest version: $TAG"

FILE="${BIN_NAME}-${OS}-${ARCH}"

DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${TAG}/${FILE}"

echo ">>> Downloading: $DOWNLOAD_URL"
curl -L "$DOWNLOAD_URL" -o "$BIN_NAME"

echo ">>> Making executable..."
chmod +x "$BIN_NAME"

echo ">>> Installing to $INSTALL_DIR..."
sudo mv "$BIN_NAME" "$INSTALL_DIR/"

echo ">>> Installation complete!"

echo ""
echo "Run 'gpad init' to start."

