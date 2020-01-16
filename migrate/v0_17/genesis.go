package v0_17

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/migrate/v0_17/asset"
	"github.com/irisnet/irishub/migrate/v0_17/auth"
	"github.com/irisnet/irishub/migrate/v0_17/coinswap"
	"github.com/irisnet/irishub/migrate/v0_17/distribution"
	"github.com/irisnet/irishub/migrate/v0_17/gov"
	"github.com/irisnet/irishub/migrate/v0_17/guardian"
	"github.com/irisnet/irishub/migrate/v0_17/htlc"
	"github.com/irisnet/irishub/migrate/v0_17/mint"
	"github.com/irisnet/irishub/migrate/v0_17/rand"
	"github.com/irisnet/irishub/migrate/v0_17/service"
	"github.com/irisnet/irishub/migrate/v0_17/slashing"
	"github.com/irisnet/irishub/migrate/v0_17/stake"
	"github.com/irisnet/irishub/migrate/v0_17/upgrade"
)

type GenesisState struct {
	Accounts     []GenesisAccount          `json:"accounts"`
	AuthData     auth.GenesisState         `json:"auth"`
	StakeData    stake.GenesisState        `json:"stake"`
	MintData     mint.GenesisState         `json:"mint"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	UpgradeData  upgrade.GenesisState      `json:"upgrade"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	ServiceData  service.GenesisState      `json:"service"`
	GuardianData guardian.GenesisState     `json:"guardian"`
	AssetData    asset.GenesisState        `json:"asset"`
	RandData     rand.GenesisState         `json:"rand"`
	SwapData     coinswap.GenesisState     `json:"swap"`
	HtlcData     htlc.GenesisState         `json:"htlc"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
}

type GenesisFileState struct {
	Accounts     []GenesisFileAccount      `json:"accounts"`
	AuthData     auth.GenesisState         `json:"auth"`
	StakeData    stake.GenesisState        `json:"stake"`
	MintData     mint.GenesisState         `json:"mint"`
	DistrData    distribution.GenesisState `json:"distr"`
	GovData      gov.GenesisState          `json:"gov"`
	UpgradeData  upgrade.GenesisState      `json:"upgrade"`
	SlashingData slashing.GenesisState     `json:"slashing"`
	ServiceData  service.GenesisState      `json:"service"`
	GuardianData guardian.GenesisState     `json:"guardian"`
	AssetData    asset.GenesisState        `json:"asset"`
	RandData     rand.GenesisState         `json:"rand"`
	SwapData     coinswap.GenesisState     `json:"swap"`
	HtlcData     htlc.GenesisState         `json:"htlc"`
	GenTxs       []json.RawMessage         `json:"gentxs"`
}

type GenesisFileAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         []string       `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`
}

func convertToGenesisState(genesisFileState GenesisFileState) GenesisState {
	var genesisAccounts []GenesisAccount
	for _, gacc := range genesisFileState.Accounts {
		var coins sdk.Coins
		for _, coinStr := range gacc.Coins {
			coin, err := sdk.ParseCoin(coinStr)
			if err != nil {
				panic(err)
			}
			coins = append(coins, coin)
		}

		acc := GenesisAccount{
			Address:       gacc.Address,
			Coins:         coins,
			AccountNumber: gacc.AccountNumber,
			Sequence:      gacc.Sequence,
		}
		genesisAccounts = append(genesisAccounts, acc)
	}

	return GenesisState{
		Accounts:     genesisAccounts,
		AuthData:     genesisFileState.AuthData,
		StakeData:    genesisFileState.StakeData,
		MintData:     genesisFileState.MintData,
		DistrData:    genesisFileState.DistrData,
		GovData:      genesisFileState.GovData,
		UpgradeData:  genesisFileState.UpgradeData,
		SlashingData: genesisFileState.SlashingData,
		ServiceData:  genesisFileState.ServiceData,
		GuardianData: genesisFileState.GuardianData,
		AssetData:    genesisFileState.AssetData,
		RandData:     genesisFileState.RandData,
		GenTxs:       genesisFileState.GenTxs,
		SwapData:     genesisFileState.SwapData,
		HtlcData:     genesisFileState.HtlcData,
	}
}
