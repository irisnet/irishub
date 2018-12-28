package upgrade

import (
	protocol "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/stake"
	sdk "github.com/irisnet/irishub/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	pk       protocol.Keeper
	// The ValidatorSet to get information about validators
	sk stake.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk stake.Keeper, pk protocol.Keeper) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
		pk:       pk,
	}
	return keeper
}

func (k Keeper) AddNewVersion(ctx sdk.Context, appVersion AppVersion) {
	kvStore := ctx.KVStore(k.storeKey)

	appVersionBytes, err := k.cdc.MarshalBinaryLengthPrefixed(appVersion)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetProposalIDKey(appVersion.ProposalID), appVersionBytes)

	proposalIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(appVersion.ProposalID)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetAppVersionKey(appVersion.Protocol.Version, appVersion.ProposalID), proposalIDBytes)

	if appVersion.Success {
		kvStore.Set(GetSuccessAppVersionKey(appVersion.Protocol.Version), proposalIDBytes)
	}
}

func (k Keeper) GetVersionByProposalId(ctx sdk.Context, proposalId uint64) *AppVersion {
	kvStore := ctx.KVStore(k.storeKey)
	versionIDBytes := kvStore.Get(GetProposalIDKey(proposalId))
	if versionIDBytes == nil {
		return nil
	}
	var versionID uint64
	err := k.cdc.UnmarshalBinaryLengthPrefixed(versionIDBytes, &versionID)
	if err != nil {
		panic(err)
	}
	versionBytes := kvStore.Get(GetAppVersionKey(versionID,proposalId))
	if versionBytes != nil {
		var version AppVersion
		err := k.cdc.UnmarshalBinaryLengthPrefixed(versionBytes, &version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
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