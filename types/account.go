package types

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// There are examples of this in: examples/basecoin/types/account.go.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       AccAddress    `json:"address"`
	Coins         Coins         `json:"coins"`
	PubKey        crypto.PubKey `json:"public_key"`
	AccountNumber uint64        `json:"account_number"`
	Sequence      uint64        `json:"sequence"`
	MemoRegexp    string        `json:"memo_regexp"`
}

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
  Address:         %s
  Pubkey:          %s
  Coins:           %s
  Account Number:  %d
  Sequence:        %d
  Memo Regexp:     %s`,
		acc.Address,
		pubkey,
		acc.Coins.String(),
		acc.AccountNumber,
		acc.Sequence,
		acc.MemoRegexp,
	)
}

// String implements human.Stringer
func (acc BaseAccount) HumanString(converter CoinsConverter) string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
  Address:         %s
  Pubkey:          %s
  Coins:           %s
  Account Number:  %d
  Sequence:        %d 
  Memo Regexp:     %s`,
		acc.Address,
		pubkey,
		converter.ToMainUnit(acc.Coins),
		acc.AccountNumber,
		acc.Sequence,
		acc.MemoRegexp,
	)
}

// Implements Account.
func (acc *BaseAccount) GetAddress() AccAddress {
	return acc.Address
}

// Implements Account.
func (acc *BaseAccount) SetAddress(addr AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// Implements Account.
func (acc *BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// Implements Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// Implements Account.
func (acc *BaseAccount) GetCoins() Coins {
	return acc.Coins
}

// Implements Account.
func (acc *BaseAccount) SetCoins(coins Coins) error {
	acc.Coins = coins
	return nil
}

// Implements Account.
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// Implements Account.
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// Implements Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// Implements Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// Implements Account.
func (acc *BaseAccount) GetMemoRegexp() string {
	return acc.MemoRegexp
}

// Implements Account.
func (acc *BaseAccount) SetMemoRegexp(regexp string) error {
	acc.MemoRegexp = regexp
	return nil
}
