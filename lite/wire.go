package lite

import (
	"github.com/irisnet/irishub/codec"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var cdc = codec.New()

func init() {
	ctypes.RegisterAmino(cdc)
}
