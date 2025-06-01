#!/usr/bin/env bash

set -e

REPO="Pjdeveloper896/VoltJs"
BIN_NAME="voltjs"

# Detect OS and architecture
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Normalize architecture
if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "❌ Unsupported architecture: $ARCH"
  exit 1
fi

# Determine install directory
if [[ -n "$PREFIX" ]]; then
  INSTALL_DIR="$PREFIX/bin"  # Termux
else
  INSTALL_DIR="/usr/local/bin"  # Linux/macOS
fi

echo "📦 Installing to $INSTALL_DIR ..."

# Fetch latest release info
LATEST_RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")

TAG_NAME=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"tag_name": "\K(.*)(?=")')
if [[ -z "$TAG_NAME" ]]; then
  echo "❌ Could not fetch the latest release tag."
  exit 1
fi

echo "🔖 Latest release: $TAG_NAME"

FILENAME="${BIN_NAME}-${OS}-${ARCH}"
DOWNLOAD_URL=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"browser_download_url": "\K(.*)(?=")' | grep "$FILENAME")

if [[ -z "$DOWNLOAD_URL" ]]; then
  echo "❌ No binary found for ${FILENAME} in latest release."
  exit 1
fi

echo "⬇️  Downloading $FILENAME from:"
echo "    $DOWNLOAD_URL"

TMP_FILE="/tmp/$FILENAME"
curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"

chmod +x "$TMP_FILE"

# Move binary to install directory
echo "⚙️  Installing..."
mv "$TMP_FILE" "$INSTALL_DIR/$BIN_NAME"

echo "✅ Installed $BIN_NAME to $INSTALL_DIR"

# Check if in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "⚠️  $INSTALL_DIR is not in your PATH."
  echo "➡️  Add this to your shell config (e.g., ~/.bashrc or ~/.zshrc):"
  echo "    export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "🚀 Run it with:"
echo "    $BIN_NAME --help"
