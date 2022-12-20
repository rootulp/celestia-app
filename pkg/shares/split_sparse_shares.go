package shares

import (
	"encoding/binary"
	"fmt"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/nmt/namespace"
	coretypes "github.com/tendermint/tendermint/types"
	"golang.org/x/exp/slices"
)

// SparseShareSplitter lazily splits blobs into shares that will eventually be
// included in a data square. It also has methods to help progressively count
// how many shares the blobs written take up.
type SparseShareSplitter struct {
	shares []Share
}

func NewSparseShareSplitter() *SparseShareSplitter {
	return &SparseShareSplitter{}
}

// Write writes the provided blob to this sparse share splitter. It returns an
// error or nil if no error is encountered.
func (sss *SparseShareSplitter) Write(blob coretypes.Blob) error {
	if !slices.Contains(appconsts.SupportedShareVersions, blob.ShareVersion) {
		return fmt.Errorf("unsupported share version: %d", blob.ShareVersion)
	}
	shares, err := split(blob)
	if err != nil {
		return err
	}
	sss.shares = append(sss.shares, shares...)
	return nil
}

// RemoveBlob will remove a blob from the underlying blob state. If
// there is namespaced padding after the blob, then that is also removed.
func (sss *SparseShareSplitter) RemoveBlob(i int) (int, error) {
	j := 1
	initialCount := len(sss.shares)
	if len(sss.shares) > i+1 {
		sequenceLen, err := sss.shares[i+1].SequenceLen()
		if err != nil {
			return 0, err
		}
		// 0 means that there is padding after the share that we are about to
		// remove. to remove this padding, we increase j by 1
		// with the blob
		if sequenceLen == 0 {
			j++
		}
	}
	copy(sss.shares[i:], sss.shares[i+j:])
	sss.shares = sss.shares[:len(sss.shares)-j]
	newCount := len(sss.shares)
	return initialCount - newCount, nil
}

// WriteNamespacedPaddedShares adds empty shares using the namespace of the last written share.
// This is useful to follow the message layout rules. It assumes that at least
// one share has already been written, if not it panics.
func (sss *SparseShareSplitter) WriteNamespacedPaddedShares(count int) {
	if len(sss.shares) == 0 {
		panic("cannot write empty namespaced shares on an empty SparseShareSplitter")
	}
	if count < 0 {
		panic("cannot write negative namespaced shares")
	}
	if count == 0 {
		return
	}
	lastBlob := sss.shares[len(sss.shares)-1]
	sss.shares = append(sss.shares, namespacedPaddedShares(lastBlob.NamespaceID(), count)...)
}

// Export returns the underlying shares written to this sparse share splitter.
func (sss *SparseShareSplitter) Export() []Share {
	return sss.shares
}

// Count returns the current number of shares that will be made if exporting.
func (sss *SparseShareSplitter) Count() int {
	return len(sss.shares)
}

func split(blob coretypes.Blob) (shares []Share, err error) {
	if len(blob.Data) == 0 {
		return shares, nil
	}
	dataToSplit := blob.Data
	for len(dataToSplit) > 0 {
		// the first share is special because it contains the sequence length
		if len(shares) == 0 {
			sequenceLen := blobSequenceLen(blob)
			dataToWrite := dataToSplit[:appconsts.FirstSparseShareContentSize]
			share, err := BuildShare(blob.NamespaceID, blob.ShareVersion, true, sequenceLen, dataToWrite)
			if err != nil {
				return shares, err
			}
			dataToSplit = dataToSplit[appconsts.FirstSparseShareContentSize:]
			shares = append(shares, share)
		} else {
			// all other shares are just the data
			dataToWrite := dataToSplit[:appconsts.ContinuationSparseShareContentSize]
			share, err := BuildShare(blob.NamespaceID, blob.ShareVersion, true, nil, dataToWrite)
			if err != nil {
				return shares, err
			}
			dataToSplit = dataToSplit[appconsts.ContinuationSparseShareContentSize:]
			shares = append(shares, share)
		}
	}
	return shares, nil
}

// blobSequenceLen returns a byte slice of appconsts.SequenceLenBytes with a big
// endian encoded sequence length of the provided blob.
func blobSequenceLen(blob coretypes.Blob) []byte {
	buf := make([]byte, appconsts.SequenceLenBytes)
	len := uint32(len(blob.Data))
	binary.BigEndian.PutUint32(buf, len)
	return buf
}

func namespacedPaddedShares(ns namespace.ID, count int) []Share {
	shares := make([]Share, count)
	for i := 0; i < count; i++ {
		shares[i] = namespacedPaddedShare(ns)
	}
	return shares
}

func namespacedPaddedShare(ns namespace.ID) Share {
	infoByte, err := NewInfoByte(appconsts.ShareVersionZero, true)
	if err != nil {
		panic(err)
	}

	sequenceLen := make([]byte, appconsts.SequenceLenBytes)
	binary.BigEndian.PutUint32(sequenceLen, uint32(0))

	share := make([]byte, 0, appconsts.ShareSize)
	share = append(share, ns...)
	share = append(share, byte(infoByte))
	share = append(share, sequenceLen...)
	share = append(share, appconsts.NameSpacedPaddedShareBytes...)
	return share
}
