#!/bin/sh
set -e

# Configuration
REPO="glass-cms/glasscms"
BINARY_NAME="glasscms"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
USE_SUDO="${USE_SUDO:-true}"

# Colors for output (only if terminal supports it)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    NC='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    NC=''
fi

# Helper functions
log_info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1" >&2
}

log_warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1" >&2
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin*) OS="darwin" ;;
        Linux*)  OS="linux" ;;
        *)       log_error "Unsupported OS: $(uname -s)" && exit 1 ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) ARCH="amd64" ;;
        arm64|aarch64) ARCH="arm64" ;;
        *) log_error "Unsupported architecture: $(uname -m)" && exit 1 ;;
    esac
}

# Get latest version from GitHub
get_latest_version() {
    log_info "Getting latest version..."
    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | cut -d'"' -f4)
    if [ -z "$VERSION" ]; then
        log_error "Failed to get latest version"
        exit 1
    fi
    log_info "Latest version: $VERSION"
}

# Download binary
download_binary() {
    binary_name="${BINARY_NAME}-${OS}-${ARCH}"
    download_url="https://github.com/$REPO/releases/download/$VERSION/$binary_name"
    tmp_file="/tmp/$binary_name"
    
    log_info "Downloading $binary_name..."
    if ! curl -sL "$download_url" -o "$tmp_file"; then
        log_error "Failed to download binary from $download_url"
        exit 1
    fi
    
    chmod +x "$tmp_file"
    printf "%s" "$tmp_file"
}

# Install binary
install_binary() {
    tmp_file="$1"
    install_path="$INSTALL_DIR/$BINARY_NAME"
    
    log_info "Installing to $install_path..."
    
    # Create install directory if it doesn't exist
    if [ ! -d "$INSTALL_DIR" ]; then
        if [ "$USE_SUDO" = "true" ]; then
            sudo mkdir -p "$INSTALL_DIR"
        else
            mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    # Install the binary
    if [ "$USE_SUDO" = "true" ] && [ ! -w "$INSTALL_DIR" ]; then
        sudo mv "$tmp_file" "$install_path"
    else
        mv "$tmp_file" "$install_path"
    fi
    
    log_info "Installation complete!"
}

# Verify installation
verify_installation() {
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        log_info "Verification successful:"
        "$BINARY_NAME" --version
    else
        log_warn "Binary installed but not in PATH. You may need to add $INSTALL_DIR to your PATH."
    fi
}

# Main installation flow
main() {
    log_info "Installing $BINARY_NAME..."
    
    detect_os
    detect_arch
    get_latest_version
    
    tmp_file=$(download_binary)
    install_binary "$tmp_file"
    verify_installation
    
    log_info "🎉 $BINARY_NAME has been successfully installed!"
}

# Always run main when script is executed directly
main "$@"