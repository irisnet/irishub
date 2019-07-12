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

func TestExportGatewayGenesis(t *testing.T) {
	ms, accountKey, assetKey, paramskey, paramsTkey := tests.SetupMultiStore()

	cdc := codec.New()
	RegisterCodec(cdc)
	auth.RegisterBaseAccount(cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	paramsKeeper := params.NewKeeper(cdc, paramskey, paramsTkey)
	ak := auth.NewAccountKeeper(cdc, accountKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(cdc, ak)
	keeper := NewKeeper(cdc, assetKey, bk, DefaultCodespace, paramsKeeper.Subspace(DefaultParamSpace))

	// init params
	keeper.SetParamSet(ctx, DefaultParams())

	// define variables
	owners := []sdk.AccAddress{
		ak.NewAccountWithAddress(ctx, []byte("owner1")).GetAddress(),
		ak.NewAccountWithAddress(ctx, []byte("owner2")).GetAddress(),
		ak.NewAccountWithAddress(ctx, []byte("owner3")).GetAddress(),
		ak.NewAccountWithAddress(ctx, []byte("owner4")).GetAddress(),
	}
	monikers := []string{"firstmk", "secondmk", "thirdmk", "fouthmk"}
	identities := []string{"identity1", "identity2", "identity3", "identity4"}
	details := []string{"details1", "details2", "details3", "details4"}
	websites := []string{"website1", "website2", "website3", "website4"}

	var gateways []Gateway

	// construct and store gateways with data above
	for i := 0; i < 4; i++ {
		gateway := NewGateway(owners[i], monikers[i], identities[i], details[i], websites[i])
		gateways = append(gateways, gateway)

		keeper.SetGateway(ctx, gateway)
		keeper.SetOwnerGateway(ctx, owners[i], monikers[i])
	}

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	acc := ak.NewAccountWithAddress(ctx, addr)
	ft := NewFungibleToken(NATIVE, "", "bch", "bch", 1, "", "satoshi", sdk.NewIntWithDecimal(1, 0), sdk.NewIntWithDecimal(1, 0), true, acc.GetAddress())
	keeper.AddToken(ctx, ft)

	// query all gateways
	var storedGateways []Gateway
	keeper.IterateGateways(ctx, func(g Gateway) bool {
		storedGateways = append(storedGateways, g)
		return false
	})

	// query all token
	var tokens Tokens
	keeper.IterateTokens(ctx, func(token FungibleToken) (stop bool) {
		tokens = append(tokens, token)
		return false
	})

	require.Equal(t, len(gateways), len(storedGateways))
	require.Equal(t, len(tokens), 1)

	// export gateways
	genesisState := ExportGenesis(ctx, keeper)
	exportedGateways := genesisState.Gateways

	// assert that exported gateways are consistant with the stored gateways
	require.Equal(t, len(storedGateways), len(exportedGateways))
	for i, eg := range exportedGateways {
		require.Equal(t, storedGateways[i], eg)
	}

	for _, token := range genesisState.Tokens {
		require.Equal(t, token, ft)
	}
}
