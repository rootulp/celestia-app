package posthandler_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/app/encoding"
	"github.com/celestiaorg/celestia-app/pkg/user"
	"github.com/celestiaorg/celestia-app/test/util/testnode"
	upgradetypes "github.com/celestiaorg/celestia-app/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	utia = 1
	tia  = 1e6
)

func TestRefundGasRemaining(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping refund gas remaining test in short mode.")
	}
	suite.Run(t, new(RefundGasRemainingSuite))
}

type RefundGasRemainingSuite struct {
	suite.Suite

	ctx      testnode.Context
	encCfg   encoding.Config
	signer   *user.Signer
	feePayer *user.Signer
}

func (s *RefundGasRemainingSuite) SetupSuite() {
	require := s.Require()
	s.encCfg = encoding.MakeConfig(app.ModuleEncodingRegisters...)
	s.ctx, _, _ = testnode.NewNetwork(s.T(), testnode.DefaultConfig().WithFundedAccounts("a", "b"))
	_, err := s.ctx.WaitForHeight(1)
	require.NoError(err)

	recordA, err := s.ctx.Keyring.Key("a")
	require.NoError(err)
	addrA, err := recordA.GetAddress()
	require.NoError(err)
	s.signer, err = user.SetupSigner(s.ctx.GoContext(), s.ctx.Keyring, s.ctx.GRPCClient, addrA, s.encCfg)
	require.NoError(err)

	recordB, err := s.ctx.Keyring.Key("b")
	require.NoError(err)
	addrB, err := recordB.GetAddress()
	require.NoError(err)
	s.feePayer, err = user.SetupSigner(s.ctx.GoContext(), s.ctx.Keyring, s.ctx.GRPCClient, addrB, s.encCfg)
	require.NoError(err)

	msg, err := feegrant.NewMsgGrantAllowance(&feegrant.BasicAllowance{}, s.feePayer.Address(), s.signer.Address())
	require.NoError(err)
	options := []user.TxOption{user.SetGasLimit(1e6), user.SetFee(tia)}
	resp, err := s.feePayer.SubmitTx(s.ctx.GoContext(), []sdk.Msg{msg}, options...)
	require.NoError(err)
	require.Equal(abci.CodeTypeOK, resp.Code)
}

func (s *RefundGasRemainingSuite) TestDecorator() {
	t := s.T()

	type testCase struct {
		name                string
		gasLimit            uint64
		fee                 uint64
		wantRefund          int64
		feePayer            sdk.AccAddress
		wantRefundRecipient sdk.AccAddress
	}

	testCases := []testCase{
		{
			// Note: gasPrice * gasLimit = fee. So gasPrice = 1 utia.
			name:                "part of the fee should be refunded",
			gasLimit:            1e5, // 100_000
			fee:                 1e5, // 100_000 utia
			wantRefund:          23069,
			wantRefundRecipient: s.signer.Address(),
		},
		{
			// Note: gasPrice * gasLimit = fee. So gasPrice = 10 utia.
			name:                "refund should vary based on gasPrice",
			gasLimit:            1e5, // 100_000
			fee:                 tia, // 1_000_000 utia
			wantRefund:          229730,
			wantRefundRecipient: s.signer.Address(),
		},
		{
			name:                "refund should be at most half of the fee",
			gasLimit:            1e6, // 1_000_000 is way higher than gas consumed by this tx
			fee:                 tia,
			wantRefund:          tia * .5,
			wantRefundRecipient: s.signer.Address(),
		},
		{
			name:                "refund should be sent to fee payer if specified",
			gasLimit:            1e6,
			fee:                 tia,
			feePayer:            s.feePayer.Address(),
			wantRefund:          tia * .5,
			wantRefundRecipient: s.feePayer.Address(),
		},
		{
			name:                "no refund should be sent if gasLimit isn't high enough to pay for the refund gas cost",
			gasLimit:            65000,
			fee:                 65000,
			wantRefund:          0,
			wantRefundRecipient: s.signer.Address(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := []user.TxOption{user.SetGasLimit(tc.gasLimit), user.SetFee(tc.fee)}
			if tc.feePayer != nil {
				// Cosmos SDK has confusing naming but invoke SetFeeGranter
				// instead of SetFeePayer.
				//
				// https://github.com/cosmos/cosmos-sdk/issues/18886
				options = append(options, user.SetFeeGranter(tc.feePayer))
			}
			msg := upgradetypes.NewMsgTryUpgrade(s.signer.Address())

			resp, err := s.signer.SubmitTx(s.ctx.GoContext(), []sdk.Msg{msg}, options...)
			require.NoError(t, err)
			require.EqualValues(t, abci.CodeTypeOK, resp.Code)

			got := getRefund(t, resp, tc.wantRefundRecipient.String())
			assert.Equal(t, tc.wantRefund, got)
		})
	}
}

// getRefund returns the amount refunded to the feePayer based on the events in the TxResponse.
func getRefund(t *testing.T, resp *sdk.TxResponse, feePayer string) (refund int64) {
	assert.NotNil(t, resp)
	transfers := getTransfers(t, resp.Events)
	for _, transfer := range transfers {
		if transfer.recipient == feePayer {
			return transfer.amount
		}
	}
	return refund
}

// getTransfers returns all the transfer events in the slice of events.
func getTransfers(t *testing.T, events []abci.Event) (transfers []transferEvent) {
	for _, event := range events {
		if event.Type == banktypes.EventTypeTransfer {
			amount, err := strconv.ParseInt(strings.TrimSuffix(string(event.Attributes[2].Value), "utia"), 10, 64)
			assert.NoError(t, err)
			transfer := transferEvent{
				recipient: string(event.Attributes[0].Value),
				from:      string(event.Attributes[1].Value),
				amount:    amount,
			}
			transfers = append(transfers, transfer)
		}
	}
	return transfers
}

// transferEvent is a struct based on the transfer event type emitted by the
// bank module.
type transferEvent struct {
	recipient string
	from      string
	amount    int64
}
