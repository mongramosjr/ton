#!/bin/bash

# Deployment Script for TON Token Contract

# Set strict mode for better error handling
set -euo pipefail

# Project and Tool Paths
PROJECT_DIR=/home/mong/mong_projects/ton/token
FUNC_COMPILER=/home/mong/mong_projects/ton-dev/ton/build/crypto/func
TON_CLIENT=/home/mong/mong_projects/ton-dev/ton-client  # Update with actual path
TONOS_CLI=/home/mong/mong_projects/tonos-cli/bin/tonos-cli  # Update with actual path

# Token Contract Parameters
TOKEN_NAME="TanTan"
TOKEN_SYMBOL="TNT"
TOTAL_SUPPLY=200000000000000  # 200T with 9 decimals
DECIMALS=9

# Wallet and Network Configuration
OWNER_WALLET="0QC_Xqe8BRH4sr58r5z3AZuYSQRH-GhA0KnSYtUtWCYzKRTx"
NETWORK="testnet"  # Use "mainnet" for prod deployment

# Compilation Function (same as in compile.sh)
compile_contract() {
    echo "Compiling TON Smart Contract..."

    $FUNC_COMPILER \
        -o "$PROJECT_DIR/build/tantan_token.fif" \
        -SPA \
        "$PROJECT_DIR/contracts/tantan_token.fc"

    echo "Compilation completed."
}

# Deployment Preparation
prepare_deployment() {
    echo "Preparing deployment parameters..."

    # Generate deployment parameters
    $TONOS_CLI convert address \
        --addr "$OWNER_WALLET" \
        --url "$NETWORK"
}

# Deploy Smart Contract
deploy_contract() {
    echo "Deploying token contract..."

    $TONOS_CLI deploy \
        --abi "$PROJECT_DIR/build/tantan_token.abi.json" \
        --sign "$PROJECT_DIR/keys/tantan_token.genkey.json" \
        --data "{
            'name': '$TOKEN_NAME',
            'symbol': '$TOKEN_SYMBOL',
            'decimals': $DECIMALS,
            'total_supply': $TOTAL_SUPPLY
        }" \
        "$PROJECT_DIR/build/tantan_token.fif"
}

# Mint Initial Tokens
mint_initial_tokens() {
    echo "Minting initial tokens to owner's wallet..."

    $TONOS_CLI call \
        --abi "$PROJECT_DIR/build/tantan_token.abi.json" \
        --sign "$PROJECT_DIR/keys/tantan_token.genkey.json" \
        --method "mint" \
        --data "{
            'recipient': '$OWNER_WALLET',
            'amount': $TOTAL_SUPPLY
        }"
}

# Validation and Checks
validate_prerequisites() {
    # Check if required tools are installed
    command -v $FUNC_COMPILER >/dev/null 2>&1 || {
        echo "FunC compiler not found. Exiting."; exit 1;
    }

    command -v $TONOS_CLI >/dev/null 2>&1 || {
        echo "TONOS CLI not found. Exiting."; exit 1;
    }
}

# Error Handling
error_handling() {
    echo "Deployment failed at line $1"
    exit 1
}

# Main Deployment Process
main() {
    # Trap errors
    trap 'error_handling $LINENO' ERR

    # Run validation
    validate_prerequisites

    # Deployment steps
    compile_contract
    prepare_deployment
    deploy_contract
    mint_initial_tokens

    echo "Token contract deployment completed successfully!"
}

# Execute main deployment process
main