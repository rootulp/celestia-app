package shares

import (
	"bytes"
	"fmt"
	"sort"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/testutil/testfactory"
	"github.com/celestiaorg/nmt/namespace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coretypes "github.com/tendermint/tendermint/types"
)

// TODO: move this test to a file named parse_sparse_shares_test.go
func Test_parseSparseShares(t *testing.T) {
	type test struct {
		name      string
		blobSize  int
		blobCount int
	}

	// each test is ran twice, once using blobSize as an exact size, and again
	// using it as a cap for randomly sized leaves
	tests := []test{
		{
			name:      "single small blob",
			blobSize:  10,
			blobCount: 1,
		},
		{
			name:      "ten small blobs",
			blobSize:  10,
			blobCount: 10,
		},
		{
			name:      "single big blob",
			blobSize:  appconsts.ContinuationSparseShareContentSize * 4,
			blobCount: 1,
		},
		// {
		// 	name:      "many big blobs",
		// 	blobSize:  appconsts.ContinuationSparseShareContentSize * 4,
		// 	blobCount: 10,
		// },
		// {
		// 	name:      "single exact size blob",
		// 	blobSize:  appconsts.FirstSparseShareContentSize,
		// 	blobCount: 1,
		// },
		// {
		// 	name:      "many exact size blobs",
		// 	blobSize:  appconsts.ContinuationSparseShareContentSize,
		// 	blobCount: 10,
		// },
	}

	for _, tc := range tests {
		// run the tests with identically sized blobs
		t.Run(fmt.Sprintf("%s identically sized ", tc.name), func(t *testing.T) {
			blobs := make([]coretypes.Blob, tc.blobCount)
			for i := 0; i < tc.blobCount; i++ {
				blobs[i] = testfactory.GenerateRandomBlob(tc.blobSize)
			}

			sort.Sort(coretypes.BlobsByNamespace(blobs))

			shares, _ := SplitBlobs(0, nil, blobs, false)
			rawShares := ToBytes(shares)

			parsedBlobs, err := parseSparseShares(rawShares, appconsts.SupportedShareVersions)
			if err != nil {
				t.Error(err)
			}

			// check that the namespaces and data are the same
			for i := 0; i < len(blobs); i++ {
				assert.Equal(t, blobs[i].NamespaceID, parsedBlobs[i].NamespaceID, "parsed blob namespace does not match")
				assert.Equal(t, blobs[i].Data, parsedBlobs[i].Data, "parsed blob data does not match")
			}
		})

		// run the same tests using randomly sized blobs with caps of tc.blobSize
		// t.Run(fmt.Sprintf("%s randomly sized", tc.name), func(t *testing.T) {
		// 	blobs := testfactory.GenerateRandomlySizedBlobs(tc.blobCount, tc.blobSize)
		// 	shares, _ := SplitBlobs(0, nil, blobs, false)
		// 	rawShares := make([][]byte, len(shares))
		// 	for i, share := range shares {
		// 		rawShares[i] = []byte(share)
		// 	}

		// 	parsedBlobs, err := parseSparseShares(rawShares, appconsts.SupportedShareVersions)
		// 	if err != nil {
		// 		t.Error(err)
		// 	}

		// 	// check that the namespaces and data are the same
		// 	for i := 0; i < len(blobs); i++ {
		// 		assert.Equal(t, blobs[i].NamespaceID, parsedBlobs[i].NamespaceID)
		// 		assert.Equal(t, blobs[i].Data, parsedBlobs[i].Data)
		// 	}
		// })
	}
}

func Test_parseSparseSharesErrors(t *testing.T) {
	type testCase struct {
		name      string
		rawShares [][]byte
	}

	unsupportedShareVersion := 5
	infoByte, _ := NewInfoByte(uint8(unsupportedShareVersion), true)

	rawShare := []byte{}
	rawShare = append(rawShare, namespace.ID{1, 1, 1, 1, 1, 1, 1, 1}...)
	rawShare = append(rawShare, byte(infoByte))
	rawShare = append(rawShare, bytes.Repeat([]byte{0}, appconsts.ShareSize-len(rawShare))...)

	tests := []testCase{
		{
			"share with unsupported share version",
			[][]byte{rawShare},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			_, err := parseSparseShares(tt.rawShares, appconsts.SupportedShareVersions)
			assert.Error(t, err)
		})
	}
}

func TestParsePaddedBlob(t *testing.T) {
	sss := NewSparseShareSplitter()
	randomSmallBlob := testfactory.GenerateRandomBlob(appconsts.ContinuationSparseShareContentSize / 2)
	randomLargeBlob := testfactory.GenerateRandomBlob(appconsts.ContinuationSparseShareContentSize * 4)
	blobs := []coretypes.Blob{
		randomSmallBlob,
		randomLargeBlob,
	}
	sort.Sort(coretypes.BlobsByNamespace(blobs))
	err := sss.Write(blobs[0])
	assert.NoError(t, err)
	sss.WriteNamespacedPaddedShares(4)
	err = sss.Write(blobs[1])
	assert.NoError(t, err)
	sss.WriteNamespacedPaddedShares(10)
	shares := sss.Export()
	rawShares := ToBytes(shares)
	pblobs, err := parseSparseShares(rawShares, appconsts.SupportedShareVersions)
	require.NoError(t, err)
	require.Equal(t, blobs, pblobs)
}

func TestSparseShareContainsInfoByte(t *testing.T) {
	blob := generateRandomBlobOfShareCount(4)

	sequenceStartInfoByte, err := NewInfoByte(appconsts.ShareVersionZero, true)
	require.NoError(t, err)

	sequenceContinuationInfoByte, err := NewInfoByte(appconsts.ShareVersionZero, false)
	require.NoError(t, err)

	type testCase struct {
		name       string
		shareIndex int
		expected   InfoByte
	}
	testCases := []testCase{
		{
			name:       "first share of blob",
			shareIndex: 0,
			expected:   sequenceStartInfoByte,
		},
		{
			name:       "second share of blob",
			shareIndex: 1,
			expected:   sequenceContinuationInfoByte,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sss := NewSparseShareSplitter()
			err := sss.Write(blob)
			assert.NoError(t, err)
			shares := sss.Export()
			got, err := shares[tc.shareIndex].InfoByte()
			require.NoError(t, err)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestSparseShareSplitterCount(t *testing.T) {
	type testCase struct {
		name     string
		blob     coretypes.Blob
		expected int
	}
	testCases := []testCase{
		{
			name:     "one share",
			blob:     generateRandomBlobOfShareCount(1),
			expected: 1,
		},
		{
			name:     "two shares",
			blob:     generateRandomBlobOfShareCount(2),
			expected: 2,
		},
		{
			name:     "ten shares",
			blob:     generateRandomBlobOfShareCount(10),
			expected: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sss := NewSparseShareSplitter()
			err := sss.Write(tc.blob)
			assert.NoError(t, err)
			got := sss.Count()
			assert.Equal(t, tc.expected, got)
		})
	}
}
