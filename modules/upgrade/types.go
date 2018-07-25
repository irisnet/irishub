package upgrade

import (
sdk "github.com/cosmos/cosmos-sdk/types"
)

type ModuleLifeTime struct {
    Start		int64
    End			int64
    Handler 	string
    Store		sdk.KVStoreKey
}

func NewModuleLifeTime(start int64, end	int64, handler string,store sdk.KVStoreKey) ModuleLifeTime {
    return ModuleLifeTime{
        Start:start,
        End:end,
        Handler:handler,
        Store:store,
    }
}

type ModulesLifeTime []ModuleLifeTime

func NewModulesLifeTime() ModulesLifeTime {
    return ModulesLifeTime{}
}

func (m ModulesLifeTime) AddModuleLifeTime(start int64, end	int64, handler string,store sdk.KVStoreKey){
    m = append(m,NewModuleLifeTime(start, end, handler,store))
}


type Version struct {
    Id			int64		// should be equal with corresponding upgradeProposalID
    Start		int64
    Votes       []sdk.AccAddress
    ModuleList	ModulesLifeTime
}

func NewVersion(id int64,start int64,moduleList ModulesLifeTime) Version{
    return Version{
        Id:id,
        Start:start,
        Votes:[]sdk.AccAddress{},
        ModuleList:moduleList,
    }
}

func (v Version) getMsgType() string{
	return " "
}

// run when app star
func (v Version) updateCurrentVersion(moduleList ModulesLifeTime){

}


type VersionList []Version

func NewVersionList() VersionList {
	return VersionList{}
}
