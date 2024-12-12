package v300

import (
	"cosmossdk.io/math"
)

var (
	// ValidatorBondFactor dictates the cap on the liquid shares
	// for a validator - determined as a multiple to their validator bond
	// (e.g. ValidatorBondShares = 1000, BondFactor = 250 -> LiquidSharesCap: 250,000)
	ValidatorBondFactor = math.LegacyNewDec(250)
	// ValidatorLiquidStakingCap represents a cap on the portion of stake that
	// comes from liquid staking providers for a specific validator
	ValidatorLiquidStakingCap = math.LegacyMustNewDecFromStr("1") // 100%
	// GlobalLiquidStakingCap represents the percentage cap on
	// the portion of a chain's total stake can be liquid
	GlobalLiquidStakingCap = math.LegacyMustNewDecFromStr("0.25") // 25%

	// BeaconContractAddress is the address of the beacon contract
	BeaconContractAddress = "0xce3d3e91a49ff35b316e7eb84d9fecb067611150"

	// MinDepositRatio is the minimum deposit ratio
	MinDepositRatio = math.LegacyMustNewDecFromStr("0.01")

	// EvmMinGasPrice is the minimum gas price for the EVM
	EvmMinGasPrice = math.LegacyNewDec(50000000000)

	allowMessages = []string{"*"}
)
