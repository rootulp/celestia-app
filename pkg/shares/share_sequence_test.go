package shares

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
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
				Namespace: appns.TxNamespace,
				Shares:    []Share{},
			},
			want:    []byte{},
			wantErr: false,
		},
		{
			name: "one empty share",
			shareSequence: ShareSequence{
				Namespace: appns.TxNamespace,
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
				Namespace: appns.TxNamespace,
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
				Namespace: appns.TxNamespace,
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
		{100000, 208},
	}
	for _, tc := range testCases {
		got := SparseSharesNeeded(tc.sequenceLen)
		assert.Equal(t, tc.want, got)
	}
}

func shareWithData(namespace appns.Namespace, isSequenceStart bool, sequenceLen uint32, data []byte) (rawShare Share) {
	infoByte, _ := NewInfoByte(appconsts.ShareVersionZero, isSequenceStart)
	rawShareBytes := make([]byte, 0, appconsts.ShareSize)
	rawShareBytes = append(rawShareBytes, namespace.Bytes()...)
	rawShareBytes = append(rawShareBytes, byte(infoByte))
	if isSequenceStart {
		sequenceLenBuf := make([]byte, appconsts.SequenceLenBytes)
		binary.BigEndian.PutUint32(sequenceLenBuf, sequenceLen)
		rawShareBytes = append(rawShareBytes, sequenceLenBuf...)
	}
	rawShareBytes = append(rawShareBytes, data...)

	return padShare(Share{data: rawShareBytes})
}

func TestJaved(t *testing.T) {
	base64Shares := []string{
		"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",
	}
	shares := []Share{}
	for _, share := range base64Shares {
		share, err := base64.StdEncoding.DecodeString(share)
		assert.NoError(t, err)
		shares = append(shares, share)
	}

	for _, s := range shares {
		fmt.Printf("%v\n", s)
	}

	// problemShares := shares[4:]

	// namespace := problemShares[0].NamespaceID()
	// sequence := ShareSequence{
	// 	NamespaceID: namespace,
	// 	Shares:      shares,
	// }

	// sequenceLen, err := sequence.SequenceLen()
	// assert.NoError(t, err)
	// fmt.Printf("sequence length %v\n", sequenceLen)

	// sparseSharesNeedd := SparseSharesNeeded(sequenceLen)
	// fmt.Printf("sparse shares needed %v\n", sparseSharesNeedd)

	// numSequenceShares := len(sequence.Shares)
	// fmt.Printf("num sequence shares %v\n", numSequenceShares)

	// for i := 0; i < numSequenceShares; i++ {
	// 	share := sequence.Shares[i]
	// 	isPadding, err := share.IsPadding()
	// 	assert.NoError(t, err)
	// 	fmt.Printf("share %v isPadding %v\n", i, isPadding)
	// }

	got, err := ParseShares(ToBytes(shares[3:]))
	assert.NoError(t, err)
	fmt.Printf("got %v\n", got)
}
