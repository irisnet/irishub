package auth

import (
	"errors"
	"fmt"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
	"github.com/tendermint/tendermint/crypto"
)

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type Account interface {
	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error

	GetFrozenCoins() sdk.Coins
	GetFrozenCoinByDenom(denom string) sdk.Coin
	SetFrozenCoin(sdk.Coin) error
	DeductFrozenCoin(sdk.Coin) error
}

// AccountDecoder unmarshals account bytes
type AccountDecoder func(accountBytes []byte) (Account, error)

//-----------------------------------------------------------
// BaseAccount

var _ Account = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// There are examples of this in: examples/basecoin/types/account.go.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address        sdk.AccAddress `json:"address"`
	Coins          sdk.Coins      `json:"coins"`
	PubKey         crypto.PubKey  `json:"public_key"`
	AccountNumber  uint64         `json:"account_number"`
	Sequence       uint64         `json:"sequence"`
	FrozenCoins    sdk.Coins      `json:"frozen_coins"`
}

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
  Address:         %s
  Pubkey:          %s
  Coins:           %s
  Frozen coins:    %s
  Account Number:  %d
  Sequence:        %d`,
		acc.Address, pubkey, acc.Coins.MainUnitString(),acc.FrozenCoins.MainUnitString(), acc.AccountNumber, acc.Sequence,
	)
}

// Prototype function for BaseAccount
func ProtoBaseAccount() Account {
	return &BaseAccount{}
}

func NewBaseAccountWithAddress(addr sdk.AccAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetFrozenCoinByDenom(denom string) sdk.Coin {
	for _, coin := range acc.FrozenCoins {
		if coin.Denom == denom {
			return coin
		}
	}
	return sdk.Coin{}
}

// Implements sdk.Account.
func (acc *BaseAccount) GetFrozenCoins() sdk.Coins {
	return acc.FrozenCoins
}

// Implements sdk.Account. plus frozen token to account.FrozenToken.
func (acc *BaseAccount) SetFrozenCoin(coin sdk.Coin) error {
	acc.FrozenCoins = acc.FrozenCoins.Plus(sdk.Coins{coin})
	return nil
}

// Implements sdk.Account.deduct frozen token from frozen coin
func (acc *BaseAccount) DeductFrozenCoin(coin sdk.Coin) error {
	if acc.FrozenCoins == nil {
		return errors.New("this account has no frozen coin")
	}
	oldCoins := acc.FrozenCoins
	diff, hasNeg := oldCoins.SafeMinus(sdk.Coins{coin})
	if hasNeg {
		return errors.New("this account has not enough frozen coin")
	}
	acc.FrozenCoins = diff
	return nil
}

// Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

//----------------------------------------
// Wire

// Most users shouldn't use this, but this comes in handy for tests.
func RegisterBaseAccount(cdc *codec.Codec) {
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "irishub/bank/BaseAccount", nil)
	codec.RegisterCrypto(cdc)
}
