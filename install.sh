#!/bin/bash
set -e

# Claude Code TTS Plugin Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/currenthandle/claude-code-tts/main/install.sh | bash

REPO="currenthandle/claude-code-tts"
INSTALL_DIR="$HOME/.claude/plugins/claude-code-tts"

echo "Installing Claude Code TTS Plugin..."
echo ""

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux)  OS="linux" ;;
    *)
        echo "Error: Unsupported OS: $OS"
        echo "Supported: macOS (darwin), Linux"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo "Error: Unsupported architecture: $ARCH"
        echo "Supported: x86_64/amd64, arm64/aarch64"
        exit 1
        ;;
esac

echo "Detected: ${OS}/${ARCH}"

# Get latest release tag
echo "Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
    echo "Error: Could not determine latest release."
    echo "Check https://github.com/$REPO/releases"
    exit 1
fi
echo "Latest release: $LATEST"

# Download binaries
mkdir -p "$INSTALL_DIR/bin"

echo "Downloading tts-server..."
curl -fsSL "https://github.com/$REPO/releases/download/$LATEST/tts-server-${OS}-${ARCH}" -o "$INSTALL_DIR/bin/tts-server"
chmod +x "$INSTALL_DIR/bin/tts-server"

echo "Downloading speak-text..."
curl -fsSL "https://github.com/$REPO/releases/download/$LATEST/speak-text-${OS}-${ARCH}" -o "$INSTALL_DIR/bin/speak-text"
chmod +x "$INSTALL_DIR/bin/speak-text"

echo ""
echo "Installation complete! Binaries installed to $INSTALL_DIR/bin/"
echo ""
echo "Next steps:"
echo ""
echo "  1. Set your API key (choose one):"
echo ""
echo "     # OpenAI:"
echo "     export OPENAI_API_KEY=\"sk-...\""
echo ""
echo "     # Azure OpenAI:"
echo "     export TTS_PROVIDER=azure"
echo "     export AZURE_OPENAI_API_KEY=\"your-key\""
echo "     export AZURE_OPENAI_ENDPOINT=\"https://your-resource.openai.azure.com\""
echo "     export AZURE_OPENAI_DEPLOYMENT=\"gpt-4o-mini-tts\""
echo "     export AZURE_OPENAI_API_VERSION=\"2025-03-01-preview\""
echo ""
echo "  2. Register with Claude Code:"
echo "     claude mcp add tts $INSTALL_DIR/bin/tts-server -s user"
echo ""
