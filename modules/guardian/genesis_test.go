package guardian_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/simapp"
)

type TestSuite struct {
	suite.Suite

	cdc    *codec.Codec
	ctx    sdk.Context
	keeper guardian.Keeper
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = app.Codec()
	suite.ctx = app.BaseApp.NewContext(false, abci.Header{})
	suite.keeper = app.GuardianKeeper
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestExportGenesis() {
	exportedGenesis := guardian.ExportGenesis(suite.ctx, suite.keeper)
	defaultGenesis := guardian.DefaultGenesisState()
	suite.Equal(exportedGenesis, defaultGenesis)
}
