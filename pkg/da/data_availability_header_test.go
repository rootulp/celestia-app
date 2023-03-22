package da

import (
	"bytes"
	"strings"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilDataAvailabilityHeaderHashDoesntCrash(t *testing.T) {
	// This follows RFC-6962, i.e. `echo -n '' | sha256sum`
	emptyBytes := []byte{
		0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8,
		0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b,
		0x78, 0x52, 0xb8, 0x55,
	}

	assert.Equal(t, emptyBytes, (*DataAvailabilityHeader)(nil).Hash())
	assert.Equal(t, emptyBytes, new(DataAvailabilityHeader).Hash())
}

func TestMinDataAvailabilityHeader(t *testing.T) {
	dah := MinDataAvailabilityHeader()
	expectedHash := []byte{0xad, 0x23, 0x6a, 0x2e, 0x4c, 0x5f, 0xca, 0x6c, 0xdb, 0xae, 0x5d, 0x5e, 0xdf, 0x79, 0xe8, 0x8e, 0x84, 0xc5, 0x2e, 0xed, 0x62, 0xeb, 0xd0, 0xb6, 0x5d, 0x18, 0xb2, 0x7c, 0x32, 0xa8, 0xbc, 0x58}
	require.Equal(t, expectedHash, dah.hash)
	require.NoError(t, dah.ValidateBasic())
}

func TestNewDataAvailabilityHeader(t *testing.T) {
	type test struct {
		name         string
		expectedHash []byte
		squareSize   uint64
		shares       [][]byte
	}

	tests := []test{
		{
			name: "typical",
			expectedHash: []byte{
				0x71, 0x3d, 0x9, 0x9c, 0x2e, 0xd1, 0xfe, 0xed, 0x64, 0x8d, 0xb0, 0x6f, 0xb0, 0xf2, 0x4b, 0xe,
				0xcd, 0x86, 0x37, 0x53, 0xb5, 0x40, 0x6c, 0x72, 0x3d, 0xf5, 0xa7, 0xe2, 0x90, 0xb4, 0x70, 0x32,
			},
			squareSize: 2,
			shares:     generateShares(4, 1),
		},
		{
			name: "max square size",
			expectedHash: []byte{
				0xe2, 0x47, 0xdc, 0x2b, 0xa4, 0xac, 0x57, 0xb8, 0x27, 0xc6, 0xb5, 0xc0, 0x2c, 0x7a, 0xed, 0xfb,
				0x30, 0x25, 0xe8, 0xa7, 0x8b, 0xde, 0x75, 0xb4, 0xc0, 0xd4, 0xaf, 0xe8, 0x10, 0xa1, 0xd9, 0xc4,
			},
			squareSize: appconsts.DefaultMaxSquareSize,
			shares:     generateShares(appconsts.DefaultMaxSquareSize*appconsts.DefaultMaxSquareSize, 99),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eds, err := ExtendShares(tt.squareSize, tt.shares)
			require.NoError(t, err)
			resdah := NewDataAvailabilityHeader(eds)
			require.Equal(t, tt.squareSize*2, uint64(len(resdah.ColumnRoots)), tt.name)
			require.Equal(t, tt.squareSize*2, uint64(len(resdah.RowsRoots)), tt.name)
			require.Equal(t, tt.expectedHash, resdah.hash, tt.name)
		})
	}
}

func TestExtendShares(t *testing.T) {
	type test struct {
		name        string
		expectedErr bool
		squareSize  uint64
		shares      [][]byte
	}

	tests := []test{
		{
			name:        "too large square size",
			expectedErr: true,
			squareSize:  appconsts.DefaultMaxSquareSize + 1,
			shares:      generateShares((appconsts.DefaultMaxSquareSize+1)*(appconsts.DefaultMaxSquareSize+1), 1),
		},
		{
			name:        "invalid number of shares",
			expectedErr: true,
			squareSize:  2,
			shares:      generateShares(5, 1),
		},
	}

	for _, tt := range tests {
		tt := tt
		eds, err := ExtendShares(tt.squareSize, tt.shares)
		if tt.expectedErr {
			require.NotNil(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, tt.squareSize*2, eds.Width(), tt.name)
	}
}

func TestDataAvailabilityHeaderProtoConversion(t *testing.T) {
	type test struct {
		name string
		dah  DataAvailabilityHeader
	}

	shares := generateShares(appconsts.DefaultMaxSquareSize*appconsts.DefaultMaxSquareSize, 1)
	eds, err := ExtendShares(appconsts.DefaultMaxSquareSize, shares)
	require.NoError(t, err)
	bigdah := NewDataAvailabilityHeader(eds)

	tests := []test{
		{
			name: "min",
			dah:  MinDataAvailabilityHeader(),
		},
		{
			name: "max",
			dah:  bigdah,
		},
	}

	for _, tt := range tests {
		tt := tt
		pdah, err := tt.dah.ToProto()
		require.NoError(t, err)
		resDah, err := DataAvailabilityHeaderFromProto(pdah)
		require.NoError(t, err)
		resDah.Hash() // calc the hash to make the comparisons fair
		require.Equal(t, tt.dah, *resDah, tt.name)
	}
}

func Test_DAHValidateBasic(t *testing.T) {
	type test struct {
		name      string
		dah       DataAvailabilityHeader
		expectErr bool
		errStr    string
	}

	shares := generateShares(appconsts.DefaultMaxSquareSize*appconsts.DefaultMaxSquareSize, 1)
	eds, err := ExtendShares(appconsts.DefaultMaxSquareSize, shares)
	require.NoError(t, err)
	bigdah := NewDataAvailabilityHeader(eds)

	// make a mutant dah that has too many roots
	var tooBigDah DataAvailabilityHeader
	tooBigDah.ColumnRoots = make([][]byte, appconsts.DefaultMaxSquareSize*appconsts.DefaultMaxSquareSize)
	tooBigDah.RowsRoots = make([][]byte, appconsts.DefaultMaxSquareSize*appconsts.DefaultMaxSquareSize)
	copy(tooBigDah.ColumnRoots, bigdah.ColumnRoots)
	copy(tooBigDah.RowsRoots, bigdah.RowsRoots)
	tooBigDah.ColumnRoots = append(tooBigDah.ColumnRoots, bytes.Repeat([]byte{1}, 32))
	tooBigDah.RowsRoots = append(tooBigDah.RowsRoots, bytes.Repeat([]byte{1}, 32))
	// make a mutant dah that has too few roots
	var tooSmallDah DataAvailabilityHeader
	tooSmallDah.ColumnRoots = [][]byte{bytes.Repeat([]byte{2}, 32)}
	tooSmallDah.RowsRoots = [][]byte{bytes.Repeat([]byte{2}, 32)}
	// use a bad hash
	badHashDah := MinDataAvailabilityHeader()
	badHashDah.hash = []byte{1, 2, 3, 4}
	// dah with not equal number of roots
	mismatchDah := MinDataAvailabilityHeader()
	mismatchDah.ColumnRoots = append(mismatchDah.ColumnRoots, bytes.Repeat([]byte{2}, 32))

	tests := []test{
		{
			name: "min",
			dah:  MinDataAvailabilityHeader(),
		},
		{
			name: "max",
			dah:  bigdah,
		},
		{
			name:      "too big dah",
			dah:       tooBigDah,
			expectErr: true,
			errStr:    "maximum valid DataAvailabilityHeader has at most",
		},
		{
			name:      "too small dah",
			dah:       tooSmallDah,
			expectErr: true,
			errStr:    "minimum valid DataAvailabilityHeader has at least",
		},
		{
			name:      "bash hash",
			dah:       badHashDah,
			expectErr: true,
			errStr:    "wrong hash",
		},
		{
			name:      "mismatched roots",
			dah:       mismatchDah,
			expectErr: true,
			errStr:    "unequal number of row and column roots",
		},
	}

	for _, tt := range tests {
		tt := tt
		err := tt.dah.ValidateBasic()
		if tt.expectErr {
			require.True(t, strings.Contains(err.Error(), tt.errStr), tt.name)
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
	}
}

func generateShares(count int, repeatByte byte) [][]byte {
	shares := make([][]byte, count)
	for i := 0; i < count; i++ {
		shares[i] = bytes.Repeat([]byte{repeatByte}, appconsts.ShareSize)
	}
	return shares
}
