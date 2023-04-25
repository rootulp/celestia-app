package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/stretchr/testify/require"
)

func TestAminoCodecFullDecodeAndEncode(t *testing.T) {
	// This tx comes from https://github.com/cosmos/cosmos-sdk/issues/8117.
	txSigned := `{"type":"cosmos-sdk/StdTx","value":{"msg":[{"type":"cosmos-sdk/MsgCreateValidator","value":{"description":{"moniker":"fulltest","identity":"satoshi","website":"example.com","details":"example inc"},"commission":{"rate":"0.500000000000000000","max_rate":"1.000000000000000000","max_change_rate":"0.200000000000000000"},"min_self_delegation":"1000000","delegator_address":"cosmos14pt0q5cwf38zt08uu0n6yrstf3rndzr5057jys","validator_address":"cosmosvaloper14pt0q5cwf38zt08uu0n6yrstf3rndzr52q28gr","pubkey":{"type":"tendermint/PubKeyEd25519","value":"CYrOiM3HtS7uv1B1OAkknZnFYSRpQYSYII8AtMMtev0="},"value":{"denom":"umuon","amount":"700000000"}}}],"fee":{"amount":[{"denom":"umuon","amount":"6000"}],"gas":"160000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"AwAOXeWgNf1FjMaayrSnrOOKz+Fivr6DiI/i0x0sZCHw"},"signature":"RcnfS/u2yl7uIShTrSUlDWvsXo2p2dYu6WJC8VDVHMBLEQZWc8bsINSCjOnlsIVkUNNe1q/WCA9n3Gy1+0zhYA=="}],"memo":"","timeout_height":"0"}}`
	legacyCdc := makeBlobEncodingConfig().Amino
	var tx legacytx.StdTx
	err := legacyCdc.UnmarshalJSON([]byte(txSigned), &tx)
	require.NoError(t, err)

	// Marshalling/unmarshalling the tx should work.
	marshaledTx, err := legacyCdc.MarshalJSON(tx)
	require.NoError(t, err)
	require.Equal(t, string(marshaledTx), txSigned)

	// Marshalling/unmarshalling the tx wrapped in a struct should work.
	txRequest := &cli.BroadcastReq{
		Mode: "block",
		Tx:   tx,
	}
	_, err = legacyCdc.MarshalJSON(txRequest)
	require.NoError(t, err)
}
