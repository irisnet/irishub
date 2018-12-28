package params

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"strings"
)

var ParamSetMapping = make(map[string]ParamSet)

func RegisterParamSet(ps ...ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			if _, ok := ParamSetMapping[ps.GetParamsKey()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamsKey()))
			}
			ParamSetMapping[ps.GetParamsKey()] = ps
		}
	}
}

func GetParamKey(keystr string) string {
	strs := strings.Split(keystr, "/")
	if len(strs) != 2 {
		return ""
	}
	return strs[1]
}

func GetParamSpaceFromKey(keystr string) string {
	strs := strings.Split(keystr, "/")
	if len(strs) != 2 {
		return ""
	}
	return strs[0]
}

var ParamMapping = make(map[string]GovParameter)

func RegisterGovParamMapping(gps ...GovParameter) {
	for _, gp := range gps {
		if gp != nil {
			ParamMapping[GovParamspace+"/"+string(gp.GetStoreKey())] = gp
		}
	}
}

func InitGenesisParameter(p Parameter, ctx sdk.Context, genesisData interface{}) {
	if p != nil {
		find := p.LoadValue(ctx)
		if !find {
			p.InitGenesis(genesisData)
			p.SaveValue(ctx)
		}
	}
}

func SetParamReadWriter(paramSpace Subspace, ps ...Parameter) {
	for _, p := range ps {
		if p != nil {
			p.SetReadWriter(paramSpace)
		}
	}
}
