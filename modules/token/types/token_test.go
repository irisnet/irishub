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
		Owner:         "",
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

func TestCheckKeywords(t *testing.T) {
	type args struct {
		denom string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "right case", args: args{denom: "stake"}, wantErr: false},
		{name: "denom is peg", args: args{denom: "peg"}, wantErr: true},
		{name: "denom begin with peg", args: args{denom: "pegtoken"}, wantErr: true},
		{name: "denom is ibc", args: args{denom: "ibc"}, wantErr: true},
		{name: "denom begin with ibc", args: args{denom: "ibctoken"}, wantErr: true},
		{name: "denom is swap", args: args{denom: "swap"}, wantErr: true},
		{name: "denom begin with swap", args: args{denom: "swaptoken"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckKeywords(tt.args.denom); (err != nil) != tt.wantErr {
				t.Errorf("CheckKeywords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
