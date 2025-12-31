#!/bin/env bash

set -e # Exit on any error

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
fi

TARGET_FILE="$TARGET_DIR/run"

if [ -f "$TARGET_FILE" ]; then
    echo "Uninstalling from $SCOPE ($TARGET_DIR)"
    if $USE_SUDO; then
        sudo rm -f "$TARGET_FILE"
    else
        rm -f "$TARGET_FILE"
    fi
    echo "Uninstalled!"
else
    echo "ERROR: $TARGET_FILE not found"
    echo "It may have already been uninstalled"
fi
