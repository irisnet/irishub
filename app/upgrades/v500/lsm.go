package v500

import (
	"encoding/binary"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/protobuf/encoding/protowire"

	"github.com/irisnet/irishub/v5/app/upgrades"
)

// migrateLSM removes the legacy staking extension at the upgrade height. The
// tokenized-share delegations and distribution state remain available to back
// the derivative coins held in bank until a later redemption is agreed upon.
func migrateLSM(ctx sdk.Context, box upgrades.Toolbox) error {
	store := ctx.KVStore(box.GetKey(stakingtypes.StoreKey))
	rewardOwners, err := legacyRewardOwners(ctx, store, box)
	if err != nil {
		return err
	}
	for _, claim := range rewardOwners {
		if err := box.DistrKeeper.SetDelegatorWithdrawAddr(ctx, claim.source, claim.owner); err != nil {
			return fmt.Errorf("preserve tokenize-share reward owner for record %d: %w", claim.id, err)
		}
	}

	// The old LSM fork added protobuf fields that the upstream staking types no
	// longer define. Reading and writing through the upstream codec drops them.
	if store.Has(stakingtypes.ParamsKey) {
		params, err := box.StakingKeeper.GetParams(ctx)
		if err != nil {
			return fmt.Errorf("read staking params: %w", err)
		}
		if err := box.StakingKeeper.SetParams(ctx, params); err != nil {
			return fmt.Errorf("write staking params: %w", err)
		}
	}

	validators, err := box.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return fmt.Errorf("read staking validators: %w", err)
	}
	for _, validator := range validators {
		// The old fork did not enforce stored minimum self delegation values.
		validator.MinSelfDelegation = math.ZeroInt()
		if err := box.StakingKeeper.SetValidator(ctx, validator); err != nil {
			return fmt.Errorf("write validator %s: %w", validator.OperatorAddress, err)
		}
	}

	delegations, err := box.StakingKeeper.GetAllDelegations(ctx)
	if err != nil {
		return fmt.Errorf("read staking delegations: %w", err)
	}
	for _, delegation := range delegations {
		if err := box.StakingKeeper.SetDelegation(ctx, delegation); err != nil {
			return fmt.Errorf("write delegation %s/%s: %w", delegation.DelegatorAddress, delegation.ValidatorAddress, err)
		}
	}
	if err := clearHistoricalLSMFields(store, box); err != nil {
		return err
	}

	// SDK v0.53 reserves these seven old LSM prefixes and does not migrate
	// them. Delete their records and indexes without touching normal stake.
	for prefix := byte(0x81); prefix <= 0x87; prefix++ {
		iterator := storetypes.KVStorePrefixIterator(store, []byte{prefix})
		var keys [][]byte
		for ; iterator.Valid(); iterator.Next() {
			keys = append(keys, append([]byte(nil), iterator.Key()...))
		}
		if err := iterator.Close(); err != nil {
			return fmt.Errorf("close LSM prefix %x iterator: %w", prefix, err)
		}
		for _, key := range keys {
			store.Delete(key)
		}
	}

	return nil
}

type rewardOwner struct {
	id     uint64
	source sdk.AccAddress
	owner  sdk.AccAddress
}

// The removed record is the only on-chain link from a tokenized delegation to
// its reward owner. Carry that link into the standard distribution withdraw map.
func legacyRewardOwners(ctx sdk.Context, store storetypes.KVStore, box upgrades.Toolbox) ([]rewardOwner, error) {
	iterator := storetypes.KVStorePrefixIterator(store, []byte{0x81})
	defer iterator.Close()
	var claims []rewardOwner
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		if len(key) != 9 {
			return nil, fmt.Errorf("invalid tokenize-share record key %x", key)
		}
		record, err := decodeLegacyRecord(iterator.Value())
		if err != nil {
			return nil, fmt.Errorf("tokenize-share record %x: %w", key, err)
		}
		if record.id == 0 || record.id != binary.BigEndian.Uint64(key[1:]) || record.module != "tokenizeshare_"+strconv.FormatUint(record.id, 10) {
			return nil, fmt.Errorf("tokenize-share record %x has invalid id or module account", key)
		}
		ownerBytes, err := box.AccountKeeper.AddressCodec().StringToBytes(record.owner)
		if err != nil || len(ownerBytes) == 0 {
			return nil, fmt.Errorf("tokenize-share record %d has invalid owner %q: %w", record.id, record.owner, err)
		}
		validatorBytes, err := box.StakingKeeper.ValidatorAddressCodec().StringToBytes(record.validator)
		if err != nil {
			return nil, fmt.Errorf("tokenize-share record %d has invalid validator: %w", record.id, err)
		}
		source := sdk.AccAddress(address.Module("lsm", []byte(record.module)))
		owner := sdk.AccAddress(ownerBytes)
		if _, err := box.StakingKeeper.GetDelegation(ctx, source, sdk.ValAddress(validatorBytes)); err != nil {
			return nil, fmt.Errorf("tokenize-share record %d has no backing delegation: %w", record.id, err)
		}
		withdraw, err := box.DistrKeeper.GetDelegatorWithdrawAddr(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("tokenize-share record %d withdrawal address: %w", record.id, err)
		}
		if !withdraw.Equals(source) && !withdraw.Equals(owner) {
			return nil, fmt.Errorf("tokenize-share record %d has an unrelated withdrawal address", record.id)
		}
		claims = append(claims, rewardOwner{record.id, source, owner})
	}
	return claims, nil
}

type legacyRecord struct {
	id        uint64
	owner     string
	module    string
	validator string
}

func decodeLegacyRecord(data []byte) (legacyRecord, error) {
	var record legacyRecord
	var seen [5]bool
	for len(data) > 0 {
		field, wireType, n := protowire.ConsumeTag(data)
		if n < 0 || field < 1 || field > 4 || seen[field] {
			return record, fmt.Errorf("invalid or duplicate protobuf field")
		}
		seen[field] = true
		data = data[n:]
		switch field {
		case 1:
			if wireType != protowire.VarintType {
				return record, fmt.Errorf("invalid id wire type")
			}
			record.id, n = protowire.ConsumeVarint(data)
		case 2, 3, 4:
			if wireType != protowire.BytesType {
				return record, fmt.Errorf("invalid string wire type")
			}
			var value []byte
			value, n = protowire.ConsumeBytes(data)
			switch field {
			case 2:
				record.owner = string(value)
			case 3:
				record.module = string(value)
			case 4:
				record.validator = string(value)
			}
		}
		if n < 0 {
			return record, fmt.Errorf("invalid protobuf value")
		}
		data = data[n:]
	}
	for field := 1; field <= 4; field++ {
		if !seen[field] {
			return record, fmt.Errorf("missing protobuf field %d", field)
		}
	}
	return record, nil
}

func clearHistoricalLSMFields(store storetypes.KVStore, box upgrades.Toolbox) error {
	iterator := storetypes.KVStorePrefixIterator(store, stakingtypes.HistoricalInfoKey)
	var keys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		keys = append(keys, append([]byte(nil), iterator.Key()...))
	}
	if err := iterator.Close(); err != nil {
		return fmt.Errorf("close historical staking iterator: %w", err)
	}
	for _, key := range keys {
		var info stakingtypes.HistoricalInfo
		if err := box.AppCodec.Unmarshal(store.Get(key), &info); err != nil {
			return fmt.Errorf("decode historical staking info %x: %w", key, err)
		}
		value, err := box.AppCodec.Marshal(&info)
		if err != nil {
			return fmt.Errorf("encode historical staking info %x: %w", key, err)
		}
		store.Set(key, value)
	}
	return nil
}
