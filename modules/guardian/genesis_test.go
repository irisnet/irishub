package guardian_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/v4/modules/guardian"
	"github.com/irisnet/irishub/v4/modules/guardian/keeper"
	"github.com/irisnet/irishub/v4/modules/guardian/types"
	"github.com/irisnet/irishub/v4/testutil"
)

type TestSuite struct {
	suite.Suite

	cdc    codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *TestSuite) SetupTest() {
	app := testutil.CreateApp(suite.T())

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContextLegacy(false, tmproto.Header{})
	suite.keeper = app.GuardianKeeper
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestExportGenesis() {
	exportedGenesis := guardian.ExportGenesis(suite.ctx, suite.keeper)
	defaultGenesis := types.DefaultGenesisState()
	suite.Equal(exportedGenesis, defaultGenesis)
}
