celestia-appd keys add destination --keyring-backend test

# This assumes single-node.sh script is running
export FROM=$(celestia-appd keys show validator --keyring-backend test --address)
export TO=$(celestia-appd keys show destination --keyring-backend test --address)
export AMOUNT=1utia
export FEES=210utia

celestia-appd tx bank send $FROM $TO $AMOUNT --yes --fees $FEES
