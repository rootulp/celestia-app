package txsim

import (
	"context"
	"math/rand"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/grpc"
)

var _ Sequence = &SendSequence{}

const (
	sendGasLimit = 100000
	sendFee      = sendGasLimit * appconsts.DefaultMinGasPrice
)

// SendSequence sets up an endless sequence of send transactions, moving tokens
// between a set of accounts
type SendSequence struct {
	numAccounts    int
	sendAmount     int
	maxHeightDelay int
	accounts       []types.AccAddress
	index          int
	numIterations  int
}

func NewSendSequence(numAccounts, sendAmount, numIterations int) *SendSequence {
	return &SendSequence{
		numAccounts:    numAccounts,
		sendAmount:     sendAmount,
		maxHeightDelay: 5,
		numIterations:  numIterations,
	}
}

func (s *SendSequence) Clone(n int) []Sequence {
	sequenceGroup := make([]Sequence, n)
	for i := 0; i < n; i++ {
		sequenceGroup[i] = NewSendSequence(s.numAccounts, s.sendAmount, s.numIterations)
	}
	return sequenceGroup
}

// Init sets up the accounts involved in the sequence. It calculates the necessary balance as the fees per transaction
// multiplied by the number of expected iterations plus the amount to be sent from one account to another
func (s *SendSequence) Init(_ context.Context, _ grpc.ClientConn, allocateAccounts AccountAllocator, _ *rand.Rand) {
	amount := s.sendAmount + (s.numIterations * int(sendFee))
	s.accounts = allocateAccounts(s.numAccounts, amount)
}

// Next sumbits a transaction to remove funds from one account to the next
func (s *SendSequence) Next(_ context.Context, _ grpc.ClientConn, rand *rand.Rand) (Operation, error) {
	if s.index >= s.numIterations {
		return Operation{}, EndOfSequence
	}
	op := Operation{
		Msgs: []types.Msg{
			bank.NewMsgSend(s.accounts[s.index%s.numAccounts], s.accounts[(s.index+1)%s.numAccounts], types.NewCoins(types.NewInt64Coin(app.BondDenom, int64(s.sendAmount)))),
		},
		Delay:    rand.Int63n(int64(s.maxHeightDelay)),
		GasLimit: sendGasLimit,
	}
	s.index++
	return op, nil
}
