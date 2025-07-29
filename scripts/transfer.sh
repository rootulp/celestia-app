#!/bin/sh

# This script is used to test IBC between Celestia and IBC-0 (a Gaia chain managed via GM).
# Run celestia-app via ./scripts/single-node.sh or ./scripts/single-node-upgrades.sh
# Then run ./scripts/hermes.sh to set up Hermes.
# Then run this script to transfer tokens from Celestia to IBC-0.

set -o errexit # Stop script execution if an error is encountered
set -o nounset # Stop script execution if an undefined variable is used

echo "--> Transferring tokens from Celestia to IBC-0"
hermes tx ft-transfer --timeout-seconds 1000 --dst-chain ibc-0 --src-chain test --src-port transfer --src-channel channel-0 --amount 100000 --denom utia

echo "--> Waiting for transfer to complete"
sleep 10

echo "--> Querying balance of IBC-0"
gaiad --node tcp://localhost:27030 query bank balances $(gaiad --home ~/.gm/ibc-0 keys --keyring-backend="test" show wallet -a)
