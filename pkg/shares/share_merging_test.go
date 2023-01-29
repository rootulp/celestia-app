package shares

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"reflect"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	testns "github.com/celestiaorg/celestia-app/testutil/namespace"
	"github.com/celestiaorg/nmt/namespace"
	"github.com/stretchr/testify/assert"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"github.com/tendermint/tendermint/types"
)

func TestParseShares(t *testing.T) {
	type testCase struct {
		name      string
		shares    [][]byte
		want      []ShareSequence
		expectErr bool
	}

	start := true
	blobOneNamespace := namespace.ID{1, 1, 1, 1, 1, 1, 1, 1}
	blobTwoNamespace := namespace.ID{2, 2, 2, 2, 2, 2, 2, 2}

	txShares, _, _ := SplitTxs(generateRandomTxs(2, 1000))
	txShareStart := txShares[0]
	txShareContinuation := txShares[1]

	blobOneShares, err := SplitBlobs(0, []uint32{}, []types.Blob{generateRandomBlobWithNamespace(blobOneNamespace, 1000)}, false)
	if err != nil {
		t.Fatal(err)
	}
	blobOneStart := blobOneShares[0]
	blobOneContinuation := blobOneShares[1]

	blobTwoShares, err := SplitBlobs(0, []uint32{}, []types.Blob{generateRandomBlobWithNamespace(blobTwoNamespace, 1000)}, false)
	if err != nil {
		t.Fatal(err)
	}
	blobTwoStart := blobTwoShares[0]
	blobTwoContinuation := blobTwoShares[1]

	invalidShare := generateRawShare(blobOneNamespace, start, 1)
	invalidShare = append(invalidShare, []byte{0}...) // invalidShare is now longer than the length of a valid share

	largeSequenceLen := 1000 // it takes more than one share to store a sequence of 1000 bytes
	oneShareWithTooLargeSequenceLen := generateRawShare(blobOneNamespace, start, uint32(largeSequenceLen))

	shortSequenceLen := 0
	oneShareWithTooShortSequenceLen := generateRawShare(blobOneNamespace, start, uint32(shortSequenceLen))

	tests := []testCase{
		{
			"empty",
			[][]byte{},
			[]ShareSequence{},
			false,
		},
		{
			"one transaction share",
			[][]byte{txShareStart},
			[]ShareSequence{{NamespaceID: appconsts.TxNamespaceID, Shares: []Share{txShareStart}}},
			false,
		},
		{
			"two transaction shares",
			[][]byte{txShareStart, txShareContinuation},
			[]ShareSequence{{NamespaceID: appconsts.TxNamespaceID, Shares: []Share{txShareStart, txShareContinuation}}},
			false,
		},
		{
			"one blob share",
			[][]byte{blobOneStart},
			[]ShareSequence{{NamespaceID: blobOneNamespace, Shares: []Share{blobOneStart}}},
			false,
		},
		{
			"two blob shares",
			[][]byte{blobOneStart, blobOneContinuation},
			[]ShareSequence{{NamespaceID: blobOneNamespace, Shares: []Share{blobOneStart, blobOneContinuation}}},
			false,
		},
		{
			"two blobs with two shares each",
			[][]byte{blobOneStart, blobOneContinuation, blobTwoStart, blobTwoContinuation},
			[]ShareSequence{
				{NamespaceID: blobOneNamespace, Shares: []Share{blobOneStart, blobOneContinuation}},
				{NamespaceID: blobTwoNamespace, Shares: []Share{blobTwoStart, blobTwoContinuation}},
			},
			false,
		},
		{
			"one transaction, one blob",
			[][]byte{txShareStart, blobOneStart},
			[]ShareSequence{
				{NamespaceID: appconsts.TxNamespaceID, Shares: []Share{txShareStart}},
				{NamespaceID: blobOneNamespace, Shares: []Share{blobOneStart}},
			},
			false,
		},
		{
			"one transaction, two blobs",
			[][]byte{txShareStart, blobOneStart, blobTwoStart},
			[]ShareSequence{
				{NamespaceID: appconsts.TxNamespaceID, Shares: []Share{txShareStart}},
				{NamespaceID: blobOneNamespace, Shares: []Share{blobOneStart}},
				{NamespaceID: blobTwoNamespace, Shares: []Share{blobTwoStart}},
			},
			false,
		},
		{
			"one share with invalid size",
			[][]byte{invalidShare},
			[]ShareSequence{},
			true,
		},
		{
			"blob one start followed by blob two continuation",
			[][]byte{blobOneStart, blobTwoContinuation},
			[]ShareSequence{},
			true,
		},
		{
			"one share with too large sequence length",
			[][]byte{oneShareWithTooLargeSequenceLen},
			[]ShareSequence{},
			true,
		},
		{
			"one share with too short sequence length",
			[][]byte{oneShareWithTooShortSequenceLen},
			[]ShareSequence{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseShares(tt.shares)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseShares() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShareSequenceRawData(t *testing.T) {
	type testCase struct {
		name          string
		shareSequence ShareSequence
		want          []byte
		wantErr       bool
	}
	blobNamespace := testns.RandomBlobNamespace()

	testCases := []testCase{
		{
			name: "empty share sequence",
			shareSequence: ShareSequence{
				NamespaceID: appconsts.TxNamespaceID,
				Shares:      []Share{},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "one empty share",
			shareSequence: ShareSequence{
				NamespaceID: appconsts.TxNamespaceID,
				Shares: []Share{
					shareWithData(blobNamespace, true, 0, []byte{}),
				},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "one share with one byte",
			shareSequence: ShareSequence{
				NamespaceID: appconsts.TxNamespaceID,
				Shares: []Share{
					shareWithData(blobNamespace, true, 1, []byte{0x0f}),
				},
			},
			want:    []byte{0xf},
			wantErr: false,
		},
		{
			name: "removes padding from last share",
			shareSequence: ShareSequence{
				NamespaceID: appconsts.TxNamespaceID,
				Shares: []Share{
					shareWithData(blobNamespace, true, appconsts.FirstSparseShareContentSize+1, bytes.Repeat([]byte{0xf}, appconsts.FirstSparseShareContentSize)),
					shareWithData(blobNamespace, false, 0, []byte{0x0f}),
				},
			},
			want:    bytes.Repeat([]byte{0xf}, appconsts.FirstSparseShareContentSize+1),
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.shareSequence.RawData()
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_compactSharesNeeded(t *testing.T) {
	type testCase struct {
		sequenceLen int
		want        int
	}
	testCases := []testCase{
		{0, 0},
		{1, 1},
		{2, 1},
		{appconsts.FirstCompactShareContentSize, 1},
		{appconsts.FirstCompactShareContentSize + 1, 2},
		{appconsts.FirstCompactShareContentSize + appconsts.ContinuationCompactShareContentSize, 2},
		{appconsts.FirstCompactShareContentSize + appconsts.ContinuationCompactShareContentSize*100, 101},
	}
	for _, tc := range testCases {
		got := CompactSharesNeeded(tc.sequenceLen)
		assert.Equal(t, tc.want, got)
	}
}

func Test_sparseSharesNeeded(t *testing.T) {
	type testCase struct {
		sequenceLen uint32
		want        int
	}
	testCases := []testCase{
		{0, 0},
		{1, 1},
		{2, 1},
		{appconsts.FirstSparseShareContentSize, 1},
		{appconsts.FirstSparseShareContentSize + 1, 2},
		{appconsts.FirstSparseShareContentSize + appconsts.ContinuationSparseShareContentSize, 2},
		{appconsts.FirstSparseShareContentSize + appconsts.ContinuationCompactShareContentSize*2, 3},
		{appconsts.FirstSparseShareContentSize + appconsts.ContinuationCompactShareContentSize*99, 100},
		{1000, 2},
		{10000, 20},
		{100000, 199},
	}
	for _, tc := range testCases {
		got := SparseSharesNeeded(tc.sequenceLen)
		assert.Equal(t, tc.want, got)
	}
}

func generateRawShare(namespace namespace.ID, isSequenceStart bool, sequenceLen uint32) (rawShare []byte) {
	infoByte, _ := NewInfoByte(appconsts.ShareVersionZero, isSequenceStart)

	sequenceLenBuf := make([]byte, appconsts.SequenceLenBytes)
	binary.BigEndian.PutUint32(sequenceLenBuf, sequenceLen)

	rawShare = append(rawShare, namespace...)
	rawShare = append(rawShare, byte(infoByte))
	rawShare = append(rawShare, sequenceLenBuf...)

	return padWithRandomBytes(rawShare)
}

func shareWithData(namespace namespace.ID, isSequenceStart bool, sequenceLen uint32, data []byte) (rawShare []byte) {
	infoByte, _ := NewInfoByte(appconsts.ShareVersionZero, isSequenceStart)
	rawShare = append(rawShare, namespace...)
	rawShare = append(rawShare, byte(infoByte))
	if isSequenceStart {
		sequenceLenBuf := make([]byte, appconsts.SequenceLenBytes)
		binary.BigEndian.PutUint32(sequenceLenBuf, sequenceLen)
		rawShare = append(rawShare, sequenceLenBuf...)
	}
	rawShare = append(rawShare, data...)

	return padShare(rawShare)
}

func padWithRandomBytes(partialShare Share) (paddedShare Share) {
	paddedShare = make([]byte, appconsts.ShareSize)
	copy(paddedShare, partialShare)
	rand.Read(paddedShare[len(partialShare):])
	return paddedShare
}

func generateRandomTxs(count, size int) types.Txs {
	txs := make(types.Txs, count)
	for i := 0; i < count; i++ {
		tx := make([]byte, size)
		_, err := rand.Read(tx)
		if err != nil {
			panic(err)
		}
		txs[i] = tx
	}
	return txs
}

func generateRandomBlobWithNamespace(namespace namespace.ID, size int) types.Blob {
	blob := types.Blob{
		NamespaceID: namespace,
		Data:        tmrand.Bytes(size),
	}
	return blob
}
