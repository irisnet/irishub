package upgrade

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// The ValidatorSet to get information about validators
	protocolKeeper sdk.ProtocolKeeper
	sk             stake.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, protocolKeeper sdk.ProtocolKeeper, sk stake.Keeper) Keeper {
	keeper := Keeper{
		key,
		cdc,
		protocolKeeper,
		sk,
	}
	return keeper
}

func (k Keeper) AddNewVersionInfo(ctx sdk.Context, versionInfo VersionInfo) {
	kvStore := ctx.KVStore(k.storeKey)

	versionInfoBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetProposalIDKey(versionInfo.UpgradeInfo.ProposalID), versionInfoBytes)

	proposalIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo.UpgradeInfo.ProposalID)
	if err != nil {
		panic(err)
	}

	if versionInfo.Success {
		kvStore.Set(GetSuccessVersionKey(versionInfo.UpgradeInfo.Protocol.Version), proposalIDBytes)
	} else {
		kvStore.Set(GetFailedVersionKey(versionInfo.UpgradeInfo.Protocol.Version,versionInfo.UpgradeInfo.ProposalID), proposalIDBytes)
	}
}

func (k Keeper) SetSignal(ctx sdk.Context, protocol uint64, address string) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes, err := k.cdc.MarshalBinaryLengthPrefixed(true)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetSignalKey(protocol, address), cmsgBytes)
}

func (k Keeper) GetSignal(ctx sdk.Context, protocol uint64, address string) bool {
	kvStore := ctx.KVStore(k.storeKey)
	flagBytes := kvStore.Get(GetSignalKey(protocol, address))
	if flagBytes != nil {
		var flag bool
		err := k.cdc.UnmarshalBinaryLengthPrefixed(flagBytes, &flag)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

func (k Keeper) DeleteSignal(ctx sdk.Context, protocol uint64, address string) {
	if ok := k.GetSignal(ctx, protocol, address); ok {
		kvStore := ctx.KVStore(k.storeKey)
		kvStore.Delete(GetSignalKey(protocol, address))
	}
}