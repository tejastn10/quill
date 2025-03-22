#!/usr/bin/env bash

set -e

# Determine OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Set binary name based on OS and ARCH
if [[ "$OS" == "linux" ]]; then
  BINARY="quill-linux-amd64"
elif [[ "$OS" == "darwin" ]]; then
  BINARY="quill-darwin-arm64"
elif [[ "$OS" == "msys" || "$OS" == "cygwin" ]]; then
  BINARY="quill-windows-amd64.exe"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

VERSION=$(curl -s https://api.github.com/repos/tejastn10/quill/releases/latest | jq -r '.tag_name')
BASE_URL="https://github.com/tejastn10/quill/releases/download/$VERSION"
FILE="$BINARY.zip"

# Download the binary
echo "Downloading $BINARY from $BASE_URL/$FILE ..."
curl -LO "$BASE_URL/$FILE"

# Extract if on Unix-like system
if [[ "$OS" != "msys" && "$OS" != "cygwin" ]]; then
  unzip "$FILE" -d quill_temp
  BINARY_PATH=$(find quill_temp -type f -name "$BINARY")
else
  # Windows users need a different extraction method
  BINARY_PATH="$BINARY"
fi

# Install: move binary to /usr/local/bin (requires sudo)
echo "Installing quill to /usr/local/bin ..."
sudo mv "$BINARY_PATH" /usr/local/bin/quill
sudo chmod +x /usr/local/bin/quill

# Cleanup
rm -rf "$FILE" quill_temp

echo "quill installed successfully!"
