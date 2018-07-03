package app

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Context struct {
	ctx context.CoreContext
}

func (c Context) BroadcastTxAsync(tx []byte) (*ctypes.ResultBroadcastTx, error) {
	return c.ctx.Client.BroadcastTxAsync(tx)
}

func (c Context) GetCosmosCtx() (context.CoreContext) {
	return c.ctx
}

func NewContext() Context {
	return Context{
		context.NewCoreContextFromViper(),
	}
}
