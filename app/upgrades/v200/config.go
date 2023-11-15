package v200

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/v2/types"
)

// NOTE: Before the release of irishub 2.0.0, the configuration in this file must be modified

const (
	maxBlockGas = 40000000
)

var (
	evmToken  = types.EvmToken
	evmParams = etherminttypes.Params{
		EvmDenom:            evmToken.MinUnit,
		EnableCreate:        true,
		EnableCall:          true,
		ChainConfig:         etherminttypes.DefaultChainConfig(),
		ExtraEIPs:           nil,
		AllowUnprotectedTxs: etherminttypes.DefaultAllowUnprotectedTxs,
	}

	feemarketParams = feemarkettypes.Params{
		NoBaseFee:                false,
		BaseFeeChangeDenominator: 8,
		ElasticityMultiplier:     4,
		BaseFee:                  math.NewInt(500000000000),
		MinGasPrice:              sdk.NewDecFromInt(math.NewInt(500000000000)),
		MinGasMultiplier:         sdk.NewDecWithPrec(5, 1),
	}
)

func generateFeemarketParams(enableHeight int64) feemarkettypes.Params {
	feemarketParams.EnableHeight = enableHeight
	return feemarketParams
}
