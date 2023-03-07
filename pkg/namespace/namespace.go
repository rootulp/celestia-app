package namespace

import (
	"bytes"
	"fmt"

	"github.com/celestiaorg/celestia-app/pkg/appconsts"
)

const (
	NamespaceSize = appconsts.NamespaceSize
	VersionZero   = uint8(0)
)

var VersionZeroPrefix = bytes.Repeat([]byte{0}, 22)

type Namespace struct {
	Version uint8
	ID      []byte
}

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
	if version != VersionZero {
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

	if version == VersionZero && !bytes.HasPrefix(id, VersionZeroPrefix) {
		return fmt.Errorf("unsupported namespace id %v must start with prefix %v", id, VersionZeroPrefix)
	}
	return nil
}
