package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Execute is responsible for checking the permission to execute patches
type Execute = func(ctx sdk.Context) error

// BlockTimerExecutor is responsible for checking the permission to execute patches
type BlockTimerExecutor struct {
	patches map[int64]Execute
}

func NewBlockTimerExecutor() *BlockTimerExecutor {
	return &BlockTimerExecutor{
		patches: make(map[int64]Execute),
	}
}

func (bte *BlockTimerExecutor) add(height int64, patch Execute) {
	bte.patches[height] = patch
}

func (bte *BlockTimerExecutor) Start(ctx sdk.Context) {
	exec := bte.patches[ctx.BlockHeight()]
	if exec != nil {
		ctx.Logger().Info("Execute patch", "height", ctx.BlockHeight())
		if err := exec(ctx); err != nil {
			panic(err)
		}
	}
}
