package auth

import (
	"fmt"
	"testing"

	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_queryAccount(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", "acc", QueryAccount),
		Data: []byte{},
	}

	res, err := queryAccount(input.ctx, req, input.ak)
	require.NotNil(t, err)
	require.Nil(t, res)

	req.Data = input.cdc.MustMarshalJSON(NewQueryAccountParams([]byte("")))
	res, err = queryAccount(input.ctx, req, input.ak)
	require.NotNil(t, err)
	require.Nil(t, res)

	_, _, addr := keyPubAddr()
	req.Data = input.cdc.MustMarshalJSON(NewQueryAccountParams(addr))
	res, err = queryAccount(input.ctx, req, input.ak)
	require.NotNil(t, err)
	require.Nil(t, res)

	input.ak.SetAccount(input.ctx, input.ak.NewAccountWithAddress(input.ctx, addr))
	res, err = queryAccount(input.ctx, req, input.ak)
	require.Nil(t, err)
	require.NotNil(t, res)

	var account Account
	err2 := input.cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
}

func Test_queryTokenStats(t *testing.T) {
	input := setupTestInput()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", "acc", QueryTokenStats),
		Data: []byte{},
	}
	res, err := queryAccount(input.ctx, req, input.ak)
	require.NotNil(t, err)
	require.Nil(t, res)

	loosenToken := sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(100))}
	input.ak.IncreaseTotalLoosenToken(input.ctx, loosenToken)

	burnedToken := sdk.Coins{sdk.NewCoin("iris", sdk.NewInt(50))}
	input.ak.IncreaseBurnedToken(input.ctx, burnedToken)

	res, err = queryTokenStats(input.ctx, input.ak)
	require.Nil(t, err)
	require.NotNil(t, res)

	var tokenStats TokenStats
	require.Nil(t, input.cdc.UnmarshalJSON(res, &tokenStats))
	require.Equal(t, loosenToken.String(), tokenStats.LooseTokens.String())
	require.Equal(t, burnedToken.String(), tokenStats.BurnedTokens.String())
}
