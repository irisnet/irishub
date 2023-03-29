package v200

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	tokenv1 "github.com/irisnet/irismod/modules/token/types/v1"

	"github.com/irisnet/irishub/types"
)

// NOTE: Before the release of irishub 2.0.0, the configuration in this file must be modified

const (
	MaxBlockGas = 20000000 // TODO
)

var (
	evmToken  = types.EvmToken // TODO
	evmParams = etherminttypes.Params{
		EvmDenom:            evmToken.MinUnit,
		EnableCreate:        true,
		EnableCall:          true,
		ChainConfig:         etherminttypes.DefaultChainConfig(),
		ExtraEIPs:           nil,
		AllowUnprotectedTxs: etherminttypes.DefaultAllowUnprotectedTxs,
	}

	feemarketParams = feemarkettypes.Params{
		NoBaseFee:                false,                    // TODO
		BaseFeeChangeDenominator: 8,                        // TODO
		ElasticityMultiplier:     2,                        // TODO
		BaseFee:                  math.NewInt(1000000000),  // TODO
		MinGasPrice:              sdk.ZeroDec(),            // TODO
		MinGasMultiplier:         sdk.NewDecWithPrec(5, 1), // TODO
	}
)

func GetEvmToken() tokenv1.Token {
	return evmToken
}

func GenerateEvmParams() etherminttypes.Params {
	return evmParams
}

func GenerateFeemarketParams(enableHeight int64) feemarkettypes.Params {
	feemarketParams.EnableHeight = enableHeight
	return feemarketParams
}
