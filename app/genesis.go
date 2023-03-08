package app

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	etherminttypes "github.com/evmos/ethermint/x/evm/types"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/types"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() types.GenesisState {
	encCfg := MakeEncodingConfig()
	return ModuleBasics.DefaultGenesis(encCfg.Marshaler)
}

func InitGenesis(appCodec codec.Codec, genesisState types.GenesisState) {
	// add evm genesis
	if _, ok := genesisState[etherminttypes.ModuleName]; !ok {
		evmGenState := etherminttypes.GenesisState{
			Accounts: []etherminttypes.GenesisAccount{},
			Params: etherminttypes.Params{
				EvmDenom:            types.EvmToken.MinUnit,
				EnableCreate:        true,
				EnableCall:          true,
				ChainConfig:         etherminttypes.DefaultChainConfig(),
				ExtraEIPs:           nil,
				AllowUnprotectedTxs: etherminttypes.DefaultAllowUnprotectedTxs,
			},
		}
		genesisState[etherminttypes.ModuleName] = appCodec.MustMarshalJSON(&evmGenState)
	}

	// add feemarket genesis
	if _, ok := genesisState[feemarkettypes.ModuleName]; !ok {
		evmGenState := feemarkettypes.GenesisState{
			Params: feemarkettypes.Params{
				NoBaseFee:                false,
				BaseFeeChangeDenominator: 8,                        // TODO
				ElasticityMultiplier:     2,                        // TODO
				EnableHeight:             0,                        // TODO
				BaseFee:                  math.NewInt(1000000000),  // TODO
				MinGasPrice:              sdk.ZeroDec(),            // TODO
				MinGasMultiplier:         sdk.NewDecWithPrec(5, 1), // TODO
			},
			BlockGas: 0,
		}
		genesisState[feemarkettypes.ModuleName] = appCodec.MustMarshalJSON(&evmGenState)
	}

	// add token genesis
	{
		var tokenGenState tokentypes.GenesisState
		appCodec.MustUnmarshalJSON(genesisState[tokentypes.ModuleName], &tokenGenState)

		evmTokenExist := false
		for _, token := range tokenGenState.Tokens {
			if token.MinUnit == types.EvmToken.MinUnit {
				break
			}
		}
		if !evmTokenExist {
			tokenGenState.Tokens = append(tokenGenState.Tokens, types.EvmToken)
		}
		genesisState[tokentypes.ModuleName] = appCodec.MustMarshalJSON(&tokenGenState)
	}

	// add service genesis
	{
		var serviceGenState servicetypes.GenesisState
		appCodec.MustUnmarshalJSON(genesisState[servicetypes.ModuleName], &serviceGenState)

		var (
			oracleServiceExist = false
			oracleBindingExist = false
			randomServiceExist = false

			oracleDefinition = servicetypes.GenOraclePriceSvcDefinition()
			randomDefinition = randomtypes.GetSvcDefinition()
		)

		for _, definition := range serviceGenState.Definitions {
			switch definition.Name {
			case oracleDefinition.Name:
				oracleServiceExist = true
			case randomDefinition.Name:
				randomServiceExist = true
			}
		}
		if !oracleServiceExist {
			serviceGenState.Definitions = append(serviceGenState.Definitions, oracleDefinition)
		}
		if !randomServiceExist {
			serviceGenState.Definitions = append(serviceGenState.Definitions, randomDefinition)
		}

		for _, binding := range serviceGenState.Bindings {
			if binding.ServiceName == oracleDefinition.Name {
				oracleBindingExist = true
				break
			}
		}
		if !oracleBindingExist {
			serviceGenState.Bindings = append(serviceGenState.Bindings, servicetypes.GenOraclePriceSvcBinding(types.NativeToken.MinUnit))
		}
		genesisState[servicetypes.ModuleName] = appCodec.MustMarshalJSON(&serviceGenState)
	}
}
