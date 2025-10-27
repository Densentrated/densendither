#!/bin/bash

# Densendither Uninstall Script
# This script removes densendither from your system

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BINARY_NAME="densendither"
CONFIG_DIR="$HOME/.config/densendither"

echo -e "${BLUE}Densendither Uninstall Script${NC}"
echo "============================="

# Function to remove binary from a directory
remove_binary() {
    local dir="$1"
    local binary_path="${dir}/${BINARY_NAME}"

    if [ -f "$binary_path" ]; then
        if [ -w "$dir" ]; then
            rm "$binary_path"
            echo -e "${GREEN}✓ Removed ${binary_path}${NC}"
            return 0
        else
            echo -e "${YELLOW}Need sudo to remove ${binary_path}${NC}"
            sudo rm "$binary_path"
            echo -e "${GREEN}✓ Removed ${binary_path}${NC}"
            return 0
        fi
    fi
    return 1
}

# Look for the binary in common locations
FOUND_BINARY=false
BINARY_LOCATIONS=(
    "/usr/local/bin"
    "/usr/bin"
    "$HOME/.local/bin"
    "$HOME/bin"
)

echo -e "${YELLOW}Searching for ${BINARY_NAME} binary...${NC}"

for location in "${BINARY_LOCATIONS[@]}"; do
    if remove_binary "$location"; then
        FOUND_BINARY=true
    fi
done

# Also check if it exists in current directory
if [ -f "./${BINARY_NAME}" ]; then
    rm "./${BINARY_NAME}"
    echo -e "${GREEN}✓ Removed ${BINARY_NAME} from current directory${NC}"
    FOUND_BINARY=true
fi

if [ "$FOUND_BINARY" = false ]; then
    echo -e "${YELLOW}No ${BINARY_NAME} binary found in common locations${NC}"
fi

# Ask about configuration files
if [ -d "$CONFIG_DIR" ]; then
    echo
    echo -e "${YELLOW}Configuration directory found: ${CONFIG_DIR}${NC}"
    echo "This contains your color palettes and settings."
    echo
    read -p "Do you want to remove configuration files? (y/N): " remove_config

    case "$remove_config" in
        [Yy]|[Yy][Ee][Ss])
            rm -rf "$CONFIG_DIR"
            echo -e "${GREEN}✓ Removed configuration directory${NC}"
            ;;
        *)
            echo -e "${BLUE}Configuration files preserved at ${CONFIG_DIR}${NC}"
            ;;
    esac
else
    echo -e "${BLUE}No configuration directory found${NC}"
fi

# Verify removal
echo
echo -e "${YELLOW}Verifying removal...${NC}"

if command -v "$BINARY_NAME" &> /dev/null; then
    echo -e "${YELLOW}Warning: ${BINARY_NAME} is still available in PATH${NC}"
    echo "You may need to:"
    echo "  1. Restart your terminal"
    echo "  2. Check for additional installations"
    echo "  3. Manually remove from custom locations"
else
    echo -e "${GREEN}✓ ${BINARY_NAME} successfully removed from PATH${NC}"
fi

echo
if [ "$FOUND_BINARY" = true ]; then
    echo -e "${GREEN}🗑️  Densendither uninstallation complete!${NC}"
else
    echo -e "${YELLOW}⚠️  No binaries were removed (may not have been installed)${NC}"
fi

echo
echo "Thank you for using Densendither!"
