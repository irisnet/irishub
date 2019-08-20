package subspace

import (
	"fmt"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

// Used for associating paramsubspace key and field of param structs
type KeyValuePair struct {
	Key   []byte
	Value interface{}
}

// Slice of KeyFieldPair
type KeyValuePairs []KeyValuePair

// Interface for structs containing parameters for a module
type ParamSet interface {
	KeyValuePairs() KeyValuePairs
	Validate(key string, value string) (interface{}, sdk.Error)
	GetParamSpace() string
	StringFromBytes(*codec.Codec, string, []byte) (string, error)
	String() string
	ReadOnly() bool
}
type ParamSets []ParamSet

func (pss ParamSets) String() string {
	var s = ""
	var p ParamSet
	for _, ps := range pss {
		if !ps.ReadOnly() {
			s += fmt.Sprintf("%s\n", ps.String())
		} else {
			p = ps
		}
	}
	s += fmt.Sprintf("\n%s", p.String())
	return s
}
