package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/irisnet/irishub/utils/protoidl"
)

const (
	outputPrivacy = "output_privacy"
	outputCached  = "output_cached"
	description   = "description"
)

// ParseMethods
// TODO
func ParseMethods(content string) (methods []string, err error) {
	return
}

// TODO New MethodToMethodProperty process
//func MethodToMethodProperty(index int, method string) (methodProperty MethodProperty, err error) {
//	return
//}

// MethodToMethodProperty
func MethodToMethodProperty(index int, method protoidl.Method) (methodProperty MethodProperty, err error) {
	// set default value
	opp := NoPrivacy
	opc := NoCached
	var e error

	if _, ok := method.Attributes[outputPrivacy]; ok {
		if opp, e = OutputPrivacyEnumFromString(method.Attributes[outputPrivacy]); e != nil {
			err = sdkerrors.Wrap(ErrInvalidOutputPrivacyEnum, method.Attributes[outputPrivacy])
			return
		}
	}

	if _, ok := method.Attributes[outputCached]; ok {
		if opc, e = OutputCachedEnumFromString(method.Attributes[outputCached]); e != nil {
			err = sdkerrors.Wrap(ErrInvalidOutputCachedEnum, method.Attributes[outputCached])
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
