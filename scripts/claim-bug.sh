#!/bin/sh

set -o errexit # Stop script execution if an error is encountered
set -o nounset # Stop script execution if an undefined variable is used

APP_HOME="${HOME}/.celestia-app"
CHAIN_ID="test"
KEYRING_BACKEND="test"
FEES="500utia"

# Key names
KEY_NAME="validator"
KEY_NAME_2="validator2"
DELEGATOR_KEY_NAME="delegator"

echo "Adding delegator key to the keyring..."
celestia-appd keys add ${DELEGATOR_KEY_NAME} --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}"

VALIDATOR_ADDRESS=$(celestia-appd keys show ${KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")
DELEGATOR_ADDRESS=$(celestia-appd keys show ${DELEGATOR_KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")

echo "Validator address: $VALIDATOR_ADDRESS"
echo "Delegator address: $DELEGATOR_ADDRESS"

echo "Sending funds from validator to delegator..."
celestia-appd tx bank send $VALIDATOR_ADDRESS $DELEGATOR_ADDRESS 1000000000utia --fees 100000utia --yes

sleep 1

echo "Querying delegator balance..."
celestia-appd query bank balances $DELEGATOR_ADDRESS
