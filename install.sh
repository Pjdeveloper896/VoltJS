#!/usr/bin/env bash

set -e

REPO="Pjdeveloper896/VoltJs"
BIN_NAME="voltjs"

# Detect OS and arch
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Termux detection
if [ "$PREFIX" != "" ]; then
  # We are in Termux
  INSTALL_DIR="$PREFIX/bin"
else
  # Default to /usr/local/bin for other Linux/macOS
  INSTALL_DIR="/usr/local/bin"
fi

echo "Installing to $INSTALL_DIR ..."

LATEST_RELEASE_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")

TAG_NAME=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"tag_name": "\K(.*)(?=")')

if [ -z "$TAG_NAME" ]; then
  echo "Could not find latest release tag."
  exit 1
fi

echo "Latest release: $TAG_NAME"

FILENAME="${BIN_NAME}-${OS}-${ARCH}"

DOWNLOAD_URL=$(echo "$LATEST_RELEASE_JSON" | grep -oP '"browser_download_url": "\K(.*)(?=")' | grep "$FILENAME")

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Could not find a binary for ${FILENAME} in the latest release."
  exit 1
fi

echo "Downloading $FILENAME from $DOWNLOAD_URL ..."

TMP_FILE="/tmp/$FILENAME"
curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"

chmod +x "$TMP_FILE"

# Move without sudo in Termux or normal user context
mv "$TMP_FILE" "$INSTALL_DIR/$BIN_NAME"

echo "Installed $BIN_NAME successfully to $INSTALL_DIR!"

# Check if INSTALL_DIR is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "Warning: $INSTALL_DIR is not in your PATH."
  echo "Add this line to your shell config (e.g. ~/.bashrc or ~/.zshrc):"
  echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "You can now run: $BIN_NAME --help"
