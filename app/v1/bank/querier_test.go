package bank

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/irisnet/irishub/app/protocol"
	"github.com/irisnet/irishub/app/v1/auth"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryAccount(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, QueryAccount),
		Data: []byte{},
	}

	res, err := queryAccount(input.ctx, req, input.bk, input.cdc)
	require.NotNil(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAccountParams([]byte("")))
	res, err = queryAccount(input.ctx, req, input.bk, input.cdc)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := keyPubAddr()
	req.Data = input.cdc.MustMarshalJSON(NewQueryAccountParams(addr))
	res, err = queryAccount(input.ctx, req, input.bk, input.cdc)
	require.NotNil(t, err)
	require.Nil(t, res)

	input.bk.am.SetAccount(input.ctx, input.bk.am.NewAccountWithAddress(input.ctx, addr))
	res, err = queryAccount(input.ctx, req, input.bk, input.cdc)
	require.Nil(t, err)
	require.NotNil(t, res)

	var account auth.Account
	err2 := input.cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
}

func TestQueryTokenStats(t *testing.T) {
	input := setupTestInput()

	// Test IrisCoinType Start ---------------
	params := QueryTokenStatsParams{
		TokenId: sdk.Iris,
	}

	bz, err := json.Marshal(params)
	require.Nil(t, err)

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, QueryTokenStats),
		Data: bz,
	}
	res, err := queryAccount(input.ctx, req, input.bk, input.cdc)
	require.NotNil(t, err)
	require.Nil(t, res)

	totalToken := sdk.Coins{sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(100))}
	input.bk.IncreaseLoosenToken(input.ctx, totalToken)

	_, _, addr := keyPubAddr()
	input.bk.am.SetAccount(input.ctx, input.bk.am.NewAccountWithAddress(input.ctx, addr))
	input.bk.AddCoins(input.ctx, addr, totalToken)

	burnedToken := sdk.Coins{sdk.NewCoin(sdk.IrisAtto, sdk.NewInt(50))}
	input.bk.BurnCoins(input.ctx, addr, burnedToken)

	res, err = queryTokenStats(input.ctx, req, input.bk, input.cdc)
	require.Nil(t, err)
	require.NotNil(t, res)

	var tokenStats TokenStats
	require.Nil(t, input.cdc.UnmarshalJSON(res, &tokenStats))
	require.Equal(t, totalToken.String(), (tokenStats.LooseTokens.Plus(burnedToken)).String())
	require.Equal(t, burnedToken.String(), tokenStats.BurnedTokens.String())

	// Test IrisCoinType End ---------------

	// Test !IrisCoinType Start ---------------
	params = QueryTokenStatsParams{
		TokenId: "abc",
	}

	denom := "abc-min"
	bz, err = json.Marshal(params)
	require.Nil(t, err)

	req = abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", protocol.AccountRoute, QueryTokenStats),
		Data: bz,
	}
	res, err = queryAccount(input.ctx, req, input.bk, input.cdc)
	require.NotNil(t, err)
	require.Nil(t, res)

	totalToken = sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(100))}
	input.bk.SetTotalSupply(input.ctx, sdk.Coin{Denom: denom, Amount: totalToken.AmountOf(denom)})

	_, _, addr = keyPubAddr()
	input.bk.am.SetAccount(input.ctx, input.bk.am.NewAccountWithAddress(input.ctx, addr))
	input.bk.AddCoins(input.ctx, addr, totalToken)

	burnedToken = sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(50))}
	input.bk.BurnCoins(input.ctx, addr, burnedToken)

	res, err = queryTokenStats(input.ctx, req, input.bk, input.cdc)
	require.Nil(t, err)
	require.NotNil(t, res)

	tokenStats = TokenStats{}
	require.Nil(t, input.cdc.UnmarshalJSON(res, &tokenStats))
	require.Equal(t, totalToken.String(), (tokenStats.TotalSupply.Plus(burnedToken)).String())
	require.Equal(t, burnedToken.String(), tokenStats.BurnedTokens.String())
}
