#!/bin/sh

set -o errexit -o nounset

HOME_DIR="/Users/rootulp/.celestia-app-mocha-3"
mkdir -p $HOME_DIR
echo "Home directory: ${HOME_DIR}"

CHAINID="mocha-3"

celestia-appd init "rootulp-1" --chain-id mocha-3 --home $HOME_DIR
cp ~/git/rootulp/celestia/networks/mocha-3/genesis.json $HOME_DIR/config
celestia-appd start --home ${HOME_DIR}

SEEDS="3314051954fc072a0678ec0cbac690ad8676ab98@65.108.66.220:26656"
PEERS="ec11f3be74010b78882de2cbd170d7ad4458d8ac@157.245.250.63:26656"

sed -i -e 's|^seeds *=.*|seeds = "'$SEEDS'"|; s|^persistent_peers *=.*|persistent_peers = "'$PEERS'"|' $HOME_DIR/config/config.toml
sed -i -e "s/^seed_mode *=.*/seed_mode = \"$SEED_MODE\"/" $HOME_DIR/config/config.toml
