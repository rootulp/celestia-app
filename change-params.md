# Change params

This document describes how to change the auth params of a celestia-appd network.

## Steps

1. Start a single node devnet

    ```bash
    ./scripts/single-node.sh
    ```

1. Verify the auth params before making any changes

    ```bash
    celestia-appd query auth params
    ```

    ```txt
    max_memo_characters: "256"
    sig_verify_cost_ed25519: "590"
    sig_verify_cost_secp256k1: "1000"
    tx_sig_limit: "7"
    tx_size_cost_per_byte: "10"
    ```

1. Submit gov proposal

    ```bash
    celestia-appd tx gov submit-legacy-proposal param-change proposal.json --from validator --fees 10000utia --yes
    ```

1. Vote on the proposal

    ```bash
    celestia-appd tx gov vote 1 yes --from validator --fees 10000utia --yes
    ```

1. Verify the gov proposal was submitted

    ```bash
    celestia-appd query gov proposals
    ```

1. Wait 1 minute for voting period to elapse
1. Observe logs

    ```log
    11:01PM INF attempt to set new parameter value; key: MaxMemoCharacters, value: "16" module=x/params
    11:01PM INF attempt to set new parameter value; key: TxSizeCostPerByte, value: "16" module=x/params

    11:01PM INF proposal tallied module=x/gov proposal=1 results=passed
    ```

1. Query the auth params. Observe they have been updated.

    ```bash
    celestia-appd query auth params
    ```

    ```txt
    max_memo_characters: "16"
    sig_verify_cost_ed25519: "590"
    sig_verify_cost_secp256k1: "1000"
    tx_sig_limit: "7"
    tx_size_cost_per_byte: "16"
    ```
