package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/tendermint/tendermint/crypto"
)

var _ sdk.Tx = (*StdTx)(nil)

type StdTx struct {
	Msgs       []sdk.Msg       `json:"msg"`
	Fee        legacytx.StdFee `json:"fee"`
	Signatures []StdSignature  `json:"signatures"`
	Memo       string          `json:"memo"`
}

// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsgs() []sdk.Msg { return tx.Msgs }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() error { return nil }

// Standard Signature
type StdSignature struct {
	crypto.PubKey `json:"pub_key"` // optional
	Signature     []byte           `json:"signature"`
	AccountNumber uint64           `json:"account_number"`
	Sequence      uint64           `json:"sequence"`
}
