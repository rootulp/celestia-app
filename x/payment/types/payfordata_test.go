package types

import (
	"bytes"
	"testing"

	sdkerrors "cosmossdk.io/errors"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	shares "github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/nmt/namespace"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_merkleMountainRangeHeights(t *testing.T) {
	type test struct {
		totalSize  uint64
		squareSize uint64
		expected   []uint64
	}
	tests := []test{
		{
			totalSize:  11,
			squareSize: 4,
			expected:   []uint64{4, 4, 2, 1},
		},
		{
			totalSize:  2,
			squareSize: 64,
			expected:   []uint64{2},
		},
		{
			totalSize:  64,
			squareSize: 8,
			expected:   []uint64{8, 8, 8, 8, 8, 8, 8, 8},
		},
		// Height
		// 3              x                               x
		//              /    \                         /    \
		//             /      \                       /      \
		//            /        \                     /        \
		//           /          \                   /          \
		// 2        x            x                 x            x
		//        /   \        /   \             /   \        /   \
		// 1     x     x      x     x           x     x      x     x         x
		//      / \   / \    / \   / \         / \   / \    / \   / \      /   \
		// 0   0   1 2   3  4   5 6   7       8   9 10  11 12 13 14  15   16   17    18
		{
			totalSize:  19,
			squareSize: 8,
			expected:   []uint64{8, 8, 2, 1},
		},
	}
	for _, tt := range tests {
		res := merkleMountainRangeSizes(tt.totalSize, tt.squareSize)
		assert.Equal(t, tt.expected, res)
	}
}

// TestCreateCommitment only shows if something changed, it doesn't actually
// show that the commitment bytes are being created correctly.
// TODO: verify the commitment bytes
func TestCreateCommitment(t *testing.T) {
	type test struct {
		name       string
		squareSize uint64
		namespace  []byte
		message    []byte
		expected   []byte
		expectErr  bool
	}
	tests := []test{
		{
			name:       "squareSize 4, message of 11 shares succeeds",
			squareSize: 4,
			namespace:  bytes.Repeat([]byte{0xFF}, 8),
			message:    bytes.Repeat([]byte{0xFF}, 11*ShareSize),
			expected:   []byte{0x1e, 0xdc, 0xc4, 0x69, 0x8f, 0x47, 0xf6, 0x8d, 0xfc, 0x11, 0xec, 0xac, 0xaa, 0x37, 0x4a, 0x3d, 0xbd, 0xfc, 0x1a, 0x9b, 0x6e, 0x87, 0x6f, 0xba, 0xd3, 0x6c, 0x6, 0x6c, 0x9f, 0x5b, 0x65, 0x38},
		},
		{
			name:       "squareSize 4, message of 12 shares succeeds",
			squareSize: 12,
			namespace:  bytes.Repeat([]byte{0xFF}, 8),
			message:    bytes.Repeat([]byte{0xFF}, 12*ShareSize),
			expected:   []byte{0x81, 0x5e, 0xf9, 0x52, 0x2a, 0xfa, 0x40, 0x67, 0x63, 0x64, 0x4a, 0x82, 0x7, 0xcd, 0x1d, 0x7d, 0x1f, 0xae, 0xe5, 0xd3, 0xb1, 0x91, 0x8a, 0xb8, 0x90, 0x51, 0xfc, 0x1, 0xd, 0xa7, 0xf3, 0x1a},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := CreateCommitment(tt.namespace, tt.message)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

// TestSignMalleatedTxs checks to see that the signatures that are generated for
// the PayForDatas malleated from the original WirePayForData are actually
// valid.
func TestSignMalleatedTxs(t *testing.T) {
	type test struct {
		name    string
		ns, msg []byte
		ss      []uint64
		options []TxBuilderOption
	}

	kb := generateKeyring(t, "test")

	signer := NewKeyringSigner(kb, "test", "test-chain-id")

	tests := []test{
		{
			name:    "single share",
			ns:      []byte{1, 1, 1, 1, 1, 1, 1, 1},
			msg:     bytes.Repeat([]byte{1}, appconsts.SparseShareContentSize),
			ss:      []uint64{2, 4, 8, 16},
			options: []TxBuilderOption{SetGasLimit(2000000)},
		},
		{
			name: "12 shares",
			ns:   []byte{1, 1, 1, 1, 1, 1, 1, 2},
			msg:  bytes.Repeat([]byte{2}, (appconsts.SparseShareContentSize*12)-4), // subtract a few bytes for the delimiter
			ss:   []uint64{4, 8, 16, 64},
			options: []TxBuilderOption{
				SetGasLimit(123456789),
				SetFeeAmount(sdk.NewCoins(sdk.NewCoin("ucls", sdk.NewInt(987654321)))),
			},
		},
		{
			name: "12 shares",
			ns:   []byte{1, 1, 1, 1, 1, 1, 1, 2},
			msg:  bytes.Repeat([]byte{1, 2, 3, 4, 5}, 10000), // subtract a few bytes for the delimiter
			ss:   AllSquareSizes(50000),
			options: []TxBuilderOption{
				SetGasLimit(123456789),
				SetFeeAmount(sdk.NewCoins(sdk.NewCoin("ucls", sdk.NewInt(987654321)))),
			},
		},
	}

	for _, tt := range tests {
		wpfd, err := NewWirePayForData(tt.ns, tt.msg)
		require.NoError(t, err, tt.name)
		err = wpfd.SignShareCommitments(signer, tt.options...)
		// there should be no error
		assert.NoError(t, err)
		// the signature should exist
		assert.Equal(t, len(wpfd.MessageShareCommitment.Signature), 64)

		sData, err := signer.GetSignerData()
		require.NoError(t, err)

		wpfdTx, err := signer.BuildSignedTx(signer.NewTxBuilder(tt.options...), wpfd)
		require.NoError(t, err)

		// VerifyPFDSigs goes through the entire malleation process for every
		// square size, creating PfDs from the wirePfD and check that the
		// signature is valid
		valid, err := VerifyPFDSigs(sData, signer.encCfg.TxConfig, wpfdTx)
		assert.NoError(t, err)
		assert.True(t, valid, tt.name)
	}
}

func TestProcessMessage(t *testing.T) {
	type test struct {
		name       string
		namespace  []byte
		msg        []byte
		squareSize uint64
		expectErr  bool
		modify     func(*MsgWirePayForData) *MsgWirePayForData
	}

	dontModify := func(in *MsgWirePayForData) *MsgWirePayForData {
		return in
	}

	kb := generateKeyring(t, "test")

	signer := NewKeyringSigner(kb, "test", "chain-id")

	tests := []test{
		{
			name:       "single share square size 2",
			namespace:  []byte{1, 1, 1, 1, 1, 1, 1, 1},
			msg:        bytes.Repeat([]byte{1}, totalMsgSize(appconsts.SparseShareContentSize)),
			squareSize: 2,
			modify:     dontModify,
		},
		{
			name:       "12 shares square size 4",
			namespace:  []byte{1, 1, 1, 1, 1, 1, 1, 2},
			msg:        bytes.Repeat([]byte{2}, totalMsgSize(appconsts.SparseShareContentSize*12)),
			squareSize: 4,
			modify:     dontModify,
		},
	}

	for _, tt := range tests {
		wpfd, err := NewWirePayForData(tt.namespace, tt.msg)
		require.NoError(t, err, tt.name)
		err = wpfd.SignShareCommitments(signer)
		assert.NoError(t, err)

		wpfd = tt.modify(wpfd)

		message, spfd, sig, err := ProcessWirePayForData(wpfd)
		if tt.expectErr {
			assert.Error(t, err, tt.name)
			continue
		}

		// ensure that the shared fields are identical
		assert.Equal(t, tt.msg, message.Data, tt.name)
		assert.Equal(t, tt.namespace, message.NamespaceId, tt.name)
		assert.Equal(t, wpfd.Signer, spfd.Signer, tt.name)
		assert.Equal(t, wpfd.MessageNamespaceId, spfd.MessageNamespaceId, tt.name)
		assert.Equal(t, wpfd.MessageShareCommitment.ShareCommitment, spfd.MessageShareCommitment, tt.name)
		assert.Equal(t, wpfd.MessageShareCommitment.Signature, sig, tt.name)
	}
}

func TestValidateBasic(t *testing.T) {
	type test struct {
		name    string
		msg     *MsgPayForData
		wantErr *sdkerrors.Error
	}

	validMsg := validMsgPayForData(t)

	// MsgPayForData that uses parity shares namespace id
	paritySharesMsg := validMsgPayForData(t)
	paritySharesMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	// MsgPayForData that uses tail padding namespace id
	tailPaddingMsg := validMsgPayForData(t)
	tailPaddingMsg.MessageNamespaceId = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE}

	// MsgPayForData that uses transaction namespace id
	txNamespaceMsg := validMsgPayForData(t)
	txNamespaceMsg.MessageNamespaceId = namespace.ID{0, 0, 0, 0, 0, 0, 0, 1}

	// MsgPayForData that uses intermediateStateRoots namespace id
	intermediateStateRootsNamespaceMsg := validMsgPayForData(t)
	intermediateStateRootsNamespaceMsg.MessageNamespaceId = namespace.ID{0, 0, 0, 0, 0, 0, 0, 2}

	// MsgPayForData that uses evidence namespace id
	evidenceNamespaceMsg := validMsgPayForData(t)
	evidenceNamespaceMsg.MessageNamespaceId = namespace.ID{0, 0, 0, 0, 0, 0, 0, 3}

	// MsgPayForData that uses the max reserved namespace id
	maxReservedNamespaceMsg := validMsgPayForData(t)
	maxReservedNamespaceMsg.MessageNamespaceId = namespace.ID{0, 0, 0, 0, 0, 0, 0, 255}

	// MsgPayForData that has no message share commitments
	noMessageShareCommitments := validMsgPayForData(t)
	noMessageShareCommitments.MessageShareCommitment = []byte{}

	tests := []test{
		{
			name:    "valid msg",
			msg:     validMsg,
			wantErr: nil,
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
			name:    "transaction namspace namespace id",
			msg:     txNamespaceMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "intermediate state root namespace id",
			msg:     intermediateStateRootsNamespaceMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "evidence namspace namespace id",
			msg:     evidenceNamespaceMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "max reserved namespace id",
			msg:     maxReservedNamespaceMsg,
			wantErr: ErrReservedNamespace,
		},
		{
			name:    "no message share commitments",
			msg:     noMessageShareCommitments,
			wantErr: ErrNoMessageShareCommitments,
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

// totalMsgSize subtracts the delimiter size from the desired total size. this
// is useful for testing for messages that occupy exactly so many shares.
func totalMsgSize(size int) int {
	return size - shares.DelimLen(uint64(size))
}

func validWirePayForData(t *testing.T) *MsgWirePayForData {
	message := bytes.Repeat([]byte{1}, 2000)
	msg, err := NewWirePayForData(
		[]byte{1, 2, 3, 4, 5, 6, 7, 8},
		message,
	)
	if err != nil {
		panic(err)
	}

	signer := generateKeyringSigner(t)

	err = msg.SignShareCommitments(signer)
	if err != nil {
		panic(err)
	}
	return msg
}

func validMsgPayForData(t *testing.T) *MsgPayForData {
	kb := generateKeyring(t, "test")
	signer := NewKeyringSigner(kb, "test", "chain-id")
	ns := []byte{1, 1, 1, 1, 1, 1, 1, 2}
	msg := bytes.Repeat([]byte{2}, totalMsgSize(appconsts.SparseShareContentSize*12))

	wpfd, err := NewWirePayForData(ns, msg)
	assert.NoError(t, err)

	err = wpfd.SignShareCommitments(signer)
	assert.NoError(t, err)

	_, spfd, _, err := ProcessWirePayForData(wpfd)
	require.NoError(t, err)

	return spfd
}
