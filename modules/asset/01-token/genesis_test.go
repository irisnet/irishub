package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	token "github.com/irisnet/irishub/modules/asset/01-token"
	"github.com/irisnet/irishub/simapp"
)

type GenesisSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context

	keeper token.Keeper
}

func (suite *GenesisSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.keeper = app.AssetKeeper.TokenKeeper
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(GenesisSuite))
}

func (suite GenesisSuite) TestExportGenesis() {
	defaultGenesis := token.DefaultGenesisState()

	// add token
	addr := sdk.AccAddress([]byte("addr1"))
	ft := token.NewFungibleToken("Bitcoin", "btc", 18, "satoshi", sdk.NewInt(21000000), sdk.NewInt(21000000), false, addr)
	err := suite.keeper.AddToken(suite.ctx, ft)
	suite.NoError(err)

	// query all token
	tokens := suite.keeper.GetAllTokens(suite.ctx)

	suite.Equal(len(tokens), len(defaultGenesis.Tokens)+1)

	// export genesis
	genesisState := token.ExportGenesis(suite.ctx, suite.keeper)

	for _, t := range genesisState.Tokens {
		if t.Symbol == ft.Symbol {
			suite.Equal(ft, t)
		}
	}
}
