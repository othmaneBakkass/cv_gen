#!/bin/bash
set -e

# =========================
# Configuration
# =========================

VERSION=${1:-latest}           # Default to latest release if no version specified
REPO="othmaneBakkass/cv_gen"
BIN_NAME="cv_gen"

# =========================
# Detect OS and Architecture
# =========================

OS=$(uname | tr '[:upper:]' '[:lower:]')

# =========================
# Determine binary name for download
# =========================

if [[ "$OS" == "linux" ]]; then
    BIN_DOWNLOAD_NAME="cv_gen"
    BIN_DIR="/usr/local/bin"
    if [ ! -w "$BIN_DIR" ]; then
        echo "Warning: You may need to run this script with sudo to install in $BIN_DIR"
    fi
elif [[ "$OS" == "darwin" ]]; then
    BIN_DOWNLOAD_NAME="cv_gen_darwin"
    BIN_DIR="/usr/local/bin"
    if [ ! -w "$BIN_DIR" ]; then
        echo "Warning: You may need to run this script with sudo to install in $BIN_DIR"
    fi
elif [[ "$OS" == "windows" ]]; then
    BIN_DOWNLOAD_NAME="cv_gen.exe"
    BIN_DIR="$HOME/bin"
    mkdir -p "$BIN_DIR"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

# =========================
# Download binary
# =========================

URL="https://github.com/$REPO/releases/download/$VERSION/$BIN_DOWNLOAD_NAME"

echo "Downloading $URL..."
curl -sSL -o "$BIN_DIR/$BIN_NAME" "$URL"

# =========================
# Make executable (Linux/macOS)
# =========================

if [[ "$OS" != "windows" ]]; then
    chmod +x "$BIN_DIR/$BIN_NAME"
fi

# =========================
# Done
# =========================

echo "Installed $BIN_NAME successfully to $BIN_DIR"
echo "Make sure $BIN_DIR is in your PATH"
