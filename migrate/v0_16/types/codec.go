package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var CodeC = codec.New()

func init() {
	cryptocodec.RegisterCrypto(CodeC)
}
