package namespace

import (
	"bytes"

	tmrand "github.com/tendermint/tendermint/libs/rand"
)

func RandomBlobNamespace() Namespace {
	for {
		randomID := tmrand.Bytes(NamespaceIDSize)
		randomNs := MustNew(NamespaceVersionZero, randomID)

		isReservedNS := bytes.Compare(randomNs.Bytes(), MaxReservedNamespace.Bytes()) <= 0
		isParityNS := bytes.Equal(randomNs.Bytes(), ParitySharesNamespaceID.Bytes())
		isTailPaddingNS := bytes.Equal(randomNs.Bytes(), TailPaddingNamespaceID.Bytes())
		if isReservedNS || isParityNS || isTailPaddingNS {
			continue
		}

		return randomNs
	}
}

func RandomBlobNamespaces(count int) [][]byte {
	namespaces := make([][]byte, count)
	for i := 0; i < count; i++ {
		namespaces[i] = RandomBlobNamespace().Bytes()
	}
	return namespaces
}
