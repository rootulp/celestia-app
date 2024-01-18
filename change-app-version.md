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

1. Query the consensus params via the REST API. Observe that the app version has been updated to 3.

    ```bash
    curl -X GET http://0.0.0.0:26657/consensus_params  | jq .
    ```

    ```json
    {
    "jsonrpc": "2.0",
    "id": -1,
    "result": {
        "block_height": "34",
        "consensus_params": {
        "block": {
            "max_bytes": "1974272",
            "max_gas": "-1",
            "time_iota_ms": "1"
        },
        "evidence": {
            "max_age_num_blocks": "120961",
            "max_age_duration": "1814400000000000",
            "max_bytes": "1048576"
        },
        "validator": {
            "pub_key_types": [
            "ed25519"
            ]
        },
        "version": {
            "app_version": "3"
        }
        }
    }
    }
    ```

1. Query the most recent block. Observe that the app version in the header has been updated.

    ```bash
    curl -X GET http://0.0.0.0:26657/block?height=36  | jq .
    ```

    ```json
    {
    "jsonrpc": "2.0",
    "id": -1,
    "result": {
        "block_id": {
        "hash": "ED84FD959FCC69CF943113EFE1A0AF29B23C179C5AE252980840A59F9BFA2CA5",
        "parts": {
            "total": 1,
            "hash": "EEB79C3A23EC93B8AB4BC9B414E457CEFE5C72D7C8872D70F6D3C8AD70E2CBEF"
        }
        },
        "block": {
        "header": {
            "version": {
            "block": "11",
            "app": "3"
            },
            "chain_id": "private",
            "height": "36",
            "time": "2024-01-18T21:41:18.676555Z",
            "last_block_id": {
            "hash": "CE19AEE5447C0FFC5064861F602C192A087662D2443D58B028C40118714ECB7F",
            "parts": {
                "total": 1,
                "hash": "7CD7A0490ED1616028D7D1C0E08D0700D19F45DDA8406E2D3A1ABB06B89A47CB"
            }
            },
            "last_commit_hash": "3950A7ACAAD71F29CCE0D048997D7FA6FE42CC7D5A86073AA6318BC4CCD258DC",
            "data_hash": "3D96B7D238E7E0456F6AF8E7CDF0A67BD6CF9C2089ECB559C659DCAA1F880353",
            "validators_hash": "12811879BCAF510B19DF19DDC66E12FA2B71E76B7426FE29D36D695AE847B795",
            "next_validators_hash": "12811879BCAF510B19DF19DDC66E12FA2B71E76B7426FE29D36D695AE847B795",
            "consensus_hash": "C0B6A634B72AE9687EA53B6D277A73ABA1386BA3CFC6D0F26963602F7F6FFCD6",
            "app_hash": "D29EC9DCE455B31A939BB36366FB7BECA65B62ACB7596AA85EF907AB3D9DA5CD",
            "last_results_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "evidence_hash": "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855",
            "proposer_address": "64D9A517E6B8A1488B4C0D20CA3593B94C8DA87F"
        },
        "data": {
            "txs": [],
            "square_size": "1"
        },
        "evidence": {
            "evidence": []
        },
        "last_commit": {
            "height": "35",
            "round": 0,
            "block_id": {
            "hash": "CE19AEE5447C0FFC5064861F602C192A087662D2443D58B028C40118714ECB7F",
            "parts": {
                "total": 1,
                "hash": "7CD7A0490ED1616028D7D1C0E08D0700D19F45DDA8406E2D3A1ABB06B89A47CB"
            }
            },
            "signatures": [
            {
                "block_id_flag": 2,
                "validator_address": "64D9A517E6B8A1488B4C0D20CA3593B94C8DA87F",
                "timestamp": "2024-01-18T21:41:18.676555Z",
                "signature": "225ZL2Ve2IRRNCk5fYEO2zcNcdk/LkNdZmoe7QNqQhuuVejcPHNs28eOHbamkiV6/54rYTrqBA4BDjp+RxVkDA=="
            }
            ]
        }
        }
    }
    }
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
