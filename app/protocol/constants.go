package protocol

import sdk "github.com/irisnet/irishub/types"

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

	// all query route
	AccountRouter  = AccountStore
	StakeRouter    = StakeStore
	DistrRouter    = DistrStore
	SlashingRouter = SlashingStore
	GovRouter      = GovStore
	ParamsRouter   = ParamsStore
	ServiceRouter  = ServiceStore
	GuardianRouter = GuardianStore
	UpgradeRouter  = UpgradeStore
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
)
