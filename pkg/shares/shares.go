package shares

import (
	"encoding/binary"
	"fmt"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
)

// Share contains the raw share data (including namespace ID).
type Share []byte

func NewShare(data []byte) (Share, error) {
	if len(data) != appconsts.ShareSize {
		return nil, fmt.Errorf("share data must be %d bytes, got %d", appconsts.ShareSize, len(data))
	}
	return Share(data), nil
}

func (s Share) Namespace() (appns.Namespace, error) {
	if len(s) < appns.NamespaceSize {
		panic(fmt.Sprintf("share %s is too short to contain a namespace", s))
	}
	return appns.From(s[:appns.NamespaceSize])
}

func (s Share) InfoByte() (InfoByte, error) {
	if len(s) < appns.NamespaceSize+appconsts.ShareInfoBytes {
		return 0, fmt.Errorf("share %s is too short to contain an info byte", s)
	}
	// the info byte is the first byte after the namespace
	unparsed := s[appns.NamespaceSize]
	return ParseInfoByte(unparsed)
}

func (s Share) Version() (uint8, error) {
	infoByte, err := s.InfoByte()
	if err != nil {
		return 0, err
	}
	return infoByte.Version(), nil
}

// IsSequenceStart returns true if this is the first share in a sequence.
func (s Share) IsSequenceStart() (bool, error) {
	infoByte, err := s.InfoByte()
	if err != nil {
		return false, err
	}
	return infoByte.IsSequenceStart(), nil
}

// IsCompactShare returns true if this is a compact share.
func (s Share) IsCompactShare() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	isCompact := ns.IsTx() || ns.IsPayForBlob()
	return isCompact, nil
}

// SequenceLen returns the sequence length of this share and optionally an
// error. It returns 0, nil if this is a continuation share (i.e. doesn't
// contain a sequence length).
func (s Share) SequenceLen() (sequenceLen uint32, err error) {
	isSequenceStart, err := s.IsSequenceStart()
	if err != nil {
		return 0, err
	}
	if !isSequenceStart {
		return 0, nil
	}

	start := appconsts.NamespaceSize + appconsts.ShareInfoBytes
	end := start + appconsts.SequenceLenBytes
	if len(s) < end {
		return 0, fmt.Errorf("share %s is too short to contain a sequence length", s)
	}
	return binary.BigEndian.Uint32(s[start:end]), nil
}

// IsPadding returns whether this share is padding or not.
func (s Share) IsPadding() (bool, error) {
	isNamespacePadding, err := s.isNamespacePadding()
	if err != nil {
		return false, err
	}
	isTailPadding, err := s.isTailPadding()
	if err != nil {
		return false, err
	}
	isReservedPadding, err := s.isReservedPadding()
	if err != nil {
		return false, err
	}
	return isNamespacePadding || isTailPadding || isReservedPadding, nil
}

func (s Share) isNamespacePadding() (bool, error) {
	isSequenceStart, err := s.IsSequenceStart()
	if err != nil {
		return false, err
	}
	sequenceLen, err := s.SequenceLen()
	if err != nil {
		return false, err
	}

	return isSequenceStart && sequenceLen == 0, nil
}

func (s Share) isTailPadding() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	return ns.IsTailPadding(), nil
}

func (s Share) isReservedPadding() (bool, error) {
	ns, err := s.Namespace()
	if err != nil {
		return false, err
	}
	return ns.IsReservedPadding(), nil
}

func (s Share) ToBytes() []byte {
	return []byte(s)
}

// RawData returns the raw share data. The raw share data does not contain the
// namespace ID, info byte, sequence length, or reserved bytes.
func (s Share) RawData() (rawData []byte, err error) {
	if len(s) < s.rawDataStartIndex() {
		return rawData, fmt.Errorf("share %s is too short to contain raw data", s)
	}

	return s[s.rawDataStartIndex():], nil
}

func (s Share) rawDataStartIndex() int {
	isStart, err := s.IsSequenceStart()
	if err != nil {
		panic(err)
	}
	isCompact, err := s.IsCompactShare()
	if err != nil {
		panic(err)
	}
	if isStart && isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.SequenceLenBytes + appconsts.CompactShareReservedBytes
	} else if isStart && !isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.SequenceLenBytes
	} else if !isStart && isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes + appconsts.CompactShareReservedBytes
	} else if !isStart && !isCompact {
		return appconsts.NamespaceSize + appconsts.ShareInfoBytes
	} else {
		panic(fmt.Sprintf("unable to determine the rawDataStartIndex for share %s", s))
	}
}

func ToBytes(shares []Share) (bytes [][]byte) {
	bytes = make([][]byte, len(shares))
	for i, share := range shares {
		bytes[i] = []byte(share)
	}
	return bytes
}

func FromBytes(bytes [][]byte) (shares []Share) {
	shares = make([]Share, len(bytes))
	for i, b := range bytes {
		shares[i] = Share(b)
	}
	return shares
}
