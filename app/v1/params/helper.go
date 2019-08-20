package params

import (
	"fmt"
	"strings"
)

func RegisterParamSet(paramSets map[string]ParamSet, ps ...ParamSet) {
	for _, ps := range ps {
		if ps != nil {
			if _, ok := paramSets[ps.GetParamSpace()]; ok {
				panic(fmt.Sprintf("<%s> already registered ", ps.GetParamSpace()))
			}
			paramSets[ps.GetParamSpace()] = ps
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
