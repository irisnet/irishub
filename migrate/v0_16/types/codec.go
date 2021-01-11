package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var CodeC = codec.NewLegacyAmino()

func init() {
	cryptocodec.RegisterCrypto(CodeC)
}
