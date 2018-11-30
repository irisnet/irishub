package guardian

import (
	sdk "github.com/irisnet/irishub/types"
)

type Profiler struct {
	Name      string         `json:"name"`
	Addr      sdk.AccAddress `json:"addr"`
	AddedAddr sdk.AccAddress `json:"added_addr"`
}

func NewProfiler(addr, addedAddr sdk.AccAddress) Profiler {
	return Profiler{
		Addr:      addr,
		AddedAddr: addedAddr,
	}
}
func ProfilerEqual(profilerA, profilerB Profiler) bool {
	return profilerA.Addr.Equals(profilerB.Addr) &&
		profilerA.AddedAddr.Equals(profilerB.AddedAddr) &&
		profilerA.Name == profilerB.Name
}
