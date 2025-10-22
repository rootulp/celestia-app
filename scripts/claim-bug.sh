#!/bin/sh

set -o errexit # Stop script execution if an error is encountered
set -o nounset # Stop script execution if an undefined variable is used

APP_HOME="${HOME}/.celestia-app"
CHAIN_ID="test"
KEYRING_BACKEND="test"
FEES="500utia"

# Key names
VALIDATOR_KEY_NAME="validator"
VALIDATOR_2_KEY_NAME="validator2"
DELEGATOR_KEY_NAME="delegator"

echo "Adding delegator key to the keyring..."
celestia-appd keys add ${DELEGATOR_KEY_NAME} --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}"

VALIDATOR_ADDRESS=$(celestia-appd keys show ${VALIDATOR_KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")
VALIDATOR_2_ADDRESS=$(celestia-appd keys show ${VALIDATOR_2_KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")
DELEGATOR_ADDRESS=$(celestia-appd keys show ${DELEGATOR_KEY_NAME} -a --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")

VALIDATOR_VALOPER_ADDRESS=$(celestia-appd keys show ${VALIDATOR_KEY_NAME} -a --bech val --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")
VALIDATOR_2_VALOPER_ADDRESS=$(celestia-appd keys show ${VALIDATOR_2_KEY_NAME} -a --bech val --keyring-backend=${KEYRING_BACKEND} --home "${APP_HOME}")

echo "Validator address: $VALIDATOR_ADDRESS"
echo "Delegator address: $DELEGATOR_ADDRESS"

echo "Sending funds from validator to delegator..."
celestia-appd tx bank send $VALIDATOR_ADDRESS $DELEGATOR_ADDRESS 1000000000utia --fees 100000utia --yes

sleep 1

echo "Querying delegator balance..."
celestia-appd query bank balances $DELEGATOR_ADDRESS


echo "Delegating funds to validator..."
celestia-appd tx staking delegate $VALIDATOR_VALOPER_ADDRESS 100000000utia --from $DELEGATOR_ADDRESS --fees 100000utia --yes
sleep 1

echo "Querying delegation..."
celestia-appd query staking delegation $DELEGATOR_ADDRESS $VALIDATOR_VALOPER_ADDRESS

echo "Redelegating funds from validator to validator2..."
celestia-appd tx staking redelegate $VALIDATOR_VALOPER_ADDRESS $VALIDATOR_2_VALOPER_ADDRESS 100000000utia --from $DELEGATOR_ADDRESS --fees 100000utia  --yes
sleep 2

echo "Querying redelegation..."
celestia-appd query staking delegation $DELEGATOR_ADDRESS $VALIDATOR_2_VALOPER_ADDRESS
