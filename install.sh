#!/bin/bash

# MCP-Think Installer Script
# This script detects your OS and architecture and installs the latest release of MCP-Think.

set -e

# Color setup
COLOR_GREEN="\033[0;32m"
COLOR_RED="\033[0;31m"
COLOR_YELLOW="\033[0;33m"
COLOR_BLUE="\033[0;34m"
COLOR_RESET="\033[0m"

# Functions
success() { echo -e "${COLOR_GREEN}✓ $1${COLOR_RESET}"; }
warn() { echo -e "${COLOR_YELLOW}! $1${COLOR_RESET}"; }
error() { echo -e "${COLOR_RED}✗ $1${COLOR_RESET}" >&2; exit 1; }
info() { echo -e "${COLOR_BLUE}i $1${COLOR_RESET}"; }

# Display header
info "MCP-Think Installer"
info "==================="

# Auto-detect OS and Architecture
OS_KERNEL=$(uname -s)
OS_ARCH=$(uname -m)

info "Detecting system configuration..."

case $OS_KERNEL in
    Linux)
        OS='linux'
        success "Detected OS: Linux"
        ;;
    Darwin)
        OS='darwin'
        success "Detected OS: macOS"
        ;;
    MINGW* | MSYS* | CYGWIN*)
        OS='windows'
        success "Detected OS: Windows"
        ;;
    *)
        error "Unsupported operating system: $OS_KERNEL"
        ;;
esac

case $OS_ARCH in
    x86_64)
        ARCH='amd64'
        success "Detected Architecture: x86_64 (using amd64)"
        ;;
    arm64 | aarch64)
        ARCH='arm64'
        success "Detected Architecture: ARM64"
        ;;
    *)
        error "Unsupported architecture: $OS_ARCH"
        ;;
esac

# Fetch the latest version from GitHub
info "Fetching the latest release version..."
LATEST_VERSION=$(curl -s https://api.github.com/repos/iamwavecut/MCP-Think/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    LATEST_VERSION="v0.0.2"  # Fallback version if GitHub API fails
    warn "Could not determine latest version, using $LATEST_VERSION as fallback"
else
    success "Found latest version: $LATEST_VERSION"
fi

# Add .exe extension for Windows
EXE_SUFFIX=""
if [ "$OS" = "windows" ]; then
    EXE_SUFFIX=".exe"
fi

BINARY_NAME="think-tool-${LATEST_VERSION}-${OS}-${ARCH}${EXE_SUFFIX}"
RELEASE_URL="https://github.com/iamwavecut/MCP-Think/releases/latest/download/${BINARY_NAME}"
TMP_DOWNLOAD="/tmp/think-tool"

info "Preparing to download MCP-Think ($BINARY_NAME)..."

# Download the binary
info "Downloading from: $RELEASE_URL"
if curl -L --progress-bar "$RELEASE_URL" -o "$TMP_DOWNLOAD"; then
    success "Download completed successfully"
else
    error "Download failed. Please check your internet connection or the URL."
fi

# Make it executable
chmod +x "$TMP_DOWNLOAD"
success "Set executable permissions"

# Determine available installation paths
SYSTEM_BIN="/usr/local/bin"
LOCAL_USER_BIN=""

# Check if ~/.local/bin exists and is in PATH
if [[ -d "$HOME/.local/bin" && "$PATH" == *"$HOME/.local/bin"* ]]; then
    LOCAL_USER_BIN="$HOME/.local/bin"
elif [[ -d "$HOME/.bin" && "$PATH" == *"$HOME/.bin"* ]]; then
    LOCAL_USER_BIN="$HOME/.bin"
fi

# Ask for installation location
echo ""
info "Installation options:"
echo "  1. System-wide ($SYSTEM_BIN) [requires sudo]"

DEFAULT_OPTION=2

if [ -n "$LOCAL_USER_BIN" ]; then
    echo "  2. User-specific ($LOCAL_USER_BIN) [default]"
    echo "  3. Current directory (./think-tool)"
    echo "  4. Cancel installation"
else
    echo "  2. Current directory (./think-tool) [default]"
    echo "  3. Cancel installation"
fi

if [ -t 0 ]; then
    # Terminal is interactive
    read -r -p "Enter your choice [default: $DEFAULT_OPTION]: " CHOICE
else
    # Non-interactive mode (e.g., piped)
    echo "Non-interactive mode detected, using default option: $DEFAULT_OPTION"
    CHOICE=""
fi

if [ -z "$CHOICE" ]; then
    CHOICE=$DEFAULT_OPTION
fi

# Process the user's choice
if [ "$CHOICE" = "1" ]; then
    # System-wide installation
    INSTALL_PATH="$SYSTEM_BIN/think-tool"
    info "Installing to $INSTALL_PATH..."
    sudo mv "$TMP_DOWNLOAD" "$INSTALL_PATH"
    EXEC_CMD="think-tool"
    success "Installation completed!"
elif [ "$CHOICE" = "2" ]; then
    # User-specific or current directory (default)
    if [ -n "$LOCAL_USER_BIN" ]; then
        INSTALL_PATH="$LOCAL_USER_BIN/think-tool"
        EXEC_CMD="think-tool"
    else
        INSTALL_PATH="./think-tool"
        EXEC_CMD="./think-tool"
    fi
    info "Installing to $INSTALL_PATH..."
    mv "$TMP_DOWNLOAD" "$INSTALL_PATH"
    success "Installation completed!"
elif [ "$CHOICE" = "3" ]; then
    # Current directory or cancel
    if [ -n "$LOCAL_USER_BIN" ]; then
        INSTALL_PATH="./think-tool"
        EXEC_CMD="./think-tool"
        info "Installing to current directory..."
        mv "$TMP_DOWNLOAD" "$INSTALL_PATH"
        success "Installation completed!"
    else
        warn "Installation cancelled"
        rm -f "$TMP_DOWNLOAD"
        exit 0
    fi
elif [ "$CHOICE" = "4" ]; then
    # Cancel installation
    if [ -n "$LOCAL_USER_BIN" ]; then
        warn "Installation cancelled"
        rm -f "$TMP_DOWNLOAD"
        exit 0
    else
        error "Invalid option: $CHOICE"
    fi
else
    error "Invalid option: $CHOICE"
fi

echo ""
info "To run MCP-Think:"
echo "  $ $EXEC_CMD"
echo ""
info "Thank you for installing MCP-Think!" 