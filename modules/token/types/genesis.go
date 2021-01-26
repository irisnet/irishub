package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var nativeToken = Token{
	Symbol:        sdk.DefaultBondDenom,
	Name:          "Network staking token",
	Scale:         0,
	MinUnit:       sdk.DefaultBondDenom,
	InitialSupply: 2000000000,
	MaxSupply:     10000000000,
	Mintable:      true,
	Owner:         sdk.AccAddress(crypto.AddressHash([]byte(ModuleName))).String(),
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens []Token) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

//SetNativeToken reset the system's default native token
func SetNativeToken(
	symbol string,
	name string,
	minUnit string,
	decimal uint32,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) {
	nativeToken = NewToken(symbol, name, minUnit, decimal, initialSupply, maxSupply, mintable, owner)
}

//GetNativeToken return the system's default native token
func GetNativeToken() Token {
	return nativeToken
}

// ValidateGenesis validates the provided token genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	// validate token
	for _, token := range data.Tokens {
		if err := ValidateToken(token); err != nil {
			return err
		}
	}

	// validate token
	for _, coin := range data.BurnedCoins {
		if err := coin.Validate(); err != nil {
			return err
		}
	}
	return nil
}
