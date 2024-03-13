package testutil

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/irisnet/irishub/v3/app"
)

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(t *testing.T, _ bool) *AppBuilder {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(
		senderPrivKey.PubKey().Address().Bytes(),
		senderPrivKey.PubKey(),
		0,
		0,
	)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000000000))),
	}
	return SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)
}

// NewConfig returns a new app config
func NewConfig() network.Config {
	cfg := network.DefaultConfig(NewTestNetworkFixture)
	encCfg := app.RegisterEncodingConfig()
	cfg.Codec = encCfg.Marshaler
	cfg.TxConfig = encCfg.TxConfig
	cfg.LegacyAmino = encCfg.Amino
	cfg.InterfaceRegistry = encCfg.InterfaceRegistry
	cfg.AppConstructor = func(val network.ValidatorI) servertypes.Application {
		builder := &AppBuilder{}
		return builder.build(
			nil,
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(cfg.ChainID),
		)
	}
	cfg.GenesisState = app.ModuleBasics.DefaultGenesis(cfg.Codec)
	return cfg
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(
	t *testing.T,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) *AppBuilder {
	t.Helper()

	builder := &AppBuilder{}
	app := builder.build(nil)

	genesisState, err := GenesisStateWithValSet(
		app.AppCodec(),
		app.DefaultGenesis(),
		valSet,
		genAccs,
		balances...,
	)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simtestutil.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return &AppBuilder{app}
}

// GenesisStateWithValSet returns a new genesis state with the validator set
func GenesisStateWithValSet(codec codec.Codec, genesisState map[string]json.RawMessage,
	valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress: sdk.ValAddress(val.Address).String(),
			ConsensusPubkey: pkAny,
			Jailed:          false,
			Status:          stakingtypes.Bonded,
			Tokens:          bondAmt,
			DelegatorShares: sdk.OneDec(),
			Description:     stakingtypes.Description{},
			UnbondingHeight: int64(0),
			UnbondingTime:   time.Unix(0, 0).UTC(),
			Commission: stakingtypes.NewCommission(
				sdk.ZeroDec(),
				sdk.ZeroDec(),
				sdk.ZeroDec(),
			),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(
			delegations,
			stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()),
		)

	}
	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(
		stakingtypes.DefaultParams(),
		validators,
		delegations,
	)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
		nil,
	)
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}

// NewTestNetworkFixture returns a new simapp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	appCtr := func(val network.ValidatorI) servertypes.Application {
		builder := &AppBuilder{}
		return builder.build(
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir),
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}
	ec := MakeCodecs()

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   DefaultGenesis(ec.Marshaler),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: ec.InterfaceRegistry,
			Codec:             ec.Marshaler,
			TxConfig:          ec.TxConfig,
			Amino:             ec.Amino,
		},
	}
}

// ExecTxCmdWithResult executes a tx command and returns the result
func ExecTxCmdWithResult(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	cmd *cobra.Command,
	extraArgs []string,
) *coretypes.ResultTx {
	buf, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, extraArgs)
	require.NoError(t, err, "ExecTestCLICmd failed")
	require.NoError(t, network.WaitForNextBlock(), "network.WaitForNextBlock failed")

	respType := proto.Message(&sdk.TxResponse{})
	require.NoError(t, clientCtx.Codec.UnmarshalJSON(buf.Bytes(), respType), buf.String())

	txResp := respType.(*sdk.TxResponse)
	require.Equal(t, uint32(0), txResp.Code)
	return WaitForTx(t, network, clientCtx, txResp.TxHash)
}

// WaitForTx waits for a tx to be included in a block
func WaitForTx(
	t *testing.T,
	network *network.Network,
	clientCtx client.Context,
	txHash string,
) *coretypes.ResultTx {
	var (
		result *coretypes.ResultTx
		err    error
		tryCnt = 3
	)

	txHashBz, err := hex.DecodeString(txHash)
	require.NoError(t, err, "hex.DecodeString failed")

reTry:
	result, err = clientCtx.Client.Tx(context.Background(), txHashBz, false)
	if err != nil && strings.Contains(err.Error(), "not found") && tryCnt > 0 {
		require.NoError(t, network.WaitForNextBlock(), "network.WaitForNextBlock failed")
		tryCnt--
		goto reTry
	}

	require.NoError(t, err, "query tx failed")
	return result
}
