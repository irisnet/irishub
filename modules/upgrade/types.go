package upgrade

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "math"
)

type ModuleLifeTime struct {
    Start		int64
    End			int64
    Handler 	string
    Store		[]string
}

func NewModuleLifeTime(start int64, end	int64, handler string, store []string) ModuleLifeTime {
    return ModuleLifeTime{
        Start:      start,
        End:        end,
        Handler:    handler,
        Store:      store,
    }
}

type ModuleLifeTimeList []ModuleLifeTime

func NewModuleLifeTimeList() ModuleLifeTimeList {
    return ModuleLifeTimeList{}
}

func (mlist ModuleLifeTimeList) BuildModuleLifeTime(start int64, handler string, store []string) ModuleLifeTimeList {
    return append(mlist, NewModuleLifeTime(start, math.MaxInt64, handler, store))
}

type Version struct {
    Id			int64
    ProposalID  uint64
    Start		int64
    ModuleList	ModuleLifeTimeList
}

func NewVersion(id int64, proposalID uint64, start int64, moduleList ModuleLifeTimeList) Version {
    return Version{
        Id:         id,
        ProposalID: proposalID,
        Start:      start,
        ModuleList: moduleList,
    }
}

func (v Version) getMsgType(msg sdk.Msg) (string, sdk.Error) {
    msgType := msg.Route()

    for _, module := range v.ModuleList {
        if msgType == module.Handler {
            return msgType, nil
        }
    }

    return "", NewError(DefaultCodespace, CodeUnSupportedMsgType, "")
}

type VersionList []Version

func NewVersionList() VersionList {
	return VersionList{}
}

func (m VersionList) AddVersion(v Version) {
	m = append(m,v)
}