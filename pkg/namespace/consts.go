package namespace

import (
	"bytes"
)

const (
	// NamespaveVersionSize is the size of a namespace version in bytes.
	NamespaceVersionSize = 1
	// NamespaceIDSize is the size of a namespace ID in bytes.
	NamespaceIDSize = 32
	// NamespaceSize is the size of a namespace in bytes.
	NamespaceSize = NamespaceIDSize + NamespaceVersionSize

	// NamespaceVersionZero is the first namespace version.
	NamespaceVersionZero = uint8(0)
)

var (
	VersionZeroPrefix = bytes.Repeat([]byte{0}, 22)

	// TxNamespaceID is the namespace reserved for transaction data.
	TxNamespaceID = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1})

	// IntermediateStateRootsNamespaceID is the namespace reserved for
	// intermediate state root data.
	IntermediateStateRootsNamespaceID = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 2})

	// EvidenceNamespaceID is the namespace reserved for evidence.
	EvidenceNamespaceID = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 3})

	// PayForBlobNamespaceID is the namespace reserved for PayForBlobs transactions.
	PayForBlobNamespaceID = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 4})

	// ReservedPaddingNamespaceID is the namespace used for padding after all
	// reserved namespaces. In practice this padding is after transactions
	// (ordinary and PFBs) but before blobs.
	ReservedPaddingNamespaceID = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 255})

	// MaxReservedNamespace is lexicographically the largest namespace that is
	// reserved for protocol use.
	MaxReservedNamespace = MustNewV0([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 255})

	// TailPaddingNamespaceID is the namespace reserved for tail padding. All data
	// with this namespace will be ignored.
	TailPaddingNamespaceID = MustNewV0([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE})

	// ParitySharesNamespaceID is the namespace reserved for erasure coded data.
	ParitySharesNamespaceID = MustNewV0([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
)
