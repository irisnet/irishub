package auth

import (
	"github.com/irisnet/irishub/modules/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
)

var _ sdk.Tx = (*StdTx)(nil)

type StdTx = auth.StdTx
type StdFee = auth.StdFee
type StdSignature = auth.StdSignature

var (
	StdSignBytes     = auth.StdSignBytes
	DefaultTxDecoder = auth.DefaultTxDecoder
	NewStdFee        = auth.NewStdFee
	NewStdTx         = auth.NewStdTx
)

func countSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(*multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}
	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += countSubKeys(subkey)
	}
	return numKeys
}
