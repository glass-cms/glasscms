#!/usr/bin/env zsh

# Exit when a command fails.
set -euo pipefail

BREW_PACKAGES=(go-task jq golangci-lint git-cliff sqlc)
NODE_VERSION="20.17.0"

BLACK="\033[30m"
RED="\033[31m"
GREEN="\033[32m"
YELLOW="\033[33m"
BLUE="\033[34m"
MAGENTA="\033[35m"
CYAN="\033[36m"
WHITE="\033[37m"
DEFAULT="\033[39m"
RESET="\033[0m"
BOLD="\033[1m"
UNDERLINE="\033[4m"
REVERSED="\033[7m"

pretty_print() {
    local message="$1"
    local color="${2:-$DEFAULT}"
    local style="${3:-}"

    echo -e "${color}${style}$message${RESET}"
}

header() {
    pretty_print "\n========================================================================"
    pretty_print "$1" $BLUE $BOLD
    pretty_print "========================================================================\n"
}

success() {
    pretty_print ""
    pretty_print "✅ $1" $GREEN
}

error() {
    pretty_print "❌ $1" $RED
    exit 1
}

# Check prerequisites
if ! command -v go >/dev/null; then
    error "Go is not installed. Please install Go first."
fi

if ! command -v brew >/dev/null; then
    error "Homebrew is not installed. Please install Homebrew first."
fi

# Install or upgrade a package with Homebrew
brew_install() {
    if brew ls --versions "$1"; then
        brew upgrade "$1"
    else
        brew install "$1"
    fi
}

# Install packages with Homebrew
install_packages() {
    header "Installing packages with Homebrew 🍺"

    for package in "${BREW_PACKAGES[@]}"; do
        echo "Installing or upgrading $package..."
        brew_install "$package"
        echo "Installed or upgraded $package"
    done

    success "Finished installing packages with Homebrew!"
}

echo "Installing Homebrew packages..."
install_packages
pretty_print "\nFinished installing Homebrew packages! 🎉" $GREEN $BOLD