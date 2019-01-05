package protocol

import sdk "github.com/irisnet/irishub/types"

var (
	KeyMain          = sdk.NewKVStoreKey(sdk.MainStore)
	KeyAccount       = sdk.NewKVStoreKey("acc")
	KeyStake         = sdk.NewKVStoreKey("stake")
	TkeyStake        = sdk.NewTransientStoreKey("transient_stake")
	KeyMint          = sdk.NewKVStoreKey("mint")
	KeyDistr         = sdk.NewKVStoreKey("distr")
	TkeyDistr        = sdk.NewTransientStoreKey("transient_distr")
	KeySlashing      = sdk.NewKVStoreKey("slashing")
	KeyGov           = sdk.NewKVStoreKey("gov")
	KeyFee           = sdk.NewKVStoreKey("fee")
	KeyParams        = sdk.NewKVStoreKey("params")
	TkeyParams       = sdk.NewTransientStoreKey("transient_params")
	KeyService       = sdk.NewKVStoreKey("service")
	KeyGuardian      = sdk.NewKVStoreKey("guardian")
	KeyUpgrade       = sdk.NewKVStoreKey("upgrade")
)
