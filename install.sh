#!/usr/bin/env sh
# Install deadsniper from GitHub Releases.
# Usage: curl -fsSL https://raw.githubusercontent.com/shutterscripter/deadsniper/main/install.sh | sh

set -e

REPO="shutterscripter/deadsniper"
BIN="deadsniper"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64)   ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported arch: $ARCH"; exit 1 ;;
esac

# Windows binary not installed via this script (run on Windows or use scoop/other)
if [ "$OS" = "windows" ]; then
  echo "Use GitHub Releases or scoop to install on Windows."
  exit 1
fi

URL="https://github.com/${REPO}/releases/latest/download/${BIN}-${OS}-${ARCH}"

INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$INSTALL_DIR"

echo "Downloading ${BIN} (${OS}/${ARCH})..."
curl -fsSL "$URL" -o "$INSTALL_DIR/$BIN"
chmod +x "$INSTALL_DIR/$BIN"

echo "Installed $BIN to $INSTALL_DIR"

if ! echo ":$PATH:" | grep -q ":$INSTALL_DIR:"; then
  echo "Add to your shell config (e.g. ~/.zshrc or ~/.bashrc):"
  echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo "Run: $BIN -u https://example.com"
