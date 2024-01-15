package app_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/pkg/user"
	"github.com/celestiaorg/celestia-app/test/util/blobfactory"
	"github.com/celestiaorg/celestia-app/test/util/testfactory"
	"github.com/celestiaorg/celestia-app/test/util/testnode"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestGasPricingSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping gas pricing test in short mode.")
	}
	suite.Run(t, new(GasPricingSuite))
}

type GasPricingSuite struct {
	suite.Suite

	accounts []string
	cfg      *testnode.Config
	cctx     testnode.Context
	ecfg     encoding.Config

	mut            sync.Mutex
	accountCounter int
}

func (s *GasPricingSuite) SetupSuite() {
	t := s.T()
	t.Log("setting up test suite")

	s.accounts = testnode.RandomAccounts(10)
	s.cfg = testnode.DefaultConfig().WithFundedAccounts(s.accounts...)
	s.cctx, _, _ = testnode.NewNetwork(t, s.cfg)
	s.ecfg = encoding.MakeConfig(app.ModuleEncodingRegisters...)
}

func (s *GasPricingSuite) unusedAccount() string {
	s.mut.Lock()
	acc := s.accounts[s.accountCounter]
	s.accountCounter++
	s.mut.Unlock()
	return acc
}

func (s *GasPricingSuite) getValidatorName() string {
	return s.cfg.Genesis.Validators()[0].Name
}

func (s *GasPricingSuite) getValidatorAccount() sdk.ValAddress {
	record, err := s.cfg.Genesis.Keyring().Key(s.getValidatorName())
	s.Require().NoError(err)
	address, err := record.GetAddress()
	s.Require().NoError(err)
	return sdk.ValAddress(address)
}

func (s *GasPricingSuite) TestGasPricing() {
	t := s.T()
	memoOptions := []user.TxOption{}
	memoOptions = append(memoOptions, blobfactory.DefaultTxOpts()...)
	memoOptions = append(memoOptions, user.SetMemo(strings.Repeat("a", 256)))

	type testCase struct {
		name         string
		msgFunc      func() (msgs []sdk.Msg, signer string)
		txOptions    []user.TxOption
		expectedCode uint32
		wantGasUsed  int64
	}

	testCases := []testCase{
		{
			name: "send 1 utia",
			msgFunc: func() (msgs []sdk.Msg, signer string) {
				account1, account2 := s.unusedAccount(), s.unusedAccount()
				msgSend := banktypes.NewMsgSend(
					testfactory.GetAddress(s.cctx.Keyring, account1),
					testfactory.GetAddress(s.cctx.Keyring, account2),
					sdk.NewCoins(sdk.NewCoin(app.BondDenom, sdk.NewInt(1))),
				)
				return []sdk.Msg{msgSend}, account1
			},
			txOptions:    blobfactory.DefaultTxOpts(),
			expectedCode: abci.CodeTypeOK,
			wantGasUsed:  77004,
			// When auth.TxSizeCostPerByte = 10, gasUsed by tx size is 3170.
			// So fixed cost = 77004 - 3170 = 73834
			// When auth.TxSizeCostPerByte = 100, gasUsed by tx size is 31700.
			// So total cost is 73734 + 31700 = 105434
			// When auth.TxSizeCostPerByte = 1000, gasUsed by tx size is 317000.
			// So total cost is 73734 + 317000 = 320734
		},
		{
			name: "send 1 utia with 256 character memo",
			msgFunc: func() (msgs []sdk.Msg, signer string) {
				account1, account2 := s.unusedAccount(), s.unusedAccount()
				msgSend := banktypes.NewMsgSend(
					testfactory.GetAddress(s.cctx.Keyring, account1),
					testfactory.GetAddress(s.cctx.Keyring, account2),
					sdk.NewCoins(sdk.NewCoin(app.BondDenom, sdk.NewInt(1))),
				)
				return []sdk.Msg{msgSend}, account1
			},
			txOptions:    memoOptions,
			expectedCode: abci.CodeTypeOK,
			wantGasUsed:  79594,
			// gasUsed by tx size is 5760
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msgs, account := tc.msgFunc()
			addr := testfactory.GetAddress(s.cctx.Keyring, account)
			signer, err := user.SetupSigner(s.cctx.GoContext(), s.cctx.Keyring, s.cctx.GRPCClient, addr, s.ecfg)
			require.NoError(t, err)
			fmt.Printf("submitting %v\n", tc.name)
			res, err := signer.SubmitTx(s.cctx.GoContext(), msgs, tc.txOptions...)
			if tc.expectedCode != abci.CodeTypeOK {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, res)
			assert.Equal(t, tc.expectedCode, res.Code, res.RawLog)
			assert.Equal(t, tc.wantGasUsed, res.GasUsed)
		})
	}
}
