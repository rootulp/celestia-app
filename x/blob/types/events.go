package types

import (
	appns "github.com/celestiaorg/celestia-app/pkg/namespace"
	"github.com/cosmos/gogoproto/proto"
)

var EventTypePayForBlob = proto.MessageName(&EventPayForBlobs{})

// NewPayForBlobsEvent returns a new EventPayForBlobs
func NewPayForBlobsEvent(signer string, blobSizes []uint32, namespaces []appns.Namespace) *EventPayForBlobs {
	rawNamespaces := make([][]byte, len(namespaces))
	for _, ns := range namespaces {
		rawNamespaces = append(rawNamespaces, ns.Bytes())
	}
	return &EventPayForBlobs{
		Signer:     signer,
		BlobSizes:  blobSizes,
		Namespaces: rawNamespaces,
	}
}
