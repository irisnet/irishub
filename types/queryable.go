package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ParseParamsErr(err error) sdk.Error {
	return sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
}

func MarshalResultErr(err error) sdk.Error {
	return sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
}
