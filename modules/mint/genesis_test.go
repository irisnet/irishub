package mint_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/simapp"
)

type TestSuite struct {
	suite.Suite

	cdc *codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.app = app
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestExportGenesis() {
	defaultGenesis := mint.DefaultGenesisState()
	exportedGenesis := mint.ExportGenesis(suite.ctx, suite.app.MintKeeper)
	suite.Equal(defaultGenesis, exportedGenesis)
}
