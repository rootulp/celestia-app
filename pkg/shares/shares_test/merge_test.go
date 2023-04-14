package shares_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/testutil"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	core "github.com/tendermint/tendermint/proto/tendermint/types"
	coretypes "github.com/tendermint/tendermint/types"
)

func Test_merge_new(t *testing.T) {
	encConf := encoding.MakeConfig(app.ModuleEncodingRegisters...)
	accounts := generateAccounts(200) // 100 for creating blob txs, 100 for creating send txs
	testApp, kr := testutil.SetupTestAppWithGenesisValSet(accounts...)
	blobCount := 100
	blobSize := 400

	timer := time.After(time.Second * 20)
	for {
		select {
		case <-timer:
			return
		default:
			blobTxs := testutil.RandBlobTxsWithAccounts(
				t,
				testApp,
				encConf.TxConfig.TxEncoder(),
				kr,
				blobSize,
				blobCount,
				true,
				"testapp",
				accounts[:blobCount],
			)
			// create 100 send transactions
			sendTxs := testutil.SendTxsWithAccounts(
				t,
				testApp,
				encConf.TxConfig.TxEncoder(),
				kr,
				1000, // amount
				accounts[0],
				accounts[len(accounts)-100:],
				"testapp",
			)
			txs := append(sendTxs, blobTxs...)
			resp := testApp.PrepareProposal(abci.RequestPrepareProposal{
				BlockData: &core.Data{
					Txs: coretypes.Txs(txs).ToSliceOfBytes(),
				},
			})
			fmt.Printf("%v", resp)
		}
	}
}

func generateAccounts(count int) []string {
	accounts := make([]string, 200)
	for i := range accounts {
		accounts[i] = tmrand.Str(20)
	}
	return accounts
}
