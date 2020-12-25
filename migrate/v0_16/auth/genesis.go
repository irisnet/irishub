package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	//system default special address
	BurnedCoinsAccAddr       = sdk.AccAddress(crypto.AddressHash([]byte("burnedCoins")))
	GovDepositCoinsAccAddr   = sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	CommunityTaxCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("communityTaxCoins")))

	ServiceDepositCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceDepositedCoins")))
	ServiceRequestCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("serviceRequestCoins")))
	ServiceTaxCoinsAccAddr     = sdk.AccAddress(crypto.AddressHash([]byte("serviceTaxCoins")))

	HTLCLockedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("HTLCLockedCoins")))
)

type GenesisState struct {
	CollectedFees sdk.Coins `json:"collected_fee"`
	FeeAuth       FeeAuth   `json:"data"`
	Params        Params    `json:"params"`
}

type Params struct {
	GasPriceThreshold sdk.Int `json:"gas_price_threshold"`
	TxSizeLimit       uint64  `json:"tx_size"`
}

type FeeAuth struct {
	NativeFeeDenom string `json:"native_fee_denom"`
}
