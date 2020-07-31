package types

import (
	"github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
)

var CodeC = amino.NewCodec()

func init() {
	tmtypes.RegisterBlockAmino(CodeC)
}
