package lsm_test

import (
	"encoding/json"
	"strings"
	"testing"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	irisapp "github.com/irisnet/irishub/v5/app"
	lsmgenesis "github.com/irisnet/irishub/v5/app/genesis/lsm"
	"github.com/irisnet/irishub/v5/app/params"
)

func legacyStakingGenesis(t *testing.T, current json.RawMessage) json.RawMessage {
	t.Helper()
	var state map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(current, &state))
	var stakingParams map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(state["params"], &stakingParams))
	stakingParams["validator_bond_factor"] = json.RawMessage(`"250.000000000000000000"`)
	stakingParams["global_liquid_staking_cap"] = json.RawMessage(`"0.250000000000000000"`)
	stakingParams["validator_liquid_staking_cap"] = json.RawMessage(`"0.500000000000000000"`)
	state["params"], _ = json.Marshal(stakingParams)
	state["tokenize_share_records"] = json.RawMessage(`[]`)
	state["last_tokenize_share_record_id"] = json.RawMessage(`"3"`)
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

func TestLegacyStakingGenesisCompatibility(t *testing.T) {
	encoding := params.MakeEncodingConfig()
	base := staking.AppModuleBasic{}
	current := base.DefaultGenesis(encoding.Codec)
	legacy := legacyStakingGenesis(t, current)
	// Reproduce the failure with the upstream decoder, before applying compatibility.
	require.ErrorContains(t, base.ValidateGenesis(encoding.Codec, encoding.TxConfig, legacy), "unknown field")
	compat := lsmgenesis.StakingModule{}
	require.NoError(t, compat.ValidateGenesis(encoding.Codec, encoding.TxConfig, legacy))
	converted, err := lsmgenesis.NormalizeStakingGenesis(legacy)
	require.NoError(t, err)
	require.JSONEq(t, string(current), string(converted))
	again, err := lsmgenesis.NormalizeStakingGenesis(converted)
	require.NoError(t, err)
	require.Equal(t, converted, again)
	unchanged, err := lsmgenesis.NormalizeStakingGenesis(current)
	require.NoError(t, err)
	require.Equal(t, current, unchanged)
	camel := strings.NewReplacer(
		"tokenize_share_records", "tokenizeShareRecords", "last_tokenize_share_record_id", "lastTokenizeShareRecordId",
		"total_liquid_staked_tokens", "totalLiquidStakedTokens", "tokenize_share_locks", "tokenizeShareLocks",
		"validator_bond_factor", "validatorBondFactor", "global_liquid_staking_cap", "globalLiquidStakingCap",
		"validator_liquid_staking_cap", "validatorLiquidStakingCap",
	).Replace(string(legacy))
	require.NoError(t, compat.ValidateGenesis(encoding.Codec, encoding.TxConfig, json.RawMessage(camel)))

	for _, test := range []struct {
		name, field, value, want string
	}{
		{"records", "tokenize_share_records", `[{"id":"1"},{"id":"2"},{"id":"3"}]`, "3 live LSM entries"},
		{"locks", "tokenize_share_locks", `[{"address":"delegator","status":"LOCKED"}]`, "live LSM entries"},
		{"camel records", "tokenizeShareRecords", `[{"id":"1"}]`, "live LSM entries"},
		{"camel locks", "tokenizeShareLocks", `[{}]`, "live LSM entries"},
		{"camel total", "totalLiquidStakedTokens", `"1"`, "live LSM state"},
		{"camel liquid", "validators", `[{"liquidShares":"1"}]`, "live LSM state"},
		{"camel bond shares", "validators", `[{"validatorBondShares":"1"}]`, "live LSM state"},
		{"camel bond flag", "delegations", `[{"validatorBond":true}]`, "live LSM state"},
		{"total", "total_liquid_staked_tokens", `"159952000"`, "live LSM state"},
		{"liquid", "validators", `[{"liquid_shares":"1"}]`, "live LSM state"},
		{"bond shares", "validators", `[{"validator_bond_shares":"1"}]`, "live LSM state"},
		{"bond flag", "delegations", `[{"validator_bond":true}]`, "live LSM state"},
		{"unknown", "unrecognized_new_field", `true`, "unknown field"},
		{"unknown param", "params", `{"unrecognized_new_field":true}`, "unknown field"},
		{"unknown validator", "validators", `[{"unrecognized_new_field":true}]`, "unknown field"},
		{"invalid total", "total_liquid_staked_tokens", `"garbage"`, "total_liquid_staked_tokens"},
		{"invalid bond", "delegations", `[{"validator_bond":"false"}]`, "validator_bond"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var state map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(legacy, &state))
			state[test.field] = json.RawMessage(test.value)
			input, err := json.Marshal(state)
			require.NoError(t, err)
			require.ErrorContains(t, compat.ValidateGenesis(encoding.Codec, encoding.TxConfig, input), test.want)
			if strings.Contains(test.want, "live LSM") {
				_, err := lsmgenesis.NormalizeStakingGenesis(input)
				require.Error(t, err)
				require.PanicsWithError(t, err.Error(), func() {
					compat.InitGenesis(sdk.Context{}, encoding.Codec, input)
				})
			}
		})
	}
}

func TestLegacyStakingGenesisPreservesOrdinaryState(t *testing.T) {
	input := json.RawMessage(`{
		"validators":[{"min_self_delegation":"100", "tokens":"42", "delegator_shares":"42.0", "liquid_shares":"0.000000000000000000", "validator_bond_shares":"0"}],
		"delegations":[{"delegator_address":"delegator", "validator_address":"validator", "shares":"42.0", "validator_bond":false}],
		"unbonding_delegations":[], "redelegations":[], "exported":true
	}`)
	converted, err := lsmgenesis.NormalizeStakingGenesis(input)
	require.NoError(t, err)
	require.JSONEq(t, `{
		"validators":[{"min_self_delegation":"0", "tokens":"42", "delegator_shares":"42.0"}],
		"delegations":[{"delegator_address":"delegator", "validator_address":"validator", "shares":"42.0"}],
		"unbonding_delegations":[], "redelegations":[], "exported":true
	}`, string(converted))
	native := json.RawMessage(`{"validators":[{"min_self_delegation":"100"}]}`)
	unchanged, err := lsmgenesis.NormalizeStakingGenesis(native)
	require.NoError(t, err)
	require.Equal(t, native, unchanged)
}

func TestLegacyStakingGenesisAppInitialization(t *testing.T) {
	app := irisapp.NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(pubKey, 1)})
	account := authtypes.NewBaseAccountWithAddress(sdk.AccAddress(pubKey.Address()))
	genesis, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), app.DefaultGenesis(), valSet, []authtypes.GenesisAccount{account})
	require.NoError(t, err)
	genesis[stakingtypes.ModuleName] = legacyStakingGenesis(t, genesis[stakingtypes.ModuleName])
	// This is the BasicManager entry point used by genesis validate-genesis.
	require.NoError(t, app.BasicManager().ValidateGenesis(app.AppCodec(), app.EncodingConfig().TxConfig, genesis))
	data, err := json.Marshal(genesis)
	require.NoError(t, err)
	_, err = app.InitChain(&abci.RequestInitChain{
		ConsensusParams: simtestutil.DefaultConsensusParams,
		AppStateBytes:   data,
	})
	require.NoError(t, err)
}
