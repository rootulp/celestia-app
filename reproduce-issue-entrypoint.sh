#!/bin/bash

# Entrypoint script to reproduce the genesis.json parsing issue
# Based on the GitHub issue: https://github.com/celestiaorg/celestia-app/issues/5793

set -e

echo "=== Reproducing celestia-app genesis.json parsing issue ==="
echo "Issue: https://github.com/celestiaorg/celestia-app/issues/5793"
echo "Docker image: ghcr.io/celestiaorg/celestia-app:v5.0.5-mocha"
echo ""

# Configuration
NODE_NAME="${NODE_NAME:-reproduce-issue-validator}"
CHAIN_ID="${CHAIN_ID:-mocha-4}"
CELESTIA_HOME="${CELESTIA_HOME:-/celestia}"

SEEDS="b402fe40f3474e9e208840702e1b7aa37f2edc4b@celestia-testnet-seed.itrocket.net:14656"

echo "Node name: ${NODE_NAME}"
echo "Chain ID: ${CHAIN_ID}"
echo "Celestia home: ${CELESTIA_HOME}"
echo "celestia-appd version: $(celestia-appd version 2>&1)"
echo ""

# Ensure the celestia home directory exists and has proper permissions
echo "Setting up celestia home directory..."
mkdir -p "$CELESTIA_HOME"
chmod 755 "$CELESTIA_HOME"

# Clean up any existing data
if [ "$(ls -A $CELESTIA_HOME 2>/dev/null)" ]; then
    echo "Cleaning up existing celestia home directory..."
    rm -rf "$CELESTIA_HOME"/*
fi

echo "Initializing config files..."
celestia-appd init ${NODE_NAME} --chain-id ${CHAIN_ID} --home ${CELESTIA_HOME}

echo "Setting seeds in config.toml..."
sed -i "s/^seeds *=.*/seeds = \"$SEEDS\"/" ${CELESTIA_HOME}/config/config.toml

echo "Downloading genesis file for ${CHAIN_ID}..."
celestia-appd download-genesis ${CHAIN_ID} --home ${CELESTIA_HOME}

echo "Genesis download completed. Checking genesis file..."
ls -la ${CELESTIA_HOME}/config/genesis.json

echo "Attempting to start celestia-appd (this should fail with the genesis parsing error)..."
celestia-appd start --home ${CELESTIA_HOME} --force-no-bbr
