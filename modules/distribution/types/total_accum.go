package types

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/libs/log"
)

// total accumulation tracker
type TotalAccum struct {
	UpdateHeight int64   `json:"update_height"`
	Accum        sdk.Dec `json:"accum"`
}

func NewTotalAccum(height int64) TotalAccum {
	return TotalAccum{
		UpdateHeight: height,
		Accum:        sdk.ZeroDec(),
	}
}

// update total accumulation factor for the new height
// CONTRACT: height should be greater than the old height
func (ta TotalAccum) UpdateForNewHeight(logger log.Logger, height int64, accumCreatedPerBlock sdk.Dec) TotalAccum {
	blocks := height - ta.UpdateHeight
	if blocks < 0 {
		panic("reverse updated for new height")
	}
	accumDelta := accumCreatedPerBlock.MulInt(sdk.NewInt(blocks))
	ta.Accum = ta.Accum.Add(accumDelta)
	logger.Debug("Update accumulation", "height", height, "last_update_height", ta.UpdateHeight, "accum_delta", accumDelta.String())
	ta.UpdateHeight = height
	return ta
}

// get total accumulation factor for the given height
// CONTRACT: height should be greater than the old height
func (ta TotalAccum) GetAccum(height int64, accumCreatedPerBlock sdk.Dec) sdk.Dec {
	blocks := height - ta.UpdateHeight
	if blocks < 0 {
		panic("reverse updated for new height")
	}
	return ta.Accum.Add(accumCreatedPerBlock.MulInt(sdk.NewInt(blocks)))
}
