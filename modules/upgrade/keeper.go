package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/irisnet/irishub/modules/gov"
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
	defaultSwichPeriod     int64 = 200
)


func (k Keeper) SetCurrentVersion(version Version,proposal gov.Proposal){

}

func (k Keeper) GetCurrentVersion() VersionList{
	return  nil
}

func (k Keeper) SetVersion(version Version){

}
func (k Keeper) GetVersion() VersionList{
	return  nil
}

func (k Keeper) SetVersionList(versionList VersionList){

}
func (k Keeper) GetVersionList() VersionList{
	return  nil
}