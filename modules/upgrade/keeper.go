package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/irisnet/irishub/modules/stake"
	"math"
	"fmt"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// The ValidatorSet to get information about validators
	sk stake.Keeper
}

var VersionListCached VersionList

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk stake.Keeper) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
		sk:       sk,
	}
	return keeper
}

func (k Keeper) GetCurrentVersion(ctx sdk.Context) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	return k.GetCurrentVersionByStore(kvStore)
}

func (k Keeper) GetCurrentVersionByStore(kvStore sdk.KVStore) *Version {
	versionIDBytes := kvStore.Get(GetCurrentVersionKey())
	if versionIDBytes != nil {
		var versionID int64
		err := k.cdc.UnmarshalBinaryLengthPrefixed(versionIDBytes, &versionID)
		if err != nil {
			panic(err)
		}
		curVersionBytes := kvStore.Get(GetVersionIDKey(versionID))
		if curVersionBytes == nil {
			return nil
		}
		var version Version
		err = k.cdc.UnmarshalBinaryLengthPrefixed(curVersionBytes, &version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) AddNewVersion(ctx sdk.Context, version Version) {
	kvStore := ctx.KVStore(k.storeKey)
	curVersion := k.GetCurrentVersion(ctx)

	if curVersion == nil {
		version.Id = 0
	} else {
		version.Id = curVersion.Id + 1
		if version.ProposalID == curVersion.ProposalID {
			return
		}
	}

	for _, module := range version.ModuleList {
		module.Start = version.Start
	}

	versionBytes, err := k.cdc.MarshalBinaryLengthPrefixed(version)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetVersionIDKey(version.Id), versionBytes)
	VersionListCached = append(VersionListCached, version)

	versionIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(version.Id)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetCurrentVersionKey(), versionIDBytes)
	kvStore.Set(GetProposalIDKey(version.ProposalID), versionIDBytes)
	kvStore.Set(GetStartHeightKey(version.Start), versionIDBytes)
}

func (k Keeper) GetVersionByHeight(ctx sdk.Context, blockHeight int64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := kvStore.ReverseIterator(GetStartHeightKey(0), GetStartHeightKey(blockHeight+1))
	defer iterator.Close()

	if iterator.Valid() {
		versionIDBytes := iterator.Value()
		if versionIDBytes == nil {
			return nil
		}
		var versionID int64
		err := k.cdc.UnmarshalBinaryLengthPrefixed(versionIDBytes, &versionID)
		if err != nil {
			panic(err)
		}
		versionBytes := kvStore.Get(GetVersionIDKey(versionID))
		if versionBytes == nil {
			return nil
		}
		var version Version
		err = k.cdc.UnmarshalBinaryLengthPrefixed(versionBytes, &version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) GetVersionByProposalId(ctx sdk.Context, proposalId uint64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	versionIDBytes := kvStore.Get(GetProposalIDKey(proposalId))
	if versionIDBytes == nil {
		return nil
	}
	var versionID int64
	err := k.cdc.UnmarshalBinaryLengthPrefixed(versionIDBytes, &versionID)
	if err != nil {
		panic(err)
	}
	versionBytes := kvStore.Get(GetVersionIDKey(versionID))
	if versionBytes != nil {
		var version Version
		err := k.cdc.UnmarshalBinaryLengthPrefixed(versionBytes, &version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) GetVersionByVersionId(versionId int64) *Version {
	len := len(VersionListCached)
	if versionId < 0 || versionId >= int64(len) {
		panic(fmt.Errorf("version id %d doesn't exist", versionId))
	}

	return &(VersionListCached[versionId])
}

func (k Keeper) RefreshVersionList(kvStore sdk.KVStore) {
	VersionListCached = k.GetVersionListByStore(kvStore)
}

func (k Keeper) GetVersionList(ctx sdk.Context) VersionList {
	kvStore := ctx.KVStore(k.storeKey)
	return k.GetVersionListByStore(kvStore)
}

func (k Keeper) GetVersionListByStore(kvStore sdk.KVStore) VersionList {

	iterator := kvStore.Iterator(GetVersionIDKey(0), GetVersionIDKey(math.MaxInt64))
	defer iterator.Close()

	var versionList VersionList
	for iterator.Valid() {
		versionBytes := iterator.Value()
		iterator.Next()
		if versionBytes == nil {
			continue
		}
		var version Version
		err := k.cdc.UnmarshalBinaryLengthPrefixed(versionBytes, &version)
		if err != nil {
			panic(err)
		}
		versionList = append(versionList, version)
	}

	return versionList
}

func (k Keeper) GetMsgTypeInCurrentVersion(ctx sdk.Context, msg sdk.Msg) (string, sdk.Error) {
	currentVersion := k.GetCurrentVersion(ctx)
	return currentVersion.getMsgType(msg)
}

func (k Keeper) SetSwitch(ctx sdk.Context, propsalID uint64, address sdk.AccAddress, cmsg MsgSwitch) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes, err := k.cdc.MarshalBinaryLengthPrefixed(cmsg)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetSwitchKey(propsalID, address), cmsgBytes)
}

func (k Keeper) GetSwitch(ctx sdk.Context, propsalID uint64, address sdk.AccAddress) (MsgSwitch, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes := kvStore.Get(GetSwitchKey(propsalID, address))
	if cmsgBytes != nil {
		var cmsg MsgSwitch
		err := k.cdc.UnmarshalBinaryLengthPrefixed(cmsgBytes, &cmsg)
		if err != nil {
			panic(err)
		}
		return cmsg, true
	}
	return MsgSwitch{}, false
}
