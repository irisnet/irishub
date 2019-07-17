package subspace

import (
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
}
