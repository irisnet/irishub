package mint_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/mint"
	"github.com/irisnet/irishub/modules/mint/types"
	"github.com/irisnet/irishub/simapp"
)

type TestSuite struct {
	suite.Suite

	cdc codec.Codec
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.app = app
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestExportGenesis() {
	defaultGenesis := types.DefaultGenesisState()
	exportedGenesis := mint.ExportGenesis(suite.ctx, suite.app.MintKeeper)
	suite.Equal(defaultGenesis, exportedGenesis)
}
