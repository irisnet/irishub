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
	curVersionBytes := kvStore.Get(GetCurrentVersionKey())
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

func (k Keeper) AddNewVersion(ctx sdk.Context, version Version) {
	kvStore := ctx.KVStore(k.storeKey)
	curVersion := k.GetCurrentVersion(ctx)

	version.Id = curVersion.Id + 1
	versionBytes,err := k.cdc.MarshalBinary(version)
	if err != nil {
		panic(err)
	}

	kvStore.Set(GetVersionIDKey(version.Id),versionBytes)
	kvStore.Set(GetProposalIDKey(version.proposalID),versionBytes)
	kvStore.Set(GetStartHeightKey(version.Start),versionBytes)
}

func (k Keeper) GetVersionByHeight(ctx sdk.Context, blockHeight int64) *Version {
	kvStore := ctx.KVStore(k.storeKey)
	iterator := kvStore.ReverseIterator(GetStartHeightKey(blockHeight),GetStartHeightKey(0))
	defer iterator.Close()

	if iterator.Valid() {
		versionBytes := iterator.Value()
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
	curVersionBytes := kvStore.Get(GetProposalIDKey(proposalId))
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
	iterator := kvStore.Iterator(GetProposalIDKey(0),GetProposalIDKey(math.MaxInt64))
	defer iterator.Close()

	var versionList VersionList
	if iterator.Valid() {
		versionBytes := iterator.Value()
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

func (k Keeper) SetSwitch(propsalID int64, address sdk.AccAddress,cmsg MsgSwitch) {

}

func (k Keeper) GetSwitch(propsalID int64, address sdk.AccAddress) (MsgSwitch, bool) {
	return MsgSwitch{}, true
}
