package namespace

import (
	tmrand "github.com/tendermint/tendermint/libs/rand"
)

func RandomBlobNamespace() Namespace {
	for {
		randomID := tmrand.Bytes(NamespaceIDSize - len(VersionZeroPrefix))
		namespace := MustNewV0(randomID)

		if namespace.IsReserved() || namespace.IsParityShares() || namespace.IsTailPadding() {
			continue
		}

		return namespace
	}
}

func RandomBlobNamespaces(count int) []Namespace {
	namespaces := make([]Namespace, count)
	for i := 0; i < count; i++ {
		namespaces[i] = RandomBlobNamespace()
	}
	return namespaces
}
