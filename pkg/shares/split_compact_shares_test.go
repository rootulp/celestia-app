package shares

import (
	"bytes"
	"testing"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/stretchr/testify/assert"
	coretypes "github.com/tendermint/tendermint/types"
)

func TestCount(t *testing.T) {
	type testCase struct {
		transactions   []coretypes.Tx
		wantShareCount int
	}
	testCases := []testCase{
		{transactions: []coretypes.Tx{}, wantShareCount: 0},
		{transactions: []coretypes.Tx{[]byte{0}}, wantShareCount: 1},
		{transactions: []coretypes.Tx{bytes.Repeat([]byte{1}, 100)}, wantShareCount: 1},
		// Test with 1 byte over 1 share
		{transactions: []coretypes.Tx{bytes.Repeat([]byte{1}, rawTxSize(appconsts.FirstCompactShareContentSize+1))}, wantShareCount: 2},
		{transactions: []coretypes.Tx{generateTx(1)}, wantShareCount: 1},
		{transactions: []coretypes.Tx{generateTx(2)}, wantShareCount: 2},
		{transactions: []coretypes.Tx{generateTx(20)}, wantShareCount: 20},
	}
	for _, tc := range testCases {
		css := NewCompactShareSplitter(appconsts.TxNamespaceID, appconsts.ShareVersionZero)
		for _, transaction := range tc.transactions {
			css.WriteTx(transaction)
		}
		got := css.Count()
		if got != tc.wantShareCount {
			t.Errorf("count got %d want %d", got, tc.wantShareCount)
		}
	}
}

// generateTx generates a transaction that occupies exactly numShares number of
// shares.
func generateTx(numShares int) coretypes.Tx {
	if numShares == 0 {
		return coretypes.Tx{}
	}
	if numShares == 1 {
		return bytes.Repeat([]byte{1}, rawTxSize(appconsts.FirstCompactShareContentSize))
	}
	return bytes.Repeat([]byte{1}, rawTxSize(appconsts.FirstCompactShareContentSize+(numShares-1)*appconsts.ContinuationCompactShareContentSize))
}

// rawTxSize returns the raw tx size that can be used to construct a
// tx of desiredSize bytes. This function is useful in tests to account for
// the length delimiter that is prefixed to a tx.
func rawTxSize(desiredSize int) int {
	return desiredSize - DelimLen(uint64(desiredSize))
}

func TestExport(t *testing.T) {
	type testCase struct {
		name       string
		want       []Share
		writeBytes [][]byte
	}

	oneShare, _ := zeroPadIfNecessary([]byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
		0x1,                // info byte
		0x0, 0x0, 0x0, 0x1, // sequence len
		0x0, 0x0, 0x0, 17, // reserved bytes
		0xf, // data
	}, appconsts.ShareSize)

	firstShare, _ := zeroPadIfNecessary([]byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
		0x1,                // info byte
		0x0, 0x0, 0x2, 0x0, // sequence len
		0x0, 0x0, 0x0, 17, // reserved bytes
		0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, // data
	}, appconsts.ShareSize)

	continuationShare, _ := zeroPadIfNecessary([]byte{
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, // namespace id
		0x0,                // info byte
		0x0, 0x0, 0x0, 0x0, // reserved bytes
		0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, 0xf, // data
	}, appconsts.ShareSize)

	testCases := []testCase{
		{
			name: "empty",
			want: []Share{},
		},
		{
			name: "one share with small sequence len",
			want: []Share{
				oneShare,
			},
			writeBytes: [][]byte{{0xf}},
		},
		{
			name: "two shares with big sequence len",
			want: []Share{
				firstShare,
				continuationShare,
			},
			writeBytes: [][]byte{bytes.Repeat([]byte{0xf}, 512)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			css := NewCompactShareSplitter(appconsts.TxNamespaceID, appconsts.ShareVersionZero)
			for _, bytes := range tc.writeBytes {
				css.WriteBytes(bytes)
			}
			got := css.Export()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWriteAfterExportV2(t *testing.T) {
	type testCase struct {
		name    string
		txs     []coretypes.Tx
		wantLen int
	}
	testCases := []testCase{
		{
			name:    "one tx that occupies exactly one share",
			txs:     []coretypes.Tx{generateTx(1)},
			wantLen: 1,
		},
		{
			name:    "one tx that occupies exactly two shares",
			txs:     []coretypes.Tx{generateTx(2)},
			wantLen: 2,
		},
		{
			name:    "one tx that occupies exactly three shares",
			txs:     []coretypes.Tx{generateTx(3)},
			wantLen: 3,
		},
		{
			name: "two txs that occupy exactly two shares",
			txs: []coretypes.Tx{
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.FirstCompactShareContentSize)),
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.ContinuationCompactShareContentSize)),
			},
			wantLen: 2,
		},
		{
			name: "three txs that occupy exactly three shares",
			txs: []coretypes.Tx{
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.FirstCompactShareContentSize)),
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.ContinuationCompactShareContentSize)),
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.ContinuationCompactShareContentSize)),
			},
			wantLen: 3,
		},
		{
			name: "four txs that occupy three full shares and one partial share",
			txs: []coretypes.Tx{
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.FirstCompactShareContentSize)),
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.ContinuationCompactShareContentSize)),
				bytes.Repeat([]byte{0xf}, rawTxSize(appconsts.ContinuationCompactShareContentSize)),
				[]byte{0xf},
			},
			wantLen: 4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			css := NewCompactShareSplitter(appconsts.TxNamespaceID, appconsts.ShareVersionZero)

			for _, tx := range tc.txs {
				css.WriteTx(tx)
			}

			assert.Equal(t, tc.wantLen, css.Count())
			assert.Equal(t, tc.wantLen, len(css.Export()))
		})
	}
}
