#!/bin/bash

# Densendither Installation Script
# This script builds the densendither binary and installs it to your system PATH

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default installation directory
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="densendither"

echo -e "${BLUE}Densendither Installation Script${NC}"
echo "=================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo -e "${GREEN}✓ Go found: $(go version)${NC}"

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
    echo -e "${RED}Error: Must run this script from the densendither project root directory${NC}"
    echo "Expected files: go.mod, main.go"
    exit 1
fi

echo -e "${GREEN}✓ Project files found${NC}"

# Clean any previous builds
echo -e "${YELLOW}Cleaning previous builds...${NC}"
go clean

# Download dependencies
echo -e "${YELLOW}Downloading dependencies...${NC}"
go mod tidy
go mod download

# Build the binary
echo -e "${YELLOW}Building densendither...${NC}"
CGO_ENABLED=0 go build -ldflags="-w -s" -o "${BINARY_NAME}" .

if [ ! -f "${BINARY_NAME}" ]; then
    echo -e "${RED}Error: Build failed - binary not created${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Build successful${NC}"

# Test the binary
echo -e "${YELLOW}Testing binary...${NC}"
if ./"${BINARY_NAME}" --help > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Binary test passed${NC}"
else
    echo -e "${RED}Error: Binary test failed${NC}"
    exit 1
fi

# Check if running with sudo for system installation
if [ "$EUID" -eq 0 ]; then
    # Running as root - install to system directory
    echo -e "${YELLOW}Installing to system directory: ${INSTALL_DIR}${NC}"

    # Create install directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"

    # Copy binary to install directory
    cp "${BINARY_NAME}" "${INSTALL_DIR}/"
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    echo -e "${GREEN}✓ Installed ${BINARY_NAME} to ${INSTALL_DIR}${NC}"

else
    # Not running as root - offer alternatives
    echo -e "${YELLOW}Not running as root. Choose installation option:${NC}"
    echo "1. Install to ~/.local/bin (user directory - recommended)"
    echo "2. Install to /usr/local/bin (system-wide - requires sudo)"
    echo "3. Keep binary in current directory"

    read -p "Enter choice (1-3): " choice

    case $choice in
        1)
            USER_BIN_DIR="$HOME/.local/bin"
            echo -e "${YELLOW}Installing to user directory: ${USER_BIN_DIR}${NC}"

            # Create user bin directory if it doesn't exist
            mkdir -p "${USER_BIN_DIR}"

            # Copy binary to user bin directory
            cp "${BINARY_NAME}" "${USER_BIN_DIR}/"
            chmod +x "${USER_BIN_DIR}/${BINARY_NAME}"

            echo -e "${GREEN}✓ Installed ${BINARY_NAME} to ${USER_BIN_DIR}${NC}"

            # Check if ~/.local/bin is in PATH
            if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
                echo -e "${YELLOW}Warning: ${USER_BIN_DIR} is not in your PATH${NC}"
                echo "Add the following line to your ~/.bashrc or ~/.zshrc:"
                echo -e "${BLUE}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
                echo "Then run: source ~/.bashrc (or ~/.zshrc)"
            fi
            ;;
        2)
            echo -e "${YELLOW}Re-running with sudo...${NC}"
            sudo cp "${BINARY_NAME}" "${INSTALL_DIR}/"
            sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
            echo -e "${GREEN}✓ Installed ${BINARY_NAME} to ${INSTALL_DIR}${NC}"
            ;;
        3)
            echo -e "${YELLOW}Binary remains in current directory: $(pwd)/${BINARY_NAME}${NC}"
            echo "You can run it with: ./${BINARY_NAME}"
            ;;
        *)
            echo -e "${RED}Invalid choice. Binary remains in current directory.${NC}"
            ;;
    esac
fi

# Final verification
echo -e "${YELLOW}Verifying installation...${NC}"

if command -v "${BINARY_NAME}" &> /dev/null; then
    echo -e "${GREEN}✓ ${BINARY_NAME} is now available in PATH${NC}"
    echo -e "${GREEN}✓ Installation completed successfully!${NC}"
    echo
    echo -e "${BLUE}Quick start:${NC}"
    echo "  ${BINARY_NAME} --help                    # Show help"
    echo "  ${BINARY_NAME} palette list              # List available palettes"
    echo "  ${BINARY_NAME} palette add -n test -c \"#000000,#FFFFFF\"  # Add a palette"
    echo "  ${BINARY_NAME} dither -i image.png -p test               # Dither an image"
else
    echo -e "${YELLOW}Installation completed, but ${BINARY_NAME} is not in PATH${NC}"
    echo "You may need to:"
    echo "  1. Restart your terminal"
    echo "  2. Run: source ~/.bashrc (or ~/.zshrc)"
    echo "  3. Add the install directory to your PATH"
fi

# Clean up build artifact from current directory
if [ -f "${BINARY_NAME}" ] && command -v "${BINARY_NAME}" &> /dev/null; then
    rm "${BINARY_NAME}"
    echo -e "${GREEN}✓ Cleaned up build artifacts${NC}"
fi

echo
echo -e "${GREEN}🎉 Densendither installation complete!${NC}"
