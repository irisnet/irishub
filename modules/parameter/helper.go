package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesisParameter(p Parameter, ctx sdk.Context, genesisData interface{}) {

	find := p.LoadValue(ctx)
	if !find {
		p.InitGenesis(genesisData)
		p.SaveValue(ctx)
	}

}

