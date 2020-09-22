package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens []Token) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

//SetNativeToken reset the system's default native token
func SetNativeToken(symbol,
	name,
	minUnit string,
	decimal uint32,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress) {
	nativeToken = NewToken(symbol, name, minUnit, decimal, initialSupply, maxSupply, mintable, owner)
}

func GetNativeToken() Token {
	return nativeToken
}

var nativeToken = Token{
	Symbol:        sdk.DefaultBondDenom,
	Name:          "Network staking token",
	Scale:         0,
	MinUnit:       sdk.DefaultBondDenom,
	InitialSupply: 2000000000,
	MaxSupply:     10000000000,
	Mintable:      true,
	Owner:         sdk.AccAddress(crypto.AddressHash([]byte(ModuleName))),
}
