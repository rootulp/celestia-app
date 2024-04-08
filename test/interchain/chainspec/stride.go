package chainspec

import (
	"encoding/json"
	"fmt"

	"github.com/icza/dyno"
	"github.com/strangelove-ventures/interchaintest/v6"
	"github.com/strangelove-ventures/interchaintest/v6/ibc"
)

const (
	trustingPeriod = "199s"
)

var Stride = &interchaintest.ChainSpec{
	Name: "stride",
	ChainConfig: ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "stride",
		ChainID: "stride-1",
		Images: []ibc.DockerImage{{
			Repository: "ghcr.io/strangelove-ventures/heighliner/stride",
			Version:    "v21.0.0",
			UidGid:     "1025:1025",
		}},
		Bin:            "strided",
		Bech32Prefix:   "stride",
		Denom:          "ustrd",
		GasPrices:      "0.01ustrd",
		TrustingPeriod: trustingPeriod,
		GasAdjustment:  1.1,
		ModifyGenesis:  ModifyGenesisStride(),
	},
	NumFullNodes:  numFullNodes(),
	NumValidators: numValidators(),
}

var AllowMessages = []string{
	"/cosmos.bank.v1beta1.MsgSend",
	"/cosmos.bank.v1beta1.MsgMultiSend",
	"/cosmos.staking.v1beta1.MsgDelegate",
	"/cosmos.staking.v1beta1.MsgUndelegate",
	"/cosmos.staking.v1beta1.MsgRedeemTokensforShares",
	"/cosmos.staking.v1beta1.MsgTokenizeShares",
	"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward",
	"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress",
	"/ibc.applications.transfer.v1.MsgTransfer",
}

const (
	StrideAdminAccount  = "admin"
	StrideAdminMnemonic = "tone cause tribe this switch near host damage idle fragile antique tail soda alien depth write wool they rapid unfold body scan pledge soft"
)

const (
	DayEpochIndex    = 1
	DayEpochLen      = "100s"
	StrideEpochIndex = 2
	StrideEpochLen   = "40s"
	IntervalLen      = 1
	VotingPeriod     = "30s"
	MaxDepositPeriod = "30s"
	UnbondingTime    = "200s"
	TrustingPeriod   = "199s"
)

func ModifyGenesisStride() func(ibc.ChainConfig, []byte) ([]byte, error) {
	return func(cfg ibc.ChainConfig, genbz []byte) ([]byte, error) {
		g := make(map[string]interface{})
		if err := json.Unmarshal(genbz, &g); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genesis file: %w", err)
		}

		if err := dyno.Set(g, DayEpochLen, "app_state", "epochs", "epochs", DayEpochIndex, "duration"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, StrideEpochLen, "app_state", "epochs", "epochs", StrideEpochIndex, "duration"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, UnbondingTime, "app_state", "staking", "params", "unbonding_time"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, IntervalLen, "app_state", "stakeibc", "params", "rewards_interval"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, IntervalLen, "app_state", "stakeibc", "params", "delegate_interval"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, IntervalLen, "app_state", "stakeibc", "params", "deposit_interval"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, IntervalLen, "app_state", "stakeibc", "params", "redemption_rate_interval"); err != nil {
			return nil, err
		}
		if err := dyno.Set(g, IntervalLen, "app_state", "stakeibc", "params", "reinvest_interval"); err != nil {
			return nil, err
		}
		out, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal genesis bytes to json: %w", err)
		}
		return out, nil
	}
}
