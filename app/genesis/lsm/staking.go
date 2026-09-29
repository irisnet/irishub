package lsm

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// StakingModule accepts the removed LSM fields only when they describe
// no live LSM economic state. Live tokenized shares require a separate asset
// migration; silently dropping their records would strand their holders' claims.
type StakingModule struct {
	staking.AppModule
}

func (m StakingModule) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, data json.RawMessage) error {
	data, err := NormalizeStakingGenesis(data)
	if err != nil {
		return err
	}
	return m.AppModule.ValidateGenesis(cdc, config, data)
}

func (m StakingModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	data, err := NormalizeStakingGenesis(data)
	if err != nil {
		panic(err)
	}
	return m.AppModule.InitGenesis(ctx, cdc, data)
}

// NormalizeStakingGenesis removes disabled legacy LSM fields with no live state.
func NormalizeStakingGenesis(data json.RawMessage) (json.RawMessage, error) {
	result, _, err := normalizeStakingGenesis(data, false)
	return result, err
}

// allowLive is used only by the application-wide preflight, which checks the
// bank/auth/distribution state before permitting any live records to be removed.
func normalizeStakingGenesis(data json.RawMessage, allowLive bool) (json.RawMessage, bool, error) {
	var state map[string]json.RawMessage
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, false, fmt.Errorf("staking genesis: %w", err)
	}
	changed := false
	remove := func(object map[string]json.RawMessage, key string) {
		if _, ok := object[key]; ok {
			delete(object, key)
			changed = true
		}
	}
	requireZero := func(object map[string]json.RawMessage, key, path string) error {
		value, ok := object[key]
		if !ok {
			return nil
		}
		var text string
		if err := json.Unmarshal(value, &text); err != nil {
			return fmt.Errorf("staking genesis %s: expected a decimal string: %w", path, err)
		}
		number, err := math.LegacyNewDecFromStr(text)
		if err != nil {
			return fmt.Errorf("staking genesis %s: %w", path, err)
		}
		if number.IsNegative() || (!allowLive && !number.IsZero()) {
			return fmt.Errorf("staking genesis %s contains live LSM state; migrate LSM assets before importing", path)
		}
		remove(object, key)
		return nil
	}
	// Accept both protobuf JSON spellings, while leaving every unrelated field
	// for the upstream strict decoder to validate.
	for _, key := range []string{"tokenize_share_records", "tokenizeShareRecords", "tokenize_share_locks", "tokenizeShareLocks"} {
		if raw, ok := state[key]; ok {
			var entries []json.RawMessage
			if err := json.Unmarshal(raw, &entries); err != nil {
				return nil, false, fmt.Errorf("staking genesis %s: %w", key, err)
			}
			if !allowLive && len(entries) != 0 {
				return nil, false, fmt.Errorf("staking genesis %s contains %d live LSM entries; migrate LSM assets before importing", key, len(entries))
			}
			remove(state, key)
		}
	}
	for _, key := range []string{"total_liquid_staked_tokens", "totalLiquidStakedTokens"} {
		if err := requireZero(state, key, key); err != nil {
			return nil, false, err
		}
	}
	remove(state, "last_tokenize_share_record_id")
	remove(state, "lastTokenizeShareRecordId")
	if raw, ok := state["params"]; ok {
		var params map[string]json.RawMessage
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, false, fmt.Errorf("staking genesis params: %w", err)
		}
		for _, key := range []string{"validator_bond_factor", "validatorBondFactor", "global_liquid_staking_cap", "globalLiquidStakingCap", "validator_liquid_staking_cap", "validatorLiquidStakingCap"} {
			remove(params, key)
		}
		state["params"], _ = json.Marshal(params)
	}
	for _, field := range []string{"validators", "delegations"} {
		raw, ok := state[field]
		if !ok {
			continue
		}
		var entries []map[string]json.RawMessage
		if err := json.Unmarshal(raw, &entries); err != nil {
			return nil, false, fmt.Errorf("staking genesis %s: %w", field, err)
		}
		for i, entry := range entries {
			if field == "validators" {
				for _, key := range []string{"liquid_shares", "liquidShares", "validator_bond_shares", "validatorBondShares"} {
					if err := requireZero(entry, key, fmt.Sprintf("%s[%d].%s", field, i, key)); err != nil {
						return nil, false, err
					}
				}
			} else {
				for _, key := range []string{"validator_bond", "validatorBond"} {
					if raw, ok := entry[key]; ok {
						var bonded bool
						if err := json.Unmarshal(raw, &bonded); err != nil {
							return nil, false, fmt.Errorf("staking genesis delegations[%d].%s: %w", i, key, err)
						}
						if bonded && !allowLive {
							return nil, false, fmt.Errorf("staking genesis delegations[%d].%s contains live LSM state; migrate LSM assets before importing", i, key)
						}
						remove(entry, key)
					}
				}
			}
		}
		state[field], _ = json.Marshal(entries)
	}
	if !changed {
		return data, false, nil
	}
	// The LSM SDK ignored the stored threshold. Preserve that effective behavior
	// for imported legacy validators; native v0.53 genesis is left unchanged.
	if raw, ok := state["validators"]; ok {
		var validators []map[string]json.RawMessage
		if err := json.Unmarshal(raw, &validators); err != nil {
			return nil, false, err
		}
		for _, validator := range validators {
			if validator != nil {
				delete(validator, "minSelfDelegation")
				validator["min_self_delegation"] = json.RawMessage(`"0"`)
			}
		}
		state["validators"], _ = json.Marshal(validators)
	}
	result, err := json.Marshal(state)
	return result, true, err
}
