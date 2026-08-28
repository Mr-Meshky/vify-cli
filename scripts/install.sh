#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${CYAN}⚡ Installing Vify CLI (Terminal VPN Client)...${NC}"

REPO="Mr-Meshky/vify-cli"
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
    darwin|linux) ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

LATEST_TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name"' | head -1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
if [ -z "$LATEST_TAG" ]; then
    LATEST_TAG="v1.0.0"
fi
VERSION="${LATEST_TAG#v}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/vify_${VERSION}_${OS}_${ARCH}.tar.gz"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT
cd "$TMP_DIR"

install_binary() {
    # $1 = path to the built/extracted "vify" binary
    if [ "$(id -u)" = "0" ]; then
        mv "$1" /usr/local/bin/vify
        echo "/usr/local/bin/vify"
        return
    fi
    if sudo mv "$1" /usr/local/bin/vify 2>/dev/null; then
        echo "/usr/local/bin/vify"
        return
    fi
    mkdir -p "$HOME/.local/bin"
    mv "$1" "$HOME/.local/bin/vify"
    echo "$HOME/.local/bin/vify"
    case ":$PATH:" in
        *":$HOME/.local/bin:"*) ;;
        *) echo -e "${RED}Note: $HOME/.local/bin is not on your PATH. Add it, e.g.:${NC}" >&2
           echo "  echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc" >&2 ;;
    esac
}

DOWNLOAD_OK=false
if command -v curl >/dev/null 2>&1; then
    if curl -fsSL "$DOWNLOAD_URL" -o vify.tar.gz; then
        DOWNLOAD_OK=true
    fi
elif command -v wget >/dev/null 2>&1; then
    if wget -q "$DOWNLOAD_URL" -O vify.tar.gz; then
        DOWNLOAD_OK=true
    fi
fi

if [ "$DOWNLOAD_OK" = true ] && [ -s vify.tar.gz ]; then
    tar -xzf vify.tar.gz
    INSTALL_PATH="$(install_binary "$TMP_DIR/vify")"
else
    echo "Falling back to building from source..."
    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}✗ Could not download a prebuilt binary for ${OS}/${ARCH}, and Go is not installed.${NC}"
        echo "Install Go from https://go.dev/dl/ and re-run this script, or grab a binary manually from:"
        echo "  https://github.com/${REPO}/releases"
        exit 1
    fi
    GOBIN="$TMP_DIR" go install "github.com/${REPO}@latest"
    mv "$TMP_DIR/vify-cli" "$TMP_DIR/vify"
    INSTALL_PATH="$(install_binary "$TMP_DIR/vify")"
fi

echo -e "${GREEN}✓ Vify CLI installed successfully at ${INSTALL_PATH}${NC}"
echo "Run 'vify connect' to establish VPN connection."
