#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}⚡ Installing Vify CLI (Terminal VPN Client)...${NC}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

LATEST_RELEASE="v1.0.0"
DOWNLOAD_URL="https://github.com/Mr-Meshky/vify-cli/releases/download/${LATEST_RELEASE}/vify_${LATEST_RELEASE}_${OS}_${ARCH}.tar.gz"

TMP_DIR="$(mktemp -d)"
cd "$TMP_DIR"

if command -v curl >/dev/null 2>&1; then
    curl -sL "$DOWNLOAD_URL" -o vify.tar.gz || true
elif command -v wget >/dev/null 2>&1; then
    wget -q "$DOWNLOAD_URL" -O vify.tar.gz || true
fi

if [ -f vify.tar.gz ]; then
    tar -xzf vify.tar.gz
    sudo mv vify /usr/local/bin/vify || mv vify "$HOME/.local/bin/vify"
else
    echo "Falling back to building from source..."
    go install github.com/Mr-Meshky/vify-cli@latest
fi

echo -e "${GREEN}✓ Vify CLI installed successfully!${NC}"
echo "Run 'vify connect' to establish VPN connection."
