package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ModuleService struct {
	ServiceName     string
	Provider        sdk.AccAddress
	ReuquestService func(ctx sdk.Context, input string) (result string, output string)
}
