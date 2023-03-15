package namespace

import (
	"bytes"
	"fmt"
)

type Namespace struct {
	Version uint8
	ID      []byte
}

// New returns a new namespace with the provided version and id.
func New(version uint8, id []byte) (Namespace, error) {
	err := validateVersion(version)
	if err != nil {
		return Namespace{}, err
	}

	err = validateID(version, id)
	if err != nil {
		return Namespace{}, err
	}

	return Namespace{
		Version: version,
		ID:      id,
	}, nil
}

// MustNew returns a new namespace with the provided version and id. It panics
// if the provided version or id are not supported.
func MustNew(version uint8, id []byte) Namespace {
	ns, err := New(version, id)
	if err != nil {
		panic(err)
	}
	return ns
}

// From returns a namespace from the provided byte slice.
func From(b []byte) (Namespace, error) {
	if len(b) != NamespaceSize+1 {
		return Namespace{}, fmt.Errorf("invalid namespace length: %v must be %v", len(b), NamespaceSize+1)
	}
	rawVersion := b[0]
	rawNamespace := b[1:]
	return New(rawVersion, rawNamespace)
}

// Bytes returns this namespace as a byte slice.
func (n Namespace) Bytes() []byte {
	return append([]byte{n.Version}, n.ID...)
}

// validateVersion returns an error if the version is not supported.
func validateVersion(version uint8) error {
	if version != NamespaceVersionZero {
		return fmt.Errorf("unsupported namespace version %v", version)
	}
	return nil
}

// validateID returns an error if the provided id does not meet the requirements
// for the provided version.
func validateID(version uint8, id []byte) error {
	if len(id) != NamespaceSize {
		return fmt.Errorf("unsupported namespace id length: id %v must be %v bytes ", id, NamespaceSize)
	}

	if version == NamespaceVersionZero && !bytes.HasPrefix(id, VersionZeroPrefix) {
		return fmt.Errorf("unsupported namespace id %v must start with prefix %v", id, VersionZeroPrefix)
	}
	return nil
}

func (n Namespace) IsReserved() bool {
	return bytes.Compare(n.Bytes(), MaxReservedNamespace.Bytes()) < 1
}

func (n Namespace) IsParityShares() bool {
	return bytes.Equal(n.Bytes(), ParitySharesNamespaceID.Bytes())
}

func (n Namespace) IsTailPadding() bool {
	return bytes.Equal(n.Bytes(), TailPaddingNamespaceID.Bytes())
}

func (n Namespace) IsReservedPadding() bool {
	return bytes.Equal(n.Bytes(), ReservedPaddingNamespaceID.Bytes())
}

func (n Namespace) IsTx() bool {
	return bytes.Equal(n.Bytes(), TxNamespaceID.Bytes())
}

func (n Namespace) IsPayForBlob() bool {
	return bytes.Equal(n.Bytes(), PayForBlobNamespaceID.Bytes())
}
