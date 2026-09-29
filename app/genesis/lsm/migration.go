package lsm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// LegacyRecord is the JSON shape of an SDK v0.50 LSM tokenize-share record.
type LegacyRecord struct {
	ID            uint64 `json:"id,string"`
	Owner         string `json:"owner"`
	ModuleAccount string `json:"module_account"`
	Validator     string `json:"validator"`
}

// Record contains a validated tokenize-share record and its conversion inputs.
type Record struct {
	LegacyRecord
	Source     sdk.AccAddress
	Holder     sdk.AccAddress
	ValAddress sdk.ValAddress
	Derivative sdk.Coin
	Shares     math.LegacyDec
}

// Plan is the deterministic conversion prepared from a legacy genesis.
type Plan []Record

// Keepers contains the module keepers required to execute a Plan.
type Keepers struct {
	Bank         bankkeeper.Keeper
	Staking      *stakingkeeper.Keeper
	Distribution distrkeeper.Keeper
	Crisis       *crisiskeeper.Keeper
}

// ValidateGenesis uses the same application-wide LSM preflight as InitChainer.
// A module-local validator cannot check token ownership in bank/auth genesis.
func ValidateGenesis(cdc codec.Codec, config client.TxEncodingConfig, basics module.BasicManager, state map[string]json.RawMessage) error {
	prepared, _, err := Prepare(cdc, state)
	if err != nil {
		return err
	}
	return basics.ValidateGenesis(cdc, config, prepared)
}

// Prepare never mutates the input. Live migration deliberately supports
// only a full-supply, signable BaseAccount holder who also owns the record's
// reward rights. Other custody/split-ownership arrangements need their own plan.
func Prepare(cdc codec.Codec, state map[string]json.RawMessage) (map[string]json.RawMessage, Plan, error) {
	raw, ok := state[stakingtypes.ModuleName]
	if !ok {
		return state, nil, nil
	}
	normalized, legacy, err := normalizeStakingGenesis(raw, true)
	if err != nil || !legacy {
		return state, nil, err
	}
	prepared := make(map[string]json.RawMessage, len(state))
	for name, data := range state {
		prepared[name] = data
	}
	prepared[stakingtypes.ModuleName] = normalized
	var old map[string]json.RawMessage
	if err := json.Unmarshal(raw, &old); err != nil {
		return nil, nil, err
	}
	if _, a := old["tokenize_share_records"]; a {
		if _, b := old["tokenizeShareRecords"]; b {
			return nil, nil, fmt.Errorf("LSM genesis has duplicate tokenize-share record fields")
		}
	}
	recordJSON := old["tokenize_share_records"]
	if recordJSON == nil {
		recordJSON = old["tokenizeShareRecords"]
	}
	var rawRecords []map[string]json.RawMessage
	if recordJSON != nil {
		if err := json.Unmarshal(recordJSON, &rawRecords); err != nil {
			return nil, nil, err
		}
	}
	var records []LegacyRecord
	for _, item := range rawRecords {
		if camel, ok := item["moduleAccount"]; ok {
			if _, duplicate := item["module_account"]; duplicate {
				return nil, nil, fmt.Errorf("LSM record has duplicate module account fields")
			}
			item["module_account"] = camel
			delete(item, "moduleAccount")
		}
		data, err := json.Marshal(item)
		if err != nil {
			return nil, nil, err
		}
		var record LegacyRecord
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&record); err != nil {
			return nil, nil, fmt.Errorf("LSM record: %w", err)
		}
		records = append(records, record)
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })

	var stakingState stakingtypes.GenesisState
	var bankState banktypes.GenesisState
	if err := cdc.UnmarshalJSON(normalized, &stakingState); err != nil {
		return nil, nil, err
	}
	if err := cdc.UnmarshalJSON(state[banktypes.ModuleName], &bankState); err != nil {
		return nil, nil, fmt.Errorf("LSM bank genesis: %w", err)
	}
	balances := make(map[string]sdk.Coins, len(bankState.Balances))
	supplyMap := sdk.NewMapCoins(sdk.NewCoins())
	for _, balance := range bankState.Balances {
		if _, duplicate := balances[balance.Address]; duplicate {
			return nil, nil, fmt.Errorf("LSM bank genesis has duplicate balance address %s", balance.Address)
		}
		if err := balance.Coins.Validate(); err != nil {
			return nil, nil, err
		}
		balances[balance.Address] = balance.Coins
		supplyMap.Add(balance.Coins...)
	}
	supply := supplyMap.ToCoins()
	if !bankState.Supply.Empty() && !bankState.Supply.Equal(supply) {
		return nil, nil, fmt.Errorf("LSM bank genesis supply does not match balances")
	}
	denoms := make(map[string]bool, len(records))
	for _, record := range records {
		denoms[strings.ToLower(record.Validator)+"/"+strconv.FormatUint(record.ID, 10)] = true
	}
	// A dangling LSM denom must not pass just because its record was omitted.
	for _, coin := range supply {
		validator, id, found := strings.Cut(coin.Denom, "/")
		if !found {
			continue
		}
		_, valErr := sdk.ValAddressFromBech32(validator)
		_, idErr := strconv.ParseUint(id, 10, 64)
		if valErr == nil && idErr == nil && !denoms[coin.Denom] {
			return nil, nil, fmt.Errorf("LSM derivative %s has supply but no tokenize-share record", coin.Denom)
		}
	}
	if len(records) == 0 {
		return prepared, nil, nil
	}
	if !stakingState.Exported {
		return nil, nil, fmt.Errorf("live LSM conversion requires exported staking genesis")
	}
	var genutil genutiltypes.GenesisState
	if data := state[genutiltypes.ModuleName]; data != nil {
		if err := cdc.UnmarshalJSON(data, &genutil); err != nil {
			return nil, nil, err
		}
		if len(genutil.GenTxs) != 0 {
			return nil, nil, fmt.Errorf("live LSM conversion does not support genesis transactions")
		}
	}
	var authState authtypes.GenesisState
	var distribution distrtypes.GenesisState
	if err := cdc.UnmarshalJSON(state[authtypes.ModuleName], &authState); err != nil {
		return nil, nil, err
	}
	if err := cdc.UnmarshalJSON(state[distrtypes.ModuleName], &distribution); err != nil {
		return nil, nil, err
	}
	accounts, err := authtypes.UnpackAccounts(authState.Accounts)
	if err != nil {
		return nil, nil, err
	}
	accountByAddress := make(map[string]authtypes.GenesisAccount, len(accounts))
	for _, account := range accounts {
		key := account.GetAddress().String()
		if _, duplicate := accountByAddress[key]; duplicate {
			return nil, nil, fmt.Errorf("LSM auth genesis has duplicate account %s", key)
		}
		accountByAddress[key] = account
	}
	validators := make(map[string]stakingtypes.Validator, len(stakingState.Validators))
	for _, validator := range stakingState.Validators {
		validators[validator.OperatorAddress] = validator
	}
	delegations := make(map[string]stakingtypes.Delegation, len(stakingState.Delegations))
	delegationCounts := make(map[string]int)
	for _, delegation := range stakingState.Delegations {
		key := delegation.DelegatorAddress + "/" + delegation.ValidatorAddress
		if _, duplicate := delegations[key]; duplicate {
			return nil, nil, fmt.Errorf("LSM staking genesis has duplicate delegation %s", key)
		}
		delegations[key] = delegation
		delegationCounts[delegation.DelegatorAddress]++
	}
	startingInfo := make(map[string]bool, len(distribution.DelegatorStartingInfos))
	for _, info := range distribution.DelegatorStartingInfos {
		startingInfo[info.DelegatorAddress+"/"+info.ValidatorAddress] = true
	}
	withdraws := make(map[string]string)
	for _, info := range distribution.DelegatorWithdrawInfos {
		withdraws[info.DelegatorAddress] = info.WithdrawAddress
	}
	var plan Plan
	sources := make(map[string]bool)
	for i, record := range records {
		if record.ID == 0 || (i > 0 && records[i-1].ID == record.ID) || record.ModuleAccount != "tokenizeshare_"+strconv.FormatUint(record.ID, 10) {
			return nil, nil, fmt.Errorf("invalid or duplicate LSM record %d or module account", record.ID)
		}
		owner, err := sdk.AccAddressFromBech32(record.Owner)
		if err != nil {
			return nil, nil, err
		}
		val, err := sdk.ValAddressFromBech32(record.Validator)
		if err != nil {
			return nil, nil, err
		}
		if _, exists := validators[record.Validator]; !exists {
			return nil, nil, fmt.Errorf("LSM record %d validator is missing", record.ID)
		}
		source := sdk.AccAddress(address.Module("lsm", []byte(record.ModuleAccount)))
		sources[source.String()] = true
		account, plain := accountByAddress[record.Owner].(*authtypes.BaseAccount)
		if !plain || len(owner) != 20 || account.GetPubKey() == nil || !bytes.Equal(account.GetPubKey().Address(), owner) {
			return nil, nil, fmt.Errorf("LSM record %d owner must be a signable BaseAccount; escrow, module, contract, and vesting custody require a separate migration", record.ID)
		}
		if _, plain := accountByAddress[source.String()].(*authtypes.BaseAccount); !plain {
			return nil, nil, fmt.Errorf("LSM record %d source must be a BaseAccount", record.ID)
		}
		for _, addr := range []string{source.String(), record.Owner} {
			if withdraw := withdraws[addr]; withdraw != "" && withdraw != addr {
				return nil, nil, fmt.Errorf("LSM record %d requires self reward withdrawal addresses", record.ID)
			}
		}
		denom := strings.ToLower(record.Validator) + "/" + strconv.FormatUint(record.ID, 10)
		amount := supply.AmountOf(denom)
		if !amount.IsPositive() || !balances[record.Owner].AmountOf(denom).Equal(amount) {
			return nil, nil, fmt.Errorf("LSM record %d requires its reward owner to hold the entire derivative supply; split holders or different principal owners require a separate migration", record.ID)
		}
		key := source.String() + "/" + record.Validator
		delegation, found := delegations[key]
		if !found || delegationCounts[source.String()] != 1 || !delegation.Shares.IsPositive() || !delegation.Shares.TruncateInt().Equal(amount) {
			return nil, nil, fmt.Errorf("LSM record %d source delegation does not match derivative supply", record.ID)
		}
		if !startingInfo[key] {
			return nil, nil, fmt.Errorf("LSM record %d source reward starting info is missing", record.ID)
		}
		ownerKey := record.Owner + "/" + record.Validator
		if _, found := delegations[ownerKey]; found && !startingInfo[ownerKey] {
			return nil, nil, fmt.Errorf("LSM record %d owner reward starting info is missing", record.ID)
		}
		plan = append(plan, Record{record, source, owner, val, sdk.NewCoin(denom, amount), delegation.Shares})
	}
	for _, unbonding := range stakingState.UnbondingDelegations {
		if sources[unbonding.DelegatorAddress] {
			return nil, nil, fmt.Errorf("LSM record source has an unbonding delegation")
		}
	}
	for _, redelegation := range stakingState.Redelegations {
		if sources[redelegation.DelegatorAddress] {
			return nil, nil, fmt.Errorf("LSM record source has a redelegation")
		}
	}
	return prepared, plan, nil
}

// Migrate runs after all module genesis has initialized, inside the
// InitChainer cache. Transfer exact staking shares, rather than unbonding/rebonding
// integer tokens, so validator power and every fractional principal claim survive.
func Migrate(ctx sdk.Context, keepers Keepers, records Plan) error {
	if len(records) == 0 {
		return nil
	}
	if err := checkInvariants(ctx, keepers); err != nil {
		return err
	}
	ordinarySupply := sdk.NewCoins()
	keepers.Bank.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
		for _, record := range records {
			if coin.Denom == record.Derivative.Denom {
				return false
			}
		}
		ordinarySupply = append(ordinarySupply, coin)
		return false
	})
	bonded := keepers.Bank.GetAllBalances(ctx, authtypes.NewModuleAddress(stakingtypes.BondedPoolName))
	unbonded := keepers.Bank.GetAllBalances(ctx, authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName))
	strict := stakingtypes.WithStrictWithdraw(ctx)
	for _, record := range records {
		validator, err := keepers.Staking.GetValidator(ctx, record.ValAddress)
		if err != nil {
			return err
		}
		source, err := keepers.Staking.GetDelegation(ctx, record.Source, record.ValAddress)
		if err != nil || !source.Shares.Equal(record.Shares) {
			return fmt.Errorf("LSM record %d source delegation changed during genesis initialization", record.ID)
		}
		if !keepers.Bank.GetBalance(ctx, record.Holder, record.Derivative.Denom).Amount.Equal(record.Derivative.Amount) || !keepers.Bank.GetSupply(ctx, record.Derivative.Denom).Amount.Equal(record.Derivative.Amount) {
			return fmt.Errorf("LSM record %d derivative supply changed during genesis initialization", record.ID)
		}
		for _, addr := range []sdk.AccAddress{record.Source, record.Holder} {
			has, err := keepers.Distribution.HasDelegatorStartingInfo(ctx, record.ValAddress, addr)
			if err != nil {
				return err
			}
			if has {
				info, err := keepers.Distribution.GetDelegatorStartingInfo(ctx, record.ValAddress, addr)
				if err != nil {
					return err
				}
				if ctx.BlockHeight() < 0 || info.Height > uint64(ctx.BlockHeight()) {
					return fmt.Errorf("LSM genesis initial height precedes reward starting height for %s", addr)
				}
			}
		}
		// This hook settles rewards and removes starting info/reference counts.
		// WithdrawDelegationRewards would reinitialize the soon-to-be-removed source.
		if err := keepers.Staking.Hooks().BeforeDelegationSharesModified(strict, record.Source, record.ValAddress); err != nil {
			return fmt.Errorf("LSM record %d source rewards: %w", record.ID, err)
		}
		target, err := keepers.Staking.GetDelegation(ctx, record.Holder, record.ValAddress)
		if err == nil {
			if err := keepers.Staking.Hooks().BeforeDelegationSharesModified(strict, record.Holder, record.ValAddress); err != nil {
				return err
			}
			target.Shares = target.Shares.Add(source.Shares)
		} else if errors.Is(err, stakingtypes.ErrNoDelegation) {
			if err := keepers.Staking.Hooks().BeforeDelegationCreated(strict, record.Holder, record.ValAddress); err != nil {
				return err
			}
			target = stakingtypes.NewDelegation(record.Owner, record.Validator, source.Shares)
		} else {
			return err
		}
		if err := keepers.Staking.RemoveDelegation(ctx, source); err != nil {
			return err
		}
		if err := keepers.Staking.SetDelegation(ctx, target); err != nil {
			return err
		}
		if err := keepers.Staking.Hooks().AfterDelegationModified(strict, record.Holder, record.ValAddress); err != nil {
			return err
		}
		// Reward rights and principal rights coincide in the supported shape.
		// Include pre-existing source dust as the old LSM reward withdrawal did.
		coins := keepers.Bank.GetAllBalances(ctx, record.Source)
		if !coins.Empty() {
			if err := keepers.Bank.SendCoins(ctx, record.Source, record.Holder, coins); err != nil {
				return err
			}
		}
		if err := keepers.Distribution.DeleteDelegatorWithdrawAddr(ctx, record.Source, record.Source); err != nil {
			return err
		}
		if err := keepers.Bank.SendCoinsFromAccountToModule(ctx, record.Holder, stakingtypes.NotBondedPoolName, sdk.NewCoins(record.Derivative)); err != nil {
			return err
		}
		if err := keepers.Bank.BurnCoins(ctx, stakingtypes.NotBondedPoolName, sdk.NewCoins(record.Derivative)); err != nil {
			return err
		}
		current, err := keepers.Staking.GetValidator(ctx, record.ValAddress)
		if err != nil {
			return err
		}
		if !current.Tokens.Equal(validator.Tokens) || !current.DelegatorShares.Equal(validator.DelegatorShares) || !keepers.Bank.GetSupply(ctx, record.Derivative.Denom).IsZero() {
			return fmt.Errorf("LSM record %d failed conservation checks", record.ID)
		}
	}
	afterSupply := sdk.NewCoins()
	keepers.Bank.IterateTotalSupply(ctx, func(coin sdk.Coin) bool {
		afterSupply = append(afterSupply, coin)
		return false
	})
	if !ordinarySupply.Equal(afterSupply) {
		return fmt.Errorf("LSM migration changed ordinary supply or left derivative supply")
	}
	if !bonded.Equal(keepers.Bank.GetAllBalances(ctx, authtypes.NewModuleAddress(stakingtypes.BondedPoolName))) || !unbonded.Equal(keepers.Bank.GetAllBalances(ctx, authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName))) {
		return fmt.Errorf("LSM migration changed staking pool balances")
	}
	if err := checkInvariants(ctx, keepers); err != nil {
		return err
	}
	keepers.Crisis.AssertInvariants(ctx)
	return nil
}
