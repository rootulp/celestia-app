package shares

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/stretchr/testify/assert"
)

func TestShareSequenceRawData(t *testing.T) {
	type testCase struct {
		name          string
		shareSequence ShareSequence
		want          []byte
		wantErr       bool
	}
	blobNamespace := appns.RandomBlobNamespace()

	testCases := []testCase{
		{
			name: "empty share sequence",
			shareSequence: ShareSequence{
				Namespace: appns.TxNamespaceID,
				Shares:    []Share{},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "one empty share",
			shareSequence: ShareSequence{
				Namespace: appns.TxNamespaceID,
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
				Namespace: appns.TxNamespaceID,
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
				Namespace: appns.TxNamespaceID,
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
		{1000, 3},
		{10000, 21},
		{100000, 210},
	}
	for _, tc := range testCases {
		got := SparseSharesNeeded(tc.sequenceLen)
		assert.Equal(t, tc.want, got)
	}
}

func shareWithData(namespace appns.Namespace, isSequenceStart bool, sequenceLen uint32, data []byte) (rawShare []byte) {
	infoByte, _ := NewInfoByte(appconsts.ShareVersionZero, isSequenceStart)
	rawShare = append(rawShare, namespace.Bytes()...)
	rawShare = append(rawShare, byte(infoByte))
	if isSequenceStart {
		sequenceLenBuf := make([]byte, appconsts.SequenceLenBytes)
		binary.BigEndian.PutUint32(sequenceLenBuf, sequenceLen)
		rawShare = append(rawShare, sequenceLenBuf...)
	}
	rawShare = append(rawShare, data...)

	return padShare(rawShare)
}
