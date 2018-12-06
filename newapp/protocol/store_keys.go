package protocol

import sdk "github.com/irisnet/irishub/types"

var (
	keyMain          = sdk.NewKVStoreKey("main")
	keyAccount       = sdk.NewKVStoreKey("acc")
	keyStake         = sdk.NewKVStoreKey("stake")
	tkeyStake        = sdk.NewTransientStoreKey("transient_stake")
	keyMint          = sdk.NewKVStoreKey("mint")
	keyDistr         = sdk.NewKVStoreKey("distr")
	tkeyDistr        = sdk.NewTransientStoreKey("transient_distr")
	keySlashing      = sdk.NewKVStoreKey("slashing")
	keyGov           = sdk.NewKVStoreKey("gov")
	keyRecord        = sdk.NewKVStoreKey("record")
	keyFeeCollection = sdk.NewKVStoreKey("fee")
	keyParams        = sdk.NewKVStoreKey("params")
	tkeyParams       = sdk.NewTransientStoreKey("transient_params")
	//keyUpgrade       = sdk.NewKVStoreKey("upgrade")
	keyService       = sdk.NewKVStoreKey("service")
	keyGuardian      = sdk.NewKVStoreKey("guardian")
)
