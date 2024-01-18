# Change app version via gov proposal

## Steps

1. Start a single node devnet

    ```bash
    ./scripts/single-node.sh
    ```

1. Submit gov proposal

    ```bash
    celestia-appd tx gov submit-legacy-proposal param-change proposal.json --from validator --fees 10000utia --yes
    ```

1. Vote on the proposal

    ```bash
    # Vote on the proposal
    celestia-appd tx gov vote 1 yes --from validator --fees 10000utia --yes
    ```

1. Verify the gov proposal was submitted

    ```bash
    # Verify gov proposal was submitted
    celestia-appd query gov proposals
    ```

1. Wait 2 minutes for voting period to elapse
1. Observe logs

    ```log
    3:46PM INF attempt to set new parameter value; key: VersionParams, value: {"app_version": "3"} module=x/params
    3:46PM INF proposal tallied module=x/gov proposal=1 results=passed
    ```

1. Stop the node and export the state. Observe

    ```json
        "consensus_params": {
            "block": {
                "max_bytes": "1974272",
                "max_gas": "-1",
                "time_iota_ms": "1"
            },
            "evidence": {
                "max_age_duration": "1814400000000000",
                "max_age_num_blocks": "120961",
                "max_bytes": "1048576"
            },
            "validator": {
                "pub_key_types": [
                    "ed25519"
                ]
            },
            "version": {}
        },
    ```
