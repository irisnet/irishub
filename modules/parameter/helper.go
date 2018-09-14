package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesisParameter(p Parameter, ctx sdk.Context) {

	find := p.LoadValue(ctx)
	if !find {
		p.InitGenesis(nil)
		p.SaveValue(ctx)
	}

}

