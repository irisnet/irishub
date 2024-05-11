package v300

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// ValidatorBondFactor dictates the cap on the liquid shares
	// for a validator - determined as a multiple to their validator bond
	// (e.g. ValidatorBondShares = 1000, BondFactor = 250 -> LiquidSharesCap: 250,000)
	ValidatorBondFactor = sdk.NewDec(250)
	// ValidatorLiquidStakingCap represents a cap on the portion of stake that
	// comes from liquid staking providers for a specific validator
	ValidatorLiquidStakingCap = sdk.MustNewDecFromStr("1") // 100%
	// GlobalLiquidStakingCap represents the percentage cap on
	// the portion of a chain's total stake can be liquid
	GlobalLiquidStakingCap = sdk.MustNewDecFromStr("0.25") // 25%

	// BeaconContractAddress is the address of the beacon contract
	BeaconContractAddress = "0xce3d3e91a49ff35b316e7eb84d9fecb067611150"

	// MinDepositRatio is the minimum deposit ratio
	MinDepositRatio = sdk.MustNewDecFromStr("0.01")

	allowMessages = []string{"*"}
)
