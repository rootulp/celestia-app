package types

import (
	"testing"

	sdkerrors "cosmossdk.io/errors"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/stretchr/testify/assert"
)

func TestWirePayForData_ValidateBasic(t *testing.T) {
	type test struct {
		name    string
		msg     *MsgWirePayForData
		wantErr *sdkerrors.Error
	}

	// valid pfd
	validMsg := validWirePayForData(t)

	// pfd with bad ns id
	badIDMsg := validWirePayForData(t)
	badIDMsg.MessageNamespaceId = []byte{1, 2, 3, 4, 5, 6, 7}

	// pfd that uses reserved ns id
	reservedMsg := validWirePayForData(t)
	reservedMsg.MessageNamespaceId = []byte{0, 0, 0, 0, 0, 0, 0, 100}

	// pfd that uses parity shares namespace id
	paritySharesMsg := validWirePayForData(t)
	paritySharesMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// pfd that uses parity shares namespace id
	tailPaddingMsg := validWirePayForData(t)
	tailPaddingMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}

	// pfd that has a wrong msg size
	invalidDeclaredMsgSizeMsg := validWirePayForData(t)
	invalidDeclaredMsgSizeMsg.MessageSize = 999

	// pfd with bad commitment
	badCommitMsg := validWirePayForData(t)
	badCommitMsg.MessageShareCommitment.ShareCommitment = []byte{1, 2, 3, 4}

	// pfd that has invalid square size (not power of 2)
	invalidSquareSizeMsg := validWirePayForData(t)
	invalidSquareSizeMsg.MessageShareCommitment.SquareSize = 15

	// pfd that signs over all squares but the first one
	missingCommitmentForOneSquareSize := validWirePayForData(t)
	missingCommitmentForOneSquareSize.MessageShareCommitment = missingCommitmentForOneSquareSize.MessageShareCommitment

	// pfd that signed over no squares
	noMessageShareCommitments := validWirePayForData(t)
	noMessageShareCommitments.MessageShareCommitment = &ShareCommitAndSignature{}

	tests := []test{
		{
			name:    "valid msg",
			msg:     validMsg,
			wantErr: nil,
		},
		{
			name:    "bad ns ID",
			msg:     badIDMsg,
			wantErr: ErrInvalidNamespaceLen,
		},
		{
			name:    "reserved ns id",
			msg:     reservedMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "bad declared message size",
			msg:     invalidDeclaredMsgSizeMsg,
			wantErr: ErrDeclaredActualDataSizeMismatch,
		},
		{
			name:    "bad commitment",
			msg:     badCommitMsg,
			wantErr: ErrInvalidShareCommit,
		},
		{
			name:    "invalid square size",
			msg:     invalidSquareSizeMsg,
			wantErr: ErrCommittedSquareSizeNotPowOf2,
		},
		{
			name:    "parity shares namespace id",
			msg:     paritySharesMsg,
			wantErr: ErrParitySharesNamespace,
		},
		{
			name:    "tail padding namespace id",
			msg:     tailPaddingMsg,
			wantErr: ErrTailPaddingNamespace,
		},
		{
			name:    "no message share commitments",
			msg:     noMessageShareCommitments,
			wantErr: ErrNoMessageShareCommitments,
		},
		{
			name:    "missing commitment for one square size",
			msg:     missingCommitmentForOneSquareSize,
			wantErr: ErrInvalidShareCommitments,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr != nil {
				assert.Contains(t, err.Error(), tt.wantErr.Error())
				space, code, log := sdkerrors.ABCIInfo(err, false)
				assert.Equal(t, tt.wantErr.Codespace(), space)
				assert.Equal(t, tt.wantErr.ABCICode(), code)
				t.Log(log)
			}
		})
	}
}

func TestMsgMinSquareSize(t *testing.T) {
	type testCase struct {
		msgLen   uint64
		expected uint64
	}
	tests := []testCase{
		{
			msgLen:   1,
			expected: 2,
		},
		{
			msgLen:   100,
			expected: 2,
		},
		{
			msgLen:   appconsts.SparseShareContentSize * 4,
			expected: 4,
		},
		{
			msgLen:   appconsts.SparseShareContentSize * 16,
			expected: 8,
		},
	}
	for _, tc := range tests {
		got := MsgMinSquareSize(tc.msgLen)
		assert.Equal(t, tc.expected, got)
	}
}
