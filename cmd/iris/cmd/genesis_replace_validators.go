package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cometbft/cometbft/types"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/irisnet/irishub/v4/app/params"
)

const (
	sourceGenesisFile = "source-genesis-file"
	targetGenesisFile = "target-genesis-file"
)

// ReplaceValidatorsResult reports how the replacement was performed.
type ReplaceValidatorsResult struct {
	// Gentx is true when the source validators come from collected gentxs
	// instead of exported staking state.
	Gentx bool
	// Validators counts the source validators, or gentxs when Gentx is true.
	Validators int
}

// replaceValidatorsCmd returns the `genesis replace-validators` command, a
// simulation-test helper that swaps the validator set of a target genesis
// (typically a mainnet export whose validator keys are controlled by others)
// for validators under the operator's own control, while preserving the
// remaining target genesis data.
func replaceValidatorsCmd(basics module.BasicManager, encodingConfig params.EncodingConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace-validators",
		Short: "Replace the validator set of a genesis with your own validators for simulation testing",
		Long: `Replace the staking and distribution (rewards) validator state of a target
genesis with the one from a source genesis you control, keeping every other
module's target data intact. Typical use: a mainnet export contains real
validators whose consensus keys nobody controls, so the chain cannot produce
blocks locally; with this command the exported data can still be exercised
against your own validators.

The source genesis must define its validators in exactly one of two ways:

  - an exported genesis of a chain you run: staking validators with
    last_validator_powers and distribution reward records, or
  - a genesis with collected gentxs (e.g. the output of "iris testnet"),
    whose create-validator transactions carry your keys.

The staking, distribution and genutil module states are taken from the source,
the bonded/not-bonded pool and distribution module account balances follow the
imported modules, the source's regular (non-module) accounts and balances are
added when missing, and the bank supply is recomputed. Everything else (auth,
bank, gov, mint, slashing, evidence, ibc, token, ...) is kept from the target.

The target validators are removed entirely, so the target genesis must not
contain tokenize-share records: dropping the records would strand the
derivative coins left in bank. Start from a genesis without live
tokenize-share records.

When the source provides gentxs, the output takes the source chain-id (gentx
signatures verify against it) and resets the initial height so genesis-time
signature verification skips account numbers; otherwise the target chain-id
and initial height are kept.

Example:
	iris genesis replace-validators \
		--source-genesis-file ./mytestnet/node0/iris/config/genesis.json \
		--target-genesis-file ./mainnet-export.json \
		--output-genesis-file ./simulation-genesis.json`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sourcePath, err := cmd.Flags().GetString(sourceGenesisFile)
			if err != nil {
				return err
			}
			targetPath, err := cmd.Flags().GetString(targetGenesisFile)
			if err != nil {
				return err
			}
			outputPath, err := cmd.Flags().GetString(outputFile)
			if err != nil {
				return err
			}

			source, err := types.GenesisDocFromFile(sourcePath)
			if err != nil {
				return err
			}
			target, err := types.GenesisDocFromFile(targetPath)
			if err != nil {
				return err
			}

			result, err := ReplaceValidators(encodingConfig.Codec, encodingConfig.TxConfig, basics, source, target, outputPath)
			if err != nil {
				return err
			}

			if result.Gentx {
				cmd.Printf("Replaced the target validator set with %d gentx validator(s); the output uses the source chain-id %q and resets the initial height\n", result.Validators, source.ChainID)
			} else {
				cmd.Printf("Replaced the target validator set with %d exported validator(s); the target chain-id and initial height are kept\n", result.Validators)
			}
			cmd.Printf("Merged genesis saved to %s\n", outputPath)
			return nil
		},
	}
	cmd.Flags().String(sourceGenesisFile, "", "genesis file you control that provides the replacement validators (exported chain genesis or collected-gentxs genesis)")
	cmd.Flags().String(targetGenesisFile, "", "genesis file (e.g. a mainnet export) whose remaining data is preserved")
	cmd.Flags().String(outputFile, "", "path to write the merged genesis file to")
	return cmd
}

// ReplaceValidators swaps the staking and distribution validator state of the
// target genesis with the source one and writes the merged genesis to output.
func ReplaceValidators(
	cdc codec.Codec,
	txConfig client.TxEncodingConfig,
	basics module.BasicManager,
	source, target *types.GenesisDoc,
	output string,
) (*ReplaceValidatorsResult, error) {
	var sourceState, targetState map[string]json.RawMessage
	if err := json.Unmarshal(source.AppState, &sourceState); err != nil {
		return nil, fmt.Errorf("source app state: %w", err)
	}
	if err := json.Unmarshal(target.AppState, &targetState); err != nil {
		return nil, fmt.Errorf("target app state: %w", err)
	}

	// The target validators are removed entirely, so tokenize-share records
	// in the target are dropped along with them while the derivative coins
	// stay in bank. Reject such targets instead of silently stranding the
	// tokenize-share claims.
	liveRecords, err := liveTokenizeShareCount(targetState[stakingtypes.ModuleName])
	if err != nil {
		return nil, err
	}
	if liveRecords != 0 {
		return nil, fmt.Errorf("target genesis contains %d live tokenize-share records whose validators would be removed; start from a genesis without live tokenize-share records", liveRecords)
	}

	sourceStakingRaw, ok := sourceState[stakingtypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("source genesis has no %s module state", stakingtypes.ModuleName)
	}
	var sourceStaking stakingtypes.GenesisState
	if err := cdc.UnmarshalJSON(sourceStakingRaw, &sourceStaking); err != nil {
		return nil, fmt.Errorf("source staking genesis: %w", err)
	}
	sourceGenutilRaw, ok := sourceState[genutiltypes.ModuleName]
	if !ok {
		return nil, fmt.Errorf("source genesis has no %s module state", genutiltypes.ModuleName)
	}
	var sourceGenutil genutiltypes.GenesisState
	if err := cdc.UnmarshalJSON(sourceGenutilRaw, &sourceGenutil); err != nil {
		return nil, fmt.Errorf("source genutil genesis: %w", err)
	}
	if _, ok := sourceState[distrtypes.ModuleName]; !ok {
		return nil, fmt.Errorf("source genesis has no %s module state", distrtypes.ModuleName)
	}

	hasValidators := len(sourceStaking.Validators) != 0
	hasGenTxs := len(sourceGenutil.GenTxs) != 0
	switch {
	case hasValidators && hasGenTxs:
		return nil, fmt.Errorf("source genesis defines validators both in %s state and in gentxs; only one mechanism is supported", stakingtypes.ModuleName)
	case !hasValidators && !hasGenTxs:
		return nil, fmt.Errorf("source genesis defines no validators: %s has no validators and genutil has no gentxs", stakingtypes.ModuleName)
	}

	if err := checkSourceConsistency(cdc, sourceState, &sourceStaking); err != nil {
		return nil, err
	}

	result := &ReplaceValidatorsResult{}
	if hasValidators {
		// Imported staking is restored like a chain export: the validator set
		// comes from last_validator_powers and distribution runs before
		// staking, whose non-exported hooks would reinitialize the imported
		// reward records.
		bondedPower := int64(0)
		for _, power := range sourceStaking.LastValidatorPowers {
			if power.Power > bondedPower {
				bondedPower = power.Power
			}
		}
		if bondedPower <= 0 {
			return nil, fmt.Errorf("source staking genesis has no bonded validator power in last_validator_powers")
		}
		sourceStaking.Exported = true
		stakingRaw, err := cdc.MarshalJSON(&sourceStaking)
		if err != nil {
			return nil, err
		}
		sourceState[stakingtypes.ModuleName] = stakingRaw
		result.Validators = len(sourceStaking.Validators)
	} else {
		result.Gentx = true
		result.Validators = len(sourceGenutil.GenTxs)
	}

	merged := make(map[string]json.RawMessage, len(targetState))
	for name, raw := range targetState {
		merged[name] = raw
	}
	merged[stakingtypes.ModuleName] = sourceState[stakingtypes.ModuleName]
	merged[distrtypes.ModuleName] = sourceState[distrtypes.ModuleName]
	merged[genutiltypes.ModuleName] = sourceGenutilRaw

	if err := mergeValidatorAccounts(cdc, sourceState, merged); err != nil {
		return nil, err
	}

	target.Validators = nil
	if hasGenTxs {
		// Gentx signatures verify against the chain-id they were signed with,
		// and only at height 0 does genesis signature verification skip
		// account numbers (imported accounts conflicting with target accounts
		// are renumbered by x/auth InitGenesis).
		target.ChainID = source.ChainID
		target.InitialHeight = 0
	}

	appState, err := json.Marshal(merged)
	if err != nil {
		return nil, err
	}
	target.AppState = appState

	// Validate before writing so a failed run never leaves an invalid
	// genesis file behind.
	if err := basics.ValidateGenesis(cdc, txConfig, merged); err != nil {
		return nil, fmt.Errorf("merged genesis is invalid: %w", err)
	}
	if err := target.SaveAs(output); err != nil {
		return nil, err
	}
	return result, nil
}

// liveTokenizeShareCount counts the tokenize-share records left in a staking
// genesis state, accepting both protobuf JSON spellings. The records'
// validators are looked up at chain start, so removing the validators strands
// the records' claims.
func liveTokenizeShareCount(stakingRaw json.RawMessage) (int, error) {
	if len(stakingRaw) == 0 {
		return 0, nil
	}
	var state map[string]json.RawMessage
	if err := json.Unmarshal(stakingRaw, &state); err != nil {
		return 0, fmt.Errorf("target staking genesis: %w", err)
	}
	raw, ok := state["tokenize_share_records"]
	if !ok {
		raw, ok = state["tokenizeShareRecords"]
	}
	if !ok {
		return 0, nil
	}
	var records []json.RawMessage
	if err := json.Unmarshal(raw, &records); err != nil {
		return 0, fmt.Errorf("target staking genesis has an invalid tokenize_share_records field: %w", err)
	}
	return len(records), nil
}

// mergeValidatorAccounts adapts the merged auth and bank states to the
// imported staking and distribution modules: the bonded/not-bonded pool and
// distribution module account balances must match the imported module state,
// and the source's regular accounts need their balances to be usable on the
// simulated chain.
func mergeValidatorAccounts(cdc codec.Codec, source, merged map[string]json.RawMessage) error {
	for module, state := range map[string]map[string]json.RawMessage{
		"source": source,
		"target": merged,
	} {
		for _, name := range []string{authtypes.ModuleName, banktypes.ModuleName} {
			if _, ok := state[name]; !ok {
				return fmt.Errorf("%s genesis has no %s module state", module, name)
			}
		}
	}
	var sourceBank, mergedBank banktypes.GenesisState
	if err := cdc.UnmarshalJSON(source[banktypes.ModuleName], &sourceBank); err != nil {
		return fmt.Errorf("source bank genesis: %w", err)
	}
	if err := cdc.UnmarshalJSON(merged[banktypes.ModuleName], &mergedBank); err != nil {
		return fmt.Errorf("target bank genesis: %w", err)
	}
	var sourceAuth, mergedAuth authtypes.GenesisState
	if err := cdc.UnmarshalJSON(source[authtypes.ModuleName], &sourceAuth); err != nil {
		return fmt.Errorf("source auth genesis: %w", err)
	}
	if err := cdc.UnmarshalJSON(merged[authtypes.ModuleName], &mergedAuth); err != nil {
		return fmt.Errorf("target auth genesis: %w", err)
	}

	pools := map[string]bool{
		authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String():    true,
		authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String(): true,
		authtypes.NewModuleAddress(distrtypes.ModuleName).String():          true,
	}

	sourceAccounts, err := authtypes.UnpackAccounts(sourceAuth.Accounts)
	if err != nil {
		return fmt.Errorf("source auth genesis: %w", err)
	}
	sourceModuleAccounts := make(map[string]bool)
	for _, account := range sourceAccounts {
		if _, isModule := account.(*authtypes.ModuleAccount); isModule {
			sourceModuleAccounts[account.GetAddress().String()] = true
		}
	}

	sourceBalances := make(map[string]sdk.Coins, len(sourceBank.Balances))
	for _, balance := range sourceBank.Balances {
		sourceBalances[balance.Address] = balance.Coins
	}

	// Keep every target balance except the staking/distribution pool
	// accounts, which follow the imported module state.
	balances := make([]banktypes.Balance, 0, len(mergedBank.Balances))
	keptAddrs := make(map[string]bool, len(mergedBank.Balances))
	for _, balance := range mergedBank.Balances {
		keptAddrs[balance.Address] = true
		if pools[balance.Address] {
			continue
		}
		balances = append(balances, balance)
	}
	for pool := range pools {
		if coins := sourceBalances[pool]; !coins.IsZero() {
			balances = append(balances, banktypes.Balance{Address: pool, Coins: coins.Sort()})
			keptAddrs[pool] = true
		}
	}
	// Add the source's regular account balances; balances of source module
	// accounts are source-internal and must not leak into the target.
	for _, balance := range sourceBank.Balances {
		addr := balance.Address
		if pools[addr] || sourceModuleAccounts[addr] || keptAddrs[addr] {
			continue
		}
		balances = append(balances, banktypes.Balance{Address: addr, Coins: balance.Coins.Sort()})
		keptAddrs[addr] = true
	}

	mergedBank.Balances = banktypes.SanitizeGenesisBalances(balances)
	supplyMap := sdk.NewMapCoins(sdk.NewCoins())
	for _, balance := range mergedBank.Balances {
		supplyMap.Add(balance.Coins...)
	}
	mergedBank.Supply = supplyMap.ToCoins()

	bankRaw, err := cdc.MarshalJSON(&mergedBank)
	if err != nil {
		return err
	}
	merged[banktypes.ModuleName] = bankRaw

	// Add the source's regular accounts when missing; account numbers
	// conflicting with target accounts are renumbered by x/auth InitGenesis.
	mergedAccounts, err := authtypes.UnpackAccounts(mergedAuth.Accounts)
	if err != nil {
		return fmt.Errorf("target auth genesis: %w", err)
	}
	mergedAddrs := make(map[string]bool, len(mergedAccounts))
	for _, account := range mergedAccounts {
		mergedAddrs[account.GetAddress().String()] = true
	}
	var additions []authtypes.GenesisAccount
	for _, account := range sourceAccounts {
		addr := account.GetAddress().String()
		if sourceModuleAccounts[addr] || mergedAddrs[addr] {
			continue
		}
		additions = append(additions, account)
	}
	if len(additions) != 0 {
		packed, err := authtypes.PackAccounts(additions)
		if err != nil {
			return err
		}
		mergedAuth.Accounts = append(mergedAuth.Accounts, packed...)
		authRaw, err := cdc.MarshalJSON(&mergedAuth)
		if err != nil {
			return err
		}
		merged[authtypes.ModuleName] = authRaw
	}
	return nil
}

// checkSourceConsistency mirrors the x/staking and x/distribution InitGenesis
// balance checks, so an inconsistent source fails here with a clear error
// instead of panicking when the merged chain starts.
func checkSourceConsistency(cdc codec.Codec, sourceState map[string]json.RawMessage, staking *stakingtypes.GenesisState) error {
	var bankState banktypes.GenesisState
	if err := cdc.UnmarshalJSON(sourceState[banktypes.ModuleName], &bankState); err != nil {
		return fmt.Errorf("source bank genesis: %w", err)
	}
	balances := make(map[string]sdk.Coins, len(bankState.Balances))
	for _, balance := range bankState.Balances {
		balances[balance.Address] = balance.Coins
	}

	bonded, notBonded := math.ZeroInt(), math.ZeroInt()
	validatorAddrs := make(map[string]bool, len(staking.Validators))
	for _, validator := range staking.Validators {
		validatorAddrs[validator.OperatorAddress] = true
		switch validator.GetStatus() {
		case stakingtypes.Bonded:
			bonded = bonded.Add(validator.GetTokens())
		case stakingtypes.Unbonding, stakingtypes.Unbonded:
			notBonded = notBonded.Add(validator.GetTokens())
		default:
			return fmt.Errorf("source validator %s has an invalid status", validator.OperatorAddress)
		}
	}
	// x/staking InitGenesis panics on power entries whose validator is missing;
	// a hand-edited source must fail here instead of at chain start.
	poweredAddrs := make(map[string]bool, len(staking.LastValidatorPowers))
	for _, power := range staking.LastValidatorPowers {
		if !validatorAddrs[power.Address] {
			return fmt.Errorf("source last_validator_powers references missing validator %s", power.Address)
		}
		if poweredAddrs[power.Address] {
			return fmt.Errorf("source last_validator_powers has a duplicate entry for %s", power.Address)
		}
		poweredAddrs[power.Address] = true
	}
	for _, unbonding := range staking.UnbondingDelegations {
		for _, entry := range unbonding.Entries {
			notBonded = notBonded.Add(entry.Balance)
		}
	}
	bondDenom := staking.Params.BondDenom
	checkPool := func(name, addr string, amount math.Int) error {
		expected := sdk.NewCoins(sdk.NewCoin(bondDenom, amount))
		if !balances[addr].Equal(expected) {
			return fmt.Errorf("source genesis %s balance is %s, but the imported staking state requires %s", name, balances[addr], expected)
		}
		return nil
	}
	if err := checkPool("bonded pool", authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(), bonded); err != nil {
		return err
	}
	if err := checkPool("not bonded pool", authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String(), notBonded); err != nil {
		return err
	}

	var distrState distrtypes.GenesisState
	if err := cdc.UnmarshalJSON(sourceState[distrtypes.ModuleName], &distrState); err != nil {
		return fmt.Errorf("source distribution genesis: %w", err)
	}
	holdings := distrState.FeePool.CommunityPool
	for _, reward := range distrState.OutstandingRewards {
		holdings = holdings.Add(reward.OutstandingRewards...)
	}
	expectedHoldings, _ := holdings.TruncateDecimal()
	addr := authtypes.NewModuleAddress(distrtypes.ModuleName).String()
	if !balances[addr].Equal(expectedHoldings) {
		return fmt.Errorf("source genesis distribution module account balance is %s, but the imported rewards require %s", balances[addr], expectedHoldings)
	}
	return nil
}
