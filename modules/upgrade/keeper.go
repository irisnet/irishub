package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"math"
	"encoding/binary"
		)

const (
	defaultSwichPeriod     int64 = 57600	// 2 days
)

type Keeper struct {
	storeKey   		sdk.StoreKey
	cdc        		*wire.Codec
	coinKeeper 		bank.Keeper
	// The ValidatorSet to get information about validators
	sk              stake.Keeper
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, ck bank.Keeper, sk stake.Keeper) Keeper {
	keeper := Keeper {
		storeKey:   key,
		cdc:        cdc,
		coinKeeper: ck,
		sk:        sk,
	}
	return keeper
}

func (k Keeper) GetCurrentVersion(ctx sdk.Context) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	versionIDBytes := kvStore.Get(GetCurrentVersionKey())
	if versionIDBytes != nil {
		versionID := int64(binary.BigEndian.Uint64(versionIDBytes))
		curVersionBytes := kvStore.Get(GetVersionIDKey(versionID))
		if curVersionBytes == nil {
			return nil
		}
		var version Version
		err := k.cdc.UnmarshalBinary(curVersionBytes,&version)
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
		version.Id =0
	} else {
		version.Id = curVersion.Id + 1
		if version.ProposalID == curVersion.ProposalID {
			return
		}
	}
	versionBytes,err := k.cdc.MarshalBinary(version)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetVersionIDKey(version.Id),versionBytes)

	versionIDBytes := make([]byte,20)
	binary.BigEndian.PutUint64(versionIDBytes,uint64(version.Id))

	kvStore.Set(GetCurrentVersionKey(),versionIDBytes)
	kvStore.Set(GetProposalIDKey(version.ProposalID),versionIDBytes)
	kvStore.Set(GetStartHeightKey(version.Start),versionIDBytes)
}

func (k Keeper) GetVersionByHeight(ctx sdk.Context, blockHeight int64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := kvStore.ReverseIterator(GetStartHeightKey(0),GetStartHeightKey(blockHeight+1))
	defer iterator.Close()

	if iterator.Valid() {
		versionIDBytes := iterator.Value()
		if versionIDBytes == nil {
			return nil
		}
		versionID := int64(binary.BigEndian.Uint64(versionIDBytes))
		versionBytes := kvStore.Get(GetVersionIDKey(versionID))
		if versionBytes == nil  {
			return nil
		}
		var version Version
		err := k.cdc.UnmarshalBinary(versionBytes,&version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) GetVersionByProposalId(ctx sdk.Context, proposalId int64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	versionIDBytes := kvStore.Get(GetProposalIDKey(proposalId))
	if versionIDBytes == nil  {
		return nil
	}
	versionID := int64(binary.BigEndian.Uint64(versionIDBytes))
	versionBytes := kvStore.Get(GetVersionIDKey(versionID))
	if versionBytes != nil {
		var version Version
		err := k.cdc.UnmarshalBinary(versionBytes,&version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) GetVersionByVersionId(ctx sdk.Context, versionId int64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	curVersionBytes := kvStore.Get(GetVersionIDKey(versionId))
	if curVersionBytes != nil {
		var version Version
		err := k.cdc.UnmarshalBinary(curVersionBytes,&version)
		if err != nil {
			panic(err)
		}
		return &version
	}
	return nil
}

func (k Keeper) GetVersionList(ctx sdk.Context) VersionList {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := kvStore.Iterator(GetVersionIDKey(0),GetVersionIDKey(math.MaxInt64))
	defer iterator.Close()

	var versionList VersionList
	for iterator.Valid() {
		versionBytes := iterator.Value()
		iterator.Next()
		if versionBytes == nil {
			continue
		}
		var version Version
		err := k.cdc.UnmarshalBinary(versionBytes,&version)
		if err != nil {
			panic(err)
		}
		versionList = append(versionList, version)
	}

	return versionList
}

func (k Keeper) GetCurrentProposalID(ctx sdk.Context) int64 {
	kvStore := ctx.KVStore(k.storeKey)
	proposalIdBytes := kvStore.Get(GetCurrentProposalIdKey())
	if proposalIdBytes != nil {
		return int64(binary.BigEndian.Uint64(proposalIdBytes))
	}
	return -1
}

func (k Keeper) SetCurrentProposalID(ctx sdk.Context, proposalID int64) {
	kvStore := ctx.KVStore(k.storeKey)
	bytes := make([]byte,16)
	binary.BigEndian.PutUint64(bytes,uint64(proposalID))
	kvStore.Set(GetCurrentProposalIdKey(),bytes)
}

func (k Keeper) GetMsgTypeInCurrentVersion(ctx sdk.Context, msg sdk.Msg) (string, sdk.Error) {
	currentVersion := k.GetCurrentVersion(ctx)
	return currentVersion.getMsgType(msg)
}

func (k Keeper) SetSwitch(ctx sdk.Context ,propsalID int64, address sdk.AccAddress,cmsg MsgSwitch) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes,err := k.cdc.MarshalBinary(cmsg)
	if err != nil {
		panic(err)
	}
	kvStore.Set(GetSwitchKey(propsalID,address),cmsgBytes)
}

func (k Keeper) GetSwitch(ctx sdk.Context ,propsalID int64, address sdk.AccAddress) (MsgSwitch, bool) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes := kvStore.Get(GetSwitchKey(propsalID,address))
	if cmsgBytes != nil {
		var cmsg MsgSwitch
		err := k.cdc.UnmarshalBinary(cmsgBytes,&cmsg)
		if err != nil {
			panic(err)
		}
		return cmsg, true
	}
	return MsgSwitch{}, false
}

func (k Keeper) SetCurrentProposalAcceptHeight(ctx sdk.Context, height int64) {
	kvStore := ctx.KVStore(k.storeKey)
	bytes := make([]byte,16)
	binary.BigEndian.PutUint64(bytes,uint64(height))
	kvStore.Set(GetCurrentProposalAcceptHeightKey(),bytes)
}

func (k Keeper) GetCurrentProposalAcceptHeight(ctx sdk.Context) int64 {
	kvStore := ctx.KVStore(k.storeKey)
	proposalAcceptHeightBytes := kvStore.Get(GetCurrentProposalAcceptHeightKey())
	if proposalAcceptHeightBytes != nil {
		return int64(binary.BigEndian.Uint64(proposalAcceptHeightBytes))
	}
	return -1
}

func (k Keeper) SetDoingSwitch(ctx sdk.Context, doing bool) {
	kvStore := ctx.KVStore(k.storeKey)

	bytes := k.cdc.MustMarshalBinary(doing)
	kvStore.Set(GetDoingSwitchKey(), bytes)
}

func (k Keeper) GetDoingSwitch(ctx sdk.Context) (doing bool) {
	kvStore := ctx.KVStore(k.storeKey)

	bytes := kvStore.Get(GetDoingSwitchKey())
	k.cdc.MustUnmarshalBinary(bytes, doing)

	return
}

func (k Keeper) DoSwitchBegin(ctx sdk.Context) {
	k.SetDoingSwitch(ctx, true)

}
