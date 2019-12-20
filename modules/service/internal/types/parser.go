package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/utils/protoidl"
)

const (
	outputPrivacy = "output_privacy"
	outputCached  = "output_cached"
	description   = "description"
)

// TODO
func ParseMethods(content string) (methods []string, err sdk.Error) {
	return
}

// TODO New MethodToMethodProperty process
//func MethodToMethodProperty(index int, method string) (methodProperty MethodProperty, err sdk.Error) {
//	return
//}

func MethodToMethodProperty(index int, method protoidl.Method) (methodProperty MethodProperty, err sdk.Error) {
	// set default value
	opp := NoPrivacy
	opc := NoCached
	var e error

	if _, ok := method.Attributes[outputPrivacy]; ok {
		if opp, e = OutputPrivacyEnumFromString(method.Attributes[outputPrivacy]); e != nil {
			err = ErrInvalidOutputPrivacyEnum(DefaultCodespace, method.Attributes[outputPrivacy])
			return
		}
	}

	if _, ok := method.Attributes[outputCached]; ok {
		if opc, e = OutputCachedEnumFromString(method.Attributes[outputCached]); e != nil {
			err = ErrInvalidOutputCachedEnum(DefaultCodespace, method.Attributes[outputCached])
			return
		}
	}

	methodProperty = MethodProperty{
		ID:            int16(index),
		Name:          method.Name,
		Description:   method.Attributes[description],
		OutputPrivacy: opp,
		OutputCached:  opc,
	}

	return
}
