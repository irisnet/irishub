package v200

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/types"
)

const (
	MaxBlockGas = 20000000
	genesisTime = "2023-04-15T00:00:00.0Z"
)

var (
	evmParams = etherminttypes.Params{
		EvmDenom:            types.EvmToken.MinUnit,
		EnableCreate:        true,
		EnableCall:          true,
		ChainConfig:         etherminttypes.DefaultChainConfig(),
		ExtraEIPs:           nil,
		AllowUnprotectedTxs: etherminttypes.DefaultAllowUnprotectedTxs,
	}

	feemarketParams = feemarkettypes.Params{
		NoBaseFee:                false,
		BaseFeeChangeDenominator: 8,                        // TODO
		ElasticityMultiplier:     2,                        // TODO
		BaseFee:                  math.NewInt(1000000000),  // TODO
		MinGasPrice:              sdk.ZeroDec(),            // TODO
		MinGasMultiplier:         sdk.NewDecWithPrec(5, 1), // TODO
	}
)

func GenerateEvmParams() etherminttypes.Params {
	return evmParams
}

func GenerateFeemarketParams(enableHeight int64) feemarkettypes.Params {
	feemarketParams.EnableHeight = enableHeight
	return feemarketParams
}

func GenerateGenesisTime() time.Time {
	genTime, err := time.Parse(time.RFC3339Nano, genesisTime)
	if err != nil {
		panic("parse genesis time error: " + err.Error())
	}
	return genTime
}
