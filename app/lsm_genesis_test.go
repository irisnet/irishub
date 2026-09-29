package app

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	lsmgenesis "github.com/irisnet/irishub/v5/app/genesis/lsm"
)

type lsmFixture struct {
	state    map[string]json.RawMessage
	app      *IrisApp
	ctx      sdk.Context
	owners   []sdk.AccAddress
	records  []lsmgenesis.LegacyRecord
	val      sdk.ValAddress
	expected map[string]sdk.Coins
}

func legacyStakingGenesis(t *testing.T, current json.RawMessage) json.RawMessage {
	t.Helper()
	var state map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(current, &state))
	var params map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(state["params"], &params))
	params["validator_bond_factor"] = json.RawMessage(`"250.000000000000000000"`)
	params["global_liquid_staking_cap"] = json.RawMessage(`"0.250000000000000000"`)
	params["validator_liquid_staking_cap"] = json.RawMessage(`"0.500000000000000000"`)
	state["params"], _ = json.Marshal(params)
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

func newLSMFixture(t *testing.T) lsmFixture {
	t.Helper()
	app := NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
	privVal := mock.NewPV()
	pub, err := privVal.GetPubKey()
	require.NoError(t, err)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{tmtypes.NewValidator(pub, 1)})
	var accounts []authtypes.GenesisAccount
	var owners []sdk.AccAddress
	for i := 0; i < 2; i++ {
		key := secp256k1.GenPrivKey().PubKey()
		owner := sdk.AccAddress(key.Address())
		accounts = append(accounts, authtypes.NewBaseAccount(owner, key, uint64(i), 0))
		owners = append(owners, owner)
	}
	state, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), app.DefaultGenesis(), valSet, accounts)
	require.NoError(t, err)
	data, err := json.Marshal(state)
	require.NoError(t, err)
	_, err = app.InitChain(&abci.RequestInitChain{AppStateBytes: data, ConsensusParams: simtestutil.DefaultConsensusParams})
	require.NoError(t, err)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 10, Time: time.Unix(1000, 0)})
	val := sdk.ValAddress(pub.Address())
	denom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	fund := func(addr sdk.AccAddress, coins sdk.Coins) {
		require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins))
		require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins))
	}
	delegate := func(addr sdk.AccAddress, amount math.Int) math.LegacyDec {
		validator, err := app.StakingKeeper.GetValidator(ctx, val)
		require.NoError(t, err)
		shares, err := app.StakingKeeper.Delegate(ctx, addr, amount, stakingtypes.Unbonded, validator, true)
		require.NoError(t, err)
		return shares
	}
	// Both targets already delegate to the same validator, as on mainnet.
	amount := sdk.DefaultPowerReduction.MulRaw(7)
	fund(owners[1], sdk.NewCoins(sdk.NewCoin(denom, amount)))
	delegate(owners[1], amount)
	var records []lsmgenesis.LegacyRecord
	for i, units := range []int64{100, 10, 50} {
		owner := owners[i%2]
		record := lsmgenesis.LegacyRecord{ID: uint64(i + 1), Owner: owner.String(), ModuleAccount: fmt.Sprintf("tokenizeshare_%d", i+1), Validator: val.String()}
		source := sdk.AccAddress(address.Module("lsm", []byte(record.ModuleAccount)))
		account := app.AccountKeeper.NewAccount(ctx, authtypes.NewBaseAccountWithAddress(source))
		app.AccountKeeper.SetAccount(ctx, account)
		// Fractional shares exercise preservation beyond the integer token supply.
		amount := sdk.DefaultPowerReduction.MulRaw(units).Add(sdk.DefaultPowerReduction.QuoRaw(2))
		fund(source, sdk.NewCoins(sdk.NewCoin(denom, amount.AddRaw(int64(i+1)))))
		shares := delegate(source, amount)
		fund(owner, sdk.NewCoins(sdk.NewCoin(val.String()+fmt.Sprintf("/%d", i+1), shares.TruncateInt())))
		records = append(records, record)
	}
	ctx = ctx.WithBlockHeight(20)
	validator, err := app.StakingKeeper.GetValidator(ctx, val)
	require.NoError(t, err)
	validator.MinSelfDelegation = math.NewInt(123)
	require.NoError(t, app.StakingKeeper.SetValidator(ctx, validator))
	rewards := sdk.NewCoins(sdk.NewInt64Coin(denom, 33333), sdk.NewInt64Coin("weris", 77777))
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, rewards))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, distrtypes.ModuleName, rewards))
	require.NoError(t, app.DistrKeeper.AllocateTokensToValidator(ctx, validator, sdk.NewDecCoinsFromCoins(rewards...)))

	expected := make(map[string]sdk.Coins)
	queryRewards := func(addr sdk.AccAddress) sdk.Coins {
		cache, _ := ctx.CacheContext()
		response, err := distrkeeper.NewQuerier(app.DistrKeeper).DelegationRewards(cache, &distrtypes.QueryDelegationRewardsRequest{DelegatorAddress: addr.String(), ValidatorAddress: val.String()})
		require.NoError(t, err)
		coins, _ := response.Rewards.TruncateDecimal()
		require.True(t, coins.IsAllPositive())
		return coins
	}
	for _, owner := range owners {
		expected[owner.String()] = app.BankKeeper.GetAllBalances(ctx, owner).Add(queryRewards(owner)...)
	}
	for _, record := range records {
		source := sdk.AccAddress(address.Module("lsm", []byte(record.ModuleAccount)))
		expected[record.Owner] = expected[record.Owner].Add(queryRewards(source)...).Add(app.BankKeeper.GetAllBalances(ctx, source)...)
		derivative := sdk.NewCoin(record.Validator+fmt.Sprintf("/%d", record.ID), app.BankKeeper.GetSupply(ctx, record.Validator+fmt.Sprintf("/%d", record.ID)).Amount)
		expected[record.Owner] = expected[record.Owner].Sub(derivative)
	}
	state, err = app.mm.ExportGenesisForModules(ctx, app.AppCodec(), nil)
	require.NoError(t, err)
	legacy := legacyStakingGenesis(t, state[stakingtypes.ModuleName])
	var stakingJSON map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(legacy, &stakingJSON))
	stakingJSON["tokenize_share_records"], err = json.Marshal(records)
	require.NoError(t, err)
	stakingJSON["total_liquid_staked_tokens"] = json.RawMessage(`"159952000"`)
	state[stakingtypes.ModuleName], err = json.Marshal(stakingJSON)
	require.NoError(t, err)
	return lsmFixture{state, app, ctx, owners, records, val, expected}
}

func TestLSMGenesisConvertsPrincipalAndRewards(t *testing.T) {
	fixture := newLSMFixture(t)
	app := NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
	require.NoError(t, ValidateGenesis(app.AppCodec(), app.EncodingConfig().TxConfig, app.BasicManager(), fixture.state))
	prepared, plan, err := lsmgenesis.Prepare(app.AppCodec(), fixture.state)
	require.NoError(t, err)
	require.Len(t, plan, 3)
	normalized, emptyPlan, err := lsmgenesis.Prepare(app.AppCodec(), prepared)
	require.NoError(t, err)
	require.Empty(t, emptyPlan)
	require.Equal(t, prepared, normalized)
	data, err := json.Marshal(fixture.state)
	require.NoError(t, err)
	_, err = app.InitChain(&abci.RequestInitChain{AppStateBytes: data, ConsensusParams: simtestutil.DefaultConsensusParams, InitialHeight: 21})
	require.NoError(t, err)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 21})
	oldValidator, err := fixture.app.StakingKeeper.GetValidator(fixture.ctx, fixture.val)
	require.NoError(t, err)
	validator, err := app.StakingKeeper.GetValidator(ctx, fixture.val)
	require.NoError(t, err)
	require.Equal(t, oldValidator.Tokens, validator.Tokens)
	require.Equal(t, oldValidator.DelegatorShares, validator.DelegatorShares)
	require.True(t, validator.MinSelfDelegation.IsZero())
	for _, owner := range fixture.owners {
		expectedShares, err := fixture.app.StakingKeeper.GetDelegation(fixture.ctx, owner, fixture.val)
		require.NoError(t, err)
		for _, record := range plan {
			if record.Owner == owner.String() {
				expectedShares.Shares = expectedShares.Shares.Add(record.Shares)
			}
		}
		actual, err := app.StakingKeeper.GetDelegation(ctx, owner, fixture.val)
		require.NoError(t, err)
		require.Equal(t, expectedShares.Shares, actual.Shares)
		require.Equal(t, fixture.expected[owner.String()], app.BankKeeper.GetAllBalances(ctx, owner))
		// Merged delegations retain usable F1 reward state.
		_, err = app.DistrKeeper.WithdrawDelegationRewards(ctx, owner, fixture.val)
		require.NoError(t, err)
	}
	for _, record := range plan {
		_, err := app.StakingKeeper.GetDelegation(ctx, record.Source, record.ValAddress)
		require.ErrorIs(t, err, stakingtypes.ErrNoDelegation)
		has, err := app.DistrKeeper.HasDelegatorStartingInfo(ctx, record.ValAddress, record.Source)
		require.NoError(t, err)
		require.False(t, has)
		require.Empty(t, app.BankKeeper.GetAllBalances(ctx, record.Source))
		require.True(t, app.BankKeeper.GetSupply(ctx, record.Derivative.Denom).IsZero())
	}
	for _, pool := range []string{stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName} {
		addr := authtypes.NewModuleAddress(pool)
		require.Equal(t, fixture.app.BankKeeper.GetAllBalances(fixture.ctx, addr), app.BankKeeper.GetAllBalances(ctx, addr))
	}
	fixture.app.BankKeeper.IterateTotalSupply(fixture.ctx, func(coin sdk.Coin) bool {
		for _, record := range plan {
			if coin.Denom == record.Derivative.Denom {
				return false
			}
		}
		require.Equal(t, coin, app.BankKeeper.GetSupply(ctx, coin.Denom))
		return false
	})
	// Independently recompute bank supply and staking shares from the resulting state.
	bank := app.BankKeeper.ExportGenesis(ctx)
	sum := sdk.NewCoins()
	for _, balance := range bank.Balances {
		sum = sum.Add(balance.Coins...)
	}
	require.Equal(t, bank.Supply, sum)
	delegations, err := app.StakingKeeper.GetAllDelegations(ctx)
	require.NoError(t, err)
	shares := math.LegacyZeroDec()
	for _, delegation := range delegations {
		shares = shares.Add(delegation.Shares)
	}
	require.Equal(t, validator.DelegatorShares, shares)
}

func TestLSMGenesisRejectsUnsupportedHoldings(t *testing.T) {
	fixture := newLSMFixture(t)
	for _, name := range []string{"split", "different owner", "missing records", "bad shares", "withdraw address", "unsigned holder"} {
		t.Run(name, func(t *testing.T) {
			state := make(map[string]json.RawMessage, len(fixture.state))
			for key, data := range fixture.state {
				state[key] = append(json.RawMessage(nil), data...)
			}
			var bank banktypes.GenesisState
			require.NoError(t, fixture.app.AppCodec().UnmarshalJSON(state[banktypes.ModuleName], &bank))
			var stakingJSON map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(state[stakingtypes.ModuleName], &stakingJSON))
			switch name {
			case "split", "different owner":
				coin := sdk.NewInt64Coin(fixture.val.String()+"/1", 1)
				for i := range bank.Balances {
					if bank.Balances[i].Address == fixture.owners[0].String() {
						if name == "different owner" {
							coin.Amount = bank.Balances[i].Coins.AmountOf(coin.Denom)
						}
						bank.Balances[i].Coins = bank.Balances[i].Coins.Sub(coin)
					}
				}
				for i := range bank.Balances {
					if bank.Balances[i].Address == fixture.owners[1].String() {
						bank.Balances[i].Coins = bank.Balances[i].Coins.Add(coin)
					}
				}
			case "missing records":
				stakingJSON["tokenize_share_records"] = json.RawMessage(`[]`)
			case "bad shares":
				var delegations []map[string]json.RawMessage
				require.NoError(t, json.Unmarshal(stakingJSON["delegations"], &delegations))
				for _, delegation := range delegations {
					delegation["shares"] = json.RawMessage(`"0.5"`)
				}
				stakingJSON["delegations"], _ = json.Marshal(delegations)
			case "withdraw address":
				var distribution distrtypes.GenesisState
				require.NoError(t, fixture.app.AppCodec().UnmarshalJSON(state[distrtypes.ModuleName], &distribution))
				distribution.DelegatorWithdrawInfos = append(distribution.DelegatorWithdrawInfos, distrtypes.DelegatorWithdrawInfo{DelegatorAddress: fixture.owners[0].String(), WithdrawAddress: fixture.owners[1].String()})
				state[distrtypes.ModuleName] = fixture.app.AppCodec().MustMarshalJSON(&distribution)
			case "unsigned holder":
				var auth authtypes.GenesisState
				require.NoError(t, fixture.app.AppCodec().UnmarshalJSON(state[authtypes.ModuleName], &auth))
				accounts, err := authtypes.UnpackAccounts(auth.Accounts)
				require.NoError(t, err)
				for _, account := range accounts {
					if account.GetAddress().Equals(fixture.owners[0]) {
						require.NoError(t, account.SetPubKey(nil))
					}
				}
				auth.Accounts, err = authtypes.PackAccounts(accounts)
				require.NoError(t, err)
				state[authtypes.ModuleName] = fixture.app.AppCodec().MustMarshalJSON(&auth)
			}
			state[banktypes.ModuleName] = fixture.app.AppCodec().MustMarshalJSON(&bank)
			state[stakingtypes.ModuleName], _ = json.Marshal(stakingJSON)
			require.Error(t, ValidateGenesis(fixture.app.AppCodec(), fixture.app.EncodingConfig().TxConfig, fixture.app.BasicManager(), state))
			app := NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
			data, err := json.Marshal(state)
			require.NoError(t, err)
			_, err = app.InitChain(&abci.RequestInitChain{AppStateBytes: data, ConsensusParams: simtestutil.DefaultConsensusParams, InitialHeight: 21})
			require.Error(t, err)
		})
	}
}

func TestLSMGenesisDeterminismAndAtomicFailure(t *testing.T) {
	fixture := newLSMFixture(t)
	var first map[string]json.RawMessage
	for pass := 0; pass < 2; pass++ {
		state := make(map[string]json.RawMessage, len(fixture.state))
		for key, raw := range fixture.state {
			state[key] = raw
		}
		if pass == 1 {
			var stakingJSON map[string]json.RawMessage
			require.NoError(t, json.Unmarshal(state[stakingtypes.ModuleName], &stakingJSON))
			records := append([]lsmgenesis.LegacyRecord(nil), fixture.records...)
			records[0], records[2] = records[2], records[0]
			stakingJSON["tokenize_share_records"], _ = json.Marshal(records)
			state[stakingtypes.ModuleName], _ = json.Marshal(stakingJSON)
			var bank banktypes.GenesisState
			require.NoError(t, fixture.app.AppCodec().UnmarshalJSON(state[banktypes.ModuleName], &bank))
			for i, j := 0, len(bank.Balances)-1; i < j; i, j = i+1, j-1 {
				bank.Balances[i], bank.Balances[j] = bank.Balances[j], bank.Balances[i]
			}
			state[banktypes.ModuleName] = fixture.app.AppCodec().MustMarshalJSON(&bank)
		}
		app := NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
		data, err := json.Marshal(state)
		require.NoError(t, err)
		_, err = app.InitChain(&abci.RequestInitChain{AppStateBytes: data, ConsensusParams: simtestutil.DefaultConsensusParams, InitialHeight: 21})
		require.NoError(t, err)
		ctx := app.NewContextLegacy(false, tmproto.Header{Height: 21})
		after, err := app.mm.ExportGenesisForModules(ctx, app.AppCodec(), []string{banktypes.ModuleName, stakingtypes.ModuleName, distrtypes.ModuleName})
		require.NoError(t, err)
		if pass == 0 {
			first = after
		} else {
			require.Equal(t, first, after)
		}
	}
	app := NewIrisApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.AppOptionsMap{})
	firstRecordReachedBurn := false
	app.BankKeeper.AppendSendRestriction(func(_ context.Context, _, to sdk.AccAddress, coins sdk.Coins) (sdk.AccAddress, error) {
		if coins.AmountOf(fixture.val.String() + "/1").IsPositive() {
			firstRecordReachedBurn = true
		}
		if coins.AmountOf(fixture.val.String() + "/2").IsPositive() {
			return nil, fmt.Errorf("injected second-record send failure")
		}
		return to, nil
	})
	data, err := json.Marshal(fixture.state)
	require.NoError(t, err)
	_, err = app.InitChain(&abci.RequestInitChain{AppStateBytes: data, ConsensusParams: simtestutil.DefaultConsensusParams, InitialHeight: 21})
	require.ErrorContains(t, err, "injected second-record send failure")
	require.True(t, firstRecordReachedBurn)
	ctx := app.NewContextLegacy(false, tmproto.Header{Height: 21})
	// Both module initialization and the already-converted first record roll back.
	require.Nil(t, app.AccountKeeper.GetAccount(ctx, fixture.owners[0]))
	require.True(t, app.BankKeeper.GetSupply(ctx, "weris").IsZero())
	_, err = app.StakingKeeper.GetValidator(ctx, fixture.val)
	require.ErrorIs(t, err, stakingtypes.ErrNoValidatorFound)
}
