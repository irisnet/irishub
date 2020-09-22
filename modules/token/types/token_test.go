package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestToken_ToMinCoin(t *testing.T) {
	token := Token{
		Symbol:        "iris",
		Name:          "irisnet",
		Scale:         18,
		MinUnit:       "atto",
		InitialSupply: 1000000,
		MaxSupply:     10000000,
		Mintable:      true,
		Owner:         nil,
	}

	amt, err := sdk.NewDecFromStr("1.500000000000000001")
	require.NoError(t, err)
	coin := sdk.NewDecCoinFromDec(token.Symbol, amt)

	c, err := token.ToMinCoin(coin)
	require.NoError(t, err)
	require.Equal(t, "1500000000000000001atto", c.String())

	coin1, err := token.ToMainCoin(c)
	require.NoError(t, err)
	require.Equal(t, coin, coin1)
}
