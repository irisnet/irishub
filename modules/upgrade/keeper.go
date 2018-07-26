package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	storeKey   		sdk.StoreKey
	cdc        		*wire.Codec
	coinKeeper 		bank.Keeper
	// The ValidatorSet to get information about validators
	vs              sdk.ValidatorSet
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, ck bank.Keeper, ds sdk.DelegationSet) Keeper {
	keeper := Keeper {
		storeKey:   key,
		cdc:        cdc,
		coinKeeper: ck,
		vs:        ds.GetValidatorSet(),
	}
	return keeper
}

var (
	defaultSwichPeriod     int64 = 57600	// 2 days
)

func (k Keeper) GetCurrentVersionID() int64 {
	return -1	// return -1 if current version not found
}

func (k Keeper) IncreaseCurrentVersionID() int64 {
	return 0
}

func (k Keeper) GetCurrentVersion() Version {
	return  Version{}
}

func (k Keeper) AddNewVersion(version Version) {
	version.Id = k.IncreaseCurrentVersionID()

}

func (k Keeper) GetVersion(blockHeight int64) Version {
	return  Version{}
}

func (k Keeper) GetVersionList() VersionList {
	return  nil
}

func (k Keeper) GetMsgTypeInCurrentVersion(msg sdk.Msg) (string, sdk.Error) {
	currentVersion := k.GetCurrentVersion()
	return currentVersion.getMsgType(msg)
}
