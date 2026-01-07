#!/bin/bash

# This script initializes celestia-app for a single-node testnet
# Similar to scripts/single-node.sh but adapted for Docker

set -o errexit
set -o nounset

# Constants
CHAIN_ID="test"
KEY_NAME="validator"
KEYRING_BACKEND="test"
FEES="500utia"

APP_HOME="${CELESTIA_APP_HOME:-/home/celestia/.celestia-app}"
GENESIS_FILE="${APP_HOME}/config/genesis.json"

echo "celestia-app home: ${APP_HOME}"
echo "celestia-app genesis file: ${GENESIS_FILE}"
echo ""

createGenesis() {
    echo "Initializing validator and node config files..."
    /bin/celestia-appd init ${CHAIN_ID} \
      --chain-id ${CHAIN_ID} \
      --home "${APP_HOME}" \
      > /dev/null 2>&1

    echo "Adding a new key to the keyring..."
    /bin/celestia-appd keys add ${KEY_NAME} \
      --keyring-backend=${KEYRING_BACKEND} \
      --home "${APP_HOME}" \
      > /dev/null 2>&1

    echo "Adding genesis account..."
    /bin/celestia-appd genesis add-genesis-account \
      "$(/bin/celestia-appd keys show ${KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")" \
      "1000000000000000utia" \
      --home "${APP_HOME}"

    echo "Creating a genesis tx..."
    /bin/celestia-appd genesis gentx ${KEY_NAME} 5000000000utia \
      --fees ${FEES} \
      --keyring-backend=${KEYRING_BACKEND} \
      --chain-id ${CHAIN_ID} \
      --home "${APP_HOME}" \
      --commission-rate=0.05 \
      --commission-max-rate=1.0 \
      --commission-max-change-rate=1.0 \
      > /dev/null 2>&1

    echo "Collecting genesis txs..."
    /bin/celestia-appd genesis collect-gentxs \
      --home "${APP_HOME}" \
        > /dev/null 2>&1

    # Override the default RPC server listening address
    sed -i.bak 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' "${APP_HOME}"/config/config.toml

    # Override the default API server listening address to allow access from other containers
    sed -i.bak 's#address = "tcp://localhost:1317"#address = "tcp://0.0.0.0:1317"#g' "${APP_HOME}"/config/app.toml

    # Enable transaction indexing
    sed -i.bak 's#"null"#"kv"#g' "${APP_HOME}"/config/config.toml

    # Persist ABCI responses
    sed -i.bak 's#discard_abci_responses = true#discard_abci_responses = false#g' "${APP_HOME}"/config/config.toml

    # Override the log level to reduce noisy logs
    sed -i.bak 's#log_level = "info"#log_level = "*:error,p2p:info,state:info"#g' "${APP_HOME}"/config/config.toml

    # Override the VotingPeriod from 1 week to 30 seconds
    sed -i.bak 's#"604800s"#"30s"#g' "${APP_HOME}"/config/genesis.json

    # Fix genesis structure for indexer compatibility
    # The indexer expects deposit_params, voting_params, and tally_params to be populated
    # Copy only the relevant fields from params to these structures
    if command -v jq >/dev/null 2>&1; then
        jq '.app_state.gov.deposit_params = {
                min_deposit: .app_state.gov.params.min_deposit,
                max_deposit_period: .app_state.gov.params.max_deposit_period
            } |
            .app_state.gov.voting_params = {
                voting_period: .app_state.gov.params.voting_period
            } |
            .app_state.gov.tally_params = {
                quorum: .app_state.gov.params.quorum,
                threshold: .app_state.gov.params.threshold,
                veto_threshold: .app_state.gov.params.veto_threshold
            }' \
            "${APP_HOME}"/config/genesis.json > "${APP_HOME}"/config/genesis.json.tmp && \
            mv "${APP_HOME}"/config/genesis.json.tmp "${APP_HOME}"/config/genesis.json
    fi

    trace_type="local"
    sed -i.bak -e "s/^trace_type *=.*/trace_type = \"$trace_type\"/" ${APP_HOME}/config/config.toml

    trace_pull_address=":26661"
    sed -i.bak -e "s/^trace_pull_address *=.*/trace_pull_address = \"$trace_pull_address\"/" ${APP_HOME}/config/config.toml

    trace_push_batch_size=1000
    sed -i.bak -e "s/^trace_push_batch_size *=.*/trace_push_batch_size = \"$trace_push_batch_size\"/" ${APP_HOME}/config/config.toml

    echo "Genesis created successfully!"
}

# Check if genesis file exists
if [ ! -f "$GENESIS_FILE" ]; then
    echo "Genesis file not found. Creating new genesis..."
    createGenesis
else
    echo "Genesis file already exists. Skipping initialization."
fi
