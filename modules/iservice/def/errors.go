package def

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	CodeNotEmpty   sdk.CodeType = 10000
	CodeHasExisted sdk.CodeType = 10001
)

func NotEmpty(msg string) sdk.Error {
	return sdk.NewError(CodeNotEmpty, msg)
}

func HasExisted(msg string) sdk.Error {
	return sdk.NewError(CodeHasExisted, msg)
}
