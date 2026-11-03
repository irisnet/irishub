package v300

import (
	"cosmossdk.io/math"
)

var (
	// BeaconContractAddress is the address of the beacon contract
	BeaconContractAddress = "0xce3d3e91a49ff35b316e7eb84d9fecb067611150"

	// MinDepositRatio is the minimum deposit ratio
	MinDepositRatio = math.LegacyMustNewDecFromStr("0.01")

	// EvmMinGasPrice is the minimum gas price for the EVM
	EvmMinGasPrice = math.LegacyNewDec(50000000000)

	allowMessages = []string{"*"}
)
