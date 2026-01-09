#!/bin/bash

# This script initializes celestia-app and configures it for state sync to Mocha testnet

set -o errexit
set -o nounset

CHAIN_ID="mocha-4"
NODE_NAME="monitoring-node"
SEEDS="b402fe40f3474e9e208840702e1b7aa37f2edc4b@celestia-testnet-seed.itrocket.net:14656,ee9f90974f85c59d3861fc7f7edb10894f6ac3c8@seed-mocha.pops.one:26656"
PEERS="daf2cecee2bd7f1b3bf94839f993f807c6b15fbf@celestia-testnet-peer.itrocket.net:11656,96b2761729cea90ee7c61206433fc0ba40c245bf@57.128.141.126:11656,f4f75a55bfc5f302ef34435ef096a4551ecb6804@152.53.33.96:12056,31bb1c9c1be7743d1115a8270bd1c83d01a9120a@148.72.141.31:26676,3e30bcfc55e7d351f18144aab4b0973e9e9bf987@65.108.226.183:11656,7a0d5818c0e5b0d4fbd86a9921f413f5e4e4ac1e@65.109.83.40:28656,43e9da043318a4ea0141259c17fcb06ecff816af@164.132.247.253:43656,5a7566aa030f7e5e7114dc9764f944b2b1324bcd@65.109.23.114:11656,c17c0cbf05e98656fee5f60fad469fc528f6d6de@65.109.25.113:11656,fb5e0b9efacc11916c58bbcd3606cbaa7d43c99f@65.108.234.84:28656,45504fb31eb97ea8778c920701fc8076e568a9cd@188.214.133.100:26656,edafdf47c443344fb940a32ab9d2067c482e59df@84.32.71.47:26656,ae7d00d6d70d9b9118c31ac0913e0808f2613a75@177.54.156.69:26656,7c841f59c35d70d9f1472d7d2a76a11eefb7f51f@136.243.69.100:43656"
RPC1="https://celestia-testnet-rpc.itrocket.net:443"
RPC2="https://public-celestia-mocha4-consensus.numia.xyz:443"

CELESTIA_APP_HOME="${CELESTIA_APP_HOME:-/home/celestia/.celestia-app}"
INIT_FLAG_FILE="${CELESTIA_APP_HOME}/config/.initialized"

echo "Initializing celestia-app for state sync to Mocha testnet..."
echo "Home directory: ${CELESTIA_APP_HOME}"
echo ""

# Check if already initialized
if [ -f "${INIT_FLAG_FILE}" ]; then
    echo "Already initialized, skipping..."
    exit 0
fi

# Initialize if not already done
if [ ! -f "${CELESTIA_APP_HOME}/config/config.toml" ]; then
    echo "Initializing config files..."
    echo "Running: celestia-appd init ${NODE_NAME} --chain-id ${CHAIN_ID}"
    if ! celestia-appd init ${NODE_NAME} --chain-id ${CHAIN_ID}; then
        echo "ERROR: Failed to initialize celestia-app"
        exit 1
    fi
    if [ ! -f "${CELESTIA_APP_HOME}/config/config.toml" ]; then
        echo "ERROR: Config file was not created after init"
        exit 1
    fi
    echo "Config files initialized successfully"
else
    echo "Config file already exists, skipping init"
fi

echo "Downloading genesis file..."
celestia-appd download-genesis ${CHAIN_ID} > /dev/null 2>&1 || echo "Warning: Genesis download failed, continuing..."

echo "Setting seeds in config.toml..."
sed -i.bak -e "s/^seeds *=.*/seeds = \"$SEEDS\"/" "${CELESTIA_APP_HOME}/config/config.toml"

echo "Setting persistent peers in config.toml..."
sed -i -e "/^\[p2p\]/,/^\[/{s/^[[:space:]]*persistent_peers *=.*/persistent_peers = \"$PEERS\"/;}" "${CELESTIA_APP_HOME}/config/config.toml"

echo "Fetching latest block height and trust information..."
# Try first RPC, fallback to second if it fails
LATEST_HEIGHT=$(curl -s --max-time 10 "${RPC1}/block" | jq -r .result.block.header.height 2>/dev/null || \
                curl -s --max-time 10 "${RPC2}/block" | jq -r .result.block.header.height 2>/dev/null)

if [ -z "$LATEST_HEIGHT" ] || [ "$LATEST_HEIGHT" == "null" ]; then
    echo "Warning: Could not fetch latest height, using default trust height"
    BLOCK_HEIGHT=1
    TRUST_HASH=""
else
    BLOCK_HEIGHT=$((LATEST_HEIGHT - 2000))
    if [ $BLOCK_HEIGHT -lt 1 ]; then
        BLOCK_HEIGHT=1
    fi

    echo "Latest height: $LATEST_HEIGHT"
    echo "Trust height: $BLOCK_HEIGHT"

    # Get trust hash
    TRUST_HASH=$(curl -s --max-time 10 "${RPC1}/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash 2>/dev/null || \
                 curl -s --max-time 10 "${RPC2}/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash 2>/dev/null)

    if [ -z "$TRUST_HASH" ] || [ "$TRUST_HASH" == "null" ]; then
        echo "Warning: Could not fetch trust hash"
        TRUST_HASH=""
    else
        echo "Trust hash: $TRUST_HASH"
    fi
fi

echo "Configuring state sync in config.toml..."
# Enable state sync
sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true|" "${CELESTIA_APP_HOME}/config/config.toml"

# Set RPC servers
sed -i -E "s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"${RPC1},${RPC2}\"|" "${CELESTIA_APP_HOME}/config/config.toml"

# Set trust height if we have it
if [ -n "$TRUST_HASH" ] && [ "$BLOCK_HEIGHT" -gt 0 ]; then
    sed -i -E "s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1${BLOCK_HEIGHT}|" "${CELESTIA_APP_HOME}/config/config.toml"
    sed -i -E "s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"${TRUST_HASH}\"|" "${CELESTIA_APP_HOME}/config/config.toml"
fi

# Set trust period (1 week)
sed -i -E "s|^(trust_period[[:space:]]+=[[:space:]]+).*$|\1\"168h0m0s\"|" "${CELESTIA_APP_HOME}/config/config.toml"

# Set discovery time
sed -i -E "s|^(discovery_time[[:space:]]+=[[:space:]]+).*$|\1\"5s\"|" "${CELESTIA_APP_HOME}/config/config.toml"

echo "Enabling Prometheus metrics in config.toml..."
# Enable Prometheus
sed -i -E "s|^(prometheus[[:space:]]+=[[:space:]]+).*$|\1true|" "${CELESTIA_APP_HOME}/config/config.toml"

# Set Prometheus listen address
sed -i -E "s|^(prometheus_listen_addr[[:space:]]+=[[:space:]]+).*$|\1\":26660\"|" "${CELESTIA_APP_HOME}/config/config.toml"

echo "Configuration complete!"
echo "State sync enabled: true"
echo "RPC servers: ${RPC1}, ${RPC2}"
if [ -n "$TRUST_HASH" ]; then
    echo "Trust height: ${BLOCK_HEIGHT}"
    echo "Trust hash: ${TRUST_HASH}"
fi
echo "Prometheus metrics enabled on :26660"

# Create flag file to indicate initialization is complete
touch "${INIT_FLAG_FILE}"
echo "Initialization flag created at ${INIT_FLAG_FILE}"
