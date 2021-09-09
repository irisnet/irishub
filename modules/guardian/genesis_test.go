package guardian_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/modules/guardian/keeper"
	"github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/simapp"
)

type TestSuite struct {
	suite.Suite

	cdc    codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(false)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
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
