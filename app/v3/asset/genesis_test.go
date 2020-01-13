package asset

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/bank"
	"github.com/irisnet/irishub/app/v1/params"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

func TestExportGenesis(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	acc := ak.NewAccountWithAddress(ctx, addr)
	ft := NewFungibleToken("bch", "Bitcoin Network", "satoshi", 1, sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, acc.GetAddress())

	genesis := GenesisState{
		Params: DefaultParams(),
		Tokens: Tokens{ft},
	}

	//InitGenesis
	InitGenesis(ctx, keeper, genesis)

	// query all token
	var tokens Tokens
	keeper.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	//check param
	require.Equal(t, len(tokens), 1)

	// export gateways
	genesisState := ExportGenesis(ctx, keeper)

	require.Equal(t, DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}

}
