package app_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/blob"
	"github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/celestiaorg/celestia-app/pkg/user"
	"github.com/celestiaorg/celestia-app/test/util/blobfactory"
	"github.com/celestiaorg/celestia-app/test/util/testfactory"
	"github.com/celestiaorg/celestia-app/test/util/testnode"
	blobtypes "github.com/celestiaorg/celestia-app/x/blob/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	oldgov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
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
	app      *app.App

	sourcePort    string
	sourceChannel string

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
		name        string
		msgFunc     func() (msgs []sdk.Msg, signer string)
		txOptions   []user.TxOption
		wantGasUsed int64
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
			txOptions:   blobfactory.DefaultTxOpts(),
			wantGasUsed: 77004,
			// tx size is 317 bytes
			// When auth.TxSizeCostPerByte = 10, gasUsed by tx size is 3170. So fixed cost = 77004 - 3170 = 73834.
			// When auth.TxSizeCostPerByte = 16, gasUsed by tx size is 5072. So total cost = 73834 + 5072 = 78906.
			// When auth.TxSizeCostPerByte = 100, gasUsed by tx size is 31700. So total cost is 73834 + 31700 = 105534.
			// When auth.TxSizeCostPerByte = 1000, gasUsed by tx size is 317000. So total cost is 73834 + 317000 = 390834.
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
			txOptions:   memoOptions,
			wantGasUsed: 79594,
			// tx size is 576 bytes
			// When auth.TxSizeCostPerByte = 10, gasUsed by tx size is 5760. So fixed cost = 79594 - 5760 = 73834.
			// When auth.TxSizeCostPerByte = 16, gasUsed by tx size is 9216. So total cost = 73834 + 9216 = 83050.
			// When auth.TxSizeCostPerByte = 100, gasUsed by tx size is 57600. So total cost is 73834 + 57600 = 131434.
			// When auth.TxSizeCostPerByte = 1000, gasUsed by tx size is 576000. So total cost is 73834 + 576000 = 649834.
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

			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, tc.wantGasUsed, res.GasUsed)
		})
	}
}

func (s *GasPricingSuite) TestGasPricingBlobTx() {
	t := s.T()

	type testCase struct {
		name        string
		blobs       []*blob.Blob
		txOptions   []user.TxOption
		wantGasUsed int64
	}

	b, err := blobtypes.NewBlob(namespace.RandomNamespace(), tmrand.Bytes(256), appconsts.ShareVersionZero)
	require.NoError(t, err)

	testCases := []testCase{
		{
			name:        "Blob with 256 bytes",
			blobs:       []*blob.Blob{b},
			txOptions:   blobfactory.DefaultTxOpts(),
			wantGasUsed: 67765,
			// tx size is 333 bytes
			// When auth.TxSizeCostPerByte = 10, gasUsed by tx size is 3330. So fixed cost = 67765 - 3330 = 64435.
			// When auth.TxSizeCostPerByte = 16, gasUsed by tx size is 5328. So total cost = 64435 + 5328 = 69763.
			// When auth.TxSizeCostPerByte = 100, gasUsed by tx size is 33300. So total cost is 64435 + 33300 = 97735.
			// When auth.TxSizeCostPerByte = 1000, gasUsed by tx size is 333000. So total cost is 64435 + 333000 = 397435.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			account := s.unusedAccount()
			addr := testfactory.GetAddress(s.cctx.Keyring, account)
			signer, err := user.SetupSigner(s.cctx.GoContext(), s.cctx.Keyring, s.cctx.GRPCClient, addr, s.ecfg)
			require.NoError(t, err)
			fmt.Printf("submitting %v\n", tc.name)
			res, err := signer.SubmitPayForBlob(s.cctx.GoContext(), tc.blobs, tc.txOptions...)

			require.NoError(t, err)
			require.NotNil(t, res)
			assert.Equal(t, tc.wantGasUsed, res.GasUsed)
		})
	}
}

func (s *GasPricingSuite) setTxCostPerByte(t *testing.T, txCostPerByte uint64) {
	account := s.getValidatorName()
	record, err := s.cfg.Genesis.Keyring().Key(account)
	s.Require().NoError(err)
	addr, err := record.GetAddress()
	s.Require().NoError(err)

	paramChange := proposal.NewParamChange(
		authtypes.ModuleName,
		string(authtypes.KeyTxSizeCostPerByte),
		fmt.Sprintf("\"%d\"", txCostPerByte),
	)
	content := proposal.NewParameterChangeProposal("title", "description", []proposal.ParamChange{paramChange})

	msg, err := oldgov.NewMsgSubmitProposal(
		content,
		sdk.NewCoins(
			sdk.NewCoin(appconsts.BondDenom, sdk.NewInt(1000000000))),
		addr,
	)
	require.NoError(t, err)

	signer, err := user.SetupSigner(s.cctx.GoContext(), s.cctx.Keyring, s.cctx.GRPCClient, addr, s.ecfg)
	require.NoError(t, err)

	res, err := signer.SubmitTx(s.cctx.GoContext(), []sdk.Msg{msg}, blobfactory.DefaultTxOpts()...)
	require.NoError(t, err)
	require.Equal(t, res.Code, abci.CodeTypeOK, res.RawLog)

	require.NoError(t, s.cctx.WaitForNextBlock())

	// query the proposal to get the id
	gqc := v1.NewQueryClient(s.cctx.GRPCClient)
	gresp, err := gqc.Proposals(s.cctx.GoContext(), &v1.QueryProposalsRequest{ProposalStatus: v1.ProposalStatus_PROPOSAL_STATUS_VOTING_PERIOD})
	require.NoError(t, err)
	require.Len(t, gresp.Proposals, 1)

	// create and submit a new vote
	vote := v1.NewMsgVote(testfactory.GetAddress(s.cctx.Keyring, account), gresp.Proposals[0].Id, v1.VoteOption_VOTE_OPTION_YES, "")
	res, err = signer.SubmitTx(s.cctx.GoContext(), []sdk.Msg{vote}, blobfactory.DefaultTxOpts()...)
	require.NoError(t, err)
	require.Equal(t, abci.CodeTypeOK, res.Code)

	// wait for the voting period to complete
	time.Sleep(time.Second * 6)

	// check that the parameters got updated as expected
	bqc := authtypes.NewQueryClient(s.cctx.GRPCClient)
	presp, err := bqc.Params(s.cctx.GoContext(), &authtypes.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, txCostPerByte, presp.Params.TxSizeCostPerByte)
}
