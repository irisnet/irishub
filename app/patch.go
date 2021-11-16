package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// excute is responsible for checking the permission to execute MsgCallService
type Execute = func(ctx sdk.Context) error

// BlockTimerExecutor is responsible for checking the permission to execute MsgCallService
type BlockTimerExecutor struct {
	patchs map[int64]Execute
}

func NewBlockTimerExecutor() *BlockTimerExecutor {
	return &BlockTimerExecutor{
		patchs: make(map[int64]Execute),
	}
}

func (bte *BlockTimerExecutor) add(height int64, patch Execute) {
	bte.patchs[height] = patch
}

func (bte *BlockTimerExecutor) Start(ctx sdk.Context) {
	exec := bte.patchs[ctx.BlockHeight()]
	if exec != nil {
		ctx.Logger().Info("Execute patch", "height", ctx.BlockHeight())
		if err := exec(ctx); err != nil {
			panic(err)
		}
	}
}
