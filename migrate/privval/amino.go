package privval

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cryptoAmino "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var cdc = codec.NewLegacyAmino()

func init() {
	cryptoAmino.RegisterCrypto(cdc)
}
