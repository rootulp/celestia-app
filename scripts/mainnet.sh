#!/bin/sh

# Start a celestia-app consensus full node on mainnet using snap sync. Partially
# inspired by https://polkachu.com/state_sync/celestia

set -o errexit # Stop script execution if an error is encountered
set -o nounset # Stop script execution if an undefined variable is used

CHAIN_ID="celestia"
NODE_NAME="node-name"

CELESTIA_APP_HOME="${HOME}/.celestia-app"
CELESTIA_APP_VERSION=$(celestia-appd version 2>&1)

echo "celestia-app home: ${CELESTIA_APP_HOME}"
echo "celestia-app version: ${CELESTIA_APP_VERSION}"
echo ""

# Ask the user for confirmation before deleting the existing celestia-app home
# directory.
read -p "Are you sure you want to delete: $CELESTIA_APP_HOME? [y/n] " response

# Check the user's response
if [ "$response" != "y" ]; then
    # Exit if the user did not respond with "y"
    echo "You must delete $CELESTIA_APP_HOME to continue."
    exit 1
fi

echo "Deleting $CELESTIA_APP_HOME..."
rm -r "$CELESTIA_APP_HOME"

echo "Initializing config files..."
celestia-appd init ${NODE_NAME} --chain-id ${CHAIN_ID} > /dev/null 2>&1 # Hide output to reduce terminal noise

echo "Settings seeds in config.toml..."
SEEDS=$(curl -sL https://raw.githubusercontent.com/celestiaorg/networks/master/celestia/seeds.txt | tr '\n' ',')
sed -i.bak -e "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CELESTIA_APP_HOME/config/config.toml

echo "Setting peers in config.toml..."
PERSISTENT_PEERS=$(curl -sL https://raw.githubusercontent.com/celestiaorg/networks/master/celestia/peers.txt | tr '\n' ',')
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PERSISTENT_PEERS\"/" $CELESTIA_APP_HOME/config/config.toml

echo "Enabling state sync..."
SNAP_RPC="https://celestia-rpc.polkachu.com:443"
LATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \
BLOCK_HEIGHT=$((LATEST_HEIGHT - 2000)); \
TRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)
sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \
s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SNAP_RPC,$SNAP_RPC\"| ; \
s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \
s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"|" $CELESTIA_APP_HOME/config/config.toml

echo "Downloading genesis file..."
celestia-appd download-genesis ${CHAIN_ID} > /dev/null 2>&1 # Hide output to reduce terminal noise

echo "Starting celestia-appd in the background and piping logs to mainnet.log"
nohup celestia-appd start > "${HOME}/mainnet.log" 2>&1 &

echo "You can check the node's status via: celestia-appd status"
