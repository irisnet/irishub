package parameter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

func InitGenesisParameter(p Parameter, ctx sdk.Context, genesisData interface{}) {
	if p != nil {
		find := p.LoadValue(ctx)
		if !find {
			p.InitGenesis(genesisData)
			p.SaveValue(ctx)
		}
	}
}

func SetParamReadWriter(setter params.Setter, ps ...Parameter) {
	for _, p := range ps {
		if p != nil {
			p.SetReadWriter(setter)
		}
	}
}
