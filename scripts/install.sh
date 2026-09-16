#!/usr/bin/env bash
set -euo pipefail

# =========================
# Configuration
# =========================

VERSION="${1:-latest}"          # release tag (e.g. v1.0.0) or "latest"
REPO="othmaneBakkass/cv_gen"
BIN_NAME="cv_gen"

# =========================
# Detect OS and architecture
# =========================

case "$(uname -s)" in
    Linux*)   OS="linux" ;;
    Darwin*)  OS="darwin" ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT) OS="windows" ;;
    *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

# =========================
# Resolve download URL
# =========================

ASSET="cv_gen_${OS}_${ARCH}"
EXT=""
if [ "$OS" = "windows" ]; then
    ASSET="${ASSET}.exe"
    EXT=".exe"
fi

if [ "$VERSION" = "latest" ]; then
    URL="https://github.com/${REPO}/releases/latest/download/${ASSET}"
else
    URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"
fi

# =========================
# Choose install directory
# =========================

if [ "$OS" = "windows" ]; then
    BIN_DIR="${BIN_DIR:-$HOME/bin}"
else
    BIN_DIR="${BIN_DIR:-/usr/local/bin}"
fi

if [ ! -d "$BIN_DIR" ]; then
    mkdir -p "$BIN_DIR" 2>/dev/null || {
        echo "Cannot create $BIN_DIR (try: sudo BIN_DIR=$BIN_DIR $0 $VERSION)" >&2
        exit 1
    }
fi

if [ ! -w "$BIN_DIR" ]; then
    echo "No write permission for $BIN_DIR (try running with sudo)." >&2
    exit 1
fi

# =========================
# Download
# =========================

DEST="${BIN_DIR}/${BIN_NAME}${EXT}"
echo "Downloading ${URL}"
if ! curl -fSL -o "$DEST" "$URL"; then
    rm -f "$DEST"
    echo "Download failed. Check that release '${VERSION}' contains asset '${ASSET}'." >&2
    exit 1
fi

if [ "$OS" != "windows" ]; then
    chmod +x "$DEST"
fi

echo "Installed ${BIN_NAME} to ${DEST}"
echo "Make sure ${BIN_DIR} is in your PATH."
