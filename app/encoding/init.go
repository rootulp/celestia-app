package encoding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const AccountAddressPrefix = "celestia"

var (
	Bech32PrefixAccAddr  = AccountAddressPrefix
	Bech32PrefixAccPub   = AccountAddressPrefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = AccountAddressPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func init() {
	// cfg := sdk.GetConfig()
	// cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	// cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	// cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	// cfg.Seal()
}
