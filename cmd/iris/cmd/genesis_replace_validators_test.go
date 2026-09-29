package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/v5/app"
)

const (
	targetChainID = "irishub-1"
	sourceChainID = "user-chain-1"
	targetHeight  = 25_000_000
)

type replaceFixture struct {
	state    map[string]json.RawMessage
	consPub  cryptotypes.PubKey
	valAddr  sdk.ValAddress
	operator sdk.AccAddress
	extra    sdk.AccAddress
	tokens   math.Int
}

func newReplaceValidatorsApp(t *testing.T) *app.IrisApp {
	t.Helper()
	return app.NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
}

// newChainApp returns an app bound to a chain-id, mirroring how the server
// wires baseapp.SetChainID before InitChain.
func newChainApp(t *testing.T, chainID string) *app.IrisApp {
	t.Helper()
	return app.NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{}, baseapp.SetChainID(chainID))
}

func copyGenesisState(state map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(state))
	for name, raw := range state {
		out[name] = raw
	}
	return out
}

func marshalGenesisState(t *testing.T, state map[string]json.RawMessage) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(state)
	require.NoError(t, err)
	return data
}

func bankGenesisState(t *testing.T, balances []banktypes.Balance) *banktypes.GenesisState {
	t.Helper()
	supplyMap := sdk.NewMapCoins(sdk.NewCoins())
	for _, balance := range balances {
		supplyMap.Add(balance.Coins...)
	}
	return banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, supplyMap.ToCoins(), nil, nil)
}

func bondedValidator(t *testing.T, consPub cryptotypes.PubKey, valAddr sdk.ValAddress, tokens math.Int) stakingtypes.Validator {
	t.Helper()
	pkAny, err := codectypes.NewAnyWithValue(consPub)
	require.NoError(t, err)
	return stakingtypes.Validator{
		OperatorAddress:   valAddr.String(),
		ConsensusPubkey:   pkAny,
		Status:            stakingtypes.Bonded,
		Tokens:            tokens,
		DelegatorShares:   math.LegacyNewDecFromInt(tokens),
		Description:       stakingtypes.NewDescription("validator", "", "", "", ""),
		Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
		MinSelfDelegation: math.OneInt(),
	}
}

// buildMainnetLikeTarget returns a mainnet-export-like genesis: an exported
// staking state with one bonded validator whose keys nobody controls, plus
// accounts and balances that must survive the validator replacement.
func buildMainnetLikeTarget(t *testing.T, cdc codec.Codec, defaults map[string]json.RawMessage) replaceFixture {
	t.Helper()
	consPriv := ed25519.GenPrivKey()
	consPub := consPriv.PubKey()
	valAddr := sdk.ValAddress(consPub.Address())
	opPub := secp256k1.GenPrivKey().PubKey()
	operator := sdk.AccAddress(opPub.Address())
	extraPub := secp256k1.GenPrivKey().PubKey()
	extra := sdk.AccAddress(extraPub.Address())
	tokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	state := copyGenesisState(defaults)
	state[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingtypes.GenesisState{
		Params:              stakingtypes.DefaultParams(),
		LastTotalPower:      math.NewInt(10),
		LastValidatorPowers: []stakingtypes.LastValidatorPower{{Address: valAddr.String(), Power: 10}},
		Validators:          []stakingtypes.Validator{bondedValidator(t, consPub, valAddr, tokens)},
		Delegations:         []stakingtypes.Delegation{stakingtypes.NewDelegation(operator.String(), valAddr.String(), math.LegacyNewDecFromInt(tokens))},
		Exported:            true,
	})
	state[authtypes.ModuleName] = cdc.MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(), []authtypes.GenesisAccount{
		authtypes.NewBaseAccount(operator, opPub, 0, 0),
		authtypes.NewBaseAccount(extra, extraPub, 1, 0),
		authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName),
	}))
	state[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesisState(t, []banktypes.Balance{
		{Address: operator.String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(2000)))},
		{Address: extra.String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(999)))},
		{Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))},
		{Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tokens))},
	}))
	return replaceFixture{state: state, consPub: consPub, valAddr: valAddr, operator: operator, extra: extra, tokens: tokens}
}

// buildExportedSource returns the user's own chain export: one bonded
// validator, a self-delegation, distribution reward records and the pool
// balances backing them.
func buildExportedSource(t *testing.T, cdc codec.Codec, defaults map[string]json.RawMessage) replaceFixture {
	t.Helper()
	consPriv := ed25519.GenPrivKey()
	consPub := consPriv.PubKey()
	valAddr := sdk.ValAddress(consPub.Address())
	opPub := secp256k1.GenPrivKey().PubKey()
	operator := sdk.AccAddress(opPub.Address())
	tokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	rewards := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, math.LegacyMustNewDecFromStr("100.5"))}

	state := copyGenesisState(defaults)
	state[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingtypes.GenesisState{
		Params:              stakingtypes.DefaultParams(),
		LastTotalPower:      math.NewInt(10),
		LastValidatorPowers: []stakingtypes.LastValidatorPower{{Address: valAddr.String(), Power: 10}},
		Validators:          []stakingtypes.Validator{bondedValidator(t, consPub, valAddr, tokens)},
		Delegations:         []stakingtypes.Delegation{stakingtypes.NewDelegation(operator.String(), valAddr.String(), math.LegacyNewDecFromInt(tokens))},
		Exported:            true,
	})
	state[authtypes.ModuleName] = cdc.MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(),
		[]authtypes.GenesisAccount{authtypes.NewBaseAccount(operator, opPub, 0, 0)}))
	state[distrtypes.ModuleName] = cdc.MustMarshalJSON(&distrtypes.GenesisState{
		Params:  distrtypes.DefaultParams(),
		FeePool: distrtypes.InitialFeePool(),
		OutstandingRewards: []distrtypes.ValidatorOutstandingRewardsRecord{{
			ValidatorAddress:   valAddr.String(),
			OutstandingRewards: rewards,
		}},
		ValidatorAccumulatedCommissions: []distrtypes.ValidatorAccumulatedCommissionRecord{{
			ValidatorAddress: valAddr.String(),
		}},
		ValidatorHistoricalRewards: []distrtypes.ValidatorHistoricalRewardsRecord{{
			ValidatorAddress: valAddr.String(),
			Period:           0,
			Rewards:          distrtypes.ValidatorHistoricalRewards{ReferenceCount: 1},
		}},
		ValidatorCurrentRewards: []distrtypes.ValidatorCurrentRewardsRecord{{
			ValidatorAddress: valAddr.String(),
			Rewards:          distrtypes.ValidatorCurrentRewards{Rewards: rewards, Period: 1},
		}},
		DelegatorStartingInfos: []distrtypes.DelegatorStartingInfoRecord{{
			ValidatorAddress: valAddr.String(),
			DelegatorAddress: operator.String(),
			StartingInfo:     distrtypes.DelegatorStartingInfo{Height: 1},
		}},
	})
	state[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesisState(t, []banktypes.Balance{
		{Address: operator.String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(500)))},
		{Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tokens))},
		// outstanding rewards of 100.5 truncate to a whole 100 token holding
		{Address: authtypes.NewModuleAddress(distrtypes.ModuleName).String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))},
	}))
	return replaceFixture{state: state, consPub: consPub, valAddr: valAddr, operator: operator, tokens: tokens}
}

// buildGentxSource returns a genesis whose validators arrive through a signed
// create-validator gentx, like the output of "iris testnet".
func buildGentxSource(t *testing.T, cdc codec.Codec, txConfig client.TxConfig, defaults map[string]json.RawMessage) replaceFixture {
	t.Helper()
	consPriv := ed25519.GenPrivKey()
	consPub := consPriv.PubKey()
	opPriv := secp256k1.GenPrivKey()
	opPub := opPriv.PubKey()
	operator := sdk.AccAddress(opPub.Address())
	selfDelegation := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	state := copyGenesisState(defaults)
	state[authtypes.ModuleName] = cdc.MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(),
		[]authtypes.GenesisAccount{authtypes.NewBaseAccount(operator, opPub, 0, 0)}))
	state[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesisState(t, []banktypes.Balance{
		{Address: operator.String(), Coins: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, selfDelegation.MulRaw(10)))},
	}))

	msg, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(operator).String(),
		consPub,
		sdk.NewCoin(sdk.DefaultBondDenom, selfDelegation),
		stakingtypes.NewDescription("my-validator", "", "", "", ""),
		stakingtypes.NewCommissionRates(math.LegacyOneDec(), math.LegacyOneDec(), math.LegacyOneDec()),
		math.OneInt(),
	)
	require.NoError(t, err)
	builder := txConfig.NewTxBuilder()
	require.NoError(t, builder.SetMsgs(msg))

	// Sign the gentx for the source chain, account 0 sequence 0, mirroring
	// what "iris testnet" produces for its operators.
	signerData := authsigning.SignerData{
		ChainID:       sourceChainID,
		AccountNumber: 0,
		Sequence:      0,
		PubKey:        opPub,
		Address:       operator.String(),
	}
	emptySig := signing.SignatureV2{
		PubKey:   opPub,
		Data:     &signing.SingleSignatureData{SignMode: signing.SignMode_SIGN_MODE_DIRECT},
		Sequence: 0,
	}
	require.NoError(t, builder.SetSignatures(emptySig))
	signBytes, err := authsigning.GetSignBytesAdapter(context.Background(), txConfig.SignModeHandler(), signing.SignMode_SIGN_MODE_DIRECT, signerData, builder.GetTx())
	require.NoError(t, err)
	sigBytes, err := opPriv.Sign(signBytes)
	require.NoError(t, err)
	require.NoError(t, builder.SetSignatures(signing.SignatureV2{
		PubKey:   opPub,
		Data:     &signing.SingleSignatureData{SignMode: signing.SignMode_SIGN_MODE_DIRECT, Signature: sigBytes},
		Sequence: 0,
	}))
	gentxJSON, err := txConfig.TxJSONEncoder()(builder.GetTx())
	require.NoError(t, err)
	state[genutiltypes.ModuleName] = cdc.MustMarshalJSON(&genutiltypes.GenesisState{
		GenTxs: []json.RawMessage{gentxJSON},
	})
	return replaceFixture{state: state, consPub: consPub, valAddr: sdk.ValAddress(consPub.Address()), operator: operator, tokens: selfDelegation}
}

func saveGenesisDoc(t *testing.T, dir, name string, doc *tmtypes.GenesisDoc) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, doc.SaveAs(path))
	return path
}

func targetGenesisDoc(t *testing.T, fixture replaceFixture) *tmtypes.GenesisDoc {
	t.Helper()
	cmtPub, err := cryptocodec.ToCmtPubKeyInterface(fixture.consPub)
	require.NoError(t, err)
	return &tmtypes.GenesisDoc{
		ChainID:       targetChainID,
		InitialHeight: targetHeight,
		GenesisTime:   time.Now().UTC(),
		Validators:    []tmtypes.GenesisValidator{{Address: cmtPub.Address(), PubKey: cmtPub, Power: 10}},
		AppState:      marshalGenesisState(t, fixture.state),
	}
}

func sourceGenesisDoc(t *testing.T, fixture replaceFixture) *tmtypes.GenesisDoc {
	t.Helper()
	return &tmtypes.GenesisDoc{
		ChainID:     sourceChainID,
		GenesisTime: time.Now().UTC(),
		AppState:    marshalGenesisState(t, fixture.state),
	}
}

// replaceStates runs the replacement on the given states and returns the
// merged genesis document.
func replaceStates(
	t *testing.T,
	cdc codec.Codec,
	txConfig client.TxConfig,
	basics module.BasicManager,
	sourceState, targetState map[string]json.RawMessage,
) (*tmtypes.GenesisDoc, error) {
	t.Helper()
	dir := t.TempDir()
	sourceDoc := sourceGenesisDoc(t, replaceFixture{state: sourceState})
	targetDoc := &tmtypes.GenesisDoc{
		ChainID:       targetChainID,
		InitialHeight: targetHeight,
		GenesisTime:   time.Now().UTC(),
		AppState:      marshalGenesisState(t, targetState),
	}
	outputPath := filepath.Join(dir, "merged.json")
	if _, err := ReplaceValidators(cdc, txConfig, basics, sourceDoc, targetDoc, outputPath); err != nil {
		return nil, err
	}
	merged, err := tmtypes.GenesisDocFromFile(outputPath)
	require.NoError(t, err)
	return merged, nil
}

func genesisBalances(t *testing.T, cdc codec.Codec, state map[string]json.RawMessage) map[string]sdk.Coins {
	t.Helper()
	var bankState banktypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[banktypes.ModuleName], &bankState))
	balances := make(map[string]sdk.Coins, len(bankState.Balances))
	for _, balance := range bankState.Balances {
		balances[balance.Address] = balance.Coins
	}
	return balances
}

func genesisAccountAddrs(t *testing.T, cdc codec.Codec, state map[string]json.RawMessage) map[string]bool {
	t.Helper()
	var authState authtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[authtypes.ModuleName], &authState))
	accounts, err := authtypes.UnpackAccounts(authState.Accounts)
	require.NoError(t, err)
	addrs := make(map[string]bool, len(accounts))
	for _, account := range accounts {
		addrs[account.GetAddress().String()] = true
	}
	return addrs
}

func TestGenesisReplaceValidatorsExportedSource(t *testing.T) {
	irisApp := newReplaceValidatorsApp(t)
	cdc := irisApp.AppCodec()
	target := buildMainnetLikeTarget(t, cdc, irisApp.DefaultGenesis())
	source := buildExportedSource(t, cdc, irisApp.DefaultGenesis())

	dir := t.TempDir()
	targetPath := saveGenesisDoc(t, dir, "target.json", targetGenesisDoc(t, target))
	sourcePath := saveGenesisDoc(t, dir, "source.json", sourceGenesisDoc(t, source))
	outputPath := filepath.Join(dir, "merged.json")

	// run through the CLI to cover the command registration and flags
	command := NewRootCmd()
	// server/cmd.Execute installs this flag around NewRootCmd in the binary.
	command.PersistentFlags().String(flags.FlagHome, t.TempDir(), "node home")
	var output bytes.Buffer
	command.SetOut(&output)
	command.SetErr(&output)
	command.SetArgs([]string{
		"genesis", "replace-validators",
		"--source-genesis-file", sourcePath,
		"--target-genesis-file", targetPath,
		"--output-genesis-file", outputPath,
		"--home", t.TempDir(),
	})
	require.NoError(t, command.Execute())
	require.Contains(t, output.String(), "exported validator")

	merged, err := tmtypes.GenesisDocFromFile(outputPath)
	require.NoError(t, err)
	require.Equal(t, targetChainID, merged.ChainID)
	require.Equal(t, int64(targetHeight), merged.InitialHeight)
	require.Empty(t, merged.Validators)

	var mergedState map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(merged.AppState, &mergedState))

	// staking comes from the source and keeps export semantics
	var stakingState stakingtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(mergedState[stakingtypes.ModuleName], &stakingState))
	require.Len(t, stakingState.Validators, 1)
	require.Equal(t, source.valAddr.String(), stakingState.Validators[0].OperatorAddress)
	require.True(t, stakingState.Exported)
	require.Equal(t, source.valAddr.String(), stakingState.LastValidatorPowers[0].Address)

	// distribution rewards come from the source verbatim (the merged file
	// reindents raw JSON, so compare the parsed state)
	var mergedDistr, sourceDistr distrtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(mergedState[distrtypes.ModuleName], &mergedDistr))
	require.NoError(t, cdc.UnmarshalJSON(source.state[distrtypes.ModuleName], &sourceDistr))
	require.Equal(t, sourceDistr, mergedDistr)

	// pool balances follow the imported modules, target data survives, the
	// source operator is funded, and the supply matches the balances
	balances := genesisBalances(t, cdc, mergedState)
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, source.tokens)), balances[authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))), balances[authtypes.NewModuleAddress(distrtypes.ModuleName).String()])
	require.Empty(t, balances[authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(500))), balances[source.operator.String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(2000))), balances[target.operator.String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(999))), balances[target.extra.String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100))), balances[authtypes.NewModuleAddress(authtypes.FeeCollectorName).String()])
	var bankState banktypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(mergedState[banktypes.ModuleName], &bankState))
	supplyMap := sdk.NewMapCoins(sdk.NewCoins())
	for _, balance := range bankState.Balances {
		supplyMap.Add(balance.Coins...)
	}
	require.Equal(t, supplyMap.ToCoins(), bankState.Supply)

	// both the source operator and the target accounts are importable
	addrs := genesisAccountAddrs(t, cdc, mergedState)
	require.True(t, addrs[source.operator.String()])
	require.True(t, addrs[target.operator.String()])
	require.True(t, addrs[target.extra.String()])

	// the merged genesis starts a chain driven by the source validator
	simApp := newChainApp(t, merged.ChainID)
	res, err := simApp.InitChain(&abci.RequestInitChain{
		ChainId:         merged.ChainID,
		AppStateBytes:   merged.AppState,
		ConsensusParams: simtestutil.DefaultConsensusParams,
		InitialHeight:   merged.InitialHeight,
	})
	require.NoError(t, err)
	require.Len(t, res.Validators, 1)
	require.Equal(t, source.consPub.Bytes(), res.Validators[0].PubKey.GetEd25519())
	require.Equal(t, int64(10), res.Validators[0].Power)
}

func TestGenesisReplaceValidatorsGentxSource(t *testing.T) {
	irisApp := newReplaceValidatorsApp(t)
	cdc := irisApp.AppCodec()
	txConfig := irisApp.EncodingConfig().TxConfig
	defaults := irisApp.DefaultGenesis()
	target := buildMainnetLikeTarget(t, cdc, defaults)
	source := buildGentxSource(t, cdc, txConfig, defaults)

	merged, err := replaceStates(t, cdc, txConfig, irisApp.BasicManager(), source.state, target.state)
	require.NoError(t, err)

	// gentx signatures verify against the source chain-id, and a genesis-time
	// initial height (0 written; cometbft defaults it to 1 on load) makes
	// signature verification skip renumbered account numbers
	require.Equal(t, sourceChainID, merged.ChainID)
	require.LessOrEqual(t, merged.InitialHeight, int64(1))
	require.Empty(t, merged.Validators)

	var mergedState map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(merged.AppState, &mergedState))

	// staking stays default: the gentx creates the validator at chain start
	var stakingState stakingtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(mergedState[stakingtypes.ModuleName], &stakingState))
	require.Empty(t, stakingState.Validators)
	require.False(t, stakingState.Exported)

	var genutilState genutiltypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(mergedState[genutiltypes.ModuleName], &genutilState))
	require.Len(t, genutilState.GenTxs, 1)

	// the target pool balances are dropped: the gentx funds the bonded pool
	balances := genesisBalances(t, cdc, mergedState)
	require.Empty(t, balances[authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String()])
	require.Empty(t, balances[authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String()])
	require.Empty(t, balances[authtypes.NewModuleAddress(distrtypes.ModuleName).String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, source.tokens.MulRaw(10))), balances[source.operator.String()])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(999))), balances[target.extra.String()])

	addrs := genesisAccountAddrs(t, cdc, mergedState)
	require.True(t, addrs[source.operator.String()])
	require.True(t, addrs[target.operator.String()])

	// the signed gentx is delivered at chain start and produces the user's
	// validator, even though its account number was renumbered behind the
	// target accounts
	simApp := newChainApp(t, merged.ChainID)
	res, err := simApp.InitChain(&abci.RequestInitChain{
		ChainId:         merged.ChainID,
		AppStateBytes:   merged.AppState,
		ConsensusParams: simtestutil.DefaultConsensusParams,
	})
	require.NoError(t, err)
	require.Len(t, res.Validators, 1)
	require.Equal(t, source.consPub.Bytes(), res.Validators[0].PubKey.GetEd25519())
	require.Positive(t, res.Validators[0].Power)
}

// withLegacyLSMFields adds the disabled legacy LSM fields that the
// app/genesis/lsm preflight strips or converts.
func withLegacyLSMFields(t *testing.T, current json.RawMessage, records string) json.RawMessage {
	t.Helper()
	var state map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(current, &state))
	var params map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(state["params"], &params))
	params["validator_bond_factor"] = json.RawMessage(`"250.000000000000000000"`)
	params["global_liquid_staking_cap"] = json.RawMessage(`"0.250000000000000000"`)
	params["validator_liquid_staking_cap"] = json.RawMessage(`"0.500000000000000000"`)
	state["params"], _ = json.Marshal(params)
	state["tokenize_share_records"] = json.RawMessage(records)
	state["last_tokenize_share_record_id"] = json.RawMessage(`"1"`)
	state["total_liquid_staked_tokens"] = json.RawMessage(`"0"`)
	state["tokenize_share_locks"] = json.RawMessage(`[]`)
	for _, field := range []string{"validators", "delegations"} {
		var entries []map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(state[field], &entries))
		for _, entry := range entries {
			if field == "validators" {
				entry["liquid_shares"] = json.RawMessage(`"0.000000000000000000"`)
				entry["validator_bond_shares"] = json.RawMessage(`"0.000000000000000000"`)
			} else {
				entry["validator_bond"] = json.RawMessage(`false`)
			}
		}
		state[field], _ = json.Marshal(entries)
	}
	result, err := json.Marshal(state)
	require.NoError(t, err)
	return result
}

// buildLiveLSMTarget extends the mainnet-like target with one valid live
// tokenize-share record whose claims the lsm genesis migration would convert.
func buildLiveLSMTarget(t *testing.T, cdc codec.Codec, defaults map[string]json.RawMessage) map[string]json.RawMessage {
	t.Helper()
	base := buildMainnetLikeTarget(t, cdc, defaults)
	state := copyGenesisState(base.state)

	ownerPub := secp256k1.GenPrivKey().PubKey()
	owner := sdk.AccAddress(ownerPub.Address())
	shareSource := sdk.AccAddress(address.Module("lsm", []byte("tokenizeshare_1")))
	denom := strings.ToLower(base.valAddr.String()) + "/1"
	tokenized := math.NewInt(50)

	var stakingState stakingtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[stakingtypes.ModuleName], &stakingState))
	// move 50 shares of the validator to the tokenize-share module account
	shares := stakingState.Validators[0].DelegatorShares.Sub(math.LegacyNewDecFromInt(tokenized))
	stakingState.Validators[0].DelegatorShares = stakingState.Validators[0].DelegatorShares.Add(math.LegacyNewDecFromInt(tokenized))
	stakingState.Delegations[0].Shares = shares
	stakingState.Delegations = append(stakingState.Delegations, stakingtypes.NewDelegation(shareSource.String(), base.valAddr.String(), math.LegacyNewDecFromInt(tokenized)))
	state[stakingtypes.ModuleName] = withLegacyLSMFields(t, cdc.MustMarshalJSON(&stakingState),
		fmt.Sprintf(`[{"id":"1","owner":%q,"module_account":"tokenizeshare_1","validator":%q}]`, owner.String(), base.valAddr.String()))

	var distrState distrtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[distrtypes.ModuleName], &distrState))
	distrState.DelegatorStartingInfos = append(distrState.DelegatorStartingInfos, distrtypes.DelegatorStartingInfoRecord{
		ValidatorAddress: base.valAddr.String(),
		DelegatorAddress: shareSource.String(),
		StartingInfo:     distrtypes.DelegatorStartingInfo{Height: 1},
	})
	state[distrtypes.ModuleName] = cdc.MustMarshalJSON(&distrState)

	var authState authtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[authtypes.ModuleName], &authState))
	accounts, err := authtypes.UnpackAccounts(authState.Accounts)
	require.NoError(t, err)
	accounts = append(accounts,
		authtypes.NewBaseAccount(owner, ownerPub, 5, 0),
		authtypes.NewBaseAccount(shareSource, nil, 6, 0),
	)
	state[authtypes.ModuleName] = cdc.MustMarshalJSON(authtypes.NewGenesisState(authtypes.DefaultParams(), accounts))

	var bankState banktypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(state[banktypes.ModuleName], &bankState))
	bankState.Balances = append(bankState.Balances, banktypes.Balance{
		Address: owner.String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(denom, tokenized)),
	})
	state[banktypes.ModuleName] = cdc.MustMarshalJSON(bankGenesisState(t, bankState.Balances))
	return state
}

func TestGenesisReplaceValidatorsRejectsLiveLSMTarget(t *testing.T) {
	irisApp := newReplaceValidatorsApp(t)
	cdc := irisApp.AppCodec()
	txConfig := irisApp.EncodingConfig().TxConfig
	defaults := irisApp.DefaultGenesis()
	target := buildMainnetLikeTarget(t, cdc, defaults)
	source := buildGentxSource(t, cdc, txConfig, defaults)

	// any tokenize-share record in the target is rejected: its validators are
	// removed, and the records cannot be converted offline
	bad := copyGenesisState(target.state)
	bad[stakingtypes.ModuleName] = withLegacyLSMFields(t, bad[stakingtypes.ModuleName],
		`[{"id":"1","owner":"invalid","module_account":"tokenizeshare_1","validator":"invalid"}]`)
	_, err := replaceStates(t, cdc, txConfig, irisApp.BasicManager(), source.state, bad)
	require.Error(t, err)
	require.Contains(t, err.Error(), "live tokenize-share records")

	// a valid live record would dangle once its validators are removed
	_, err = replaceStates(t, cdc, txConfig, irisApp.BasicManager(), source.state, buildLiveLSMTarget(t, cdc, defaults))
	require.Error(t, err)
	require.Contains(t, err.Error(), "live tokenize-share records")

	// empty legacy LSM state carries no claims and is accepted
	empty := copyGenesisState(target.state)
	empty[stakingtypes.ModuleName] = withLegacyLSMFields(t, empty[stakingtypes.ModuleName], `[]`)
	merged, err := replaceStates(t, cdc, txConfig, irisApp.BasicManager(), source.state, empty)
	require.NoError(t, err)
	require.Equal(t, sourceChainID, merged.ChainID)
}

func TestGenesisReplaceValidatorsSourceErrors(t *testing.T) {
	irisApp := newReplaceValidatorsApp(t)
	cdc := irisApp.AppCodec()
	txConfig := irisApp.EncodingConfig().TxConfig
	defaults := irisApp.DefaultGenesis()
	target := buildMainnetLikeTarget(t, cdc, defaults)
	basics := irisApp.BasicManager()

	// validators both in staking state and in gentxs
	both := copyGenesisState(buildExportedSource(t, cdc, defaults).state)
	both[genutiltypes.ModuleName] = cdc.MustMarshalJSON(&genutiltypes.GenesisState{
		GenTxs: []json.RawMessage{json.RawMessage(`{}`)},
	})
	_, err := replaceStates(t, cdc, txConfig, basics, both, target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "only one mechanism")

	// no validators at all
	_, err = replaceStates(t, cdc, txConfig, basics, copyGenesisState(defaults), target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "defines no validators")

	// validators without last_validator_powers cannot restore a validator set
	noPowers := copyGenesisState(buildExportedSource(t, cdc, defaults).state)
	var stakingState stakingtypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(noPowers[stakingtypes.ModuleName], &stakingState))
	stakingState.LastValidatorPowers = nil
	noPowers[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingState)
	_, err = replaceStates(t, cdc, txConfig, basics, noPowers, target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no bonded validator power")

	// powers referencing a missing validator would panic at chain start
	danglingPower := copyGenesisState(noPowers)
	missingValAddr := sdk.ValAddress(ed25519.GenPrivKey().PubKey().Address())
	require.NoError(t, cdc.UnmarshalJSON(danglingPower[stakingtypes.ModuleName], &stakingState))
	stakingState.LastValidatorPowers = []stakingtypes.LastValidatorPower{{Address: missingValAddr.String(), Power: 10}}
	danglingPower[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingState)
	_, err = replaceStates(t, cdc, txConfig, basics, danglingPower, target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "references missing validator")

	// duplicate power entries are rejected as well
	duplicatePower := copyGenesisState(buildExportedSource(t, cdc, defaults).state)
	require.NoError(t, cdc.UnmarshalJSON(duplicatePower[stakingtypes.ModuleName], &stakingState))
	duplicate := stakingState.LastValidatorPowers[0]
	stakingState.LastValidatorPowers = append(stakingState.LastValidatorPowers, duplicate)
	duplicatePower[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingState)
	_, err = replaceStates(t, cdc, txConfig, basics, duplicatePower, target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate entry")

	// pool balances not backing the staking state
	wrongPool := copyGenesisState(buildExportedSource(t, cdc, defaults).state)
	var bankState banktypes.GenesisState
	require.NoError(t, cdc.UnmarshalJSON(wrongPool[banktypes.ModuleName], &bankState))
	for i := range bankState.Balances {
		if bankState.Balances[i].Address == authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String() {
			bankState.Balances[i].Coins = bankState.Balances[i].Coins.Add(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt()))
		}
	}
	wrongPool[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankState)
	_, err = replaceStates(t, cdc, txConfig, basics, wrongPool, target.state)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bonded pool")
}
