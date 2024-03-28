#!/bin/sh

# Stop script execution if an error is encountered
set -o errexit
# Stop script execution if an undefined variable is used
set -o nounset

CHAIN_ID="celestia"
NODE_NAME="node-name"
SEEDS="e6116822e1a5e283d8a85d3ec38f4d232274eaf3@consensus-full-seed-1.celestia-bootstrap.net:26656,cf7ac8b19ff56a9d47c75551bd4864883d1e24b5@consensus-full-seed-2.celestia-bootstrap.net:26656"
CELESTIA_APP_HOME="${HOME}/.celestia-app"
CELESTIA_APP_VERSION=$(celestia-appd version 2>&1)
RPC_SERVERS="https://rpc.lunaroasis.net:26657,https://public-celestia-rpc.numia.xyz:26657"
RPC_RESPONSE=$(curl -s https://rpc.lunaroasis.net/status?)
TRUST_HEIGHT=$(echo $RPC_RESPONSE | jq -r '.result.sync_info.latest_block_height')
TRUST_HASH=$(echo $RPC_RESPONSE | jq -r '.result.sync_info.latest_block_hash')

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
sed -i.bak -e "s/^seeds *=.*/seeds = \"$SEEDS\"/" $CELESTIA_APP_HOME/config/config.toml
echo "Enabling state sync..."
sed -i.bak -e "s/^enable = false/enable = true/" $CELESTIA_APP_HOME/config/config.toml
echo "Setting RPC servers..."
sed -i.bak -e "s|^rpc_servers *=.*|rpc_servers = \"$RPC_SERVERS\"|" $CELESTIA_APP_HOME/config/config.toml
echo "Setting trust height to $TRUST_HEIGHT..."
sed -i.bak -e "s/^trust_height = 0/trust_height = $TRUST_HEIGHT/" $CELESTIA_APP_HOME/config/config.toml
echo "Setting trust hash to $TRUST_HASH..."
sed -i.bak -e "s/^trust_hash = \"\"/trust_hash = \"$TRUST_HASH\"/" $CELESTIA_APP_HOME/config/config.toml

echo "Downloading genesis file..."
celestia-appd download-genesis ${CHAIN_ID} > /dev/null 2>&1 # Hide output to reduce terminal noise

echo "Starting celestia-appd in the background and piping logs to mainnet.log"
nohup celestia-appd start > "${HOME}/mainnet.log" 2>&1 &

echo "You can check the node's status via: celestia-appd status"
