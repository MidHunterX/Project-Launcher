#!/bin/env bash

set -e # Exit on any error

SOURCE_FILE="./run"
if [ ! -f "$SOURCE_FILE" ]; then
    echo "Error: $SOURCE_FILE not found in current directory"
    exit 1
fi

if [[ ":$PATH:" == *":$HOME/.local/bin:"* ]]; then
    # LOCAL BIN
    TARGET_DIR="$HOME/.local/bin"
    SCOPE="local bin"
    USE_SUDO=false
elif [[ ":$PATH:" == *":/usr/local/bin:"* ]]; then
    # USER BIN
    TARGET_DIR="/usr/local/bin"
    SCOPE="user bin"
    USE_SUDO=true
else
    echo "ERROR: No bin found"
    exit
fi

TARGET_FILE="$TARGET_DIR/run"

[ -f "$TARGET_FILE" ] && UPDATE=true || UPDATE=false

if $UPDATE; then
    echo "Updating existing installation..."
else
    echo "Installing to $SCOPE ($TARGET_DIR)"
fi

# Create directory if it doesn't exist
if $USE_SUDO; then
    sudo mkdir -p "$TARGET_DIR"
    sudo cp "$SOURCE_FILE" "$TARGET_FILE"
    sudo chmod +x "$TARGET_FILE"
else
    mkdir -p "$TARGET_DIR"
    cp "$SOURCE_FILE" "$TARGET_FILE"
    chmod +x "$TARGET_FILE"
fi

if $UPDATE; then
    echo "Updated successfully!"
else
    echo "Installed successfully!"
fi
