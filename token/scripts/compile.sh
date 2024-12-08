#!/bin/bash

# Compilation Script for TON Token Contract

# Set strict mode for better error handling
set -euo pipefail

# Project Directories
PROJECT_DIR=/home/mong/mong_projects/ton/token
FUNC_COMPILER=/home/mong/mong_projects/ton-dev/ton/build/crypto/func

# Source and Output Paths
SOURCE_CONTRACT=$PROJECT_DIR/contracts/tantan_token.fc
OUTPUT_FIF=$PROJECT_DIR/build/tantan_token.fif
OUTPUT_BOC=$PROJECT_DIR/build/tantan_token.boc

# Compilation Function
compile_contract() {
    echo "Compiling TON Smart Contract..."

    # Compile FunC to Fift Assembler
    $FUNC_COMPILER \
        -o "$OUTPUT_FIF" \
        -SPA \
        "$SOURCE_CONTRACT"

    echo "Compilation completed. Output: $OUTPUT_FIF"
}

# Validate Compiler
validate_compiler() {
    if [ ! -x "$FUNC_COMPILER" ]; then
        echo "Error: FunC compiler not found or not executable"
        exit 1
    fi
}

# Main Compilation Process
main() {
    validate_compiler
    compile_contract
}

# Run the main compilation process
main