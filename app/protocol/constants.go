package protocol

import (
	sdk "github.com/irisnet/irishub/types"
)

const (
	// all store name
	AccountStore         = "acc"
	StakeStore           = "stake"
	StakeTransientStore  = "transient_stake"
	MintStore            = "mint"
	DistrStore           = "distr"
	DistrTransientStore  = "transient_distr"
	SlashingStore        = "slashing"
	GovStore             = "gov"
	FeeStore             = "fee"
	ParamsStore          = "params"
	ParamsTransientStore = "transient_params"
	ServiceStore         = "service"
	GuardianStore        = "guardian"
	UpgradeStore         = "upgrade"
	AssetStore           = "asset"
	RandStore            = "rand"
	SwapStore            = "coinswap"
	HtlcStore            = "htlc"

	// all route for query and handler
	BankRoute     = "bank"
	AccountRoute  = AccountStore
	StakeRoute    = StakeStore
	DistrRoute    = DistrStore
	SlashingRoute = SlashingStore
	GovRoute      = GovStore
	ParamsRoute   = ParamsStore
	ServiceRoute  = ServiceStore
	GuardianRoute = GuardianStore
	UpgradeRoute  = UpgradeStore
	AssetRoute    = AssetStore
	RandRoute     = RandStore
	SwapRoute     = SwapStore
	HtlcRoute     = HtlcStore
)

var (
	KeyMain     = sdk.NewKVStoreKey(sdk.MainStore)
	KeyAccount  = sdk.NewKVStoreKey(AccountStore)
	KeyStake    = sdk.NewKVStoreKey(StakeStore)
	TkeyStake   = sdk.NewTransientStoreKey(StakeTransientStore)
	KeyMint     = sdk.NewKVStoreKey(MintStore)
	KeyDistr    = sdk.NewKVStoreKey(DistrStore)
	TkeyDistr   = sdk.NewTransientStoreKey(DistrTransientStore)
	KeySlashing = sdk.NewKVStoreKey(SlashingStore)
	KeyGov      = sdk.NewKVStoreKey(GovStore)
	KeyFee      = sdk.NewKVStoreKey(FeeStore)
	KeyParams   = sdk.NewKVStoreKey(ParamsStore)
	TkeyParams  = sdk.NewTransientStoreKey(ParamsTransientStore)
	KeyService  = sdk.NewKVStoreKey(ServiceStore)
	KeyGuardian = sdk.NewKVStoreKey(GuardianStore)
	KeyUpgrade  = sdk.NewKVStoreKey(UpgradeStore)
	KeyAsset    = sdk.NewKVStoreKey(AssetStore)
	KeyRand     = sdk.NewKVStoreKey(RandStore)
	KeySwap     = sdk.NewKVStoreKey(SwapStore)
	KeyHtlc     = sdk.NewKVStoreKey(HtlcStore)
)
