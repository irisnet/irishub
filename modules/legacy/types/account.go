package types

import (
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// There are examples of this in: examples/basecoin/types/account.go.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       sdk.AccAddress     `json:"address"`
	Coins         sdk.Coins          `json:"coins"`
	PubKey        cryptotypes.PubKey `json:"public_key"`
	AccountNumber uint64             `json:"account_number"`
	Sequence      uint64             `json:"sequence"`
	MemoRegexp    string             `json:"memo_regexp"`
}
