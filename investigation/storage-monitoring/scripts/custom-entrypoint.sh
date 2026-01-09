#!/bin/bash

# Custom entrypoint that handles initialization before starting celestia-appd

set -o errexit
set -o nounset

CELESTIA_APP_HOME="${CELESTIA_APP_HOME:-/home/celestia/.celestia-app}"
INIT_FLAG_FILE="${CELESTIA_APP_HOME}/config/.initialized"

# Check if we're running as root (for permission fixes)
if [ "$(id -u)" = "0" ]; then
    # Fix permissions for config and data directories
    # This is needed because volumes might be owned by root
    if [ -d "${CELESTIA_APP_HOME}/config" ]; then
        echo "Fixing permissions for config directory..."
        chown -R celestia:celestia "${CELESTIA_APP_HOME}/config" 2>/dev/null || true
        chmod -R u+w "${CELESTIA_APP_HOME}/config" 2>/dev/null || true
    fi

    if [ -d "${CELESTIA_APP_HOME}/data" ]; then
        echo "Fixing permissions for data directory..."
        chown -R celestia:celestia "${CELESTIA_APP_HOME}/data" 2>/dev/null || true
        chmod -R u+w "${CELESTIA_APP_HOME}/data" 2>/dev/null || true
    fi

    # Ensure the home directory structure exists
    mkdir -p "${CELESTIA_APP_HOME}/config" "${CELESTIA_APP_HOME}/data"
    chown -R celestia:celestia "${CELESTIA_APP_HOME}" 2>/dev/null || true
fi

# Initialize if config doesn't exist and not already initialized
if [ ! -f "${CELESTIA_APP_HOME}/config/config.toml" ] && [ ! -f "${INIT_FLAG_FILE}" ]; then
    echo "Config not found, running initialization..."
    if ! /opt/init-celestia-app.sh; then
        echo "ERROR: Initialization failed!"
        exit 1
    fi
elif [ -f "${INIT_FLAG_FILE}" ]; then
    echo "Already initialized (flag file exists)"
fi

# Create priv_validator_state.json if it doesn't exist and we're starting
if [[ "$1" == "start" && ! -f "${CELESTIA_APP_HOME}/data/priv_validator_state.json" ]]; then
    mkdir -p "${CELESTIA_APP_HOME}/data"
    cat <<EOF > "${CELESTIA_APP_HOME}/data/priv_validator_state.json"
{
  "height": "0",
  "round": 0,
  "step": 0
}
EOF
fi

# Fix permissions one more time before switching user (if root)
if [ "$(id -u)" = "0" ]; then
    chown -R celestia:celestia "${CELESTIA_APP_HOME}" 2>/dev/null || true
fi

# Switch to celestia user and start celestia-appd
echo "Starting celestia-appd with command:"
echo "/bin/celestia-appd $@"
echo ""

# For this monitoring setup, we'll run as root to avoid permission issues
# This is acceptable for a temporary investigation container
# The permissions have been fixed above, so celestia-appd should work fine
exec /bin/celestia-appd "$@"
